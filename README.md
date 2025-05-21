# rmtk - A lightweight CLI e-book browser

Raamatukogu: A simple and minimalistic terminal-based e-book browser. Written in Go and using Zathura to display the ebooks.

## Features

- **Lightweight Terminal UI**: Navigate your documents without leaving the terminal
- **Vim-style Navigation**: Use familiar Vim keybindings to move through files
- **Terminal Browsing**: Works with the 'devour' utility to replace your terminal with Zathura
- **Search Functionality**: Filter files by name in real-time as you type

## Requirements

- Go programming environment (for building)
- [Zathura](https://pwmt.org/projects/zathura/) document viewer (and the revelant extensions for the file types you're planning to read)
- [termbox-go](https://github.com/nsf/termbox-go) library
- [devour](https://github.com/salman-abedin/devour) (optional, for swallowing terminals)

## Installation

1. Clone the repository:
   ```
   git clone 
   cd rmtk
   ```

2. Build the application:
   ```
   go build
   ```

3. Install to your path (optional):
   ```
   cp rmtk /usr/local/bin/
   ```

## Usage

### Basic Usage

```
rmtk [directory]
```

If no directory is specified, RMTK will open in the current directory.

### Navigation

| Key           | Action                          |
|---------------|----------------------------------|
| ↑ or k        | Move up one file                |
| ↓ or j        | Move down one file              |
| Enter         | Open directory or file          |
| gg            | Go to first file                |
| G             | Go to last file                 |
| Ctrl+U        | Move up half a page             |
| Ctrl+D        | Move down half a page           |
| /             | Enter search mode               |
| Esc           | Exit search mode                |
| q             | Quit RMTK                       |

### Search Mode

1. Press `/` to enter search mode
2. Type your search query
3. Files will be filtered in real-time as you type
4. Press Enter to accept the filtered results or Esc to cancel and return to the full list

## Supported File Types

RMTK can open the following file types with Zathura:

- PDF (`.pdf`)
- DJVU (`.djvu`)
- PostScript (`.ps`)
- EPUB (`.epub`)
- Comic Books (`.cb`, `.cbz`, `.cbr`)

## Configuration

Currently, RMTK doesn't use a configuration file. All preferences must be set by modifying the source code.

## Future Plans

- Add custom configuration file support
- Add bookmarks for frequently used directories
- Add functionality to sort directories