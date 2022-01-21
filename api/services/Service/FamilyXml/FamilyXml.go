package FamilyXml

import (
	"api/services/Enum"
	"api/services/Service/CvsShipping"
	fm "api/services/Service/FamilyMartLogistics"
	"api/services/dao/Cvs"
	"api/services/dao/Mart"
	"api/services/database"
	"api/services/entity"
	"api/services/util/log"
	"api/services/util/tools"
	tw_address "api/services/util/tw-address"
	"encoding/xml"
	"strconv"
	"time"
)

type F100Xml struct {
	Xml    xml.Name   `xml:"doc"`
	Header F100Header `xml:"HEADER"`
	I00    []I00      `xml:"BODY>I00"`
}

type F100Header struct {
	RDFMT string // 區別碼
	SNCD  string
	PRDT  string
}

type I00 struct {
	RDFMT          string
	StoreId        string
	StoreName      string
	MdcStareDate   string
	MdcEndDate     string
	ROUTE          string
	STEP           string
	StoreAddress   string
	TelNo          string
	OldStore       string
	StoreCloseDate string
	Area           string
	EquipmentID    string
}

func getMartFamilyStore(engine *database.MysqlSession) error {
	t := "I00"
	var client fm.Client
	client.GetClient(false)

	if client.ConnectFtpError != nil {
		log.Error("MartFamilyFetchShipping [%v]", client.ConnectFtpError.Error())
		return client.ConnectFtpError
	}
	defer client.ConnectFtp.Quit()

	fileNames, err := client.FetchFolder(t)
	if err != nil {
		return err
	}

	// 下載解壓縮讀取資料
	data, err := client.RetrFileAndUnzip(t, fileNames[0])
	if err != nil {
		return err
	}
	var xmlData F100Xml

	err = xml.Unmarshal([]byte(data), &xmlData)
	if err != nil {
		client.MoveFileToDest(t, fileNames[0], `ERR`)
		return err
	}

	var inS []interface{}
	for _, data := range xmlData.I00 {
		country, district, address, flag := tw_address.AddressSplit(data.StoreAddress)
		if flag {
			store := entity.MartFamilyStoreData{
				StoreId:        data.StoreId,
				StoreName:      data.StoreName,
				StoreAddress:   address,
				StoreCloseDate: data.StoreCloseDate,
				MdcStareDate:   data.MdcStareDate,
				MdcEndDate:     data.MdcEndDate,
				Route:          data.ROUTE,
				Step:           data.STEP,
				TelNo:          data.TelNo,
				OldStore:       data.OldStore,
				Area:           data.Area,
				EquipmentId:    data.EquipmentID,
				City:           country,
				District:       district,
				UpdateTime:     time.Now(),
			}
			exists := new(entity.MartFamilyStoreData)

			query := map[string]interface{}{}
			query["store_id"] = data.StoreId
			exist, err := engine.Engine.Where(query).Exist(exists)

			if err != nil {
				log.Error("Exist Writing Error:%v", err)
				continue
			}

			if exist {
				_, _ = engine.Session.Where(query).Update(store)
			} else {
				inS = append(inS, store)
			}
		}

	}
	_, _ = engine.Session.Insert(inS...)

	client.MoveFileToDest(t, fileNames[0], `OK`)

	return nil
}

func MartFamilyFetchStoreList() {
	log.Debug("MartFamilyFetchStoreList")

	var engine = database.GetMysqlEngine()
	defer engine.Close()

	// 從 Ftp 更新店舖資訊進 DB
	_ = getMartFamilyStore(engine)

	var newStores []entity.MartFamilyStoreData
	err := engine.Engine.Find(&newStores)
	if err != nil {
		log.Error("MartFamilyFetchStoreList:", err)
		return
	}

	_ = Mart.WriteFamilyStoreData(newStores)
}

// 全家帳務
func MartFamilyFetchAccounting()  {
	log.Info("MartFamilyFetchAccounting")
	var client fm.Client
	client.GetClient(false)

	if client.ConnectFtpError != nil {
		log.Error("MartFamilyFetchShipping [%v]", client.ConnectFtpError.Error())
		return
	}
	defer client.ConnectFtp.Quit()

	engine := database.GetMysqlEngine()
	defer engine.Close()

	MartFamilyFetchR89(engine, client)
	MartFamilyFetchR98(engine, client)
	MartFamilyFetchR99(engine, client)
}

