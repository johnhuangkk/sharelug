package Shipment

import (
	"api/services/Enum"
	"api/services/Service/FamilyApi"
	"api/services/Service/HiLifeApi"
	"api/services/Service/IPost"
	"api/services/Service/OKApi"
	"api/services/Service/PostBag"
	"api/services/Service/SevenMyshipApi"
	"api/services/dao/Cvs"
	"api/services/model"
	"api/services/util/tools"
	"strings"
	"time"

	//"api/services/Enum"
	"api/services/Service/UserAddressService"
	"api/services/VO/ShipmentVO"
	"api/services/dao/Orders"
	"api/services/dao/member"
	"api/services/database"
	"api/services/entity"
	"api/services/util/log"
	"fmt"
)

// 確認訂單賣家擁有者
func checkOrderSellerOwner(engine *database.MysqlSession, orders ShipmentVO.Orders, storeId string) error {
	ss := tools.StringArrayToInterface(orders.OrderId)

	d, err := Orders.DistinctSellerOrdersOwner(engine, ss)

	if err != nil || len(d) > 1 {
		log.Error("checkOrderSellerOwner Error:[%s]", d)
		return fmt.Errorf("系統錯誤")
	}

	if d[0] != storeId {
		log.Error("checkOrderSellerOwner Error:[%s]", "非此賣家不可操作動作")
		return fmt.Errorf("非此賣家不可操作動作")
	}

	return nil
}

// 檢查訂單
func checkOrder(engine *database.MysqlSession, orderId, uid, sid string) (entity.OrderData, error) {

	log.Debug(`checkOrder [%s]`, orderId, uid, sid)

	var orderData entity.OrderData
	orderData, _ = Orders.GetOrderByOrderId(engine, orderId)

	if len(orderData.OrderId) == 0 {
		log.Error("無此單號: ", orderId)
		return orderData, fmt.Errorf("%s: [%s]", "無此單號", orderId)
	}

	if orderData.StoreId != sid && orderData.BuyerId != uid {
		return orderData, fmt.Errorf("非訂單擁有者")
	}

	return orderData, nil
}

// i郵箱＆宅配
func getSellerSendAddress(engine *database.MysqlSession, orderData entity.OrderData) (*ShipmentVO.SellerSenderAddress, error, int64) {
	ship := orderData.ShipType
	ship = "DELIVERY" // 特規因露天未上線需設定為面交地址 todo check

	sellerAddress := &ShipmentVO.SellerSenderAddress{}
	// 檢查是否填寫寄送地址
	if boolean := UserAddressService.CheckShipSendAddressExist(engine, ship, orderData.SellerId); !boolean {
		return sellerAddress, fmt.Errorf(`請先設定買家未取貨時的退件地址與收件人姓名。`), 1002105
	}

	switch ship {
	case Enum.I_POST:
		address := UserAddressService.GetSendDefaultAddressByShip(engine, orderData.ShipType, orderData.SellerId)
		iPostZipAddress := model.GetPostBoxAddressById(engine, address.Address)
		if iPostZipAddress.Status == "N" {
			log.Error("郵局 i郵箱 狀態: ", address)
			return sellerAddress, fmt.Errorf("i郵箱 [%s] 不提供服務 ", iPostZipAddress.Alias), 200
		}
		sellerAddress.Id = iPostZipAddress.Id
		sellerAddress.Zip = iPostZipAddress.Zip
		sellerAddress.Address = iPostZipAddress.Address
	default:
		address := UserAddressService.GetSendDefaultAddressByShip(engine, "DELIVERY", orderData.SellerId)
		addressSplit := strings.Split(address.Address, ",")
		sellerAddress.Id = ""
		sellerAddress.Zip = addressSplit[0]
		sellerAddress.Address = strings.Join(addressSplit[1:], "")
	}

	return sellerAddress, nil, 200
}

