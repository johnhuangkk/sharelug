package model

import (
	"api/services/Enum"
	"api/services/Service/Balance"
	"api/services/Service/Erp"
	"api/services/Service/Excel"
	"api/services/Service/Invoice"
	"api/services/Service/MemberService"
	"api/services/Service/Notification"
	"api/services/Service/OrderService"
	"api/services/Service/Sms"
	"api/services/Service/Upgrade"
	"api/services/VO/ExcelVo"
	"api/services/VO/InvoiceVo"
	"api/services/VO/Request"
	"api/services/VO/Response"
	"api/services/dao/Credit"
	"api/services/dao/Cvs"
	"api/services/dao/InvoiceDao"
	"api/services/dao/Orders"
	postBag "api/services/dao/PostBag"
	"api/services/dao/Store"
	"api/services/dao/UserAddressData"
	"api/services/dao/Withdraw"
	"api/services/dao/member"
	"api/services/dao/product"
	"api/services/dao/transfer"
	"api/services/database"
	"api/services/entity"
	"api/services/util/log"
	"api/services/util/tools"
	"fmt"
	"strings"
	"time"
)

//信用卡審單處理
func HandleCreditAudit(params Request.ErpRequest) error {
	engine := database.GetMysqlEngine()
	defer engine.Close()
	if err := engine.Session.Begin(); err != nil {
		log.Error("Set Begin Error", err)
		return fmt.Errorf("系統錯誤")
	}
	//取訂單資料
	for _, OrderId := range params.OrderId {
		data, err := Credit.GetCreditByOrderId(engine, OrderId)
		if err != nil {
			log.Error("Get Order Data Error", err)
			return fmt.Errorf("系統錯誤")
		}
		if params.AuditType == Enum.CreditAuditRelease {
			if err := HandleCreditAuditRelease(engine, data); err != nil {
				if err := engine.Session.Rollback(); err != nil {
					log.Error("Rollback Error")
				}
				return fmt.Errorf("系統錯誤")
			}
		}
		if params.AuditType == Enum.CreditAuditRefused {
			if err := HandleCreditAuditRefused(engine, data); err != nil {
				if err := engine.Session.Rollback(); err != nil {
					log.Error("Rollback Error")
				}
				return fmt.Errorf("系統錯誤")
			}
		}
		if params.AuditType == Enum.CreditAuditNote {
			if err := HandleCreditAuditNote(engine, data); err != nil {
				if err := engine.Session.Rollback(); err != nil {
					log.Error("Rollback Error")
				}
				return fmt.Errorf("系統錯誤")
			}
		}
		if params.AuditType == Enum.CreditAuditPending {
			if err := HandleCreditAuditPending(engine, data); err != nil {
				if err := engine.Session.Rollback(); err != nil {
					log.Error("Rollback Error")
				}
				return fmt.Errorf("系統錯誤")
			}
		}
	}
	if err := engine.Session.Commit(); err != nil {
		return fmt.Errorf("系統錯誤")
	}
	return nil
}

//審單放行
func HandleCreditAuditRelease(engine *database.MysqlSession, data entity.GwCreditAuthData) error {
	//_ = Notification.SendWaitShipMessage(engine, OrderData)
	//判斷目前狀態
	if data.AuditStatus == Enum.CreditAuditRefused || data.AuditStatus == Enum.CreditAuditInit {
		return fmt.Errorf("此訂單不得更改狀態")
	}
	//GW信用卡改變裝態
	data.AuditStatus = Enum.CreditAuditRelease
	data.ReleaseTime = time.Now()
	data.AuditStaff = "john" //fixme  修改取管理者
	if err := Credit.ChangeGwCreditStatus(engine, data); err != nil {
		log.Error("Change Gw Credit Status Error", err)
		return err
	}
	//更改會員信用卡已過卡的紀錄
	if err := MemberService.ChangeCreditCardIsRelease(engine, data.CardId); err != nil {
		log.Error("Change Credit Card Release Error", err)
		return err
	}
	//訂單改變裝態
	OrderData, err := Orders.GetOrderByOrderId(engine, data.OrderId)
	if err != nil {
		log.Error("Get Order Error", err)
		return err
	}
	if err := OrderService.ChangeOrderStatus(engine, OrderData, Enum.OrderSuccess); err != nil {
		return err
	}
	return nil
}

//審單照會
func HandleCreditAuditNote(engine *database.MysqlSession, data entity.GwCreditAuthData) error {
	//訂單改變裝態
	data.AuditStatus = Enum.CreditAuditNote
	data.NoteTime = time.Now()
	data.AuditStaff = "john" //fixme  修改取管理者
	if err := Credit.ChangeGwCreditStatus(engine, data); err != nil {
		return fmt.Errorf("系統錯誤")
	}
	return nil
}