func MartFamilyFetchShipping() {
	log.Info("MartFamilyFetchShipping")

	var client fm.Client
	client.GetClient(false)

	if client.ConnectFtpError != nil {
		log.Error("MartFamilyFetchShipping [%v]", client.ConnectFtpError.Error())
		return
	}
	defer client.ConnectFtp.Quit()

	engine := database.GetMysqlEngine()
	defer engine.Close()

	// 商品寄件檔/多批次
	MartFamilyFetchR22(engine, client)
	// 寄件離店檔/多批次
	MartFamilyFetchR23(engine, client)
	// 物流進貨驗收檔/多批次
	MartFamilyFetchR25(engine, client)

	// 貨態-物流出貨 (批次）
	MartFamilyFetchR27(engine, client)
	// 物流出貨通知檔 22:30  R04 與 R27 相同 R04為日檔 之後全家會捨棄
	MartFamilyFetchR04(engine, client)

	// 物流驗收異常檔 21:30
	MartFamilyFetchRS9(engine, client)

	// 店舖進貨驗收檔/多批次
	MartFamilyFetchR28(engine, client)
	//店舖進貨驗收檔 10:30
	MartFamilyFetchRS4(engine, client)

	// 商品取貨檔/多批次
	MartFamilyFetchR29(engine, client)
	// 商品取貨檔 12:00
	MartFamilyFetchR96(engine, client)

	// 店舖未取退回物流驗退檔 18:30
	MartFamilyFetchR08(engine, client)

}

// 商品寄件檔/多批次
func MartFamilyFetchR22(engine *database.MysqlSession, client fm.Client) {

	t := "R22"
	fileNames, err := client.FetchFolder(t)
	if err != nil {
		return
	}

	var UpdateCvsShipping CvsShipping.UpdateCvsShipping
	UpdateCvsShipping.Type = t
	UpdateCvsShipping.ShipType = Enum.CVS_FAMILY

	for _, f := range fileNames {
		data, err := client.RetrFileAndUnzip(t, f)
		if err != nil {
			return
		}
		body, err := client.GetR22Xml(data)
		log.Info(t+`:`+f+`body:`, body)
		if err != nil {
			client.MoveFileToDest(t, f, `ERR`)
			return
		}
		for _, d := range body.Body.Data {
			UpdateCvsShipping.ShipNo = d.OrderNo
			UpdateCvsShipping.DateTime = d.GetDateTime()
			UpdateCvsShipping.DetailStatus = ""
			UpdateCvsShipping.Log = tools.XmlToString(d)
			UpdateCvsShipping.FileName = f
			err := UpdateCvsShipping.UpdateCvsShippingShipment(engine)
			if err != nil {
				log.Error(t+" Writing Error Data [%s]:", d)
				log.Error(t+" Writing Error [%s]:", err.Error())
			}
		}

		client.MoveFileToDest(t, f, `OK`)
	}

}

// 寄件離店檔/多批次
func MartFamilyFetchR23(engine *database.MysqlSession, client fm.Client) {
	t := "R23"
	fileNames, err := client.FetchFolder(t)
	if err != nil {
		return
	}

	var UpdateCvsShipping CvsShipping.UpdateCvsShipping
	UpdateCvsShipping.Type = t
	UpdateCvsShipping.ShipType = Enum.CVS_FAMILY

	for _, f := range fileNames {
		data, err := client.RetrFileAndUnzip(t, f)
		if err != nil {
			return
		}
		body, err := client.GetR23Xml(data)
		log.Info(t+`:`+f+`body:`, body)
		if err != nil {
			client.MoveFileToDest(t, f, `ERR`)
			return
		}
		for _, d := range body.Body.Data {
			UpdateCvsShipping.ShipNo = d.OrderNo
			UpdateCvsShipping.DateTime = d.GetDateTime()
			UpdateCvsShipping.DetailStatus = ""
			UpdateCvsShipping.Log = tools.XmlToString(d)
			UpdateCvsShipping.FileName = f
			err := UpdateCvsShipping.UpdateCvsShippingTransit(engine)

			if err != nil {
				log.Error(t+" Writing Error Data [%s]:", d)
				log.Error(t+" Writing Error [%s]:", err.Error())
			}
		}

		client.MoveFileToDest(t, f, `OK`)
	}
}

