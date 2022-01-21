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

func BankNew() *excelFormats {
	fs := &excelFormats{}
	tag := []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M"}
	title := []string{"賣家帳號", "會員代碼", "種類", "姓名", "證號", "負責人", "負責人證號", "地址", "收銀機名稱1", "收銀機名稱2", "收銀機名稱3", "收銀機名稱4", "收銀機名稱5"}
	column := []string{"Account", "TerminalId", "Category", "Name", "Identity", "Head", "HeadIdentity", "CompanyAddr", "Store1", "Store2", "Store3", "Store4", "Store5"}
	for k, v := range tag {
		SetCell(v, column[k], title[k], fs)
	}
	return fs
}

func MemberNew() *excelFormats {
	fs := &excelFormats{}
	tag := []string{"A", "B", "C"}
	title := []string{"會員帳號", "會員代碼", "會員餘額"}
	column := []string{"MemberId", "TerminalId", "Balance"}
	for k, v := range tag {
		SetCell(v, column[k], title[k], fs)
	}
	return fs
}

func OrderNew() *excelFormats {
	fs := &excelFormats{}
	tag := []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L"}
	title := []string{"序號", "交易日期", "付款方式", "付款日期", "出貨方式", "出貨時間", "訂單編號", "賣方會員代碼", "買方會員代碼", "交易金額", "平台費用", "是否已付平台費用"}
	column := []string{"Id", "OrderTime", "PayWay", "PayWayTime", "ShipType", "ShipTime", "OrderId", "SellerId", "BuyerId", "Amount", "PlatformFee", "IsFee"}
	for k, v := range tag {
		SetCell(v, column[k], title[k], fs)
	}
	return fs
}

//
func (fs excelFormats) ToMemberReportFile(report []ExcelVo.MemberReportVo) (string, error) {
	f := excelize.NewFile()
	for _, v := range fs.excelFormats {
		err := f.SetCellValue("Sheet1", fmt.Sprintf("%s%v", v.Tag, 6), v.value)
		if err != nil {
			log.Error("SetCellValue Error")
			return "", err
		}
	}
	for key, value := range report {
		column := reflect.ValueOf(value)
		for _, v := range fs.excelFormats {
			row := key+7
			log.Debug("row", row)
			err := f.SetCellValue("Sheet1", fmt.Sprintf("%s%v", v.Tag, row), column.FieldByName(v.Column))
			if err != nil {
				log.Error("SetCellValue Error")
				return "", err
			}
		}
	}
	path := tools.GetFilePath("/temp/", "", 0)
	filename := fmt.Sprintf("%s%s.%s", path, time.Now().Format("MEM20060102150405"), "xlsx")
	if err := f.SaveAs(filename); err != nil {
		log.Error("Save xlsx Error", err)
	}
	return filename, nil
}

func (fs excelFormats) ToOrderReportFile(report []ExcelVo.OrderReportVo) (string, error) {
	f := excelize.NewFile()
	for _, v := range fs.excelFormats {
		err := f.SetCellValue("Sheet1", fmt.Sprintf("%s%v", v.Tag, 6), v.value)
		if err != nil {
			log.Error("SetCellValue Error")
			return "", err
		}
	}
	for key, value := range report {
		column := reflect.ValueOf(value)
		for _, v := range fs.excelFormats {
			row := key+7
			//log.Debug("row", row)
			err := f.SetCellValue("Sheet1", fmt.Sprintf("%s%v", v.Tag, row), column.FieldByName(v.Column))
			if err != nil {
				log.Error("SetCellValue Error")
				return "", err
			}
		}
	}
	path := tools.GetFilePath("/temp/", "", 0)
	filename := fmt.Sprintf("%s%s.%s", path, time.Now().Format("Order20060102150405"), "xlsx")
	if err := f.SaveAs(filename); err != nil {
		log.Error("Save xlsx Error", err)
	}
	return filename, nil
}

func (fs excelFormats) ToBankReportFile(report []ExcelVo.BankReportVo) (string, error) {
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
			row := key+2
			//log.Debug("row", row)
			err := f.SetCellValue("Sheet1", fmt.Sprintf("%s%v", v.Tag, row), column.FieldByName(v.Column))
			if err != nil {
				log.Error("SetCellValue Error")
				return "", err
			}
		}
	}
	path := tools.GetFilePath("/temp/", "", 0)
	filename := fmt.Sprintf("%s%s.%s", path, time.Now().Format("Bank20060102150405"), "xlsx")
	if err := f.SaveAs(filename); err != nil {
		log.Error("Save xlsx Error", err)
	}
	return filename, nil
}