package main

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func ShowResults(results []string) {
	list := container.NewVBox()
	for i, r := range results {
		list.Add(widget.NewLabel(fmt.Sprintf("%d️⃣ %s", i+1, r)))
	}

	MainWindow.SetContent(
		container.NewCenter(
			container.NewVBox(
				widget.NewLabelWithStyle("🏆 Classement", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
				list,
			),
		),
	)
}
