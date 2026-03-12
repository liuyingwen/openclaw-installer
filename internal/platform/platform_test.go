package platform

import "testing"

func TestDetectsDarwinWithBrew(t *testing.T) {
	profile := Detect(Inputs{
		GOOS:           "darwin",
		AvailableTools: map[string]bool{"brew": true},
	})

	if profile.OS != OSMac {
		t.Fatalf("expected %q, got %q", OSMac, profile.OS)
	}

	if len(profile.PackageManagers) != 1 || profile.PackageManagers[0] != "brew" {
		t.Fatalf("expected brew package manager, got %#v", profile.PackageManagers)
	}
}

func TestDetectsLinuxPackageManagerPriority(t *testing.T) {
	profile := Detect(Inputs{
		GOOS:           "linux",
		IDLike:         []string{"debian"},
		AvailableTools: map[string]bool{"apt": true, "dnf": true},
	})

	if profile.OS != OSLinux {
		t.Fatalf("expected %q, got %q", OSLinux, profile.OS)
	}

	if len(profile.PackageManagers) == 0 || profile.PackageManagers[0] != "apt" {
		t.Fatalf("expected apt to be first, got %#v", profile.PackageManagers)
	}
}

func TestDetectsWindowsPackageManagersInPriorityOrder(t *testing.T) {
	profile := Detect(Inputs{
		GOOS:           "windows",
		AvailableTools: map[string]bool{"scoop": true, "winget": true},
	})

	if profile.OS != OSWindows {
		t.Fatalf("expected %q, got %q", OSWindows, profile.OS)
	}

	expected := []string{"winget", "scoop"}
	if len(profile.PackageManagers) != len(expected) {
		t.Fatalf("expected %d package managers, got %d", len(expected), len(profile.PackageManagers))
	}

	for index, name := range expected {
		if profile.PackageManagers[index] != name {
			t.Fatalf("expected %q at index %d, got %q", name, index, profile.PackageManagers[index])
		}
	}
}