//審單待決
func HandleCreditAuditPending(engine *database.MysqlSession, data entity.GwCreditAuthData) error {
	//訂單改變裝態
	data.AuditStatus = Enum.CreditAuditPending
	data.PendingTime = time.Now()
	data.AuditStaff = "john" //fixme  修改取管理者
	if err := Credit.ChangeGwCreditStatus(engine, data); err != nil {
		return fmt.Errorf("系統錯誤")
	}
	return nil
}

//審單拒決
func HandleCreditAuditRefused(engine *database.MysqlSession, data entity.GwCreditAuthData) error {
	OrderData, err := Orders.GetOrderByOrderId(engine, data.OrderId)
	if err != nil {
		log.Error("Get Order Error", err)
		return err
	}
	if err := OrderService.ChangeOrderStatus(engine, OrderData, Enum.OrderCancel); err != nil {
		return err
	}
	//退款
	if err := Balance.OrderCancelPaymentRefund(engine, OrderData); err != nil {
		return fmt.Errorf("系統錯誤")
	}
	//訂單改變裝態
	data.AuditStatus = Enum.CreditAuditRefused
	data.RefusedTime = time.Now()
	data.AuditStaff = "john" //fixme  修改取管理者
	if err := Credit.ChangeGwCreditStatus(engine, data); err != nil {
		return fmt.Errorf("系統錯誤")
	}
	return nil
}

//審單備註
func HandleCreditAuditCreditAuditMemo(params Request.ErpAuditMemoRequest) error {
	engine := database.GetMysqlEngine()
	defer engine.Close()
	data, err := Credit.GetGwCreditByOrderId(engine, params.OrderId)
	if err != nil {
		log.Error("Get Gw Credit Data Error", err)
		return fmt.Errorf("系統錯誤")
	}
	if len(data.OrderId) == 0 {
		return fmt.Errorf("系統無此記錄")
	}
	data.Memo = fmt.Sprintf("%s \n %s", data.Memo, params.Memo)
	if err := Credit.UpdateGwCreditData(engine, data); err != nil {
		log.Error("Update Gw Credit Error", err)
		return fmt.Errorf("系統錯誤")
	}
	return nil
}

//審單列表
func HandleCreditAuditCreditAuditList(params Request.ErpAuditListRequest) (Response.AuditListResponse, error) {
	engine := database.GetMysqlEngine()
	defer engine.Close()
	var resp Response.AuditListResponse
	data, err := Credit.GetAllGwCreditDataByAuditStatus(engine, params.Tabs)
	if err != nil {
		log.Error("Get Gw Credit Data Error", err)
		return resp, fmt.Errorf("系統錯誤")
	}
	if err := Erp.GetAuditList(engine, data, &resp); err != nil {
		return resp, fmt.Errorf("系統錯誤")
	}
	resp.AuditTabs.Wait = Erp.CountGwCreditByAuditStatus(engine, Enum.CreditTransTypeAuth, Enum.CreditAuditWait)
	resp.AuditTabs.Note = Erp.CountGwCreditByAuditStatus(engine, Enum.CreditTransTypeAuth, Enum.CreditAuditNote)
	resp.AuditTabs.Pending = Erp.CountGwCreditByAuditStatus(engine, Enum.CreditTransTypeAuth, Enum.CreditAuditPending)
	resp.AuditTabs.Capture = Erp.CountGwCreditByAuditStatus(engine, Enum.CreditTransTypeAuth, Enum.CreditAuditRelease)
	resp.AuditTabs.Refund = Erp.CountGwCreditByAuditStatus(engine, Enum.CreditTransTypeRefund, Enum.CreditAuditRelease)
	resp.AuditTabs.Void = Erp.CountGwCreditByAuditStatus(engine, Enum.CreditTransTypeVoid, Enum.CreditAuditRelease)

	return resp, nil
}

//會員中止升級服務
func HandleMemberSuspendUpgrade(params Request.ErpDemoteRequest) error {
	engine := database.GetMysqlEngine()
	defer engine.Close()
	if err := engine.Session.Begin(); err != nil {
		log.Error("Set Begin Error", err)
		return fmt.Errorf("1001001")
	}
	//取出會員資料
	userData, err := member.GetMemberDataByUid(engine, params.UserId)
	if err != nil {
		return fmt.Errorf("1001001")
	}
	if len(userData.Uid) == 0 {
		return fmt.Errorf("1001001")
	}
	if err := Upgrade.MemberSuspendUpgradeService(engine, userData); err != nil {
		log.Error("Member Demote Level  Error", err)
		if err := engine.Session.Rollback(); err != nil {
			log.Error("Rollback Error")
		}
		return fmt.Errorf("1001001")
	}
	if err := engine.Session.Commit(); err != nil {
		return fmt.Errorf("1001001")
	}
	if err := Notification.SendUpgradeStopMessage(engine, userData); err != nil {
		return fmt.Errorf("1001001")
	}
	return nil
}

