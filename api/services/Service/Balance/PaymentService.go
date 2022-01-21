package Balance

import (
	"api/services/Enum"
	"api/services/Service/MemberService"
	"api/services/Service/Notification"
	"api/services/Service/OrderService"
	"api/services/Service/TransferService"
	"api/services/Service/Upgrade"
	"api/services/VO/OrderVo"
	"api/services/VO/Request"
	"api/services/VO/Response"
	"api/services/dao/Orders"
	"api/services/dao/Promotion"
	"api/services/dao/Store"
	"api/services/dao/member"
	"api/services/database"
	"api/services/entity"
	"api/services/util/log"
	"api/services/util/tools"
	"encoding/json"
	"fmt"
	"github.com/spf13/viper"
	"strings"
	"time"
)

//付款完成
func TransferC2CCheckout(engine *database.MysqlSession, orderId string, Status string, now time.Time) error {
	OrderData, err := Orders.GetOrderByOrderId(engine, orderId)
	if err != nil {
		log.Error("Get Order Data Error", err)
		return err
	}
	//判斷付款完成 就別再執行
	if OrderData.OrderStatus == Enum.OrderSuccess {
		return err
	}
	if OrderData.PayWay != Enum.CvsPay && Status == Enum.OrderSuccess || Status == Enum.OrderAudit {
		OrderData.PayWayTime = now
	}
	if Status == Enum.OrderFail {
		//付款失敗 退回庫存
		if err := OrderService.ProcessReturnStock(engine, orderId); err != nil {
			log.Error("Process Return Stock Error", err)
		}
	}
	if err := OrderService.ChangeOrderStatus(engine, OrderData, Status); err != nil {
		return err
	}
	if Status == Enum.OrderSuccess || Status == Enum.OrderAudit {
		//將收款金額寫入保留款項
		if err := RetainDeposit(OrderData.SellerId, OrderData.OrderId, OrderData.TotalAmount, Enum.BalanceTypeDeposit, "訂單交易存入"); err != nil {
			return err
		}
		//發送訂購成功訊息
		if err := Notification.SendSuccessfullyOrderedMessage(engine, OrderData.OrderId); err != nil {
			log.Error("Send Success Message Error", err)
			return err
		}
	}
	return nil
}

//付款完成
func OrderCheckout(engine *database.MysqlSession, orderId string, Status string) error {
	OrderData, err := Orders.GetOrderByOrderId(engine, orderId)
	if err != nil {
		log.Error("Get Order Data Error", err)
		return err
	}
	//判斷付款完成 就別再執行
	if OrderData.OrderStatus == Enum.OrderSuccess {
		return err
	}
	if OrderData.PayWay != Enum.CvsPay && Status == Enum.OrderSuccess || Status == Enum.OrderAudit {
		OrderData.PayWayTime = time.Now()
	}
	if Status == Enum.OrderFail {
		//付款失敗 退回庫存
		if err := OrderService.ProcessReturnStock(engine, orderId); err != nil {
			log.Error("Process Return Stock Error", err)
		}
		//付款失敗 退回優惠卷
		if err := HandleDeleteCouponUsedRecord(engine, OrderData); err != nil {
			log.Error("Process Delete Coupon Used Record Error", err)
		}
	}
	if err := OrderService.ChangeOrderStatus(engine, OrderData, Status); err != nil {
		return err
	}
	if Status == Enum.OrderSuccess || Status == Enum.OrderAudit {
		//將收款金額寫入保留款項
		if err := RetainDeposit(OrderData.SellerId, OrderData.OrderId, OrderData.TotalAmount, Enum.BalanceTypeDeposit, "訂單交易存入"); err != nil {
			return err
		}
		//發送訂購成功訊息
		if err := Notification.SendSuccessfullyOrderedMessage(engine, OrderData.OrderId); err != nil {
			log.Error("Send Success Message Error", err)
			return err
		}
	}
	return nil
}
//B2C付款完成
func B2cOrderCheckout(engine *database.MysqlSession, orderId, status, payment string) error {
	OrderData, err := Orders.GetB2cOrderByOrderId(engine, orderId)
	if err != nil {
		log.Error("Get B2c Order Data Error", err)
		return err
	}
	OrderData.OrderStatus = status
	if len(payment) != 0 {
		OrderData.Payment = payment
	}
	if status == Enum.OrderSuccess {
		OrderData.AskInvoice = true
		OrderData.InvoiceStatus = Enum.InvoiceOpenStatusNot
	}
	if err := Orders.UpdateB2cOrderData(engine, OrderData); err != nil {
		log.Error("Update B2c Order Data Error", err)
		return err
	}
	if status == Enum.OrderSuccess {
		//取出帳單
		var person []entity.B2cOrderDetail
		_ = json.Unmarshal([]byte(OrderData.OrderDetail), &person)
		for _, v := range person {
			if v.ProductType == Enum.B2cOrderTypeBilling {
				if err := B2cBillChangeStatus(engine, OrderData.OrderId, v.ProductId); err != nil {
					log.Error("Update B2c Order Detail Status Error", err)
					return err
				}
			}
		}
		if err := UpdateMemberUpgradeStatus(engine, OrderData.UserId, OrderData.Expiration, OrderData.UpgradeLevel); err != nil {
			return err
		}
		if err := Upgrade.OpenUpgradeService(engine, OrderData.UserId, OrderData.UpgradeLevel); err != nil {
			log.Error("Open Upgrade data Error!!")
			return err
		}
		//發送訊息
		if err := Notification.SendUpgradeApplyMessage(engine, OrderData.StoreId, OrderData.UserId, OrderData.OrderId); err != nil {
			log.Error("Send Upgrade Message Error", err)
		}
	}
	return nil
}
//更新會員升級方案的狀態
func UpdateMemberUpgradeStatus(engine *database.MysqlSession, userId string, expire time.Time, level int64) error {
	if expire.IsZero() {
		expire = tools.NextMonth(time.Now())
	}
	if err := member.UpdateMemberLevel(engine, userId, expire, level); err != nil {
		log.Error("Update member data Error!!")
		return err
	}
	return nil
}

