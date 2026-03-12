package cli

import (
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/liuyingwen/openclaw-installer/internal/config"
	"github.com/liuyingwen/openclaw-installer/internal/installer"
	"github.com/liuyingwen/openclaw-installer/internal/platform"
	"github.com/liuyingwen/openclaw-installer/internal/repair"
)

type Dependencies interface {
	LoadConfig(path string) (config.Config, error)
	DetectProfile() (platform.Profile, error)
	DetectInstalled(prerequisites []config.PrerequisiteSpec) (map[string]bool, error)
	RunWorkflow(mode string, cfg config.Config, plan repair.Plan) (installer.Result, error)
}

type App struct {
	deps Dependencies
}

func NewApp(deps Dependencies) App {
	if deps == nil {
		deps = defaultDependencies{}
	}

	return App{deps: deps}
}

func (a App) Run(args []string, out io.Writer) error {
	if len(args) == 0 {
		return fmt.Errorf("missing command")
	}

	switch args[0] {
	case "install":
		return a.runInstall(args[1:], out)
	case "doctor":
		return a.runDoctor(args[1:], out)
	case "print-plan":
		return a.runPrintPlan(args[1:], out)
	default:
		return fmt.Errorf("unknown command: %s", args[0])
	}
}

func (a App) runInstall(args []string, out io.Writer) error {
	flags := flag.NewFlagSet("install", flag.ContinueOnError)
	flags.SetOutput(io.Discard)

	var configPath string
	var dryRun bool
	var assumeYes bool
	flags.StringVar(&configPath, "config", "", "config path")
	flags.BoolVar(&dryRun, "dry-run", false, "print plan only")
	flags.BoolVar(&assumeYes, "yes", false, "run without interactive confirmation")
	if err := flags.Parse(args); err != nil {
		return err
	}

	cfg, profile, repairPlan, err := a.prepare(configPath)
	if err != nil {
		return err
	}

	if dryRun {
		printProfile(out, profile)
		printPlan(out, "repair plan", repairPlan.Commands)
		planned, err := plannedCommands(string(profile.OS), repairPlan.Commands, cfg)
		if err != nil {
			return err
		}
		printPlan(out, "planned commands", planned)
		return nil
	}

	if !assumeYes {
		return fmt.Errorf("install requires --yes or --dry-run")
	}

	result, err := a.deps.RunWorkflow("install", cfg, repairPlan)
	if err != nil {
		return err
	}
	printPlan(out, "executed commands", result.PlannedCommands)
	if result.LogPath != "" {
		fmt.Fprintf(out, "log file: %s\n", result.LogPath)
	}
	printPostInstallGuidance(out, profile.OS)
	return nil
}

func (a App) runDoctor(args []string, out io.Writer) error {
	flags := flag.NewFlagSet("doctor", flag.ContinueOnError)
	flags.SetOutput(io.Discard)

	var configPath string
	flags.StringVar(&configPath, "config", "", "config path")
	if err := flags.Parse(args); err != nil {
		return err
	}

	cfg, profile, repairPlan, err := a.prepare(configPath)
	if err != nil {
		return err
	}

	printProfile(out, profile)
	printInstalled(out, cfg.Prerequisites, repairPlan.Commands)
	printPlan(out, "repair plan", repairPlan.Commands)
	planned, err := plannedCommands(string(profile.OS), repairPlan.Commands, cfg)
	if err != nil {
		return err
	}
	printPlan(out, "planned commands", planned)
	return nil
}

func (a App) runPrintPlan(args []string, out io.Writer) error {
	flags := flag.NewFlagSet("print-plan", flag.ContinueOnError)
	flags.SetOutput(io.Discard)

	var configPath string
	flags.StringVar(&configPath, "config", "", "config path")
	if err := flags.Parse(args); err != nil {
		return err
	}

	cfg, profile, repairPlan, err := a.prepare(configPath)
	if err != nil {
		return err
	}

	planned, err := plannedCommands(string(profile.OS), repairPlan.Commands, cfg)
	if err != nil {
		return err
	}
	printPlan(out, "planned commands", planned)
	return nil
}

func (a App) prepare(configPath string) (config.Config, platform.Profile, repair.Plan, error) {
	cfg := config.Default()
	var err error
	if configPath != "" {
		cfg, err = a.deps.LoadConfig(configPath)
		if err != nil {
			return config.Config{}, platform.Profile{}, repair.Plan{}, err
		}
	}

	profile, err := a.deps.DetectProfile()
	if err != nil {
		return config.Config{}, platform.Profile{}, repair.Plan{}, err
	}

	installed, err := a.deps.DetectInstalled(cfg.Prerequisites)
	if err != nil {
		return config.Config{}, platform.Profile{}, repair.Plan{}, err
	}

	repairPlan, err := (repair.Planner{}).Plan(profile, cfg.Prerequisites, installed)
	if err != nil {
		return config.Config{}, platform.Profile{}, repair.Plan{}, err
	}

	return cfg, profile, repairPlan, nil
}