//提領回檔處理
func HandleEachFile(filename string) error {
	engine := database.GetMysqlEngine()
	defer engine.Close()
	if err := engine.Session.Begin(); err != nil {
		log.Error("Set Begin Error", err)
		return fmt.Errorf("系統錯誤")
	}
	data, err := Excel.ReadAchResponseExcel(filename)
	if err != nil {
		log.Debug("Read Ach Response Excel Error", err)
		return err
	}
	for _, v := range data {
		log.Debug("get row", v)
		log.Debug("get transId", v.PSEQ)
		data, err := Withdraw.GetAchWithdrawDataByTransId(engine, v.PSEQ)
		if err != nil {
			log.Error("Get Withdraw Error", err)
			if err := engine.Session.Rollback(); err != nil {
				log.Error("Rollback Error")
			}
			return err
		}
		if data.ResponseStatus == Enum.WithdrawStatusInit {
			if v.RCODE == "00" {
				data.ResponseCode = v.RCODE
				data.ResponseStatus = Enum.WithdrawStatusSuccess
			} else {
				data.ResponseCode = v.RCODE
				data.ResponseStatus = Enum.WithdrawStatusFailed
				//退回會員提領金額
				if err := WithdrawFailedRefund(engine, data); err != nil {
					log.Error("Refund Balance Error", err)
					if err := engine.Session.Rollback(); err != nil {
						log.Error("Rollback Error")
					}
					return err
				}
			}
			if err := Withdraw.UpdateWithdrawDataByResponse(engine, data); err != nil {
				log.Error("Update Withdraw Response Error", err)
				if err := engine.Session.Rollback(); err != nil {
					log.Error("Rollback Error")
				}
				return err
			}
		}
	}
	if err := engine.Session.Commit(); err != nil {
		return fmt.Errorf("系統錯誤")
	}
	return nil
}

//會員餘額列表產生Excel
func HandleMemberReportExporter() (string, error) {
	engine := database.GetMysqlEngine()
	defer engine.Close()
	data, err := member.GetAllMemberData(engine)
	if err != nil {
		log.Debug("Read Ach Response Excel Error", err)
	}
	var report []ExcelVo.MemberReportVo
	for _, v := range data {
		var rep ExcelVo.MemberReportVo
		rep.MemberId = tools.MaskerPhoneLater(v.Mphone)
		rep.TerminalId = v.TerminalId
		rep.Balance = int64(Balance.GetBalanceByUid(engine, v.Uid))
		report = append(report, rep)
	}

	var rep ExcelVo.MemberReportVo
	rep.MemberId = "checkne"
	rep.Balance = int64(Balance.GetBalanceByUid(engine, "U2021010100001"))
	report = append(report, rep)

	filename, err := Excel.MemberNew().ToMemberReportFile(report)
	if err != nil {
		return "", fmt.Errorf("系統錯誤！")
	}
	log.Debug("filename", filename)
	return filename, nil
}

//發票報表
func HandleInvoiceReportExporter(Year, Month string) (string, error) {
	engine := database.GetMysqlEngine()
	defer engine.Close()
	data, err := InvoiceDao.GetInvoiceByYearMonth(engine, Year, Month)
	if err != nil {
		log.Debug("Get Invoice Data Error", err)
	}
	var report []ExcelVo.InvoiceReportVo
	for _, v := range data {
		var rep ExcelVo.InvoiceReportVo
		rep.InvoiceYm = fmt.Sprintf("%s%s", v.Year, v.Month)
		rep.InvoiceNumber = fmt.Sprintf("%s%s", v.InvoiceTrack, v.InvoiceNumber)
		rep.CreateTime = v.CreateTime.Format("2006/01/02")
		rep.Sales = v.Sales
		rep.Tax = v.Tax
		rep.Amount = v.Amount
		if v.Identifier != "0000000000" {
			rep.Identifier = v.Identifier
		}
		rep.InvoiceStatus = ""
		if v.InvoiceStatus == Enum.InvoiceStatusCancel {
			rep.InvoiceStatus = "發票作廢"
		}
		var detail []InvoiceVo.Details
		_ = tools.JsonDecode([]byte(v.Detail), &detail)
		var s []string
		for _, row := range detail {
			s = append(s, row.ProductName)
		}
		rep.Products = strings.Join(s, ",")
		report = append(report, rep)
	}
	filename, err := Excel.InvoiceNew().ToExcelReportFile(report)
	if err != nil {
		return "", fmt.Errorf("系統錯誤！")
	}
	return filename, nil
}
//取出會員發票資料
func HandleUserInvoiceReportExporter(userId string) (string, error) {
	engine := database.GetMysqlEngine()
	defer engine.Close()
	data, err := InvoiceDao.GetInvoiceAndOrderByUserid(engine, userId, "2021-07-17")
	if err != nil {
		log.Debug("Get Invoice Data Error", err)
	}
	var report []ExcelVo.UserInvoiceReportVo
	for _, v := range data {
		var rep ExcelVo.UserInvoiceReportVo
		rep.OrderId = v.OrderId
		rep.InvoiceNumber = fmt.Sprintf("%s%s", v.InvoiceTrack, v.InvoiceNumber)
		rep.CreateTime = v.CreateTime.Format("2006/01/02")
		rep.Amount = v.Amount
		report = append(report, rep)
	}
	filename, err := Excel.UserInvoiceNew().ToInvoiceExcelFile(report)
	if err != nil {
		return "", fmt.Errorf("系統錯誤！")
	}
	return filename, nil
}

