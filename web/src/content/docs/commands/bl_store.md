---
title: bl store
description: Command reference for bl store
---

Store a folder/file as snippet or stack

### Synopsis

Store a file as a snippet or directory as a stack in your Boiler store.

Files are stored as snippets with version numbers.
Directories must have a boiler.stack.json config file (run 'bl init' first).

Stacks require boiler.stack.json with:
  - id: Stack name
  - version: Version number
  - ignore: Patterns to exclude

If a stack version already exists, you'll be prompted to overwrite.

```
bl store [path] [flags]
```

### Examples

```
  # Store current directory as stack
  bl store

  # Store specific file as snippet
  bl store ./utils/logger.js

  # Store directory as stack
  bl store ./my-template

  # Store with custom name
  bl store ./config.js --name dbConfig.js
```

### Options

```
  -d, --description string   Description
  -h, --help                 help for store
      --name string          Name for the resource (auto-detected from path if not provided)
  -n, --snippet              Force store as snippet
  -k, --stack                Force store as stack
```

