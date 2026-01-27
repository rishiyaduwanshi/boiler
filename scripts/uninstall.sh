#!/bin/bash
# Boiler Uninstallation Script (Unix/Linux/macOS)

set -e

echo "Uninstalling Boiler CLI..."

INSTALL_DIR="$HOME/.boiler/bin"
CONFIG_DIR="$HOME/.boiler"

# Ask for confirmation
read -p "This will remove Boiler and all its data. Continue? (y/N) " confirm
if [ "$confirm" != "y" ] && [ "$confirm" != "Y" ]; then
    echo "Uninstallation cancelled."
    exit 0
fi

# Remove from PATH in shell configs
for config in "$HOME/.bashrc" "$HOME/.zshrc" "$HOME/.profile"; do
    if [ -f "$config" ]; then
        if grep -q "$INSTALL_DIR" "$config"; then
            echo "Removing from $config..."
            sed -i.bak "/# Boiler CLI/d" "$config"
            sed -i.bak "\|export PATH=.*$INSTALL_DIR|d" "$config"
            rm -f "$config.bak"
        fi
    fi
done

# Remove binary directory
if [ -d "$INSTALL_DIR" ]; then
    echo "Removing binary directory..."
    rm -rf "$INSTALL_DIR"
fi

# Ask about config/data
read -p "Remove config and store data? (y/N) " removeData
if [ "$removeData" = "y" ] || [ "$removeData" = "Y" ]; then
    if [ -d "$CONFIG_DIR" ]; then
        echo "Removing config directory..."
        rm -rf "$CONFIG_DIR"
    fi
fi

echo ""
echo "✓ Boiler uninstalled successfully!"
echo "Restart your terminal or run 'source ~/.bashrc' (or ~/.zshrc) for PATH changes to take effect."
