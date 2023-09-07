package fyne

import (
	"fmt"
	"strconv"

	"github.com/iotames/downutils/conf"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func RenderSpiderSetting(w fyne.Window) fyne.CanvasObject {
	c := conf.GetSpiderConf()

	userAgentInput := widget.NewEntry()
	userAgentInput.Text = c.UserAgent

	httpProxyInput := widget.NewEntry()
	httpProxyInput.Text = c.HttpProxy
	itemProxy := widget.NewFormItem("Http代理", httpProxyInput)
	itemProxy.HintText = "http://127.0.0.1:1080"

	cookieInput := widget.NewEntry()
	cookieInput.Text = c.Cookie

	imgsPathInput := widget.NewEntry()
	imgsPathInput.Text = c.ImagesPath

	timeoutInput := widget.NewEntry()
	timeoutInput.Text = fmt.Sprintf("%d", c.Timeout)

	retryInput := widget.NewEntry()
	retryInput.Text = fmt.Sprintf("%d", c.RetryTimes)

	delayInput := widget.NewEntry()
	delayInput.Text = fmt.Sprintf("%d", c.Delay)

	randomDelayInput := widget.NewEntry()
	randomDelayInput.Text = fmt.Sprintf("%d", c.RandomDelay)

	parallelismInput := widget.NewEntry()
	parallelismInput.Text = fmt.Sprintf("%d", c.Parallelism)

	form := NewFyForm(
		widget.NewFormItem("UserAgent", userAgentInput),
		itemProxy,
		widget.NewFormItem("图片下载目录", imgsPathInput),
		widget.NewFormItem("Cookie", cookieInput),
		widget.NewFormItem("超时/秒", timeoutInput),
		widget.NewFormItem("自动重试/次", retryInput),
		widget.NewFormItem("延时/秒", delayInput),
		widget.NewFormItem("额外随机延时/秒", randomDelayInput),
		widget.NewFormItem("并发数", parallelismInput),
	)

	form.OnSubmit = func() {
		_, err1 := strconv.Atoi(timeoutInput.Text) // strconv.ParseUint(timeoutInput.Text, 10, 8)
		_, err2 := strconv.Atoi(retryInput.Text)
		_, err3 := strconv.Atoi(delayInput.Text)
		_, err4 := strconv.Atoi(randomDelayInput.Text)
		_, err5 := strconv.Atoi(parallelismInput.Text)
		if err1 != nil {
			dialog.NewError(err1, w)
			return
		}
		if err2 != nil {
			dialog.NewError(err2, w)
			return
		}
		if err3 != nil {
			dialog.NewError(err3, w)
			return
		}
		if err4 != nil {
			dialog.NewError(err4, w)
			return
		}
		if err5 != nil {
			dialog.NewError(err5, w)
			return
		}

		conf.UpdateConf(map[string]string{
			"USER_AGENT":   userAgentInput.Text,
			"HTTP_PROXY":   httpProxyInput.Text,
			"IMAGES_PATH":  imgsPathInput.Text,
			"COOKIE":       cookieInput.Text,
			"TIMEOUT":      timeoutInput.Text,
			"RETRY_TIMES":  retryInput.Text,
			"DELAY":        delayInput.Text,
			"RANDOM_DELAY": randomDelayInput.Text,
			"PARALLELISM":  parallelismInput.Text,
		})

		conf.LoadEnv()
		// dialog can prevent repeat click
		msg := dialog.NewInformation("提示", "提交成功！", w)
		msg.Show()
	}

	return container.NewVBox(form)
}
