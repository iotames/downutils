package service

import (
	"bytes"
	"fmt"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"log"
	"os"
	"strings"

	"github.com/gocolly/colly/v2"
	"github.com/iotames/downutils/conf"
	"github.com/iotames/miniutils"
	"github.com/xuri/excelize/v2"
	_ "golang.org/x/image/webp"
)

func (e *ExcelService) setLocalImagesByIndex(sheetName string, colIndex, rowIndex int, baseUrl, dirname string, setLocalPath func(excelImage *ExcelImage)) error {
	f := e.ExcelFile
	cols, err := f.Cols(sheetName)
	if err != nil {
		return err
	}

	colI := 'A'
	countttt := 1
	for {
		if countttt == colIndex {
			break
		}
		colI++
		countttt += 1
	}

	var excelImages []ExcelImage

	domain := miniutils.GetDomainByUrl(baseUrl)
	if dirname == "" {
		dinfo := strings.Split(domain, ".")
		dirname = dinfo[0]
		if dinfo[0] == "www" {
			dirname = dinfo[1]
		}
	}
	sconf := conf.GetSpiderConf()
	miniutils.Mkdir(sconf.ImagesPath + "/" + dirname)
	colHeiht := conf.ImageHeight
	colWidth := conf.ColWidth
	imgwebpage := conf.ImageWebpage

	// 遍历数据列
	fmt.Printf("------setLocalImagesByIndex---countttt(%d)--colI(%s)--colIndex(%d)--rowIndex(%d)---\n", countttt, string(colI), colIndex, rowIndex)
	var colCount int = 1
	for cols.Next() {
		col, err := cols.Rows()
		if err != nil {
			// 遍历列失败
			return err
		}
		fmt.Printf("-------setLocalImagesByIndex---cols.Next----colCount(%d) == colIndex(%d)-----\n", colCount, colIndex)
		if colCount == colIndex {
			// 定位到开始的列
			lencol := len(col)
			if rowIndex > lencol {
				err = fmt.Errorf("定位行数第(%d)行开始，已超过最大行数（%d）", rowIndex, lencol)
				return err
			}
			// dt = col[rowIndex-1:]

			rowI := 1
			f.SetColWidth(sheetName, fmt.Sprintf("%c", colI), fmt.Sprintf("%c", colI), float64(colWidth))
			for _, cell := range col {
				f.SetRowHeight(sheetName, rowI, float64(colHeiht))

				if rowI < rowIndex {
					// 下载起始行数之前的数据一律跳过
					rowI++
					continue
				}
				axis := fmt.Sprintf("%c%d", colI, rowI)
				fileurl := strings.TrimSpace(cell)
				fmt.Printf("\n----setLocalImagesByIndex--axis(%s)--fileurl(%s)----\n", axis, fileurl)
				excelImg := ExcelImage{Axis: axis, Url: fileurl}
				if fileurl == "" {
					excelImg.LocalPath = ""
				} else {
					// 设置excelImg.LocalPath
					setLocalPath(&excelImg)
				}

				excelImages = append(excelImages, excelImg)

				// f.SetCellValue(sheetName, axis, excelImg.LocalPath)
				rowI++
			}

			break
		}
		colCount++

	}
	// if len(excelImages) == 0 {
	// 	log.Println("imgTitle is not Exist")
	// 	return fmt.Errorf("excel文件里%d列%d行开始，找不到任何有效的下载链接", colIndex, rowIndex)
	// }

	// AddPicture 不指定图片栏列宽度，图片无法填满整个单元格
	f.SetColWidth(sheetName, string(colI), string(colI), float64(colWidth))

	for _, excelImg := range excelImages {
		if excelImg.LocalPath == "" {
			continue
		}
		f.SetCellValue(sheetName, excelImg.Axis, "")
		imgopts := &excelize.GraphicOptions{AutoFit: true, LockAspectRatio: true, HyperlinkType: "External"}
		if imgwebpage {
			imgopts.Hyperlink = excelImg.Url
		}
		err := f.AddPicture(sheetName, excelImg.Axis, excelImg.LocalPath,
			// `{"x_scale": 0.2, "y_scale": 0.2}`,
			imgopts,
		)
		// fmt.Println(excelImg)
		if err != nil {
			fmt.Println("-------AddPicture---Error: ", err, excelImg.LocalPath)
		}
	}

	fmt.Println("DownloadImages By File: " + e.ExcelFile.Path + " Success!")
	return nil
}

