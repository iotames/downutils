package fyne

import (
	"fyne.io/fyne/v2"
	// "fyne.io/fyne/v2/cmd/fyne_settings/settings"
	"fyne.io/fyne/v2/container"
	// "fyne.io/fyne/v2/dialog"
)

func getMainMenu(a fyne.App, w fyne.Window) *fyne.MainMenu {
	// newItem := fyne.NewMenuItem("New", nil)
	// otherItem := fyne.NewMenuItem("Other", nil)
	// otherItem.ChildMenu = fyne.NewMenu("",
	// 	fyne.NewMenuItem("Project", func() { fmt.Println("Menu New->Other->Project") }),
	// 	fyne.NewMenuItem("Mail", func() { fmt.Println("Menu New->Other->Mail") }),
	// )
	// newItem.ChildMenu = fyne.NewMenu("",
	// 	fyne.NewMenuItem("File", func() { fmt.Println("Menu New->File") }),
	// 	fyne.NewMenuItem("Directory", func() { fmt.Println("Menu New->Directory") }),
	// 	otherItem,
	// )

	// renderSetItem := fyne.NewMenuItem("外观设置", func() {
	// 	w := a.NewWindow("外观设置")
	// 	w.SetContent(settings.NewSettings().LoadAppearanceScreen(w))
	// 	w.Resize(fyne.NewSize(480, 480))
	// 	w.Show()
	// 	// SetMainContant(w, settings.NewSettings().LoadAppearanceScreen(w))
	// })

	// mainConfig := fyne.NewMenuItem("系统设置", nil)
	// mainConfig.ChildMenu = fyne.NewMenu("",
	// 	fyne.NewMenuItem("下载器设置", func() { renderSpiderConf(w) }),
	// )
	spiderSet := fyne.NewMenuItem("爬虫设置", func() {
		w := a.NewWindow("爬虫设置")
		c := container.NewVBox(RenderSpiderSetting(w))
		w.SetContent(c)
		w.Resize(fyne.NewSize(880, 720))
		w.Show()
	})

	// initConfig := fyne.NewMenuItem("恢复出厂设置", func() {
	// 	dialog.NewConfirm("警告", "本操作会丢失原有设置!", func(b bool) {
	// 		if b {
	// 			f, err := os.OpenFile(conf.ENV_FILE, os.O_TRUNC|os.O_WRONLY, 0644)
	// 			if err != nil {
	// 				CheckError(err, w)
	// 				// panic("open .env Error: " + err.Error())
	// 			}
	// 			defer conf.LoadEnv()
	// 			defer f.Close()
	// 			_, err = f.WriteString("")
	// 			if err != nil {
	// 				panic(err)
	// 			}
	// 		}
	// 	}, w).Show()
	// })

	settingMenu := fyne.NewMenu("设置",
		spiderSet,
		// renderSetItem,
		fyne.NewMenuItemSeparator(),
	)

	// a quit item will be appended to our first (File) menu
	homeMenu := fyne.NewMenuItem("首页", func() { SetMainContant(w, GetWelcomeContent(a, w)) })
	startMenu := fyne.NewMenu("开始", homeMenu)

	return fyne.NewMainMenu(
		startMenu,
		// siteMenu,
		settingMenu,
		// helpMenu,
	)
}
