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

func DayStatementNew() *excelFormats {
	fs := &excelFormats{}
	tag := []string{"A", "B", "C", "D", "E", "F", "G", "H", "I"}
	title := []string{"序號", "交易日期", "入金", "付款方式", "訂單編號", "賣方會員代碼", "買方會員代碼", "交易金額", "平台費用"}
	column := []string{"Id", "TransactionDate", "TransactionType", "PaymentType", "OrderId", "SellerId", "BuyerId", "Amount", "PlatformFee"}
	for k, v := range tag {
		SetCell(v, column[k], title[k], fs)
	}
	return fs
}

func (fs excelFormats) ToDayStatementReportFile(report []ExcelVo.DayStatementVo, StartTime, EndTime string) (string, error) {
	f := excelize.NewFile()
	err := DayStatementHeader(f, StartTime, EndTime)
	if err != nil {
		log.Error("Set Header Value Error")
		return "", err
	}
	for _, v := range fs.excelFormats {
		err := f.SetCellValue("Sheet1", fmt.Sprintf("%s%v", v.Tag, 5), v.value)
		if err != nil {
			log.Error("SetCellValue Error")
			return "", err
		}
	}
	for key, value := range report {
		column := reflect.ValueOf(value)
		for _, v := range fs.excelFormats {
			row := key+6
			log.Debug("row", row)
			err := f.SetCellValue("Sheet1", fmt.Sprintf("%s%v", v.Tag, row), column.FieldByName(v.Column))
			if err != nil {
				log.Error("SetCellValue Error")
				return "", err
			}
		}
	}
	day, _ := time.Parse("2006/01/02", EndTime)
	path := tools.GetFilePath("/temp/day/", "", 0)
	filename := fmt.Sprintf("%s%s.%s", path, day.Format("Day20060102"), "xlsx")
	if err := f.SaveAs(filename); err != nil {
		log.Error("Save excel Error", err)
	}
	return filename, nil
}

func DayStatementHeader(f *excelize.File, StartTime, EndTime string) error {
	header := []string {"Check'Ne", "日結表-凱基信託帳戶", fmt.Sprintf("期間：%s ~ %s", StartTime, EndTime),
		fmt.Sprintf("報表日期：%s (單位：新台幣)", time.Now().Format("2006/01/02"))}
	for k, v := range header {
		row := k+1
		err := f.MergeCell("Sheet1", fmt.Sprintf("A%v", row),  fmt.Sprintf("I%v", row))
		if err != nil {
			log.Error("Merge Cell Error")
			return err
		}
		err = f.SetCellValue("Sheet1", fmt.Sprintf("A%v", row), v)
		if err != nil {
			log.Error("Set Cell Value Error")
		}
		style, err := f.NewStyle(`{"alignment":{"horizontal":"center"}}`)
		err = f.SetCellStyle("Sheet1", fmt.Sprintf("A%v", row), fmt.Sprintf("I%v", row), style)
	}
	return nil
}
