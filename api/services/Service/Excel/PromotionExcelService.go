package Excel

import (
	"api/services/VO/ExcelVo"
	"api/services/util/log"
	"api/services/util/tools"
	"fmt"
	"reflect"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize/v2"
)

func CouponRecordNew() *excelFormats {
	fs := &excelFormats{}
	tag := []string{"A", "B", "C", "D", "E", "F", "G", "H"}
	title := []string{"序號", "折扣碼", "狀態", "訂單編號", "訂購日期", "訂購人手機", "訂單總金額", "折扣後金額"}
	column := []string{"Id", "Code", "Status", "OrderId", "TransTime", "BuyerPhone", "Amount", "DiscountAmount"}
	for k, v := range tag {
		SetCell(v, column[k], title[k], fs)
	}
	return fs
}

func (fs excelFormats) ToCouponUsedRecordFile(report []ExcelVo.CouponUsedReportVo, promoName string) (string, error) {
	f := excelize.NewFile()
	err := CouponUsedRecordHeader(f, promoName)
	if err != nil {
		log.Error("Set Header Value Error")
		return "", err
	}
	for _, v := range fs.excelFormats {
		err := f.SetCellValue("Sheet1", fmt.Sprintf("%s%v", v.Tag, 3), v.value)
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
	if err := f.SetColWidth("Sheet1", "F", "H", 11); err != nil {
		log.Error("SetCellValue Error")
	}
	if err := f.SetColWidth("Sheet1", "D", "E", 14); err != nil {
		log.Error("SetCellValue Error")
	}
	path := tools.GetFilePath("/temp/day/coupon", "", 0)
	filename := fmt.Sprintf("%s%s.%s", path, time.Now().Format("Coupon20060102"), "xlsx")
	if err := f.SaveAs(filename); err != nil {
		log.Error("Save excel Error", err)
	}
	return filename, nil
}

func CouponUsedRecordHeader(f *excelize.File, promoName string) error {
	header := []string{fmt.Sprintf("Check'Ne, 折扣碼使用表- %s", promoName),
		fmt.Sprintf("報表日期：%s (單位：新台幣)", time.Now().Format("2006/01/02"))}
	for k, v := range header {
		row := k + 1
		err := f.MergeCell("Sheet1", fmt.Sprintf("A%v", row), fmt.Sprintf("H%v", row))
		if err != nil {
			log.Error("Merge Cell Error")
			return err
		}
		err = f.SetCellValue("Sheet1", fmt.Sprintf("A%v", row), v)
		if err != nil {
			log.Error("Set Cell Value Error")
		}
		style, err := f.NewStyle(`{"alignment":{"horizontal":"center"}}`)
		err = f.SetCellStyle("Sheet1", fmt.Sprintf("A%v", row), fmt.Sprintf("H%v", row), style)
	}
	return nil
}
