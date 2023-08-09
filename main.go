package main

import (
	"fmt"
	"image/color"
	"strconv"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

var playerTopShips = binding.BindStringList(
	&[]string{},
)
var currentUser int

var isWorking = false

func InputWidgetWatcher(entry *widget.Entry, button *widget.Button) {
	oldIsWorking := isWorking
	for {
		if oldIsWorking != isWorking {
			if isWorking {
				entry.Disable()
				button.Disable()
			} else {
				entry.Enable()
				button.Enable()
			}
			oldIsWorking = isWorking
		}
		time.Sleep(10 * time.Millisecond)
	}
}

var objectsToHide = []*canvas.Text{}

func HideObjects() {
	for _, o := range objectsToHide {
		o.Hide()
	}
}

func UpdateDetailInfo(newData [][]string) {
	HideObjects()
	objectsToHide = []*canvas.Text{}

	// Reset the table.
	gDetailInfo.Length = func() (int, int) {
		if len(newData) == 0 {
			return 0, 0
		} else {
			return len(newData), len(newData[0])
		}
	}
	gDetailInfo.CreateCell = func() fyne.CanvasObject {
		o := canvas.NewText("template", color.White)
		objectsToHide = append(objectsToHide, o)
		return o
	}
	gDetailInfo.UpdateCell = func(i widget.TableCellID, o fyne.CanvasObject) {
		o.(*canvas.Text).Text = newData[i.Row][i.Col]
		o.(*canvas.Text).Alignment = fyne.TextAlignCenter
		if i.Row == 0 {
			o.(*canvas.Text).TextStyle = fyne.TextStyle{Bold: true}
		} else {
			if newData[i.Row][i.Col] == "O" {
				o.(*canvas.Text).Color = color.RGBA{R: 255, G: 0, B: 0, A: 255}
			} else if newData[i.Row][i.Col] == "X" {
				o.(*canvas.Text).Color = color.RGBA{R: 0, G: 255, B: 0, A: 255}
			} else if _, err := strconv.Atoi(newData[i.Row][i.Col]); err == nil {
				o.(*canvas.Text).Color = color.RGBA{R: 255, G: 255, B: 0, A: 255}
			}
		}
	}
	gDetailInfo.SetColumnWidth(0, 100)
	gDetailInfo.SetColumnWidth(1, 100)

	for i := 0; i < len(newData); i++ {
		gDetailInfo.SetRowHeight(i, 30)
	}
	gDetailInfo.Refresh()
}

func main() {
	// Create a new Fyne application instance.
	a := app.New()

	// Create a new window with the title "Go Eye".
	w := a.NewWindow("Go Eye")

	// Initialize a data binding for the player name.
	playerName := binding.NewString()
	err := playerName.Set("")
	if err != nil {
		fmt.Printf("Error occurred: %v\n", err)
	}

	// Create widgets for player entry, search button, clipboard watcher switch, result list, and detail label.
	playerEntry := createPlayerEntry(playerName)
	searchButton := createSearchButton(playerName)
	resultList, detailInfo := createResultWidgets()
	gDetailInfo = detailInfo

	// Create sub-container for player entry, search button, and clipboard watcher switch.
	subContainer := createInputContainer(playerEntry, searchButton)

	// Create the main container with a horizontal split for result list and detail label.
	mainContainer := createMainContainer(subContainer, resultList, detailInfo)

	go InputWidgetWatcher(playerEntry, searchButton)

	// Set the main container as the content of the window, resize it, and show the window.
	w.SetContent(mainContainer)
	w.Resize(fyne.NewSize(800, 400))
	w.SetFixedSize(true)
	w.ShowAndRun()
}

var gDetailInfo *widget.Table

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
	playerEntry.OnSubmitted = func(s string) {
		gSearchButton.OnTapped()
	}
	return playerEntry
}

var gSearchButton *widget.Button

// createSearchButton creates a widget for the search button with the provided callback function.
func createSearchButton(playerName binding.String) *widget.Button {
	searchButton := widget.NewButton("Analyze", func() {
		isWorking = true

		playerTopShips.Set([]string{})

		playerNameString, err := playerName.Get()
		if err != nil {
			isWorking = false
			fmt.Printf("Error occurred: %v\n", err)
			return
		}

		playerID, err := ResolveNamesToCharacterIDs([]string{playerNameString})
		if err != nil {
			isWorking = false
			fmt.Printf("Error occurred: %v\n", err)
			return
		}

		if len(playerID) == 0 {
			isWorking = false
			fmt.Printf("Error occurred: %v\n", fmt.Errorf("no character found"))
			return
		}

		ships, err := GetTopShips(playerID[0])
		if err != nil {
			isWorking = false
			fmt.Printf("Error occurred: %v\n", err)
			return
		}

		shipNames, err := ResolveIdsToNames(ships)
		if err != nil {
			isWorking = false
			fmt.Printf("Error occurred: %v\n", err)
			return
		}

		err = playerTopShips.Set(shipNames)
		if err != nil {
			isWorking = false
			fmt.Printf("Error occurred: %v\n", err)
			return
		}

		currentUser = playerID[0]

		isWorking = false
	})

	gSearchButton = searchButton
	return searchButton
}

