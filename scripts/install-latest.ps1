$ErrorActionPreference = "Stop"

$Repo = if ($env:OPENCLAW_INSTALLER_REPO) {
  $env:OPENCLAW_INSTALLER_REPO
} else {
  "liuyingwen/openclaw-installer"
}

function Resolve-AssetName {
  param(
    [string]$Architecture
  )

  if (-not $Architecture) {
    $Architecture = [System.Runtime.InteropServices.RuntimeInformation]::OSArchitecture.ToString()
  }

  $architectureValue = $Architecture.ToUpperInvariant()

  switch ($architectureValue) {
    "X64" { return "openclaw-installer-windows-amd64.exe" }
    "AMD64" { return "openclaw-installer-windows-amd64.exe" }
    default { throw "unsupported Windows architecture: $architectureValue" }
  }
}

function Get-LatestReleaseUrl {
  param(
    [Parameter(Mandatory = $true)]
    [string]$AssetName
  )

  return "https://github.com/$Repo/releases/latest/download/$AssetName"
}

function Invoke-InstallerBootstrap {
  param(
    [string[]]$InstallerArgs
  )

  if (-not $InstallerArgs -or $InstallerArgs.Count -eq 0) {
    $InstallerArgs = @("install", "--yes")
  }

  $assetName = Resolve-AssetName
  $url = Get-LatestReleaseUrl -AssetName $assetName
  $tempDir = Join-Path ([System.IO.Path]::GetTempPath()) ("openclaw-installer-" + [System.Guid]::NewGuid().ToString("N"))
  $binaryPath = Join-Path $tempDir "openclaw-installer.exe"

  New-Item -ItemType Directory -Path $tempDir | Out-Null

  try {
    Write-Host "Downloading $assetName..."
    Invoke-WebRequest -Uri $url -OutFile $binaryPath
    & $binaryPath @InstallerArgs
    if ($LASTEXITCODE -ne 0) {
      exit $LASTEXITCODE
    }
  }
  finally {
    Remove-Item -Path $tempDir -Recurse -Force -ErrorAction SilentlyContinue
  }
}

if (-not $env:OPENCLAW_INSTALLER_INSTALL_LATEST_PWSH_SOURCED) {
  Invoke-InstallerBootstrap @args
}
