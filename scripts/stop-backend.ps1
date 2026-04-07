param()

$ErrorActionPreference = "Stop"

$repoRoot = Split-Path -Parent $PSScriptRoot
$pidPath = Join-Path $repoRoot "backend-go\.cache\backend.pid"

if (-not (Test-Path $pidPath)) {
    Write-Output "No backend pid file found."
    exit 0
}

$pid = Get-Content $pidPath | Select-Object -First 1
if (-not $pid) {
    Write-Output "PID file is empty."
    exit 0
}

$proc = Get-Process -Id $pid -ErrorAction SilentlyContinue
if ($proc) {
    Stop-Process -Id $pid -Force
    Write-Output "Backend stopped. PID=$pid"
}
else {
    Write-Output "Backend process not running. PID=$pid"
}

Remove-Item $pidPath -Force -ErrorAction SilentlyContinue

