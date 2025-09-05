package service

import (
	"fmt"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	"strings"

	// "github.com/iotames/downutils/conf"

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

// GetColDataBegin 获取数据列的内容
func (e ExcelService) GetColsBegin(sheetName string, colIndex, rowIndex int) (dt []string, err error) {
	f := e.ExcelFile
	var cols *excelize.Cols
	cols, err = f.Cols(sheetName)
	if err != nil {
		return
	}
	if colIndex == 0 {
		colIndex = 1
	}
	if rowIndex == 0 {
		rowIndex = 1
	}
	var colCount int = 1
	for cols.Next() {
		var col []string
		col, err = cols.Rows()
		if err != nil {
			// 遍历列失败
			fmt.Printf("----遍历列(%d)失败(%v)------\n", colCount, err)
			break
		}
		fmt.Printf("-------GetColsBegin------colCount%d == colIndex %d---------\n", colCount, colIndex)
		if colCount == colIndex {
			// 定位到开始的列
			fmt.Printf("-------colCount(%d)---rowIndex(%d)---col(%+v)--\n", colIndex, rowIndex, col)
			lencol := len(col)
			if rowIndex > lencol {
				err = fmt.Errorf("---定位行数第(%d)行开始，已超过最大行数（%d）-----", rowIndex, lencol)
				break
			}
			rowBegin := rowIndex - 1
			dt = col[rowBegin:]
			fmt.Printf("----GetColsBegin--rowBegin(%d)--dt(%v)-----\n", rowBegin, dt)
			break
		}
		colCount++
	}
	return
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
func DownloadImagesByExcel(filepath, sheetName, referer, dirname string, colIndex, rowIndex int, withImgFile bool, onResp func(furl string)) error {
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
	err = updateExcel.DownloadImagesByCollyLocation(sheetName, referer, dirname, colIndex, rowIndex, onResp, withImgFile)
	if err != nil {
		fmt.Printf("-----DownloadImagesByExcel--err(%v)\n", err)
		return err
	}
	if withImgFile {
		ffffff := strings.Replace(filepath, ".xls", "_image.xls", 1)
		fmt.Printf("-----------Save--Excel(%s)--------------\n", ffffff)
		return updateExcel.SaveAs(ffffff)
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

func GetImgsByExcelIndex(filepath, sheetName string, colIndex, rowIndex int) (imgs []string, err error) {
	var f *excelize.File
	f, err = excelize.OpenFile(filepath)
	if err != nil {
		fmt.Println("error: ", err)
		return
	}
	defer f.Close()
	ec := NewExcelService(f)
	snames := ec.ExcelFile.GetSheetList()
	for _, sn := range snames {
		if sheetName != "" && strings.TrimSpace(sn) != strings.TrimSpace(sheetName) {
			continue
		}
		cols, errc := ec.GetColsBegin(sn, colIndex, rowIndex)
		if errc != nil {
			err = errc
			break
		}
		for _, col := range cols {
			if strings.Index(col, "http") != 0 {
				fmt.Printf("----Skip----colData(%s)--NotBeginWith(http)-----\n", col)
				continue
			}
			imgs = append(imgs, col)
		}
	}
	return
}
