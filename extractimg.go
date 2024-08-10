package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/xuri/excelize/v2"
)

var filename string
var dirname string
var sheetname string
var imgcol string
var imgnamecol string

func Extractimg() {
	flag.StringVar(&filename, "filename", "runtime/111.xlsx", "输入Excel文件名")
	flag.StringVar(&dirname, "dirname", "runtime/images", "输入下载目录名")
	flag.StringVar(&sheetname, "sheetname", "Sheet1", "sheet名称")
	flag.StringVar(&imgcol, "imgcol", "A", "图片所在列")
	flag.StringVar(&imgnamecol, "imgnamecol", "B", "图片名称所在列")
	flag.Parse()

	// dirname = "runtime/images"
	// sheetname = "Sheet1"
	// imgcol = "A"     // 图片所在列
	// imgnamecol = "B" // 图片名称所在列
	// filename = "runtime/111.xlsx"

	f, err := excelize.OpenFile(filename)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	total := 0
	okk := 0
	for i := 2; i < 404; i++ {
		total++
		picin := fmt.Sprintf("%s%d", imgcol, i)
		pics, err := f.GetPictures(sheetname, picin)
		if err != nil {
			fmt.Printf("-----error:GetPictures--picin(%s)--err(%v)---\n", picin, err)
			// panic(err)
		}
		picname, err := f.GetCellValue(sheetname, fmt.Sprintf("%s%d", imgnamecol, i))
		if err != nil {
			fmt.Printf("-----GetCellValue--err(%v)-----\n", err)
		}
		if len(pics) == 0 {
			err = fmt.Errorf("------picin(%s)--pics-len=0--picname(%s)---", picin, picname)
			fmt.Println(err)
			// panic(err)
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
		// fmt.Printf("----savePics---(%s)--------\n", filepath)
		if err := os.WriteFile(filepath, pic.File, 0644); err != nil {
			fmt.Printf("------pic-write-err(%v)--writeFile(%s)------\n", err, filepath)
		}
		ok++
	}
	return ok
}
