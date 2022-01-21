package HiLifeXml

import (
	"api/services/Enum"
	"api/services/Service/CvsShipping"
	"api/services/Service/MartHiLife"
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

type R00Xml struct {
	Xml    xml.Name  `xml:"doc"`
	Header R00Header `xml:"HEADER"`
	R00    []R00     `xml:"BODY>R00"`
}

type R00Header struct {
	DocDate         string `xml:"DOCDATE"`
	FromPartnerCode string `xml:"FROMPARTNERCODE"`
	ToPartenerCode  string `xml:"TOPARTNERCODE"`
}
type R00 struct {
	RDFMT          string `xml:"RDFMT"`
	StoreId        string `xml:"STOREID"`
	StoreName      string `xml:"STORE_NAME"`
	StoreAddress   string `xml:"STORE_ADDRESS"`
	TelNo          string `xml:"TEL_NO"`
	OldStore       string `xml:"OLD_STORE"`
	StoreCloseDate string `xml:"STORE_CLOSE_DATE"`
	MdcStareDate   string `xml:"MDC_START_DATE"`
	MdcEndDate     string `xml:"MDC_END_DATE"`
	ROUTE          string `xml:"ROUTER"`
	STEP           string `xml:"STEP"`
}

// 萊爾富帳務
func MartHiLifeAccounting()  {
	log.Info("MartHiLifeAccounting[Begin]")

	var client MartHiLife.Client
	client.GetClient()

	if client.ConnectFtpError != nil {
		log.Error("MartHiLifeFetchShipping [%v]", client.ConnectFtpError.Error())
		return
	}
	defer client.ConnectFtp.Close()

	engine := database.GetMysqlEngine()
	defer engine.Close()
	// 寄件運費檔 10:00
	MartHiLifeFetchR98(engine, client)

	// 遺失賠償檔 10:00
	MartHiLifeFetchR89(engine, client)

	// 取件核帳檔 10:00
	MartHiLifeFetchR99(engine, client)
}

func MartHiLifeFetchShipping() {
	log.Info("MartHiLifeFetchShipping[Begin]")

	var client MartHiLife.Client
	client.GetClient()

	if client.ConnectFtpError != nil {
		log.Error("MartHiLifeFetchShipping [%v]", client.ConnectFtpError.Error())
		return
	}
	defer client.ConnectFtp.Close()

	engine := database.GetMysqlEngine()
	defer engine.Close()

	// 即時交寄檔 整點
	MartHiLifeFetchR27(engine, client)

	// 店舖交寄日檔 09:00
	MartHiLifeFetchR22(engine, client)

	// DC 驗收檔 18:00 / 23:00
	MartHiLifeFetchR04(engine, client)

	// 店舖驗收時檔 整點
	MartHiLifeFetchR28(engine, client)

	// 店舖驗收日檔 09:00 RS4124HIL20201126.XML.zip 有問題 無法解密 也無法移動
	MartHiLifeFetchRS4(engine, client)

	// 驗收異常檔 09:00
	MartHiLifeFetchRS9(engine, client)

	// 取貨即時檔 整點
	MartHiLifeFetchR29(engine, client)

	// 取貨日檔 09:00 有問題 無法解密 也無法移動
	MartHiLifeFetchR96(engine, client)

	// 刷退檔 20:00
	MartHiLifeFetchR08(engine, client)



	log.Info("MartHiLifeFetchShipping[End]")
}

func getMartHiLifeStore(engine *database.MysqlSession) error {
	t := `R00`

	var client MartHiLife.Client
	client.GetClient()

	if client.ConnectFtpError != nil {
		log.Error("MartHiLifeFetchShipping [%v]", client.ConnectFtpError.Error())
		return client.ConnectFtpError
	}

	defer client.ConnectFtp.Close()

	fileNames, err := client.FetchFolder(t, "zip")
	if err != nil {
		return err
	}

	data, err := client.RetrFileAndUnzip(t, fileNames[0])
	if err != nil {
		return err
	}
	var xmlData R00Xml

	err = xml.Unmarshal([]byte(data), &xmlData)
	// body, err := client.GetR00Xml(data)

	if err != nil {
		client.MoveFileToDest(t, fileNames[0], `ERR`)
		return err
	}

	for _, data := range xmlData.R00 {
		country, district, address, flag := tw_address.AddressSplit(data.StoreAddress)
		if flag {
			store := entity.MartHiLifeStoreData{
				StoreId:        data.StoreId,
				StoreName:      data.StoreName,
				StoreAddress:   address,
				StoreCloseDate: data.StoreCloseDate,
				MdcStareDate:   data.MdcStareDate,
				MdcEndDate:     data.MdcEndDate,
				TelNo:          data.TelNo,
				City:           country,
				District:       district,
				UpdateTime:     time.Now(),
			}

			exists := new(entity.MartHiLifeStoreData)

			query := map[string]interface{}{}
			query["store_id"] = data.StoreId
			exist, err := engine.Engine.Where(query).Exist(exists)

			if err != nil {
				log.Error("MartHiLifeFetchStoreList Exist", err)
				continue
			}

			if exist {
				if _, err := engine.Session.Where(query).Update(store, entity.MartHiLifeStoreData{StoreId: data.StoreId}); err != nil {
					log.Error("MartHiLifeFetchStoreList Update", err)
				}
			} else {
				_, err := engine.Session.InsertOne(store)
				if err != nil {
					log.Error("MartHiLifeFetchStoreList InsertOne", err)
				}
			}
		}

	}
	client.MoveFileToDest(t, fileNames[0], `OK`)

	return nil
}

func MartHiLifeFetchStoreList() {
	engine := database.GetMysqlEngine()
	defer engine.Close()
	// 從 Ftp 更新店舖資訊進 DB
	_ = getMartHiLifeStore(engine)

	var newStores []entity.MartHiLifeStoreData
	err := engine.Engine.Find(&newStores)
	if err != nil {
		log.Error("MartFamilyFetchFetchStoreList:", err)
		return
	}

	_ = Mart.WriteHiLifeStoreData(newStores)
}

// 即時交寄檔
func MartHiLifeFetchR27(engine *database.MysqlSession, client MartHiLife.Client) {
	t := `R27`

	fileNames, err := client.FetchFolder(t, "zip")
	if err != nil {
		return
	}

	var UpdateCvsShipping CvsShipping.UpdateCvsShipping
	UpdateCvsShipping.Type = t
	UpdateCvsShipping.ShipType = Enum.CVS_HI_LIFE

	for _, f := range fileNames {
		data, err := client.RetrFileAndUnzip(t, f)
		if err != nil {
			return
		}
		body, err := client.GetR27Xml(data)
		if err != nil {
			client.MoveFileToDest(t, f, `ERR`)
			return
		}
		log.Info(t+`:`+f+`body:`, body)
		for _, d := range body.Body.Contents {
			UpdateCvsShipping.ShipNo = d.OrderNo
			UpdateCvsShipping.DateTime = d.OrderDate
			UpdateCvsShipping.DetailStatus = ""
			UpdateCvsShipping.Log = tools.XmlToString(d)
			UpdateCvsShipping.FileName = f

			err = UpdateCvsShipping.UpdateCvsShippingShipment(engine)

			if err != nil {
				log.Error(t+":"+f+":  Data [%s]:", d)
				log.Error(t+":"+f+": Error [%s]:", err.Error())
			}
		}
		client.MoveFileToDest(t, f, `OK`)
	}
}

// 店舖交寄日檔
func MartHiLifeFetchR22(engine *database.MysqlSession, client MartHiLife.Client) {
	t := `R22`

	fileNames, err := client.FetchFolder(t, "zip")
	if err != nil {
		return
	}

	var UpdateCvsShipping CvsShipping.UpdateCvsShipping
	UpdateCvsShipping.Type = t
	UpdateCvsShipping.ShipType = Enum.CVS_HI_LIFE

	for _, f := range fileNames {
		data, err := client.RetrFileAndUnzip(t, f)
		if err != nil {
			return
		}
		body, err := client.GetR22Xml(data)
		if err != nil {
			client.MoveFileToDest(t, f, `ERR`)
			return
		}
		log.Info(t+`:`+f+`body:`, body)
		for _, d := range body.Body.Contents {
			UpdateCvsShipping.ShipNo = d.OrderNo
			UpdateCvsShipping.DateTime = d.OrderDate
			UpdateCvsShipping.DetailStatus = ""
			UpdateCvsShipping.Log = tools.XmlToString(d)
			UpdateCvsShipping.FileName = f

			// 1:寄件成功
			if d.SendStatus == `1` {
				err = UpdateCvsShipping.UpdateCvsShippingShipment(engine)
			} else {
				// todo 2:取消寄件
			}

			if err != nil {
				log.Error(t+":"+f+":  Data [%s]:", d)
				log.Error(t+":"+f+": Error [%s]:", err.Error())
			}
		}
		client.MoveFileToDest(t, f, `OK`)
	}
}

// DC 驗收檔
func MartHiLifeFetchR04(engine *database.MysqlSession, client MartHiLife.Client) {
	t := `R04`

	fileNames, err := client.FetchFolder(t, "zip")
	if err != nil {
		return
	}

	var UpdateCvsShipping CvsShipping.UpdateCvsShipping
	UpdateCvsShipping.Type = t
	UpdateCvsShipping.ShipType = Enum.CVS_HI_LIFE

	for _, f := range fileNames {
		data, err := client.RetrFileAndUnzip(t, f)
		if err != nil {
			return
		}
		body, err := client.GetR04Xml(data)
		if err != nil {
			client.MoveFileToDest(t, f, `ERR`)
			return
		}
		log.Info(t+`:`+f+`body:`, body)
		for _, d := range body.Body.Contents {
			UpdateCvsShipping.ShipNo = d.OrderNo
			UpdateCvsShipping.DateTime = d.DCReceiveDate
			UpdateCvsShipping.DetailStatus = d.DCReceiveStatus
			UpdateCvsShipping.FlowType = d.FlowType
			UpdateCvsShipping.Log = tools.XmlToString(d)
			UpdateCvsShipping.FileName = f

			log.Info(`UpdateCvsShipping :`, UpdateCvsShipping)

			err = UpdateCvsShipping.UpdateCvsShippingTransit(engine)

			if err != nil {
				log.Error(t+":"+f+":  Data [%s]:", d)
				log.Error(t+":"+f+": Error [%s]:", err.Error())
			}
		}
		client.MoveFileToDest(t, f, `OK`)
	}
}

// 店舖驗收時檔 整點
func MartHiLifeFetchR28(engine *database.MysqlSession, client MartHiLife.Client) {
	t := `R28`

	fileNames, err := client.FetchFolder(t, "zip")
	if err != nil {
		return
	}

	var UpdateCvsShipping CvsShipping.UpdateCvsShipping
	UpdateCvsShipping.Type = t
	UpdateCvsShipping.ShipType = Enum.CVS_HI_LIFE

	for _, f := range fileNames {
		data, err := client.RetrFileAndUnzip(t, f)
		if err != nil {
			return
		}
		body, err := client.GetR28Xml(data)
		if err != nil {
			client.MoveFileToDest(t, f, `ERR`)
			return
		}
		log.Info(t+`:`+f+`body:`, body)
		for _, d := range body.Body.Contents {
			UpdateCvsShipping.ShipNo = d.OrderNo
			UpdateCvsShipping.DateTime = d.DCReceiveDate
			UpdateCvsShipping.DetailStatus = d.DCReceiveStatus
			UpdateCvsShipping.FlowType = d.FlowType
			UpdateCvsShipping.Log = tools.XmlToString(d)
			UpdateCvsShipping.FileName = f

			err = UpdateCvsShipping.UpdateCvsShippingShop(engine)

			if err != nil {
				log.Error(t+":"+f+":  Data [%s]:", d)
				log.Error(t+":"+f+": Error [%s]:", err.Error())
			}
		}
		client.MoveFileToDest(t, f, `OK`)
	}
}

// 店舖驗收日檔
func MartHiLifeFetchRS4(engine *database.MysqlSession, client MartHiLife.Client) {

	t := `RS4`

	fileNames, err := client.FetchFolder(t, "zip")
	if err != nil {
		return
	}

	var UpdateCvsShipping CvsShipping.UpdateCvsShipping
	UpdateCvsShipping.Type = t
	UpdateCvsShipping.ShipType = Enum.CVS_HI_LIFE

	for _, f := range fileNames {
		data, err := client.RetrFileAndUnzip(t, f)
		if err != nil {
			return
		}
		body, err := client.GetRS4Xml(data)
		if err != nil {
			client.MoveFileToDest(t, f, `ERR`)
			return
		}
		log.Info(t+`:`+f+`body:`, body)
		for _, d := range body.Body.Contents {
			UpdateCvsShipping.ShipNo = d.OrderNo
			UpdateCvsShipping.DateTime = d.DCReceiveDate
			UpdateCvsShipping.DetailStatus = d.DCReceiveStatus
			UpdateCvsShipping.FlowType = d.FlowType
			UpdateCvsShipping.Log = tools.XmlToString(d)
			UpdateCvsShipping.FileName = f

			err = UpdateCvsShipping.UpdateCvsShippingShop(engine)

			if err != nil {
				log.Error(t+":"+f+":  Data [%s]:", d)
				log.Error(t+":"+f+": Error [%s]:", err.Error())
			}
		}
		client.MoveFileToDest(t, f, `OK`)
	}
}

// 貨態-物品領取
func MartHiLifeFetchR29(engine *database.MysqlSession, client MartHiLife.Client) {
	t := `R29`

	fileNames, err := client.FetchFolder(t, "zip")
	if err != nil {
		return
	}

	var UpdateCvsShipping CvsShipping.UpdateCvsShipping
	UpdateCvsShipping.Type = t
	UpdateCvsShipping.ShipType = Enum.CVS_HI_LIFE

	for _, f := range fileNames {
		data, err := client.RetrFileAndUnzip(t, f)
		if err != nil {
			return
		}
		body, err := client.GetR29Xml(data)
		if err != nil {
			client.MoveFileToDest(t, f, `ERR`)
			return
		}
		log.Info(t+`:`+f+`body:`, body)
		for _, d := range body.Body.Contents {
			UpdateCvsShipping.ShipNo = d.OrderNo
			UpdateCvsShipping.DateTime = d.SPDate
			UpdateCvsShipping.DetailStatus = ``
			UpdateCvsShipping.FlowType = d.FlowType
			UpdateCvsShipping.Log = tools.XmlToString(d)
			UpdateCvsShipping.FileName = f

			err = UpdateCvsShipping.UpdateCvsShippingSuccess(engine)

			if err != nil {
				log.Error(t+":"+f+":  Data [%s]:", d)
				log.Error(t+":"+f+": Error [%s]:", err.Error())
			}
		}
		client.MoveFileToDest(t, f, `OK`)
	}
}

// 貨態-物品領取(日)
func MartHiLifeFetchR96(engine *database.MysqlSession, client MartHiLife.Client) {
	t := `R96`

	fileNames, err := client.FetchFolder(t, "zip")
	if err != nil {
		return
	}

	var UpdateCvsShipping CvsShipping.UpdateCvsShipping
	UpdateCvsShipping.Type = t
	UpdateCvsShipping.ShipType = Enum.CVS_HI_LIFE

	for _, f := range fileNames {
		data, err := client.RetrFileAndUnzip(t, f)
		if err != nil {
			return
		}
		body, err := client.GetR96Xml(data)
		if err != nil {
			client.MoveFileToDest(t, f, `ERR`)
			return
		}
		log.Info(t+`:`+f+`body:`, body)
		for _, d := range body.Body.Contents {
			UpdateCvsShipping.ShipNo = d.OrderNo
			UpdateCvsShipping.DateTime = d.SPDate
			UpdateCvsShipping.DetailStatus = ``
			UpdateCvsShipping.FlowType = d.FlowType
			UpdateCvsShipping.Log = tools.XmlToString(d)
			UpdateCvsShipping.FileName = f

			err = UpdateCvsShipping.UpdateCvsShippingSuccess(engine)

			if err != nil {
				log.Error(t+":"+f+":  Data [%s]:", d)
				log.Error(t+":"+f+": Error [%s]:", err.Error())
			}
		}
		client.MoveFileToDest(t, f, `OK`)
	}
}

// 店舖未取退回物流驗退檔
func MartHiLifeFetchR08(engine *database.MysqlSession, client MartHiLife.Client) {
	t := "R08"

	fileNames, err := client.FetchFolder(t, "zip")
	if err != nil {
		return
	}

	var UpdateCvsShipping CvsShipping.UpdateCvsShipping
	UpdateCvsShipping.Type = t
	UpdateCvsShipping.ShipType = Enum.CVS_HI_LIFE

	for _, f := range fileNames {
		data, err := client.RetrFileAndUnzip(t, f)
		if err != nil {
			return
		}
		body, err := client.GetR08Xml(data)
		if err != nil {
			client.MoveFileToDest(t, f, `ERR`)
			return
		}
		log.Info(t+`:`+f+`body:`, body)
		for _, d := range body.Body.Contents {
			UpdateCvsShipping.ShipNo = d.OrderNo
			UpdateCvsShipping.DateTime = d.DCReturnDate
			UpdateCvsShipping.DetailStatus = ``
			UpdateCvsShipping.FlowType = d.FlowType
			UpdateCvsShipping.Log = tools.XmlToString(d)
			UpdateCvsShipping.FileName = f

			err = UpdateCvsShipping.UpdateCvsShippingBuyerNotPickUp(engine)

			if err != nil {
				log.Error(t+":"+f+":  Data [%s]:", d)
				log.Error(t+":"+f+": Error [%s]:", err.Error())
			}
		}
		client.MoveFileToDest(t, f, `OK`)
	}
}

// 驗收異常檔
func MartHiLifeFetchRS9(engine *database.MysqlSession, client MartHiLife.Client) {

	t := `RS9`

	fileNames, err := client.FetchFolder(t, "zip")
	if err != nil {
		return
	}

	var UpdateCvsShipping CvsShipping.UpdateCvsShipping
	UpdateCvsShipping.Type = t
	UpdateCvsShipping.ShipType = Enum.CVS_HI_LIFE

	for _, f := range fileNames {
		data, err := client.RetrFileAndUnzip(t, f)
		if err != nil {
			return
		}
		body, err := client.GetRS9Xml(data)
		if err != nil {
			client.MoveFileToDest(t, f, `ERR`)
			return
		}
		log.Info(t+`:`+f+`body:`, body)
		for _, d := range body.Body.Contents {
			UpdateCvsShipping.ShipNo = d.OrderNo
			UpdateCvsShipping.DateTime = d.DCReceiveDate
			UpdateCvsShipping.DetailStatus = d.StatusDetails
			UpdateCvsShipping.FlowType = d.FlowType
			UpdateCvsShipping.Log = tools.XmlToString(d)
			UpdateCvsShipping.FileName = f

			if UpdateCvsShipping.DetailStatus == `T01` {
				err = UpdateCvsShipping.UpdateCvsShippingSwitch(engine)
			} else {
				err = UpdateCvsShipping.UpdateCvsShippingFail(engine)
			}

			if err != nil {
				log.Error(t+":"+f+":  Data [%s]:", d)
				log.Error(t+":"+f+": Error [%s]:", err.Error())
			}
		}
		client.MoveFileToDest(t, f, `OK`)
	}
}

// 貨態-寄件運費檔
func MartHiLifeFetchR98(engine *database.MysqlSession, client MartHiLife.Client) {

	t := `R98`

	fileNames, err := client.FetchFolder(t, "zip")
	if err != nil {
		return
	}

	for _, f := range fileNames {
		data, err := client.RetrFileAndUnzip(t, f)
		if err != nil {
			return
		}
		body, err := client.GetR98Xml(data)
		if err != nil {
			client.MoveFileToDest(t, f, `ERR`)
			return
		}

		var accountingData entity.CvsAccountingData
		accountingData.SetType(t, Enum.CVS_HI_LIFE)

		log.Info(t+`:`+f+`body:`, body)
		for _, d := range body.Body.Contents {
			accountingData.DataId = d.OrderNo
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

// 貨態-核帳檔
func MartHiLifeFetchR99(engine *database.MysqlSession, client MartHiLife.Client) {

	t := `R99`

	fileNames, err := client.FetchFolder(t, "zip")
	if err != nil {
		return
	}

	for _, f := range fileNames {
		data, err := client.RetrFileAndUnzip(t, f)
		if err != nil {
			return
		}
		body, err := client.GetR99Xml(data)
		if err != nil {
			client.MoveFileToDest(t, f, `ERR`)
			return
		}

		var accountingData entity.CvsAccountingData
		accountingData.SetType(t, Enum.CVS_HI_LIFE)

		log.Info(t+`:`+f+`body:`, body)
		for _, d := range body.Body.Contents {
			accountingData.DataId = d.OrderNo
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

// 貨態-遺失
func MartHiLifeFetchR89(engine *database.MysqlSession, client MartHiLife.Client) {

	t := `R89`

	fileNames, err := client.FetchFolder(t, "zip")
	if err != nil {
		return
	}

	for _, f := range fileNames {
		data, err := client.RetrFileAndUnzip(t, f)
		if err != nil {
			return
		}
		body, err := client.GetR89Xml(data)
		if err != nil {
			client.MoveFileToDest(t, f, `ERR`)
			return
		}

		var accountingData entity.CvsAccountingData
		accountingData.SetType(t, Enum.CVS_HI_LIFE)

		log.Info(t+`:`+f+`body:`, body)
		for _, d := range body.Body.Contents {
			accountingData.DataId = d.OrderNo
			accountingData.Amount, _ = strconv.ParseFloat(d.OrderAmount, 64)
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
