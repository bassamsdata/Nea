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
nvm install nightly

# Install latest stable version
nvm install stable

# Install specific stable version
nvm install 0.11.0
```

### Use

Switch between installed versions:

```bash
# Use the latest nightly version
nvm use nightly

# Use the latest stable version
nvm use stable

# Use a specific stable version
nvm use 0.11.0
```

### Rollback

Return to a previous nightly version:

```bash
# Roll back to an earlier nightly version (e.g., 3 versions back)
nvm rollback 3
```

### List

List available Neovim versions:

```bash
# List locally installed versions
nvm ls local       # Shows all stable and up to 7 most recent nightly versions
nvm ls local 10    # Show all stable and 10 most recent nightly versions
nvm ls local -1    # Show all stable and all nightly versions

# List remotely available versions
nvm ls remote      # Shows 7 most recent stable versions
nvm ls remote 15   # Shows 15 most recent stable versions
nvm ls remote -1   # Shows all available stable versions
```

### Clean

Remove installed versions:

```bash
# Clean the oldest nightly version
nvm clean nightly

# Clean a specific nightly version by date
nvm clean 2023-05-15

# Clean all nightly versions
nvm clean nightly all

# Clean a specific stable version
nvm clean 0.9.0

# Clean all stable versions
nvm clean stable all

# Clean all versions (stable and nightly)
nvm clean all
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