//取出提領列表
func HandleGetWithdraw(params Request.SearchWithdrawRequest) (Response.SearchWithdrawResponse, error) {
	engine := database.GetMysqlEngine()
	defer engine.Close()
	var resp Response.SearchWithdrawResponse

	data, err := Withdraw.GetWithdrawData(engine, params)
	if err != nil {
		log.Debug("Get Withdraw Error", err)
		return resp, err
	}
	for _, v := range data {
		resp.WithdrawList = append(resp.WithdrawList, v.GetWithdraw())
	}
	resp.WithdrawTabs.Wait = Withdraw.CountWithdrawDataByStatus(engine, Enum.WithdrawStatusWait)
	resp.WithdrawTabs.Pending = Withdraw.CountWithdrawDataByStatus(engine, Enum.WithdrawStatusPending)
	return resp, nil
}

func HandleWithdrawChangeStatus(params Request.ChangeWithdrawRequest) error {
	engine := database.GetMysqlEngine()
	defer engine.Close()

	for _, v := range params.WithdrawId {
		data, err := Withdraw.GetAchWithdrawDataByWithdrawId(engine, v)
		if err != nil {
			log.Error("Get Withdraw Database Error", err)
			continue
		}
		if len(data.WithdrawId) == 0 {
			log.Error("Not Withdraw data")
			continue
		}
		if data.WithdrawStatus == Enum.ReturnStatusSuccess {
			log.Error("Withdraw Status Success")
			continue
		}
		switch params.Status {
		case Enum.WithdrawStatusWait:
			data.WithdrawStatus = Enum.WithdrawStatusWait
		case Enum.WithdrawStatusPending:
			data.WithdrawStatus = Enum.WithdrawStatusPending
		}
		if err := Withdraw.UpdateWithdrawData(engine, data); err != nil {
			log.Debug("Update Withdraw Error", err)
			return fmt.Errorf("1001001")
		}
	}
	return nil
}
//取訂單列表
func HandleSearchOrders(params Request.SearchOrderRequest) (Response.SearchOrderResponse, error) {
	engine := database.GetMysqlEngine()
	defer engine.Close()
	var resp Response.SearchOrderResponse
	count, err := Orders.CountOrderJoinStoreJoinRefund(engine, params)
	if err != nil {
		log.Debug("Count Order Error", err)
		return resp, fmt.Errorf("1001001")
	}
	data, err := Orders.GetOrderJoinStoreJoinRefund(engine, params)
	if err != nil {
		log.Debug("Get Order Error", err)
		return resp, fmt.Errorf("1001001")
	}
	for _, v := range data {
		var res Response.OrdersResponse
		res = v.GetSearchOrders()
		res.OrderStatusText = OrderService.GetOrderStatusText(v.Order)
		buyer, _ := member.GetMemberDataByUid(engine, v.Order.BuyerId)
		res.BuyerId = buyer.TerminalId
		resp.Orders = append(resp.Orders, res)
	}
	resp.Count = count
	return resp, nil
}

//取出訂單內容
func HandleSearchOrdersDetail(params Request.ErpSearchOrderRequest) (Response.SearchOrderDetailResponse, error) {
	engine := database.GetMysqlEngine()
	defer engine.Close()
	var resp Response.SearchOrderDetailResponse
	data, err := Orders.GetOrderAndSellerByOrderId(engine, params.OrderId)
	if err != nil {
		log.Debug("Get Order Error", err)
		return resp, fmt.Errorf("1001001")
	}
	if len(data.Order.OrderId) == 0 {
		return resp, fmt.Errorf("1001007")
	}
	resp = data.GetSearchOrder()
	//訂單狀態
	resp.OrderStatusText = OrderService.GetOrderStatusText(data.Order)
	detail, err := Orders.GetOrderDetailListByOrderId(engine, data.Order.OrderId)
	if err != nil {
		log.Debug("Get Order Detail Error", err)
		return resp, fmt.Errorf("1001001")
	}
	for _, v := range detail {
		var res Response.SearchOrderDetail
		res = v.GetOrderDetail()
		resp.Detail = append(resp.Detail, res)
	}
	return resp, nil
}

