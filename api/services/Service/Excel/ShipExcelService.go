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

func ShippingNew() *excelFormats {
	fs := &excelFormats{}
	tag := []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M"}
	title := []string{"序號", "訂單編號", "物流業者", "出貨單號", "商品名稱", "單價", "件數", "收件人-姓名", "收件人-手機", "收件人-郵遞區號", "收件人-地址", "備註", "買家備註"}
	column := []string{"Id", "OrderId", "ShipIdn", "ShipNumber", "ProductName", "Price", "Pieces", "ReceiverName", "ReceiverPhone", "ReceiverCode", "ReceiverAddress", "OrderMemo", "BuyerNotes"}
	for k, v := range tag {
		SetCell(v, column[k], title[k], fs)
	}
	return fs
}

func (fs excelFormats) ToShippingReportFile(report []ExcelVo.ShipReportVo) (string, error) {
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
	filename := fmt.Sprintf("%s%s.%s", path, time.Now().Format("checkne20060102150405"), "xlsx")
	if err := f.SaveAs(filename); err != nil {
		log.Error("Save xlsx Error", err)
	}
	return filename, nil
}

func ReadOrderShippingExcel(filename string) ([]ExcelVo.ShipReportVo, error) {
	var resp []ExcelVo.ShipReportVo
	path := fmt.Sprintf("./data/temp/%s", filename)
	f, err := excelize.OpenFile(path)
	if err != nil {
		log.Error("Open File Error", err.Error())
		return resp, err
	}
	list := f.GetSheetList()
	for _, sheet := range list {
		// 獲取 Sheet1 上所有單元格
		rows, err := f.GetRows(sheet)
		if err != nil {
			log.Error("Get Rows Error", err)
		}
		for k, row := range rows {
			if len(row) != 0 {
				if k != 0 && len(row[0]) != 0 {
					if tools.IsOrderId(row[1]) {
						var data ExcelVo.ShipReportVo
						data.Id = tools.StringToInt64(row[0])
						data.OrderId = row[1]
						data.ShipIdn = row[2]
						data.ShipNumber = row[3]
						data.ProductName = row[4]
						data.Price = row[5]
						data.Pieces = row[6]
						data.ReceiverName = row[7]
						data.ReceiverPhone = row[8]
						data.ReceiverCode = row[9]
						data.ReceiverAddress = row[10]
						resp = append(resp, data)
					}
				}
			}
		}
	}
	return resp, nil
}

func ReadSpecialStoreExcel(filename string) ([]ExcelVo.SpecialStoreReportVo, error) {
	var resp []ExcelVo.SpecialStoreReportVo
	path := fmt.Sprintf("./data/temp/%s", filename)
	f, err := excelize.OpenFile(path)
	if err != nil {
		log.Error("Open File Error", err.Error())
		return resp, err
	}
	list := f.GetSheetList()
	for _, sheet := range list {
		// 獲取 Sheet1 上所有單元格
		rows, err := f.GetRows(sheet)
		if err != nil {
			log.Error("Get Rows Error", err)
		}
		for k, row := range rows {
			if len(row) != 0 {
				if k != 0 && len(row[0]) != 0 {
					var data ExcelVo.SpecialStoreReportVo
					data.MerchantId = row[0]
					data.Terminal3DId = row[1]
					data.TerminalN3DId = row[2]
					data.ChStoreName = row[3]
					data.EnStoreName = row[4]
					resp = append(resp, data)
				}
			}
		}
	}
	return resp, nil
}

