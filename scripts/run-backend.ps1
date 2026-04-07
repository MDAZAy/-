param(
    [string]$GoBin = "D:\bin\go.exe",
    [switch]$BuildOnly
)

$ErrorActionPreference = "Stop"

$repoRoot = Split-Path -Parent $PSScriptRoot
$backendDir = Join-Path $repoRoot "backend-go"
$cacheDir = Join-Path $backendDir ".cache"
$binDir = Join-Path $cacheDir "bin"
$goCache = Join-Path $cacheDir "go-build"
$goModCache = Join-Path $cacheDir "gomod"
$binaryPath = Join-Path $binDir "backend.exe"
$pidPath = Join-Path $cacheDir "backend.pid"
$stamp = Get-Date -Format "yyyyMMdd-HHmmss"
$outLog = Join-Path $cacheDir "backend-$stamp.out.log"
$errLog = Join-Path $cacheDir "backend-$stamp.err.log"
$logMetaPath = Join-Path $cacheDir "backend.logs.txt"

New-Item -ItemType Directory -Force $binDir, $goCache, $goModCache | Out-Null

$env:GOCACHE = $goCache
$env:GOMODCACHE = $goModCache

Push-Location $backendDir
try {
    & $GoBin mod tidy
    & $GoBin build -o ".cache\bin\backend.exe" "./cmd/server"
}
finally {
    Pop-Location
}

if ($BuildOnly) {
    Write-Output "Backend built: $binaryPath"
    exit 0
}

if (Test-Path $pidPath) {
    $existingPid = (Get-Content $pidPath -ErrorAction SilentlyContinue | Select-Object -First 1)
    if ($existingPid) {
        $existing = Get-Process -Id $existingPid -ErrorAction SilentlyContinue
        if ($existing) {
            Stop-Process -Id $existingPid -Force
            Start-Sleep -Seconds 1
        }
    }
}

$proc = Start-Process -FilePath $binaryPath `
    -WorkingDirectory $backendDir `
    -RedirectStandardOutput $outLog `
    -RedirectStandardError $errLog `
    -PassThru

Set-Content -Path $pidPath -Value $proc.Id
Set-Content -Path $logMetaPath -Value @($outLog, $errLog)
Start-Sleep -Seconds 2

Write-Output "Backend started. PID=$($proc.Id)"
Write-Output "Logs: $outLog"
