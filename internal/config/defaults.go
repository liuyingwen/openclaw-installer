package config

const unixInstallScriptCommand = "curl -fsSL --proto '=https' --tlsv1.2 https://openclaw.ai/install.sh | bash -s -- --no-onboard"

func unixOpenClawCommand(args string) string {
	return "if command -v openclaw >/dev/null 2>&1; then openclaw " + args +
		`; elif [ -x /opt/homebrew/bin/openclaw ]; then /opt/homebrew/bin/openclaw ` + args +
		`; elif [ -x /usr/local/bin/openclaw ]; then /usr/local/bin/openclaw ` + args +
		`; elif [ -x "$HOME/.local/bin/openclaw" ]; then "$HOME/.local/bin/openclaw" ` + args +
		`; elif [ -x "$HOME/.npm-global/bin/openclaw" ]; then "$HOME/.npm-global/bin/openclaw" ` + args +
		`; elif command -v npm >/dev/null 2>&1 && [ -x "$(npm prefix -g 2>/dev/null)/bin/openclaw" ]; then "$(npm prefix -g 2>/dev/null)/bin/openclaw" ` + args +
		`; else openclaw ` + args + `; fi`
}

func windowsOpenClawCommand(args string) string {
	return "powershell -ExecutionPolicy Bypass -Command \"$cmd = Get-Command openclaw -ErrorAction SilentlyContinue; $appDataOpenClaw = Join-Path $env:APPDATA 'npm\\openclaw.cmd'; $localOpenClaw = Join-Path $env:USERPROFILE '.local\\bin\\openclaw.cmd'; if ($cmd -and $cmd.Definition) { & $cmd.Definition " + args + " } elseif (Test-Path $appDataOpenClaw) { & $appDataOpenClaw " + args + " } elseif (Test-Path $localOpenClaw) { & $localOpenClaw " + args + " } else { openclaw " + args + " }\""
}

func Default() Config {
	return Config{
		App: AppSpec{
			Name:       "openclaw",
			Version:    "latest",
			InstallDir: "~/.openclaw",
		},
		Install: []CommandStep{
			{
				Name: "install openclaw cli",
				RunByOS: map[string]string{
					"darwin":  unixInstallScriptCommand,
					"linux":   unixInstallScriptCommand,
					"windows": "powershell -ExecutionPolicy Bypass -Command \"& ([scriptblock]::Create((iwr -useb https://openclaw.ai/install.ps1))) -NoOnboard\"",
				},
			},
			{
				Name: "install default gateway",
				RunByOS: map[string]string{
					"darwin":  unixOpenClawCommand("gateway install"),
					"linux":   unixOpenClawCommand("gateway install"),
					"windows": windowsOpenClawCommand("gateway install"),
				},
			},
		},
		Verify: []CommandStep{
			{
				Name: "check binary",
				RunByOS: map[string]string{
					"darwin":  unixOpenClawCommand("--version"),
					"linux":   unixOpenClawCommand("--version"),
					"windows": windowsOpenClawCommand("--version"),
				},
			},
		},
	}
}
