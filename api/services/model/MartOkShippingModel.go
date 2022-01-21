package model

import (
	"api/services/Enum"
	om "api/services/Service/OKMart"
	"api/services/dao/Mart"
	"api/services/dao/Orders"
	"api/services/database"
	"api/services/entity"
	"api/services/util/log"
	"api/services/util/tools"
	"errors"
	"fmt"
	"github.com/spf13/viper"
)

var okClient om.Client

var martOkEcNo1 string
var martOkEcNo2 string
var martOkEcNo3 string
var martOkEcNo4 string
var martOkEcNo5 string

func init() {
	martOkEcNo1 = "461"
	martOkEcNo2 = "462"
	martOkEcNo3 = "463"
	martOkEcNo4 = "464"
	martOkEcNo5 = "465"

	okClient = om.Client{
		ApiHost: "https://Ecrelay2018.okmart.com.tw:8443",
		EcNo1: "461",
		EcNo2: "462",
		EcNo3: "463",
		EcNo4: "464",
		EcNo5: "465",
		F60StserCode: "951",
		F60SerCode: "952",
		F10SerCode: "953",

		FtpHost: "40.81.31.181:21",
		FtpUsername: "VSHALUG",
		FtpPassword: "2bWtZgZ6",
		KeyCode: "5aKP4mayvbnMiW/ENGHcQcXqBexzXgEGPsS6cl7PxCI=",
	}
}

//func MartOkAddShippingOrder(ecOrderNo, orderDt, rStoreId, rN, rP, sStoreId, sN, sP, amount string, needPay bool) (orderNo string, err error) {
//
//	result ,err := okClient.OrderAdd(ecOrderNo, sN, sP, rN, rP, rStoreId, amount, needPay)
//	if err != nil {
//		return "",err
//	}
//
//	if result.Body.Order.ErrCode != "000" {
//		return "",errors.New("["+result.Body.Order.ErrCode +"]" + result.Body.Order.ErrDesc)
//	}
//
//	procDateTime := result.Body.Return.ProcDate
//
//	state := "訂單成立"
//	record := generateRecord(procDateTime,"","",state)
//
//	shipNo := result.Body.Order.OdNo
//
//	de,te := tools.Now("Ymd"),tools.Now("YmdHis")
//	js,_ := json.Marshal(result)
//	lg := generateLog(de,te,state,string(js))
//
//	order := entity.MartOkShippingData{
//		ShipNo:     shipNo,
//		EcOrderNo:  ecOrderNo,
//		ParentId:   "",
//		EshopId:    "",
//		RrStoreId:  rStoreId,
//		SrStoreId:  sStoreId,
//		State:      state,
//		StateCode:  "0",
//		Record:     record,
//		Log:        lg,
//		NeedPay:    false,
//		NeedChange: false,
//		Amount:     amount,
//		ShipFee:    "0",
//		IsLose:     false,
//		RrName:     rN,
//		RrPhone:    rP,
//		SrName:     sN,
//		SrPhone:    sP,
//		OnReturn:   false,
//		CreateDT:   tools.NowYYYYMMDD(),
//	}
//	err = Mart.WriteInsertShippingOrder(order)
//	if err != nil {
//		return "", err
//	}
//
//	return shipNo,nil
//}

func MartOkPrintShippingOrder(shipNo []string) (data []byte,err error) {

	objs := []entity.MartOkShippingData{}
	for _, v := range shipNo {
		obj := entity.MartOkShippingData{ShipNo: v}

		has,err := database.GetMysqlEngineGroup().Get(&obj)
		if err != nil {
			return nil,err
		}
		if has {
			objs = append(objs,obj)
		}
	}

	if len(shipNo) != len(objs) {
		log.Error("MartOkPrintShippingOrder:多筆單號列印失敗")
		return nil,errors.New("Multiple shipNo not found")
	}


	shipNos := []string{}
	tels := []string{}

	for _,obj := range objs {
		leng := len(obj.SrPhone)
		if leng < 3  {
			return nil, errors.New("Shipment order not correct")
		}
		// 寄件人手機末三碼
		phoneSubfix := obj.SrPhone[leng-3:]
		shipNos = append(shipNos,obj.ShipNo)
		tels = append(tels,phoneSubfix)
	}

	return okClient.OrderPrint(shipNos,tels)
}

