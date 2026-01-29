# Boiler Uninstallation Script (Windows PowerShell)

# Display banner
Write-Host ""
Write-Host "================================================" -ForegroundColor Red
Write-Host "         Boiler CLI Uninstaller                " -ForegroundColor Red
Write-Host "================================================" -ForegroundColor Red
Write-Host ""

$INSTALL_DIR = "$env:USERPROFILE\.boiler\bin"
$ROOT_DIR = "$env:USERPROFILE\.boiler"
$TEMP_DIR = "$env:TEMP"

# Ask for confirmation
Write-Host "This will remove Boiler from your system." -ForegroundColor Yellow
Write-Host ""
$confirm = Read-Host "Continue? (y/N)"
if ($confirm -ne "y" -and $confirm -ne "Y") {
    Write-Host ""
    Write-Host "Uninstallation cancelled." -ForegroundColor Yellow
    Write-Host ""
    exit 0
}

Write-Host ""
Write-Host "[1/3] Removing from PATH..." -ForegroundColor Yellow
# Remove from PATH
$currentPath = [Environment]::GetEnvironmentVariable("Path", "User")
if ($currentPath -like "*$INSTALL_DIR*") {
    $newPath = ($currentPath -split ';' | Where-Object { $_ -ne $INSTALL_DIR }) -join ';'
    [Environment]::SetEnvironmentVariable("Path", $newPath, "User")
    Write-Host "      PATH cleaned" -ForegroundColor Gray
} else {
    Write-Host "      Not in PATH" -ForegroundColor Gray
}
Write-Host ""

# Ask about complete cleanup
Write-Host "[2/3] Cleanup options:" -ForegroundColor Yellow
Write-Host ""
Write-Host "  1. Keep my data (only remove binary)" -ForegroundColor White
Write-Host "  2. Delete everything (complete cleanup)" -ForegroundColor White
Write-Host ""
$cleanup = Read-Host "Choose option (1/2)"
Write-Host ""

if ($cleanup -eq "2") {
    # Complete cleanup - remove EVERYTHING
    Write-Host "[WARNING] This will permanently delete:" -ForegroundColor Red
    Write-Host "  - All snippets and stacks" -ForegroundColor Yellow
    Write-Host "  - Configuration files" -ForegroundColor Yellow
    Write-Host "  - Log files" -ForegroundColor Yellow
    Write-Host "  - Binary files" -ForegroundColor Yellow
    Write-Host ""
    Write-Host "This action cannot be undone!" -ForegroundColor Red
    Write-Host ""
    $finalConfirm = Read-Host "Type 'DELETE' to confirm"
    Write-Host ""
    
    if ($finalConfirm -eq "DELETE") {
        Write-Host "[3/3] Removing all files..." -ForegroundColor Yellow
        
        # Create cleanup script to run after current process exits
        $cleanupScript = @"
Start-Sleep -Seconds 2
if (Test-Path '$ROOT_DIR') {
    Remove-Item -Path '$ROOT_DIR' -Recurse -Force -ErrorAction SilentlyContinue
}
Get-ChildItem -Path '$TEMP_DIR' -Filter 'bl-*' -ErrorAction SilentlyContinue | Remove-Item -Force -ErrorAction SilentlyContinue
Get-ChildItem -Path '$TEMP_DIR' -Filter 'boiler-*' -ErrorAction SilentlyContinue | Remove-Item -Force -ErrorAction SilentlyContinue
"@
        
        $cleanupPath = "$TEMP_DIR\boiler-cleanup-$((Get-Random)).ps1"
        $cleanupScript | Out-File -FilePath $cleanupPath -Encoding UTF8
        
        Write-Host "      Deleting $ROOT_DIR..." -ForegroundColor Gray
        Write-Host "      Cleaning temporary files..." -ForegroundColor Gray
        
        # Start cleanup script in background and exit immediately
        Start-Process powershell -ArgumentList "-NoProfile -ExecutionPolicy Bypass -WindowStyle Hidden -File `"$cleanupPath`"" -WindowStyle Hidden
        
        Write-Host ""
        Write-Host "================================================" -ForegroundColor Green
        Write-Host "   Complete Cleanup Done!                     " -ForegroundColor Green
        Write-Host "================================================" -ForegroundColor Green
        Write-Host ""
        Write-Host "Boiler has been completely removed from your system." -ForegroundColor White
    } else {
        Write-Host "Complete cleanup cancelled." -ForegroundColor Yellow
        exit 0
    }
} else {
    Write-Host "[3/3] Removing binary..." -ForegroundColor Yellow
    # Only remove binary
    if (Test-Path $INSTALL_DIR) {
        Write-Host "      Deleting binary files..." -ForegroundColor Gray
        Remove-Item -Path $INSTALL_DIR -Recurse -Force -ErrorAction SilentlyContinue
    }
    
    Write-Host ""
    Write-Host "================================================" -ForegroundColor Green
    Write-Host "   Binary Removed!                            " -ForegroundColor Green
    Write-Host "================================================" -ForegroundColor Green
    Write-Host ""
    Write-Host "Your data is preserved in:" -ForegroundColor White
    Write-Host "  $ROOT_DIR" -ForegroundColor Cyan
}

Write-Host ""
Write-Host "Restart your terminal for PATH changes to take effect." -ForegroundColor Yellow
