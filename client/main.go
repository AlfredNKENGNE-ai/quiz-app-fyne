package main

import (
        "fyne.io/fyne/v2"
        "fyne.io/fyne/v2/app"
        "fyne.io/fyne/v2/theme"
)

var App fyne.App
var MainWindow fyne.Window

func main() {
        App = app.New()
        // Theme sombre mais moderne
        App.Settings().SetTheme(theme.DarkTheme())

        MainWindow = App.NewWindow("Quiz Battle 🕹️")
        MainWindow.Resize(fyne.NewSize(420, 720))

        InitNetwork()
        ShowLoginScreen()

        MainWindow.ShowAndRun()
}
