param(
    [string]$PythonBin = "python",
    [string]$BaseUrl = "http://127.0.0.1:8080"
)

$ErrorActionPreference = "Stop"

$repoRoot = Split-Path -Parent $PSScriptRoot
$botDir = Join-Path $repoRoot "bot-python"
$venvPython = Join-Path $botDir ".venv\Scripts\python.exe"
$envFile = Join-Path $botDir ".env"
$envExample = Join-Path $botDir ".env.example"

$ignoredProxyVars = @()
foreach ($proxyVar in @("HTTP_PROXY", "HTTPS_PROXY", "ALL_PROXY")) {
    $value = [Environment]::GetEnvironmentVariable($proxyVar)
    if ($value -eq "http://127.0.0.1:9") {
        Remove-Item "Env:$proxyVar" -ErrorAction SilentlyContinue
        $ignoredProxyVars += $proxyVar
    }
}

if (-not (Test-Path $envFile) -and (Test-Path $envExample)) {
    Copy-Item $envExample $envFile
}

$deps = @()
if (Test-Path $venvPython) {
    $deps = & $venvPython -c "import aiogram, httpx, pydantic_settings, dotenv; print('ok')"
}
else {
    $deps = "venv-missing"
}

$backendReachable = $false
$backendHealth = $null
$backendProvider = $null
try {
    $backend = Invoke-RestMethod -Uri "$BaseUrl/health" -Method Get
    $backendReachable = $true
    $backendHealth = $backend.status
    $backendProvider = $backend.payment_provider
}
catch {
    $backendReachable = $false
}

$envContent = @{}
if (Test-Path $envFile) {
    Get-Content $envFile | ForEach-Object {
        if ($_ -match '^\s*([^#=]+?)\s*=\s*(.*)\s*$') {
            $envContent[$matches[1]] = $matches[2]
        }
    }
}

[PSCustomObject]@{
    venv_present       = Test-Path $venvPython
    dependencies_ready = ($deps -eq "ok")
    env_present        = Test-Path $envFile
    bot_token_ready    = ($envContent.ContainsKey("BOT_TOKEN") -and $envContent["BOT_TOKEN"] -ne "" -and $envContent["BOT_TOKEN"] -ne "change-me")
    ignored_proxy_vars = $ignoredProxyVars
    backend_reachable  = $backendReachable
    backend_health     = $backendHealth
    backend_provider   = $backendProvider
} | ConvertTo-Json
