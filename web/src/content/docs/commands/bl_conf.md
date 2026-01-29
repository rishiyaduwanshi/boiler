---
title: bl conf
description: Command reference for bl conf
---

Manage boiler configuration

### Synopsis

View and manage Boiler configuration.

You can:
  - View current configuration (default)
  - Edit config in default editor (use -e or --edit)
  - Reset to defaults (use -r or --reset)

Configuration includes paths, preferences, and behavior settings.

```
bl conf [flags]
```

### Examples

```
  # Show configuration
  bl conf

  # Edit configuration
  bl conf --edit

  # Reset to defaults
  bl conf --reset
```

### Options

```
  -e, --edit    Edit configuration
  -h, --help    help for conf
  -r, --reset   Reset configuration to defaults
  -s, --show    Show configuration
```