// 寄件改店
func MartOkSwitchStore(shipNo, ecOrderNo, newStoreId string, isReceiveStore bool) (err error) {
	obj := entity.MartOkShippingData{ EcOrderNo: ecOrderNo, ShipNo: shipNo}
	result,err := database.Mysql().Get(&obj)
	if err != nil {
		return err
	}

	if !result {
		return errors.New("Shipment not found")
	}

	l := len(obj.RrPhone)
	nRP := obj.RrPhone[l-3:l]

	newEcNo := martOkEcNo4
	if obj.EshopId == martOkEcNo3 {
		newEcNo = martOkEcNo5
	}

	newOlder, err := okClient.OrderReSend(newEcNo ,obj.ShipNo, newStoreId, obj.RrName, nRP, obj.Amount, `obj.NeedPay`)
	fmt.Println("MartOkSwitchStore:",newOlder)
	d := newOlder.Body.Order
	dirc := "轉店"
	dt := tools.Now("Ymd")
	te := tools.Now("Hms")
	state := "成功"
	log := generateLog(dt,te,state,d)

	if d.ErrorCode != "000" && d.ErrorCode != ""{
		state := "失敗"
		detail := generateRecord(dt,te,dirc,state)
		_ = Mart.WriteOKShipStateWithRecord(obj.ShipNo, state, "1", detail, log)
		return
	}
	// 清除關轉時間
	Mart.UpdateOkSwitchTime(obj.ShipNo, "")
	Mart.UpdateOrderShippingStatus(ecOrderNo,Enum.OrderShipOnShipping)
	detail := generateRecord(dt,te,dirc,state)
	_ = Mart.WriteOKShipStateWithNeedChangeByShipNo(obj.ShipNo, state, "0", detail, false, log)

	return nil
}

func MartOkFetchShipping() {
	log.Debug("MartOkFetchShipping[Begin]")

	ecNo1 := "461"
	{
		// 更新交寄狀態
		MartOKFetchF27(ecNo1)

		// 更新交寄狀態(批次)
		MartOKFetchF25(ecNo1)

		// 更新離店狀態
		MartOKFetchF84(ecNo1)

		// 貨態-物流進貨
		MartOKFetchF71(ecNo1)
	}

	ecNo2 := "462"
	{
		// 貨態-物流出貨
		MartOKFetchF63(ecNo2)
		// 貨態-驗收異常
		MartOKFetchF67(ecNo2)
		// 貨態-到店待取
		MartOKFetchF44(ecNo2)
		// 貨態-到店待取(批次)
		MartOKFetchF64(ecNo2)
		// 貨態-到店已取
		MartOKFetchF17(ecNo2)
		// 貨態-到店已取(批次)
		MartOKFetchF65(ecNo2)
		// 未取離店
		MartOKFetchF84(ecNo2)
		// 物流驗退
		MartOKFetchF67(ecNo2)
	}

	ecNo3 := "464"
	{
		MartOKFetchF03(ecNo3)
		MartOKFetchF07(ecNo3)
		MartOKFetchF44(ecNo3)
		MartOKFetchF04(ecNo3)
		MartOKFetchF17(ecNo3)
		MartOKFetchF05(ecNo3)
		MartOKFetchF84(ecNo3)
		MartOKFetchF07(ecNo3)
	}

	ecNo4 := "463"
	{
		MartOKFetchF03(ecNo4)
		MartOKFetchF07(ecNo4)
		MartOKFetchF44(ecNo4)
		MartOKFetchF04(ecNo4)
		MartOKFetchF17(ecNo4)
		MartOKFetchF05(ecNo4)
		MartOKFetchF84(ecNo4)
		MartOKFetchF07(ecNo4)
	}

	ecNo5 := "465"
	{
		MartOKFetchF03(ecNo5)
		MartOKFetchF07(ecNo5)
		MartOKFetchF44(ecNo5)
		MartOKFetchF04(ecNo5)
		MartOKFetchF17(ecNo5)
		MartOKFetchF05(ecNo5)
		MartOKFetchF84(ecNo5)
		MartOKFetchF07(ecNo5)
	}
	log.Debug("MartOkFetchShipping[End]")

}