//取出訂單退貨退款內容
func HandleSearchOrdersRefund(params Request.ErpSearchOrderRequest) (Response.SearchOrderRefundResponse, error) {
	engine := database.GetMysqlEngine()
	defer engine.Close()
	var resp Response.SearchOrderRefundResponse
	data, err := Orders.GetReturnAndRefundByOrderId(engine, params.OrderId)
	if err != nil {
		log.Debug("Get Order Refund Error", err)
		return resp, fmt.Errorf("1001001")
	}
	for _, v := range data {
		if v.RefundType == Enum.TypeRefund {
			resp.Refund = append(resp.Refund, v.GetRefund())
		} else {
			resp.Return = append(resp.Return, v.GetReturn())
		}
	}
	return resp, nil
}

//取出訂單運送內容
func HandleSearchOrdersShipping(params Request.ErpSearchOrderRequest) (Response.ShippingResponse, error) {
	engine := database.GetMysqlEngine()
	defer engine.Close()
	var resp Response.ShippingResponse
	data, err := Orders.GetOrderByOrderId(engine, params.OrderId)
	if err != nil {
		log.Debug("Get Order Error", err)
		return resp, fmt.Errorf("1001001")
	}
	resp.ShipMode = Enum.Shipping[data.ShipType]
	resp.ShipNumber = data.ShipNumber
	resp.Receiver = data.ReceiverName
	resp.ReceiverPhone = data.ReceiverPhone
	resp.ShipTime = ""
	if !data.ShipTime.IsZero() {
		resp.ShipTime = data.ShipTime.Format("2006/01/02 15:04")
	}
	resp.ShipTrader = Erp.GetShipTrader(data.ShipType, data.ShipText)
	resp.ReceiverAddr = GetShipAddress(engine, data.ShipType, data.ReceiverAddress)
	resp.SwitchPieces = ""
	resp.SwitchStore = ""
	if tools.InArray([]string{Enum.DELIVERY_POST_BAG1, Enum.DELIVERY_POST_BAG2, Enum.DELIVERY_POST_BAG3, Enum.DELIVERY_I_POST_BAG1,
		Enum.DELIVERY_I_POST_BAG2, Enum.DELIVERY_I_POST_BAG3, Enum.I_POST}, data.ShipType) {
		post, err := postBag.GetConsignmentByOrderId(engine, data.OrderId)
		if err != nil {
			log.Debug("Get Post Consignment Error", err)
			return resp, fmt.Errorf("1001001")
		}
		resp.SellerAddr = post.SellerAddr
	}
	return resp, nil
}

