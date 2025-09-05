package fyne

import (
	// "context"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/iotames/downutils/service"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	sqdialog "github.com/sqweek/dialog"

	// "fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
)

func RenderFormDownImg(w fyne.Window) fyne.CanvasObject {
	filepathInput := widget.NewEntry()
	filepathInput.PlaceHolder = "读取Excel表格某一列，批量下载。填xlsx文件路径"

	// fpicker := hbox(filepathInput) //  container.NewHBox() widget.NewLabel("选择xlsx文件"),

	sheetInput := widget.NewEntry()
	sheetInput.Text = ""
	sheetInput.PlaceHolder = "例: Sheet1"
	sheetItem := widget.NewFormItem("Sheet名(可选)", sheetInput)
	sheetItem.HintText = "Excel的sheet子表。留空则读取所有子表"

	imgLocationInput := widget.NewEntry()
	imgLocationInput.Text = "1,2"
	imgLocationInput.PlaceHolder = "例: 1,2表示从第1列第2行开始查找整列的数据"
	downLocationItem := widget.NewFormItem("行列位置", imgLocationInput)
	downLocationItem.HintText = "表示从第几列第几行开始下载整列的数据"

	imgdirnameInput := widget.NewEntry()
	imgdirnameInput.Text = ""
	imgdirnameInput.PlaceHolder = "例: amazon"
	imgdirItem := widget.NewFormItem("下载目录(可选)", imgdirnameInput)
	imgdirItem.HintText = "文件夹名。留空则从Referer值生成。"

	refererInput := widget.NewEntry()
	refererInput.Text = "https://image.baidu.com/"
	refererInput.PlaceHolder = "填网站主页。例: https://www.amazon.com/"
	excelWithPicInput := widget.NewCheck("启用", nil) // NewFyRadio(map[bool]string{false: "否", true: "是"}, false, nil)
	withPicItem := widget.NewFormItem("另存为含图片的xlsx文件", excelWithPicInput)
	withPicItem.HintText = "不是下载图片，请勿开启此功能"
	mainForm := NewFyForm(
		widget.NewFormItem("xlsx文件路径", filepathInput),
		sheetItem,
		downLocationItem,
		imgdirItem,
		widget.NewFormItem("Referer", refererInput),
		withPicItem,
	)
	fprogress := widget.NewProgressBar()
	fprogress.TextFormatter = func() string { return fmt.Sprintf("下载进度: %.0f of %.0f", fprogress.Value, fprogress.Max) }
	var isRunning bool = false
	// ctx, cancel := context.WithCancel(context.Background())
	mainForm.SubmitText = "开始"
	mainForm.OnSubmit = func() {
		fpath := filepathInput.Text
		if fpath == "" {
			CheckError(fmt.Errorf("xlsx文件路径不能为空"), w)
			return
		}
		refererUrl := refererInput.Text
		if refererUrl == "" {
			refererUrl = "https://image.baidu.com/"
		}
		if strings.Index(refererUrl, "http") != 0 {
			CheckError(fmt.Errorf("Referer必须以http开头"), w)
			return
		}

		dialog.NewConfirm("确认", "爬虫确认执行", func(b bool) {
			if b {
				var err error
				var colIndex, rowIndex int
				if isRunning {
					CheckError(fmt.Errorf("下载进行中，请勿频繁点击"), w)
					return
				}
				isRunning = true
				newXlsxWithPic := excelWithPicInput.Checked
				imgLocationSplit := strings.Split(imgLocationInput.Text, ",")
				if len(imgLocationSplit) != 2 {
					imgLocationSplit = strings.Split(imgLocationInput.Text, "，")
					if len(imgLocationSplit) != 2 {
						fmt.Printf("---------22222222--NotIndex(2)---len(%d)---(%+v)----", len(imgLocationSplit), imgLocationSplit)
						CheckError(fmt.Errorf("行列位置必须为1,1格式。"), w)
						isRunning = false
						return
					}

				}
				colIndex, err = strconv.Atoi(strings.TrimSpace(imgLocationSplit[0]))
				if err != nil {
					CheckError(fmt.Errorf("行列位置解析错误(%v)", err), w)
					return
				}
				rowIndex, err = strconv.Atoi(strings.TrimSpace(imgLocationSplit[1]))
				if err != nil {
					CheckError(fmt.Errorf("行列位置解析错误(%v)", err), w)
					return
				}
				log.Printf("------dialog.NewConfirm----filepath(%s)--sheetName(%s)--imgLocationInput(%s)--location(%d,%d)-------fileWithPic(%v)--\n", fpath, sheetInput.Text, imgLocationInput.Text, colIndex, rowIndex, newXlsxWithPic)
				dname := imgdirnameInput.Text
				// imgs, err := service.GetImgsByExcel(fpath, sheetInput.Text, imgtitleInput.Text)
				imgs, err := service.GetImgsByExcelIndex(fpath, sheetInput.Text, colIndex, rowIndex)
				if err != nil {
					CheckError(fmt.Errorf("service.GetImgsByExcel错误:%v", err), w)
					isRunning = false
					return
				}
				for i, imgdt := range imgs {
					fmt.Printf("------imgs--(%d)--img(%s)---\n", i, imgdt)
				}
				fprogress.Max = float64(len(imgs))
				fprogress.SetValue(0)
				go func() {
					err = service.DownloadImagesByExcel(fpath, sheetInput.Text, refererUrl, strings.TrimSpace(dname), colIndex, rowIndex, newXlsxWithPic, func(furl string) {
						fprogress.SetValue(fprogress.Value + 1)
					})
					if err != nil {
						CheckError(fmt.Errorf("下载错误(%v)", err), w)
						return
					}
					isRunning = false
					dialog.NewInformation("提示", "下载完成", w).Show()
				}()
				// go func(ctx context.Context) {
				// 	defer fmt.Println("---Out--Of---Loop---")
				// 	for {
				// 		select {
				// 		case <-ctx.Done():
				// 			fmt.Println("子 协程 接受停止信号...")
				// 			isRunning = false
				// 			runtime.Goexit()
				// 			return
				// 		default:
				// 			if !isRunning {
				// 				go func() {
				// 					isRunning = true
				// 					mainForm.SubmitText = "下载中..."
				// 					service.DownloadImagesByExcel(fpath, sheetInput.Text, imgtitleInput.Text, refererUrl, strings.TrimSpace(dname), newXlsxWithPic, func(furl string) {
				// 						fprogress.SetValue(fprogress.Value + 1)
				// 					})
				// 					isRunning = false
				// 					mainForm.SubmitText = "开始"
				// 				}()
				// 			}
				// 		}
				// 	}
				// }(ctx)

				// msg := dialog.NewInformation("提示", "已提交", w)
				// msg.Show()
			}
		}, w).Show()
	}

	// mainForm.OnCancel = func() {
	// 	if isRunning {
	// 		dialog.NewConfirm("取消下载", "操作正在进行中，确认取消下载？", func(b bool) {
	// 			if b {
	// 				go func() {
	// 					fmt.Println("---send ctx Cancel-----")
	// 					cancel()
	// 				}()
	// 			}
	// 		}, w).Show()
	// 	} else {
	// 		CheckError(fmt.Errorf("下载未开始，无法取消"), w)
	// 	}
	// }
	// mainForm.CancelText = "取消"

	filePicker := widget.NewButton("选择xlsx文件", func() {
		filename, err := sqdialog.File().Filter("Excel表格(*.xlsx)", "xlsx").Load()
		if err == nil {
			filepathInput.SetText(filename)
		} else {
			CheckError(err, w)
		}
		// fd := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
		// 	if err != nil {
		// 		dialog.ShowError(err, w)
		// 		return
		// 	}
		// 	if reader == nil {
		// 		fmt.Println("Cancelled")
		// 		return
		// 	}
		// 	defer reader.Close()
		// 	filepathInput.SetText(reader.URI().Path())
		// 	// TODO 自动填写referer, imgdirName
		// }, w)
		// fd.SetFilter(storage.NewExtensionFileFilter([]string{".xlsx"}))
		// fd.Show()
	})

	return vbox(filePicker, mainForm, fprogress) // container.NewVBox(row1, form)
}
