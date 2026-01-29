---
title: bl init
description: Command reference for bl init
---

Initialize boiler in current directory

### Synopsis

Initialize Boiler by creating necessary directories and configuration files.

This sets up:
  - Store directory for snippets and stacks
  - Configuration file
  - Metadata tracking
  - Log directory

Run this once before using other Boiler commands.

```
bl init [flags]
```

### Examples

```
  # Initialize Boiler
  bl init

  # After init, you can start storing resources
  bl store ./utils/logger.js
```

### Options

```
  -g, --global   Initialize global configuration
  -h, --help     help for init
```

