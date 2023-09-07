package fyne

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
)

func hbox(objects ...fyne.CanvasObject) *fyne.Container {
	return container.NewHBox(objects...)
}

func vbox(objects ...fyne.CanvasObject) *fyne.Container {
	return container.NewVBox(objects...)
}
