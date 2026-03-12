# OpenClaw Installer

Cross-platform OpenClaw installer CLI for macOS, Linux, and Windows. The tool detects the host platform, checks required tools, prepares a repair plan using the available package manager, runs installation steps, and verifies the final result.

## Current scope

- Single Go binary entrypoint
- Supported package managers: `brew`, `apt`, `dnf`, `yum`, `pacman`, `winget`, `choco`, `scoop`
- Commands:
  - `doctor`: inspect the current machine and print missing prerequisites plus a repair plan
  - `print-plan`: print the full repair + install + verify command sequence
  - `install`: execute the repair plan, install steps, and verification steps
  - `install --dry-run`: print the planned commands without modifying the machine
  - `--config` is optional; if omitted, the binary uses a built-in OpenClaw install plan

## Quick start

```bash
go run ./cmd/openclaw-installer print-plan
go run ./cmd/openclaw-installer doctor
go run ./cmd/openclaw-installer install --dry-run
go run ./cmd/openclaw-installer install --yes
```

## Build binaries

```bash
./scripts/build-release.sh
```

When running inside a restricted environment, set a writable Go cache first:

```bash
env GOCACHE=/tmp/openclaw-gocache GOMODCACHE=/tmp/openclaw-gomodcache go test ./...
```

## Manifest

`config/openclaw.example.yaml` is now just an optional override example. The binary embeds the same default plan, so you do not need to ship a YAML file unless you want to customize behavior.

The manifest describes:

- `app`: application metadata
- `prerequisites`: optional commands that must exist before installation and the package names to install with each package manager
- `install`: ordered shell commands for fetching and running the OpenClaw installer
- `verify`: ordered shell commands that must succeed after installation

The example manifest now points at the official OpenClaw installers:

- macOS/Linux: `https://openclaw.ai/install-cli.sh --no-onboard`
- Windows: `https://openclaw.ai/install.ps1 -NoOnboard`

## Notes

- The current example manifest uses the official OpenClaw installation scripts and then installs the default gateway.
- When no supported package manager is present, the planner can bootstrap Homebrew on macOS and Scoop on Windows before installing prerequisites.
- Linux package-manager bootstrapping is intentionally conservative. The current version expects the target distro to already expose one of `apt`, `dnf`, `yum`, or `pacman`.
- `install-cli.sh --no-onboard` avoids interactive onboarding during installation, which makes it a better fit for one-command automation. Users can run onboarding later if needed.
- `install` now requires `--yes` unless you use `--dry-run`, so automation is explicit and accidental execution is harder.
- If a step does not define a command for the current OS, the installer now fails fast instead of silently skipping it.
- Successful installs write the executed command plan to a log file and print the path at the end.
- Installation commands run with `sh -c` on Unix and `cmd /C` on Windows.
