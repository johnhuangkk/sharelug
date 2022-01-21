package Excel

import (
	"api/services/Enum"
	"api/services/VO/KgiBank"
	"api/services/dao/Credit"
	"api/services/dao/Withdraw"
	"api/services/dao/member"
	"api/services/entity"
	"api/services/util/log"
	"api/services/util/tools"
	"fmt"
	"reflect"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize/v2"
)

type OrderReport struct {
	Id            string
	OrderTime     string
	OrderId       string
	ProductName   string
	ProductSpec   string
	PaymentTime   string
	PaymentType   string
	OrderStatus   string
	ShipTime      string
	ShipType      string
	ShipNumber    string
	OrderAmount   int64
	ProductAmount int64
	ShipFee       int64
	PlatformFee   int64
	CreditAmount  int64
}

type excelFormat struct {
	Tag    string
	Column string
	value  string
}

type excelFormats struct {
	excelFormats []excelFormat
}

func New() *excelFormats {
	fs := &excelFormats{}
	tag := []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P"}
	title := []string{"序號", "訂購日期", "訂單編號", "商品名稱", "規格",
		"付款日期", "付款方式", "訂單狀態", "出貨日期", "出貨方式", "出貨單號",
		"訂單金額", "商品總金額", "買方支付運費", "平台費用", "入帳金額"}
	column := []string{"Id", "OrderTime", "OrderId", "ProductName", "ProductSpec", "PaymentTime", "PaymentType", "OrderStatus",
		"ShipTime", "ShipType", "ShipNumber", "OrderAmount", "ProductAmount", "ShipFee", "PlatformFee", "CreditAmount"}

	for k, v := range tag {
		SetCell(v, column[k], title[k], fs)
	}
	return fs
}

func SetCell(tag, column, value string, fs *excelFormats) {
	format := SetValue(tag, column, value)
	fs.excelFormats = append(fs.excelFormats, format)
}

func SetValue(tag, column, value string) excelFormat {
	var ExcelFormat excelFormat
	ExcelFormat.Tag = tag
	ExcelFormat.Column = column
	ExcelFormat.value = value
	return ExcelFormat
}

func (fs excelFormats) ToReportFile(report []OrderReport, StartTime, EndTime string) (string, error) {
	f := excelize.NewFile()
	err := ReportHeader(f, StartTime, EndTime)
	if err != nil {
		log.Error("Set Header Value Error")
		return "", err
	}

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
			row := key + 7
			err := f.SetCellValue("Sheet1", fmt.Sprintf("%s%v", v.Tag, row), column.FieldByName(v.Column))
			if err != nil {
				log.Error("SetCellValue Error")
				return "", err
			}
		}
	}
	path := tools.GetFilePath("/temp/", "", 0)
	filename := fmt.Sprintf("%s%s.%s", path, time.Now().Format("20060102150405"), "xlsx")
	if err := f.SaveAs(filename); err != nil {
		log.Error("Save xlsx Error", err)
	}
	return filename, nil
}

func ReportHeader(f *excelize.File, StartTime, EndTime string) error {
	header := []string{"Check'Ne", "收銀機銷售明細表", fmt.Sprintf("銷售期間：%s ~ %s", StartTime, EndTime),
		fmt.Sprintf("報表日期：%s", time.Now().Format("2006/01/02")), "*僅包含已完成付款之訂單，且不含訂單因故取消、未取貨、退貨、或退款等情形。**單位：新台幣"}
	for k, v := range header {
		row := k + 1
		err := f.MergeCell("Sheet1", fmt.Sprintf("A%v", row), fmt.Sprintf("O%v", row))
		if err != nil {
			log.Error("Merge Cell Error")
			return err
		}
		err = f.SetCellValue("Sheet1", fmt.Sprintf("A%v", row), v)
		if err != nil {
			log.Error("Set Cell Value Error")
		}
		style, err := f.NewStyle(`{"alignment":{"horizontal":"center"}}`)
		err = f.SetCellStyle("Sheet1", fmt.Sprintf("A%v", row), fmt.Sprintf("O%v", row), style)
	}
	return nil
}

func ReadBinCodeExcel() {
	f, err := excelize.OpenFile("./data/bin-202103181638.xlsx")
	if err != nil {
		log.Error("Open File Error", err.Error())
		return
	}
	// 獲取 Sheet1 上所有單元格
	rows, err := f.GetRows("BinTable")
	for k, row := range rows {
		if k != 0 && k != 1 {
			if len(row) == 0 {
				return
			}
			if len(row) > 5 {
				var data entity.BankBinCode
				data.BankName = row[1]
				data.CardType = row[3]
				data.BinNumber = row[4]
				data.IsDebit = 0
				if row[9] == "Y" {
					data.IsDebit = 1
				}
				data.Status = 1
				//log.Debug("row", row[1], row[3], row[4], row[9])
				if err := Credit.InsertBankBinCodeData(data); err != nil {
					log.Error("Insert Bank bin Code Error", err.Error())
					return
				}
			}
		}
	}
}

