package cli

import (
	"bytes"
	"strings"
	"testing"

	"github.com/liuyingwen/openclaw-installer/internal/config"
	"github.com/liuyingwen/openclaw-installer/internal/installer"
	"github.com/liuyingwen/openclaw-installer/internal/platform"
	"github.com/liuyingwen/openclaw-installer/internal/repair"
)

type stubDependencies struct {
	config      config.Config
	profile     platform.Profile
	installed   map[string]bool
	workflowErr error
	ranInstall  bool
	lastMode    string
}

func (d *stubDependencies) LoadConfig(string) (config.Config, error) {
	return d.config, nil
}

func (d *stubDependencies) DetectProfile() (platform.Profile, error) {
	return d.profile, nil
}

func (d *stubDependencies) DetectInstalled(prerequisites []config.PrerequisiteSpec) (map[string]bool, error) {
	result := map[string]bool{}
	for _, prerequisite := range prerequisites {
		result[prerequisite.Name] = d.installed[prerequisite.Name]
	}
	return result, nil
}

func (d *stubDependencies) RunWorkflow(mode string, cfg config.Config, plan repair.Plan) (installer.Result, error) {
	d.lastMode = mode
	if mode == "install" {
		d.ranInstall = true
	}
	if d.workflowErr != nil {
		return installer.Result{}, d.workflowErr
	}
	return installer.Result{
		PlannedCommands: append(append([]string{}, plan.Commands...), commandRunsFor(string(d.profile.OS), cfg.Install, cfg.Verify)...),
		LogPath:         "/tmp/openclaw-installer.log",
	}, nil
}

func TestInstallDryRunPrintsPlanWithoutExecutingWorkflow(t *testing.T) {
	deps := &stubDependencies{
		config: sampleConfig(),
		profile: platform.Profile{
			OS:              platform.OSMac,
			PackageManagers: []string{"brew"},
		},
		installed: map[string]bool{"git": false},
	}

	app := NewApp(deps)
	var output bytes.Buffer

	if err := app.Run([]string{"install", "--config", "config/openclaw.example.yaml", "--dry-run"}, &output); err != nil {
		t.Fatalf("expected dry run to succeed, got %v", err)
	}

	if deps.ranInstall {
		t.Fatal("expected dry run to avoid install execution")
	}

	text := output.String()
	if !strings.Contains(text, "https://openclaw.ai/install.sh") {
		t.Fatalf("expected install step in output, got %q", text)
	}
	if !strings.Contains(text, "openclaw gateway install") {
		t.Fatalf("expected gateway install to use global openclaw command, got %q", text)
	}
}

func TestDoctorPrintsEnvironmentSummary(t *testing.T) {
	deps := &stubDependencies{
		config: sampleConfig(),
		profile: platform.Profile{
			OS:              platform.OSLinux,
			PackageManagers: []string{"apt"},
		},
		installed: map[string]bool{"git": false},
	}

	app := NewApp(deps)
	var output bytes.Buffer

	if err := app.Run([]string{"doctor", "--config", "config/openclaw.example.yaml"}, &output); err != nil {
		t.Fatalf("expected doctor to succeed, got %v", err)
	}

	text := output.String()
	if !strings.Contains(text, "platform: linux") {
		t.Fatalf("expected platform summary, got %q", text)
	}
	if !strings.Contains(text, "repair plan:") {
		t.Fatalf("expected repair plan heading, got %q", text)
	}
}

func TestPrintPlanOnlyPrintsPlan(t *testing.T) {
	deps := &stubDependencies{
		config: sampleConfig(),
		profile: platform.Profile{
			OS:              platform.OSWindows,
			PackageManagers: []string{"winget"},
		},
		installed: map[string]bool{"git": false},
	}

	app := NewApp(deps)
	var output bytes.Buffer

	if err := app.Run([]string{"print-plan", "--config", "config/openclaw.example.yaml"}, &output); err != nil {
		t.Fatalf("expected print-plan to succeed, got %v", err)
	}

	text := output.String()
	if !strings.Contains(text, "planned commands:") {
		t.Fatalf("expected planned commands heading, got %q", text)
	}
	if !strings.Contains(text, "powershell -ExecutionPolicy Bypass -Command") {
		t.Fatalf("expected windows install command, got %q", text)
	}
}

func TestPrintPlanUsesEmbeddedDefaultConfigWhenConfigFlagMissing(t *testing.T) {
	deps := &stubDependencies{
		config: sampleConfig(),
		profile: platform.Profile{
			OS:              platform.OSMac,
			PackageManagers: []string{"brew"},
		},
		installed: map[string]bool{},
	}

	app := NewApp(deps)
	var output bytes.Buffer

	if err := app.Run([]string{"print-plan"}, &output); err != nil {
		t.Fatalf("expected print-plan without config to succeed, got %v", err)
	}

	text := output.String()
	if !strings.Contains(text, "https://openclaw.ai/install.sh") {
		t.Fatalf("expected embedded install command, got %q", text)
	}
	if !strings.Contains(text, "openclaw --version") {
		t.Fatalf("expected embedded verify command to use global openclaw command, got %q", text)
	}
}

func TestInstallRequiresConfirmUnlessNonInteractive(t *testing.T) {
	deps := &stubDependencies{
		config: sampleConfig(),
		profile: platform.Profile{
			OS:              platform.OSMac,
			PackageManagers: []string{"brew"},
		},
	}

	app := NewApp(deps)
	var output bytes.Buffer

	err := app.Run([]string{"install", "--config", "config/openclaw.example.yaml"}, &output)
	if err == nil {
		t.Fatal("expected install to require confirmation")
	}

	if deps.ranInstall {
		t.Fatal("expected install not to run without confirmation")
	}
}

func TestInstallNonInteractiveRunsAndReportsLogPath(t *testing.T) {
	deps := &stubDependencies{
		config: sampleConfig(),
		profile: platform.Profile{
			OS:              platform.OSMac,
			PackageManagers: []string{"brew"},
		},
	}

	app := NewApp(deps)
	var output bytes.Buffer

	err := app.Run([]string{"install", "--config", "config/openclaw.example.yaml", "--yes"}, &output)
	if err != nil {
		t.Fatalf("expected install to succeed, got %v", err)
	}

	if !deps.ranInstall || deps.lastMode != "install" {
		t.Fatal("expected install workflow to run")
	}

	text := output.String()
	if !strings.Contains(text, "log file: /tmp/openclaw-installer.log") {
		t.Fatalf("expected log path in output, got %q", text)
	}
	if !strings.Contains(text, "openclaw should now be available on your PATH") {
		t.Fatalf("expected global command guidance in output, got %q", text)
	}
	if !strings.Contains(text, "try: openclaw --version") {
		t.Fatalf("expected verification guidance in output, got %q", text)
	}
}

func sampleConfig() config.Config {
	return config.Default()
}

func commandRunsFor(goos string, groups ...[]config.CommandStep) []string {
	var result []string
	for _, steps := range groups {
		for _, step := range steps {
			result = append(result, step.CommandFor(goos))
		}
	}
	return result
}
