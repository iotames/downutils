package fyne

import (
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func SetMainContant(w fyne.Window, content fyne.CanvasObject) {
	w.SetContent(container.NewVBox(getTimeLabelContent(), content))
	w.Resize(fyne.NewSize(640, 460))
	w.Show()
}

func getTimeLabelContent() fyne.CanvasObject {
	_, timeLabel, _ := MakeWelcomeUI()
	content := container.NewHBox(layout.NewSpacer(), timeLabel) // HBox 左右 VBox 上下
	// content := container.NewVBox(helloBox)
	UpdateTime(timeLabel)
	go func() {
		for range time.Tick(time.Second) {
			UpdateTime(timeLabel)
		}
	}()
	return content
}

func UpdateTime(label *widget.Label) {
	timeText := time.Now().Format("2006-01-02 15:04:05") // 当前时间: 2006-01-02 15:04:05
	label.SetText(timeText)
}

func MakeWelcomeUI() (*widget.Label, *widget.Label, *widget.Entry) {
	helloName := widget.NewLabel("您好!")
	timeLabel := widget.NewLabel("")
	helloInput := widget.NewEntry()
	helloInput.OnChanged = func(s string) {
		helloName.SetText(s + ", 您好!")
	}
	return helloName, timeLabel, helloInput
}

func NewFyForm(items ...*widget.FormItem) *widget.Form {
	form := &widget.Form{
		Items:      items,
		SubmitText: "提交",
	}
	return form
}

// NewFyRadio. 值类型 bool.
func NewFyRadio(optDic map[bool]string, value bool, changed func(string)) *widget.RadioGroup {
	var opts []string
	for _, val := range optDic {
		opts = append(opts, val)
	}
	radio := widget.NewRadioGroup(opts, changed)
	radio.Selected = optDic[value]
	radio.Horizontal = true
	return radio
}

func CheckError(err error, w fyne.Window) {
	if err != nil {
		dialog.NewError(err, w).Show()
	}
}
