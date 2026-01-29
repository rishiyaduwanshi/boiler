#!/bin/bash
# Boiler Uninstallation Script (Unix/Linux/macOS)

set -e

echo "Uninstalling Boiler CLI..."

INSTALL_DIR="$HOME/.boiler/bin"
ROOT_DIR="$HOME/.boiler"

# Ask for confirmation
read -p "This will remove Boiler from your system. Continue? (y/N) " confirm
if [ "$confirm" != "y" ] && [ "$confirm" != "Y" ]; then
    echo "Uninstallation cancelled."
    exit 0
fi

# Remove from PATH in shell configs
for config in "$HOME/.bashrc" "$HOME/.zshrc" "$HOME/.profile" "$HOME/.bash_profile"; do
    if [ -f "$config" ]; then
        if grep -q "$INSTALL_DIR" "$config"; then
            echo "Removing from $config..."
            sed -i.bak "/# Boiler CLI/d" "$config"
            sed -i.bak "\|export PATH=.*$INSTALL_DIR|d" "$config"
            rm -f "$config.bak"
        fi
    fi
done

# Ask about complete cleanup
echo ""
echo "Cleanup options:"
echo "  1. Keep my snippets and stacks (only remove binary)"
echo "  2. Remove everything (complete cleanup)"
echo ""
read -p "Choose option (1/2): " cleanup

if [ "$cleanup" = "2" ]; then
    # Complete cleanup - remove EVERYTHING
    echo ""
    echo "[WARNING] This will permanently delete:"
    echo "  • All snippets and stacks"
    echo "  • Configuration files"
    echo "  • Log files"
    echo "  • Binary files"
    echo ""
    read -p "Type 'DELETE' to confirm complete cleanup: " finalConfirm
    
    if [ "$finalConfirm" = "DELETE" ]; then
        # Remove entire .boiler directory
        if [ -d "$ROOT_DIR" ]; then
            echo "Removing all Boiler files..."
            rm -rf "$ROOT_DIR"
        fi
        
        # Clean temp files
        echo "Cleaning temporary files..."
        rm -f /tmp/bl-* 2>/dev/null || true
        rm -f /tmp/boiler-* 2>/dev/null || true
        
        echo ""
        echo "[SUCCESS] Complete cleanup done! Boiler completely removed."
    else
        echo "Complete cleanup cancelled."
        exit 0
    fi
else
    # Only remove binary
    if [ -d "$INSTALL_DIR" ]; then
        echo "Removing binary directory..."
        rm -rf "$INSTALL_DIR"
    fi
    
    echo ""
    echo "[SUCCESS] Boiler binary removed!"
    echo "Your snippets and stacks are preserved in: $ROOT_DIR"
fi

echo ""
echo "Restart your terminal or run 'source ~/.bashrc' (or ~/.zshrc) for PATH changes to take effect."