//物流進貨驗收檔/多批次
func MartFamilyFetchR25(engine *database.MysqlSession, client fm.Client) {

	t := "R25"
	fileNames, err := client.FetchFolder(t)
	if err != nil {
		return
	}

	var UpdateCvsShipping CvsShipping.UpdateCvsShipping
	UpdateCvsShipping.Type = t
	UpdateCvsShipping.ShipType = Enum.CVS_FAMILY

	for _, f := range fileNames {
		data, err := client.RetrFileAndUnzip(t, f)
		if err != nil {
			return
		}
		body, err := client.GetR25Xml(data)
		log.Info(t+`:`+f+`body:`, body)
		if err != nil {
			client.MoveFileToDest(t, f, `ERR`)
			return
		}
		for _, d := range body.Body.Data {
			UpdateCvsShipping.OrderNo = d.ShipmentNo
			UpdateCvsShipping.DateTime = d.GetDateTime()
			UpdateCvsShipping.DetailStatus = d.DCReceiveStatus
			UpdateCvsShipping.FlowType = d.FlowType
			UpdateCvsShipping.Log = tools.XmlToString(d)
			UpdateCvsShipping.FileName = f
			err := UpdateCvsShipping.OnlyWriteShippingLog(engine, true)

			if err != nil {
				log.Error(t+" Writing Error Data [%s]:", d)
				log.Error(t+" Writing Error [%s]:", err.Error())
			}
		}

		client.MoveFileToDest(t, f, `OK`)
	}
}

// 貨態-物流出貨 (批次）
func MartFamilyFetchR27(engine *database.MysqlSession, client fm.Client) {

	t := "R27"
	fileNames, err := client.FetchFolder(t)
	if err != nil {
		return
	}

	var UpdateCvsShipping CvsShipping.UpdateCvsShipping
	UpdateCvsShipping.Type = t
	UpdateCvsShipping.ShipType = Enum.CVS_FAMILY

	for _, f := range fileNames {
		data, err := client.RetrFileAndUnzip(t, f)
		if err != nil {
			return
		}
		body, err := client.GetR27Xml(data)
		log.Info(t+`:`+f+`body:`, body)
		if err != nil {
			client.MoveFileToDest(t, f, `ERR`)
			return
		}
		for _, d := range body.Body.Data {
			UpdateCvsShipping.OrderNo = d.ShipmentNo
			UpdateCvsShipping.DateTime = d.GetDateTime()
			UpdateCvsShipping.DetailStatus = ""
			UpdateCvsShipping.Log = tools.XmlToString(d)
			UpdateCvsShipping.FileName = f
			err := UpdateCvsShipping.OnlyWriteShippingLog(engine, false)

			if err != nil {
				log.Error(t+" Writing Error Data [%s]:", d)
				log.Error(t+" Writing Error [%s]:", err.Error())
			}
		}

		client.MoveFileToDest(t, f, `OK`)
	}
}

// 物流出貨通知檔 22:30
func MartFamilyFetchR04(engine *database.MysqlSession, client fm.Client) {

	t := "R04"
	fileNames, err := client.FetchFolder(t)
	if err != nil {
		return
	}

	var UpdateCvsShipping CvsShipping.UpdateCvsShipping
	UpdateCvsShipping.Type = t
	UpdateCvsShipping.ShipType = Enum.CVS_FAMILY

	for _, f := range fileNames {
		data, err := client.RetrFileAndUnzip(t, f)
		if err != nil {
			return
		}
		body, err := client.GetR04Xml(data)
		log.Info(t+`:`+f+`body:`, body)
		if err != nil {
			client.MoveFileToDest(t, f, `ERR`)
			return
		}
		for _, d := range body.Body.Data {
			UpdateCvsShipping.OrderNo = d.ShipmentNo
			UpdateCvsShipping.DateTime = d.DCReceiveDate
			UpdateCvsShipping.DetailStatus = d.DCReceiveStatus
			UpdateCvsShipping.FlowType = d.FlowType
			UpdateCvsShipping.Log = tools.XmlToString(d)
			UpdateCvsShipping.FileName = f
			err := UpdateCvsShipping.OnlyWriteShippingLog(engine, false)

			if err != nil {
				log.Error(t+" Writing Error Data [%s]:", d)
				log.Error(t+" Writing Error [%s]:", err.Error())
			}
		}

		client.MoveFileToDest(t, f, `OK`)
	}
}

