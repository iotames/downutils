package fyne

// https://github.com/fyne-io/fyne
// https://developer.fyne.io/started/

import (
	"log"
	"os"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/theme"
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
	// if conf.FyneFont == "" && lenfont == 0 {
	// 	conf.FyneFont = "resource/fonts/SourceHanSans-Bold.ttf"
	// }
	if conf.FyneLogoResource == nil {
		conf.FyneLogoResource = theme.FyneLogo()
	}
	if conf.Title == "" {
		conf.Title = "下载小工具"
	}
	return &FyneApp{Conf: conf}
}

func (f *FyneApp) Start() {
	app := app.New()
	app.SetIcon(f.Conf.FyneLogoResource)
	fyneFontEnv := os.Getenv("FYNE_FONT")
	log.Printf("---Start--fyneFontEnv(%s)----\n", fyneFontEnv)
	if !strings.Contains(fyneFontEnv, ".ttf") {
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
