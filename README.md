![Banner](https://s-christy.com/sbs/status-banner.svg?icon=action/autorenew&hue=200&title=Recompile&description=A%20file%20watcher%20that%20automatically%20runs%20your%20build%20command%20on%20changes)

## Overview

Recompile is a lightweight file watcher that monitors your project directory
and automatically runs a build command whenever a file changes. It is useful
for fast feedback loops during development.

## Features

- Watches the current directory recursively for file changes
- Runs a configurable build command on write, create, or remove events
- Runs an initial build on startup
- Ignores `.git` and `.cache` directories by default
- Supports additional ignore patterns via the `-ignore` flag
- Ignores editor temp files (files ending in `~`)
- Colored output: red for errors, grey for status messages
- Exits cleanly on `SIGINT` or `SIGTERM`

## Setup

```
go build .
```

## Usage

```
recompile [-command <cmd>] [-ignore <dir>]
```

| Flag | Default | Description |
|---|---|---|
| `-command` | `go build ./...` | Command to run on file changes |
| `-ignore` | | Additional file or directory name to ignore (can be repeated) |

## Examples

Watch a Go project and run tests on every change:

```
recompile -command "go test ./..."
```

Watch a C project using make, ignoring the build output directory:

```
recompile -command "make" -ignore build
```

Run an arbitrary command and ignore multiple directories:

```
recompile -command "make run" -ignore build -ignore dist
```

## License

This work is licensed under the GNU General Public License version 3 (GPLv3).

[<img src="https://s-christy.com/status-banner-service/GPLv3_Logo.svg" width="150" />](https://www.gnu.org/licenses/gpl-3.0.en.html)
