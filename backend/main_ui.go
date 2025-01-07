package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

var language = "en"
var translations = map[string]map[string]string{
	"en": {"team_label": "Teams", "location_label": "Locations", "assign_button": "Assign"},
	"ko": {"team_label": "팀", "location_label": "위치", "assign_button": "할당"},
}

func main() {
	a := app.New()
	w := a.NewWindow("Team Tracker")

	// UI elements
	teamList := widget.NewList(func() int { return 0 }, func() fyne.CanvasObject {
		return widget.NewLabel("")
	}, func(i widget.ListItemID, o fyne.CanvasObject) {
		o.(*widget.Label).SetText("")
	})
	locationList := widget.NewList(func() int { return 0 }, func() fyne.CanvasObject {
		return widget.NewLabel("")
	}, func(i widget.ListItemID, o fyne.CanvasObject) {
		o.(*widget.Label).SetText("")
	})

	// Assign button
	assignButton := widget.NewButton(getTranslation("assign_button"), func() {
		// Logic to assign a team to a location
		fmt.Println("Assign button clicked")
	})

	// Language toggle
	toggleLanguage := widget.NewButton("Toggle Language", func() {
		if language == "en" {
			language = "ko"
		} else {
			language = "en"
		}
		updateLabels(assignButton)
	})

	// Theme toggle
	darkMode := false
	themeToggle := widget.NewButton("Toggle Dark Mode", func() {
		darkMode = !darkMode
		if darkMode {
			a.Settings().SetTheme(theme.DarkTheme())
		} else {
			a.Settings().SetTheme(theme.LightTheme())
		}
	})

	// Layout
	w.SetContent(container.NewVBox(
		widget.NewLabel("Team Tracker"),
		toggleLanguage,
		themeToggle,
		widget.NewLabel(getTranslation("team_label")),
		teamList,
		widget.NewLabel(getTranslation("location_label")),
		locationList,
		assignButton,
	))
	w.Resize(fyne.NewSize(400, 600))
	w.ShowAndRun()

	// Fetch teams and locations
	go fetchTeams(teamList)
	go fetchLocations(locationList)
}

func getTranslation(key string) string {
	if val, ok := translations[language][key]; ok {
		return val
	}
	return key
}

func updateLabels(assignButton *widget.Button) {
	assignButton.SetText(getTranslation("assign_button"))
}

func fetchTeams(list *widget.List) {
	resp, err := http.Get("http://localhost:8080/api/teams")
	if err != nil {
		fmt.Println("Error fetching teams:", err)
		return
	}
	defer resp.Body.Close()

	data, _ := ioutil.ReadAll(resp.Body)
	var teams []string
	json.Unmarshal(data, &teams)

	list.Length = func() int { return len(teams) }
	list.UpdateItem = func(i widget.ListItemID, o fyne.CanvasObject) {
		o.(*widget.Label).SetText(teams[i])
	}
	list.Refresh()
}

func fetchLocations(list *widget.List) {
	resp, err := http.Get("http://localhost:8080/api/locations")
	if err != nil {
		fmt.Println("Error fetching locations:", err)
		return
	}
	defer resp.Body.Close()

	data, _ := ioutil.ReadAll(resp.Body)
	var locations []map[string]interface{}
	json.Unmarshal(data, &locations)

	list.Length = func() int { return len(locations) }
	list.UpdateItem = func(i widget.ListItemID, o fyne.CanvasObject) {
		o.(*widget.Label).SetText(fmt.Sprintf("%v", locations[i]["name"]))
	}
	list.Refresh()
}
