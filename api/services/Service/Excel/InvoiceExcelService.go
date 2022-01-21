package Excel

import (
	"api/services/VO/ExcelVo"
	"api/services/util/log"
	"api/services/util/tools"
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"reflect"
	"time"
)

func UserInvoiceNew() *excelFormats {
	fs := &excelFormats{}
	tag := []string{"A", "B", "C", "D"}
	title := []string{"訂單編號", "發票號碼", "發票日期", "發票金額"}
	column := []string{"OrderId", "InvoiceNumber", "CreateTime", "Amount"}
	for k, v := range tag {
		SetCell(v, column[k], title[k], fs)
	}
	return fs
}

func (fs excelFormats) ToInvoiceExcelFile(report []ExcelVo.UserInvoiceReportVo) (string, error) {
	f := excelize.NewFile()
	for _, v := range fs.excelFormats {
		err := f.SetCellValue("Sheet1", fmt.Sprintf("%s%v", v.Tag, 1), v.value)
		if err != nil {
			log.Error("SetCellValue Error")
			return "", err
		}
	}
	for key, value := range report {
		column := reflect.ValueOf(value)
		for _, v := range fs.excelFormats {
			row := key + 2
			err := f.SetCellValue("Sheet1", fmt.Sprintf("%s%v", v.Tag, row), column.FieldByName(v.Column))
			if err != nil {
				log.Error("SetCellValue Error")
				return "", err
			}
		}
	}
	path := tools.GetFilePath("/temp/", "", 0)
	filename := fmt.Sprintf("%s%s.%s", path, time.Now().Format("checkne_20060102150405"), "xlsx")
	if err := f.SaveAs(filename); err != nil {
		log.Error("Save xlsx Error", err)
	}
	return filename, nil
}



func InvoiceNew() *excelFormats {
	fs := &excelFormats{}
	tag := []string{"A", "B", "C", "D", "E", "F", "G", "H", "I"}
	title := []string{"發票年月", "發票號碼", "開立時間", "購買內容", "銷售額", "稅額", "發票金額", "公司統編", "發票狀態"}
	column := []string{"InvoiceYm", "InvoiceNumber", "CreateTime", "Products", "Sales", "Tax", "Amount", "Identifier", "InvoiceStatus"}
	for k, v := range tag {
		SetCell(v, column[k], title[k], fs)
	}
	return fs
}

func (fs excelFormats) ToExcelReportFile(report []ExcelVo.InvoiceReportVo) (string, error) {
	f := excelize.NewFile()
	for _, v := range fs.excelFormats {
		err := f.SetCellValue("Sheet1", fmt.Sprintf("%s%v", v.Tag, 1), v.value)
		if err != nil {
			log.Error("SetCellValue Error")
			return "", err
		}
	}
	for key, value := range report {
		column := reflect.ValueOf(value)
		for _, v := range fs.excelFormats {
			row := key + 2
			err := f.SetCellValue("Sheet1", fmt.Sprintf("%s%v", v.Tag, row), column.FieldByName(v.Column))
			if err != nil {
				log.Error("SetCellValue Error")
				return "", err
			}
		}
	}
	path := tools.GetFilePath("/temp/", "", 0)
	filename := fmt.Sprintf("%s%s.%s", path, time.Now().Format("checkne_20060102150405"), "xlsx")
	if err := f.SaveAs(filename); err != nil {
		log.Error("Save xlsx Error", err)
	}
	return filename, nil
}

