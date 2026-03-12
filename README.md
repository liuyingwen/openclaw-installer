# OpenClaw Installer

这是一个面向 macOS、Linux 和 Windows 的 OpenClaw 安装器 CLI。大多数用户不需要从源码构建，也不需要手动挑选二进制文件，直接执行一条命令即可安装 OpenClaw。

## 一键安装

### macOS / Linux

```bash
curl -fsSL https://raw.githubusercontent.com/liuyingwen/openclaw-installer/main/scripts/install-latest.sh | bash
```

### Windows（PowerShell）

```powershell
& ([scriptblock]::Create((irm https://raw.githubusercontent.com/liuyingwen/openclaw-installer/main/scripts/install-latest.ps1)))
```

这两个 bootstrap 脚本会自动识别当前系统，下载最新 GitHub Release 中匹配的安装器二进制，默认执行 `install --yes`，并在结束后清理临时文件。

当前 release 覆盖的二进制平台是：macOS（Apple Silicon / Intel）、Linux（x86_64）和 Windows（x86_64）。

## 安装前先检查

如果你想先看看安装器准备做什么，建议先把 bootstrap 脚本下载到本地，再把参数透传给安装器。

### macOS / Linux

```bash
curl -fsSL https://raw.githubusercontent.com/liuyingwen/openclaw-installer/main/scripts/install-latest.sh -o install-openclaw.sh
bash install-openclaw.sh doctor
bash install-openclaw.sh install --dry-run
bash install-openclaw.sh print-plan
```

### Windows（PowerShell）

```powershell
Invoke-WebRequest -Uri "https://raw.githubusercontent.com/liuyingwen/openclaw-installer/main/scripts/install-latest.ps1" -OutFile "install-openclaw.ps1"
.\install-openclaw.ps1 doctor
.\install-openclaw.ps1 install --dry-run
```

## 这个安装器会做什么

- 单文件 Go 二进制入口
- 支持的包管理器：`brew`、`apt`、`dnf`、`yum`、`pacman`、`winget`、`choco`、`scoop`
- 支持的命令：
  - `doctor`：检查当前机器环境，并输出缺失依赖和修复方案
  - `print-plan`：打印完整的修复、安装、校验命令清单
  - `install`：执行修复步骤、安装步骤和校验步骤
  - `install --dry-run`：只打印计划执行的命令，不真正修改系统
  - `--config` 为可选参数；如果不传，程序会使用内置的 OpenClaw 默认安装方案

如果你是在受限环境里做开发或测试，请先指定可写的 Go 缓存目录：

```bash
env GOCACHE=/tmp/openclaw-gocache GOMODCACHE=/tmp/openclaw-gomodcache go test ./...
```

## 高级配置

`config/openclaw.example.yaml` 现在只是一个可选覆盖示例。二进制内部已经内置了同样的默认安装计划，除非你想自定义行为，否则不需要额外分发 YAML 文件。

这个配置文件主要描述：

- `app`：应用元数据
- `prerequisites`：安装前必须存在的命令，以及各个包管理器对应的安装包名
- `install`：按顺序执行的安装命令，用来拉取并运行 OpenClaw 安装器
- `verify`：安装完成后必须成功执行的校验命令

示例配置当前指向的是 OpenClaw 官方推荐安装脚本：

- macOS/Linux: `https://openclaw.ai/install.sh --no-onboard`
- Windows: `https://openclaw.ai/install.ps1 -NoOnboard`

## 从源码构建

大多数用户不需要这一节。这里只面向本地开发和 release 维护。

```bash
go run ./cmd/openclaw-installer print-plan
go run ./cmd/openclaw-installer doctor
go run ./cmd/openclaw-installer install --dry-run
go run ./cmd/openclaw-installer install --yes
./scripts/build-release.sh
```

## 发布 Release

GitHub Release 由 `.github/workflows/release.yml` 自动发布。

- 推送一个版本标签，比如 `v1.2.3`，就会自动发布 release。
- 也可以在 GitHub Actions 里手动运行 `Release Installer` workflow，并填写目标 tag。
- `scripts/install-latest.sh` 和 `scripts/install-latest.ps1` 会从最新 release 拉取对应平台的二进制。
- 这个 workflow 会执行 `go test ./...`，然后运行 `./scripts/build-release.sh`，并把 `dist/` 里的所有二进制和 `openclaw-installer-checksums.txt` 上传到对应的 GitHub Release。
- 如果该 release 已经存在，workflow 会通过 `gh release upload --clobber` 覆盖同名产物。

## 说明

- 当前示例配置会调用 OpenClaw 官方推荐安装脚本，并安装默认 gateway。
- 如果系统里没有受支持的包管理器，安装器可以在 macOS 上引导安装 Homebrew，在 Windows 上引导安装 Scoop，然后再补齐依赖。
- Linux 上的包管理器引导策略目前比较保守，当前版本默认目标系统已经具备 `apt`、`dnf`、`yum` 或 `pacman` 其中之一。
- `install.sh --no-onboard` 会跳过交互式 onboarding，更适合一键安装场景；如果用户需要，后续可以再单独执行 onboarding。
- 官方推荐安装脚本会走全局安装路径；实际可执行文件通常会出现在 `/usr/local/bin`、`/opt/homebrew/bin`、`%APPDATA%\\npm` 或其他全局 npm bin 目录里，具体位置取决于当前系统和 Node/npm 前缀。
- `install` 现在要求显式传入 `--yes`，除非你使用的是 `--dry-run`，这样可以减少误执行。
- 如果某个步骤没有为当前操作系统定义命令，安装器会直接失败，而不是静默跳过。
- 安装成功后，程序会把实际执行的命令计划写入日志文件，并在最后输出日志路径。
- Unix 平台上的安装命令通过 `sh -c` 执行，Windows 平台通过 `cmd /C` 执行。
