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

func UI() {

	App = app.New()
	w := App.NewWindow("Pokecatcher v2.4.0")
	v, _ := box.Find("icon\\icons.png")
	icon := fyne.NewStaticResource("pokecatcher.png", v)
	w.SetIcon(icon)

	labellog := widget.NewLabel("Logs :")

	Logs = widget.NewVBox()
	LogScroll = widget.NewVScrollContainer(Logs)
	ProgressBar = widget.NewProgressBar()
	ProgressBar.Resize(fyne.NewSize(200, 20))
	ProgressBar.Min, ProgressBar.Max = 0, 1
	ProgressBar.SetValue(0)

	labelimg := widget.NewLabel("Last Pokemon Spawn :")
	c, _ := box.FindString("icon\\nothing.png")
	img, _ := loadImg(c)
	LastPokemonImg = canvas.NewImageFromImage(img)
	LastPokemonLabel = widget.NewLabelWithStyle("", fyne.TextAlignCenter, fyne.TextStyle{})
	LastPokemonImg.FillMode = 2

	bottom1 := makeCell()
	left1 := makeCell()
	right1 := makeCell()
	bottom2 := makeCell()
	left2 := makeCell()
	right2 := makeCell()

	w.SetContent(widget.NewHBox(
		widget.NewHBox(
			fyne.NewContainerWithLayout(layout.NewBorderLayout(labellog, bottom1, left1, right1),
				labellog,
				bottom1,
				left1,
				right1,
				LogScroll,
			),
		),
		widget.NewHBox(
			fyne.NewContainerWithLayout(layout.NewGridLayout(1),
				widget.NewVBox(widget.NewLabel("Refresh Pokemon List :"), ProgressBar),
				widget.NewVBox(
					widget.NewButton("Clear Logs", func() {
						Logs.Children = []fyne.CanvasObject{}
						Logs.Refresh()
						LogBlueLn(Logs, "The console has been cleared successfully !")
					}),
					widget.NewButton("Catch latest", func() {
						CatchLatest()
					}),
					fyne.NewContainerWithLayout(layout.NewBorderLayout(labelimg, bottom2, left2, right2),
						labelimg,
						bottom2,
						left2,
						right2,
						fyne.NewContainerWithLayout(layout.NewGridLayout(1),
							widget.NewVBox(LastPokemonImg, LastPokemonLabel)),
					)),
			))))
	w.Resize(fyne.NewSize(500, 500))
	w.SetFixedSize(true)
	go UsefulVariables()
	w.ShowAndRun()

}