/**
準備訂單相關資訊
*/
func prepareOrderRelationData(engine *database.MysqlSession, orderData entity.OrderData) (error, int64) {

	var shipNo string
	var err error = nil
	var code int64 = 200
	var sellerAddress *ShipmentVO.SellerSenderAddress
	var expire = 7
	// 非正式環境取號 過期時間一年
	if !tools.EnvIsProduction() {
		expire = 365
	}

	sellerData, _ := member.GetMemberDataByUid(engine, orderData.SellerId)

	// UserAddress 寄件地址連動 member SendName 欄位
	if len(sellerData.SendName) == 0 {
		return fmt.Errorf(`請先設定買家未取貨時的退件地址與收件人姓名。`), 5555
	}

	// 寄件期限 i郵箱會用到 需提前填入
	tNow := time.Now().AddDate(0, 0, expire)
	orderData.ShipExpire = time.Date(tNow.Year(), tNow.Month(), tNow.Day(), 23, 59, 59, 999999999, tNow.Location())

	log.Debug("ShipType : [%s]", orderData.ShipType)

	switch orderData.ShipType {
	case Enum.CVS_7_ELEVEN:
		shipNo, err = SevenMyshipApi.CreateShipOrder(orderData, sellerData)
	case Enum.CVS_FAMILY:
		shipNo, err = FamilyApi.C2cOrderAdd(engine, orderData, sellerData)
	case Enum.CVS_HI_LIFE:
		shipNo, err = HiLifeApi.GetShippingOrderNo(engine, orderData, sellerData)
	case Enum.CVS_OK_MART:
		shipNo, err = OKApi.GetShippingOrderNo(engine, orderData, sellerData)
	case Enum.I_POST, Enum.DELIVERY_I_POST_BAG1:
		sellerAddress, err, code = getSellerSendAddress(engine, orderData)
		if err != nil {
			return err, 5555
		}
		shipNo, err = IPost.AddShippingOrder(engine, orderData, sellerData, sellerAddress)
	case Enum.DELIVERY_POST_BAG1, Enum.DELIVERY_POST_BAG2, Enum.DELIVERY_POST_BAG3:
		sellerAddress, err, code = getSellerSendAddress(engine, orderData)
		if err != nil {
			log.Error(err.Error(), "便利包取寄件地址失敗")
			return err, 5555
		}
		shipNo, err = PostBag.CreateShipNumber(engine, orderData, sellerData, sellerAddress)

	default:
		log.Debug("ShipType : [%s]", orderData.ShipType)
		err = fmt.Errorf("不支援此物流")
	}

	log.Debug("shipNo : [%s]", shipNo)

	if err != nil || len(shipNo) == 0 {
		log.Error("取號失敗 %v", orderData)
		log.Error("取號失敗 error %v", err)
		return err, 66666
	}

	// 寫入訂單託運單號
	_, err = model.WriteOrderDataShipNumber(engine, shipNo, orderData)

	if err != nil {
		return fmt.Errorf("[%s] 取號成功更新訂單失敗", orderData.OrderId), 77777
	}

	return nil, code
}

// 取號
func GetShipNumber(orders ShipmentVO.Orders, uid, sid string) (error, int64) {
	engine := database.GetMysqlEngine()
	defer engine.Close()
	err := checkOrderSellerOwner(engine, orders, sid)
	if err != nil {
		return err, 2222
	}
	for _, o := range orders.OrderId {
		orderData, err := checkOrder(engine, o, uid, sid)
		if err != nil {
			return fmt.Errorf(err.Error()), 23456
		}
		if len(orderData.ShipNumber) != 0 {
			log.Error("已產生寄件編號不可重複產生: ", o)
			return fmt.Errorf("%s :已產生寄件編號不可重複產生", o), 343434
		}
		err, code := prepareOrderRelationData(engine, orderData)
		if err != nil {
			return err, code
		}
	}
	return nil, 200
}

// 確認要相同配送方式
func checkShipType(engine *database.MysqlSession, orders ShipmentVO.Orders) (string, error) {
	ss := tools.StringArrayToInterface(orders.OrderId)

	d, err := Orders.DistinctShipOrders(engine, ss)

	if err != nil || len(d) > 1 {
		return "", fmt.Errorf("不允許不同配送同時取得托運單")
	}
	return d[0], nil
}

