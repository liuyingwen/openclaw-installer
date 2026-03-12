package platform

type OS string

const (
	OSLinux   OS = "linux"
	OSMac     OS = "darwin"
	OSWindows OS = "windows"
)

type Inputs struct {
	GOOS           string
	IDLike         []string
	AvailableTools map[string]bool
}

type Profile struct {
	OS              OS
	PackageManagers []string
}

func Detect(input Inputs) Profile {
	profile := Profile{
		OS: OS(input.GOOS),
	}

	switch profile.OS {
	case OSMac:
		profile.PackageManagers = availableInOrder(input.AvailableTools, []string{"brew"})
	case OSLinux:
		profile.PackageManagers = detectLinuxManagers(input)
	case OSWindows:
		profile.PackageManagers = availableInOrder(input.AvailableTools, []string{"winget", "choco", "scoop"})
	}

	return profile
}

func detectLinuxManagers(input Inputs) []string {
	if contains(input.IDLike, "debian") || contains(input.IDLike, "ubuntu") {
		if managers := availableInOrder(input.AvailableTools, []string{"apt", "dnf", "yum", "pacman"}); len(managers) > 0 {
			return managers
		}
	}

	if contains(input.IDLike, "rhel") || contains(input.IDLike, "fedora") {
		if managers := availableInOrder(input.AvailableTools, []string{"dnf", "yum", "apt", "pacman"}); len(managers) > 0 {
			return managers
		}
	}

	return availableInOrder(input.AvailableTools, []string{"apt", "dnf", "yum", "pacman"})
}

func availableInOrder(available map[string]bool, names []string) []string {
	var result []string
	for _, name := range names {
		if available[name] {
			result = append(result, name)
		}
	}
	return result
}

func contains(values []string, want string) bool {
	for _, value := range values {
		if value == want {
			return true
		}
	}
	return false
}
