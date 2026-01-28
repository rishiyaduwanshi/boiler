# Boiler Installation Script (Windows PowerShell)

$ErrorActionPreference = "Stop"

Write-Host "Installing Boiler CLI..." -ForegroundColor Cyan

# Configuration
$INSTALL_DIR = "$env:USERPROFILE\.boiler\bin"
$REPO = "rishiyaduwanshi/boiler"
$BINARY_NAME = "bl.exe"

# Create installation directory
if (!(Test-Path $INSTALL_DIR)) {
    New-Item -ItemType Directory -Path $INSTALL_DIR -Force | Out-Null
}

# Get latest release info
Write-Host "Fetching latest release..." -ForegroundColor Yellow
$releaseUrl = "https://api.github.com/repos/$REPO/releases/latest"
try {
    $release = Invoke-RestMethod -Uri $releaseUrl
    $version = $release.tag_name
    $asset = $release.assets | Where-Object { $_.name -like "*windows*.zip" -or $_.name -like "*windows*.exe" -or $_.name -like "*bl.exe*" } | Select-Object -First 1
    
    if (!$asset) {
        throw "No Windows binary found in release"
    }
    
    $downloadUrl = $asset.browser_download_url
} catch {
    Write-Host "Error: No releases found. Please check https://github.com/$REPO/releases" -ForegroundColor Red
    exit 1
}

# Download binary
$extension = if ($asset.name -like "*.zip") { ".zip" } else { ".exe" }
$tempFile = "$env:TEMP\bl-download$extension"
Write-Host "Downloading $version..." -ForegroundColor Yellow
Invoke-WebRequest -Uri $downloadUrl -OutFile $tempFile

# Install binary
if ($asset.name -like "*.zip") {
    Write-Host "Extracting..." -ForegroundColor Yellow
    Expand-Archive -Path $tempFile -DestinationPath $INSTALL_DIR -Force
} else {
    Copy-Item $tempFile -Destination "$INSTALL_DIR\$BINARY_NAME" -Force
}
Remove-Item $tempFile -Force

# Create wrapper script for 'boiler' alias
$wrapperCmd = @"
@echo off
"%~dp0bl.exe" %*
"@
$wrapperCmd | Out-File -FilePath "$INSTALL_DIR\boiler.cmd" -Encoding ASCII -Force

# Add to PATH if not already present
$currentPath = [Environment]::GetEnvironmentVariable("Path", "User")
if ($currentPath -notlike "*$INSTALL_DIR*") {
    Write-Host "Adding to PATH..." -ForegroundColor Yellow
    [Environment]::SetEnvironmentVariable(
        "Path",
        "$currentPath;$INSTALL_DIR",
        "User"
    )
    $env:Path = "$env:Path;$INSTALL_DIR"
    Write-Host "Added to PATH. Restart your terminal for changes to take effect." -ForegroundColor Green
} else {
    Write-Host "Already in PATH." -ForegroundColor Green
}

# Create default config directory
$CONFIG_DIR = "$env:USERPROFILE\.boiler"
if (!(Test-Path $CONFIG_DIR)) {
    New-Item -ItemType Directory -Path $CONFIG_DIR -Force | Out-Null
}

Write-Host ""
Write-Host "✓ Boiler installed successfully!" -ForegroundColor Green
Write-Host ""
Write-Host "Installation directory: $INSTALL_DIR" -ForegroundColor Cyan
Write-Host "Config directory: $CONFIG_DIR" -ForegroundColor Cyan
Write-Host ""
Write-Host "Run 'bl version' or 'boiler version' to verify installation." -ForegroundColor Yellow
Write-Host "Run 'bl --help' to get started." -ForegroundColor Yellow
