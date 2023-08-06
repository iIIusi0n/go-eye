package main

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func main() {
	// Create a new Fyne application instance.
	a := app.New()

	// Create a new window with the title "Go Eye".
	w := a.NewWindow("Go Eye")

	// Initialize a data binding for the player name.
	playerName := binding.NewString()
	playerName.Set("")

	// Initialize a data binding for the clipboard watcher switch state.
	clipboardWatcher := binding.NewBool()
	clipboardWatcher.Set(false)

	// Create widgets for player entry, search button, clipboard watcher switch, result list, and detail label.
	playerEntry := createPlayerEntry(playerName)
	searchButton := createSearchButton(playerName)
	clipboardWatcherSwitch := createClipboardWatcherSwitch(playerEntry, searchButton, clipboardWatcher)
	resultList, detailLabel := createResultWidgets()

	// Create sub-container for player entry, search button, and clipboard watcher switch.
	subContainer := createInputContainer(playerEntry, searchButton, clipboardWatcherSwitch)

	// Create the main container with a horizontal split for result list and detail label.
	mainContainer := createMainContainer(subContainer, resultList, detailLabel)

	// Set the main container as the content of the window, resize it, and show the window.
	w.SetContent(mainContainer)
	w.Resize(fyne.NewSize(800, 600))
	w.ShowAndRun()
}

// createPlayerEntry creates a widget for player name entry and binds it to the provided data binding.
func createPlayerEntry(playerName binding.String) *widget.Entry {
	playerEntry := widget.NewEntry()
	playerEntry.Bind(playerName)
	playerEntry.Validator = func(s string) error {
		if len(s) < 3 {
			return fmt.Errorf("name must be at least 3 characters")
		} else if len(s) > 37 {
			return fmt.Errorf("name must be at most 37 characters")
		}
		return nil
	}
	playerEntry.SetPlaceHolder("Enter the pilot's name...")
	return playerEntry
}

// createSearchButton creates a widget for the search button with the provided callback function.
func createSearchButton(playerName binding.String) *widget.Button {
	searchButton := widget.NewButton("Analyze", func() {
		name, _ := playerName.Get()
		fmt.Printf("Searching for %s\n", name)
	})
	return searchButton
}

// createClipboardWatcherSwitch creates a widget for the clipboard watcher switch and handles its state changes.
func createClipboardWatcherSwitch(playerEntry *widget.Entry, searchButton *widget.Button, clipboardWatcher binding.Bool) *widget.Check {
	clipboardWatcherSwitch := widget.NewCheck("Capture Clipboard", func(checked bool) {
		clipboardWatcher.Set(checked)
		if checked {
			playerEntry.Disable()
			searchButton.Disable()
		} else {
			playerEntry.Enable()
			searchButton.Enable()
		}
	})
	return clipboardWatcherSwitch
}

// createResultWidgets creates widgets for the result list and detail label.
func createResultWidgets() (*widget.List, *widget.Label) {
	resultList := widget.NewList(
		func() int {
			return len(ships)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("Ships")
		},
		func(id widget.ListItemID, item fyne.CanvasObject) {
			item.(*widget.Label).SetText(ships[id])
		},
	)

	detailLabel := widget.NewLabel("")
	resultList.OnSelected = func(id widget.ListItemID) {
		detailLabel.SetText(fakeMsg)
	}

	return resultList, detailLabel
}

// createInputContainer creates a container for player entry, search button, and clipboard watcher switch.
func createInputContainer(playerEntry *widget.Entry, searchButton *widget.Button, clipboardWatcherSwitch *widget.Check) *fyne.Container {
	miscContainer := container.NewHBox(searchButton, clipboardWatcherSwitch)
	return container.New(
		layout.NewBorderLayout(nil, nil, nil, miscContainer),
		playerEntry,
		miscContainer,
	)
}

// createMainContainer creates the main container with a horizontal split for result list and detail label.
func createMainContainer(subContainer *fyne.Container, resultList *widget.List, detailLabel *widget.Label) *fyne.Container {
	resultContainer := container.NewHSplit(resultList, detailLabel)
	resultContainer.SetOffset(0.3)

	return container.New(layout.NewBorderLayout(subContainer, nil, nil, nil), subContainer, resultContainer)
}