// 物流驗收異常檔 21:30
func MartFamilyFetchRS9(engine *database.MysqlSession, client fm.Client) {

	t := "RS9"
	fileNames, err := client.FetchFolder(t)
	if err != nil {
		return
	}

	var UpdateCvsShipping CvsShipping.UpdateCvsShipping
	UpdateCvsShipping.Type = t
	UpdateCvsShipping.ShipType = Enum.CVS_FAMILY

	for _, f := range fileNames {
		data, err := client.RetrFileAndUnzip(t, f)
		if err != nil {
			return
		}
		body, err := client.GetRS9Xml(data)
		log.Info(t+`:`+f+`body:`, body)
		if err != nil {
			client.MoveFileToDest(t, f, `ERR`)
			return
		}
		for _, d := range body.Body.Data {
			UpdateCvsShipping.OrderNo = d.ShipmentNo
			UpdateCvsShipping.DateTime = d.DCReceiveDate
			UpdateCvsShipping.DetailStatus = d.DCReceiveStatus
			UpdateCvsShipping.FlowType = d.FlowType
			UpdateCvsShipping.Log = tools.XmlToString(d)
			UpdateCvsShipping.FileName = f

			if UpdateCvsShipping.DetailStatus == `T01` {
				err = UpdateCvsShipping.UpdateCvsShippingSwitch(engine)
			} else {
				err = UpdateCvsShipping.UpdateCvsShippingFail(engine)
			}

			if err != nil {
				log.Error(t+" Writing Error Data [%s]:", d)
				log.Error(t+" Writing Error [%s]:", err.Error())
			}
		}

		client.MoveFileToDest(t, f, `OK`)
	}
}

// 店舖進貨驗收檔/多批次
func MartFamilyFetchR28(engine *database.MysqlSession, client fm.Client) {

	t := "R28"
	fileNames, err := client.FetchFolder(t)
	if err != nil {
		return
	}

	var UpdateCvsShipping CvsShipping.UpdateCvsShipping
	UpdateCvsShipping.Type = t
	UpdateCvsShipping.ShipType = Enum.CVS_FAMILY

	for _, f := range fileNames {
		data, err := client.RetrFileAndUnzip(t, f)
		if err != nil {
			return
		}
		body, err := client.GetR28Xml(data)
		log.Info(t+`:`+f+`body:`, body)
		if err != nil {
			client.MoveFileToDest(t, f, `ERR`)
			return
		}
		for _, d := range body.Body.Data {
			UpdateCvsShipping.ShipNo = d.OrderNo
			UpdateCvsShipping.DateTime = d.GetDateTime()
			UpdateCvsShipping.DetailStatus = ``
			UpdateCvsShipping.FlowType = d.FlowType
			UpdateCvsShipping.Log = tools.XmlToString(d)
			UpdateCvsShipping.FileName = f

			err := UpdateCvsShipping.UpdateCvsShippingShop(engine)

			if err != nil {
				log.Error(t+" Writing Error Data [%s]:", d)
				log.Error(t+" Writing Error [%s]:", err.Error())
			}
		}

		client.MoveFileToDest(t, f, `OK`)
	}

}

// 店舖進貨驗收檔 10:30
func MartFamilyFetchRS4(engine *database.MysqlSession, client fm.Client) {

	t := "RS4"
	fileNames, err := client.FetchFolder(t)
	if err != nil {
		return
	}

	var UpdateCvsShipping CvsShipping.UpdateCvsShipping
	UpdateCvsShipping.Type = t
	UpdateCvsShipping.ShipType = Enum.CVS_FAMILY

	for _, f := range fileNames {
		data, err := client.RetrFileAndUnzip(t, f)
		if err != nil {
			return
		}
		body, err := client.GetRS4Xml(data)
		log.Info(t+`:`+f+`body:`, body)
		if err != nil {
			client.MoveFileToDest(t, f, `ERR`)
			return
		}
		for _, d := range body.Body.Data {
			UpdateCvsShipping.OrderNo = d.ShipmentNo
			UpdateCvsShipping.DateTime = d.DCReceiveDate
			UpdateCvsShipping.DetailStatus = ``
			UpdateCvsShipping.FlowType = d.FlowType
			UpdateCvsShipping.Log = tools.XmlToString(d)
			UpdateCvsShipping.FileName = f

			err := UpdateCvsShipping.UpdateCvsShippingShop(engine)

			if err != nil {
				log.Error(t+" Writing Error Data [%s]:", d)
				log.Error(t+" Writing Error [%s]:", err.Error())
			}
		}

		client.MoveFileToDest(t, f, `OK`)
	}
}