//取出訂單付款內容
func HandleSearchOrdersPayment(params Request.ErpSearchOrderRequest) (Response.PaymentInfoResponse, error) {
	engine := database.GetMysqlEngine()
	defer engine.Close()
	var resp Response.PaymentInfoResponse
	data, err := Orders.GetOrderAndBuyerByOrderId(engine, params.OrderId)
	if err != nil {
		log.Debug("Get Order Error", err)
		return resp, fmt.Errorf("1001001")
	}
	resp.BuyerName = data.Order.BuyerName
	resp.BuyerPhone = tools.MaskerPhone(data.Order.BuyerPhone)
	resp.BuyerId = data.Member.TerminalId
	resp.Receiver = data.Order.ReceiverName
	resp.ReceiverPhone = data.Order.ReceiverPhone
	resp.PaymentMode = Enum.PayWay[data.Order.PayWay]
	resp.PaymentDate = ""
	if !data.Order.PayWayTime.IsZero() {
		resp.PaymentDate = data.Order.PayWayTime.Format("2006/01/02 15:04")
	}
	switch data.Order.PayWay {
	case Enum.Credit:
		auth, err := Credit.GetGwCreditAndMemberCreditByOrderId(engine, data.Order.OrderId, Enum.CreditTransTypeAuth)
		if err != nil {
			log.Debug("Get Gw Credit Error", err)
			return resp, fmt.Errorf("1001001")
		}
		resp.BankName = auth.BankName
		resp.LastFour = fmt.Sprintf("**** **** **** %s", auth.Last4Digits)
		resp.Amount = auth.TramsAmount
		resp.AcquirerBank = "凱基銀行"
		void, err := Credit.GetGwCreditByOrderIdAndTransType(engine, data.Order.OrderId, Enum.CreditTransTypeVoid)
		if err != nil {
			log.Debug("Get Gw Credit Error", err)
			return resp, fmt.Errorf("1001001")
		}
		if len(void.OrderId) != 0 {
			resp.VoidDate = void.CreateTime.Format("2006/01/02 15:04")
		}
		refund, err := Credit.GetGwCreditByOrderIdAndTransType(engine, data.Order.OrderId, Enum.CreditTransTypeRefund)
		if err != nil {
			log.Debug("Get Gw Credit Error", err)
			return resp, fmt.Errorf("1001001")
		}
		if len(refund.OrderId) != 0 {
			resp.RefundDate = refund.CreateTime.Format("2006/01/02 15:04")
		}
	case Enum.Balance:
		resp.Amount = int64(data.Order.TotalAmount)
	case Enum.CvsPay:
		resp.Amount = int64(data.Order.TotalAmount)
	case Enum.Transfer:
		trans, err := transfer.GetTransferByOrderIds(engine, data.Order.OrderId)
		if err != nil {
			log.Debug("Get Transfer Error", err)
			return resp, fmt.Errorf("1001001")
		}
		resp.BankName = Erp.AccountGetBankName(engine, trans.RecdBankAccount)
		resp.LastFour = tools.MaskerBankAccount(trans.RecdBankAccount)
		resp.Amount = tools.StringToInt64(trans.RecdAmount)
		resp.AcquirerBank = "凱基銀行"
	}
	return resp, nil
}

//產生日結表 掛除B2C訂單 B2C訂單不會進信托帳戶
func GeneratorDayStatement(start, end string) error {
	engine := database.GetMysqlEngine()
	defer engine.Close()
	var report []ExcelVo.DayStatementVo
	credit, err := Orders.GetOrderAndCredit(engine, start, end)
	if err != nil {
		log.Error("Get Order and credit Error", err)
		return err
	}
	for k, v := range credit {
		seller, err := member.GetMemberDataByUid(engine, v.Order.SellerId)
		if err != nil {
			log.Debug("Get Member data Error", err)
		}
		buyer, err := member.GetMemberDataByUid(engine, v.Order.BuyerId)
		if err != nil {
			log.Debug("Get Member data Error", err)
		}
		TransAmount := int64(0)
		if v.Auth.TransType == Enum.CreditTransTypeRefund {
			TransAmount -= v.Auth.TramsAmount
		} else {
			TransAmount = v.Auth.TramsAmount
		}
		PlatformFee := int64(0)
		if v.Auth.TransType != Enum.CreditTransTypeRefund {
			PlatformFee = int64(v.Order.PlatformTransFee + v.Order.PlatformShipFee + v.Order.PlatformInfoFee + v.Order.PlatformPayFee)
		}
		res := ExcelVo.DayStatementVo{
			Id:              int64(k + 1),
			TransactionDate: v.Order.CreateTime.Format("2006/01/02"),
			TransactionType: Enum.AuthReport[v.Auth.TransType],
			PaymentType:     Enum.PayWayReport[v.Order.PayWay],
			OrderId:         v.Order.OrderId,
			SellerId:        seller.TerminalId,
			BuyerId:         buyer.TerminalId,
			Amount:          TransAmount,
			PlatformFee:     PlatformFee,
		}
		report = append(report, res)
	}
	transfers, err := Orders.GetOrderAndTransfer(engine, start, end)
	if err != nil {
		log.Error("Get Order and transfer Error", err)
		return err
	}

	for k, v := range transfers {
		seller, err := member.GetMemberDataByUid(engine, v.Order.SellerId)
		if err != nil {
			log.Debug("Get Member data Error", err)
		}
		buyer, err := member.GetMemberDataByUid(engine, v.Order.BuyerId)
		if err != nil {
			log.Debug("Get Member data Error", err)
		}

		res := ExcelVo.DayStatementVo{
			Id:              int64(k + 1),
			TransactionDate: v.Order.CreateTime.Format("2006/01/02"),
			TransactionType: "代收交易款項",
			PaymentType:     Enum.PayWayReport[v.Order.PayWay],
			OrderId:         v.Order.OrderId,
			SellerId:        seller.TerminalId,
			BuyerId:         buyer.TerminalId,
			Amount:          v.Transfer.Amount,
			PlatformFee:     int64(v.Order.PlatformTransFee + v.Order.PlatformShipFee + v.Order.PlatformInfoFee + v.Order.PlatformPayFee),
		}
		report = append(report, res)
	}
	filename, err := Excel.DayStatementNew().ToDayStatementReportFile(report, start, end)
	if err != nil {
		log.Debug("Generator Report File Error", err)
		return err
	}
	log.Debug("sssss", filename)
	return nil
}

