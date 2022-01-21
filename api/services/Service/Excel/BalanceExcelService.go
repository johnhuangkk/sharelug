package Excel

import (
	"api/services/VO/Response"
	"api/services/util/log"
	"api/services/util/tools"
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"reflect"
	"time"
)

func SellerBalanceNew() *excelFormats{
	fs := &excelFormats{}
	tag := []string{"A", "B", "C", "D", "E", "F"}
	title := []string{"賣家帳號(手機號碼)", "賣家會員代碼", "餘額(可提領餘額)", "待撥付餘額", "扣留餘額", "保留餘額"}
	column := []string{"Account", "SellerId", "Balance", "RetainBalance", "DetainBalance", "WithholdBalance"}
	for k, v := range tag {
		SetCell(v, column[k], title[k], fs)
	}
	return fs
}

func sellerBalanceHeader(f *excelize.File) error {
	header := []string {"賣家餘額表", fmt.Sprintf("製作日期：%s", time.Now().Format("2006/01/02"))}
	for k, v := range header {
		row := k+1
		if err := f.MergeCell("Sheet1", fmt.Sprintf("A%v", row),  fmt.Sprintf("I%v", row)); err != nil {
			log.Error("Merge Cell Error")
			return err
		}
		if err := f.SetCellValue("Sheet1", fmt.Sprintf("A%v", row), v); err != nil {
			log.Error("Set Cell Value Error")
		}
		style, err := f.NewStyle(`{"alignment":{"horizontal":"center"}}`)
		if err != nil {
			log.Error("Set Cell Value Error")
		}
		if err := f.SetCellStyle("Sheet1", fmt.Sprintf("A%v", row), fmt.Sprintf("I%v", row), style); err != nil {
			log.Error("Set Cell Value Error")
		}
	}
	return nil
}


func (fs excelFormats) ToSellerBalanceReportFile(report []Response.SellerBalanceResponse) (string, error) {
	f := excelize.NewFile()
	if err := sellerBalanceHeader(f); err != nil {
		log.Error("Set Header Error", err)
		return "", err
	}
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
			row := key + 4
			err := f.SetCellValue("Sheet1", fmt.Sprintf("%s%v", v.Tag, row), column.FieldByName(v.Column))
			if err != nil {
				log.Error("SetCellValue Error")
				return "", err
			}
		}
	}
	path := tools.GetFilePath("/temp/", "", 0)
	filename := fmt.Sprintf("%s%s.%s", path, time.Now().Format("Balance20060102150405"), "xlsx")
	if err := f.SaveAs(filename); err != nil {
		log.Error("Save xlsx Error", err)
	}
	return filename, nil
}

func BalanceNew() *excelFormats {
	fs := &excelFormats{}
	tag := []string{"A", "B", "C", "D", "E", "F"}
	title := []string{"日期", "項目", "存入金額", "支出金額", "餘額", "備註"}
	column := []string{"Date", "TransText", "In", "Out", "Balance", "Comment"}
	for k, v := range tag {
		SetCell(v, column[k], title[k], fs)
	}
	return fs
}

func (fs excelFormats) ToBalanceReportFile(report []Response.BalanceAccountList) (string, error) {
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
	filename := fmt.Sprintf("%s%s.%s", path, time.Now().Format("Balance20060102150405"), "xlsx")
	if err := f.SaveAs(filename); err != nil {
		log.Error("Save xlsx Error", err)
	}
	return filename, nil
}