package main

import (
	"quiz-app-fyne/shared"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

var CurrentUser *shared.User

func ShowLobbyScreen() {
	codeEntry := widget.NewEntry()
	codeEntry.SetPlaceHolder("Code de la partie")

	createBtn := widget.NewButtonWithIcon("Créer une partie ➕", theme.ContentAddIcon(), func() {
		SendCreateGame(CurrentUser.ID, "multi")
	})

	joinBtn := widget.NewButtonWithIcon("Rejoindre 🎯", theme.MailSendIcon(), func() {
		if codeEntry.Text != "" {
			SendJoinGame(codeEntry.Text, CurrentUser.ID)
		}
	})

	MainWindow.SetContent(
		container.NewCenter(
			container.NewVBox(
				widget.NewLabelWithStyle("🕹️ Lobby", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
				createBtn,
				codeEntry,
				joinBtn,
			),
		),
	)
}