func MartOKFetchStoreList() {
	content ,err := okClient.GetF01Document()
	if err != nil {
		log.Debug("MartOKFetchStoreList Err:%v",err.Error())
	}

	engine := database.GetMysqlEngine()
	defer engine.Close()

	if err == nil {
		stores := []entity.MartOkStoreData{}
		for _,s := range content.Contents {
			store := entity.MartOkStoreData{
				StoreId:        s.StoreId,
				StoreName:      s.StoreName,
				StoreAddress:   s.StoreAddress,
				StoreCloseDate: s.SDATE,
				MdcStareDate:   s.SDATE,
				MdcEndDate:     s.EDATE,
				TelNo:          s.StoreTel,
				City:           s.StoreCity,
				District:       s.StoreDisct,
			}
			stores = append(stores,store)

			exist, err := engine.Session.Exist(&store)
			if err != nil {
				log.Debug("MartOKFetchStoreList:%v",err)
				continue
			}

			if exist {
				_ ,err := engine.Session.Update(store,entity.MartOkStoreData{ StoreId:s.StoreId})
				if err != nil {
					log.Debug("MartOKFetchStoreList:%v",err)
				}
			} else {
				_ ,err := engine.Session.InsertOne(store)
				if err != nil {
					log.Debug("MartOKFetchStoreList:%v",err)
				}
			}
		}
	}

	newStores := []entity.MartOkStoreData{}
	err = engine.Engine.Find(&newStores)
	if err != nil {
		log.Debug("MartFamilyFetchFetchStoreList:",err)
		return
	}

	Mart.WriteOkStoreData(newStores)
}

// 貨態-到店交寄
func MartOKFetchF27(ecNo string) {
	log.Debug("MartOKFetchF27[Begin]")

	t := "F27"
	raw,fn,err := okClient.GetRawDocument(ecNo,t)
	if err != nil {
		log.Debug("GetRaw Fail:",t,err.Error())
		return
	}
	err = Mart.InsertFileOK(ecNo,"0", t, fn, tools.Now("Ymd"), raw)
	if err != nil {
		log.Error("Insert Raw Fail:",err.Error())
		return
	}

	content ,err := okClient.GetF27Document(ecNo)
	if err != nil{
		log.Debug("MartOKFetchF27 Fail:",err.Error())
		return
	}

	for _,d := range content.Body {
		state := "已到店寄件"
		dirc := Enum.OrderShipSend
		detail := generateRecord(d.UpDateTime,"",dirc,state)
		log := generateLog(d.UpDateTime,"",state,d)
		_ = Mart.WriteOKShipStateWithRecord(d.OrderNo, state, "1", detail, log)
	}
}

func MartOKFetchF25(ecNo string) {
	//log.Debug("MartOKFetchF25[Begin]")
	//
	//t := "F25"
	//raw,fn,err := okClient.GetRawDocument(ecNo,t)
	//if err != nil {
	//	log.Debug("GetRaw Fail:",t,err.Error())
	//	return
	//}
	//err = Mart.InsertFileOK(ecNo,"0", t, fn, tools.NowYYYYMMDD(), raw)
	//if err != nil {
	//	log.Error("Insert Raw Fail:",err.Error())
	//	return
	//}
	//
	//content ,err := okClient.GetF25Document(ecNo)
	//if err != nil{
	//	log.Debug("MartOKFetchF25 Fail:",err.Error())
	//	return
	//}

	//for _,d := range content.Body {
		//state := Enum.OrderShipment
		//dirc := Enum.OrderShipSend
		//detail := generateRecord(d.UpDateTime,"",dirc,state)
		//log := generateLog(d.UpDateTime,"",state,d)
		//_ = Mart.WriteOKShipStateWithRecord(d.OrderNo, state, "1", detail, log)
	//}
}

