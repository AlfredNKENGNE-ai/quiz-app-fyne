package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func ShowWaitingRoom() {
	MainWindow.SetContent(
		container.NewCenter(
			widget.NewLabelWithStyle(
				"⏳ La partie va commencer...",
				fyne.TextAlignCenter,
				fyne.TextStyle{Bold: true},
			),
		),
	)
}
