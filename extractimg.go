package main

import (
	"fmt"
	"os"
	"strings"

	"flag"

	"github.com/xuri/excelize/v2"
)

var filename string
var dirname string
var sheetname string
var imgcol string
var imgnamecol string

func Extractimg() {
	flag.StringVar(&filename, "filename", "test.xlsx", "输入Excel文件名")
	flag.StringVar(&dirname, "dirname", "images", "输入下载目录名")
	flag.StringVar(&sheetname, "sheetname", "大货基础信息表", "sheet名称")
	flag.StringVar(&imgcol, "imgcol", "J", "图片所在列")
	flag.StringVar(&imgnamecol, "imgnamecol", "K", "图片名称所在列")

	flag.Parse()
	f, err := excelize.OpenFile(filename)
	if err != nil {
		fmt.Println(err)
		return
	}
	// dirname := "imgs"
	// sheetname := "大货基础信息表"
	// imgcol := "J"     // 图片所在列
	// imgnamecol := "K" // 图片名称所在列
	total := 0
	okk := 0
	for i := 2; i < 1845; i++ {
		total++
		picin := fmt.Sprintf("%s%d", imgcol, i)
		pics, err := f.GetPictures(sheetname, picin)
		if err != nil {
			fmt.Printf("-----GetPictures--err(%v)-----\n", err)
		}
		picname, err := f.GetCellValue(sheetname, fmt.Sprintf("%s%d", imgnamecol, i))
		if err != nil {
			fmt.Printf("-----GetCellValue--err(%v)-----\n", err)
		}
		result := savePics(dirname, picname, pics)
		if result > 0 {
			okk++
			fmt.Printf("-----SUCCESS---%s--%s--%d\n", picin, picname, result)
		} else {
			fmt.Printf("-----FAIL---%s--%s\n", picin, picname)
		}
	}
	fmt.Printf("-----Done(%d/%d)------\n", okk, total)
}

func savePics(dirname string, picname string, pics []excelize.Picture) int {
	picname = strings.TrimSpace(picname)
	if picname == "" {
		return 0
	}
	ok := 0
	for _, pic := range pics {
		picname := fmt.Sprintf("%s%s", picname, pic.Extension)
		filepath := dirname + "/" + picname
		if err := os.WriteFile(filepath, pic.File, 0644); err != nil {
			fmt.Println(picname, err)
		}
		ok++
	}
	return ok
}