// 貨態-寄件離店
func MartOKFetchF84(ecNo string) {
	log.Debug("MartOKFetchF84[Begin]")

	t := "F84"
	raw,fn,err := okClient.GetRawDocument(ecNo,t)
	if err != nil {
		log.Debug("GetRaw Fail:",t,err.Error())
		return
	}
	err = Mart.InsertFileOK(ecNo,"0", t, fn, tools.Now("Ymd"), raw)
	if err != nil {
		log.Error("Insert Raw Fail:",err.Error())
		return
	}

	content ,err := okClient.GetF84Document(ecNo)
	if err != nil{
		log.Debug("MartOKFetchF84 Fail:",err.Error())
		return
	}

	for _,d := range content.Body {
		state := "已離開寄件門市"
		dirc := Enum.OrderShipSend
		detail := generateRecord(d.UpDateTime,"",dirc,state)
		log := generateLog(d.UpDateTime,"",state,d)
		_ = Mart.WriteOKShipStateWithRecord(d.OrderNo, state, "1", detail, log)
	}
}

// 貨態-物流中心驗收1
func MartOKFetchF71(ecNo string) {
	t := "F71"
	raw,fn,err := okClient.GetRawDocument(ecNo,t)
	if err != nil {
		log.Debug("GetRaw Fail:",t,err.Error())
		return
	}
	err = Mart.InsertFileOK(ecNo,"0", t, fn, tools.Now("Ymd"), raw)
	if err != nil {
		log.Error("Insert Raw Fail:",err.Error())
		return
	}

	content ,err := okClient.GetF71Document(ecNo)
	if err != nil{
		log.Debug("MartOKFetchF71 Fail:",err.Error())
		return
	}

	for _,d := range content.Body {
		state := "物流中心進貨驗收"

		dirc := Enum.OrderShipSend

		detail := generateRecord(d.UpDateTime,"",dirc,state)
		lg := generateLog(d.UpDateTime,"",state,d)
		_ = Mart.WriteOKShipStateWithRecord(d.OrderNo, state, "1", detail, lg)
	}

}

// 貨態-物流中心驗收2
func MartOKFetchF63(ecNo string) {
	t := "F63"
	raw,fn,err := okClient.GetRawDocument(ecNo,t)
	if err != nil {
		log.Debug("GetRaw Fail:",t,err.Error())
		return
	}
	err = Mart.InsertFileOK(ecNo,"0", t, fn, tools.Now("Ymd"), raw)
	if err != nil {
		log.Error("Insert Raw Fail:",err.Error())
		return
	}

	content ,err := okClient.GetF63Document(ecNo)
	if err != nil{
		log.Debug("MartOKFetchF63 Fail:",err.Error())
		return
	}

	for _,d := range content.Body {
		state := "物流中心進貨驗收"

		dirc := Enum.OrderShipSend

		detail := generateRecord(d.UpDateTime,"",dirc,state)
		lg := generateLog(d.UpDateTime,"",state,d)
		_ = Mart.WriteOKShipStateWithRecord(d.OrderNo, state, "1", detail, lg)
	}
}

// 貨態-物流中心驗收2(重出)
func MartOKFetchF03(ecNo string) {
	t := "F03"
	raw,fn,err := okClient.GetRawDocument(ecNo,t)
	if err != nil {
		log.Debug("GetRaw Fail:",t,err.Error())
		return
	}
	err = Mart.InsertFileOK(ecNo,"0", t, fn, tools.Now("Ymd"), raw)
	if err != nil {
		log.Error("Insert Raw Fail:",err.Error())
		return
	}

	content ,err := okClient.GetF03Document(ecNo)
	if err != nil{
		log.Debug("MartOKFetchF03 Fail:",err.Error())
		return
	}

	for _,d := range content.Body {
		state := "物流中心進貨驗收"
		dirc := Enum.OrderShipSend
		detail := generateRecord(d.UpDateTime,"",dirc,state)
		lg := generateLog(d.UpDateTime,"",state,d)
		_ = Mart.WriteOKShipStateWithRecord(d.OrderNo, state, "1", detail, lg)
	}
}