func (e *ExcelService) DownloadImagesByCollyLocation(sheetName, referer, dirname string, colIndex, rowIndex int, onResp func(furl string), withImgFile bool) error {
	var err error
	if dirname == "" {
		urlSplit := strings.Split(referer, "/")
		if len(urlSplit) > 2 {
			dmm := urlSplit[2]
			dmmSplit := strings.Split(dmm, ".")
			if len(dmmSplit) > 1 {
				dirname = dmmSplit[len(dmmSplit)-2]
			}
		}
	}
	spd := NewSpider(dirname)
	baseUrl := miniutils.GetBaseUrl(referer)
	spd.BaseUrl = baseUrl
	c := spd.GetCollector()
	spd.SetAsyncAndLimit(c, ".")

	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("referer", baseUrl+"/")
	})

	c.OnResponse(func(r *colly.Response) {
		LocalPath := r.Ctx.Get("LocalPath")
		// LocalPath := r.Ctx.Get(r.Request.URL.String())
		log.Println("Download To LocalPath:", LocalPath)
		var f *os.File
		f, err = os.Create(LocalPath)
		if err != nil {
			log.Printf("\n---Error Happened. os.Create:%v-------\n", err)
		}
		_, err = io.Copy(f, bytes.NewReader(r.Body))
		if err != nil {
			log.Printf("---Error Happened. io.Copy:%v-------\n", err)
		}
		onResp(r.Request.URL.String())
	})

	snames := e.ExcelFile.GetSheetList()
	// reqcount := 0
	for _, sn := range snames {
		if sheetName != "" && strings.TrimSpace(sn) != strings.TrimSpace(sheetName) {
			continue
		}
		var imgs []string
		imgs, err = e.GetColsBegin(sn, colIndex, rowIndex)
		if err != nil {
			return fmt.Errorf("DownloadImagesByCollyLocation获取列数据失败：%v", err)
		}
		for _, img := range imgs {
			imgUrl := strings.TrimSpace(img)
			if imgUrl == "" {
				continue
			}
			filepath := spd.GetLocalFile(imgUrl, dirname, LOCAL_IMAGE_FILE_EXT)
			isExist := miniutils.IsPathExists(filepath)
			if isExist {
				onResp(imgUrl)
				log.Printf("-----Skip--DownloadImagesByCollyLocation----sheetName(%s)--imgUrl(%s)--filepath(%s)--is exist---", sn, imgUrl, filepath)
				continue
			}
			if strings.Index(imgUrl, "http") != 0 {
				log.Printf("-----Skip--DownloadImagesByCollyLocation--UrlNotHttp--sheetName(%s)--imgUrl(%s)---", sn, imgUrl)
				continue
			}
			ctx := colly.NewContext()
			ctx.Put("LocalPath", filepath)
			c.Request("GET", imgUrl, nil, ctx, nil)
		}
	}
	fmt.Println("-------DownloadImagesByCollyLocation----c.Wait()----------")
	c.Wait()
	if err != nil {
		return err
	}
	var result error
	if withImgFile {
		for _, sn := range snames {
			if sheetName != "" && strings.TrimSpace(sn) != strings.TrimSpace(sheetName) {
				continue
			}
			res := e.setLocalImagesByIndex(sn, colIndex, rowIndex, baseUrl, dirname, func(excelImg *ExcelImage) {
				imgurl := excelImg.Url
				// log.Printf("---ReadRow--DownloadImagesByColly-debug-sheetName(%s)--imgurl(%s)--filepath(%s)----", sn, imgurl, filepath)
				if imgurl == "" {
					log.Println("Skip: DownloadImagesByColly Request: DownloadUrl is empty----------------")
					return
				}
				filepath := spd.GetLocalFile(imgurl, dirname, LOCAL_IMAGE_FILE_EXT)
				excelImg.LocalPath = filepath
				isExist := miniutils.IsPathExists(filepath)
				if isExist {
					log.Println("Skip: DownloadImagesByColly Request: filepath is exist---", filepath)
					return
				}
			})
			if res != nil {
				result = res
			}
		}
	}
	// c.Wait() Wait Fail
	return result
}
