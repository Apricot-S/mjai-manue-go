$ErrorActionPreference = "Stop"

$scriptDir = Split-Path -Parent $MyInvocation.MyCommand.Path

docker compose -f (Join-Path $scriptDir "compose.yaml") up --build --abort-on-container-exit

$logDir = if ($env:LOG_DIR) { $env:LOG_DIR } else { Join-Path $scriptDir "out" }
Write-Host ""
Write-Host "Logs are in: $logDir"