// 貨態-到店待取
func MartOKFetchF44(ecNo string) {
	t := "F44"
	raw,fn,err := okClient.GetRawDocument(ecNo,t)
	if err != nil {
		log.Debug("GetRaw Fail:",t,err.Error())
		return
	}
	err = Mart.InsertFileOK(ecNo,"0", t, fn, tools.Now("Ymd"), raw)
	if err != nil {
		log.Error("Insert Raw Fail:",err.Error())
		return
	}

	content ,err := okClient.GetF44Document(ecNo)
	if err != nil {
		log.Debug("MartOKFetchF44 Fail:",err.Error())
		return
	}

	for _,d := range content.Body {
		state := "到店待取"
		dirc := Enum.OrderShipSend
		dStr := d.RrInDateTime
		tStr := ""
		detail := generateRecord(dStr,tStr,dirc,state)
		lg := generateLog(dStr,tStr,state,d)
		_ = writeOkShipStateWithRecord(d.OrderNo, state,"1", detail,lg)
	}
}

func MartOKFetchF04(ecNo string) {
	t := "F04"
	raw,fn,err := okClient.GetRawDocument(ecNo,t)
	if err != nil {
		log.Debug("GetRaw Fail:",t,err.Error())
		return
	}
	err = Mart.InsertFileOK(ecNo,"0", t, fn, tools.Now("Ymd"), raw)
	if err != nil {
		log.Error("Insert Raw Fail:",err.Error())
		return
	}

	content ,err := okClient.GetF04Document(ecNo)
	if err != nil {
		log.Debug("MartOKFetchF04 Fail:",err.Error())
		return
	}

	for _,d := range content.Body {
		state := "到店待取"
		dirc := Enum.OrderShipSend
		dStr := d.UpDateTime
		tStr := ""
		detail := generateRecord(dStr,tStr,dirc,state)
		lg := generateLog(dStr,tStr,state,d)
		_ = writeOkShipStateWithRecord(d.OrderNo, state,"1", detail,lg)
	}
}

func MartOKFetchF64(ecNo string) {
	t := "F64"
	raw,fn,err := okClient.GetRawDocument(ecNo,t)
	if err != nil {
		log.Debug("GetRaw Fail:",t,err.Error())
		return
	}
	err = Mart.InsertFileOK(ecNo,"0", t, fn, tools.Now("Ymd"), raw)
	if err != nil {
		log.Error("Insert Raw Fail:",err.Error())
		return
	}

	content ,err := okClient.GetF64Document(ecNo)
	if err != nil {
		log.Debug("MartOKFetchF64 Fail:",err.Error())
		return
	}

	for _,d := range content.Body {
		state := "已取貨2"
		dirc := Enum.OrderShipSend
		dStr := tools.Now("Ymd")
		tStr := tools.Now("Hms")

		shipNo := d.OrderNo

		detail := generateRecord(dStr,tStr,dirc,state)
		lg := generateLog(dStr,tStr,state,d)
		_ = writeOkShipStateWithRecord(shipNo, state,"1", detail,lg)
	}
}

// 貨態-物品領取
func MartOKFetchF17(ecNo string) {
	t := "F17"
	raw,fn,err := okClient.GetRawDocument(ecNo,t)
	if err != nil {
		log.Debug("GetRaw Fail:",t,err.Error())
		return
	}
	err = Mart.InsertFileOK(ecNo,"0", t, fn, tools.Now("Ymd"), raw)
	if err != nil {
		log.Error("Insert Raw Fail:",err.Error())
		return
	}

	content ,err := okClient.GetF17Document(ecNo)
	if err != nil {
		log.Debug("MartOKFetchF17 Fail:",err.Error())
		return
	}

	for _,d := range content.Body {
		state := "已取貨"
		dirc := Enum.OrderShipSend
		dStr := tools.Now("Ymd")
		tStr := tools.Now("Hms")

		shipNo := om.GenerateOkBarCodeToShipNo(d.BarCode1,d.BarCode2)

		detail := generateRecord(dStr,tStr,dirc,state)
		lg := generateLog(dStr,tStr,state,d)
		_ = writeOkShipStateWithRecord(shipNo, state,"1", detail,lg)
	}
}

