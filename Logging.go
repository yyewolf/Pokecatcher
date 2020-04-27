//go:generate goversioninfo
package main

import (
	"fmt"
	"image/color"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/widget"
)

func scaleString(c fyne.Canvas) string {
	return fmt.Sprintf("%0.2f", c.Scale())
}

func LogBlueLn(g *widget.Box, s string) {
	t := canvas.NewText(s, color.RGBA{
		R: 0,
		G: 0,
		B: 255,
		A: 255,
	})
	g.Append(t)
	adjust := Logs.MinSize().Height
	LogScroll.Offset = fyne.NewPos(0, adjust)
	LogScroll.Refresh()
}

func LogRedLn(g *widget.Box, s string) {
	t := canvas.NewText(s, color.RGBA{
		R: 255,
		G: 0,
		B: 0,
		A: 255,
	})
	g.Append(t)
	adjust := Logs.MinSize().Height
	LogScroll.Offset = fyne.NewPos(0, adjust)
	LogScroll.Refresh()
}

func LogYellowLn(g *widget.Box, s string) {
	t := canvas.NewText(s, color.RGBA{
		R: 255,
		G: 255,
		B: 0,
		A: 255,
	})
	g.Append(t)
	adjust := Logs.MinSize().Height
	LogScroll.Offset = fyne.NewPos(0, adjust)
	LogScroll.Refresh()
}

func LogCyanLn(g *widget.Box, s string) {
	t := canvas.NewText(s, color.RGBA{
		R: 0,
		G: 255,
		B: 255,
		A: 255,
	})
	g.Append(t)
	adjust := Logs.MinSize().Height
	LogScroll.Offset = fyne.NewPos(0, adjust)
	LogScroll.Refresh()
}

func LogGreenLn(g *widget.Box, s string) {
	t := canvas.NewText(s, color.RGBA{
		R: 0,
		G: 255,
		B: 0,
		A: 255,
	})
	g.Append(t)
	adjust := Logs.MinSize().Height
	LogScroll.Offset = fyne.NewPos(0, adjust)
	LogScroll.Refresh()
}

func LogMagentaLn(g *widget.Box, s string) {
	t := canvas.NewText(s, color.RGBA{
		R: 255,
		G: 0,
		B: 255,
		A: 255,
	})
	g.Append(t)
	adjust := Logs.MinSize().Height
	LogScroll.Offset = fyne.NewPos(0, adjust)
	LogScroll.Refresh()
}

func GreenTXT(s string) *canvas.Text {
	return canvas.NewText(s, color.RGBA{
		R: 0,
		G: 255,
		B: 0,
		A: 255,
	})
}

func BlueTXT(s string) *canvas.Text {
	return canvas.NewText(s, color.RGBA{
		R: 0,
		G: 0,
		B: 255,
		A: 255,
	})
}
