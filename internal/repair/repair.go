package repair

import (
	"fmt"

	"github.com/liuyingwen/openclaw-installer/internal/config"
	"github.com/liuyingwen/openclaw-installer/internal/platform"
)

type Planner struct{}

type Plan struct {
	Commands []string
}

func (Planner) Plan(profile platform.Profile, prerequisites []config.PrerequisiteSpec, installed map[string]bool) (Plan, error) {
	var plan Plan
	packageManagers := append([]string(nil), profile.PackageManagers...)

	if len(packageManagers) == 0 {
		bootstrapCommands, bootstrappedManagers := bootstrapPackageManagers(profile)
		plan.Commands = append(plan.Commands, bootstrapCommands...)
		packageManagers = append(packageManagers, bootstrappedManagers...)
	}

	for _, prerequisite := range prerequisites {
		if installed[prerequisite.Name] {
			continue
		}

		command, ok := installCommand(packageManagers, prerequisite.Packages)
		if !ok {
			return Plan{}, fmt.Errorf("no supported package manager found for %s", prerequisite.Name)
		}

		plan.Commands = append(plan.Commands, command)
	}

	return plan, nil
}

func bootstrapPackageManagers(profile platform.Profile) ([]string, []string) {
	switch profile.OS {
	case platform.OSMac:
		return []string{`/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"`}, []string{"brew"}
	case platform.OSWindows:
		return []string{`powershell -ExecutionPolicy Bypass -Command "iwr -useb get.scoop.sh | iex"`}, []string{"scoop"}
	default:
		return nil, nil
	}
}

func installCommand(packageManagers []string, packages map[string]string) (string, bool) {
	for _, manager := range packageManagers {
		pkg, ok := packages[manager]
		if !ok || pkg == "" {
			continue
		}

		switch manager {
		case "brew":
			return "brew install " + pkg, true
		case "apt":
			return "sudo apt-get install -y " + pkg, true
		case "dnf":
			return "sudo dnf install -y " + pkg, true
		case "yum":
			return "sudo yum install -y " + pkg, true
		case "pacman":
			return "sudo pacman -S --noconfirm " + pkg, true
		case "winget":
			return "winget install --id " + pkg + " --silent --accept-package-agreements --accept-source-agreements", true
		case "choco":
			return "choco install " + pkg + " -y", true
		case "scoop":
			return "scoop install " + pkg, true
		}
	}

	return "", false
}