/**
超商核帳
accountingType S | P
*/
func GetCvsAccountingChecked(params Request.CvsSendCheckedRequest, accountingType string) ([]entity.CvsShippingWithAmount, error) {
	engine := database.GetMysqlEngine()
	defer engine.Close()
	log.Info(`GetCvsAccountingChecked params`, params, accountingType)

	var whereCondition []string
	var whereBindCondition []interface{}
	whereCondition = append(whereCondition, `type = ?`)
	whereBindCondition = append(whereBindCondition, accountingType)

	whereCondition = append(whereCondition, `cvs_accounting_data.cvs_type = ?`)
	whereBindCondition = append(whereBindCondition, params.Type)

	if len(params.Checked) > 0 {
		whereCondition = append(whereCondition, `checked = ?`)
		whereBindCondition = append(whereBindCondition, params.Checked)
	}

	if len(params.Duration) > 0 {
		date := strings.Split(params.Duration, "-")
		if accountingType == `S` {
			whereCondition = append(whereCondition, `cvs_shipping_data.send_time BETWEEN ? AND ?`)
		} else {
			// whereCondition = append(whereCondition, `cvs_accounting_data.service_type = ?`)
			// whereBindCondition = append(whereBindCondition, 1)
			whereCondition = append(whereCondition, `cvs_shipping_data.receive_time BETWEEN ? AND ?`)
		}
		whereBindCondition = append(whereBindCondition, date[0])
		whereBindCondition = append(whereBindCondition, date[1])
	}

	var joinCondition = `cvs_accounting_data.data_id = cvs_shipping_data.ec_order_no`

	if params.Type == Enum.CVS_HI_LIFE {
		joinCondition = `cvs_accounting_data.data_id = cvs_shipping_data.ship_no`
	}

	return Cvs.GetCvsAccountingChecked(engine, joinCondition, whereCondition, whereBindCondition)
}

func HandleSearchMember(params Request.SearchMemberRequest) (Response.SearchMemberResponse, error) {
	engine := database.GetMysqlEngine()
	defer engine.Close()
	var resp Response.SearchMemberResponse
	count, err := member.CountSearchMemberData(engine, params)
	if err != nil {
		log.Debug("Count Order Error", err)
		return resp, fmt.Errorf("1001001")
	}
	data, err := member.TakeSearchMemberData(engine, params)
	if err != nil {
		log.Debug("Get Order Error", err)
		return resp, fmt.Errorf("1001001")
	}
	for _, row := range data {
		var res Response.MembersResponse
		res.Nickname = row.Username
		res.Account = row.Mphone
		res.Email = row.Email
		res.RegisterTime = row.RegisterTime.Format("2006/01/02 15:04")
		res.Status = Enum.MemberStatus[row.MemberStatus]
		res.Uid = row.Uid
		res.TerminalId = row.TerminalId
		pro, err := product.GetUpgradeProductDataByLevel(engine, row.UpgradeLevel)
		if err != nil {
			return resp, fmt.Errorf("1001001")
		}
		res.Upgrade = pro.ProductName
		//計算加值付款期間
		if !row.UpgradeExpire.IsZero() {
			res.UpgradeCycle = fmt.Sprintf("%s ~ %s", tools.LastMonth(row.UpgradeExpire).Format("2006/01/02"), row.UpgradeExpire.Format("2006/01/02"))
		}
		//使用者目前餘額
		res.Balance = int64(Balance.GetBalanceByUid(engine, row.Uid))
		//使用者目前待撥付餘額
		res.RetainBalance = int64(Balance.GetBalanceRetainsByUid(engine, row.Uid))
		//使用者目前扣留餘額 fixme
		res.DetainBalance = 0
		//使用者目前保留餘額 fixme
		res.WithholdBalance = 0
		//取出加值未付帳單金額
		sum, err := Orders.SumB2cWaitOrdersByUserId(engine, row.Uid)
		if err != nil {
			return resp, fmt.Errorf("1001001")
		}
		res.Outstanding = sum
		//取出管理者數
		manager, err := Store.CountStoreManagerBySellerId(engine, row.Uid)
		if err != nil {
			return resp, fmt.Errorf("1001001")
		}
		res.Manage = manager
		//取出賣場資料
		store, err := Store.GetStoreBySellerId(engine, row.Uid)
		if err != nil {
			return resp, fmt.Errorf("1001001")
		}
		res.Store = int64(len(store))
		for _, v := range store {
			res.StoreName = append(res.StoreName, v.StoreName)
		}
		resp.Members = append(resp.Members, res)
	}
	resp.Count = count
	return resp, nil
}

