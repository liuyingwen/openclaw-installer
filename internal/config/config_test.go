package config

import (
	"strings"
	"testing"
)

func TestLoadParsesInstallerManifest(t *testing.T) {
	cfg, err := Load("../../config/openclaw.example.yaml")
	if err != nil {
		t.Fatalf("expected config to load, got error: %v", err)
	}

	if cfg.App.Name != "openclaw" {
		t.Fatalf("expected app name %q, got %q", "openclaw", cfg.App.Name)
	}

	if len(cfg.Prerequisites) != 0 {
		t.Fatalf("expected 0 prerequisites, got %d", len(cfg.Prerequisites))
	}

	if len(cfg.Install) != 2 {
		t.Fatalf("expected 2 install steps, got %d", len(cfg.Install))
	}

	if len(cfg.Verify) != 1 {
		t.Fatalf("expected 1 verify step, got %d", len(cfg.Verify))
	}

	if cfg.Install[0].CommandFor("darwin") == "" {
		t.Fatal("expected darwin install command")
	}

	if cfg.Install[0].CommandFor("windows") == "" {
		t.Fatal("expected windows install command")
	}
}

func TestLoadRejectsInvalidManifest(t *testing.T) {
	_, err := Load("../../config/does-not-exist.yaml")
	if err == nil {
		t.Fatal("expected load to fail for missing file")
	}
}

func TestDefaultUsesRecommendedInstallScriptAndGlobalOpenClawCommands(t *testing.T) {
	cfg := Default()

	if !strings.Contains(cfg.Install[0].CommandFor("darwin"), "https://openclaw.ai/install.sh") {
		t.Fatalf("expected darwin install to use install.sh, got %q", cfg.Install[0].CommandFor("darwin"))
	}
	if !strings.Contains(cfg.Install[1].CommandFor("darwin"), "openclaw gateway install") {
		t.Fatalf("expected darwin gateway install to use global openclaw command, got %q", cfg.Install[1].CommandFor("darwin"))
	}
	if !strings.Contains(cfg.Verify[0].CommandFor("darwin"), "openclaw --version") {
		t.Fatalf("expected darwin verify to use global openclaw command, got %q", cfg.Verify[0].CommandFor("darwin"))
	}
	if !strings.Contains(cfg.Install[1].CommandFor("windows"), "$cmd.Definition") {
		t.Fatalf("expected windows gateway install to invoke resolved command definition, got %q", cfg.Install[1].CommandFor("windows"))
	}
}

func TestResolveCommandFailsForUnsupportedOS(t *testing.T) {
	step := CommandStep{
		Name: "install openclaw",
		RunByOS: map[string]string{
			"darwin": "run-on-mac",
		},
	}

	_, err := step.ResolveCommand("windows")
	if err == nil {
		t.Fatal("expected unsupported os error")
	}
}
