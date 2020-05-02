package main

import (
	"strings"
	"fyne.io/fyne"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"
)

//AskUserForToken will display a window to the user prompting them to enter their token
func AskUserForToken(){
	w := fyne.CurrentApp().NewWindow("Enter your token")
	w.SetIcon(WindowIcon)
	
	TextInput := widget.NewEntry()
	TextInput.SetPlaceHolder("Type in your token")
	TextInput.OnChanged = func(text string) {
		Config.Token = strings.ReplaceAll(text, "\"", "")
		Config.Token = strings.ReplaceAll(Config.Token, "(", "")
		Config.Token = strings.ReplaceAll(Config.Token, ")", "")
	}
	
	empty := widget.NewLabel("")
	
	Btn := widget.NewButton("Done", func() {
		go Login()
		go SaveConfig()
		w.Close()
	})
	
	
	w.SetContent(fyne.NewContainerWithLayout(layout.NewGridLayout(1), empty, TextInput, Btn))

	w.Resize(fyne.NewSize(1000, 100))
	w.Show()
}