# Boiler Uninstallation Script (Windows PowerShell)

Write-Host "Uninstalling Boiler CLI..." -ForegroundColor Cyan

$INSTALL_DIR = "$env:USERPROFILE\.boiler\bin"
$CONFIG_DIR = "$env:USERPROFILE\.boiler"

# Ask for confirmation
$confirm = Read-Host "This will remove Boiler and all its data. Continue? (y/N)"
if ($confirm -ne "y" -and $confirm -ne "Y") {
    Write-Host "Uninstallation cancelled." -ForegroundColor Yellow
    exit 0
}

# Remove from PATH
$currentPath = [Environment]::GetEnvironmentVariable("Path", "User")
if ($currentPath -like "*$INSTALL_DIR*") {
    Write-Host "Removing from PATH..." -ForegroundColor Yellow
    $newPath = ($currentPath -split ';' | Where-Object { $_ -ne $INSTALL_DIR }) -join ';'
    [Environment]::SetEnvironmentVariable("Path", $newPath, "User")
}

# Remove binary directory
if (Test-Path $INSTALL_DIR) {
    Write-Host "Removing binary directory..." -ForegroundColor Yellow
    Remove-Item -Path $INSTALL_DIR -Recurse -Force
}

# Ask about config/data
$removeData = Read-Host "Remove config and store data? (y/N)"
if ($removeData -eq "y" -or $removeData -eq "Y") {
    if (Test-Path $CONFIG_DIR) {
        Write-Host "Removing config directory..." -ForegroundColor Yellow
        Remove-Item -Path $CONFIG_DIR -Recurse -Force
    }
}

Write-Host ""
Write-Host "✓ Boiler uninstalled successfully!" -ForegroundColor Green
Write-Host "Restart your terminal for PATH changes to take effect." -ForegroundColor Yellow