// 商品取貨檔/多批次
func MartFamilyFetchR29(engine *database.MysqlSession, client fm.Client) {

	t := "R29"
	fileNames, err := client.FetchFolder(t)
	if err != nil {
		return
	}

	var UpdateCvsShipping CvsShipping.UpdateCvsShipping
	UpdateCvsShipping.Type = t
	UpdateCvsShipping.ShipType = Enum.CVS_FAMILY

	for _, f := range fileNames {
		data, err := client.RetrFileAndUnzip(t, f)
		if err != nil {
			return
		}
		body, err := client.GetR29Xml(data)
		log.Info(t+`:`+f+`body:`, body)
		if err != nil {
			client.MoveFileToDest(t, f, `ERR`)
			return
		}
		for _, d := range body.Body.Data {
			UpdateCvsShipping.ShipNo = d.OrderNo
			UpdateCvsShipping.DateTime = d.GetDateTime()
			UpdateCvsShipping.DetailStatus = ``
			UpdateCvsShipping.FlowType = d.FlowType
			UpdateCvsShipping.Log = tools.XmlToString(d)
			UpdateCvsShipping.FileName = f

			err := UpdateCvsShipping.UpdateCvsShippingSuccess(engine)

			if err != nil {
				log.Error(t+" Writing Error Data [%s]:", d)
				log.Error(t+" Writing Error [%s]:", err.Error())
			}
		}

		client.MoveFileToDest(t, f, `OK`)
	}
}

// 商品取貨檔 12:00
func MartFamilyFetchR96(engine *database.MysqlSession, client fm.Client) {

	t := "R96"
	fileNames, err := client.FetchFolder(t)
	if err != nil {
		return
	}

	var UpdateCvsShipping CvsShipping.UpdateCvsShipping
	UpdateCvsShipping.Type = t
	UpdateCvsShipping.ShipType = Enum.CVS_FAMILY

	for _, f := range fileNames {
		data, err := client.RetrFileAndUnzip(t, f)
		if err != nil {
			return
		}
		body, err := client.GetR96Xml(data)
		log.Info(t+`:`+f+`body:`, body)
		if err != nil {
			client.MoveFileToDest(t, f, `ERR`)
			return
		}
		for _, d := range body.Body.Data {
			UpdateCvsShipping.OrderNo = d.ShipmentNo
			UpdateCvsShipping.DateTime = d.SPDate
			UpdateCvsShipping.DetailStatus = ``
			UpdateCvsShipping.FlowType = d.FlowType
			UpdateCvsShipping.Log = tools.XmlToString(d)
			UpdateCvsShipping.FileName = f

			err := UpdateCvsShipping.UpdateCvsShippingSuccess(engine)

			if err != nil {
				log.Error(t+" Writing Error Data [%s]:", d)
				log.Error(t+" Writing Error [%s]:", err.Error())
			}
		}

		client.MoveFileToDest(t, f, `OK`)
	}
}

// 店舖未取退回物流驗退檔 18:30
func MartFamilyFetchR08(engine *database.MysqlSession, client fm.Client) {

	t := "R08"
	fileNames, err := client.FetchFolder(t)
	if err != nil {
		return
	}

	var UpdateCvsShipping CvsShipping.UpdateCvsShipping
	UpdateCvsShipping.Type = t
	UpdateCvsShipping.ShipType = Enum.CVS_FAMILY

	for _, f := range fileNames {
		data, err := client.RetrFileAndUnzip(t, f)
		if err != nil {
			return
		}
		body, err := client.GetR08Xml(data)
		log.Info(t+`:`+f+`body:`, body)
		if err != nil {
			client.MoveFileToDest(t, f, `ERR`)
			return
		}
		for _, d := range body.Body.Data {
			UpdateCvsShipping.OrderNo = d.ShipmentNo
			UpdateCvsShipping.DateTime = d.DCReturnDate
			UpdateCvsShipping.DetailStatus = ``
			UpdateCvsShipping.FlowType = d.FlowType
			UpdateCvsShipping.Log = tools.XmlToString(d)
			UpdateCvsShipping.FileName = f
			err := UpdateCvsShipping.UpdateCvsShippingBuyerNotPickUp(engine)

			if err != nil {
				log.Error(t+" Writing Error Data [%s]:", d)
				log.Error(t+" Writing Error [%s]:", err.Error())
			}
		}

		client.MoveFileToDest(t, f, `OK`)
	}
}

