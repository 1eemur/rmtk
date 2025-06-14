package input

import (
	"fmt"
	"os"

	"rmtk/external"
	"rmtk/filesystem"

	"github.com/nsf/termbox-go"
)

type Handler struct {
	gPressed bool
}

func NewHandler() *Handler {
	return &Handler{
		gPressed: false,
	}
}

// HandleKey processes key events and returns (shouldExit, newFileList, error)
func (h *Handler) HandleKey(ev termbox.Event, fileList *filesystem.FileList) (bool, *filesystem.FileList, error) {
	if fileList.SearchMode {
		return h.handleSearchMode(ev, fileList)
	}
	return h.handleNormalMode(ev, fileList)
}

func (h *Handler) handleSearchMode(ev termbox.Event, fileList *filesystem.FileList) (bool, *filesystem.FileList, error) {
	switch ev.Key {
	case termbox.KeyEsc:
		fileList.ExitSearchMode()
	case termbox.KeyEnter:
		// Accept search results and exit search mode
		fileList.SearchMode = false
	case termbox.KeyBackspace, termbox.KeyBackspace2:
		fileList.RemoveLastSearchChar()
	default:
		if ev.Ch != 0 {
			fileList.UpdateSearch(ev.Ch)
		}
	}
	return false, fileList, nil
}

func (h *Handler) handleNormalMode(ev termbox.Event, fileList *filesystem.FileList) (bool, *filesystem.FileList, error) {
	switch ev.Key {
	case termbox.KeyEsc, termbox.KeyCtrlC:
		return true, fileList, nil
	case termbox.KeyArrowUp:
		fileList.MoveUp()
	case termbox.KeyArrowDown:
		fileList.MoveDown()
	case termbox.KeyCtrlD:
		fileList.PageDown()
	case termbox.KeyCtrlU:
		fileList.PageUp()
	case termbox.KeyEnter:
		return h.handleEnter(fileList)
	default:
		return h.handleCharacter(ev.Ch, fileList)
	}
	return false, fileList, nil
}

func (h *Handler) handleEnter(fileList *filesystem.FileList) (bool, *filesystem.FileList, error) {
	if len(fileList.Files) == 0 {
		return false, fileList, nil
	}

	fullPath := fileList.GetFullPath()

	if fileList.IsDirectory() {
		newFileList, err := fileList.NavigateToPath(fullPath)
		if err != nil {
			return false, fileList, fmt.Errorf("error opening directory: %v", err)
		}
		return false, newFileList, nil
	} else {
		// Check if file can be opened with Zathura
		if external.CanOpenWithZathura(fullPath) {
			termbox.Close()
			err := external.OpenWithZathura(fullPath)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error opening file with Zathura: %v\n", err)
			}
			return true, fileList, nil
		}
	}
	return false, fileList, nil
}

func (h *Handler) handleCharacter(ch rune, fileList *filesystem.FileList) (bool, *filesystem.FileList, error) {
	switch ch {
	case 'q':
		return true, fileList, nil
	case 'k':
		fileList.MoveUp()
		h.gPressed = false
	case 'j':
		fileList.MoveDown()
		h.gPressed = false
	case 'g':
		if h.gPressed {
			fileList.GoToTop()
			h.gPressed = false
		} else {
			h.gPressed = true
		}
	case 'G':
		fileList.GoToBottom()
		h.gPressed = false
	case '/':
		fileList.EnterSearchMode()
	default:
		h.gPressed = false
	}
	return false, fileList, nil
}
