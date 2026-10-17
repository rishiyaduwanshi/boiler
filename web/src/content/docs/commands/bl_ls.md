---
title: bl ls
description: Command reference for bl ls
---

List snippets or stacks

### Synopsis

List all stored snippets and stacks with their version numbers.

By default, shows both snippets and stacks. Use flags to filter by type.
All resources are shown with version numbers included.

```
bl ls [flags]
```

### Examples

```
  # List everything
  bl ls

  # List only snippets
  bl ls --snippets

  # List only stacks
  bl ls --stacks
```

### Options

```
  -h, --help       help for ls
  -n, --snippets   List snippets
  -k, --stacks     List stacks

# global flags (inherited)
      --global     Force global scope for this command
      --local      Force local scope for this command
      --no-color   Disable colors and animated progress (also honors NO_COLOR)
  -V, --verbose    Enable verbose debug output
```

