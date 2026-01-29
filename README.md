<div align="center">

# ⚡ Boiler

**Your personal code library. Store once, reuse everywhere.**

[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/badge/Go-1.25+-00ADD8?logo=go)](https://golang.org)
[![Release](https://img.shields.io/github/v/release/rishiyaduwanshi/boiler)](https://github.com/rishiyaduwanshi/boiler/releases)

[Installation](#installation) • [Usage](#usage) • [Documentation](https://boiler.iamabhinav.dev)

</div>

---

## Why Boiler?

Stop copy-pasting code between projects. Boiler lets you **store reusable code once** and **inject it anywhere** with automatic versioning.

```bash
# Store your utility function
bl store ./utils/errorHandler.js

# Add it to any project instantly
bl add errorHandler

# Done. It's in your current directory.
```

**Perfect for:**
- API middleware you use in every project
- Database configs that never change
- Logging utilities
- Authentication helpers
- Project boilerplates (Express, Django, Next.js stacks)

---

## Installation

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
```

---

## Usage

### 📦 Snippets - Single Files

Store individual files with automatic versioning:

```bash
# Store a file from anywhere
bl store ./middleware/auth.js
# → Saved as auth@1.js

# Add it to current project
bl add auth
# → Copies auth@1.js to ./

# Store multiple versions
bl store ./middleware/auth-v2.js --name auth.js
# → Saved as auth@2.js

# List all snippets
bl ls --snippets
```

**Supports all languages:** `.js`, `.py`, `.go`, `.java`, `.ts`, `.rb`, etc.

---

### 🚀 Stacks - Full Project Templates

Store entire directory structures as reusable templates:

```bash
# Store current directory as a stack
bl store
# Name: express-api → Saved as express-api@1

# Store specific directory
bl store ./my-nextjs-template
# → Saved as my-nextjs-template@1

# Add stack to new project folder
bl add express-api@1
# → Copies entire stack structure

# List all stacks
bl ls --stacks
```

**Use cases:** Express APIs, Next.js templates, microservice boilerplates, config folders

---

## Commands

```bash
bl init              # Initialize Boiler in current directory
bl store [path]      # Store file (snippet) or folder (stack)
bl add <name>        # Add snippet/stack to current directory
bl ls                # List all snippets and stacks
bl search <query>    # Search for resources by name
bl info <name>       # Show detailed info about a resource
bl clean             # Remove unused resources
bl path              # Show Boiler storage paths
bl config            # Edit configuration
bl self update       # Update Boiler to latest version
bl self uninstall    # Uninstall Boiler
```

---

## Features

- ✅ **Automatic Versioning** - Store multiple versions of the same file
- ✅ **Language Agnostic** - Works with any programming language
- ✅ **Stack Templates** - Store entire project structures
- ✅ **Smart Search** - Find snippets and stacks by name
- ✅ **Zero Config** - Works out of the box
- ✅ **Cross-Platform** - Windows, Linux, macOS
- ✅ **Lightweight** - Single binary, no dependencies
- ✅ **Secure** - SHA256 checksum verification on install

---

## Documentation

**Full docs:** [boiler.iamabhinav.dev](https://boiler.iamabhinav.dev)



## License

MIT © [Abhinav Prakash](https://github.com/rishiyaduwanshi)

---

<div align="center">

**Made with ❤️ for developers who hate repeating themselves**

[⭐ Star on GitHub](https://github.com/rishiyaduwanshi/boiler) • [📖 Read the Docs](https://boiler.iamabhinav.dev)

</div>
