package config

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
					"darwin":  "curl -fsSL https://openclaw.ai/install-cli.sh | bash -s -- --no-onboard",
					"linux":   "curl -fsSL https://openclaw.ai/install-cli.sh | bash -s -- --no-onboard",
					"windows": "powershell -ExecutionPolicy Bypass -Command \"& ([scriptblock]::Create((iwr -useb https://openclaw.ai/install.ps1))) -NoOnboard\"",
				},
			},
			{
				Name: "install default gateway",
				RunByOS: map[string]string{
					"darwin":  "~/.openclaw/bin/openclaw gateway install",
					"linux":   "~/.openclaw/bin/openclaw gateway install",
					"windows": "\"%USERPROFILE%\\.openclaw\\bin\\openclaw.exe\" gateway install",
				},
			},
		},
		Verify: []CommandStep{
			{
				Name: "check binary",
				RunByOS: map[string]string{
					"darwin":  "~/.openclaw/bin/openclaw --version",
					"linux":   "~/.openclaw/bin/openclaw --version",
					"windows": "\"%USERPROFILE%\\.openclaw\\bin\\openclaw.exe\" --version",
				},
			},
		},
	}
}
