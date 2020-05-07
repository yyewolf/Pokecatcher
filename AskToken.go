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
