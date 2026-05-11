# DAP Installer for Windows
# This script builds DAP and adds it to your User PATH.

$installDir = Join-Path $HOME ".dap\bin"
if (!(Test-Path $installDir)) {
    Write-Host "Creating installation directory at $installDir..." -ForegroundColor Cyan
    New-Item -ItemType Directory -Path $installDir -Force | Out-Null
}

Write-Host "Building DAP..." -ForegroundColor Cyan
go build -o "$installDir\dap.exe" .

if ($LASTEXITCODE -ne 0) {
    Write-Error "Build failed! Please ensure Go is installed and you are in the DAP project directory."
    exit $LASTEXITCODE
}

Write-Host "DAP built successfully at $installDir\dap.exe" -ForegroundColor Green

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
