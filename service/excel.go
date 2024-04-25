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

	"github.com/iotames/downutils/conf"

	"github.com/iotames/miniutils"

	"github.com/gocolly/colly/v2"
	"github.com/xuri/excelize/v2"
	_ "golang.org/x/image/webp"
)

type ExcelService struct {
	// filepath  string
	ExcelFile *excelize.File
}

// NewExcelService
//
// f, _ := excelize.OpenFile(filepath)
//
// updateExcel := service.NewExcelService(f) OR service.NewExcelService(excelize.NewFile()) OR service.NewExcelService(nil)
func NewExcelService(excelFile *excelize.File) *ExcelService {
	if excelFile == nil {
		excelFile = excelize.NewFile()
	}
	return &ExcelService{ExcelFile: excelFile}
}

// SetRowData. rowi startbegin 1
func (e *ExcelService) SetRowData(sheetName string, data []interface{}, rowi int) {
	if rowi < 1 {
		panic("Error in ExcelService.SetRowData: arg rowi must greater than 0")
	}
	coli := 'A'
	for _, cellValue := range data {
		axis := fmt.Sprintf("%c%d", coli, rowi)
		if coli > 'Z' {
			add := coli - 'Z'
			coll := 'A' + add - 1
			axis = fmt.Sprintf("A%c%d", coll, rowi)
		}
		e.ExcelFile.SetCellValue(sheetName, axis, cellValue)
		coli++
	}
}

type ExcelImage struct {
	Url, LocalPath, Axis string
}

func (e *ExcelService) DownloadImagesByColly(sheetName, imgTitle, referer, dirname string, onResp func(furl string), withImgFile bool) error {
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

	// snames := []string{sheetName}
	// if sheetName == "" {
	// 	snames = e.ExcelFile.GetSheetList()
	// }

	snames := e.ExcelFile.GetSheetList()
	reqcount := 0
	for _, sn := range snames {
		if sheetName != "" && strings.TrimSpace(sn) != strings.TrimSpace(sheetName) {
			continue
		}
		imgColI, _, err := e.GetColsByTitle(sn, imgTitle)
		if err != nil {
			break
		}
		err = e.ReadRows(sn, func(rowData map[rune]string, rowI int) ReadRowResult {
			imgUrl := rowData[imgColI]
			if imgUrl == "" {
				log.Printf("-------Skip--DownloadImagesByColly--ReadRows--empty--imgUrl--sheetName(%s)--Coli(%d)----rowData(%+v)-------\n", sn, imgColI, rowData)
				return ReadRowResult{SkipAndContinue: true}
			}
			filepath := spd.GetLocalFile(imgUrl, dirname, LOCAL_IMAGE_FILE_EXT)
			isExist := miniutils.IsPathExists(filepath)
			if isExist {
				onResp(imgUrl)
				log.Printf("-----Skip--DownloadImagesByColly--ReadRows--sheetName(%s)--imgUrl(%s)--filepath(%s)--is exist---", sn, imgUrl, filepath)
				return ReadRowResult{SkipAndContinue: true}
			}
			log.Printf("---colly.Collector.Request(%s)\n", imgUrl)
			log.Printf("---Before Request:LocalPath(%s)\n", filepath)
			ctx := colly.NewContext()
			ctx.Put("LocalPath", filepath)
			reqcount++
			c.Request("GET", imgUrl, nil, ctx, nil)
			return ReadRowResult{Success: true}
		})

		if err != nil {
			break
		}
	}
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
			res := e.setLocalImages(sn, imgTitle, baseUrl, dirname, func(excelImg ExcelImage) ExcelImage {
				imgurl := excelImg.Url
				filepath := spd.GetLocalFile(imgurl, dirname, LOCAL_IMAGE_FILE_EXT)
				excelImg.LocalPath = filepath
				// log.Printf("---ReadRow--DownloadImagesByColly-debug-sheetName(%s)--imgurl(%s)--filepath(%s)----", sn, imgurl, filepath)
				if imgurl == "" {
					log.Println("Skip: DownloadImagesByColly Request: DownloadUrl is empty----------------")
					return excelImg
				}
				isExist := miniutils.IsPathExists(filepath)
				if isExist {
					log.Println("Skip: DownloadImagesByColly Request: filepath is exist---", filepath)
					return excelImg
				}
				return excelImg
			})
			if res != nil {
				result = res
			}
		}
	}
	// c.Wait() Wait Fail
	return result
}

