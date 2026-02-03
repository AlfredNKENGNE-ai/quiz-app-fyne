package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func ShowModeSelectionScreen() {
	solo := widget.NewButtonWithIcon("🎮 Solo", theme.MediaPlayIcon(), func() {
		SendCreateGame(CurrentUser.ID, "solo")
	})

	multi := widget.NewButtonWithIcon("👥 Multijoueur", theme.AccountIcon(), func() {
		ShowLobbyScreen()
	})

	MainWindow.SetContent(
		container.NewVBox(
			widget.NewLabelWithStyle("Choisir le mode", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
			solo,
			multi,
		),
	)
}