func OrderBillCheckout(engine *database.MysqlSession, orderId string, Status string) error {
	OrderData, err := Orders.GetBillOrderByOrderId(engine, orderId)
	if err != nil {
		log.Error("Get Order Data Error", err)
		return err
	}
	OrderData.PayWayTime = time.Now()
	OrderData.PayWayStatus = Status
	if err := Orders.UpdateBillOrderData(engine, OrderData); err != nil {
		return err
	}
	return nil
}
//更新帳單付款完成
func B2cBillChangeStatus(engine *database.MysqlSession, orderId, billId string) error {
	data, err := Orders.GetB2cBillByOrderId(engine, billId)
	if err != nil {
		log.Error("Get B2c Order Data Error", err)
		return err
	}
	data.OrderId = orderId
	data.BillingStatus = Enum.OrderSuccess
	if err := Orders.UpdateB2cBillData(engine, data); err != nil {
		log.Error("Update B2c Bill Data Error", err)
		return err
	}
	return nil
}
//C2C餘額付款
func C2cBalanceTransaction(engine *database.MysqlSession, orderData entity.OrderData) error {
	balance := GetBalanceByUid(engine, orderData.BuyerId)
	if orderData.TotalAmount > balance {
		return fmt.Errorf("1005008")
	}
	storeData, err := Store.GetStoreDataByStoreId(engine, orderData.StoreId)
	if err != nil {
		log.Error("Get Store Error", err)
		return fmt.Errorf("1005005")
	}
	//買家扣除餘額
	comment := fmt.Sprintf("%s<br>%s", orderData.OrderId, storeData.StoreName) //賣場名稱 + 訂單號碼
	if err := Withdrawal(orderData.BuyerId, orderData.OrderId, orderData.TotalAmount, Enum.BalanceTypeBalancePay, comment); err != nil {
		log.Error("debit Error", err)
		return fmt.Errorf("1005005")
	}
	if err := OrderCheckout(engine, orderData.OrderId, Enum.OrderSuccess); err != nil {
		log.Error("Order change Status Error", err)
		return fmt.Errorf("1005005")
	}
	return nil
}
//B2C餘額付款
func B2cBalanceTransaction(engine *database.MysqlSession, orderData entity.B2cOrderData, comment string) error {
	balance := GetBalanceByUid(engine, orderData.UserId)
	if float64(orderData.Amount) > balance {
		return fmt.Errorf("1005008")
	}
	//出帳
	if err := Withdrawal(orderData.UserId, orderData.OrderId, float64(orderData.Amount), Enum.BalanceTypeService, comment); err != nil {
		log.Error("debit Error", err)
		return fmt.Errorf("1005005")
	}
	SellerId := viper.GetString("PLATFORM.USERID")
	if err := Deposit(SellerId, orderData.OrderId, float64(orderData.Amount), Enum.BalanceTypePayment, orderData.UserId); err != nil {
		log.Error("Deposit Error", err)
		return fmt.Errorf("1005005")
	}
	if err := B2cOrderCheckout(engine, orderData.OrderId, Enum.OrderSuccess, Enum.Balance); err != nil {
		log.Error("Order change Status Error", err)
		return fmt.Errorf("1005005")
	}
	return nil
}
//B2C付款
func ProcessPayment(engine *database.MysqlSession, order entity.B2cOrderData, params Request.B2CPayRequest) (Response.B2cPayResponse, error) {
	var resp Response.B2cPayResponse
	switch strings.ToUpper(order.Payment) {
	// todo 台灣PAY
	case Enum.TaiwanPay:
		return resp, nil
	case Enum.Transfer: //轉帳
		//建立轉帳資料
		if _, err := TransferService.CreateTransfer(engine, order.OrderId, order.Amount, Enum.OrderTransB2c); err != nil {
			log.Debug("Create transfer", err)
			return resp, fmt.Errorf("1005004")
		}
		resp.OrderId = order.OrderId
		return resp, nil
	case Enum.Balance: //餘額付款
		//餘額付款
		comment := ""
		if len(order.BillingTime) > 0 {
			comment = fmt.Sprintf("%s", order.BillingTime)
		} else {
			now := time.Now()
			expire := tools.NextMonth(now)
			comment = fmt.Sprintf("%s ~ %s", now.Format("2006/01/02"), expire.Format("2006/01/02"))
		}
		//扣除餘額
		if err := B2cBalanceTransaction(engine, order, comment); err != nil {
			log.Debug("Balance trans Error", err)
			return resp, fmt.Errorf("1005016")
		}
		resp.OrderId = order.OrderId
		return resp, nil
	case Enum.Credit: //信用卡
		//取信用卡資料 或新增信用卡資料
		cardData, err := MemberService.GetCreditData(engine, params.CardId, params.CardNumber, params.CardExpiration, order.UserId)
		if err != nil {
			log.Error("Get Card data Error", err)
			return resp, fmt.Errorf("1005006")
		}
		// todo 付款完成後的處理
		vo := OrderVo.CreditPaymentVo {
			OrderId: order.OrderId,
			TotalAmount: float64(order.Amount),
			OrderType: Enum.OrderTransB2c,
			AuditStatus: Enum.OrderAuditRelease,
			BuyerId: order.UserId,
		}
		param := params.GetB2cCreditPayment()
		//信用卡取授權
		res, err := HandleAuth(vo, &param, cardData)
		if err != nil {
			log.Error("Auth handle Error", err)
			return resp, fmt.Errorf("1005006")
		}
		if len(res.RtnURL) != 0 {
			resp.RtnURL = res.RtnURL
		}
		resp.OrderId = order.OrderId
		return resp, nil
	}
	return resp, fmt.Errorf("1005002")
}

