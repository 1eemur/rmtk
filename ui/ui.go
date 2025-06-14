package ui

import (
	"fmt"
	"os"
	"path/filepath"

	"rmtk/filesystem"

	"github.com/nsf/termbox-go"
)

// Render renders the file list to the terminal
func Render(fl *filesystem.FileList) {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	w, h := termbox.Size()
	fl.MaxVisible = h - 3 // Reserve space for header and footer

	// Draw header
	header := fmt.Sprintf(" RMTK - %s ", fl.CurrentPath)
	if fl.SearchMode {
		header = fmt.Sprintf(" RMTK - %s [Search: %s] ", fl.CurrentPath, fl.SearchQuery)
	}
	drawText(0, 0, header, termbox.ColorBlack, termbox.ColorWhite)
	drawLine(1, w, termbox.ColorWhite)

	// Draw files
	visibleEnd := fl.Offset + fl.MaxVisible
	if visibleEnd > len(fl.Files) {
		visibleEnd = len(fl.Files)
	}

	if len(fl.Files) == 0 {
		if fl.SearchMode && fl.SearchQuery != "" {
			drawText(2, 3, "No matching files", termbox.ColorRed, termbox.ColorDefault)
		} else {
			drawText(2, 3, "No files in this directory", termbox.ColorRed, termbox.ColorDefault)
		}
	} else {
		for i := fl.Offset; i < visibleEnd; i++ {
			y := i - fl.Offset + 2
			file := fl.Files[i]

			// Check if it's a directory and add indicator
			filePath := filepath.Join(fl.CurrentPath, file)
			fileInfo, err := os.Stat(filePath)
			var displayName string
			if err == nil && fileInfo.IsDir() {
				displayName = file + "/"
			} else {
				displayName = file
			}

			// Highlight current selection
			if i == fl.CurrentIdx {
				drawText(0, y, "> "+displayName, termbox.ColorBlack, termbox.ColorWhite)
			} else {
				drawText(2, y, displayName, termbox.ColorDefault, termbox.ColorDefault)
			}
		}
	}

	// Draw footer
	footerY := h - 1
	footer := " ↑/k: Up | ↓/j: Down | gg: Top | G: Bottom | Ctrl+U/D: Half Page | /: Search | Esc: Exit Search | q: Quit "
	drawText(0, footerY, footer, termbox.ColorBlack, termbox.ColorWhite)

	termbox.Flush()
}

func drawText(x, y int, text string, fg, bg termbox.Attribute) {
	for i, char := range text {
		termbox.SetCell(x+i, y, char, fg, bg)
	}
}

func drawLine(y, width int, color termbox.Attribute) {
	for i := 0; i < width; i++ {
		termbox.SetCell(i, y, '─', color, termbox.ColorDefault)
	}
}
