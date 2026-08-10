# DAP Installer for Windows
# This script installs DAP (downloads pre-built release or builds from source) and adds it to your User PATH.

$ErrorActionPreference = "Stop"
$repo = "rafiathallah3/dap"
$installDir = Join-Path $HOME ".dap\bin"

if (!(Test-Path $installDir)) {
    Write-Host "Creating installation directory at $installDir..." -ForegroundColor Cyan
    New-Item -ItemType Directory -Path $installDir -Force | Out-Null
}

$hasGo = $false
if (Test-Path "utama.go") {
    try {
        $null = Get-Command "go" -ErrorAction Stop
        $hasGo = $true
    } catch {
        $hasGo = $false
    }
}

if ($hasGo) {
    Write-Host "Building DAP from source..." -ForegroundColor Cyan
    go build -o "$installDir\dap.exe" .
    Write-Host "DAP built successfully!" -ForegroundColor Green
} else {
    Write-Host "Downloading pre-built DAP binary..." -ForegroundColor Cyan
    $downloadUrl = "https://github.com/$repo/releases/latest/download/dap-windows-amd64.zip"
    $zipPath = Join-Path $env:TEMP "dap-windows-amd64.zip"
    $extractPath = Join-Path $env:TEMP "dap_extract"

    [Net.ServicePointManager]::SecurityProtocol = [Net.SecurityProtocolType]::Tls12
    Invoke-WebRequest -Uri $downloadUrl -OutFile $zipPath -UseBasicParsing

    if (Test-Path $extractPath) { Remove-Item $extractPath -Recurse -Force }
    Expand-Archive -Path $zipPath -DestinationPath $extractPath -Force

    $exeFile = Get-ChildItem -Path $extractPath -Filter "dap.exe" -Recurse | Select-Object -First 1
    if ($exeFile) {
        Copy-Item -Path $exeFile.FullName -Destination "$installDir\dap.exe" -Force
    } else {
        Write-Error "Failed to locate dap.exe in downloaded archive."
        exit 1
    }

    Remove-Item $zipPath -ErrorAction SilentlyContinue
    Remove-Item $extractPath -Recurse -ErrorAction SilentlyContinue
}

Write-Host "DAP installed successfully at $installDir\dap.exe" -ForegroundColor Green

# Add to PATH if not already there
$path = [Environment]::GetEnvironmentVariable("Path", "User")
if ($path -notlike "*$installDir*") {
    Write-Host "Adding $installDir to User PATH..." -ForegroundColor Cyan
    $newPath = "$path;$installDir"
    [Environment]::SetEnvironmentVariable("Path", $newPath, "User")
    $env:Path += ";$installDir"
    Write-Host "Successfully added to PATH. You may need to restart your terminal for changes to take effect." -ForegroundColor Green
} else {
    Write-Host "PATH already contains $installDir." -ForegroundColor Yellow
}

Write-Host "`nInstallation complete! Try running 'dap' in a new terminal window." -ForegroundColor Green
