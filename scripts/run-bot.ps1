param(
    [string]$PythonBin = "python"
)

$ErrorActionPreference = "Stop"

$repoRoot = Split-Path -Parent $PSScriptRoot
$botDir = Join-Path $repoRoot "bot-python"
$venvDir = Join-Path $botDir ".venv"
$venvPython = Join-Path $venvDir "Scripts\python.exe"
$envFile = Join-Path $botDir ".env"
$envExample = Join-Path $botDir ".env.example"

foreach ($proxyVar in @("HTTP_PROXY", "HTTPS_PROXY", "ALL_PROXY")) {
    $value = [Environment]::GetEnvironmentVariable($proxyVar)
    if ($value -eq "http://127.0.0.1:9") {
        Remove-Item "Env:$proxyVar" -ErrorAction SilentlyContinue
        Write-Output "Ignored invalid $proxyVar=$value"
    }
}

Push-Location $botDir
try {
    if (-not (Test-Path $venvPython)) {
        & $PythonBin -m venv .venv
        & $venvPython -m pip install -r requirements.txt
    }

    if (-not (Test-Path $envFile) -and (Test-Path $envExample)) {
        Copy-Item $envExample $envFile
        Write-Output "Created bot-python/.env from .env.example"
    }

    $envContent = @{}
    if (Test-Path $envFile) {
        Get-Content $envFile | ForEach-Object {
            if ($_ -match '^\s*([^#=]+?)\s*=\s*(.*)\s*$') {
                $envContent[$matches[1]] = $matches[2]
            }
        }
    }

    if (-not $envContent.ContainsKey("BOT_TOKEN") -or [string]::IsNullOrWhiteSpace($envContent["BOT_TOKEN"]) -or $envContent["BOT_TOKEN"] -eq "change-me") {
        throw "Set a real BOT_TOKEN in bot-python/.env before starting the bot."
    }

    & $venvPython -m app.main
}
finally {
    Pop-Location
}
