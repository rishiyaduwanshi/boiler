#!/bin/bash
# Boiler Installation Script (Unix/Linux/macOS)

set -e

REPO="rishiyaduwanshi/boiler"
INSTALL_DIR="$HOME/.boiler/bin"
BINARY_NAME="bl"

echo "Installing Boiler CLI..."

# Create installation directory
mkdir -p "$INSTALL_DIR"

# Detect OS and architecture
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

case "$ARCH" in
    x86_64|amd64) ARCH="amd64" ;;
    arm64|aarch64) ARCH="arm64" ;;
    *) echo "Unsupported architecture: $ARCH"; exit 1 ;;
esac

# Fetch latest release
echo "Fetching latest release..."
RELEASE_URL="https://api.github.com/repos/$REPO/releases/latest"
DOWNLOAD_URL=$(curl -s $RELEASE_URL | grep "browser_download_url.*${OS}.*${ARCH}" | cut -d '"' -f 4 | head -n 1)

if [ -z "$DOWNLOAD_URL" ]; then
    echo "No binary found for $OS-$ARCH"
    echo "Please check: https://github.com/$REPO/releases"
    exit 1
fi

VERSION=$(curl -s $RELEASE_URL | grep '"tag_name"' | cut -d '"' -f 4)

# Download binary
echo "Downloading $VERSION..."
TEMP_FILE="/tmp/bl-download"
curl -fsSL "$DOWNLOAD_URL" -o "$TEMP_FILE"

# Install binary
if file "$TEMP_FILE" | grep -q "gzip"; then
    tar -xzf "$TEMP_FILE" -C "$INSTALL_DIR"
else
    mv "$TEMP_FILE" "$INSTALL_DIR/$BINARY_NAME"
fi

chmod +x "$INSTALL_DIR/$BINARY_NAME"

# Create wrapper script for 'boiler' alias
cat > "$INSTALL_DIR/boiler" << 'EOF'
#!/bin/bash
DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
"$DIR/bl" "$@"
EOF
chmod +x "$INSTALL_DIR/boiler"

# Add to PATH in shell config
SHELL_CONFIG=""
if [ -n "$BASH_VERSION" ]; then
    SHELL_CONFIG="$HOME/.bashrc"
elif [ -n "$ZSH_VERSION" ]; then
    SHELL_CONFIG="$HOME/.zshrc"
fi

if [ -n "$SHELL_CONFIG" ]; then
    if ! grep -q "$INSTALL_DIR" "$SHELL_CONFIG" 2>/dev/null; then
        echo "" >> "$SHELL_CONFIG"
        echo "# Boiler CLI" >> "$SHELL_CONFIG"
        echo "export PATH=\"\$PATH:$INSTALL_DIR\"" >> "$SHELL_CONFIG"
        echo "Added to PATH in $SHELL_CONFIG"
        echo "Run 'source $SHELL_CONFIG' or restart your terminal."
    else
        echo "Already in PATH."
    fi
fi

# Create default config directory
CONFIG_DIR="$HOME/.boiler"
mkdir -p "$CONFIG_DIR"

# Download uninstall script for offline use
echo "Downloading uninstall script..."
UNINSTALL_URL="https://raw.githubusercontent.com/$REPO/main/scripts/uninstall.sh"
if curl -fsSL "$UNINSTALL_URL" -o "$INSTALL_DIR/uninstall.sh" 2>/dev/null; then
    chmod +x "$INSTALL_DIR/uninstall.sh"
else
    echo "Warning: Could not download uninstall script. You can use 'bl self uninstall' instead."
fi

echo ""
echo "[SUCCESS] Boiler installed successfully!"
echo ""
echo "Installation directory: $INSTALL_DIR"
echo "Config directory: $CONFIG_DIR"
echo ""
echo "Run 'bl version' or 'boiler version' to verify installation."
echo "Run 'bl --help' to get started."
echo ""
echo "To uninstall: Run 'bl self uninstall' or '$INSTALL_DIR/uninstall.sh'"
