package external

import (
	"os/exec"
	"path/filepath"
	"strings"
)

var zathuraFormats = []string{".pdf", ".djvu", ".ps", ".epub", ".cb", ".cbz", ".cbr"}

// CanOpenWithZathura checks if a file can be opened with Zathura
func CanOpenWithZathura(filePath string) bool {
	ext := strings.ToLower(filepath.Ext(filePath))

	for _, format := range zathuraFormats {
		if ext == format {
			return true
		}
	}
	return false
}

// OpenWithZathura opens a file with Zathura, using devour if available
func OpenWithZathura(filePath string) error {
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
