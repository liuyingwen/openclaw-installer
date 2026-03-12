package installer

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/liuyingwen/openclaw-installer/internal/config"
	"github.com/liuyingwen/openclaw-installer/internal/repair"
)

type logRunner struct{}

func (logRunner) Run(string) error { return nil }

func TestInstallWritesExecutionLog(t *testing.T) {
	tempDir := t.TempDir()
	logPath := filepath.Join(tempDir, "install.log")

	workflow := Workflow{
		Runner:  logRunner{},
		GOOS:    "darwin",
		LogPath: logPath,
	}

	cfg := config.Config{
		Install: []config.CommandStep{{Name: "install", Run: "install-command"}},
		Verify:  []config.CommandStep{{Name: "verify", Run: "verify-command"}},
	}

	result, err := workflow.Install(cfg, repair.Plan{Commands: []string{"repair-command"}})
	if err != nil {
		t.Fatalf("expected install to succeed, got %v", err)
	}

	if result.LogPath != logPath {
		t.Fatalf("expected log path %q, got %q", logPath, result.LogPath)
	}

	content, err := os.ReadFile(logPath)
	if err != nil {
		t.Fatalf("expected log file to exist, got %v", err)
	}

	text := string(content)
	for _, expected := range []string{"repair-command", "install-command", "verify-command"} {
		if !strings.Contains(text, expected) {
			t.Fatalf("expected log to contain %q, got %q", expected, text)
		}
	}
}