type ReadRowResult struct {
	Success, SkipAndContinue, BreakEnd bool
}

func (e *ExcelService) ReadRows(sheetKey interface{}, callback func(rowData map[rune]string, rowI int) ReadRowResult) error {
	f := e.ExcelFile
	sheetName := ""
	var ok bool = true
	switch stn := sheetKey.(type) {
	case string:
		sheetName = stn
	case int:
		sheetName, ok = f.GetSheetMap()[stn]
	default:
		panic("sheetKey must be string | int")
	}
	fmt.Printf("-----ReadRows--sheetName(%s)\n", sheetName)
	if !ok {
		panic("index range out of SheetMap")
	}
	rows, err := f.Rows(sheetName)
	if err != nil {
		fmt.Printf("-----f.Rows--sheetName(%s)--err(%v)\n", sheetName, err)
		return err
	}
	rowI := 1
	for rows.Next() {
		fmt.Printf("---------------%d---\n", rowI)
		row, err := rows.Columns()
		if err != nil {
			// 遍历列失败
			return err
		}
		colI := 'A'
		rowString := ""
		rowData := make(map[rune]string, len(row))
		for _, cell := range row {
			rowData[colI] = strings.TrimSpace(cell)
			rowString += fmt.Sprintf("--%s%d: %v---", string(colI), rowI, cell)
			colI++
		}
		result := callback(rowData, rowI)
		if !result.Success {
			if result.BreakEnd {
				break
			}
			if result.SkipAndContinue {
				continue
			}
		}
		// log.Printf("\n---%d---"+rowString+"---\n", rowI)
		rowI++
	}
	return nil
}

func (e ExcelService) GetColsByTitle(sheetName, imgTitle string) (colIndex rune, dt []string, err error) {
	f := e.ExcelFile
	var cols *excelize.Cols
	cols, err = f.Cols(sheetName)
	if err != nil {
		return
	}
	colIndex = 'A'
	isok := false
	for cols.Next() {
		var col []string
		col, err = cols.Rows()
		if err != nil {
			// 遍历列失败
			break
		}
		// 跳过不是图片的数据列
		if col[0] != imgTitle {
			colIndex++
			continue
		} else {
			dt = col
			isok = true
			break
		}
	}
	if !isok {
		err = fmt.Errorf("找不到标题列:%s", imgTitle)
	}
	return
}

