# OpenClaw Installer Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Build a cross-platform one-command installer for OpenClaw that can diagnose the host environment, repair missing prerequisites, install required tools, install OpenClaw itself, and verify the result on macOS, Linux, and Windows.

**Architecture:** Implement a single-binary Go CLI with a layered design: platform detection and command execution at the bottom, package-manager and environment-repair abstractions in the middle, and an OpenClaw installation workflow driven by a YAML manifest at the top. Keep platform-specific operations declarative so the same workflow can choose `brew`, `apt`, `dnf`, `pacman`, `winget`, `choco`, or `scoop` at runtime.

**Tech Stack:** Go 1.25+, standard library, `gopkg.in/yaml.v3`, table-driven tests with `testing`.

### Task 1: Bootstrap the Go project

**Files:**
- Create: `go.mod`
- Create: `cmd/openclaw-installer/main.go`
- Create: `internal/cli/app.go`
- Create: `README.md`

**Step 1: Write the failing test**

Create a smoke test for the CLI app constructor in `internal/cli/app_test.go` that expects the app to expose `install`, `doctor`, and `print-plan` actions.

**Step 2: Run test to verify it fails**

Run: `go test ./...`
Expected: FAIL because the CLI package and constructor do not exist yet.

**Step 3: Write minimal implementation**

Add a small CLI application layer that parses the first argument and routes to stub handlers. Keep output text minimal and deterministic for testing.

**Step 4: Run test to verify it passes**

Run: `go test ./...`
Expected: PASS for the new smoke test.

**Step 5: Commit**

```bash
git add go.mod cmd/openclaw-installer/main.go internal/cli/app.go internal/cli/app_test.go README.md
git commit -m "feat: bootstrap openclaw installer cli"
```

### Task 2: Add configuration loading and validation

**Files:**
- Create: `config/openclaw.example.yaml`
- Create: `internal/config/config.go`
- Create: `internal/config/config_test.go`

**Step 1: Write the failing test**

Add tests that load an example YAML manifest and verify required fields:
- app metadata
- prerequisite tools
- installation steps
- verification commands

**Step 2: Run test to verify it fails**

Run: `go test ./internal/config -v`
Expected: FAIL because parsing and validation are not implemented.

**Step 3: Write minimal implementation**

Create a manifest model with validation methods and a `Load(path)` helper.

**Step 4: Run test to verify it passes**

Run: `go test ./internal/config -v`
Expected: PASS.

**Step 5: Commit**

```bash
git add config/openclaw.example.yaml internal/config/config.go internal/config/config_test.go
git commit -m "feat: add installer manifest loading"
```

### Task 3: Build platform and package-manager detection

**Files:**
- Create: `internal/platform/platform.go`
- Create: `internal/platform/platform_test.go`

**Step 1: Write the failing test**

Add table-driven tests covering:
- macOS detection
- Linux distro-family mapping
- Windows detection
- package-manager priority order selection

**Step 2: Run test to verify it fails**

Run: `go test ./internal/platform -v`
Expected: FAIL because detection and selection logic do not exist.

**Step 3: Write minimal implementation**

Implement platform normalization and package-manager selection from mocked command availability.

**Step 4: Run test to verify it passes**

Run: `go test ./internal/platform -v`
Expected: PASS.

**Step 5: Commit**

```bash
git add internal/platform/platform.go internal/platform/platform_test.go
git commit -m "feat: detect supported platforms and package managers"
```

### Task 4: Add environment repair logic

**Files:**
- Create: `internal/repair/repair.go`
- Create: `internal/repair/repair_test.go`

**Step 1: Write the failing test**

Add tests for a repair planner that:
- chooses the right install command for missing tools
- skips tools that already exist
- reports unsupported repair paths clearly

**Step 2: Run test to verify it fails**

Run: `go test ./internal/repair -v`
Expected: FAIL because the repair planner is missing.

**Step 3: Write minimal implementation**

Implement a planner that maps required tools to install commands per package manager.

**Step 4: Run test to verify it passes**

Run: `go test ./internal/repair -v`
Expected: PASS.

**Step 5: Commit**

```bash
git add internal/repair/repair.go internal/repair/repair_test.go
git commit -m "feat: add environment repair planner"
```

### Task 5: Add installer execution workflow

**Files:**
- Create: `internal/installer/installer.go`
- Create: `internal/installer/installer_test.go`

**Step 1: Write the failing test**

Add tests for the main workflow:
- doctor mode returns detected issues without executing commands
- install mode executes repair steps before install steps
- verification failure stops the workflow with a useful error

**Step 2: Run test to verify it fails**

Run: `go test ./internal/installer -v`
Expected: FAIL because the workflow service does not exist.

**Step 3: Write minimal implementation**

Implement an installer service using a mockable command runner and structured results.

**Step 4: Run test to verify it passes**

Run: `go test ./internal/installer -v`
Expected: PASS.

**Step 5: Commit**

```bash
git add internal/installer/installer.go internal/installer/installer_test.go
git commit -m "feat: add install workflow engine"
```

### Task 6: Wire the real CLI and documentation

**Files:**
- Modify: `cmd/openclaw-installer/main.go`
- Modify: `internal/cli/app.go`
- Modify: `README.md`

**Step 1: Write the failing test**

Add CLI integration tests that verify:
- `install --config config/openclaw.example.yaml --dry-run`
- `doctor --config config/openclaw.example.yaml`
- `print-plan --config config/openclaw.example.yaml`

**Step 2: Run test to verify it fails**

Run: `go test ./internal/cli -v`
Expected: FAIL because the flags and workflow hooks are incomplete.

**Step 3: Write minimal implementation**

Wire config loading, workflow invocation, dry-run output, and user-facing help text.

**Step 4: Run test to verify it passes**

Run: `go test ./internal/cli -v`
Expected: PASS.

**Step 5: Commit**

```bash
git add cmd/openclaw-installer/main.go internal/cli/app.go internal/cli/app_test.go README.md
git commit -m "feat: wire installer commands"
```

### Task 7: Verify end-to-end behavior

**Files:**
- Modify: `README.md`

**Step 1: Run the full test suite**

Run: `go test ./...`
Expected: PASS.

**Step 2: Run dry-run verification**

Run: `go run ./cmd/openclaw-installer install --config config/openclaw.example.yaml --dry-run`
Expected: prints the detected platform, missing tools, repair plan, install steps, and verification plan without changing the machine.

**Step 3: Run doctor verification**

Run: `go run ./cmd/openclaw-installer doctor --config config/openclaw.example.yaml`
Expected: prints environment findings and recommended repairs.

**Step 4: Update docs**

Document supported package managers, assumptions, and the remaining work needed to plug in the real OpenClaw artifact source.

**Step 5: Commit**

```bash
git add README.md
git commit -m "docs: document installer usage"
```
