//go:generate goversioninfo
package main

import (
	"fmt"
	"image/color"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/widget"
)

func logDebug(a ...interface{}) {
	if config.Debug {
		fmt.Println(a...)
	}
}

func scaleString(c fyne.Canvas) string {
	return fmt.Sprintf("%0.2f", c.Scale())
}

func logBlueLn(g *widget.Box, s string) {
	t := canvas.NewText(s, color.RGBA{
		R: 88,
		G: 127,
		B: 252,
		A: 255,
	})
	g.Append(t)
	adjust := logs.MinSize().Height
	logScroll.Offset = fyne.NewPos(0, adjust)
	logScroll.Refresh()
}

func logRedLn(g *widget.Box, s string) {
	t := canvas.NewText(s, color.RGBA{
		R: 255,
		G: 0,
		B: 0,
		A: 255,
	})
	g.Append(t)
	adjust := logs.MinSize().Height
	logScroll.Offset = fyne.NewPos(0, adjust)
	logScroll.Refresh()
}

func logYellowLn(g *widget.Box, s string) {
	t := canvas.NewText(s, color.RGBA{
		R: 255,
		G: 255,
		B: 0,
		A: 255,
	})
	g.Append(t)
	adjust := logs.MinSize().Height
	logScroll.Offset = fyne.NewPos(0, adjust)
	logScroll.Refresh()
}

func logCyanLn(g *widget.Box, s string) {
	t := canvas.NewText(s, color.RGBA{
		R: 0,
		G: 255,
		B: 255,
		A: 255,
	})
	g.Append(t)
	adjust := logs.MinSize().Height
	logScroll.Offset = fyne.NewPos(0, adjust)
	logScroll.Refresh()
}

func logGreenLn(g *widget.Box, s string) {
	t := canvas.NewText(s, color.RGBA{
		R: 0,
		G: 255,
		B: 0,
		A: 255,
	})
	g.Append(t)
	adjust := logs.MinSize().Height
	logScroll.Offset = fyne.NewPos(0, adjust)
	logScroll.Refresh()
}

func logMagentaLn(g *widget.Box, s string) {
	t := canvas.NewText(s, color.RGBA{
		R: 255,
		G: 0,
		B: 255,
		A: 255,
	})
	g.Append(t)
	adjust := logs.MinSize().Height
	logScroll.Offset = fyne.NewPos(0, adjust)
	logScroll.Refresh()
}

func greenTXT(s string) *canvas.Text {
	return canvas.NewText(s, color.RGBA{
		R: 0,
		G: 255,
		B: 0,
		A: 255,
	})
}

func blueTXT(s string) *canvas.Text {
	return canvas.NewText(s, color.RGBA{
		R: 255,
		G: 183,
		B: 50,
		A: 255,
	})
}
