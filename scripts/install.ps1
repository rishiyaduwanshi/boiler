# Boiler Installation Script (Windows PowerShell)

$ErrorActionPreference = "Stop"

# Display banner
Write-Host ""
Write-Host "================================================" -ForegroundColor Cyan
Write-Host "           Boiler CLI Installer                " -ForegroundColor Cyan
Write-Host "================================================" -ForegroundColor Cyan
Write-Host ""

# Configuration
$INSTALL_DIR = "$env:USERPROFILE\.boiler\bin"
$REPO = "rishiyaduwanshi/boiler"
$BINARY_NAME = "bl.exe"

# Create installation directory
Write-Host "[1/5] Creating installation directory..." -ForegroundColor Yellow
if (!(Test-Path $INSTALL_DIR)) {
    New-Item -ItemType Directory -Path $INSTALL_DIR -Force | Out-Null
    Write-Host "      Created: $INSTALL_DIR" -ForegroundColor Gray
} else {
    Write-Host "      Directory exists" -ForegroundColor Gray
}
Write-Host ""

# Get latest release info
Write-Host "[2/5] Fetching latest release..." -ForegroundColor Yellow
$releaseUrl = "https://api.github.com/repos/$REPO/releases/latest"
try {
    $release = Invoke-RestMethod -Uri $releaseUrl
    $version = $release.tag_name
    Write-Host "      Found version: $version" -ForegroundColor Gray
    $asset = $release.assets | Where-Object { $_.name -like "*windows*.zip" -or $_.name -like "*windows*.exe" -or $_.name -like "*bl.exe*" } | Select-Object -First 1
    
    if (!$asset) {
        throw "No Windows binary found in release"
    }
    
    $downloadUrl = $asset.browser_download_url
} catch {
    Write-Host "      [ERROR] No releases found" -ForegroundColor Red
    Write-Host "      Please check https://github.com/$REPO/releases" -ForegroundColor Red
    exit 1
}
Write-Host ""

# Download binary
Write-Host "[3/5] Downloading binary..." -ForegroundColor Yellow
$extension = if ($asset.name -like "*.zip") { ".zip" } else { ".exe" }
$tempFile = "$env:TEMP\bl-download$extension"
Write-Host "      Downloading from GitHub..." -ForegroundColor Gray
Invoke-WebRequest -Uri $downloadUrl -OutFile $tempFile
Write-Host "      Download complete" -ForegroundColor Gray
Write-Host ""

# Install binary
Write-Host "[4/5] Installing binary..." -ForegroundColor Yellow
if ($asset.name -like "*.zip") {
    Write-Host "      Extracting archive..." -ForegroundColor Gray
    Expand-Archive -Path $tempFile -DestinationPath $INSTALL_DIR -Force
} else {
    Copy-Item $tempFile -Destination "$INSTALL_DIR\$BINARY_NAME" -Force
}
Remove-Item $tempFile -Force
Write-Host "      Binary installed" -ForegroundColor Gray

# Create wrapper script for 'boiler' alias
$wrapperCmd = @"
@echo off
"%~dp0bl.exe" %*
"@
$wrapperCmd | Out-File -FilePath "$INSTALL_DIR\boiler.cmd" -Encoding ASCII -Force
Write-Host "      Alias created: boiler -> bl" -ForegroundColor Gray
Write-Host ""

# Add to PATH if not already present
Write-Host "[5/5] Configuring PATH..." -ForegroundColor Yellow
$currentPath = [Environment]::GetEnvironmentVariable("Path", "User")
if ($currentPath -notlike "*$INSTALL_DIR*") {
    [Environment]::SetEnvironmentVariable(
        "Path",
        "$currentPath;$INSTALL_DIR",
        "User"
    )
    $env:Path = "$env:Path;$INSTALL_DIR"
    Write-Host "      Added to PATH" -ForegroundColor Gray
} else {
    Write-Host "      Already in PATH" -ForegroundColor Gray
}

# Create default config directory
$CONFIG_DIR = "$env:USERPROFILE\.boiler"
if (!(Test-Path $CONFIG_DIR)) {
    New-Item -ItemType Directory -Path $CONFIG_DIR -Force | Out-Null
}

# Download uninstall script for offline use
Write-Host "      Downloading uninstall script..." -ForegroundColor Gray
try {
    $uninstallUrl = "https://raw.githubusercontent.com/$REPO/main/scripts/uninstall.ps1"
    Invoke-WebRequest -Uri $uninstallUrl -OutFile "$INSTALL_DIR\uninstall.ps1" -ErrorAction SilentlyContinue
} catch {
    # Silent fail - not critical
}
Write-Host ""

# Success message
Write-Host "================================================" -ForegroundColor Green
Write-Host "   Installation Complete!                     " -ForegroundColor Green
Write-Host "================================================" -ForegroundColor Green
Write-Host ""
Write-Host "  Version:      $version" -ForegroundColor White
Write-Host "  Install path: $INSTALL_DIR" -ForegroundColor White
Write-Host "  Config path:  $CONFIG_DIR" -ForegroundColor White
Write-Host ""
Write-Host "Quick Start:" -ForegroundColor Cyan
Write-Host "  bl version          # Verify installation" -ForegroundColor Gray
Write-Host "  bl init             # Initialize Boiler" -ForegroundColor Gray
Write-Host "  bl --help           # Show all commands" -ForegroundColor Gray
Write-Host ""
Write-Host "Note: Restart your terminal for PATH changes to take effect" -ForegroundColor Yellow
Write-Host ""
Write-Host "To uninstall: Run 'bl self uninstall' or 'powershell $INSTALL_DIR\uninstall.ps1'" -ForegroundColor Cyan