// 超商取得托運單 需轉換為托運單號再處理
func cvsShip(engine *database.MysqlSession, ship string, orderId []string) (interface{}, error) {
	var data interface{}
	var err error
	shipSevenNumbers := map[string][]string{}

	shipNumbers, err := Orders.GetShipNumber(engine, tools.StringArrayToInterface(orderId))
	if ship == Enum.CVS_7_ELEVEN {

		shipSevenNumbers, err = Orders.GetShipNumberGroupByPayWay(tools.StringArrayToInterface(orderId))

		if err != nil {
			return shipSevenNumbers, err
		}
	}
	log.Debug("shipNumbers", shipNumbers)
	log.Debug("shipSevenNumbers", shipSevenNumbers)

	if err != nil {
		return data, err
	}

	switch ship {
	case Enum.CVS_7_ELEVEN:
		data, err = SevenMyshipApi.PrintShippingOrder(shipSevenNumbers)
	case Enum.CVS_FAMILY:
		data, err = FamilyApi.PrintShippingOrder(shipNumbers)
	case Enum.CVS_HI_LIFE:
		data, err = HiLifeApi.PrintShippingOrder(shipNumbers)
	case Enum.CVS_OK_MART:
		cvsData, err := Cvs.GetCvsShippingDataByShipNo(engine, shipNumbers)
		if err == nil {
			data, err = OKApi.PrintShippingOrderX(cvsData)
		}
	}

	if err != nil {
		return data, err
	}

	return data, nil
}

// 取得托運單
func GetConsignment(orders ShipmentVO.Orders, StoreDataResp entity.StoreDataResp) (string, interface{}, error) {

	var data interface{}
	engine := database.GetMysqlEngine()
	defer engine.Close()

	err := checkOrderSellerOwner(engine, orders, StoreDataResp.StoreId)

	if err != nil {
		return "", data, err
	}

	ship, err := checkShipType(engine, orders)

	if err != nil {
		return ship, data, err
	}

	switch ship {
	case Enum.CVS_7_ELEVEN, Enum.CVS_FAMILY, Enum.CVS_HI_LIFE, Enum.CVS_OK_MART:
		data, err = cvsShip(engine, ship, orders.OrderId)
	case Enum.DELIVERY_POST_BAG1, Enum.DELIVERY_POST_BAG2, Enum.DELIVERY_POST_BAG3:
		data, err = PostBag.GetConsignment(orders.OrderId)
	default:
		data = IPost.GetPostConsignmentData(orders.OrderId, StoreDataResp)
	}

	if err != nil {
		return ship, data, err
	}

	return ship, data, nil

}

// 取得超商配送資訊
func GetCvsShippingData(params ShipmentVO.Order, uid, sid string) (entity.CvsShippingData, error) {

	var cvsShippingData entity.CvsShippingData

	engine := database.GetMysqlEngine()
	defer engine.Close()

	orderData, err := checkOrder(engine, params.OrderId, uid, sid)

	if err != nil {
		return cvsShippingData, err
	}

	return Cvs.GetCvsShippingData(engine, orderData), nil
}

// 閉轉
func SwitchOrder(params ShipmentVO.SwitchOrder, uid, sid string) error {
	engine := database.GetMysqlEngine()
	defer engine.Close()

	orderData, err := checkOrder(engine, params.OrderId, uid, sid)

	log.Info("orderData [%v]", orderData)
	if err != nil {
		return err
	}

	cvsShippingData := Cvs.GetCvsShippingData(engine, orderData)

	log.Info("cvsShippingData [%v]", cvsShippingData)

	if len(cvsShippingData.ShipNo) == 0 {
		return fmt.Errorf("系統錯誤")
	}

	switch orderData.ShipType {
	case Enum.CVS_7_ELEVEN:
	case Enum.CVS_FAMILY:
		err = FamilyApi.SwitchStore(cvsShippingData, params.StoreId)
	case Enum.CVS_HI_LIFE:
		switchLog := Cvs.GetOneCvsShippingLogData(engine, cvsShippingData.ShipNo, `Switch`)
		err = HiLifeApi.SwitchStore(cvsShippingData, switchLog, params.StoreId)
	case Enum.CVS_OK_MART:
		err = OKApi.SwitchStore(orderData, cvsShippingData, params.StoreId)
	}

	if err != nil {
		log.Error("SwitchStore Error [%v]", err)
		return fmt.Errorf("關轉換號失敗")
	}

	// 關閉 閉轉
	cvsShippingData.Switch = `0`
	if cvsShippingData.FlowType == `R` {
		//逆向閉轉需寫入店家
		cvsShippingData.SenderStoreId = params.StoreId
	} else {
		//順向閉轉需改變訂單
		orderData.ReceiverAddress = params.StoreId
		cvsShippingData.SwitchReceiverAddress = params.StoreId
		_, err = Orders.UpdateOrderData(engine, orderData.OrderId, orderData)
		if err != nil {
			log.Error("訂單更新失敗 [%v]", orderData)
		}
	}

	return Cvs.UpdateCvsShippingData(engine, cvsShippingData)

}