func (e *ExcelService) setLocalImages(sheetName, imgTitle, baseUrl, dirname string, callback func(excelImage ExcelImage) ExcelImage) error {
	f := e.ExcelFile
	cols, err := f.Cols(sheetName)
	if err != nil {
		return err
	}

	colI := 'A'
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
	imgTitleExist := false
	for cols.Next() {
		fmt.Println("------------------", string(colI))
		col, err := cols.Rows()
		if err != nil {
			// 遍历列失败
			return err
		}

		// 跳过不是图片的数据列
		if col[0] != imgTitle {
			colI++
			continue
		}
		imgTitleExist = true

		rowI := 1
		f.SetColWidth(sheetName, fmt.Sprintf("%c", colI), fmt.Sprintf("%c", colI), float64(colWidth))
		for i, cell := range col {
			f.SetRowHeight(sheetName, rowI, float64(colHeiht))

			if i == 0 {
				rowI++
				continue
			}
			fileurl := strings.TrimSpace(cell)
			fmt.Printf("\n----setLocalImages--row(%d)--fileurl(%s)----\n", rowI, fileurl)
			axis := fmt.Sprintf("%c%d", colI, rowI)
			excelImg := ExcelImage{Axis: axis, Url: fileurl}
			excelImg = callback(excelImg)
			excelImages = append(excelImages, excelImg)

			// f.SetCellValue(sheetName, axis, excelImg.LocalPath)
			rowI++
		}
	}
	if !imgTitleExist {
		log.Println("imgTitle is not Exist")
		return fmt.Errorf("excel文件里找不到标题名为%s的列，无法下载图片", imgTitle)
	}

	// AddPicture 不指定图片栏列宽度，图片无法填满整个单元格
	f.SetColWidth(sheetName, string(colI), string(colI), 11)

	for _, excelImg := range excelImages {
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

// TODO Wait download before save excel
// func (e *ExcelService) DownloadImages(sheetName, imgTitle, baseUrl, dirname string) error {
// 	spiderService := &SiteSpider{BaseUrl: baseUrl}
// 	var err error
// 	snames := []string{}
// 	if sheetName == "" {
// 		snames = e.ExcelFile.GetSheetList()
// 	} else {
// 		snames = append(snames, sheetName)
// 	}
// 	for _, sn := range snames {
// 		err = e.setLocalImages(sn, imgTitle, baseUrl, dirname, func(excelImg ExcelImage) ExcelImage {
// 			excelImg.LocalPath = spiderService.GetLocalFile(excelImg.Url, dirname, LOCAL_IMAGE_FILE_EXT)
// 			spiderService.DownloadFile(excelImg.Url, excelImg.LocalPath, "")
// 			return excelImg
// 		})
// 		if err != nil {
// 			break
// 		}
// 	}
// 	return err
// }

func (e *ExcelService) SaveAs(filepath string) error {
	defer e.ExcelFile.Close()
	return e.ExcelFile.SaveAs(filepath)
}

func (e *ExcelService) Save() error {
	defer e.ExcelFile.Close()
	return e.ExcelFile.Save()
}

// TODO Wait download before save excel
// --downloadimages=\"runtime/hello.xlsx,Sheet1,图片,https://www.amazon.com,amazon\"
func DownloadImagesByExcel(filepath, sheetName, imgTitle, referer, dirname string, withImgFile bool, onResp func(furl string)) error {
	// 下载图片并另存xlsx文件
	f, err := excelize.OpenFile(filepath)
	if err != nil {
		fmt.Println("error: ", err)
		return err
	}
	updateExcel := NewExcelService(f)
	if referer == "" {
		referer = "https://image.baidu.com/"
	}
	err = updateExcel.DownloadImagesByColly(sheetName, imgTitle, referer, dirname, onResp, withImgFile)
	// err = updateExcel.DownloadImages(sheetName, imgTitle, referer, dirname)
	if err != nil {
		fmt.Printf("-----DownloadImagesByExcel--err(%v)\n", err)
		return err
	}
	if withImgFile {
		return updateExcel.SaveAs(strings.Replace(filepath, ".xls", "_image.xls", 1))
	}
	return nil
}

func GetImgsByExcel(filepath, sheetName, imgTtile string) (imgs []string, err error) {
	var f *excelize.File
	f, err = excelize.OpenFile(filepath)
	if err != nil {
		fmt.Println("error: ", err)
		return
	}
	defer f.Close()
	ec := NewExcelService(f)

	// snames := []string{sheetName}
	// if sheetName == "" {
	// 	snames = ec.ExcelFile.GetSheetList()
	// }
	snames := ec.ExcelFile.GetSheetList()
	for _, sn := range snames {
		if sheetName != "" && strings.TrimSpace(sn) != strings.TrimSpace(sheetName) {
			continue
		}
		_, cols, errc := ec.GetColsByTitle(sn, imgTtile)
		if errc != nil {
			err = errc
			break
		}
		for i, col := range cols {
			if i == 0 {
				continue
			}
			imgs = append(imgs, col)
		}
	}
	return
}