func MartOKFetchF65(ecNo string) {
	t := "F65"
	raw,fn,err := okClient.GetRawDocument(ecNo,t)
	if err != nil {
		log.Debug("GetRaw Fail:",t,err.Error())
		return
	}
	err = Mart.InsertFileOK(ecNo,"0", t, fn,  tools.Now("Ymd"), raw)
	if err != nil {
		log.Error("Insert Raw Fail:",err.Error())
		return
	}

	content ,err := okClient.GetF65Document(ecNo)
	if err != nil {
		log.Debug("MartOKFetchF65 Fail:",err.Error())
		return
	}

	for _,d := range content.Body {
		state := "已取貨"
		dirc := Enum.OrderShipSend
		dStr := tools.Now("Ymd")
		tStr := tools.Now("Hms")

		shipNo := om.GenerateOkBarCodeToShipNo(d.BarCode1,d.BarCode2)

		detail := generateRecord(dStr,tStr,dirc,state)
		lg := generateLog(dStr,tStr,state,d)
		_ = writeOkShipStateWithRecord(shipNo, state,"1", detail,lg)
	}
}

func MartOKFetchF05(ecNo string) {
	t := "F05"
	raw,fn,err := okClient.GetRawDocument(ecNo,t)
	if err != nil {
		log.Debug("GetRaw Fail:",t,err.Error())
		return
	}
	err = Mart.InsertFileOK(ecNo,"0", t, fn,  tools.Now("Ymd"), raw)
	if err != nil {
		log.Error("Insert Raw Fail:",err.Error())
		return
	}

	content ,err := okClient.GetF05Document(ecNo)
	if err != nil {
		log.Debug("MartOKFetchF05 Fail:",err.Error())
		return
	}

	for _,d := range content.Body {
		state := "C2 已取貨"
		dirc := Enum.OrderShipSend
		dStr := tools.Now("Ymd")
		tStr := tools.Now("Hms")

		shipNo := om.GenerateOkBarCodeToShipNo(d.BarCode1,d.BarCode2)

		detail := generateRecord(dStr,tStr,dirc,state)
		lg := generateLog(dStr,tStr,state,d)
		_ = writeOkShipStateWithRecord(shipNo, state,"1", detail,lg)
	}
}

// 貨態-寄件驗退
func MartOKFetchF67(ecNo string) {
	t := "F67"
	raw,fn,err := okClient.GetRawDocument(ecNo,t)
	if err != nil {
		log.Debug("GetRaw Fail:",t,err.Error())
		return
	}
	err = Mart.InsertFileOK(ecNo,"0", t, fn,  tools.Now("Ymd"), raw)
	if err != nil {
		log.Error("Insert Raw Fail:",err.Error())
		return
	}

	content ,err := okClient.GetF67Document(ecNo)
	if err != nil{
		log.Debug("MartOKFetchF67 Fail:",err.Error())
		return
	}

	for _,d := range content.Body {
		str := om.GenerateOkReturnToString(d.ReturnCode)
		dirc,state := Enum.OrderShipSend,"[" + d.ReturnCode +"]物流中心進貨驗退(" + str + ")"
		needChange := d.ReturnCode == "T01"
		shipNo := d.OrderNo

		detail := generateRecord(d.UpDateTime,"",dirc,state)
		lg := generateLog(d.UpDateTime,"",state,d)

		if d.ReturnCode == "T00" {
			Mart.UpdateOnReturnOK(shipNo,true)
		}

		if needChange {
			state = Enum.OrderShipReceiverStoreSwitch
			_ = Mart.UpdateEshopIdOK(shipNo,d.EcNo)
			// 更新關轉等待時間為3天
			Mart.UpdateOkSwitchTime(shipNo, tools.NowYYYYMMDDHHmmss(3))
		}
		_ = Mart.WriteOKShipStateWithNeedChangeByShipNo(shipNo, state, "1", detail, needChange, lg)

		obj,err := Mart.QueryOKShippingData(d.VendorNo)
		if err != nil {
			log.Error("MartOKFetchF67:",err)
			continue
		}

		engine := database.GetMysqlEngine()
		defer engine.Close()
		shipStatus := Enum.OrderShipReceiverStoreSwitch
		if obj.OnReturn {
			shipStatus = Enum.OrderShipSenderStoreSwitch
		}
		c,err := Orders.UpdateOrderData(engine,obj.EcOrderNo, entity.OrderData{ShipStatus: shipStatus })
		if err != nil || c != 1 {
			log.Error("MartOKFetchF67:", err, c)
		}
	}
}

