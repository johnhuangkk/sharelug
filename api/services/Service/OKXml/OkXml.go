package OKXml

import (
	"api/services/Enum"
	"api/services/Service/CvsShipping"
	om "api/services/Service/OKMart"
	"api/services/dao/Cvs"
	"api/services/dao/Mart"
	"api/services/database"
	"api/services/entity"
	"api/services/util/log"
	"api/services/util/tools"
	tw_address "api/services/util/tw-address"
	"bufio"
	"bytes"
	"encoding/xml"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type F01Xml struct {
	Xml    xml.Name  `xml:"F01DOC"`
	Header F01Header `xml:"DOCHEAD"`
	F01    []F01     `xml:"F01CONTENT"`
}

type F01Header struct {
	DOCDATE         string `xml:"DOCDATE"`
	FROMPARTNERCODE string `xml:"FROMPARTNERCODE"`
	TOPARTENERCODE  string `xml:"TOPARTENERCODE"`
}
type F01 struct {
	StoreId      string `xml:"STNO"`
	StoreName    string `xml:"STNM"`
	StoreTel     string `xml:"STTEL"`
	StoreCity    string `xml:"STCITY"`
	StoreDisct   string `xml:"STCNTRY"`
	StoreAddress string `xml:"STADR"`
	Zipcode      string `xml:"ZIPCD"`
	DCRONO       string `xml:"DCRONO"`
	SDATE        string `xml:"SDATE"`
	EDATE        string `xml:"EDATE"`
}

func getOkStore(engine *database.MysqlSession) {
	var client om.Client
	client.GetClient()

	if client.ConnectFtpError != nil {
		return
	}

	defer client.ConnectFtp.Quit()

	t := "/461/F01"

	fileNames, err := client.FetchFolder(t)
	if err != nil {
		return
	}

	for _, f := range fileNames {
		data, err := client.RetrFile(t, f)
		if err != nil {
			return
		}
		// content, err := client.GetF01Xml(data)
		var xmlData F01Xml

		err = xml.Unmarshal(data, &xmlData)
		if err != nil {
			return
		}

		var stores []entity.MartOkStoreData
		for _, data := range xmlData.F01 {
			country, district, address, flag := tw_address.AddressSplit(data.StoreAddress)
			if flag {
				store := entity.MartOkStoreData{
					StoreId:        data.StoreId,
					StoreName:      data.StoreName,
					StoreAddress:   address,
					StoreCloseDate: data.SDATE,
					MdcStareDate:   data.SDATE,
					MdcEndDate:     data.EDATE,
					TelNo:          data.StoreTel,
					City:           country,
					District:       district,
				}
				stores = append(stores, store)

				exists := new(entity.MartOkStoreData)

				query := map[string]interface{}{}
				query["store_id"] = data.StoreId
				exist, err := engine.Engine.Where(query).Exist(exists)

				if err != nil {
					log.Debug("MartOKFetchStoreList:%v", err)
					continue
				}

				if exist {
					_, err := engine.Session.Update(store, entity.MartOkStoreData{StoreId: data.StoreId})
					if err != nil {
						log.Debug("MartOKFetchStoreList:%v", err)
					}
				} else {
					_, err := engine.Session.InsertOne(store)
					if err != nil {
						log.Debug("MartOKFetchStoreList:%v", err)
					}
				}
			}

		}
		client.MoveFileToBackup(t, f)
	}
}

func MartOKFetchStoreList() {

	engine := database.GetMysqlEngine()
	defer engine.Close()

	getOkStore(engine)

	var newStores []entity.MartOkStoreData
	err := engine.Engine.Find(&newStores)
	if err != nil {
		log.Debug("MartOkStoreData:", err)
		return
	}

	_ = Mart.WriteOkStoreData(newStores)
}

func okDateTimeFormat(dataTime string) string {
	t, _ := time.Parse(`20060102150405`, dataTime)
	return t.Format(`2006-01-02 15:04:05`)
}

// ＯＫ帳務
func MarkOkAccounting()  {
	log.Debug("MartOkStoreData")
	var client om.Client
	client.GetClient()

	if client.ConnectFtpError != nil {
		return
	}

	defer client.ConnectFtp.Quit()

	engine := database.GetMysqlEngine()
	defer engine.Close()

	MartOKFetchAcct(client, engine)
}

func MartOkFetchShipping() {
	log.Debug("MartOkStoreData")
	var client om.Client
	client.GetClient()

	if client.ConnectFtpError != nil {
		return
	}

	defer client.ConnectFtp.Quit()

	engine := database.GetMysqlEngine()
	defer engine.Close()

	ecNo1 := client.EcNo1
	{
		// 即時寄件代收檔
		MartOKFetchF27(ecNo1, client, engine)
		// 寄件代收檔 (00:30)
		MartOKFetchF25(ecNo1, client, engine)
		// 離店檔
		MartOKFetchF84(ecNo1, client, engine)
		// 大物流驗收檔
		MartOKFetchF71(ecNo1, client, engine)
	}

	ecNo2 := client.EcNo2
	{
		// 小物流驗收檔
		MartOKFetchF63(ecNo2, client, engine)
		// 即時進店檔
		MartOKFetchF44(ecNo2, client, engine)
		// 進店檔(13:30)
		MartOKFetchF64(ecNo2, client, engine)
		// 即時取貨代收檔
		MartOKFetchF17(ecNo2, client, engine)
		// 取貨代收檔(00 : 45)
		MartOKFetchF65(ecNo2, client, engine)
		// 未取離店
		MartOKFetchF84(ecNo2, client, engine)
		// 物流驗退檔
		MartOKFetchF67(ecNo2, client, engine)
	}

	ecNo3 := client.EcNo3
	{
		// 重出物流驗收檔
		MartOKFetchF03(ecNo3, client, engine)
		// 即時進店檔
		MartOKFetchF44(ecNo3, client, engine)
		// 重出進店檔(13:30)
		MartOKFetchF04(ecNo3, client, engine)
		// 即時取貨代收檔
		MartOKFetchF17(ecNo3, client, engine)
		// 重出取貨完成檔
		MartOKFetchF05(ecNo3, client, engine)
		// 離店檔
		MartOKFetchF84(ecNo3, client, engine)
		// 重出物流驗退檔
		MartOKFetchF07(ecNo3, client, engine)
	}

	ecNo4 := client.EcNo4
	{
		// 重出物流驗收檔
		MartOKFetchF03(ecNo4, client, engine)
		// 即時進店檔
		MartOKFetchF44(ecNo4, client, engine)
		// 重出進店檔(13:30)
		MartOKFetchF04(ecNo4, client, engine)
		// 即時取貨代收檔
		MartOKFetchF17(ecNo4, client, engine)
		// 重出取貨完成檔
		MartOKFetchF05(ecNo4, client, engine)
		// 離店檔
		MartOKFetchF84(ecNo4, client, engine)
		// 重出物流驗退檔
		MartOKFetchF07(ecNo4, client, engine)
	}

	ecNo5 := client.EcNo5
	{
		MartOKFetchF03(ecNo5, client, engine)
		MartOKFetchF07(ecNo5, client, engine)
		MartOKFetchF44(ecNo5, client, engine)
		MartOKFetchF04(ecNo5, client, engine)
		MartOKFetchF17(ecNo5, client, engine)
		MartOKFetchF05(ecNo5, client, engine)
		MartOKFetchF84(ecNo5, client, engine)
		MartOKFetchF07(ecNo5, client, engine)
	}
}

// 即時寄件代收檔
func MartOKFetchF27(ecNo string, client om.Client, engine *database.MysqlSession) {

	t := "/" + ecNo + "/F27"

	fileNames, err := client.FetchFolder(t)
	if err != nil {
		return
	}

	var UpdateCvsShipping CvsShipping.UpdateCvsShipping
	UpdateCvsShipping.Type = `F27`
	UpdateCvsShipping.ShipType = Enum.CVS_OK_MART

	for _, f := range fileNames {
		data, err := client.RetrFile(t, f)
		if err != nil {
			return
		}
		content, err := client.GetF27Xml(data)
		log.Info(t+`:`+f+`content:`, content)
		if err != nil {
			return
		}
		for _, d := range content.Body {
			UpdateCvsShipping.ShipNo = d.OrderNo
			UpdateCvsShipping.DateTime = okDateTimeFormat(d.UpDateTime)
			UpdateCvsShipping.DetailStatus = ""
			UpdateCvsShipping.Log = tools.XmlToString(d)
			UpdateCvsShipping.FileName = f

			err = UpdateCvsShipping.UpdateCvsShippingShipment(engine)

			if err != nil {
				log.Error(t+":"+f+" Error Data [%s]:", d)
				log.Error(t+" Error [%s]:", err.Error())
			}
		}
		client.MoveFileToBackup(t, f)
	}
}

// 寄件代收檔 (日）
func MartOKFetchF25(ecNo string, client om.Client, engine *database.MysqlSession) {

	t := "/" + ecNo + "/F25"

	fileNames, err := client.FetchFolder(t)
	if err != nil {
		return
	}

	var UpdateCvsShipping CvsShipping.UpdateCvsShipping
	UpdateCvsShipping.Type = `F25`
	UpdateCvsShipping.ShipType = Enum.CVS_OK_MART

	for _, f := range fileNames {
		data, err := client.RetrFile(t, f)
		if err != nil {
			return
		}
		content, err := client.GetF25Xml(data)
		log.Info(t+`:`+f+`content:`, content)
		if err != nil {
			return
		}
		for _, d := range content.Body {
			UpdateCvsShipping.OrderNo = d.VendorNo
			UpdateCvsShipping.DateTime = okDateTimeFormat(d.UpDateTime)
			UpdateCvsShipping.DetailStatus = ""
			UpdateCvsShipping.Log = tools.XmlToString(d)
			UpdateCvsShipping.FileName = f

			err = UpdateCvsShipping.UpdateCvsShippingShipment(engine)

			if err != nil {
				log.Error(t+":"+f+" Error Data [%s]:", d)
				log.Error(t+" Error [%s]:", err.Error())
			}
		}
		client.MoveFileToBackup(t, f)
	}
}

// 貨態-寄件離店
func MartOKFetchF84(ecNo string, client om.Client, engine *database.MysqlSession) {

	t := "/" + ecNo + "/F84"

	fileNames, err := client.FetchFolder(t)
	if err != nil {
		return
	}

	var UpdateCvsShipping CvsShipping.UpdateCvsShipping
	UpdateCvsShipping.Type = `F84`
	UpdateCvsShipping.ShipType = Enum.CVS_OK_MART

	for _, f := range fileNames {
		data, err := client.RetrFile(t, f)
		if err != nil {
			return
		}
		content, err := client.GetF84Xml(data)
		log.Info(t+`:`+f+`content:`, content)
		if err != nil {
			return
		}
		for _, d := range content.Body {
			UpdateCvsShipping.ShipNo = d.OrderNo
			UpdateCvsShipping.DateTime = okDateTimeFormat(d.UpDateTime)
			UpdateCvsShipping.DetailStatus = ""
			UpdateCvsShipping.Log = tools.XmlToString(d)
			UpdateCvsShipping.FileName = f

			if ecNo == client.EcNo1 {
				// 正向離店
				err = UpdateCvsShipping.UpdateCvsShippingTransit(engine)
			} else {
				// 未取貨 逆向離店
				UpdateCvsShipping.FlowType = `R`
				err = UpdateCvsShipping.OnlyWriteShippingLog(engine, true)
			}

			if err != nil {
				log.Error(t+":"+f+" Error Data [%s]:", d)
				log.Error(t+" Error [%s]:", err.Error())
			}
		}
		client.MoveFileToBackup(t, f)
	}

}

// 貨態-物流中心驗收1
func MartOKFetchF71(ecNo string, client om.Client, engine *database.MysqlSession) {
	t := "/" + ecNo + "/F71"

	fileNames, err := client.FetchFolder(t)
	if err != nil {
		return
	}

	var UpdateCvsShipping CvsShipping.UpdateCvsShipping
	UpdateCvsShipping.Type = `F71`
	UpdateCvsShipping.ShipType = Enum.CVS_OK_MART

	for _, f := range fileNames {
		data, err := client.RetrFile(t, f)
		if err != nil {
			return
		}
		content, err := client.GetF71Xml(data)
		log.Info(t+`:`+f+`content:`, content)
		if err != nil {
			return
		}
		for _, d := range content.Body {
			UpdateCvsShipping.ShipNo = d.OrderNo
			UpdateCvsShipping.DateTime = okDateTimeFormat(d.UpDateTime)
			UpdateCvsShipping.DetailStatus = ""
			UpdateCvsShipping.Log = tools.XmlToString(d)
			UpdateCvsShipping.FileName = f

			err = UpdateCvsShipping.OnlyWriteShippingLog(engine, false)
			if err != nil {
				log.Error(t+":"+f+" Error Data [%s]:", d)
				log.Error(t+" Error [%s]:", err.Error())
			}
		}
		client.MoveFileToBackup(t, f)
	}

}

// 貨態-物流中心驗收2
func MartOKFetchF63(ecNo string, client om.Client, engine *database.MysqlSession) {
	t := "/" + ecNo + "/F63"

	fileNames, err := client.FetchFolder(t)
	if err != nil {
		return
	}

	var UpdateCvsShipping CvsShipping.UpdateCvsShipping
	UpdateCvsShipping.Type = `F63`
	UpdateCvsShipping.ShipType = Enum.CVS_OK_MART

	for _, f := range fileNames {
		data, err := client.RetrFile(t, f)
		if err != nil {
			return
		}
		content, err := client.GetF63Xml(data)
		log.Info(t+`:`+f+`content:`, content)
		if err != nil {
			return
		}
		for _, d := range content.Body {
			UpdateCvsShipping.ShipNo = d.OrderNo
			UpdateCvsShipping.DateTime = okDateTimeFormat(d.UpDateTime)
			UpdateCvsShipping.DetailStatus = ""
			UpdateCvsShipping.Log = tools.XmlToString(d)
			UpdateCvsShipping.FileName = f

			err = UpdateCvsShipping.OnlyWriteShippingLog(engine, false)
			if err != nil {
				log.Error(t+":"+f+" Error Data [%s]:", d)
				log.Error(t+" Error [%s]:", err.Error())
			}
		}
		client.MoveFileToBackup(t, f)
	}
}

// 貨態-物流中心驗收2(重出)
func MartOKFetchF03(ecNo string, client om.Client, engine *database.MysqlSession) {

	t := "/" + ecNo + "/F03"

	fileNames, err := client.FetchFolder(t)
	if err != nil {
		return
	}

	var UpdateCvsShipping CvsShipping.UpdateCvsShipping
	UpdateCvsShipping.Type = `F03`
	UpdateCvsShipping.ShipType = Enum.CVS_OK_MART

	for _, f := range fileNames {
		data, err := client.RetrFile(t, f)
		if err != nil {
			return
		}
		content, err := client.GetF03Xml(data)
		log.Info(t+`:`+f+`content:`, content)
		if err != nil {
			return
		}
		for _, d := range content.Body {
			UpdateCvsShipping.ShipNo = d.OrderNo
			UpdateCvsShipping.DateTime = okDateTimeFormat(d.UpDateTime)
			UpdateCvsShipping.FlowType = client.GetFlowType(ecNo)
			UpdateCvsShipping.DetailStatus = ""
			UpdateCvsShipping.Log = tools.XmlToString(d)
			UpdateCvsShipping.FileName = f

			err = UpdateCvsShipping.OnlyWriteShippingLog(engine, false)
			if err != nil {
				log.Error(t+":"+f+" Error Data [%s]:", d)
				log.Error(t+" Error [%s]:", err.Error())
			}
		}
		client.MoveFileToBackup(t, f)
	}

}

// 貨態-寄件驗退
func MartOKFetchF67(ecNo string, client om.Client, engine *database.MysqlSession) {

	t := "/" + ecNo + "/F67"

	fileNames, err := client.FetchFolder(t)
	if err != nil {
		return
	}

	var UpdateCvsShipping CvsShipping.UpdateCvsShipping
	UpdateCvsShipping.Type = `F67`
	UpdateCvsShipping.ShipType = Enum.CVS_OK_MART

	for _, f := range fileNames {
		data, err := client.RetrFile(t, f)
		if err != nil {
			return
		}
		content, err := client.GetF67Xml(data)
		log.Info(t+`:`+f+`content:`, content)
		if err != nil {
			return
		}
		for _, d := range content.Body {
			UpdateCvsShipping.ShipNo = d.OrderNo
			UpdateCvsShipping.DateTime = okDateTimeFormat(d.UpDateTime)
			UpdateCvsShipping.DetailStatus = d.ReturnCode
			UpdateCvsShipping.Log = tools.XmlToString(d)
			UpdateCvsShipping.FileName = f

			if UpdateCvsShipping.DetailStatus == `T01` {
				err = UpdateCvsShipping.UpdateCvsShippingSwitch(engine)
			} else {
				if UpdateCvsShipping.DetailStatus == `T00` {
					err = UpdateCvsShipping.UpdateCvsShippingBuyerNotPickUp(engine)
				} else {
					err = UpdateCvsShipping.UpdateCvsShippingFail(engine)
				}
			}

			if err != nil {
				log.Error(t+":"+f+" Error Data [%s]:", d)
				log.Error(t+" Error [%s]:", err.Error())
			}
		}
		client.MoveFileToBackup(t, f)
	}
}

// 即時進店檔
func MartOKFetchF44(ecNo string, client om.Client, engine *database.MysqlSession) {

	t := "/" + ecNo + "/F44"

	fileNames, err := client.FetchFolder(t)
	if err != nil {
		return
	}

	var UpdateCvsShipping CvsShipping.UpdateCvsShipping
	UpdateCvsShipping.Type = `F44`
	UpdateCvsShipping.ShipType = Enum.CVS_OK_MART

	for _, f := range fileNames {
		data, err := client.RetrFile(t, f)
		if err != nil {
			return
		}
		content, err := client.GetF44Xml(data)
		log.Info(t+`:`+f+`content:`, content)
		if err != nil {
			return
		}
		for _, d := range content.Body {
			UpdateCvsShipping.ShipNo = d.OrderNo
			UpdateCvsShipping.DateTime = okDateTimeFormat(d.RrInDateTime)
			UpdateCvsShipping.DetailStatus = ``
			UpdateCvsShipping.FlowType = client.GetFlowType(ecNo)
			UpdateCvsShipping.Log = tools.XmlToString(d)
			UpdateCvsShipping.FileName = f

			err = UpdateCvsShipping.UpdateCvsShippingShop(engine)

			if err != nil {
				log.Error(t+":"+f+" Error Data [%s]:", d)
				log.Error(t+" Error [%s]:", err.Error())
			}
		}
		client.MoveFileToBackup(t, f)
	}
}

// 進店檔 若 F44 檔案接收有遺漏可 於此檔案補入。
func MartOKFetchF64(ecNo string, client om.Client, engine *database.MysqlSession) {

	t := "/" + ecNo + "/F64"

	fileNames, err := client.FetchFolder(t)
	if err != nil {
		return
	}

	var UpdateCvsShipping CvsShipping.UpdateCvsShipping
	UpdateCvsShipping.Type = `F64`
	UpdateCvsShipping.ShipType = Enum.CVS_OK_MART

	for _, f := range fileNames {
		data, err := client.RetrFile(t, f)
		if err != nil {
			return
		}
		content, err := client.GetF64Xml(data)
		log.Info(t+`:`+f+`content:`, content)
		if err != nil {
			return
		}
		for _, d := range content.Body {
			UpdateCvsShipping.ShipNo = d.OrderNo
			UpdateCvsShipping.DateTime = d.GetDateTime()
			UpdateCvsShipping.DetailStatus = ``
			UpdateCvsShipping.FlowType = client.GetFlowType(ecNo)
			UpdateCvsShipping.Log = tools.XmlToString(d)
			UpdateCvsShipping.FileName = f

			err = UpdateCvsShipping.UpdateCvsShippingShop(engine)

			if err != nil {
				log.Error(t+":"+f+" Error Data [%s]:", d)
				log.Error(t+" Error [%s]:", err.Error())
			}
		}
		client.MoveFileToBackup(t, f)
	}

}

//  F04 重出進店檔 若 F44 檔案接收有遺 漏可於此檔案補入。
func MartOKFetchF04(ecNo string, client om.Client, engine *database.MysqlSession) {

	t := "/" + ecNo + "/F04"

	fileNames, err := client.FetchFolder(t)
	if err != nil {
		return
	}

	var UpdateCvsShipping CvsShipping.UpdateCvsShipping
	UpdateCvsShipping.Type = `F04`
	UpdateCvsShipping.ShipType = Enum.CVS_OK_MART

	for _, f := range fileNames {
		data, err := client.RetrFile(t, f)
		if err != nil {
			return
		}
		content, err := client.GetF04Xml(data)
		log.Info(t+`:`+f+`content:`, content)
		if err != nil {
			return
		}
		for _, d := range content.Body {
			UpdateCvsShipping.ShipNo = d.OrderNo
			UpdateCvsShipping.DateTime = okDateTimeFormat(d.UpDateTime)
			UpdateCvsShipping.DetailStatus = ``
			UpdateCvsShipping.FlowType = client.GetFlowType(ecNo)
			UpdateCvsShipping.Log = tools.XmlToString(d)
			UpdateCvsShipping.FileName = f

			err = UpdateCvsShipping.UpdateCvsShippingShop(engine)

			if err != nil {
				log.Error(t+":"+f+" Error Data [%s]:", d)
				log.Error(t+" Error [%s]:", err.Error())
			}
		}
		client.MoveFileToBackup(t, f)
	}
}

// 即時取貨代收檔
func MartOKFetchF17(ecNo string, client om.Client, engine *database.MysqlSession) {

	t := "/" + ecNo + "/F17"

	fileNames, err := client.FetchFolder(t)
	if err != nil {
		return
	}

	var UpdateCvsShipping CvsShipping.UpdateCvsShipping
	UpdateCvsShipping.Type = `F17`
	UpdateCvsShipping.ShipType = Enum.CVS_OK_MART

	for _, f := range fileNames {
		data, err := client.RetrFile(t, f)
		if err != nil {
			return
		}
		content, err := client.GetF17Xml(data)
		log.Info(t+`:`+f+`content:`, content)
		if err != nil {
			return
		}
		for _, d := range content.Body {
			UpdateCvsShipping.ShipNo = d.SendNo
			UpdateCvsShipping.DateTime = okDateTimeFormat(d.RrPickDateTime)
			UpdateCvsShipping.DetailStatus = ``
			UpdateCvsShipping.FlowType = client.GetFlowType(ecNo)
			UpdateCvsShipping.Log = tools.XmlToString(d)
			UpdateCvsShipping.FileName = f

			err = UpdateCvsShipping.UpdateCvsShippingSuccess(engine)

			if err != nil {
				log.Error(t+":"+f+" Error Data [%s]:", d)
				log.Error(t+" Error [%s]:", err.Error())
			}
		}
		client.MoveFileToBackup(t, f)
	}

}

// 取貨代收檔 若 F17 檔案接收有遺漏可 於此檔案補入。
func MartOKFetchF65(ecNo string, client om.Client, engine *database.MysqlSession) {

	t := "/" + ecNo + "/F65"

	fileNames, err := client.FetchFolder(t)
	if err != nil {
		return
	}

	var UpdateCvsShipping CvsShipping.UpdateCvsShipping
	UpdateCvsShipping.Type = `F65`
	UpdateCvsShipping.ShipType = Enum.CVS_OK_MART

	for _, f := range fileNames {
		data, err := client.RetrFile(t, f)
		if err != nil {
			return
		}
		content, err := client.GetF65Xml(data)
		log.Info(t+`:`+f+`content:`, content)
		if err != nil {
			return
		}
		for _, d := range content.Body {
			UpdateCvsShipping.ShipNo = om.GenerateOkBarCodeToShipNo(d.BarCode1, d.BarCode2)
			UpdateCvsShipping.DateTime = okDateTimeFormat(d.RrPickDateTime)
			UpdateCvsShipping.DetailStatus = ``
			UpdateCvsShipping.FlowType = client.GetFlowType(ecNo)
			UpdateCvsShipping.Log = tools.XmlToString(d)
			UpdateCvsShipping.FileName = f

			err = UpdateCvsShipping.UpdateCvsShippingSuccess(engine)

			if err != nil {
				log.Error(t+":"+f+" Error Data [%s]:", d)
				log.Error(t+" Error [%s]:", err.Error())
			}
		}
		client.MoveFileToBackup(t, f)
	}

}

func MartOKFetchF05(ecNo string, client om.Client, engine *database.MysqlSession) {

	t := "/" + ecNo + "/F05"

	fileNames, err := client.FetchFolder(t)
	if err != nil {
		return
	}

	var UpdateCvsShipping CvsShipping.UpdateCvsShipping
	UpdateCvsShipping.Type = `F05`
	UpdateCvsShipping.ShipType = Enum.CVS_OK_MART

	for _, f := range fileNames {
		data, err := client.RetrFile(t, f)
		if err != nil {
			return
		}
		content, err := client.GetF05Xml(data)
		log.Info(t+`:`+f+`content:`, content)
		if err != nil {
			return
		}
		for _, d := range content.Body {
			//fixme 單號會有問題
			UpdateCvsShipping.ShipNo = om.GenerateOkBarCodeToShipNo(d.BarCode1, d.BarCode2)
			UpdateCvsShipping.DateTime = okDateTimeFormat(d.UpDateTime)
			UpdateCvsShipping.DetailStatus = ``
			UpdateCvsShipping.FlowType = client.GetFlowType(ecNo)
			UpdateCvsShipping.Log = tools.XmlToString(d)
			UpdateCvsShipping.FileName = f

			err = UpdateCvsShipping.UpdateCvsShippingSuccess(engine)

			if err != nil {
				log.Error(t+":"+f+" Error Data [%s]:", d)
				log.Error(t+" Error [%s]:", err.Error())
			}
		}
		client.MoveFileToBackup(t, f)
	}

}

// F07 重出物流驗退檔 若 F17 檔案接收有遺 漏可於此檔案補入。
func MartOKFetchF07(ecNo string, client om.Client, engine *database.MysqlSession) {

	t := "/" + ecNo + "/F07"

	fileNames, err := client.FetchFolder(t)
	if err != nil {
		return
	}

	var UpdateCvsShipping CvsShipping.UpdateCvsShipping
	UpdateCvsShipping.Type = `F07`
	UpdateCvsShipping.ShipType = Enum.CVS_OK_MART

	for _, f := range fileNames {
		data, err := client.RetrFile(t, f)
		if err != nil {
			return
		}
		content, err := client.GetF07Xml(data)
		log.Info(t+`:`+f+`content:`, content)
		if err != nil {
			return
		}
		for _, d := range content.Body {
			UpdateCvsShipping.ShipNo = d.OrderNo
			UpdateCvsShipping.DateTime = d.UpDateTime
			UpdateCvsShipping.DetailStatus = d.ReturnCode
			UpdateCvsShipping.FlowType = client.GetFlowType(ecNo)
			UpdateCvsShipping.Log = tools.XmlToString(d)
			UpdateCvsShipping.FileName = f

			if UpdateCvsShipping.DetailStatus == `T01` {
				err = UpdateCvsShipping.UpdateCvsShippingSwitch(engine)
			} else {
				if UpdateCvsShipping.DetailStatus == `T00` {
					err = UpdateCvsShipping.UpdateCvsShippingBuyerNotPickUp(engine)
				} else {
					err = UpdateCvsShipping.UpdateCvsShippingFail(engine)
				}

			}

			if err != nil {
				log.Error(t+"Error Data [%s]:", d)
				log.Error(t+"Error [%s]:", err.Error())
			}
		}
		client.MoveFileToBackup(t, f)
	}
}

// OK帳務
func MartOKFetchAcct(client om.Client, engine *database.MysqlSession) {

	t := "/462/ACCT"

	fileNames, err := client.FetchFolder(t)
	if err != nil {
		return
	}

	var accountingData entity.CvsAccountingData
	var type_ = `SEND`
	for _, f := range fileNames {
		data, err := client.RetrFile(t, f)
		if err != nil {
			return
		}

		match, _ := regexp.Match(".*SEND.*\\.(?i)CSV$", []byte(f))
		if match {
			type_ = `SEND`
		} else {
			type_ = `PICK`
		}
		accountingData.SetType(type_, Enum.CVS_OK_MART)
		scanner := bufio.NewScanner(bytes.NewReader(data))
		for scanner.Scan() {
			text := scanner.Text()
			a := strings.Split(text, `,`)
			if a[0] == `寄件店號` || a[0] == `取件店號` {
				continue
			}
			accountingData.DataId = a[2]
			accountingData.Amount, _ = strconv.ParseFloat(a[3], 64)
			accountingData.ServiceType = a[4] == `1`
			accountingData.FileDate, _ = time.Parse(`20060102`, a[5])
			accountingData.FileName = f
			accountingData.Status = `1`
			accountingData.Log = text

			_ = Cvs.InsertAccountingData(engine, accountingData)
		}
		client.MoveFileToBackup(t, f)
	}
}
