package main

import (
	"context"
	"fmt"

	"github.com/mum4k/termdash"
	"github.com/mum4k/termdash/cell"
	"github.com/mum4k/termdash/container"
	"github.com/mum4k/termdash/linestyle"
	"github.com/mum4k/termdash/terminal/termbox"
	"github.com/mum4k/termdash/terminal/terminalapi"
	"github.com/mum4k/termdash/widgets/button"
	"github.com/mum4k/termdash/widgets/gauge"
	"github.com/mum4k/termdash/widgets/text"
)

func writeLines(ctx context.Context, t *text.Text, text string) {
	t.Write(fmt.Sprintf("%s\n", text))
}

func InitUI() {
	var err error
	tmx, err = termbox.New()
	// TEMP:
	if err != nil {
		panic(err)
	}
	defer tmx.Close()
	ctx, cancel := context.WithCancel(context.Background())

	//This is where the logs are
	logBox, err = text.New(text.RollContent(), text.WrapAtWords())
	if err != nil {
		panic(err)
	}
	PrintGreenln("Terminal successfully launched.")

	//This is where the last pokemon appear
	imageBox, err = text.New(text.RollContent(), text.WrapAtWords())
	if err != nil {
		panic(err)
	}

	//This is where the pokemon list refreshes
	ProgressBar, err = gauge.New(
		gauge.Height(1),
		gauge.Color(cell.ColorBlue),
		gauge.Border(linestyle.Light),
		gauge.BorderTitle("Pokemon List Refresh"),
	)
	ProgressBar.Absolute(0, 0)

	//This is the clear log button
	ClearLogs, _ := button.New("Clear logs", func() error {
		logBox.Reset()
		return nil
	},
		button.GlobalKey('r'),
		button.FillColor(cell.ColorYellow),
	)

	//This is the catch latest pokemon button
	CatchLast, _ := button.New("Catch latest", func() error {
		CatchLatest()
		return nil
	},
		button.GlobalKey('c'),
		button.FillColor(cell.ColorYellow),
	)

	c, err := container.New(
		tmx,
		container.Border(linestyle.Light),
		container.BorderTitle("PRESS Q TO QUIT"),
		container.SplitVertical(
			container.Left(
				container.Border(linestyle.Light),
				container.BorderTitle("Logs"),
				container.PlaceWidget(logBox),
			),
			container.Right(
				container.SplitHorizontal(
					container.Top(
						container.SplitHorizontal(
							container.Top(
								container.PlaceWidget(ProgressBar),
							),
							container.Bottom(
								container.SplitVertical(
									container.Left(
										container.PlaceWidget(ClearLogs),
									),
									container.Right(
										container.PlaceWidget(CatchLast),
									),
									container.SplitPercent(50),
								),
							),
							container.SplitPercent(50),
						),
					),
					container.Bottom(
						container.Border(linestyle.Light),
						container.BorderTitle("Last Pokemon :"),
						container.PlaceWidget(imageBox),
					),
					container.SplitPercent(50),
				),
			),
		),
	)
	if err != nil {
		panic(err)
	}
	go Useful_Variables()

	quitter := func(k *terminalapi.Keyboard) {
		if k.Key == 'q' || k.Key == 'Q' {
			cancel()
		}
	}

	if err := termdash.Run(ctx, tmx, c, termdash.KeyboardSubscriber(quitter)); err != nil {
		panic(err)
	}
}