// 貨態-驗收異常含關轉
func MartOKFetchF07(ecNo string) {
	t := "F07"
	raw,fn,err := okClient.GetRawDocument(ecNo,t)
	if err != nil {
		log.Debug("GetRaw Fail:",t,err.Error())
		return
	}
	err = Mart.InsertFileOK(ecNo,"0", t, fn,  tools.Now("Ymd"), raw)
	if err != nil {
		log.Error("Insert Raw Fail:",err.Error())
		return
	}

	content ,err := okClient.GetF07Document(ecNo)
	if err != nil{
		log.Debug("MartOKFetchF07 Fail:",err.Error())
		return
	}

	for _,d := range content.Body {
		str := om.GenerateOkReturnToString(d.ReturnCode)
		dirc,state := Enum.OrderShipSend,"[" + d.ReturnCode +"]驗收異常(" + str + ")"
		needChange := d.ReturnCode == "T01"
		shipNo := d.OrderNo

		Mart.UpdateEshopIdOK(shipNo, d.EcNo)

		if d.EcNo == martOkEcNo4 {
			Mart.UpdateOnReturnOK(shipNo,true)
		}

		detail := generateRecord(d.UpDateTime,"",dirc,state)
		lg := generateLog(d.UpDateTime,"",state,d)
		_ = Mart.WriteOKShipStateWithNeedChangeByShipNo(shipNo, state, "1", detail, needChange, lg)

		obj,err := Mart.QueryOKShippingDataByShipNo(d.OrderNo)
		if err != nil {
			log.Error("MartOKFetchF07:",err)
			continue
		}

		shipStatus := Enum.OrderShipReceiverStoreSwitch
		if obj.OnReturn {
			shipStatus = Enum.OrderShipSenderStoreSwitch
		}
		Mart.UpdateOrderShippingStatus(obj.EcOrderNo,shipStatus)

	}
}

func writeOkShipStateWithRecord(shipNo, state, stateCode, record string, lg string) error {
	session := database.GetMysqlEngineGroup()

	result,err := session.Exec("UPDATE mart_ok_shipping_data SET record=CONCAT(?,record),log=CONCAT(?,log),state=?, state_code=? WHERE ship_no=?", record, lg, state, stateCode, shipNo)
	if err != nil {
		log.Debug("Update:",err.Error())
		return err
	}
	count ,err := result.RowsAffected()
	if err != nil {
		log.Debug("Update:",count,err.Error())
		return err
	}
	if count != 1 {
		log.Debug("Update:",count)
		return err
	}
	return nil
}

func getOkClient(reverseFlow bool) *om.Client {
	key := "MartOK"
	if reverseFlow {
		key = "MartOK"
	}

	apiHost := viper.GetString(key + ".ApiHost")

	ftpHost := viper.GetString(key + ".FtpHost")
	ftpUn := viper.GetString(key + ".FtpUsername")
	ftpPw := viper.GetString(key + ".FtpPassword")

	client := om.Client{
		ApiHost: apiHost,
		FtpHost: ftpHost,
		FtpUsername: ftpUn,
		FtpPassword: ftpPw,
		KeyCode:      "",
		EcNo1:        "461",
		EcNo2:        "462",
		EcNo3:        "463",
		EcNo4:        "464",
		EcNo5:        "465",
		F60StserCode: "951",
		F60SerCode:   "952",
		F10SerCode:   "953",
	}

	return &client
}