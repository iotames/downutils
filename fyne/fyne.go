package fyne

// https://github.com/fyne-io/fyne
// https://developer.fyne.io/started/

import (
	"log"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"github.com/iotames/downutils/conf"
	"github.com/iotames/miniutils"
)

type FyneConf struct {
	// FyFontPath       string
	FyneLogoResource fyne.Resource
	FyneFontResource fyne.Resource
	Title            string
}

type FyneApp struct {
	Conf         *FyneConf
	masterWindow fyne.Window
}

func NewFyneApp(conf *FyneConf) *FyneApp {
	return &FyneApp{Conf: conf}
}

func (f *FyneApp) Start() {
	app := app.New()
	app.SetIcon(f.Conf.FyneLogoResource)
	log.Printf("---Start--SetEnv(FYNE_FONT:%s)FromFile(%s)----\n", conf.FyneFont, conf.WorkEnvFile)
	if !strings.Contains(conf.FyneFont, ".ttf") || !miniutils.IsPathExists(conf.FyneFont) {
		log.Printf("---Start--SetThemeByFyneFontResource(%s)----\n", f.Conf.FyneFontResource.Name())
		app.Settings().SetTheme(NewFyneTheme(f.Conf.FyneFontResource))
	}
	window := app.NewWindow(f.Conf.Title)
	window.SetMaster()
	f.masterWindow = window
	mainMenu := getMainMenu(app, window)
	window.SetMainMenu(mainMenu)

	// go func() {
	// 	time.Sleep(2 * time.Second)
	// 	for _, item := range mainMenu.Items[0].Items {
	// 		if item.Label == "Quit" {
	// 			item.Label = "退出"
	// 		}
	// 	}
	// }()

	SetMainContant(window, GetWelcomeContent(app, window))
	// startDemo(f)
	// startDataBinding()
	// startPreferencesAPI()

	app.Run()
}
