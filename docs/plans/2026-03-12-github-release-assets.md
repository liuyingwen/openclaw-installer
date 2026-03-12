# GitHub Release Assets Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Publish the installer binaries as downloadable assets on each GitHub Release.

**Architecture:** Keep the existing `build.yml` workflow focused on CI artifacts for pushes and pull requests. Add a dedicated `release.yml` workflow that reuses `scripts/build-release.sh` to cross-compile all release binaries on Ubuntu, then creates or updates the GitHub Release for a tag and uploads the built files plus a checksum manifest.

**Tech Stack:** GitHub Actions, GitHub CLI (`gh`), Bash, Go toolchain.

### Task 1: Add the GitHub Release workflow

**Files:**
- Create: `.github/workflows/release.yml`
- Modify: `README.md`

**Step 1: Write the failing test**

Confirm the repository does not already have a release workflow or release-upload step.

**Step 2: Run test to verify it fails**

Run: `test -f .github/workflows/release.yml`
Expected: exit code `1` because the workflow does not exist yet.

Run: `rg -n "gh release|action-gh-release|release:" .github/workflows`
Expected: no matches because release publishing is not wired yet.

**Step 3: Write minimal implementation**

Create `.github/workflows/release.yml` with:
- `push` trigger for version tags such as `v1.2.3`
- `workflow_dispatch` input for a manual tag name
- `contents: write` permission
- `go test ./...`
- `./scripts/build-release.sh`
- checksum generation for all files in `dist/`
- `gh release create` when the tag has no release yet
- `gh release upload --clobber` when the release already exists

**Step 4: Run test to verify it passes**

Run: `test -f .github/workflows/release.yml`
Expected: exit code `0`.

Run: `ruby -e 'require "yaml"; YAML.load_file(".github/workflows/release.yml")'`
Expected: exits successfully, proving the workflow file is valid YAML.

**Step 5: Commit**

```bash
git add .github/workflows/release.yml README.md
git commit -m "feat: publish installer binaries with github releases"
```

### Task 2: Document how to publish a release

**Files:**
- Modify: `README.md`

**Step 1: Write the failing test**

Check that the README does not explain how GitHub Releases are created or which assets are uploaded.

**Step 2: Run test to verify it fails**

Run: `rg -n "GitHub Release|workflow_dispatch|checksums|v1\\.2\\.3" README.md`
Expected: no relevant publishing guidance.

**Step 3: Write minimal implementation**

Add a short “GitHub Releases” section explaining:
- tag-driven release publishing
- manual dispatch with a tag input
- uploaded binaries and checksum file

**Step 4: Run test to verify it passes**

Run: `rg -n "GitHub Releases|workflow_dispatch|checksums|v1\\.2\\.3" README.md`
Expected: matches the new publishing guidance.

**Step 5: Commit**

```bash
git add README.md
git commit -m "docs: document github release publishing"
```

### Task 3: Verify build outputs still work

**Files:**
- Modify: `.github/workflows/release.yml`
- Modify: `README.md`

**Step 1: Run the Go test suite**

Run: `env GOCACHE=/tmp/openclaw-gocache GOMODCACHE=/tmp/openclaw-gomodcache go test ./...`
Expected: PASS.

**Step 2: Run the release build script**

Run: `env GOCACHE=/tmp/openclaw-gocache GOMODCACHE=/tmp/openclaw-gomodcache ./scripts/build-release.sh`
Expected: all four platform binaries are written to `dist/`.

**Step 3: Inspect the output**

Run: `ls -1 dist`
Expected: `openclaw-installer-darwin-arm64`, `openclaw-installer-darwin-amd64`, `openclaw-installer-linux-amd64`, and `openclaw-installer-windows-amd64.exe`.

**Step 4: Final review**

Confirm the release workflow uploads the same `dist/` files and generates a checksum manifest.

**Step 5: Commit**

```bash
git add .github/workflows/release.yml README.md
git commit -m "chore: verify release publishing workflow"
```
