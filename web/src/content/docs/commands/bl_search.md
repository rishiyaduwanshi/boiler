---
title: bl search
description: Command reference for bl search
---

Search for snippets or stacks

### Synopsis

Search for resources in your store by name.

Searches both snippets and stacks by default. Use flags to filter:
  - Use -s or --snippets to search only snippets
  - Use -k or --stacks to search only stacks

Search is case-insensitive and matches partial names.

```
bl search [query] [flags]
```

### Examples

```
  # Search for anything with 'error'
  bl search error

  # Search only snippets
  bl search logger --snippets

  # Search only stacks
  bl search express --stacks
```

### Options

```
  -h, --help       help for search
  -r, --remote     Search remote registry
  -n, --snippets   Search only snippets
  -k, --stacks     Search only stacks
```

