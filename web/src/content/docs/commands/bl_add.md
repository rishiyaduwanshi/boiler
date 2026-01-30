---
title: bl add
description: Command reference for bl add
---

Add a snippet or stack to current directory

### Synopsis

Add a stored snippet or stack to your current directory.

The command copies resources from your store. For snippets with a single version,
you can use just the name (e.g., 'errorHandler' will auto-select version 1).
For multiple versions, you'll be prompted to choose.

Template Variables:
  Snippets can contain template variables using the format: bl__VAR_NAME
  When adding a snippet with variables, you'll be prompted to provide values:
    - Default values are shown in brackets (from __var declarations)
    - Press Enter to use default or type a custom value
    - Variables are replaced and metadata comments are removed in the final file

Stacks are also versioned and can be added by name or with explicit version.

```
bl add [resource] [flags]
```

### Examples

```
  # Add snippet (auto-detects if single version)
  bl add errorHandler

  # Add snippet with template variables
  bl add apiClient
  # Prompts: bl__API_URL [http://localhost:3000]: https://api.example.com
  #          bl__API_KEY [your-key]: abc123xyz
  # Output: Clean file with variables replaced, no metadata comments

  # Add specific version
  bl add logger@2.js

  # Add to specific directory
  bl add config --to ./src/utils

  # Add stack
  bl add express-api@1

  # Force overwrite
  bl add middleware --force
```

### Options

```
  -b, --both        Add to both local and global
  -f, --force       Force operation without confirmation
  -g, --global      Add to global store
  -h, --help        help for add
  -r, --remote      Fetch from remote registry
  -t, --to string   Destination path (default ".")
```