// createResultWidgets creates widgets for the result list and detail label.
func createResultWidgets() (*widget.List, *widget.Table) {
	resultList := widget.NewListWithData(
		playerTopShips,
		func() fyne.CanvasObject {
			return widget.NewLabel("template")
		},
		func(i binding.DataItem, item fyne.CanvasObject) {
			item.(*widget.Label).Bind(i.(binding.String))
		},
	)
	resultList.OnSelected = func(id int) {
		isWorking = true

		shipSelected, err := playerTopShips.Get()
		if err != nil {
			isWorking = false
			fmt.Printf("Error occurred: %v\n", err)
			return
		}

		shipID, err := ResolveItemNamesToIDs([]string{shipSelected[id]})
		if err != nil {
			isWorking = false
			fmt.Printf("Error occurred: %v\n", err)
			return
		}

		kms, err := GetRecentLosses(currentUser, shipID[0])
		if err != nil {
			isWorking = false
			fmt.Printf("Error occurred: %v\n", err)
			return
		}

		newData := [][]string{
			{"Date", "Prop", "Scram", "Point", "Web", "Neut", "Damp"},
		}
		for _, km := range kms {
			items, t, err := GetItemsFromKillmail(km.KillmailID, km.ZKB.Hash)
			if err != nil {
				fmt.Printf("Error occurred: %v\n", err)
				continue
			}

			var prop, scram, point, web, neut, damp int

			itemsInString, err := ResolveIdsToNames(items)
			if err != nil {
				fmt.Printf("Error occurred: %v\n", err)
				continue
			}

			var date string
			if time.Now().Sub(t).Hours() < 24*31 {
				date = fmt.Sprintf("%v days ago", int(time.Now().Sub(t).Hours()/24))
			} else {
				date = fmt.Sprintf("%v", t.Format("2006-01-02"))
			}

			for _, item := range itemsInString {
				item = strings.ToLower(item)
				if strings.Contains(item, "1mn") {
					prop = 1
				} else if strings.Contains(item, "5mn") {
					prop = 5
				} else if strings.Contains(item, "10mn") {
					prop = 10
				} else if strings.Contains(item, "50mn") {
					prop = 50
				} else if strings.Contains(item, "100mn") {
					prop = 100
				} else if strings.Contains(item, "500mn") {
					prop = 500
				}

				if strings.Contains(item, "warp scrambler") {
					scram += 1
				}

				if strings.Contains(item, "warp disruptor") {
					point += 1
				}

				if strings.Contains(item, "stasis web") {
					web += 1
				}

				if strings.Contains(item, "energy neutralizer") {
					neut += 1
				}

				if strings.Contains(item, "sensor dampener") {
					damp += 1
				}
			}

			line := []string{
				date,
			}
			if prop > 0 {
				if prop == 1 || prop == 10 || prop == 100 {
					line = append(line, fmt.Sprintf("%vMN AB", prop))
				} else {
					line = append(line, fmt.Sprintf("%vMN MWD", prop))
				}
			} else {
				line = append(line, "X")
			}

			if scram > 0 {
				line = append(line, "O")
			} else {
				line = append(line, "X")
			}

			if point > 0 {
				line = append(line, "O")
			} else {
				line = append(line, "X")
			}

			if web > 0 {
				line = append(line, strconv.Itoa(web))
			} else {
				line = append(line, "X")
			}

			if neut > 0 {
				line = append(line, "O")
			} else {
				line = append(line, "X")
			}

			if damp > 0 {
				line = append(line, "O")
			} else {
				line = append(line, "X")
			}

			newData = append(newData, line)
		}

		UpdateDetailInfo(newData)
		isWorking = false
	}

	detailInfo := widget.NewTable(
		func() (int, int) {
			return 0, 0
		},
		func() fyne.CanvasObject {
			return canvas.NewText("template", color.White)
		},
		func(id widget.TableCellID, object fyne.CanvasObject) {
			object.(*canvas.Text).Text = "template"
			object.(*canvas.Text).Alignment = fyne.TextAlignCenter
		},
	)

	return resultList, detailInfo
}

// createInputContainer creates a container for player entry, search button, and clipboard watcher switch.
func createInputContainer(playerEntry *widget.Entry, searchButton *widget.Button) *fyne.Container {
	miscContainer := container.NewHBox(searchButton)
	return container.New(
		layout.NewBorderLayout(nil, nil, nil, miscContainer),
		playerEntry,
		miscContainer,
	)
}

// createMainContainer creates the main container with a horizontal split for result list and detail label.
func createMainContainer(subContainer *fyne.Container, resultList *widget.List, detailInfo *widget.Table) *fyne.Container {
	resultContainer := container.NewHSplit(resultList, detailInfo)
	resultContainer.SetOffset(0.3)

	return container.New(layout.NewBorderLayout(subContainer, nil, nil, nil), subContainer, resultContainer)
}
