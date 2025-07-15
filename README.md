# Fudge Config Manager

A Terminal User Interface (TUI) application built with Go and BubbleTea for managing configuration files.

## Features

- **Open config**: Opens the `config.json` file in your default editor (`$EDITOR`)
- **Add a Sync file**: Add JSON file paths to synchronize with your main config
- **Sync config**: Apply configuration from `config.json` to all sync files

## Installation

```bash
go build -o fudge-configs .
```

## Usage

Run the application:

```bash
./fudge-configs
```

### Navigation

- Use **arrow keys** or **j/k** to navigate the menu
- Press **Enter** or **Space** to select an option
- Press **Ctrl+C** or **q** to quit
- Press **Esc** to go back to the main menu (when applicable)

### Configuration Files

The application stores configuration files in:

- **macOS/Linux**: `$XDG_CONFIG_HOME/fudge_configs/` (defaults to `~/.config/fudge_configs/`)

#### config.json

This is your main configuration file that will be applied to sync files. If it doesn't exist, it will be created with a default structure:

```json
{
  "servers": {}
}
```

#### sync.json

This file contains paths to JSON files that should be synchronized with your main config:

```json
{
  "paths": ["/path/to/file1.json", "/path/to/file2.json"]
}
```

### How Syncing Works

When you run "Sync config":

1. The application reads your `config.json`
2. For each file listed in `sync.json`, it:
   - Reads the target JSON file
   - Merges the config data (overwriting existing keys)
   - Preserves any keys in the target file that aren't in the config
   - Writes the updated data back to the file

### Example Workflow

1. **Set up your main config**: Select "Open config" to edit `config.json` with your preferred editor
2. **Add sync files**: Select "Add a Sync file" and enter the full path to JSON files you want to sync
3. **Sync configurations**: Select "Sync config" to apply your main config to all sync files

## Requirements

- Go 1.21 or later
- A configured `$EDITOR` environment variable (defaults to `vi` if not set)

## Dependencies

- [BubbleTea](https://github.com/charmbracelet/bubbletea) - TUI framework
- [Lipgloss](https://github.com/charmbracelet/lipgloss) - Terminal styling