func ReadDonateExcel() {
	f, err := excelize.OpenFile("./data/donate210415.xlsx")
	if err != nil {
		log.Error("Open File Error", err.Error())
		return
	}
	// 獲取 Sheet1 上所有單元格
	rows, err := f.GetRows("Sheet1")
	for k, row := range rows {
		if len(row) == 0 {
			return
		}
		if k != 0 {
			var data entity.DonateData
			data.DonateName = row[1]
			data.DonateCode = row[2]
			data.DonateShort = row[3]
			data.DonateBan = row[4]
			data.DonateCity = row[5]
			if data.DonateCode == "14697346" || data.DonateCode == "13579" || data.DonateCode == "885521" || data.DonateCode == "5299" || data.DonateCode == "119" {
				data.DonateStatus = Enum.OrderSuccess
			} else {
				data.DonateStatus = Enum.OrderFail
			}
			err := member.InsertDonateData(data)
			if err != nil {
				log.Error("Insert Bank Code Error", err.Error())
				return
			}
		}
	}
}

func ReadAchExcel() {
	f, err := excelize.OpenFile("ACH1214.xlsx")
	if err != nil {
		log.Error("Open File Error", err.Error())
		return
	}
	// 獲取 Sheet1 上所有單元格
	rows, err := f.GetRows("Sheet1")
	for _, row := range rows {
		if len(row) == 0 {
			return
		}
		var data entity.BankCodeData
		data.BankCode = row[0][:3]
		data.BankName = row[1]
		data.BranchCode = row[0][:7]
		data.BankStatus = Enum.OrderSuccess
		err := Withdraw.InsertBankCode(data)
		if err != nil {
			log.Error("Insert Bank Code Error", err.Error())
			return
		}
	}
}

func ReadAchResponseExcel(filename string) ([]KgiBank.AchBody, error) {
	var resp []KgiBank.AchBody
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
			return resp, err
		}
		for k, row := range rows {
			log.Debug("row", row, len(row))
			if k != 0 && len(row[0]) != 0 {
				var data KgiBank.AchBody
				data.TYPE = row[0]
				data.TXTYPE = row[1]
				data.TXID = row[2]
				data.SEQ = row[3]
				data.PBANK = row[4]
				data.PCLNO = row[5]
				data.RBANK = row[6]
				data.RCLNO = row[7]
				data.AMT = row[8]
				data.RCODE = row[9]
				data.SCHD = row[10]
				data.CID = row[11]
				data.PID = row[12]
				data.SID = row[13]
				data.PDATE = row[14]
				data.PSEQ = row[15]
				data.PSCHD = row[16]
				data.CNO = row[17]
				data.NOTE = row[18]
				data.MEMO = row[19]
				data.CFEE = row[20]
				if len(row) > 22 {
					data.NOTEB = row[21]
					data.FILLER = row[22]
				}
				resp = append(resp, data)
			}
		}
	}
	return resp, nil
}

func ReadIndustryExcel() {
	f, err := excelize.OpenFile("./data/temp/Industry.xlsx")
	if err != nil {
		log.Error("Open File Error", err.Error())
		return
	}
	cat := map[string]int{
		"百貨業":   11010,
		"服飾業":   12010,
		"食品業":   13010,
		"餐飲業":   14010,
		"飲料小吃業": 15010,
		"旅宿業":   16010,
		"裝潢業":   17010,
		"電器業":   18010,
		"運輸服務業": 19010,
		"醫療相關":  20010,
		"文教類":   21010,
		"電腦業":   22010,
		"娛樂類":   23010,
		"資訊服務":  24010,
		"運動類":   25010,
		"禮品業":   26010,
		"精密儀器":  27010,
		"美容業":   28010,
		"工程業":   29010,
		"商業服務":  30010,
		"其它類":   31010,
	}
	// 獲取 Sheet1 上所有單元格
	rows, err := f.GetRows("Sheet1")
	for k, row := range rows {
		if k != 0 && len(row[0]) != 0 {
			i, _ := cat[row[2]]
			var data entity.IndustryData
			data.IndustryId = row[0]
			data.Mcc = row[1]
			data.Category = tools.Nl2br(row[2])
			data.Industry = tools.Nl2br(row[3])
			data.Sort = tools.IntToString(i)
			if err := Withdraw.InsertIndustryData(data); err != nil {
				log.Error("Insert Bank Code Error", err.Error())
				return
			}
			cat[row[2]]++
		}
	}
}
