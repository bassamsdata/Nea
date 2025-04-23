<h1 align="center">🕊️ NEΛ </h1>
<p align="center"><em>The Neovim version manager</em></p>

<hr>

NEΛ is a command-line tool for managing multiple Neovim versions on Linux and macOS. Install, switch between stable releases and nightly builds, and easily roll back to previous versions.

## Features

- Install and switch between stable Neovim releases
- Track and install nightly Neovim builds
- Roll back to previously installed nightly versions
- List local and remote available versions
- Clean up old installations to save disk space

## Installation

```bash
# Clone the repository
git clone https://github.com/bassamsdata/nea.git
cd nea

# Build the binary
go build -o nea

# Move to a directory in your PATH
# if you're on MacOs, no need for sudo
sudo mv nea /usr/local/bin/
```

## Commands

### Install

Install Neovim versions:

```bash
# Install latest nightly build
nea install nightly

# Install latest stable version
nea install stable

# Install specific stable version
nea install 0.11.0
```

### Use

Switch between installed versions:

```bash
# Use the latest nightly version
nea use nightly

# Use the latest stable version
nea use stable

# Use a specific stable version
nea use 0.11.0
```

### Rollback

Return to a previous nightly version:

```bash
# Roll back to an earlier nightly version (e.g., 3 versions back)
nea rollback 3
```

### List

List available Neovim versions:

```bash
# List locally installed versions
nea ls local       # Shows all stable and up to 7 most recent nightly versions
nea ls local 10    # Show all stable and 10 most recent nightly versions
nea ls local -1    # Show all stable and all nightly versions

# List remotely available versions
nea ls remote      # Shows 7 most recent stable versions
nea ls remote 15   # Shows 15 most recent stable versions
nea ls remote -1   # Shows all available stable versions
```

### Clean

Remove installed versions:

```bash
# Clean the oldest nightly version
nea clean nightly

# Clean a specific nightly version by date
nea clean 2023-05-15

# Clean all nightly versions
nea clean nightly all

# Clean a specific stable version
nea clean 0.9.0

# Clean all stable versions
nea clean stable all

# Clean all versions (stable and nightly)
nea clean all
```

## Directory Structure

NeoVMan stores configurations and Neovim versions in the following locations:

```
~/.local/share/neoManager/
├── bin/           # Contains the symlink to the active Neovim version
├── nightly/       # Contains nightly versions and version tracking info
│   └── versions_info.json
└── stable/        # Contains stable versions organized by version number
    ├── 0.8.0/
    ├── 0.9.0/
    └── ...
```

## Version Tracking

Nightly versions are tracked with a unique identifier and creation date, allowing you to roll back to previous versions if needed.

# About the Name

Nea means “new things”
