package main

import (
	"fmt"
	"os"

	"rmtk/config"
	"rmtk/filesystem"
	"rmtk/input"
	"rmtk/ui"

	"github.com/nsf/termbox-go"
)

func main() {
	// Get directory path from args, config file, or use current directory
	var path string

	if len(os.Args) > 1 {
		// Use CLI argument if provided
		path = os.Args[1]
	} else {
		// Try to get path from config file
		configPath := config.GetDefaultPath()
		if configPath != "" {
			path = configPath
		} else {
			// Fall back to current directory
			path = "."
		}
	}

	// Initialize terminal UI
	err := termbox.Init()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing termbox: %v\n", err)
		os.Exit(1)
	}
	defer termbox.Close()

	// Create file list
	fileList, err := filesystem.NewFileList(path)
	if err != nil {
		termbox.Close()
		fmt.Fprintf(os.Stderr, "Error creating file list: %v\n", err)
		os.Exit(1)
	}

	// Set terminal mode for input handling and complete first render
	termbox.SetInputMode(termbox.InputEsc)
	ui.Render(fileList)

	// Initialize input handler
	inputHandler := input.NewHandler()

	// Main event loop
	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			shouldExit, newFileList, err := inputHandler.HandleKey(ev, fileList)
			if err != nil {
				termbox.Close()
				fmt.Fprintf(os.Stderr, "Error handling input: %v\n", err)
				return
			}
			if shouldExit {
				return
			}
			fileList = newFileList
		case termbox.EventResize:
			// Handle terminal resize
		}

		ui.Render(fileList)
	}
}
