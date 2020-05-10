package main

import (
	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"

	"image/color"
)

func makeCell() fyne.CanvasObject {
	rect := canvas.NewRectangle(&color.RGBA{0, 0, 0, 255})
	rect.SetMinSize(fyne.NewSize(2, 2))
	return rect
}

func startUI() {

	appli = app.New()
	w := appli.NewWindow("Pokecatcher v2.5.3")
	v, _ := box.Find("icon\\icons.png")
	windowIcon = fyne.NewStaticResource("pokecatcher.png", v)
	w.SetIcon(windowIcon)

	labellog := widget.NewLabel("logs :")

	logs = widget.NewVBox()
	logScroll = widget.NewVScrollContainer(logs)
	progressBar = widget.NewProgressBar()
	progressBar.Resize(fyne.NewSize(200, 20))
	progressBar.Min, progressBar.Max = 0, 1
	progressBar.SetValue(0)

	labelimg := widget.NewLabel("Last Pokemon Spawn :")
	c, _ := box.FindString("icon\\nothing.png")
	img, _ := loadImg(c)
	lastPokemonImg = canvas.NewImageFromImage(img)
	lastPokemonLabel = widget.NewLabelWithStyle("", fyne.TextAlignCenter, fyne.TextStyle{})
	currentLabel := widget.NewLabel("Current Pokemon :")
	currentPokemonImg = canvas.NewImageFromImage(img)
	currentPokemonLevel = widget.NewLabelWithStyle("Level : Not set", fyne.TextAlignCenter, fyne.TextStyle{})
	currentPokemonImg.FillMode = 2
	lastPokemonImg.FillMode = 2

	bottom1 := makeCell()
	left1 := makeCell()
	right1 := makeCell()
	bottom2 := makeCell()
	left2 := makeCell()
	right2 := makeCell()
	bottom3 := makeCell()
	left3 := makeCell()
	right3 := makeCell()

	w.SetContent(widget.NewHBox(
		widget.NewHBox(
			fyne.NewContainerWithLayout(layout.NewBorderLayout(labellog, bottom1, left1, right1),
				labellog,
				bottom1,
				left1,
				right1,
				logScroll,
			),
		),
		widget.NewHBox(
			fyne.NewContainerWithLayout(layout.NewGridLayout(1),
				widget.NewVBox(
					widget.NewLabel("Refresh Pokemon List :"),
					progressBar,
					fyne.NewContainerWithLayout(layout.NewBorderLayout(currentLabel, bottom3, left3, right3),
						currentLabel,
						bottom3,
						left3,
						right3,
						fyne.NewContainerWithLayout(layout.NewGridLayout(1),
							widget.NewVBox(currentPokemonImg, currentPokemonLevel)),
					)),
				widget.NewVBox(
					widget.NewButton("Clear logs", func() {
						logs.Children = []fyne.CanvasObject{}
						logs.Refresh()
						logBlueLn(logs, "The console has been cleared successfully !")
					}),
					widget.NewButton("Catch latest", func() {
						catchLatest()
					}),
					fyne.NewContainerWithLayout(layout.NewBorderLayout(labelimg, bottom2, left2, right2),
						labelimg,
						bottom2,
						left2,
						right2,
						fyne.NewContainerWithLayout(layout.NewGridLayout(1),
							widget.NewVBox(lastPokemonImg, lastPokemonLabel)),
					)),
			))))
	w.Resize(fyne.NewSize(500, 500))
	w.SetFixedSize(true)
	go usefulVariables()
	w.ShowAndRun()

}
