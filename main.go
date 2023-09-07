package main

import (
	"github.com/iotames/downutils/conf"
	"github.com/iotames/downutils/fyne"

	fy "fyne.io/fyne/v2"
)

func main() {
	fyne.NewFyneApp(&fyne.FyneConf{
		Title:            conf.FyneTitle,
		FyneLogoResource: fy.NewStaticResource("logo.png", logopng),
		FyneFontResource: fy.NewStaticResource("OPPOSans-H.ttf", fyfont),
	}).Start()
}

func init() {
	conf.LoadEnv()
}
