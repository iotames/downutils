package main

import (
	fy "fyne.io/fyne/v2"
	"github.com/iotames/downutils/conf"
	"github.com/iotames/downutils/fyne"
)

func main() {
	fyne.NewFyneApp(&fyne.FyneConf{
		Title:            conf.FyneTitle,
		FyneLogoResource: fy.NewStaticResource("logo.png", logopng),
		FyneFontResource: fy.NewStaticResource("OPPOSans-H.ttf", fyfont),
	}).Start()
	// Extractimg()
}

func init() {
	conf.LoadEnv()
}
