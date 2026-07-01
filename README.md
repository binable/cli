<!--
SPDX-FileCopyrightText: 2026 Tim Lochmüller <tim@binable.app>

SPDX-License-Identifier: MIT
-->

# binable CLI

Command-line tool to fetch waste collection schedules via [binable.app](https://binable.app).

## Usage

```bash
binable --street <street> --house <no> --zip <zip> --city <city>
```

**Example:**

```bash
binable --street "Schürhornweg" --house 1 --zip 33649 --city Bielefeld
```

```
Provider: Stadt Bielefeld
Address:  Schürhornweg 1, 33649 Bielefeld, DE

  Biomüll     30.06.2026 (Tue)  in 5 days
  Restmüll    07.07.2026 (Tue)  in 12 days
  Wertstoff   14.07.2026 (Tue)  in 19 days
  Papiermüll  21.07.2026 (Tue)  in 26 days
```

### Options

| Short | Long        | Description               | Default |
|-------|-------------|---------------------------|---------|
| `-s`  | `--street`  | Street name               | —       |
| `-n`  | `--house`   | House number              | —       |
| `-z`  | `--zip`     | ZIP code                  | —       |
| `-c`  | `--city`    | City                      | —       |
| `-C`  | `--country` | Country code              | `DE`    |
| `-a`  | `--all`     | Show all collection dates | —       |
| `-j`  | `--json`    | Output as JSON            | —       |
| `-v`  | `--version` | Print version             | —       |

Set `NO_COLOR=1` to disable colored output.

## Installation

### Prebuilt Binary (recommended)

Pre-built binaries for Linux, macOS, and Windows are available on the [Releases](../../releases) page.

```bash
# macOS (Apple Silicon)
curl -L https://github.com/binable/cli/releases/latest/download/binable-macos-arm64 -o binable
chmod +x binable
```

### Build from Source

**Requirement:** [Go 1.21+](https://go.dev/dl/)

```bash
git clone https://github.com/binable/cli.git
cd cli
go build -o binable .
```

**With version info:**

```bash
go build -ldflags="-s -w -X main.version=v1.0.0" -o binable .
```

### Cross-Compilation

Go can target other platforms without any additional tools:

```bash
# Linux amd64
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o binable-linux-amd64 .

# Linux arm64 (e.g. Raspberry Pi)
GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -o binable-linux-arm64 .

# macOS Intel
GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build -o binable-macos-amd64 .

# macOS Apple Silicon
GOOS=darwin GOARCH=arm64 CGO_ENABLED=0 go build -o binable-macos-arm64 .

# Windows
GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -o binable-windows-amd64.exe .
```

## Release

A new release is triggered automatically via GitHub Actions whenever a tag is pushed:

```bash
git tag v1.0.0
git push origin v1.0.0
```

The workflow builds all platforms in parallel and attaches the binaries to the GitHub release.

## No External Dependencies

This tool uses only the Go standard library — no `go get`, no vendor directory, a single static binary.
