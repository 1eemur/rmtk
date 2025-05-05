package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/nsf/termbox-go"
)

type FileList struct {
	files       []string
	currentIdx  int
	offset      int
	maxVisible  int
	currentPath string
}

func newFileList(path string) (*FileList, error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}

	files, err := listFiles(absPath)
	if err != nil {
		return nil, err
	}

	return &FileList{
		files:       files,
		currentIdx:  0,
		offset:      0,
		maxVisible:  0,
		currentPath: absPath,
	}, nil
}

func listFiles(path string) ([]string, error) {
	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	var files []string
	for _, entry := range entries {
		files = append(files, entry.Name())
	}

	return files, nil
}

func (fl *FileList) moveUp() {
	if fl.currentIdx > 0 {
		fl.currentIdx--
		if fl.currentIdx < fl.offset {
			fl.offset = fl.currentIdx
		}
	}
}

func (fl *FileList) moveDown() {
	if fl.currentIdx < len(fl.files)-1 {
		fl.currentIdx++
		if fl.currentIdx >= fl.offset+fl.maxVisible {
			fl.offset = fl.currentIdx - fl.maxVisible + 1
		}
	}
}

func (fl *FileList) currentFile() string {
	if len(fl.files) == 0 {
		return ""
	}
	return fl.files[fl.currentIdx]
}

func (fl *FileList) render() {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)

	w, h := termbox.Size()
	fl.maxVisible = h - 3 // Reserve space for header and footer

	// Draw header
	header := fmt.Sprintf(" File Navigator - %s ", fl.currentPath)
	drawText(0, 0, header, termbox.ColorBlack, termbox.ColorWhite)
	drawLine(1, w, termbox.ColorWhite)

	// Draw files
	visibleEnd := fl.offset + fl.maxVisible
	if visibleEnd > len(fl.files) {
		visibleEnd = len(fl.files)
	}

	if len(fl.files) == 0 {
		drawText(2, 3, "No files in this directory", termbox.ColorRed, termbox.ColorDefault)
	} else {
		for i := fl.offset; i < visibleEnd; i++ {
			y := i - fl.offset + 2
			file := fl.files[i]

			// Check if it's a directory and add indicator
			filePath := filepath.Join(fl.currentPath, file)
			fileInfo, err := os.Stat(filePath)
			var displayName string
			if err == nil && fileInfo.IsDir() {
				displayName = file + "/"
			} else {
				displayName = file
			}

			// Highlight current selection
			if i == fl.currentIdx {
				drawText(0, y, "> "+displayName, termbox.ColorBlack, termbox.ColorWhite)
			} else {
				drawText(2, y, displayName, termbox.ColorDefault, termbox.ColorDefault)
			}
		}
	}

	// Draw footer
	footerY := h - 1
	footer := " ↑/k: Up | ↓/j: Down | Enter: Open with Zathura | q: Quit "
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

func openWithZathura(filePath string) error {
	// First check if devour is available
	_, err := exec.LookPath("devour")
	if err == nil {
		// Use devour to launch zathura (this will close the terminal)
		cmd := exec.Command("devour", "zathura", filePath)
		return cmd.Start()
	} else {
		// Fallback to regular zathura if devour is not installed
		cmd := exec.Command("zathura", filePath)
		return cmd.Start()
	}
}

func main() {
	// Get directory path from args or use current directory
	path := "."
	if len(os.Args) > 1 {
		path = os.Args[1]
	}

	// Initialize terminal UI
	err := termbox.Init()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing termbox: %v\n", err)
		os.Exit(1)
	}
	defer termbox.Close()

	// Create file list
	fileList, err := newFileList(path)
	if err != nil {
		termbox.Close()
		fmt.Fprintf(os.Stderr, "Error creating file list: %v\n", err)
		os.Exit(1)
	}

	// Set terminal mode for input handling
	termbox.SetInputMode(termbox.InputEsc)

	// First render
	fileList.render()

mainloop:
	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			switch ev.Key {
			case termbox.KeyEsc, termbox.KeyCtrlC:
				break mainloop
			case termbox.KeyArrowUp:
				fileList.moveUp()
			case termbox.KeyArrowDown:
				fileList.moveDown()
			case termbox.KeyEnter:
				if len(fileList.files) > 0 {
					selectedFile := filepath.Join(fileList.currentPath, fileList.currentFile())

					// Check if file is a directory
					fileInfo, err := os.Stat(selectedFile)
					if err == nil && fileInfo.IsDir() {
						// Navigate into directory
						fileList, err = newFileList(selectedFile)
						if err != nil {
							// Show error briefly
							termbox.Close()
							fmt.Fprintf(os.Stderr, "Error opening directory: %v\n", err)
							return
						}
					} else {
						// Only try to open files that might be compatible with Zathura
						// Common document formats that Zathura can open
						ext := strings.ToLower(filepath.Ext(selectedFile))
						zathuraFormats := []string{".pdf", ".djvu", ".ps", ".epub", ".cb", ".cbz", ".cbr"}

						canOpen := false
						for _, format := range zathuraFormats {
							if ext == format {
								canOpen = true
								break
							}
						}

						if canOpen {
							termbox.Close()
							err = openWithZathura(selectedFile)
							if err != nil {
								fmt.Fprintf(os.Stderr, "Error opening file with Zathura: %v\n", err)
							}
							return
						}
					}
				}
			default:
				// Handle char keys for vim-style navigation
				switch ev.Ch {
				case 'q':
					break mainloop
				case 'k':
					fileList.moveUp()
				case 'j':
					fileList.moveDown()
				}
			}
		case termbox.EventResize:
			// Handle terminal resize
		}

		fileList.render()
	}
}
