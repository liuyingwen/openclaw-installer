package config

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type Config struct {
	App           AppSpec            `yaml:"app"`
	Prerequisites []PrerequisiteSpec `yaml:"prerequisites"`
	Install       []CommandStep      `yaml:"install"`
	Verify        []CommandStep      `yaml:"verify"`
}

type AppSpec struct {
	Name       string `yaml:"name"`
	Version    string `yaml:"version"`
	InstallDir string `yaml:"install_dir"`
}

type PrerequisiteSpec struct {
	Name     string            `yaml:"name"`
	Packages map[string]string `yaml:"packages"`
}

type CommandStep struct {
	Name    string            `yaml:"name"`
	Run     string            `yaml:"run"`
	RunByOS map[string]string `yaml:"run_by_os"`
}

func Load(path string) (Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return Config{}, fmt.Errorf("read config: %w", err)
	}
	defer file.Close()

	cfg, err := parseConfig(file)
	if err != nil {
		return Config{}, fmt.Errorf("parse config: %w", err)
	}

	if err := cfg.Validate(); err != nil {
		return Config{}, err
	}

	return cfg, nil
}

func parseConfig(file *os.File) (Config, error) {
	scanner := bufio.NewScanner(file)

	var (
		cfg             Config
		section         string
		currentPrereq   *PrerequisiteSpec
		currentCommand  *CommandStep
		packagesSection bool
	)

	flushPrereq := func() {
		if currentPrereq != nil {
			cfg.Prerequisites = append(cfg.Prerequisites, *currentPrereq)
			currentPrereq = nil
		}
	}

	flushCommand := func() {
		if currentCommand == nil {
			return
		}

		switch section {
		case "install":
			cfg.Install = append(cfg.Install, *currentCommand)
		case "verify":
			cfg.Verify = append(cfg.Verify, *currentCommand)
		}

		currentCommand = nil
	}

	for scanner.Scan() {
		rawLine := scanner.Text()
		line := strings.TrimSpace(rawLine)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		indent := len(rawLine) - len(strings.TrimLeft(rawLine, " "))
		switch {
		case indent == 0 && strings.HasSuffix(line, ":"):
			flushPrereq()
			flushCommand()
			packagesSection = false
			section = strings.TrimSuffix(line, ":")
		case section == "app" && indent == 2:
			key, value, ok := splitKeyValue(line)
			if !ok {
				return Config{}, fmt.Errorf("invalid app line: %q", rawLine)
			}
			assignAppValue(&cfg.App, key, value)
		case section == "prerequisites" && indent == 2 && strings.HasPrefix(line, "- "):
			flushPrereq()
			key, value, ok := splitKeyValue(strings.TrimPrefix(line, "- "))
			if !ok || key != "name" {
				return Config{}, fmt.Errorf("invalid prerequisite line: %q", rawLine)
			}
			currentPrereq = &PrerequisiteSpec{
				Name:     value,
				Packages: map[string]string{},
			}
			packagesSection = false
		case section == "prerequisites" && indent == 4 && line == "packages:":
			if currentPrereq == nil {
				return Config{}, fmt.Errorf("packages declared before prerequisite")
			}
			packagesSection = true
		case section == "prerequisites" && indent == 6 && packagesSection:
			if currentPrereq == nil {
				return Config{}, fmt.Errorf("package declared before prerequisite")
			}
			key, value, ok := splitKeyValue(line)
			if !ok {
				return Config{}, fmt.Errorf("invalid package line: %q", rawLine)
			}
			currentPrereq.Packages[key] = value
		case (section == "install" || section == "verify") && indent == 2 && strings.HasPrefix(line, "- "):
			flushCommand()
			key, value, ok := splitKeyValue(strings.TrimPrefix(line, "- "))
			if !ok || key != "name" {
				return Config{}, fmt.Errorf("invalid command line: %q", rawLine)
			}
			currentCommand = &CommandStep{Name: value}
		case (section == "install" || section == "verify") && indent == 4:
			if currentCommand == nil {
				return Config{}, fmt.Errorf("command detail declared before command")
			}
			if line == "run_by_os:" {
				currentCommand.RunByOS = map[string]string{}
				continue
			}
			key, value, ok := splitKeyValue(line)
			if !ok {
				return Config{}, fmt.Errorf("invalid command detail: %q", rawLine)
			}
			if key == "run" {
				currentCommand.Run = value
			}
		case (section == "install" || section == "verify") && indent == 6:
			if currentCommand == nil || currentCommand.RunByOS == nil {
				return Config{}, fmt.Errorf("os command declared before run_by_os")
			}
			key, value, ok := splitKeyValue(line)
			if !ok {
				return Config{}, fmt.Errorf("invalid os command detail: %q", rawLine)
			}
			currentCommand.RunByOS[key] = value
		default:
			return Config{}, fmt.Errorf("unsupported line: %q", rawLine)
		}
	}

	if err := scanner.Err(); err != nil {
		return Config{}, err
	}

	flushPrereq()
	flushCommand()

	return cfg, nil
}

func splitKeyValue(line string) (string, string, bool) {
	parts := strings.SplitN(line, ":", 2)
	if len(parts) != 2 {
		return "", "", false
	}

	return strings.TrimSpace(parts[0]), strings.Trim(strings.TrimSpace(parts[1]), `"`), true
}

func assignAppValue(app *AppSpec, key string, value string) {
	switch key {
	case "name":
		app.Name = value
	case "version":
		app.Version = value
	case "install_dir":
		app.InstallDir = value
	}
}

func (c Config) Validate() error {
	if c.App.Name == "" {
		return fmt.Errorf("app.name is required")
	}

	if len(c.Install) == 0 {
		return fmt.Errorf("at least one install step is required")
	}

	if len(c.Verify) == 0 {
		return fmt.Errorf("at least one verify step is required")
	}

	return nil
}

func (c CommandStep) CommandFor(goos string) string {
	if c.RunByOS != nil {
		if command := c.RunByOS[goos]; command != "" {
			return command
		}
	}
	return c.Run
}

func (c CommandStep) ResolveCommand(goos string) (string, error) {
	command := c.CommandFor(goos)
	if strings.TrimSpace(command) == "" {
		return "", fmt.Errorf("step %q does not support os %q", c.Name, goos)
	}
	return command, nil
}