func UpgradeAutoPayment(engine *database.MysqlSession, data entity.B2cBillingData) error {
	balance := GetBalanceByUid(engine, data.UserId)
	//檢查餘額是否足夠
	if data.Amount <= int64(balance) {
		//產生DETAIL
		var details []entity.B2cOrderDetail
		var detail entity.B2cOrderDetail
		detail.ProductName = data.BillName
		detail.ProductId = data.BillingId
		detail.ProductAmount = data.Amount
		detail.ProductType = "BILLING"
		details = append(details, detail)
		orderDetail, _ := tools.JsonEncode(details)
		//fixme 少發票資料
		last, err := Orders.GetUpgradeProductDataByOrder(engine, data.UserId)
		if err != nil {
			log.Error("Get Upgrade Order Error", err)
			return err
		}
		//產生訂單
		orderVo := entity.B2cOrderVo{
			ProductId: data.ProductId,
			ProductName: data.ProductName,
			UserId: data.UserId,
			StoreId: data.StoreId,
			OrderDetail: orderDetail,
			BillingTime: data.BillingTime,
			UpgradeLevel: data.BillingLevel,
			Amount: data.Amount,
			Expire: data.Expiration,
			InvoiceType: last.InvoiceType,
			CompanyBan: last.CompanyBan,
			CompanyName: last.CompanyName,
			DonateBan: last.DonateBan,
			CarrierType: last.CarrierType,
			CarrierId: last.CarrierId,
		}
		order, err := Upgrade.GeneratorB2cOrder(engine, orderVo)
		if err != nil {
			log.Error("Generator Order Error", err)
			return err
		}
		comment := fmt.Sprintf("%s", data.BillingTime)
		if err := B2cBalanceTransaction(engine, order, comment); err != nil {
			log.Error("Balance trans Error", err)
			return err
		}
	} else {
		//發通知付款
		if err := Notification.SendUpgradePaymentMessage(engine, data.StoreId, data.UserId, data.BillingId); err != nil {
			return err
		}
	}
	return nil
}

func HandleDeleteCouponUsedRecord(engine *database.MysqlSession, order entity.OrderData) error {
	if len(order.CouponNumber) == 0 {
		return nil
	}
	data, err := Promotion.GetPromotionCodeByCodeAndStoreId(engine, order.StoreId, order.CouponNumber)
	if err != nil {
		log.Debug("Get Promotion Code Error", err)
		return err
	}
	data.PromotionCode.IsUsed = false
	if err := Promotion.UpdatePromotionCode(engine, data.PromotionCode); err != nil {
		log.Debug("Update Promotion Code Error", err)
		return err
	}
	record, err := Promotion.GetUsedCouponRecordByOrderId(engine, order.OrderId)
	if err != nil {
		log.Error("Insert Used Coupon Record Error", err)
		return err
	}
	record.RecordStatus = Enum.RecordStatusDelete
	if err := Promotion.UpdateUsedCouponRecord(engine, record); err != nil {
		log.Error("Update Used Coupon Record Error", err)
		return err
	}
	return nil
}