package Upgrade

import (
	"api/services/Enum"
	"api/services/Service/Carts"
	"api/services/VO/Request"
	"api/services/VO/Response"
	"api/services/dao/Orders"
	"api/services/dao/Store"
	"api/services/dao/member"
	"api/services/dao/product"
	"api/services/dao/transfer"
	"api/services/database"
	"api/services/entity"
	"api/services/util/log"
	"api/services/util/tools"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

func GetUpgradePlan(userData entity.MemberData, storeData entity.StoreDataResp) (Response.UpgradePlanResponse, error) {
	engine := database.GetMysqlEngine()
	defer engine.Close()
	var resp Response.UpgradePlanResponse
	var current string
	data, err := product.GetUpgradeProductData(engine)
	if err != nil {
		return resp, fmt.Errorf("系統錯誤")
	}
	for _, v := range data {
		var res Response.UpgradePlanList
		res.ProductId = v.ProductId
		res.ProductName = v.ProductName
		res.Description = strings.Split(v.Description, ",")
		res.Note = v.Note
		res.Amount = v.Amount
		res.IsPay = false
		if userData.UpgradeLevel >= v.UpgradeLevel {
			res.IsPay = true
		}
		if userData.UpgradeLevel == v.UpgradeLevel {
			resp.CurrentPlan = res
			text := ""
			if userData.UpgradeLevel == 1 {
				text = fmt.Sprintf("可指派%v位管理員。", v.Manager)
			}
			if userData.UpgradeLevel == 2 || userData.UpgradeLevel == 3 {
				text = fmt.Sprintf("最多可開設%v個收銀機，每個收銀機可指派%v個管理員。", v.Store, v.Manager)
			}
			current = fmt.Sprintf("你目前使用每月%v元方案，%s<br>每月%v元方案到期時間：%s", v.Amount, text, v.Amount, userData.UpgradeExpire.Format("2006/01/02"))
		}
		resp.UpgradePlanList = append(resp.UpgradePlanList, res)
	}

	if len(current) == 0 {
		resp.StoreUpgradeText = fmt.Sprintf("你目前的免費方案，只能由一個主帳號負責管理，無法建立管理帳號讓其他人協助管理收銀機。每月支付最低99元起就能立即升級，讓更多人協助你管理收銀機。")
	} else {
		resp.StoreUpgradeText = current
	}
	resp.StoreName = storeData.StoreName
	resp.StorePicture = storeData.StorePicture
	ManagerLimit, ManagerCurrent, err := ComputeManager(engine, userData, storeData)
	if err != nil {
		return resp, fmt.Errorf("系統錯誤")
	}
	resp.StoreManagerLimit = ManagerLimit
	resp.StoreManagerCurrent = ManagerCurrent
	return resp, nil
}
func GeneratorNewOrder(engine *database.MysqlSession, userData entity.MemberData, storeData entity.StoreDataResp, params Request.GetB2CPayRequest, cookie string, resp *Response.GetB2cPayResponse) (entity.B2cOrder, error) {
	//升級和未付帳單 合併 做法要改 多加 DETAIL 把商品都加在內
	// 包含 商品名稱、單價、類別(升級或帳單)
	var detail entity.B2cOrder
	var pId string
	var pName string
	var uLevel int64
	var expiration time.Time
	var billingTime string
	var productId string
	//抓取未付帳單
	billData, orderSum, err := getUnpaidBill(engine, userData, &detail)
	if err != nil {
		return detail, fmt.Errorf("1001001")
	}
	//帳單總金額
	resp.OrderSum = orderSum
	switch params.ProductId {
		case Enum.ServiceTypeUpgrade:
		case Enum.ServiceTypeShop:
		default:
			productId = params.ProductId
	}
	if len(productId) != 0 {
		productData, err := product.GetUpgradeProductByProductId(engine, productId)
		if err != nil {
			log.Error("Get B2c Product Error", err)
			return detail, fmt.Errorf("1001001")
		}
		pId = productData.ProductId
		pName = productData.ProductName
		uLevel = productData.UpgradeLevel
		expiration = time.Time{}
		upgradeSum, err := getProduct(engine, productData, userData, &detail)
		if err != nil {
			log.Error("Get Product Error", err)
			return detail, fmt.Errorf("1001001")
		}
		//升級方案總金額
		resp.UpgradeSum = upgradeSum
	} else {
		if len(billData) != 0 {
			data := billData[len(billData)-1]
			pId = data.ProductId
			pName = data.ProductName
			uLevel = data.BillingLevel
			expiration = data.Expiration
			billingTime = data.BillingTime
		}
	}
	OrderDetail, _ := tools.JsonEncode(detail.Detail)
	//此筆訂單總金額
	resp.PriceTotal = resp.UpgradeSum + resp.OrderSum
	order, err := GeneratorOrderToCarts(cookie, pId, pName, userData.Uid, storeData.StoreId, string(OrderDetail), billingTime, uLevel, resp.PriceTotal, expiration)
	if err != nil {
		log.Error("Generator B2c Carts Error", err)
		return detail, fmt.Errorf("1001001")
	}
	//此筆訂單編號
	resp.OrderId = order.OrderId
	return detail, nil
}
//取出未付訂單
func getUnpaidBill(engine *database.MysqlSession, UserData entity.MemberData, resp *entity.B2cOrder) ([]entity.B2cBillingData, int64, error) {
	billData, err := Orders.GetB2cUnpaidBillByUserId(engine, UserData.Uid)
	if err != nil {
		return billData, 0, fmt.Errorf("系統錯誤")
	}
	orderTotal := 0
	if len(billData) > 0 {
		for _, v := range billData {
			var detail entity.B2cOrderDetail
			detail.ProductId = v.BillingId
			detail.ProductName = v.BillName
			detail.BillingTime = v.BillingTime
			detail.ProductAmount = v.Amount
			detail.ProductType = Enum.B2cOrderTypeBilling
			detail.ProductDetail = v.ProductDesc
			orderTotal += int(v.Amount)
			resp.Detail = append(resp.Detail, detail)
		}
	}
	return billData, int64(orderTotal), nil
}
//取出購買方案
func getProduct(engine *database.MysqlSession, productData entity.UpgradeProductData, UserData entity.MemberData, resp *entity.B2cOrder) (int64, error) {

	if productData.UpgradeLevel <= UserData.UpgradeLevel {
		return 0, fmt.Errorf("不能選擇此方案")
	}
	data1 := entity.B2cOrderDetail {
		ProductId: productData.ProductId,
		ProductName: fmt.Sprintf("升級加值服務 $%v 方案", productData.Amount),
		ProductAmount: productData.Amount,
		ProductType: "UPGRADE",
	}
	resp.Detail = append(resp.Detail, data1)
	upgradeTotal := productData.Amount
	//方案未過期 UpgradeExpire
	if UserData.UpgradeLevel > 0 {
		data, err := product.GetUpgradeProductDataByLevel(engine, UserData.UpgradeLevel)
		if err != nil {
			return 0, fmt.Errorf("系統錯誤")
		}
		data2 := entity.B2cOrderDetail{
			ProductId: "",
			ProductName:  "當期加值服務費抵扣",
			ProductAmount: data.Amount,
			ProductType: "UPGRADE",
		}
		resp.Detail = append(resp.Detail, data2)
		upgradeTotal -= data.Amount
	}
	return upgradeTotal, nil
}

func GeneratorOrderToCarts(cookie, ProductId, ProductName, UserId, StoreId, OrderDetail, BillingTime string, UpgradeLevel, Amount int64, expire time.Time) (entity.B2cOrderData, error) {
	data := entity.B2cOrderData{}
	data.OrderId = tools.GeneratorB2COrderId()
	data.UserId = UserId
	data.StoreId = StoreId
	data.ProductId = ProductId
	data.ProductName = ProductName
	data.UpgradeLevel = UpgradeLevel
	data.OrderDetail = OrderDetail
	data.Amount = Amount
	data.OrderStatus = Enum.OrderInit
	data.Expiration = expire
	data.BillingTime = BillingTime
	data.OrderSys = 0
	data.AskInvoice = false
	if err := Carts.NewB2cRedisCarts(data, cookie, Enum.StyleB2C); err != nil {
		return data, fmt.Errorf("1001001")
	}
	return data, nil
}
//func GetOrder(engine *database.MysqlSession, UserId string) (entity.B2cOrderData, error) {
//	data, err := Orders.GetB2cOrderLastByUserId(engine, UserId)
//	if err != nil {
//		return data, err
//	}
//	if len(data.OrderId) == 0 {
//		return data, fmt.Errorf("尚未升級")
//	}
//	//檢查是否已升級
//	if !time.Now().Before(data.Expiration) {
//		return data, fmt.Errorf("尚未升級")
//	}
//	return data, nil
//}
//
//func GetUpgradeOrder(engine *database.MysqlSession, UserId string) (entity.B2cOrderData, error) {
//	data, err := Orders.GetB2cOrderLastByUserId(engine, UserId)
//	if err != nil {
//		return data, err
//	}
//	if len(data.OrderId) == 0 {
//		return entity.B2cOrderData{}, err
//	}
//	//檢查是否已升級
//	now := time.Now()
//	if !now.Before(data.Expiration) {
//		return entity.B2cOrderData{}, err
//	}
//	return data, nil
//}
//
//func GetPlanByManager(engine *database.MysqlSession, data entity.B2cOrderData) (int64, error) {
//	resp, err := product.GetUpgradeProductByProductId(engine, data.ProductId)
//	if err != nil {
//		log.Error("Auth handle Error", err)
//		return 0, fmt.Errorf("系統錯誤")
//	}
//	return resp.Manager, nil
//}
//
//func GetPlanByStore(engine *database.MysqlSession, data entity.B2cOrderData) (int64, error) {
//	resp, err := product.GetUpgradeProductByProductId(engine, data.ProductId)
//	if err != nil {
//		log.Error("Auth handle Error", err)
//		return 0, fmt.Errorf("系統錯誤")
//	}
//	return resp.Store, nil
//}

func ComputeManager(engine *database.MysqlSession, UserData entity.MemberData, StoreData entity.StoreDataResp) (int64, int64, error) {
	ManagerMax := 0
	//計算管理者是數量
	ManagerCount := Store.CountStoreManager(engine, StoreData.StoreId)
	if time.Now().Before(UserData.UpgradeExpire) {
		if UserData.UpgradeLevel != 0 {
			//取出升級訂單
			order, err := product.GetUpgradeProductDataByLevel(engine, UserData.UpgradeLevel)
			if err != nil {
				return int64(ManagerMax), ManagerCount, err
			}
			ManagerMax = int(order.Manager)
		}
	}
	return int64(ManagerMax), ManagerCount, nil
}

func ComputeStore(engine *database.MysqlSession, UserData entity.MemberData) (int64, int64, error) {
	var StoreMax int64
	var StoreCurrent int64
	StoreMax = 1
	if time.Now().Before(UserData.UpgradeExpire) {
		if UserData.UpgradeLevel != 0 {
			//取出升級訂單
			order, err := product.GetUpgradeProductDataByLevel(engine, UserData.UpgradeLevel)
			if err != nil {
				return StoreMax, StoreCurrent, err
			}
			StoreMax = order.Store
		}
	}
	//計算收銀機是數量
	StoreCurrent = Store.CountStoreByUserId(engine, UserData.Uid)
	log.Debug("count", StoreMax, StoreCurrent)
	return StoreMax, StoreCurrent, nil
}

//取出未付款加值服務帳單
func HandleCountUnpaidUpgradeOrder(UserData entity.MemberData) (Response.CountUpgradeOrderResponse, error) {
	engine := database.GetMysqlEngine()
	defer engine.Close()
	var resp Response.CountUpgradeOrderResponse
	//計算未付款加值服務
	count, err := Orders.CountB2cUnpaidBillsByUserId(engine, UserData.Uid)
	if err != nil {
		return resp, fmt.Errorf("系統錯誤")
	}
	resp.UpgradeCount = count
	//取出轉帳加值服務訂單
	data, err := Orders.GetB2cOrderTransferByUserId(engine, UserData.Uid)
	if err != nil {
		return resp, fmt.Errorf("系統錯誤")
	}
	resp.BillingOrderId = data.OrderId
	return resp, nil
}


func HandleGetUnpaidUpgradeOrder(UserData entity.MemberData) ([]Response.UnpaidUpgradeOrderResponse, error) {
	engine := database.GetMysqlEngine()
	defer engine.Close()
	var resp []Response.UnpaidUpgradeOrderResponse
	data, err := Orders.GetB2cUnpaidOrdersByUserId(engine, UserData.Uid)
	if err != nil {
		return resp, fmt.Errorf("系統錯誤")
	}
	for _, v := range data {
		rsp := Response.UnpaidUpgradeOrderResponse{
			OrderId: v.OrderId,
			UserId: v.UserId,
			StoreId: v.StoreId,
			ProductId: v.ProductId,
			ProductName: v.ProductName,
			ProductDetail: v.ProductDetail,
			BillingTime: v.BillingTime,
			Amount: v.Amount,
			Payment: v.Payment,
			Status: v.OrderStatus,
			CreateTime: v.CreateTime.Format("2006/01/02"),
		}
		resp = append(resp, rsp)
	}
	return resp, nil
}
//產生B2C的帳單
func GeneratorB2cBill(engine *database.MysqlSession, UserData entity.MemberData, storeData entity.StoreData, UpgradeExpire time.Time) (entity.B2cBillingData, error) {
	start := UpgradeExpire.Add(time.Hour * time.Duration(24))
	expire := tools.NextMonth(start)
	var data entity.B2cBillingData
	//取出商品資料
	productData, err := product.GetUpgradeProductDataByLevel(engine, UserData.UpgradeLevel)
	if err != nil {
		return data, err
	}
	data.BillingId = tools.GeneratorB2CBillId()
	data.UserId = UserData.Uid
	data.StoreId = storeData.StoreId
	data.ProductId = productData.ProductId
	data.ProductName = productData.ProductName
	data.ProductDesc = productData.Description
	data.BillName = fmt.Sprintf("%s月份應付加值服務費", start.Format("01"))
	data.BillingTime = fmt.Sprintf("%s-%s", start.Format("2006/01/02"), expire.Format("2006/01/02"))
	data.BillingLevel = productData.UpgradeLevel
	data.Amount = productData.Amount
	data.BillingStatus = Enum.OrderWait
	data.Expiration = expire
	err = Orders.InsertB2cBillData(engine, data)
	if err != nil {
		return data, err
	}
	return data, nil
}
//產生B2C的訂單
func GeneratorB2cOrder(engine *database.MysqlSession, params entity.B2cOrderVo) (entity.B2cOrderData, error) {
	data := entity.B2cOrderData{}
	data.OrderId = tools.GeneratorB2COrderId()
	data.UserId = params.UserId
	data.StoreId = params.StoreId
	data.ProductId = params.ProductId
	data.ProductName = params.ProductName
	data.UpgradeLevel = params.UpgradeLevel
	data.OrderDetail = params.OrderDetail
	data.Amount = params.Amount
	data.OrderStatus = Enum.OrderInit
	data.Expiration = params.Expire
	data.BillingTime = params.BillingTime
	data.OrderSys = 3
	data.InvoiceType = params.InvoiceType
	data.CompanyBan = params.CompanyBan
	data.CompanyName = params.CompanyName
	data.DonateBan = params.DonateBan
	data.CarrierType = params.CarrierType
	data.CarrierId = params.CarrierId
	data.AskInvoice = false
	if err := Orders.InsertB2cOrderData(engine, data); err != nil {
		return data, err
	}
	return data, nil
}
//開啟收銀機
func OpenUpgradeService(engine *database.MysqlSession, UserId string, Level int64) error {
	//取出收銀機
	data, err := Store.GetUserAllStoreIdStore(engine, UserId)
	if err != nil {
		return err
	}
	for _, v := range data {
		if Level == 1 && v.StoreDefault != 1 {
			continue
		}
		//下架商品
		err := product.UpdateProductStatus(engine, v.StoreId, Enum.ProductStatusSuccess, Enum.ProductStatusPending)
		if err != nil {
			log.Error("Update Product Status Error", err)
		}
		//開收銀機
		v.StoreStatus = Enum.StoreStatusSuccess
		if err := Store.UpdateStoreData(engine, v.StoreId, v); err != nil {
			log.Error("Update Store Status Error", err)
		}
		//取出管理者
		storeData, err := Store.GetStoreManagerByStoreId(engine, v.StoreId)
		if err != nil {
			log.Error("Get Store Manager Error", err)
		}
		for _, v := range storeData {
			//開啟管理者
			if v.RankStatus != Enum.StoreRankDelete {
				v.RankStatus = Enum.StoreRankSuccess
				if err := Store.UpdateStoreManagerData(engine, v); err != nil {
					log.Debug("Update Store Manager Error")
				}
			}
		}
	}
	return nil
}

//關閉收銀機
func CloseUpgradeService(engine *database.MysqlSession, UserData entity.MemberData) error {
	//取出收銀機
	data, err := Store.GetUserAllStoreIdStore(engine, UserData.Uid)
	if err != nil {
		return err
	}
	for _, v := range data {
		if v.StoreDefault != 1 {
			//下架商品
			err := product.UpdateProductStatus(engine, v.StoreId, Enum.ProductStatusPending, Enum.ProductStatusSuccess)
			if err != nil {
				log.Error("Update Product Status Error", err)
			}
			//關閉收銀機
			v.StoreStatus = Enum.StoreStatusSuspend
			if err := Store.UpdateStoreData(engine, v.StoreId, v); err != nil {
				log.Error("Update Store Status Error", err)
			}
		}
		//取出管理者
		storeData, err := Store.GetStoreManagerByStoreId(engine, v.StoreId)
		if err != nil {
			log.Error("Get Store Manager Error", err)
		}
		for _, v := range storeData {
			//關閉管理者
			v.RankStatus = Enum.StoreRankSuspend
			if err := Store.UpdateStoreManagerData(engine, v); err != nil {
				log.Debug("Update Store Manager Error")
			}
		}
	}
	return nil
}

func MemberSuspendUpgradeService(engine *database.MysqlSession, UserData entity.MemberData) error {
	UserData.UpgradeType = Enum.UpgradeTypeSuspend
	UserData.UpdateTime = time.Now()
	if _, err := member.UpdateMember(engine, &UserData); err != nil {
		log.Error("Update member data Error!!")
		return err
	}
	return nil
}

func MemberDemoteLevelService(engine *database.MysqlSession, UserData entity.MemberData, Level int64) error {
	if err := member.UpdateMemberLevel(engine, UserData.Uid, time.Time{}, Level); err != nil {
		log.Error("Update member data Error!!")
		return err
	}
	data, err := Store.GetUserAllStoreIdStore(engine, UserData.Uid)
	if err != nil {
		return err
	}
	for _, v := range data {
		//非主收銀機
		if v.StoreDefault != 1 {
			//下架商品
			err := product.UpdateProductStatus(engine, v.StoreId, Enum.ProductStatusPending, Enum.ProductStatusSuccess)
			if err != nil {
				log.Error("Update Product Status Error", err)
			}
			//關閉收銀機
			v.StoreStatus = Enum.StoreStatusSuspend
			if err := Store.UpdateStoreData(engine, v.StoreId, v); err != nil {
				log.Error("Update Store Status Error", err)
			}
		}
		//取出管理者
		storeData, err := Store.GetStoreManagerByStoreId(engine, v.StoreId)
		if err != nil {
			log.Error("Get Store Manager Error", err)
		}
		for _, v := range storeData {
			//關閉管理者
			v.RankStatus = Enum.StoreRankDelete
			if err := Store.UpdateStoreManagerData(engine, v); err != nil {
				log.Debug("Update Store Manager Error")
			}
		}
	}
	return nil
}
//取得升級方案訂單結果
func HandleGetUpgradeOrder(orderId string, userData entity.MemberData) (Response.B2COrderResponse, error) {
	engine := database.GetMysqlEngine()
	defer engine.Close()
	var resp Response.B2COrderResponse
	data, err := Orders.GetB2cOrderByOrderId(engine, orderId)
	if err != nil {
		log.Error("Get Store Manager Error", err)
		return resp, fmt.Errorf("1001001")
	}
	if data.UserId != userData.Uid {
		return resp, fmt.Errorf("1001010")
	}
	if len(data.Payment) == 0 || data.OrderStatus == Enum.OrderInit {
		return resp, fmt.Errorf("1001010")
	}
	resp.OrderId = data.OrderId
	resp.OrderTime = data.CreateTime.Format("2006/01/02 15:04")
	resp.ExpireTime = ""
	if !data.Expiration.IsZero() {
		resp.ExpireTime = data.Expiration.Format("2006/01/02 15:04")
	}
	resp.PriceTotal = data.Amount
	resp.Payment = data.Payment
	//取得方案資訊
	result, err := product.GetUpgradeProductDataByLevel(engine, data.UpgradeLevel)
	if err != nil {
		return resp, fmt.Errorf("1001001")
	}

	resp.ProductName = fmt.Sprintf("每月 %v 元方案", result.Amount)
	resp.ProductDetail.Store = result.Store
	resp.ProductDetail.Manager = result.Manager
	resp.ProductDetail.ProductAmt = result.Amount

	var detail []Response.B2cOrderDetail
	if err := json.Unmarshal([]byte(data.OrderDetail), &detail); err != nil {
		return resp, fmt.Errorf("1001001")
	}

	var unpaidOrder []Response.OrderList
	var upgradeList []Response.UpgradeList
	var orderTotal int64
	var upgradeTotal int64
	for _, v := range detail {
		if v.ProductType == "UPGRADE" {
			var sign bool
			if len(v.ProductId) == 0 {
				sign = false
				upgradeTotal -= v.ProductAmount
			} else {
				sign = true
				upgradeTotal += v.ProductAmount
			}
			OrderList := Response.UpgradeList{
				UpgradeText:  v.ProductName,
				UpgradePrice: v.ProductAmount,
				SignType: sign,
			}
			upgradeList = append(upgradeList, OrderList)
		}
		if v.ProductType == "BILLING" {
			OrderList := Response.OrderList{
				OrderText:  v.ProductName,
				OrderPrice: v.ProductAmount,
			}
			orderTotal += v.ProductAmount
			unpaidOrder = append(unpaidOrder, OrderList)
		}
	}
	resp.UpgradeSum = upgradeTotal
	resp.OrderSum = orderTotal
	resp.UpgradeList = upgradeList
	resp.OrderList = unpaidOrder
	if data.Payment == Enum.Transfer {
		transferData, err := transfer.GetTransferByOrderId(engine, data.OrderId)
		if err != nil {
			return resp, fmt.Errorf("1001001")
		}
		resp.BankCode = transferData.BankCode + " " + transferData.BankName
		resp.BankAccount = transferData.BankAccount
		resp.AtmExpire = transferData.ExpireDate.Format("2006/01/02 15:04")
	}
	return resp, nil
}

func B2cOrderMigration() error {
	engine := database.GetMysqlEngine()
	defer engine.Close()
	order, err := Orders.GetB2cBillOrders(engine)
	if err != nil {
		log.Error("Get B2c Bill order Error", err)
	}
	for _, v := range order {
		var detail []entity.B2cOrderDetail
		if err := tools.JsonDecode([]byte(v.OrderDetail), &detail); err != nil {
			log.Error("Json Decode Error", err)
		}
		if len(detail) == 0 {
			log.Error("ssss", v.OrderId)
			continue
		}
		productData, err := product.GetUpgradeProductByProductId(engine, detail[0].ProductId)
		if err != nil {
			return fmt.Errorf("1001001")
		}
		var data entity.B2cBillingData
		data.BillingId = tools.GeneratorB2CBillId()
		data.UserId = v.UserId
		data.StoreId = v.StoreId
		data.ProductId = detail[0].ProductId
		data.ProductName = productData.ProductName
		data.ProductDesc = detail[0].ProductDetail
		data.BillName = detail[0].ProductName
		data.BillingTime = v.BillingTime
		data.BillingLevel = productData.UpgradeLevel
		data.Amount = detail[0].ProductAmount
		data.BillingStatus = v.OrderStatus
		data.OrderId = v.UpgradeOrder
		data.Expiration = v.Expiration
		data.ServiceType = "Upgrade"
		err = Orders.InsertB2cBillData(engine, data)
		if err != nil {
			return err
		}
	}
	return nil
}

func GetChangeProduct(productId int64) string {
	switch productId {
		case 1:
			return "每月99元方案"
		case 2:
			return "每月199元方案"
		case 3:
			return "每月299元方案"
		default:
			return "每月99元方案"
	}
}


