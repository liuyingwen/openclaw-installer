package repair

import (
	"testing"

	"github.com/liuyingwen/openclaw-installer/internal/config"
	"github.com/liuyingwen/openclaw-installer/internal/platform"
)

func TestPlannerCreatesInstallCommandsForMissingTools(t *testing.T) {
	planner := Planner{}

	plan, err := planner.Plan(platform.Profile{
		OS:              platform.OSMac,
		PackageManagers: []string{"brew"},
	}, []config.PrerequisiteSpec{
		{
			Name: "git",
			Packages: map[string]string{
				"brew": "git",
			},
		},
	}, map[string]bool{"git": false})
	if err != nil {
		t.Fatalf("expected plan without error, got %v", err)
	}

	if len(plan.Commands) != 1 {
		t.Fatalf("expected 1 repair command, got %d", len(plan.Commands))
	}

	if plan.Commands[0] != "brew install git" {
		t.Fatalf("expected brew install command, got %q", plan.Commands[0])
	}
}

func TestPlannerSkipsInstalledTools(t *testing.T) {
	planner := Planner{}

	plan, err := planner.Plan(platform.Profile{
		OS:              platform.OSMac,
		PackageManagers: []string{"brew"},
	}, []config.PrerequisiteSpec{
		{
			Name: "git",
			Packages: map[string]string{
				"brew": "git",
			},
		},
	}, map[string]bool{"git": true})
	if err != nil {
		t.Fatalf("expected plan without error, got %v", err)
	}

	if len(plan.Commands) != 0 {
		t.Fatalf("expected 0 repair commands, got %d", len(plan.Commands))
	}
}

func TestPlannerFailsWithoutSupportedPackageManager(t *testing.T) {
	planner := Planner{}

	_, err := planner.Plan(platform.Profile{
		OS:              platform.OSWindows,
		PackageManagers: nil,
	}, []config.PrerequisiteSpec{
		{
			Name: "git",
			Packages: map[string]string{
				"winget": "Git.Git",
			},
		},
	}, map[string]bool{"git": false})
	if err == nil {
		t.Fatal("expected unsupported package manager error")
	}
}

func TestPlannerBootstrapsHomebrewOnMacWhenMissing(t *testing.T) {
	planner := Planner{}

	plan, err := planner.Plan(platform.Profile{
		OS:              platform.OSMac,
		PackageManagers: nil,
	}, []config.PrerequisiteSpec{
		{
			Name: "git",
			Packages: map[string]string{
				"brew": "git",
			},
		},
	}, map[string]bool{"git": false})
	if err != nil {
		t.Fatalf("expected plan without error, got %v", err)
	}

	if len(plan.Commands) != 2 {
		t.Fatalf("expected 2 commands, got %d", len(plan.Commands))
	}

	if plan.Commands[0] != "/bin/bash -c \"$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)\"" {
		t.Fatalf("unexpected bootstrap command: %q", plan.Commands[0])
	}
	if plan.Commands[1] != "brew install git" {
		t.Fatalf("unexpected install command: %q", plan.Commands[1])
	}
}

func TestPlannerBootstrapsScoopOnWindowsWhenMissing(t *testing.T) {
	planner := Planner{}

	plan, err := planner.Plan(platform.Profile{
		OS:              platform.OSWindows,
		PackageManagers: nil,
	}, []config.PrerequisiteSpec{
		{
			Name: "git",
			Packages: map[string]string{
				"scoop": "git",
			},
		},
	}, map[string]bool{"git": false})
	if err != nil {
		t.Fatalf("expected plan without error, got %v", err)
	}

	if len(plan.Commands) != 2 {
		t.Fatalf("expected 2 commands, got %d", len(plan.Commands))
	}

	if plan.Commands[0] != "powershell -ExecutionPolicy Bypass -Command \"iwr -useb get.scoop.sh | iex\"" {
		t.Fatalf("unexpected bootstrap command: %q", plan.Commands[0])
	}
	if plan.Commands[1] != "scoop install git" {
		t.Fatalf("unexpected install command: %q", plan.Commands[1])
	}
}
