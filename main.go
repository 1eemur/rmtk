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
	files         []string
	currentIdx    int
	offset        int
	maxVisible    int
	currentPath   string
	searchMode    bool
	searchQuery   string
	filteredFiles []string
	originalFiles []string
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
		files:         files,
		originalFiles: files,
		filteredFiles: []string{},
		currentIdx:    0,
		offset:        0,
		maxVisible:    0,
		currentPath:   absPath,
		searchMode:    false,
		searchQuery:   "",
	}, nil
}

func listFiles(path string) ([]string, error) {
	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	var files []string

	// Add ".." if not at root
	parent := filepath.Dir(path)
	if parent != path {
		files = append(files, "..")
	}

	for _, entry := range entries {
		files = append(files, entry.Name())
	}

	return files, nil
}

func (fl *FileList) pageUp() {
	half := fl.maxVisible / 2
	if fl.currentIdx-half < 0 {
		fl.currentIdx = 0
	} else {
		fl.currentIdx -= half
	}
	if fl.currentIdx < fl.offset {
		fl.offset = fl.currentIdx
	}
}

func (fl *FileList) pageDown() {
	half := fl.maxVisible / 2
	if fl.currentIdx+half >= len(fl.files) {
		fl.currentIdx = len(fl.files) - 1
	} else {
		fl.currentIdx += half
	}
	if fl.currentIdx >= fl.offset+fl.maxVisible {
		fl.offset = fl.currentIdx - fl.maxVisible + 1
	}
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

func (fl *FileList) goToTop() {
	fl.currentIdx = 0
	fl.offset = 0
}

func (fl *FileList) goToBottom() {
	fl.currentIdx = len(fl.files) - 1
	if fl.currentIdx >= fl.maxVisible {
		fl.offset = fl.currentIdx - fl.maxVisible + 1
	} else {
		fl.offset = 0
	}
}

func (fl *FileList) currentFile() string {
	if len(fl.files) == 0 {
		return ""
	}
	return fl.files[fl.currentIdx]
}

func (fl *FileList) enterSearchMode() {
	fl.searchMode = true
	fl.searchQuery = ""
	// Store the original list before filtering
	fl.originalFiles = fl.files
	fl.filteredFiles = fl.files
}

func (fl *FileList) exitSearchMode() {
	fl.searchMode = false
	fl.searchQuery = ""
	fl.files = fl.originalFiles
	fl.currentIdx = 0
	fl.offset = 0
}

func (fl *FileList) updateSearch(r rune) {
	if r == 8 { // Backspace
		if len(fl.searchQuery) > 0 {
			fl.searchQuery = fl.searchQuery[:len(fl.searchQuery)-1]
		}
	} else {
		fl.searchQuery += string(r)
	}

	// Filter files based on the query
	fl.filterFiles()

	// Reset cursor position
	fl.currentIdx = 0
	fl.offset = 0
}

func (fl *FileList) filterFiles() {
	if fl.searchQuery == "" {
		fl.files = fl.originalFiles
		return
	}

	fl.filteredFiles = []string{}
	query := strings.ToLower(fl.searchQuery)

	for _, file := range fl.originalFiles {
		if strings.Contains(strings.ToLower(file), query) {
			fl.filteredFiles = append(fl.filteredFiles, file)
		}
	}

	fl.files = fl.filteredFiles
}

func (fl *FileList) render() {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	w, h := termbox.Size()
	fl.maxVisible = h - 3 // Reserve space for header and footer

	// Draw header
	header := fmt.Sprintf(" RMTK - %s ", fl.currentPath)
	if fl.searchMode {
		header = fmt.Sprintf(" RMTK - %s [Search: %s] ", fl.currentPath, fl.searchQuery)
	}
	drawText(0, 0, header, termbox.ColorBlack, termbox.ColorWhite)
	drawLine(1, w, termbox.ColorWhite)

	// Draw files
	visibleEnd := fl.offset + fl.maxVisible
	if visibleEnd > len(fl.files) {
		visibleEnd = len(fl.files)
	}

	if len(fl.files) == 0 {
		if fl.searchMode && fl.searchQuery != "" {
			drawText(2, 3, "No matching files", termbox.ColorRed, termbox.ColorDefault)
		} else {
			drawText(2, 3, "No files in this directory", termbox.ColorRed, termbox.ColorDefault)
		}
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
	var gPressed bool
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

	// Set terminal mode for input handling and complete first render
	termbox.SetInputMode(termbox.InputEsc)
	fileList.render()

mainloop:
	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			if fileList.searchMode {
				switch ev.Key {
				case termbox.KeyEsc:
					fileList.exitSearchMode()
				case termbox.KeyEnter:
					// Accept search results and exit search mode
					fileList.searchMode = false
				case termbox.KeyBackspace, termbox.KeyBackspace2:
					if len(fileList.searchQuery) > 0 {
						fileList.searchQuery = fileList.searchQuery[:len(fileList.searchQuery)-1]
						fileList.filterFiles()
					} else {
						fileList.exitSearchMode()
					}
				default:
					if ev.Ch != 0 {
						fileList.updateSearch(ev.Ch)
					}
				}
			} else {
				switch ev.Key {
				case termbox.KeyEsc, termbox.KeyCtrlC:
					break mainloop
				case termbox.KeyArrowUp:
					fileList.moveUp()
				case termbox.KeyArrowDown:
					fileList.moveDown()
				case termbox.KeyCtrlD:
					fileList.pageDown()
				case termbox.KeyCtrlU:
					fileList.pageUp()
				case termbox.KeyEnter:
					if len(fileList.files) > 0 {
						selected := fileList.currentFile()
						var newPath string

						if selected == ".." {
							newPath = filepath.Dir(fileList.currentPath)
						} else {
							newPath = filepath.Join(fileList.currentPath, selected)
						}

						fileInfo, err := os.Stat(newPath)
						if err == nil && fileInfo.IsDir() {
							fileList, err = newFileList(newPath)
							if err != nil {
								termbox.Close()
								fmt.Fprintf(os.Stderr, "Error opening directory: %v\n", err)
								return
							}
						} else {
							ext := strings.ToLower(filepath.Ext(newPath))
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
								err = openWithZathura(newPath)
								if err != nil {
									fmt.Fprintf(os.Stderr, "Error opening file with Zathura: %v\n", err)
								}
								return
							}
						}
					}
				default:
					switch ev.Ch {
					case 'q':
						break mainloop
					case 'k':
						fileList.moveUp()
						gPressed = false
					case 'j':
						fileList.moveDown()
						gPressed = false
					case 'g':
						if gPressed {
							fileList.goToTop()
							gPressed = false
						} else {
							gPressed = true
						}
					case 'G':
						fileList.goToBottom()
						gPressed = false
					case '/':
						fileList.enterSearchMode()
					default:
						gPressed = false
					}
				}
			}
		case termbox.EventResize:
			// Handle terminal resize
		}

		fileList.render()
	}
}
