$ErrorActionPreference = "Stop"
$env:OPENCLAW_INSTALLER_INSTALL_LATEST_PWSH_SOURCED = "1"

. "$PSScriptRoot/install-latest.ps1"

function Assert-Equal {
  param(
    [Parameter(Mandatory = $true)]
    [string]$Expected,
    [Parameter(Mandatory = $true)]
    [string]$Actual
  )

  if ($Expected -ne $Actual) {
    throw "expected: $Expected`nactual:   $Actual"
  }
}

function Assert-Throws {
  param(
    [Parameter(Mandatory = $true)]
    [scriptblock]$ScriptBlock
  )

  try {
    & $ScriptBlock
  }
  catch {
    return
  }

  throw "expected script block to throw"
}

Assert-Equal "openclaw-installer-windows-amd64.exe" (Resolve-AssetName "X64")
Assert-Equal "openclaw-installer-windows-amd64.exe" (Resolve-AssetName "AMD64")
Assert-Equal `
  "https://github.com/liuyingwen/openclaw-installer/releases/latest/download/openclaw-installer-windows-amd64.exe" `
  (Get-LatestReleaseUrl -AssetName "openclaw-installer-windows-amd64.exe")

Assert-Throws { Resolve-AssetName "ARM64" }

Write-Host "install-latest.ps1 tests passed"
