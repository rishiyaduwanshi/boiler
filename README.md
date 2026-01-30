<div align="center">

<img src="web/src/assets/logo.svg" alt="Boiler Logo" width="120" />

# Boiler

**Code Once. Reuse Forever.**

Your personal code snippet and stack manager with automatic versioning, template variables, and zero configuration.

[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/badge/Go-1.25+-00ADD8?logo=go)](https://golang.org)
[![Release](https://img.shields.io/github/v/release/rishiyaduwanshi/boiler)](https://github.com/rishiyaduwanshi/boiler/releases)
[![Platform](https://img.shields.io/badge/Platform-Windows%20%7C%20Linux%20%7C%20macOS-lightgrey)](https://github.com/rishiyaduwanshi/boiler)

[Installation](#-installation) • [Quick Start](#-quick-start) • [Documentation](https://boiler.iamabhinav.dev) • [Use Cases](https://boiler.iamabhinav.dev/guides/usecases/)

</div>

---

## 🎯 Why Boiler?

**Have you ever thought:** Sometimes we add entire packages just for a single utility function?

```bash
npm install lodash    # 1.4 MB for one debounce function?
pip install requests  # Just for a simple HTTP wrapper?
```

**Boiler solves this.** Store reusable code snippets and project templates locally with automatic versioning, then inject them anywhere with a single command.

### The Problem
- 📋 Copy-pasting code between projects
- 📦 Installing bloated packages for simple utilities  
- 🔄 Manually versioning repeated code
- 🗂️ Losing track of reusable templates
- ⚙️ No standardized project setup across teams

### The Solution
```bash
# Store once
bl store ./utils/errorHandler.js
# → Saved as errorHandler@1.js

# Reuse everywhere
bl add errorHandler
# → Instant copy to current directory

# Update and version automatically
bl store ./utils/errorHandler-v2.js --name errorHandler.js
# → Saved as errorHandler@2.js
```

**Perfect for:**
- 🛠️ **Utility Functions** - Error handlers, loggers, validators you use in every project
- 🔧 **Configuration Files** - Database configs, middleware, authentication helpers
- 🚀 **Project Boilerplates** - Express APIs, Django templates, Next.js starters
- 📝 **Code Templates** - REST controllers, database models, API clients
- 🎨 **Custom Snippets** - Team-specific patterns and best practices

---

## 🚀 Installation

### Quick Install

**Windows (PowerShell)**
```powershell
iwr -useb https://boiler.iamabhinav.dev/install | iex
```

**Linux / macOS**
```bash
curl -fsSL https://boiler.iamabhinav.dev/install | bash
```

### Verify Installation
```bash
bl version
# Boiler v0.0.11

bl --help
# Shows all available commands
```

**System Requirements:**
- Windows 10+, Linux, or macOS
- No dependencies required (single binary)

---

## ⚡ Quick Start

### Initialize Boiler
```bash
bl init
# Creates configuration and storage directories
```

### Store & Reuse Snippets

**Store a file:**
```bash
bl store ./middleware/auth.js
# ✓ Saved as auth@1.js
```

**Add to any project:**
```bash
cd ~/new-project
bl add auth
# ✓ auth.js copied to current directory
```

**Update with new version:**
```bash
bl store ./middleware/auth-updated.js --name auth.js
# Choose: [o]verwrite, [n]ew version, or [c]ancel
# → Saved as auth@2.js
```

### 🎨 Template Variables (NEW!)

Create reusable code with configurable variables:

**Store with variables:**
```js
// errorHandler.js
// __author: Your Name
// __desc: Centralized error handler with configurable logging
// __version: 1.0.0
// __var bl__LOG_LEVEL = error
// __var bl__NOTIFY_EMAIL = admin@example.com

function handleError(err) {
  console[bl__LOG_LEVEL](err.message);
  sendEmail('bl__NOTIFY_EMAIL', err);
}
```

**Add and customize:**
```bash
bl add errorHandler
#   bl__LOG_LEVEL [error]: warn
#   bl__NOTIFY_EMAIL [admin@example.com]: dev@myapp.com
# ✓ errorHandler.js created (metadata stripped, variables replaced)
```

**Output (clean):**
```js
// errorHandler.js
function handleError(err) {
  console.warn(err.message);
  sendEmail('dev@myapp.com', err);
}
```

---

## 📚 Usage

### 📦 Snippets - Single Files

Store and version individual files:

```bash
# Store from anywhere
bl store ~/utils/logger.js
# → Saved as logger@1.js

# Add to current project
bl add logger
# → Copies latest version

# Add specific version
bl add logger@1
# → Copies logger@1.js

# List all snippets
bl ls --snippets

# Search snippets
bl search logger

# View details
bl info logger
```

**Supported languages:** JavaScript, Python, Go, Java, TypeScript, Rust, C, C++, Ruby, PHP, and more.

---

### 🚀 Stacks - Full Project Templates

Store entire directory structures with all dependencies:

```bash
# Store current directory as stack
bl store
# Name: express-api
# → Saved as express-api@1 with all files

# Store specific directory
bl store ./templates/nextjs-starter
# → Saved as nextjs-starter@1

# Add stack to new project
mkdir my-new-api && cd my-new-api
bl add express-api
# → Copies entire project structure

# List all stacks
bl ls --stacks

# View stack info
bl info express-api
```

**Use cases:** Express APIs, Django apps, microservices, React templates, config directories, documentation structures

---

## 🔥 Features

### Core Capabilities
- ✅ **Automatic Versioning** - Store multiple versions with `@1`, `@2`, etc.
- ✅ **Template Variables** - `bl__VAR_NAME` syntax with default values and prompts
- ✅ **Smart Version Management** - Choose to overwrite, create new version, or cancel
- ✅ **Language Agnostic** - Works with any file type (`.js`, `.py`, `.go`, `.java`, etc.)
- ✅ **Stack Templates** - Store entire directory structures as reusable projects
- ✅ **Metadata Support** - `__author`, `__desc`, `__version`, `__var` (auto-stripped from output)
- ✅ **Fuzzy Search** - Find snippets and stacks by partial name matching

### Developer Experience
- ✅ **Zero Configuration** - Works immediately after install
- ✅ **Interactive Prompts** - User-friendly CLI with defaults
- ✅ **Cross-Platform** - Windows, Linux, macOS support
- ✅ **Lightweight** - Single binary (~8MB), no runtime dependencies
- ✅ **Fast** - Instant file operations, no network calls for local resources
- ✅ **Self-Updating** - `bl self update` keeps CLI up-to-date
- ✅ **Secure** - SHA256 checksum verification on install

### Advanced Features
- 🔍 **Smart Search** - `bl search <query>` with fuzzy matching
- 📊 **Resource Info** - Detailed metadata view with `bl info <name>`
- 🧹 **Cleanup** - Remove unused resources with `bl clean`
- 📍 **Path Management** - `bl path` shows storage locations
- ⚙️ **Configurable** - `bl config` to customize settings
- 🗑️ **Easy Uninstall** - `bl self uninstall` removes everything cleanly

---

## 📖 Commands Reference

```bash
# Initialization
bl init                    # Initialize Boiler (first-time setup)

# Storage
bl store [path]            # Store file as snippet or folder as stack
bl store --name <name>     # Store with custom name

# Retrieval  
bl add <name>              # Add latest version
bl add <name@version>      # Add specific version

# Discovery
bl ls                      # List all snippets and stacks
bl ls --snippets           # List only snippets
bl ls --stacks             # List only stacks
bl search <query>          # Search by name (fuzzy matching)

# Information
bl info <name>             # Show detailed resource info
bl path                    # Show storage paths
bl config                  # Edit configuration file

# Maintenance
bl clean                   # Remove unused resources
bl self update             # Update Boiler to latest version
bl self uninstall          # Uninstall Boiler completely
bl version                 # Show current version
```

---

## 🎨 Template Variables Syntax

Create dynamic snippets with replaceable variables:

### Declaration Syntax
```js
// __var VARIABLE_NAME = default_value
```

### Example
```js
// apiClient.js
// __author: John Doe
// __desc: HTTP client with configurable base URL
// __var bl__API_URL = http://localhost:3000
// __var bl__TIMEOUT = 5000

const client = axios.create({
  baseURL: 'bl__API_URL',
  timeout: bl__TIMEOUT
});
```

### On Add
```bash
bl add apiClient
#   bl__API_URL [http://localhost:3000]: https://api.prod.com
#   bl__TIMEOUT [5000]: 10000
# ✓ apiClient.js created
```

### Output (metadata stripped)
```js
// apiClient.js
const client = axios.create({
  baseURL: 'https://api.prod.com',
  timeout: 10000
});
```

**Rules:**
- Variable names: `bl__[A-Za-z_][A-Za-z0-9_]*`
- Default values: Any string (spaces allowed)
- Metadata lines: `// __author`, `// __desc`, `// __version`, `// __var`
- All metadata automatically removed from output

**Full docs:** [Syntax Guide](https://boiler.iamabhinav.dev/guides/syntax/)

---

## 🌟 Use Cases

### 1. **Avoid Package Bloat**
Stop installing entire libraries for simple functions:
```bash
# Instead of: npm install lodash (1.4 MB)
bl store debounce.js       # Store your 20-line debounce
bl add debounce            # Add only what you need
```

### 2. **Standardize Team Workflows**  
Share snippets across your organization:
```bash
bl store company-eslint-config.js
bl store api-error-handler.js
# Team members: bl add company-eslint-config
```

### 3. **Bootstrap Projects Instantly**
```bash
bl add express-api         # Full Express setup in 1 second
bl add react-component     # Boilerplate component structure
```

### 4. **Version Control Your Utils**
```bash
bl store logger@1.js       # Basic logger
bl store logger@2.js       # Enhanced with colors
bl add logger@1            # Use stable version
```

**More examples:** [Use Cases Documentation](https://boiler.iamabhinav.dev/guides/usecases/)

---

## 📂 Project Structure

```
boiler/
├── cmd/boiler/          # CLI entry point
├── internal/
│   ├── cli/             # Command implementations (add, store, ls, etc.)
│   ├── config/          # Configuration management
│   ├── store/           # Storage operations
│   └── utils/           # Helpers (logger, prompts, fs)
├── scripts/             # Install/uninstall scripts
├── store/               # Default snippet/stack storage
├── web/                 # Documentation website (Starlight)
└── main.go              # Application entry
```

---

## 🤝 Contributing

Contributions are welcome! See [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

**Quick start:**
```bash
git clone https://github.com/rishiyaduwanshi/boiler.git
cd boiler
go mod download
go build -o bl main.go
./bl --help
```

**Areas to contribute:**
- 🐛 Bug fixes
- ✨ New features (plugin system, remote sync, etc.)
- 📝 Documentation improvements
- 🧪 Test coverage
- 🎨 UI/UX enhancements

---

## 📄 License

MIT © [Abhinav Prakash](https://github.com/rishiyaduwanshi)

See [LICENSE](LICENSE) for full details.

---

## 🔗 Links

- 📖 **Documentation:** [boiler.iamabhinav.dev](https://boiler.iamabhinav.dev)
- 🐛 **Issues:** [GitHub Issues](https://github.com/rishiyaduwanshi/boiler/issues)
- 💬 **Discussions:** [GitHub Discussions](https://github.com/rishiyaduwanshi/boiler/discussions)
- 📦 **Releases:** [GitHub Releases](https://github.com/rishiyaduwanshi/boiler/releases)

---

<div align="center">

**Made with 💜 for developers who value efficiency**

Stop repeating yourself. Start reusing.

[⭐ Star on GitHub](https://github.com/rishiyaduwanshi/boiler) • [📖 Documentation](https://boiler.iamabhinav.dev) • [🚀 Get Started](#-installation)

</div>