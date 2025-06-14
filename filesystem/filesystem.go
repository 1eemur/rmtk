package filesystem

import (
	"os"
	"path/filepath"
	"strings"
)

type FileList struct {
	Files         []string
	CurrentIdx    int
	Offset        int
	MaxVisible    int
	CurrentPath   string
	SearchMode    bool
	SearchQuery   string
	FilteredFiles []string
	OriginalFiles []string
}

func NewFileList(path string) (*FileList, error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}

	files, err := listFiles(absPath)
	if err != nil {
		return nil, err
	}

	return &FileList{
		Files:         files,
		OriginalFiles: files,
		FilteredFiles: []string{},
		CurrentIdx:    0,
		Offset:        0,
		MaxVisible:    0,
		CurrentPath:   absPath,
		SearchMode:    false,
		SearchQuery:   "",
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

func (fl *FileList) PageUp() {
	half := fl.MaxVisible / 2
	if fl.CurrentIdx-half < 0 {
		fl.CurrentIdx = 0
	} else {
		fl.CurrentIdx -= half
	}
	if fl.CurrentIdx < fl.Offset {
		fl.Offset = fl.CurrentIdx
	}
}

func (fl *FileList) PageDown() {
	half := fl.MaxVisible / 2
	if fl.CurrentIdx+half >= len(fl.Files) {
		fl.CurrentIdx = len(fl.Files) - 1
	} else {
		fl.CurrentIdx += half
	}
	if fl.CurrentIdx >= fl.Offset+fl.MaxVisible {
		fl.Offset = fl.CurrentIdx - fl.MaxVisible + 1
	}
}

func (fl *FileList) MoveUp() {
	if fl.CurrentIdx > 0 {
		fl.CurrentIdx--
		if fl.CurrentIdx < fl.Offset {
			fl.Offset = fl.CurrentIdx
		}
	}
}

func (fl *FileList) MoveDown() {
	if fl.CurrentIdx < len(fl.Files)-1 {
		fl.CurrentIdx++
		if fl.CurrentIdx >= fl.Offset+fl.MaxVisible {
			fl.Offset = fl.CurrentIdx - fl.MaxVisible + 1
		}
	}
}

func (fl *FileList) GoToTop() {
	fl.CurrentIdx = 0
	fl.Offset = 0
}

func (fl *FileList) GoToBottom() {
	fl.CurrentIdx = len(fl.Files) - 1
	if fl.CurrentIdx >= fl.MaxVisible {
		fl.Offset = fl.CurrentIdx - fl.MaxVisible + 1
	} else {
		fl.Offset = 0
	}
}

func (fl *FileList) CurrentFile() string {
	if len(fl.Files) == 0 {
		return ""
	}
	return fl.Files[fl.CurrentIdx]
}

func (fl *FileList) EnterSearchMode() {
	fl.SearchMode = true
	fl.SearchQuery = ""
	// Store the original list before filtering
	fl.OriginalFiles = fl.Files
	fl.FilteredFiles = fl.Files
}

func (fl *FileList) ExitSearchMode() {
	fl.SearchMode = false
	fl.SearchQuery = ""
	fl.Files = fl.OriginalFiles
	fl.CurrentIdx = 0
	fl.Offset = 0
}

func (fl *FileList) UpdateSearch(r rune) {
	if r == 8 { // Backspace
		if len(fl.SearchQuery) > 0 {
			fl.SearchQuery = fl.SearchQuery[:len(fl.SearchQuery)-1]
		}
	} else {
		fl.SearchQuery += string(r)
	}

	// Filter files based on the query
	fl.filterFiles()

	// Reset cursor position
	fl.CurrentIdx = 0
	fl.Offset = 0
}

func (fl *FileList) filterFiles() {
	if fl.SearchQuery == "" {
		fl.Files = fl.OriginalFiles
		return
	}

	fl.FilteredFiles = []string{}
	query := strings.ToLower(fl.SearchQuery)

	for _, file := range fl.OriginalFiles {
		if strings.Contains(strings.ToLower(file), query) {
			fl.FilteredFiles = append(fl.FilteredFiles, file)
		}
	}

	fl.Files = fl.FilteredFiles
}

func (fl *FileList) RemoveLastSearchChar() {
	if len(fl.SearchQuery) > 0 {
		fl.SearchQuery = fl.SearchQuery[:len(fl.SearchQuery)-1]
		fl.filterFiles()
	} else {
		fl.ExitSearchMode()
	}
}

// NavigateToPath navigates to a new directory and returns a new FileList
func (fl *FileList) NavigateToPath(newPath string) (*FileList, error) {
	return NewFileList(newPath)
}

// GetFullPath returns the full path of the currently selected file
func (fl *FileList) GetFullPath() string {
	if len(fl.Files) == 0 {
		return ""
	}

	selected := fl.CurrentFile()
	if selected == ".." {
		return filepath.Dir(fl.CurrentPath)
	}
	return filepath.Join(fl.CurrentPath, selected)
}

// IsDirectory checks if the currently selected item is a directory
func (fl *FileList) IsDirectory() bool {
	fullPath := fl.GetFullPath()
	if fullPath == "" {
		return false
	}

	fileInfo, err := os.Stat(fullPath)
	return err == nil && fileInfo.IsDir()
}
