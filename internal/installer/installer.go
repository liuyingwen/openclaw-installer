package installer

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/liuyingwen/openclaw-installer/internal/config"
	"github.com/liuyingwen/openclaw-installer/internal/repair"
)

type Runner interface {
	Run(command string) error
}

type Workflow struct {
	Runner Runner
	GOOS   string
	LogPath string
}

type Result struct {
	PlannedCommands []string
	LogPath         string
}

func (w Workflow) Doctor(cfg config.Config, repairPlan repair.Plan) (Result, error) {
	planned, err := appendCommands(w.GOOS, repairPlan.Commands, cfg.Install, cfg.Verify)
	if err != nil {
		return Result{}, err
	}
	return Result{
		PlannedCommands: planned,
		LogPath:         w.LogPath,
	}, nil
}

func (w Workflow) Install(cfg config.Config, repairPlan repair.Plan) (Result, error) {
	planned, err := appendCommands(w.GOOS, repairPlan.Commands, cfg.Install, cfg.Verify)
	if err != nil {
		return Result{}, err
	}
	logPath, err := w.ensureLogPath()
	if err != nil {
		return Result{}, err
	}
	if err := writeLog(logPath, planned); err != nil {
		return Result{}, err
	}
	for _, command := range planned {
		if err := w.Runner.Run(command); err != nil {
			return Result{}, fmt.Errorf("run %q: %w", command, err)
		}
	}

	return Result{PlannedCommands: planned, LogPath: logPath}, nil
}

func appendCommands(goos string, repairCommands []string, installSteps []config.CommandStep, verifySteps []config.CommandStep) ([]string, error) {
	planned := append([]string(nil), repairCommands...)
	for _, step := range installSteps {
		command, err := step.ResolveCommand(goos)
		if err != nil {
			return nil, err
		}
		planned = append(planned, command)
	}
	for _, step := range verifySteps {
		command, err := step.ResolveCommand(goos)
		if err != nil {
			return nil, err
		}
		planned = append(planned, command)
	}
	return planned, nil
}

func (w Workflow) ensureLogPath() (string, error) {
	if w.LogPath != "" {
		return w.LogPath, nil
	}

	dir, err := os.MkdirTemp("", "openclaw-installer-*")
	if err != nil {
		return "", fmt.Errorf("create log dir: %w", err)
	}

	return filepath.Join(dir, "install.log"), nil
}

func writeLog(path string, commands []string) error {
	content := strings.Join(append([]string{time.Now().UTC().Format(time.RFC3339)}, commands...), "\n") + "\n"
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		return fmt.Errorf("write log: %w", err)
	}
	return nil
}
