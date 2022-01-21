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

const (
	defaultValue   = ""
	officeLink     = "checkne.com"
	enable         = "Y"
	disable        = "N"
	signingfeeRate = "1.8"
	maxWidth       = 255
)

func SpecialStoreRecordNew() *excelFormats {
	fs := &excelFormats{}
	tag := []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z", "AA"}
	title := []string{"特店數", "特店代號", "端末機代號非3D", "端末機代號3D", "3D生效日", "EC簽約手續費率(%)", "MCC CODE",
		"區域代碼", "行業別代碼", "次特店統編", "次特店中文名稱", "次特店英文名稱", "次特店網址", "次特店營業(中文)地址", "次特店英文城市別",
		"次特店營業(英文)地址", "分期", "紅利", "3D", "登記名稱", "負責人姓名", "負責人姓氏(英文)", "負責人名(英文)", "負責人ID", "資本額(全數字)以萬元為單位", "設立日期(民國)(YYYMMDD)", "登記地址郵遞區號", "登記地址"}
	column := []string{"Id", "MerchantId", "TerminalId", "Terminal3DId", "3DTime", "Signingfee", "MccCode", "CityCode", "JobCode", "RepresentId", "StoreName", "StoreNameEn", "Link", "Addr", "CityName", "AddrEn", "Installments", "Bonus", "Enable3D", "CompanyName",
		"Represent", "RepresentLast", "RepresentFirst", "RepresentId", "Capital", "Establish", "ZipCode", "AddrEstablish"}
	for k, v := range tag {
		SetCell(v, column[k], title[k], fs)
	}
	return fs
}

func (fs excelFormats) ToSpecialStoreRecordFile(report []ExcelVo.SpecialStoreRecordVo) (string, string, error) {
	f := excelize.NewFile()
	err := SpecialStoreRecordHeader(f)
	if err != nil {
		log.Error("Set Header Value Error")
		return "", "", err
	}
	maxAddrWidth := 22
	maxAddrEnWidth := 22
	style, err := f.NewStyle(`{"alignment":{"horizontal":"center"}}`)
	for _, v := range fs.excelFormats {
		err := f.SetCellValue("Sheet1", fmt.Sprintf("%s%v", v.Tag, 3), v.value)
		if err != nil {
			log.Error("SetCellValue Error")
			return "", "", err
		}
		err = f.SetCellStyle("Sheet1", fmt.Sprintf("A%v", 3), fmt.Sprintf("AA%v", 3), style)
	}
	if err := f.SetColWidth("Sheet1", "A", "A", 7); err != nil {
		log.Error("Set Special Store Column Width Value Error", err)
	}
	if err := f.SetColWidth("Sheet1", "B", "X", 22); err != nil {
		log.Error("Set Special Store Column Width Value Error", err)
	}
	if err := f.SetColWidth("Sheet1", "Y", "AA", 32); err != nil {
		log.Error("Set Special Store Column Width Value Error", err)
	}
	for key, value := range report {

		column := reflect.ValueOf(value)
		if (len([]rune(value.Addr)) * 12 / 6) > maxAddrWidth {
			maxAddrWidth = len([]rune(value.Addr)) * 12 / 6

		}
		if (len([]rune(value.AddrEn)) * 12 / 6) > maxAddrEnWidth {
			maxAddrEnWidth = len([]rune(value.AddrEn)) * 12 / 6
		}

		var err error
		for _, v := range fs.excelFormats {
			row := key + 4
			switch v.Column {
			case "3DTime":
				err = f.SetCellValue("Sheet1", fmt.Sprintf("%s%v", v.Tag, row), defaultValue)
			case "Signingfee":
				err = f.SetCellValue("Sheet1", fmt.Sprintf("%s%v", v.Tag, row), signingfeeRate)
			case "Installments":
				err = f.SetCellValue("Sheet1", fmt.Sprintf("%s%v", v.Tag, row), disable)
			case "Link":
				err = f.SetCellValue("Sheet1", fmt.Sprintf("%s%v", v.Tag, row), officeLink)
			case "Bonus":
				err = f.SetCellValue("Sheet1", fmt.Sprintf("%s%v", v.Tag, row), disable)
			case "Enable3D":
				err = f.SetCellValue("Sheet1", fmt.Sprintf("%s%v", v.Tag, row), enable)
			case "Capital":
				err = f.SetCellValue("Sheet1", fmt.Sprintf("%s%v", v.Tag, row), defaultValue)
			case "Establish":
				err = f.SetCellValue("Sheet1", fmt.Sprintf("%s%v", v.Tag, row), defaultValue)
			case "ZipCode":
				err = f.SetCellValue("Sheet1", fmt.Sprintf("%s%v", v.Tag, row), defaultValue)
			case "AddrEstablish":
				err = f.SetCellValue("Sheet1", fmt.Sprintf("%s%v", v.Tag, row), defaultValue)
			case "CompanyName":
				err = f.SetCellValue("Sheet1", fmt.Sprintf("%s%v", v.Tag, row), defaultValue)

			default:
				err = f.SetCellValue("Sheet1", fmt.Sprintf("%s%v", v.Tag, row), column.FieldByName(v.Column))
			}

			if err != nil {
				log.Error("SetCellValue Error")
				continue
			}
			err = f.SetCellStyle("Sheet1", fmt.Sprintf("A%v", row), fmt.Sprintf("AA%v", row), style)
		}
	}
	if maxAddrWidth > maxWidth {
		maxAddrWidth = maxWidth
	}
	if maxAddrEnWidth > maxWidth {
		maxAddrEnWidth = maxWidth
	}
	if err := f.SetColWidth("Sheet1", "N", "N", float64(maxAddrWidth)); err != nil {
		log.Error("Set Special Store Column Width Value Error", err)
	}
	if err := f.SetColWidth("Sheet1", "P", "P", float64(maxAddrEnWidth)); err != nil {
		log.Error("Set Special Store Column Width Value Error", err)
	}
	// filename := "test"
	path := tools.GetFilePath("/temp/day/specialstore", "", 0)
	filename := fmt.Sprintf("%s.%s", time.Now().Format("KgiSpecialStore_20060102"), "xlsx")
	fileWithPath := fmt.Sprintf("%s%s", path, filename)
	if err := f.SaveAs(fileWithPath); err != nil {
		log.Error("Save excel Error", err)
	}
	return filename, path, nil
}

func SpecialStoreRecordHeader(f *excelize.File) error {
	header := []string{fmt.Sprintf("Check'Ne, 特店資料表"),
		fmt.Sprintf("報表日期：%s", time.Now().Format("2006/01/02"))}
	for k, v := range header {
		row := k + 1
		err := f.MergeCell("Sheet1", fmt.Sprintf("A%v", row), fmt.Sprintf("AA%v", row))
		if err != nil {
			log.Error("Merge Cell Error")
			return err
		}
		err = f.SetCellValue("Sheet1", fmt.Sprintf("A%v", row), v)
		if err != nil {
			log.Error("Set Cell Value Error")
		}
		style, err := f.NewStyle(`{"alignment":{"horizontal":"center"}}`)
		err = f.SetCellStyle("Sheet1", fmt.Sprintf("A%v", row), fmt.Sprintf("AA%v", row), style)
	}
	return nil
}
