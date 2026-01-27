# Boiler Installation Guide

## Quick Install (One-Liner)

**Single command for all platforms:**

```powershell
# Windows (PowerShell)
iwr boiler.iamabhinav.dev/install -useb | iex
```

```bash
# Linux/macOS
curl -fsSL boiler.iamabhinav.dev/install | bash
```

> The same URL works for all platforms! OS is auto-detected and the appropriate binary is downloaded from GitHub Releases.

## What Happens During Installation

1. **Detects your OS and architecture** (Windows/Linux/macOS, amd64/arm64)
2. **Downloads latest release** from GitHub Releases (no Go required!)
3. **Installs to** `~/.boiler/bin/` (Windows: `C:\Users\YourName\.boiler\bin\`)
4. **Adds to PATH** automatically
5. **Creates wrapper** so both `bl` and `boiler` commands work
6. **Creates config directory** at `~/.boiler/`

## After Installation

Restart your terminal and verify:
```bash
bl version
boiler version
```

Both commands work!

## Alternative: Manual Download

If the one-liner doesn't work:

1. Go to [GitHub Releases](https://github.com/rishiyaduwanshi/boiler/releases)
2. Download the binary for your OS:
   - Windows: `boiler-windows-amd64.zip` or `bl.exe`
   - Linux: `boiler-linux-amd64.tar.gz`
   - macOS: `boiler-darwin-amd64.tar.gz` or `boiler-darwin-arm64.tar.gz`
3. Extract and place in a directory in your PATH
4. Test: `bl version`

## Development Setup

For contributors who want to work on the code:

```bash
# Clone repository
git clone https://github.com/rishiyaduwanshi/boiler.git
cd boiler

# Install mise (optional but recommended)
curl https://mise.run | sh

# Setup and build
mise install        # Installs Go 1.25.5
mise run build      # Builds binary
./bl.exe version    # Windows
./bl version        # Unix
```

**Or without mise:**
```bash
go mod download
go build -o bl.exe .    # Windows
go build -o bl .        # Unix
```

See [CONTRIBUTING.md](CONTRIBUTING.md) and [HOW_IT_WORKS.md](HOW_IT_WORKS.md) for more details.

## Uninstall

### One-Liner Uninstall

**Windows**:
```powershell
iwr https://raw.githubusercontent.com/rishiyaduwanshi/boiler/main/scripts/uninstall.ps1 -useb | iex
```

**Linux/macOS**:
```bash
curl -fsSL https://raw.githubusercontent.com/rishiyaduwanshi/boiler/main/scripts/uninstall.sh | bash
```

This will:
- Remove binary from `~/.boiler/bin`
- Remove from PATH
- Optionally remove config and store data

### Manual Uninstall

**Windows**:
```powershell
Remove-Item -Recurse -Force "$env:USERPROFILE\.boiler"
# Then remove from PATH via System Environment Variables
```

**Linux/macOS**:
```bash
rm -rf ~/.boiler
# Remove from PATH (edit ~/.bashrc or ~/.zshrc)
```

## Troubleshooting

### Command not found after installation

**Windows**: Restart PowerShell/CMD completely

**Linux/macOS**: 
```bash
source ~/.bashrc  # or ~/.zshrc
# Or restart your terminal
```

### Permission denied (Linux/macOS)

```bash
chmod +x ~/.boiler/bin/bl
chmod +x ~/.boiler/bin/boiler
```

### Download fails

- Check internet connection
- Verify GitHub is accessible
- Try manual download from [Releases](https://github.com/rishiyaduwanshi/boiler/releases)

### No binary for my OS/architecture

Currently supported:
- Windows (amd64)
- Linux (amd64, arm64)
- macOS (amd64, arm64)

If your system isn't supported, build from source (requires Go 1.25+):
```bash
git clone https://github.com/rishiyaduwanshi/boiler.git
cd boiler
go build -o bl .
```

### Manual PATH setup

**Windows**:
1. Search "Environment Variables" in Start Menu
2. Click "Environment Variables"
3. Under "User variables", select "Path" and click "Edit"
4. Click "New" and add: `C:\Users\YourName\.boiler\bin`
5. Click OK and restart terminal

**Linux/macOS**:
```bash
echo 'export PATH="$PATH:$HOME/.boiler/bin"' >> ~/.bashrc
source ~/.bashrc
```

For zsh:
```bash
echo 'export PATH="$PATH:$HOME/.boiler/bin"' >> ~/.zshrc
source ~/.zshrc
```

## System Requirements

- **OS**: Windows 10+, macOS 10.15+, Linux (any modern distro)
- **Architecture**: amd64 (x86_64) or arm64
- **Disk**: ~10-20 MB for binary
- **Internet**: Required for installation only
- **No Go required** - pre-built binaries provided

## What Gets Installed

```
~/.boiler/
├── bin/
│   ├── bl.exe (or bl)          # Main binary
│   └── boiler.cmd (or boiler)  # Wrapper alias
├── boiler.conf.json            # Config (created on first run)
├── store/                      # Snippets and stacks
│   └── boiler.meta.json
└── logs/                       # Application logs
```

## What Gets Installed

- **Binary**: `~/.boiler/bin/bl.exe` (Windows) or `~/.boiler/bin/bl` (Unix)
- **Wrapper**: `~/.boiler/bin/boiler` (for `boiler` command)
- **Config**: `~/.boiler/boiler.conf.json`
- **Store**: `~/.boiler/store/`
- **PATH**: Automatically added to user PATH

## After Installation

Restart your terminal and verify:
```bash
bl version
boiler version
```

Both commands should work!

## Manual Installation

If automatic install fails:

### Windows
```powershell
# Build
go build -o bl.exe .

# Move to a directory in your PATH
Move-Item bl.exe C:\Windows\System32\
# Or add current directory to PATH
```

### Linux/macOS
```bash
# Build
go build -o bl .

# Move to a directory in your PATH
sudo mv bl /usr/local/bin/

# Or add to user bin
mkdir -p ~/.local/bin
mv bl ~/.local/bin/
echo 'export PATH="$PATH:$HOME/.local/bin"' >> ~/.bashrc
```

## Uninstall

### Windows
```powershell
.\uninstall.ps1
```

### Linux/macOS
```bash
chmod +x uninstall.sh
./uninstall.sh
```

This will:
1. Remove binary from `~/.boiler/bin`
2. Remove from PATH
3. Optionally remove config and store data

## Development Installation

For developers who want to work on the code:

```bash
# Clone and setup
git clone https://github.com/rishiyaduwanshi/boiler.git
cd boiler

# Install mise (if not installed)
curl https://mise.run | sh

# Setup project
mise install
mise trust

# Build and test
mise run build
./bl.exe version  # Windows
./bl version      # Unix
```

See [CONTRIBUTING.md](CONTRIBUTING.md) for more details.

## Troubleshooting

### Command not found after installation

**Windows**: Restart PowerShell/CMD

**Linux/macOS**: 
```bash
source ~/.bashrc  # or ~/.zshrc
```

### Permission denied (Linux/macOS)

```bash
chmod +x ~/.boiler/bin/bl
chmod +x ~/.boiler/bin/boiler
```

### PATH not updated

Manually add to PATH:

**Windows**:
1. Search "Environment Variables"
2. Edit User PATH
3. Add `C:\Users\YourName\.boiler\bin`

**Linux/macOS**:
```bash
echo 'export PATH="$PATH:$HOME/.boiler/bin"' >> ~/.bashrc
source ~/.bashrc
```