func MartFamilyFetchR89(engine *database.MysqlSession, client fm.Client) {

	t := "R89"
	fileNames, err := client.FetchFolder(t)
	if err != nil {
		return
	}

	var accountingData entity.CvsAccountingData
	accountingData.SetType(t, Enum.CVS_FAMILY)
	for _, f := range fileNames {
		data, err := client.RetrFileAndUnzip(t, f)
		if err != nil {
			return
		}
		body, err := client.GetR89Xml(data)
		log.Info(t+`:`+f+`body:`, body)
		if err != nil {
			client.MoveFileToDest(t, f, `ERR`)
			return
		}
		for _, d := range body.Body.Data {
			accountingData.DataId = d.ShipmentNo
			accountingData.Amount, _ = strconv.ParseFloat(d.SPAmount, 64)
			accountingData.ServiceType = d.ServiceType == `1`
			accountingData.Status = d.SPAstatus
			accountingData.FileDate, _ = time.Parse(`2006-01-02`, d.SPAdate)
			accountingData.FileName = f
			accountingData.Log = tools.XmlToString(d)

			_ = Cvs.InsertAccountingData(engine, accountingData)
		}

		client.MoveFileToDest(t, f, `OK`)
	}
}

func MartFamilyFetchR98(engine *database.MysqlSession, client fm.Client) {

	t := "R98"

	fileNames, err := client.FetchFolder(t)
	if err != nil {
		return
	}

	var accountingData entity.CvsAccountingData
	accountingData.SetType(t, Enum.CVS_FAMILY)

	for _, f := range fileNames {
		data, err := client.RetrFileAndUnzip(t, f)
		if err != nil {
			return
		}
		body, err := client.GetR98Xml(data)
		log.Info(t+`:`+f+`body:`, body)
		if err != nil {
			client.MoveFileToDest(t, f, `ERR`)
			return
		}
		for _, d := range body.Body.Data {
			accountingData.DataId = d.ShipmentNo
			accountingData.Amount, _ = strconv.ParseFloat(d.SPAmount, 64)
			accountingData.ServiceType = d.ServiceType == `1`
			accountingData.Status = d.SPAstatus
			accountingData.FileDate, _ = time.Parse(`20060102`, d.SPAdate)
			accountingData.FileName = f
			accountingData.Log = tools.XmlToString(d)

			_ = Cvs.InsertAccountingData(engine, accountingData)
		}

		client.MoveFileToDest(t, f, `OK`)
	}

}

func MartFamilyFetchR99(engine *database.MysqlSession, client fm.Client) {

	t := "R99"

	fileNames, err := client.FetchFolder(t)
	if err != nil {
		return
	}

	var accountingData entity.CvsAccountingData
	accountingData.SetType(t, Enum.CVS_FAMILY)

	for _, f := range fileNames {
		data, err := client.RetrFileAndUnzip(t, f)
		if err != nil {
			return
		}
		body, err := client.GetR99Xml(data)
		log.Info(t+`:`+f+`body:`, body)
		if err != nil {
			client.MoveFileToDest(t, f, `ERR`)
			return
		}
		for _, d := range body.Body.Data {
			accountingData.DataId = d.ShipmentNo
			accountingData.Amount, _ = strconv.ParseFloat(d.SPAmount, 64)
			accountingData.ServiceType = d.ServiceType == `1`
			accountingData.Status = d.SPAstatus
			accountingData.FileDate, _ = time.Parse(`20060102`, d.SPAdate)
			accountingData.FileName = f
			accountingData.Log = tools.XmlToString(d)

			_ = Cvs.InsertAccountingData(engine, accountingData)
		}

		client.MoveFileToDest(t, f, `OK`)
	}
}