func HandleMemberStores(account string) (Response.MemberMainResponse, error) {
	engine := database.GetMysqlEngine()
	defer engine.Close()
	var resp Response.MemberMainResponse
	var mainStoreList []string
	user, stores, managers, err := member.GetMemberStores(engine, account)
	if err != nil {
		log.Debug("Get Order Error", err)
		return resp, fmt.Errorf("1001001")
	}
	for _, store := range stores {
		mainStoreList = append(mainStoreList, store.StoreId)
		var s Response.MemberStore
		s.StoreName = store.StoreName
		s.StoreStatus = store.StoreStatus
		s.Created = store.CreateTime.Format("2006/01/02 15:04")
		productCount, instantCount, err := product.CountProductAndInstatOrderByStore(engine, store.StoreId)
		if err != nil {
			log.Error(err.Error())
		} else {
			s.ProductCount = productCount
			s.InstantCount = instantCount
		}
		ms := managers[store.StoreId]
		for _, m := range ms {
			mu, err := member.GetMemberDataByUid(engine, m.UserId)
			if err != nil {
				log.Debug("Get Order Error", err)
				continue
			}
			var mra Response.RankAccount
			mra.StoreName = store.StoreName
			mra.NickName = mu.Username
			mra.Account = mu.Mphone
			mra.Email = m.Email
			mra.MainAccount = account
			mra.CreateTime = m.CreateTime.Format("2006/01/02 15:04")
			mra.DeleteTime = m.UpdateTime.Format("2006/01/02 15:04")
			switch m.RankStatus {
			case Enum.StoreRankDelete:
				s.Deleted = append(s.Deleted, mra)
			case Enum.StoreRankSuccess:
				s.Operated = append(s.Operated, mra)
			}
		}
		resp.Master.Stores = append(resp.Master.Stores, s)

	}
	// resp.Manager.Managers = append(resp.Manager.Managers, mra)
	slaveManagers, err := Store.GetStoreSlaveManagerByUserId(engine, mainStoreList, user.Uid)
	for _, sm := range slaveManagers {
		mu, err := member.GetMemberDataByUid(engine, sm.SellerId)
		if err != nil {
			log.Debug("Get Member Error", err)
			continue
		}
		var mra Response.RankAccount
		mra.StoreName = sm.StoreName
		mra.NickName = mu.Username
		mra.MainAccount = mu.Mphone
		mra.Status = sm.StoreRank.RankStatus
		mra.Email = sm.StoreRank.Email
		mra.Account = user.Mphone
		mra.CreateTime = sm.StoreRank.CreateTime.Format("2006/01/02 15:04")
		mra.DeleteTime = sm.StoreRank.UpdateTime.Format("2006/01/02 15:04")
		resp.Manager.Managers = append(resp.Manager.Managers, mra)
	}

	return resp, nil
}

func HandleSendPlatformMessage(params Request.PlatformMessageRequest) error {
	engine := database.GetMysqlEngine()
	defer engine.Close()
	data, err := Store.GetUserAllStoreData(engine)
	if err != nil {
		return err
	}
	for _, v := range data {
		if params.Type == "STORE" {
			count, _ := product.CountProductByStoreId(engine, v.StoreId)
			if count != 0 {
				if UserAddressData.CountOnlineNotifySystemByUserId(engine, v.SellerId) == 0 {
					if err := Notification.SendSystemNotify(engine, v.SellerId, params.Message, Enum.NotifyMsgTypePlaPlatform, ""); err != nil {
						return err
					}
					user, err := member.GetMemberDataByUid(engine, v.SellerId)
					if err != nil {
						return err
					}
					if err := Sms.SendSystemSms(user); err != nil {
						return err
					}
				}
			}
		}
	}
	return nil
}

func HandleRefundPlatformFee(orderId string) error {
	engine := database.GetMysqlEngine()
	defer engine.Close()
	OrderData, err := Orders.GetOrderByOrderId(engine, orderId)
	if err != nil {
		log.Error("Get Order Data Error", err)
		return err
	}
	comment := fmt.Sprintf("訂單退回費用")
	fee := OrderData.PlatformShipFee + OrderData.PlatformTransFee + OrderData.PlatformInfoFee + OrderData.PlatformPayFee
	if err := Balance.Deposit(OrderData.SellerId, OrderData.OrderId, fee, Enum.BalanceTypeAdjustment, comment); err != nil {
		log.Error("Balance Deposit Error", err)
		return err
	}
	if err := Invoice.ProcessCancelInvoice(engine, OrderData.OrderId, "平台取消訂單"); err != nil {
		return err
	}
	return nil
}
