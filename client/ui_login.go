package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func ShowLoginScreen() {
	title := widget.NewLabelWithStyle(
		"🕹️ QUIZ BATTLE",
		fyne.TextAlignCenter,
		fyne.TextStyle{Bold: true},
	)

	email := widget.NewEntry()
	email.SetPlaceHolder("Email")

	password := widget.NewPasswordEntry()
	password.SetPlaceHolder("Mot de passe")

	loginBtn := widget.NewButtonWithIcon("Se connecter 🚀", theme.ConfirmIcon(), func() {
		SendLogin(email.Text, password.Text)
	})

	card := widget.NewCard(
		"Connexion",
		"Entre dans la partie",
		container.NewVBox(
			email,
			password,
			loginBtn,
		),
	)

	MainWindow.SetContent(
		container.NewCenter(
			container.NewVBox(
				title,
				card,
			),
		),
	)
}
