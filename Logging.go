//go:generate goversioninfo
package main

import (
	"github.com/mum4k/termdash/cell"
	"github.com/mum4k/termdash/widgets/text"
)

func PrintBlueln(txt string) {
	logBox.Write(txt+"\n", text.WriteCellOpts(
		cell.FgColor(cell.ColorBlue)),
	)
}

func PrintRedln(txt string) {
	logBox.Write(txt+"\n", text.WriteCellOpts(
		cell.FgColor(cell.ColorRed)),
	)
}

func PrintYellowln(txt string) {
	logBox.Write(txt+"\n", text.WriteCellOpts(
		cell.FgColor(cell.ColorYellow)),
	)
}

func PrintCyanln(txt string) {
	logBox.Write(txt+"\n", text.WriteCellOpts(
		cell.FgColor(cell.ColorCyan)),
	)
}

func PrintGreenln(txt string) {
	logBox.Write(txt+"\n", text.WriteCellOpts(
		cell.FgColor(cell.ColorGreen)),
	)
}

func PrintMagentaln(txt string) {
	logBox.Write(txt+"\n", text.WriteCellOpts(
		cell.FgColor(cell.ColorGreen)),
	)
}

func PrintGreen(txt string) {
	logBox.Write(txt, text.WriteCellOpts(
		cell.FgColor(cell.ColorGreen)),
	)
}

func PrintBlue(txt string) {
	logBox.Write(txt, text.WriteCellOpts(
		cell.FgColor(cell.ColorBlue)),
	)
}

func LPrintGreen(txt string) {
	imageBox.Write(txt, text.WriteCellOpts(cell.FgColor(cell.ColorGreen)))
}

func LPrintBlue(txt string) {
	imageBox.Write(txt, text.WriteCellOpts(cell.FgColor(cell.ColorBlue)))
}
