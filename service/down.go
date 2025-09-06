package service

import (
	// "context"
	"fmt"
	// "log"
	// "strconv"
	// "strings"

	// "fyne.io/fyne/v2/dialog"

	"fyne.io/fyne/v2/widget"
)

type DownUtil struct {
	fromFile           string
	sheet              string
	colStart, rowStart int
	downDirName        string
	referer            string
	imgStickToFile     bool
	progressBar        *widget.ProgressBar
}

func NewDownUtil(downDirName string, progressBar *widget.ProgressBar) *DownUtil {
	return &DownUtil{
		downDirName: downDirName,
		progressBar: progressBar,
	}
}

func (d *DownUtil) SetImgsInFile() {
	d.imgStickToFile = true
}
func (d *DownUtil) SetReferer(referer string) {
	d.referer = referer
}

func (d *DownUtil) SetExcel(fpath, sheet string, colStart, rowStart int) {
	d.fromFile = fpath
	d.sheet = sheet
	d.colStart = colStart
	d.rowStart = rowStart
}

func (d DownUtil) getColData() (data []string, err error) {
	data, err = GetImgsByExcelIndex(d.fromFile, d.sheet, d.colStart, d.rowStart)
	return
}

func (d *DownUtil) Run() error {
	imgs, err := d.getColData()
	if err != nil {
		return fmt.Errorf("service.GetImgsByExcel错误:%v", err)
	}
	for i, imgdt := range imgs {
		fmt.Printf("------imgs--i(%d)--img(%s)---\n", i, imgdt)
	}
	d.progressBar.Max = float64(len(imgs))
	d.progressBar.SetValue(0)
	return d.Download()
}

func (d *DownUtil) Download() error {
	err := DownloadImagesByExcel(d.fromFile, d.sheet, d.referer, d.downDirName, d.colStart, d.rowStart, d.imgStickToFile, func(furl string) {
		d.progressBar.SetValue(d.progressBar.Value + 1)
	})
	if err != nil {
		return fmt.Errorf("下载错误(%v)", err)
	}
	return err
}
