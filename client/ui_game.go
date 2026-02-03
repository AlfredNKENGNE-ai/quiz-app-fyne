package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func ShowQuestionScreen(question string, options []string, questionID int) {
	questionLabel := widget.NewLabelWithStyle(
		question,
		fyne.TextAlignCenter,
		fyne.TextStyle{Bold: true},
	)

	var buttons []fyne.CanvasObject
	for i, opt := range options {
		index := i
		btn := widget.NewButton(opt, func() {
			SendAnswer(questionID, index)
		})
		buttons = append(buttons, btn)
	}

	MainWindow.SetContent(
		container.NewVBox(
			questionLabel,
			container.NewGridWithRows(2, buttons...),
		),
	)
}

func ShowRiddleScreen(text string) {
	answer := widget.NewEntry()
	answer.SetPlaceHolder("Ta réponse...")

	submit := widget.NewButtonWithIcon("Valider ✅", theme.ConfirmIcon(), func() {
		SendRiddleAnswer(answer.Text)
	})

	hint1 := widget.NewButton("Indice -25 pts 💡", func() {
		RequestHint(1)
	})
	hint2 := widget.NewButton("Indice -50 pts 💡", func() {
		RequestHint(2)
	})

	MainWindow.SetContent(
		container.NewVBox(
			widget.NewLabelWithStyle(text, fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
			answer,
			submit,
			container.NewGridWithColumns(2, hint1, hint2),
		),
	)
}
