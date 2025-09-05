package fyne

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func GetWelcomeContent(a fyne.App, w fyne.Window) fyne.CanvasObject {
	// chromedpTest := widget.NewButton("浏览器爬虫", func() {
	// 	manualSpider(a, w)
	// })
	// webManager := widget.NewButton("Web控制台", func() {
	// 	go func() {
	// 		webserver.Run()
	// 	}()
	// 	webc := config.NewWebserverConfigData()
	// 	util.StartBrowserByUrl(webc.BaseUrl + "/client")
	// })

	downfiles := widget.NewButton("文件批量下载器", func() {
		subw := a.NewWindow("文件批量下载器")
		c := container.NewVBox(RenderFormDownImg(subw))
		subw.SetContent(c)
		subw.Resize(fyne.NewSize(720, 500))
		subw.Show()
		// SetMainContant(w, c)
	})
	content := hbox(
		// chromedpTest,
		// webManager,
		downfiles,
	)
	return content
}
