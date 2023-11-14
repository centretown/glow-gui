package main

import (
	"glow-gui/data"
	"glow-gui/resources"
	"glow-gui/settings"
	"glow-gui/store"
	"glow-gui/ui"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

func main() {
	app := app.NewWithID(resources.AppID)

	icon, err := resources.DarkGanderImage.Load()
	if err == nil {
		app.SetIcon(icon)
	}

	preferences := app.Preferences()
	theme := settings.NewGlowTheme(preferences)
	app.Settings().SetTheme(theme)

	store := store.NewStore(preferences)
	model := data.NewModel(store)

	window := app.NewWindow(resources.GlowLabel.String())
	ui := ui.NewUi(app, window, model, theme)

	window.SetContent(ui.BuildContent())
	window.SetCloseIntercept(func() {
		store.OnExit()
		ui.OnExit()
		size := window.Canvas().Size()
		preferences.SetInt(settings.ContentWidth.String(), int(size.Width))
		preferences.SetInt(settings.ContentHeight.String(), int(size.Height))
		window.Close()
	})

	width := preferences.IntWithFallback(settings.ContentWidth.String(), 0)
	height := preferences.IntWithFallback(settings.ContentHeight.String(), 0)
	if height > 0 && width > 0 {
		window.Resize(fyne.Size{Width: float32(width), Height: float32(height)})
	}

	window.Show()
	effect := preferences.StringWithFallback(settings.Effect.String(), "")
	if len(effect) > 0 {
		model.LoadFrame(effect)
	}

	app.Run()
}
