package model

import (
	"api/services/Enum"
	"api/services/Service/Balance"
	"api/services/Service/MemberService"
	"api/services/Service/Notification"
	"api/services/Service/OrderService"
	"api/services/Service/Product"
	"api/services/Service/TransferService"
	"api/services/Service/UserAddressService"
	"api/services/VO/OrderVo"
	"api/services/VO/Request"
	"api/services/VO/Response"
	"api/services/dao/Orders"
	"api/services/dao/member"
	"api/services/database"
	"api/services/entity"
	"api/services/util/images"
	"api/services/util/log"
	"api/services/util/qrcode"
	"api/services/util/tools"
	"api/services/util/upload"
	"fmt"
	"strings"
	"time"
)

//買家預覽帳單內容
func HandleGetReviewBillOrder(UserData entity.MemberData, OrderId string) (Response.BillResponse, error) {
	engine := database.GetMysqlEngine()
	defer engine.Close()
	var resp Response.BillResponse
	data, err := Orders.GetBillOrderByOrderId(engine, OrderId)
	if err != nil {
		log.Error("Ger Bill Order Error", err)
		return resp, fmt.Errorf("1001001")
	}
	if len(data.BillId) == 0 {
		return resp, fmt.Errorf("1004001")
	}
	if data.BuyerId != UserData.Uid {
		return resp, fmt.Errorf("1001001")
	}
	resp = data.GetBillToResponse()
	resp.Qrcode = qrcode.GetQrcodeImageLink(resp.BillId)
	resp.ReceiverAddress = GetShipAddress(engine, data.ShipType, data.ReceiverAddress)
	buyer, err := member.GetMemberDataByUid(engine, data.BuyerId)
	if err != nil {
		log.Error("Ger Buyer Data Error", err)
		return resp, fmt.Errorf("1001001")
	}
	resp.MemberInfo = buyer.GetMemberLoginInfo()
	return resp, nil
}
//賣家查看帳單內容
func HandleGetBillOrder(OrderId string) (Response.BillResponse, error) {
	engine := database.GetMysqlEngine()
	defer engine.Close()
	var resp Response.BillResponse
	data, err := Orders.GetBillOrderByOrderId(engine, OrderId)
	if err != nil {
		log.Error("Ger Bill Order Error", err)
		return resp, fmt.Errorf("1001001")
	}
	if len(data.BillId) == 0 {
		return resp, fmt.Errorf("1004001")
	}
	resp = data.GetBillToResponse()
	resp.BuyerName = tools.MaskerName(resp.BuyerName)
	resp.ReceiverName = tools.MaskerName(resp.ReceiverName)
	resp.Qrcode = qrcode.GetQrcodeImageLink(resp.BillId)
	resp.ReceiverAddress = GetShipAddress(engine, data.ShipType, data.ReceiverAddress)
	if !tools.InArray([]string{Enum.I_POST, Enum.CVS_FAMILY, Enum.CVS_HI_LIFE, Enum.CVS_OK_MART, Enum.CVS_7_ELEVEN}, resp.ShipType) {
		resp.ReceiverAddress = tools.MaskerAddress(resp.ReceiverAddress)
	}
	buyer, err := member.GetMemberDataByUid(engine, data.BuyerId)
	if err != nil {
		log.Error("Ger Buyer Data Error", err)
		return resp, fmt.Errorf("1001001")
	}
	resp.MemberInfo = buyer.GetMemberLoginInfo()
	return resp, nil
}
//賣家確認接受帳單
func HandleBillConfirm(storeData entity.StoreDataResp, params Request.BillConfirmRequest) (Response.BillConfirmResponse, error) {
	engine := database.GetMysqlEngine()
	defer engine.Close()
	var resp Response.BillConfirmResponse
	bill, err := Orders.GetBillOrderByOrderId(engine, params.BillId)
	if err != nil {
		log.Error("Ger Bill Order Error", err)
		return resp, fmt.Errorf("1001001")
	}
	if bill.BuyerId == storeData.SellerId {
		return resp, fmt.Errorf("1007006")
	}
	if bill.BillStatus == Enum.BillStatusOverdue {
		return resp, fmt.Errorf("1007002")
	}
	if bill.BillStatus == Enum.BillStatusClose {
		return resp, fmt.Errorf("1007003")
	}
	//建立訂單
	order, err := CreateBillOrder(engine, storeData,  bill)
	if err != nil {
		log.Error("Ger Bill Order Error", err)
		return resp, fmt.Errorf("1001001")
	}
	//付款
	if err := processBillConfirmPayment(engine, order); err != nil {
		return resp, err
	}
	resp.OrderId = order.OrderId
	return resp, nil
}
//付款處理
func processBillConfirmPayment(engine *database.MysqlSession, orderData entity.OrderData) error {
	switch strings.ToUpper(orderData.PayWay) {
	case Enum.Transfer: //轉帳
		_, err := TransferService.CreateTransfer(engine, orderData.OrderId, int64(orderData.TotalAmount), Enum.OrderTransC2c)
		if err != nil {
			log.Error("Create transfer", err)
			return fmt.Errorf("1005004")
		}
		if err := OrderService.OrderWaitPayment(engine, orderData, Enum.OrderWait); err != nil {
			log.Error("Order change Status Error", err)
			return fmt.Errorf("1005004")
		}
		//發送訊息
		if err := Notification.SendBillOrderAtmMessage(engine, orderData); err != nil {
			log.Error("Send Atm Successfully Ordered Message Error", err)
			return fmt.Errorf("1005004")
		}
	case Enum.Balance: //餘額付款
		if err := Balance.OrderCheckout(engine, orderData.OrderId, Enum.OrderSuccess); err != nil {
			log.Error("Order change Status Error", err)
			return fmt.Errorf("1005005")
		}
		//發送訊息
		if err := Notification.SendBillOrderSuccessMessage(engine, orderData); err != nil {
			log.Error("Send Atm Successfully Ordered Message Error", err)
			return fmt.Errorf("1005004")
		}
	case Enum.Credit: //信用卡
		//信用卡改為可請款 fixme
		if err := Balance.OrderCheckout(engine, orderData.OrderId, Enum.OrderSuccess); err != nil {
			log.Error("Order change Status Error", err)
			return fmt.Errorf("1005006")
		}
		//發送訊息
		if err := Notification.SendBillOrderSuccessMessage(engine, orderData); err != nil {
			log.Error("Send Atm Successfully Ordered Message Error", err)
			return fmt.Errorf("1005004")
		}
		return nil
	case Enum.CvsPay:
		if err := Balance.OrderCheckout(engine, orderData.OrderId, Enum.OrderSuccess); err != nil {
			log.Error("超商取貨付款 Order Checkout Error", err)
			return fmt.Errorf("1005007")
		}
		if err := Notification.SendBillOrderSuccessMessage(engine, orderData); err != nil {
			log.Error("Send Atm Successfully Ordered Message Error", err)
			return fmt.Errorf("1005004")
		}
	}
	//結束帳單
	if err:= ChangeBillOrderStatus(engine, orderData.OrderId, Enum.BillStatusClose); err != nil {
		log.Error("Change Bill Order status Error", err)
		return fmt.Errorf("1005002")
	}
	return nil
}
//建立反向帳單
func HandleNewBillOrder(host string, userData entity.MemberData, params entity.BillOrderParams) (Response.PayResponse, error) {
	engine := database.GetMysqlEngine()
	defer engine.Close()
	var resp Response.PayResponse
	if params.ProductPrice == 0 {
		return resp, fmt.Errorf("1007008")
	}
	//圖片上傳
	filename := ""
	if len(params.ProductImage) != 0 {
		file, err := processBillProductImage(params.ProductImage)
		if err != nil {
			log.Error("Upload Product image", err)
			return resp, fmt.Errorf("1007001")
		}
		filename = file
	}
	if params.ProductQty == 0 {
		log.Error("Product Qty Set 0")
		return resp, fmt.Errorf("1007007")
	}
	//產生反向訂單
	order, err := generateOrder(engine, userData, params, filename)
	if err != nil {
		return resp, err
	}
	//產生QRCODE
	var url = fmt.Sprintf("https://%s/bill/%s", host, order.BillId)

	var imgPath = fmt.Sprintf("./www%s", userData.Picture)
	img, err := images.GetImageFromFilePath(imgPath)
	if err != nil {
		log.Error("Get image Error", err, imgPath)
	}
	if err := qrcode.GeneratorQrCode(url, order.BillId, img); err != nil {
		log.Error("generator qr code err ", err)
	}
	//信用卡取授權、扣除餘額 其他不處理
	res, err := processBillPayment(engine, userData, params, order)
	if err != nil {
		return resp, err
	}
	resp.OrderId = order.BillId
	if res.Status == Enum.OrderWait {
		resp.RtnURL = res.RtnHtml
		return resp, nil
	}
	return resp, nil
}
//圖片上傳
func processBillProductImage(Image string) (string, error) {
	filename := ""
	if len(Image) != 0 {
		if strings.Index(Image, ".") < 0 {
			name, err := upload.ProductImage(Image)
			if err != nil {
				return filename, err
			}
			filename = name
		}
	}
	return filename, nil
}
//產生反向訂單
func generateOrder(engine *database.MysqlSession, userData entity.MemberData, params entity.BillOrderParams, filename string) (entity.BillOrderData, error) {
	var data entity.BillOrderData
	if len(params.ReceiverId) != 0 {
		params.Address = UserAddressService.GetReceiverAddress(engine, params.ReceiverId)
	}
	if (strings.ToLower(params.ShipType) == Enum.DELIVERY_POST ||  strings.ToLower(params.ShipType) == Enum.DELIVERY_E_CAN ||
		strings.ToLower(params.ShipType) == Enum.DELIVERY_T_CAT || strings.ToLower(params.ShipType) == Enum.DELIVERY_OTHER ) &&
		len(strings.Split(params.Address, ",")) != 4 {
		return data, fmt.Errorf("1005003")
	}
	data = params.GeneratorBillOrderData(userData)
	data.ProductImage = filename
	//產生平台費用
	vo := data.GetBillOrder()
	fee := Balance.CalculatePlatFormFee(vo)
	data.PlatformShipFee = fee.PlatformShipFee
	data.PlatformTransFee = fee.PlatformTransFee
	data.PlatformPayFee = fee.PlatformPayFee
	data.PlatformInfoFee = fee.PlatformInfoFee
	data.CaptureAmount = fee.CaptureAmount
	//產生短網址
	url := fmt.Sprintf("/bill/%s", data.BillId)
	tiny, err := Product.GeneratorShortUrl(engine, url)
	if err != nil {
		return data, fmt.Errorf("1001001")
	}
	data.TinyUrl = tiny
	if len(params.ReceiverId) == 0 {
		if err := UserAddressService.NewShipReceiverAddress(engine, userData, params.ShipType, params.BuyerName, params.BuyerPhone, params.ReceiverName, params.ReceiverPhone, params.Address);
			err != nil {
			return data, fmt.Errorf("1001001")
		}
	}
	if err := Orders.InsertBillOrderData(engine, data); err != nil {
		log.Error("Insert Bill Order Error", err)
		return data, fmt.Errorf("1001001")
	}
	return data, nil
}
//處理付款
func processBillPayment(engine *database.MysqlSession, userData entity.MemberData, params entity.BillOrderParams, data entity.BillOrderData) (Response.PaymentResponse, error) {
	var resp Response.PaymentResponse
	switch strings.ToUpper(data.PayWayType) {
		case Enum.Credit:
			cardData, err := MemberService.GetCreditData(engine, params.CardId, params.CardNumber, params.CardExpiration, userData.Uid)
			if err != nil {
				log.Error("Get Card data Error", err)
				return resp, fmt.Errorf("1005002")
			}
			vo := OrderVo.CreditPaymentVo {
				OrderId: data.BillId,
				TotalAmount: data.TotalAmount,
				OrderType: Enum.OrderTransBill,
				BuyerId: userData.Uid,
				AuditStatus: Enum.OrderAuditInit,
			}
			param := params.GetCreditPayment()
			res, err := Balance.HandleAuth(vo, &param, cardData)
			if err != nil {
				log.Error("Auth handle Error", err)
				return resp, fmt.Errorf("1005002")
			}
			resp.Status = res.Status
			resp.Message = "你已透過Check'Ne成功開立帳單。"
			resp.RtnHtml = res.RtnURL
			return resp, nil
		case Enum.Balance:
			balance := Balance.GetBalanceByUid(engine, userData.Uid)
			if data.TotalAmount > balance {
				return resp, fmt.Errorf("1005008")
			}
			comment := fmt.Sprintf("%s", data.BillId) //賣場名稱 + 訂單號碼
			if err := Balance.Withdrawal(userData.Uid, data.BillId, data.TotalAmount, Enum.BalanceTypeBalancePay, comment); err != nil {
				log.Error("debit Error", err)
				return resp, fmt.Errorf("1005005")
			}
			if err := Balance.OrderBillCheckout(engine, data.BillId, Enum.OrderSuccess); err != nil {
				return resp, fmt.Errorf("1005005")
			}
		default:
			resp.Status = Enum.OrderSuccess
			resp.Message = "你已透過Check'Ne成功開立帳單。"
			return resp, nil
	}
	return resp, nil
}
//變更反向帳單狀態
func ChangeBillOrderStatus(engine *database.MysqlSession, billId, status string) error {
	data, err := Orders.GetBillOrderByOrderId(engine, billId)
	if err != nil {
		log.Error("Ger Bill Order Error", err)
		return err
	}
	data.BillStatus = status
	data.UpdateTime = time.Now()
	if err := Orders.UpdateBillOrderData(engine, data); err != nil{
		log.Error("Update Bill Order Error", err)
		return err
	}
	return nil
}
//反向帳單列表
func HandleBillList(userData entity.MemberData, params Request.BillListRequest) (Response.BillListResponse, error) {
	engine := database.GetMysqlEngine()
	defer engine.Close()
	var resp Response.BillListResponse
	Tab := ""
	switch params.Tab {
		case "Valid":
			Tab = Enum.BillStatusInit
		case "Cancel":
			Tab = Enum.BillStatusCancel
		case "Overdue":
			Tab = Enum.BillStatusOverdue
		default:
			Tab = Enum.BillStatusInit
	}
	count, err := Orders.CountBillOrderByUserId(engine, userData.Uid, Tab)
	if err != nil {
		log.Error("Count Invoice List Error", err)
		return resp, fmt.Errorf("1001001")
	}
	data, err := Orders.GetBillOrderListByUserId(engine, userData.Uid, Tab, int(params.Limit), int(params.Start))
	if err != nil {
		log.Error("Get Invoice List Error", err)
		return resp, fmt.Errorf("1001001")
	}
	validCount, err := Orders.CountBillOrderByUserId(engine, userData.Uid, Enum.BillStatusInit)
	if err != nil {
		log.Error("Count Invoice List Error", err)
		return resp, fmt.Errorf("1001001")
	}
	cancelCount, err := Orders.CountBillOrderByUserId(engine, userData.Uid, Enum.BillStatusCancel)
	if err != nil {
		log.Error("Count Invoice List Error", err)
		return resp, fmt.Errorf("1001001")
	}
	overdueCount, err := Orders.CountBillOrderByUserId(engine, userData.Uid, Enum.BillStatusOverdue)
	if err != nil {
		log.Error("Count Invoice List Error", err)
		return resp, fmt.Errorf("1001001")
	}
	for _, v := range data {
		var res Response.BillResponse
		res = v.GetBillToResponse()
		res.Qrcode = qrcode.GetQrcodeImageLink(v.BillId)
		res.ReceiverAddress = GetShipAddress(engine, v.ShipType, v.ReceiverAddress)
		resp.BillList = append(resp.BillList, res)
	}
	resp.Count = count
	resp.Tabs.Valid = validCount
	resp.Tabs.Cancel = cancelCount
	resp.Tabs.Overdue = overdueCount
	return resp, nil
}
//帳單延期24小時
func HandleBillExtension(userData entity.MemberData, params Request.BillConfirmRequest) error {
	engine := database.GetMysqlEngine()
	defer engine.Close()
	data, err := Orders.GetBillOrderByOrderIdAndUserId(engine, params.BillId, userData.Uid)
	if err != nil {
		log.Error("Ger Bill Order Error", err)
		return fmt.Errorf("1001001")
	}
	if len(data.BillId) == 0 {
		return fmt.Errorf("1007004")
	}
	if data.IsExtension == true {
		return fmt.Errorf("1007005")
	}
	data.IsExtension = true
	data.BillExpire = data.BillExpire.Add(time.Hour * time.Duration(24))
	if err := Orders.UpdateBillOrderData(engine, data); err != nil{
		log.Error("Update Bill Order Error", err)
		return fmt.Errorf("1001001")
	}
	return nil
}
//買家帳單取消
func HandleBillCancel(userData entity.MemberData, params Request.BillConfirmRequest) error {
	engine := database.GetMysqlEngine()
	defer engine.Close()
	data, err := Orders.GetBillOrderByOrderIdAndUserId(engine, params.BillId, userData.Uid)
	if err != nil {
		log.Error("Ger Bill Order Error", err)
		return fmt.Errorf("1001001")
	}
	if len(data.BillId) == 0 {
		return fmt.Errorf("1007004")
	}
	if data.PayWayType == Enum.Balance {
		comment := fmt.Sprintf("%s", data.BillId) //訂購單編號
		err = Balance.Deposit(data.BuyerId, data.BillId, data.TotalAmount, Enum.BalanceTypeBillFail, comment)
		if err != nil {
			log.Error("Balance Deposit Error", err)
		}
	}
	if data.PayWayType == Enum.Credit {
		var vo entity.CancelRequest
		vo.OrderId = data.BillId
		if err := Balance.VoidProcess(engine, vo); err != nil {
			return fmt.Errorf("退款失敗！")
		}
	}
	data.BillStatus = Enum.BillStatusCancel
	if err := Orders.UpdateBillOrderData(engine, data); err != nil{
		log.Error("Update Bill Order Error", err)
		return fmt.Errorf("1001001")
	}
	return nil
}
//反向帳單列表
func HandleAllBillList(userData entity.MemberData, params Request.BuyerBillListRequest) (Response.BuyerBillListResponse, error) {
	engine := database.GetMysqlEngine()
	defer engine.Close()
	var resp Response.BuyerBillListResponse
	count, err := Orders.CountAllBillOrderByUserId(engine, userData.Uid)
	if err != nil {
		log.Error("Count Invoice List Error", err)
		return resp, fmt.Errorf("1001001")
	}
	data, err := Orders.GetAllBillListByUserId(engine, userData.Uid, int(params.Limit), int(params.Start))
	if err != nil {
		log.Error("Get Invoice List Error", err)
		return resp, fmt.Errorf("1001001")
	}
	for _, v := range data {
		var res Response.BillResponse
		res = v.GetBillToResponse()
		res.Qrcode = qrcode.GetQrcodeImageLink(v.BillId)
		res.ReceiverAddress = GetShipAddress(engine, v.ShipType, v.ReceiverAddress)
		resp.BillList = append(resp.BillList, res)
	}
	resp.Count = count
	return resp, nil
}