func plannedCommands(goos string, repairCommands []string, cfg config.Config) ([]string, error) {
	planned := append([]string(nil), repairCommands...)
	for _, step := range cfg.Install {
		command, err := step.ResolveCommand(goos)
		if err != nil {
			return nil, err
		}
		planned = append(planned, command)
	}
	for _, step := range cfg.Verify {
		command, err := step.ResolveCommand(goos)
		if err != nil {
			return nil, err
		}
		planned = append(planned, command)
	}
	return planned, nil
}

func printProfile(out io.Writer, profile platform.Profile) {
	fmt.Fprintf(out, "platform: %s\n", profile.OS)
	if len(profile.PackageManagers) > 0 {
		fmt.Fprintf(out, "package managers: %s\n", strings.Join(profile.PackageManagers, ", "))
	}
}

func printInstalled(out io.Writer, prerequisites []config.PrerequisiteSpec, repairCommands []string) {
	missingCount := len(repairCommands)
	fmt.Fprintf(out, "prerequisites declared: %d\n", len(prerequisites))
	fmt.Fprintf(out, "missing prerequisites: %d\n", missingCount)
}

func printPlan(out io.Writer, title string, commands []string) {
	fmt.Fprintf(out, "%s:\n", title)
	if len(commands) == 0 {
		fmt.Fprintln(out, "- none")
		return
	}
	for _, command := range commands {
		fmt.Fprintf(out, "- %s\n", command)
	}
}

func printPostInstallGuidance(out io.Writer, osName platform.OS) {
	fmt.Fprintln(out, "openclaw should now be available on your PATH")
	fmt.Fprintln(out, "try: openclaw --version")

	switch osName {
	case platform.OSMac, platform.OSLinux:
		fmt.Fprintln(out, "if your current shell still cannot find it, open a new terminal and run:")
		fmt.Fprintln(out, "command -v openclaw")
	case platform.OSWindows:
		fmt.Fprintln(out, "if your current PowerShell still cannot find it, open a new terminal and run:")
		fmt.Fprintln(out, "Get-Command openclaw")
	}
}

type defaultDependencies struct{}

func (defaultDependencies) LoadConfig(path string) (config.Config, error) {
	return config.Load(path)
}

func (defaultDependencies) DetectProfile() (platform.Profile, error) {
	availableTools := map[string]bool{}
	for _, name := range []string{"brew", "apt", "dnf", "yum", "pacman", "winget", "choco", "scoop"} {
		_, err := exec.LookPath(name)
		availableTools[name] = err == nil
	}

	return platform.Detect(platform.Inputs{
		GOOS:           runtime.GOOS,
		IDLike:         linuxIDLike(),
		AvailableTools: availableTools,
	}), nil
}

func (defaultDependencies) DetectInstalled(prerequisites []config.PrerequisiteSpec) (map[string]bool, error) {
	result := map[string]bool{}
	for _, prerequisite := range prerequisites {
		_, err := exec.LookPath(prerequisite.Name)
		result[prerequisite.Name] = err == nil
	}
	return result, nil
}

func (defaultDependencies) RunWorkflow(mode string, cfg config.Config, plan repair.Plan) (installer.Result, error) {
	workflow := installer.Workflow{
		Runner: shellRunner{},
		GOOS:   runtime.GOOS,
	}
	if mode == "doctor" {
		return workflow.Doctor(cfg, plan)
	}
	return workflow.Install(cfg, plan)
}

type shellRunner struct{}

func (shellRunner) Run(command string) error {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/C", command)
	} else {
		cmd = exec.Command("sh", "-c", command)
	}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func linuxIDLike() []string {
	if runtime.GOOS != "linux" {
		return nil
	}

	content, err := os.ReadFile("/etc/os-release")
	if err != nil {
		return nil
	}

	for _, line := range strings.Split(string(content), "\n") {
		if !strings.HasPrefix(line, "ID_LIKE=") {
			continue
		}
		value := strings.TrimPrefix(line, "ID_LIKE=")
		value = strings.Trim(value, `"`)
		if value == "" {
			return nil
		}
		return strings.Fields(value)
	}

	return nil
}
