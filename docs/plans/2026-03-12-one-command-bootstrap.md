# One-Command Bootstrap Installer Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Let users install OpenClaw with a single command that auto-detects the current platform, downloads the correct release binary, and runs the installer.

**Architecture:** Add a small Unix bootstrap script and a PowerShell bootstrap script in `scripts/`. The Unix script will detect `uname -s` and `uname -m`, map them to the existing release asset names, download the correct binary from `releases/latest/download`, execute it, and remove the temporary file. README becomes bootstrap-first, while source-build guidance moves to a secondary section.

**Tech Stack:** Bash, PowerShell, GitHub Releases, existing Go installer binary.

### Task 1: Add failing bootstrap tests

**Files:**
- Create: `scripts/test-install-latest.sh`
- Test: `scripts/install-latest.sh`

**Step 1: Write the failing test**

Create a shell test script that expects `scripts/install-latest.sh` to expose helper functions for asset resolution and latest-release download URLs.

**Step 2: Run test to verify it fails**

Run: `bash scripts/test-install-latest.sh`
Expected: FAIL because `scripts/install-latest.sh` does not exist yet.

**Step 3: Write minimal implementation**

Add a sourceable Bash script that defines the required helpers and only runs `main` when executed directly.

**Step 4: Run test to verify it passes**

Run: `bash scripts/test-install-latest.sh`
Expected: PASS.

**Step 5: Commit**

```bash
git add scripts/test-install-latest.sh scripts/install-latest.sh
git commit -m "test: cover bootstrap installer detection"
```

### Task 2: Add bootstrap installers

**Files:**
- Create: `scripts/install-latest.sh`
- Create: `scripts/install-latest.ps1`

**Step 1: Write the failing test**

Re-run the shell test and manually verify the Windows script path exists.

**Step 2: Run test to verify it fails**

Run: `test -f scripts/install-latest.ps1`
Expected: exit code `1` before the PowerShell script exists.

**Step 3: Write minimal implementation**

Implement:
- Unix bootstrap script with platform detection, latest-release download, cleanup, and argument pass-through
- PowerShell bootstrap script with Windows asset download and argument pass-through
- default behavior of `install --yes` when no extra arguments are provided

**Step 4: Run test to verify it passes**

Run: `test -f scripts/install-latest.ps1`
Expected: exit code `0`.

Run: `bash scripts/test-install-latest.sh`
Expected: PASS.

**Step 5: Commit**

```bash
git add scripts/install-latest.sh scripts/install-latest.ps1 scripts/test-install-latest.sh
git commit -m "feat: add one-command bootstrap installers"
```

### Task 3: Update docs and verify behavior

**Files:**
- Modify: `README.md`
- Modify: `.github/workflows/build.yml`
- Modify: `.github/workflows/release.yml`

**Step 1: Write the failing test**

Confirm README does not yet document the bootstrap scripts and workflows do not run the shell bootstrap test.

**Step 2: Run test to verify it fails**

Run: `rg -n "install-latest\\.sh|install-latest\\.ps1" README.md .github/workflows`
Expected: no matches or incomplete coverage.

**Step 3: Write minimal implementation**

Update README to show one-command bootstrap usage:
- Unix-like systems: `curl ... | bash`
- Windows: `irm ... | iex`

Add a shell-test step to the CI/release workflows so bootstrap detection stays verified.

**Step 4: Run test to verify it passes**

Run: `rg -n "install-latest\\.sh|install-latest\\.ps1" README.md .github/workflows`
Expected: matches in README and both workflows.

Run: `env GOCACHE=/tmp/openclaw-gocache GOMODCACHE=/tmp/openclaw-gomodcache go test ./...`
Expected: PASS.

Run: `bash scripts/test-install-latest.sh`
Expected: PASS.

**Step 5: Commit**

```bash
git add README.md .github/workflows/build.yml .github/workflows/release.yml
git commit -m "docs: document one-command installer bootstrap"
```
