package main

import (
	"strings"

	"fyne.io/fyne"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"
)

//askUserForToken will display a window to the user prompting them to enter their token
func askUserForToken() {
	w := fyne.CurrentApp().NewWindow("Enter your token")
	w.SetIcon(windowIcon)

	//Where the user types/paste his token
	textInput := widget.NewEntry()
	textInput.SetPlaceHolder("Type in your token")
	textInput.OnChanged = func(text string) {
		config.Token = strings.ReplaceAll(text, "\"", "")
		config.Token = strings.ReplaceAll(config.Token, "(", "")
		config.Token = strings.ReplaceAll(config.Token, ")", "")
	}
	//Used to center the input
	empty := widget.NewLabel("")

	btn := widget.NewButton("Done", func() {
		go discordLogin()
		go saveConfig()
		w.Close()
	})

	w.SetContent(fyne.NewContainerWithLayout(layout.NewGridLayout(1), empty, textInput, btn))

	w.Resize(fyne.NewSize(1000, 100))
	w.Show()
}

//askUserForAIToken will display a window to the user prompting them to enter their token for the AI API
func askUserForAIToken() {
	w := fyne.CurrentApp().NewWindow("Enter your AIPokedex Token")
	w.SetIcon(windowIcon)

	//Where the user types/paste his token
	textInput := widget.NewEntry()
	textInput.SetPlaceHolder("Type in your AIPokedex token")
	textInput.OnChanged = func(text string) {
		config.AIToken = strings.ReplaceAll(text, "\"", "")
		config.AIToken = strings.ReplaceAll(config.AIToken, "(", "")
		config.AIToken = strings.ReplaceAll(config.AIToken, ")", "")
	}
	//Used to center the input
	empty := widget.NewLabel("")

	btn := widget.NewHBox(
		widget.NewButton("Done", func() {
			go saveConfig()
			w.Close()
		}),
		widget.NewButton("Cancel", func() {
			w.Close()
		}),
	)

	w.SetContent(fyne.NewContainerWithLayout(layout.NewGridLayout(1), textInput, fyne.NewContainerWithLayout(layout.NewGridLayout(3), empty, btn, empty)))

	w.Resize(fyne.NewSize(1000, 100))
	w.Show()
}
