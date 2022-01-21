package controllers

import (
	"api/services/Service/Edm"
	"api/services/Service/Invoice"
	"api/services/Service/Notification"
	"api/services/VO/Request"
	"api/services/VO/Response"
	"api/services/dao/Customer"
	"api/services/database"
	"api/services/entity"
	Resp "api/services/entity/Response"
	"api/services/errorMessage"
	"api/services/model"
	"api/services/util/log"
	"api/services/util/response"
	"api/services/util/tools"
	"fmt"
	"strconv"
	"time"
	"unicode/utf8"

	"github.com/gin-gonic/gin"
)

//取得訂單列表
func SearchOrdersAction(ctx *gin.Context) {
	resp := response.New(ctx)
	var params Request.SearchOrderRequest
	if err := ctx.ShouldBindJSON(&params); err != nil {
		resp.Fail(errorMessage.GetMessageByCode(1001002)).Send()
		return
	}
	data, err := model.HandleSearchOrders(params)
	if err != nil {
		resp.Fail(errorMessage.GetMessageByCodeString(err.Error())).Send()
		return
	}
	resp.Success("OK").SetData(data).Send()
}

//取訂單內容
func GetOrderDetailAction(ctx *gin.Context) {
	resp := response.New(ctx)
	var params Request.ErpSearchOrderRequest
	if err := ctx.BindQuery(&params); err != nil {
		resp.Fail(errorMessage.GetMessageByCode(1001002)).Send()
		return
	}
	data, err := model.HandleSearchOrdersDetail(params)
	if err != nil {
		resp.Fail(errorMessage.GetMessageByCodeString(err.Error())).Send()
		return
	}
	resp.Success("OK").SetData(data).Send()
}

//取退貨退款資料
func GetRefundAndReturnAction(ctx *gin.Context) {
	resp := response.New(ctx)
	var params Request.ErpSearchOrderRequest
	if err := ctx.BindQuery(&params); err != nil {
		resp.Fail(errorMessage.GetMessageByCode(1001002)).Send()
		return
	}
	data, err := model.HandleSearchOrdersRefund(params)
	if err != nil {
		resp.Fail(errorMessage.GetMessageByCodeString(err.Error())).Send()
		return
	}
	resp.Success("OK").SetData(data).Send()
}

//取運送資訊
func GetShippingAction(ctx *gin.Context) {
	resp := response.New(ctx)
	var params Request.ErpSearchOrderRequest
	if err := ctx.BindQuery(&params); err != nil {
		resp.Fail(errorMessage.GetMessageByCode(1001002)).Send()
		return
	}
	data, err := model.HandleSearchOrdersShipping(params)
	if err != nil {
		resp.Fail(errorMessage.GetMessageByCodeString(err.Error())).Send()
		return
	}
	resp.Success("OK").SetData(data).Send()
}

func GetSearchPaymentAction(ctx *gin.Context) {
	resp := response.New(ctx)
	var params Request.ErpSearchOrderRequest
	if err := ctx.BindQuery(&params); err != nil {
		resp.Fail(errorMessage.GetMessageByCode(1001002)).Send()
		return
	}
	data, err := model.HandleSearchOrdersPayment(params)
	if err != nil {
		resp.Fail(errorMessage.GetMessageByCodeString(err.Error())).Send()
		return
	}
	resp.Success("OK").SetData(data).Send()
}

type CustomerUser struct {
	Customer entity.CustomerData `xorm:"extends"`
	Member   entity.MemberData   `xorm:"extends"`
}
type OrderFull struct {
	Order   entity.OrderData             `xorm:"extends"`
	Product entity.OrderDetail           `xorm:"extends"`
	Message entity.OrderMessageBoardData `xorm:"extends"`
	Refund  entity.OrderRefundData       `xorm:"extends"`
}

//審單執行
func CreditAuditAction(ctx *gin.Context) {
	resp := response.New(ctx)
	var params Request.ErpRequest
	if err := ctx.ShouldBindJSON(&params); err != nil {
		resp.Fail(errorMessage.GetMessageByCode(1001002)).Send()
		return
	}
	err := model.HandleCreditAudit(params)
	if err != nil {
		resp.Fail(1001001, err.Error()).Send()
		return
	}
	resp.Success("OK").SetData(true).Send()
}

//審單備註
func CreditAuditMemoAction(ctx *gin.Context) {
	resp := response.New(ctx)
	var params Request.ErpAuditMemoRequest
	if err := ctx.ShouldBindJSON(&params); err != nil {
		resp.Fail(errorMessage.GetMessageByCode(1001002)).Send()
		return
	}
	err := model.HandleCreditAuditCreditAuditMemo(params)
	if err != nil {
		resp.Fail(1001001, err.Error()).Send()
		return
	}
	resp.Success("OK").SetData(true).Send()
}

//審單列表
func CreditAuditListAction(ctx *gin.Context) {
	resp := response.New(ctx)
	var params Request.ErpAuditListRequest
	if err := ctx.ShouldBindJSON(&params); err != nil {
		resp.Fail(errorMessage.GetMessageByCode(1001002)).Send()
		return
	}
	data, err := model.HandleCreditAuditCreditAuditList(params)
	if err != nil {
		resp.Fail(1001001, err.Error()).Send()
		return
	}
	resp.Success("OK").SetData(data).Send()
}

//中止升級服務
func SuspendUpgradeAction(ctx *gin.Context) {
	resp := response.New(ctx)
	var params Request.ErpDemoteRequest
	if err := ctx.ShouldBindJSON(&params); err != nil {
		resp.Fail(errorMessage.GetMessageByCode(1001002)).Send()
		return
	}
	if err := model.HandleMemberSuspendUpgrade(params); err != nil {
		resp.Fail(errorMessage.GetMessageByCodeString(err.Error())).Send()
		return
	}
	resp.Success("OK").SetData(true).Send()
}

//上傳提領結果檔
func EachResponseAction(ctx *gin.Context) {
	resp := response.New(ctx)
	file, header, err := ctx.Request.FormFile("file")
	if err != nil {
		resp.Fail(errorMessage.GetMessageByCode(1001002)).Send()
		return
	}
	filename, err := tools.UploadFile(file, header)
	if err != nil {
		resp.Fail(errorMessage.GetMessageByCode(1001002)).Send()
		return
	}
	if err := model.HandleEachFile(filename); err != nil {
		resp.Fail(1001001, err.Error()).Send()
		return
	}
	log.Debug("file", file)
	resp.Success("OK").SetData(true).Send()
}

//取出提領資料
func GetWithdrawAction(ctx *gin.Context) {
	resp := response.New(ctx)
	params := Request.SearchWithdrawRequest{}
	if err := ctx.ShouldBindJSON(&params); err != nil {
		resp.Fail(errorMessage.GetMessageByCode(1001002)).Send()
		return
	}
	data, err := model.HandleGetWithdraw(params)
	if err != nil {
		resp.Fail(errorMessage.GetMessageByCode(1001002)).Send()
		return
	}
	resp.Success("OK").SetData(data).Send()
}

//變更提領狀態
func GetWithdrawChangeStatusAction(ctx *gin.Context) {
	resp := response.New(ctx)
	params := Request.ChangeWithdrawRequest{}
	if err := ctx.ShouldBindJSON(&params); err != nil {
		resp.Fail(errorMessage.GetMessageByCode(1001002)).Send()
		return
	}
	if err := model.HandleWithdrawChangeStatus(params); err != nil {
		resp.Fail(errorMessage.GetMessageByCode(1001002)).Send()
		return
	}
	resp.Success("OK").SetData(true).Send()
}

func ResendInvoiceAction(ctx *gin.Context) {
	resp := response.New(ctx)
	var params Request.ErpSearchOrderRequest
	if err := ctx.ShouldBindJSON(&params); err != nil {
		resp.Fail(errorMessage.GetMessageByCode(1001002)).Send()
		return
	}
	err := Invoice.ResendInvoice(params)
	if err != nil {
		resp.Fail(errorMessage.GetMessageByCodeString(err.Error())).Send()
		return
	}
	resp.Success("OK").SetData(true).Send()
}

func FetchOrderWithMultipleField(ctx *gin.Context) {
	// resp := response.New(ctx)
	orderStartTime, orderStartTimeFlag := ctx.GetQuery("orderStartTime")
	payStartTime, paySartTimeFlag := ctx.GetQuery("payStartTime")
	orderEndTime, orderEndTimeFlag := ctx.GetQuery("orderEndTime")
	payEndTime, payEndFlag := ctx.GetQuery("payEndTime")

	fmt.Println(orderStartTime, orderStartTimeFlag)
	fmt.Println(payStartTime, paySartTimeFlag)
	fmt.Println(orderEndTime, orderEndTimeFlag)
	fmt.Println(payEndTime, payEndFlag)
}

// 超商寄件核帳
func GetCvsSendCheckedAction(ctx *gin.Context) {
	resp := response.New(ctx)
	var params Request.CvsSendCheckedRequest
	if err := ctx.BindQuery(&params); err != nil {
		resp.Fail(errorMessage.GetMessageByCode(1001002)).Send()
		return
	}
	data, err := model.GetCvsAccountingChecked(params, `S`)
	if err != nil {
		resp.Fail(errorMessage.GetMessageByCodeString(err.Error())).Send()
		return
	}
	resp.Success("OK").SetData(data).Send()
}

// 超商取件核帳
func GetCvsPickUpCheckedAction(ctx *gin.Context) {
	resp := response.New(ctx)
	var params Request.CvsSendCheckedRequest
	if err := ctx.BindQuery(&params); err != nil {
		resp.Fail(errorMessage.GetMessageByCode(1001002)).Send()
		return
	}
	data, err := model.GetCvsAccountingChecked(params, `P`)
	if err != nil {
		resp.Fail(errorMessage.GetMessageByCodeString(err.Error())).Send()
		return
	}
	resp.Success("OK").SetData(data).Send()
}

func GetProductsAction(ctx *gin.Context) {
	resp := response.New(ctx)
	params := Request.ErpSearchProductRequest{}
	if err := ctx.ShouldBindJSON(&params); err != nil {
		resp.Fail(errorMessage.GetMessageByCode(1001002)).Send()
		return
	}
	data, err := model.HandleProductList(params)
	if err != nil {
		resp.Fail(errorMessage.GetMessageByCode(1001002)).Send()
		return
	}
	resp.Success("OK").SetData(data).Send()
}

func FetchCustomerContact(ctx *gin.Context) {
	engine := database.GetMysqlEngine()
	defer engine.Close()
	var orderData []entity.OrderData
	err := engine.Engine.Table(entity.OrderData{}).Limit(30, 0).Find(&orderData)
	if err != nil {
		log.Error(err.Error())
		return
	}
	ctx.JSON(200, orderData)
}

func FetchCountCustomerQuestionsByStatus(ctx *gin.Context) {
	engine := database.GetMysqlEngine()
	defer engine.Close()
	countOfnone, err := engine.Engine.Table(entity.CustomerData{}).Where(`status = ?`, `none`).Count()
	if err != nil {
		log.Error("erp fecth customer limit error ", err)
		countOfnone = 0
	}
	countOfpending, err := engine.Engine.Table(entity.CustomerData{}).Where(`status = ?`, `pending`).Count()
	if err != nil {
		log.Error("erp fecth customer limit error ", err)
		countOfpending = 0
	}
	ctx.JSON(200, gin.H{
		"none":    countOfnone,
		"pending": countOfpending,
	})

}
func FetchCustomerRelatedQuestions(ctx *gin.Context) {
	questionId := ctx.Query(`questionId`)
	data, err := model.GetCustomerByRelatedId(questionId)
	if err != nil {
		log.Error("erp fecth customer limit error ", err)
		return
	}
	ctx.JSON(200, data)
}
func FetchCustomerAction(ctx *gin.Context) {
	status := ctx.Query("status")
	offset, err := strconv.Atoi(ctx.Query("offset"))

	if offset <= 1 || err != nil {
		log.Error("erp fecth customer offset error ", err)
		offset = 0
	} else {
		offset = (offset - 1) * 5
	}

	limit, err := strconv.Atoi(ctx.Query("limit"))
	if err != nil {
		log.Error("erp fecth customer limit error ", err)
	}
	if limit <= 5 || err != nil {
		limit = 5
	}
	engine := database.GetMysqlEngine()
	defer engine.Close()
	var orderData []CustomerUser

	count, err := engine.Engine.Table(entity.CustomerData{}).
		Where(`customer_data.status=?`, status).OrderBy(`customer_data.create_time DESC`).
		Join(`INNER`, entity.MemberData{}, "member_data.uid = customer_data.user_id").
		Limit(limit, offset).FindAndCount(&orderData)
	if err != nil {
		log.Error(err.Error(), "fetch erp customer error")
		return
	}
	var resp []Response.CustomerContactResponse

	for _, data := range orderData {
		member := Response.CustomerMemberInfo{
			Email:   data.Member.Email,
			Name:    data.Member.Username,
			Account: data.Member.Mphone,
		}
		res := Response.CustomerContactResponse{
			ID:           strconv.Itoa(int(data.Customer.Id)),
			RelatedId:    data.Customer.RelatedId,
			QuestionId:   data.Customer.QuestionId,
			OrderID:      data.Customer.OrderId,
			Type:         data.Customer.Question,
			Content:      data.Customer.Contents,
			CreateTime:   data.Customer.CreateTime.Format("2006/01/02 15:04"),
			ReplyContent: data.Customer.Reply,
			ReplyTime:    data.Customer.ReplyTime.Format("2006/01/02 15:04"),
			Member:       member,
		}
		var memos []entity.CustomerMemo
		err := engine.Engine.Table(entity.CustomerMemo{}).Where(`customer_data_id=?`, data.Customer.Id).Find(&memos)
		defer engine.Close()
		if err != nil {
			log.Error(err.Error(), "customer memos get error")
			continue
		}
		if len(memos) > 0 {
			var respMemos []Response.CustomerMemoInfo
			for _, memo := range memos {
				m := Response.CustomerMemoInfo{
					Staff:      memo.Staff,
					ID:         strconv.Itoa(int(memo.CustomerDataId)),
					Content:    memo.Content,
					CreateTime: memo.CreateTime.Format("2006/01/02 15:04"),
				}
				respMemos = append(respMemos, m)
			}
			res.Memos = respMemos
		}

		resp = append(resp, res)

	}
	ctx.JSON(200, gin.H{
		"questions": resp,
		"count":     count,
	})
}
func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}

func FetchErpOrder(ctx *gin.Context) {
	engine := database.GetMysqlEngine()
	defer engine.Close()

	orderId := ctx.Param("orderId")
	if utf8.RuneCountInString(orderId) != 13 {

		ctx.JSON(200, map[string]string{})
		return
	}
	var orders []OrderFull

	err := engine.Engine.Table(entity.OrderData{}).
		Join(`INNER`, entity.OrderDetail{}, "order_detail.order_id = order_data.order_id").Where("order_data.order_id = ?", orderId).
		Join(`LEFT OUTER`, entity.OrderMessageBoardData{}, "order_message_board_data.order_id = order_data.order_id").
		Join(`LEFT OUTER`, entity.OrderRefundData{}, "order_refund_data.order_id = order_data.order_id").
		Desc("order_message_board_data.id").Find(&orders)

	if err != nil {
		return
	}

	var OrderProduct Response.OrderWithProducts
	var products []Response.OrderProduct
	var inputlist []string
	var productList []string
	var returnProductList []string
	var refundList []string
	var messages []Response.MessageData
	var returnProducts []Response.OrderReturnProduct
	var refunds []Response.OrderRefund
	var emptyMessageFlag = false
	shipMerge := false

	for _, order := range orders {
		if order.Product.ShipMerge == 1 {
			shipMerge = true
		}
		product := Response.OrderProduct{
			ShipFee:  strconv.Itoa(int(order.Product.ShipFee)),
			Name:     order.Product.ProductName,
			Type:     order.Product.ProductSpecName,
			Price:    strconv.Itoa(int(order.Product.ProductPrice)),
			Quantity: strconv.Itoa(int(order.Product.ProductQuantity)),
		}

		message := Response.MessageData{
			OrderID:    order.Message.OrderId,
			Message:    order.Message.Message,
			Role:       order.Message.MessageRole,
			CreateTime: order.Message.CreateTime.Format("2006/01/02 15:04"),
		}

		if !(order.Message == entity.OrderMessageBoardData{}) {
			if !contains(inputlist, strconv.Itoa(int(order.Message.Id))) {
				inputlist = append(inputlist, strconv.Itoa(int(order.Message.Id)))
				messages = append(messages, message)

			}
		} else {
			emptyMessageFlag = true
		}
		if !(order.Refund == entity.OrderRefundData{}) {
			if order.Refund.RefundType == "RETURN" {
				if !contains(returnProductList, order.Refund.ProductSpecId) {
					returnProduct := Response.OrderReturnProduct{
						Status:     order.Refund.Status,
						RefundTime: order.Refund.RefundTime.Format("2006/01/02 15:04"),
						Quantity:   int(order.Refund.Qty),
						ProductID:  order.Refund.ProductSpecId,
						Name:       order.Refund.ProductName,
						ID:         order.Refund.RefundId,
						Price:      strconv.Itoa(int(order.Refund.Amount)),
					}
					returnProductList = append(returnProductList, order.Refund.ProductSpecId)
					returnProducts = append(returnProducts, returnProduct)
				}
			}
			if order.Refund.RefundType == "REFUND" {
				if !contains(refundList, order.Refund.RefundId) {
					refund := Response.OrderRefund{
						Status:     order.Refund.Status,
						RefundTime: order.Refund.RefundTime.Format("2006/01/02 15:04"),
						ID:         order.Refund.RefundId,
						Amount:     strconv.Itoa(int(order.Refund.Total)),
					}
					refundList = append(refundList, order.Refund.RefundId)
					refunds = append(refunds, refund)
				}
			}
		}

		if !contains(productList, order.Product.ProductName+order.Product.ProductName) {
			productList = append(productList, order.Product.ProductName+order.Product.ProductName)
			products = append(products, product)
		}
	}
	OrderProduct.Products = products
	OrderProduct.CreateDate = orders[0].Order.CreateTime.Format("2006/01/02 15:04")
	OrderProduct.OrderID = orders[0].Order.OrderId
	OrderProduct.PayFee = strconv.Itoa(int(orders[0].Order.PlatformPayFee))
	OrderProduct.PlatFee = strconv.Itoa(int(orders[0].Order.PlatformInfoFee + orders[0].Order.PlatformShipFee + orders[0].Order.PlatformPayFee + orders[0].Order.PlatformTransFee))
	OrderProduct.InfoFee = strconv.Itoa(int(orders[0].Order.PlatformInfoFee))
	OrderProduct.TransFee = strconv.Itoa(int(orders[0].Order.PlatformTransFee))
	OrderProduct.TotoalAmount = strconv.Itoa(int(orders[0].Order.TotalAmount))
	OrderProduct.ProductTotalPrice = strconv.Itoa(int(orders[0].Order.SubTotal))
	OrderProduct.PlatShipFee = strconv.Itoa(int(orders[0].Order.PlatformShipFee))
	OrderProduct.ShipFee = strconv.Itoa(int(orders[0].Order.ShipFee))
	OrderProduct.ShipMerge = shipMerge
	OrderProduct.Status = orders[0].Order.OrderStatus
	orderMessage := Response.OrderMessage{}
	if emptyMessageFlag {
		orderMessage.Count = 0
		orderMessage.OrderMessage = messages
	} else {
		orderMessage.Count = len(messages)
		orderMessage.OrderMessage = messages
	}
	OrderProduct.Messages = orderMessage
	OrderProduct.Rufund = refunds
	OrderProduct.ReturnProduct = returnProducts
	ctx.JSON(200, OrderProduct)
}

func ReplyCustomer(ctx *gin.Context) {
	engine := database.GetMysqlEngine()
	defer engine.Close()
	t := time.Now()
	response := Response.ReplyCustomer{}
	customerQuestionId := ctx.Param("questionId")
	message, _ := ctx.GetPostForm("message")
	data, err := Customer.GetCustomerQuestionById(engine, customerQuestionId)
	if err != nil {
		response.Status = "fail"
		ctx.JSON(200, response)
	}
	if message != "" {
		_, err := engine.Session.Table(entity.CustomerData{}).Where("id =?", customerQuestionId).Cols("id", "status", "reply", "reply_time").Update(entity.CustomerData{Status: `finish`, Reply: message, ReplyTime: t})
		if err != nil {
			response.Status = "fail"
			ctx.JSON(200, response)
		} else {
			response.Status = "success"
			err := Notification.SendCustomerReplyMessage(engine, data.UserId, data.OrderId, data.Question, data.QuestionId)
			if err != nil {
				log.Error("Send Customer Reply Message error", err)
			}
			ctx.JSON(200, response)
		}
	} else {
		response.Status = "fail"
		ctx.JSON(200, response)
	}
}

func PendingCustomer(ctx *gin.Context) {
	engine := database.GetMysqlEngine()
	defer engine.Close()

	response := Response.ReplyCustomer{}
	customerQuestionId := ctx.Param("questionId")

	if customerQuestionId != "" {
		_, err := engine.Session.Table(entity.CustomerData{}).Where("id =?", customerQuestionId).Cols("id", "status").Update(entity.CustomerData{Status: `pending`})
		if err != nil {
			response.Status = "fail"
			ctx.JSON(200, response)
		} else {
			response.Status = "success"
			ctx.JSON(200, response)
		}
	} else {
		response.Status = "fail"
		ctx.JSON(200, response)
	}
}
func FinishCustomer(ctx *gin.Context) {
	engine := database.GetMysqlEngine()
	defer engine.Close()

	response := Response.ReplyCustomer{}
	customerQuestionId := ctx.Param("questionId")

	if customerQuestionId != "" {
		_, err := engine.Session.Table(entity.CustomerData{}).Where("id =?", customerQuestionId).Cols("id", "status").Update(entity.CustomerData{Status: `finish`})
		if err != nil {
			response.Status = "fail"
			ctx.JSON(200, response)
		} else {
			response.Status = "success"
			ctx.JSON(200, response)
		}
	} else {
		response.Status = "fail"
		ctx.JSON(200, response)
	}
}

func CustomerMemo(ctx *gin.Context) {
	engine := database.GetMysqlEngine()
	defer engine.Close()
	response := Response.CustomerMemo{}
	customerQuestionId := ctx.Param("questionId")
	Id, _ := strconv.Atoi(customerQuestionId)
	message, _ := ctx.GetPostForm("message")
	var data entity.CustomerMemo
	if customerQuestionId != "" {
		data.Content = message
		data.CustomerDataId = int64(Id)
		data.Staff = `CS`
		_, err := engine.Session.Table("customer_memo").Insert(&data)
		if err != nil {
			response.Status = "fail"
			ctx.JSON(200, response)
		} else {
			response.Status = "success"
			ctx.JSON(200, response)
		}
	} else {
		response.Status = "fail"
		ctx.JSON(200, response)
	}
}

//查詢會員資料
func SearchMemberAction(ctx *gin.Context) {
	resp := response.New(ctx)
	var params Request.SearchMemberRequest
	if err := ctx.ShouldBindJSON(&params); err != nil {
		resp.Fail(errorMessage.GetMessageByCode(1001002)).Send()
		return
	}
	data, err := model.HandleSearchMember(params)
	if err != nil {
		resp.Fail(errorMessage.GetMessageByCodeString(err.Error())).Send()
		return
	}
	resp.Success("OK").SetData(data).Send()
}

//手動執行轉帳被動查詢
func SearchTransferAction(ctx *gin.Context) {
	resp := response.New(ctx)
	now := ctx.Param("date")
	engine := database.GetMysqlEngine()
	defer engine.Close()
	log.Debug("date", now)
	if len(now) == 0 {
		now = time.Now().AddDate(0, 0, -1).Format("20060102")
	}
	temp := ""
	var res []Resp.DETAIL
	for true {
		smx, data, err := model.QueryC2cTransDateTransfer(now, now, temp)
		if err != nil {
			log.Error("Query TransDate Transfer", err, data)
			resp.Fail(errorMessage.GetMessageByCode(1001002)).Send()
			return
		}
		for _, v := range data {
			log.Debug("Trans Date", v)
			if err = model.ProcessTransfer(engine, v); err != nil {
				log.Error("Process Transfer Error", v)
				resp.Fail(errorMessage.GetMessageByCode(1001002)).Send()
				return
			}
			res = append(res, v)
		}
		if len(smx.SvcRs.TEMPDATA) == 0 {
			break
		} else {
			temp = smx.SvcRs.TEMPDATA
		}
	}
	resp.Success("OK").SetData(res).Send()
}

//手動執行信用查詢
func SearchCreditAction(ctx *gin.Context) {
	resp := response.New(ctx)
	orderId := ctx.Param("OrderId")
	if len(orderId) == 0 {
		resp.Fail(errorMessage.GetMessageByCode(1001002)).Send()
		return
	}
	if err := model.QueryC2cCreditTrans(orderId); err != nil {
		resp.Fail(errorMessage.GetMessageByCodeString(err.Error())).Send()
		return
	}
	resp.Success("OK").SetData(true).Send()
}

func CancelCreditOrderAction(ctx *gin.Context) {
	resp := response.New(ctx)
	orderId := ctx.Param("OrderId")
	if len(orderId) == 0 {
		resp.Fail(errorMessage.GetMessageByCode(1001002)).Send()
		return
	}
	if err := model.HandleCreditCancelOrder(orderId); err != nil {
		resp.Fail(errorMessage.GetMessageByCodeString(err.Error())).Send()
		return
	}
	resp.Success("OK").SetData(true).Send()
}

func SetOrderShipAction(ctx *gin.Context) {
	resp := response.New(ctx)
	orderId := ctx.Param("OrderId")
	//fixme 運費問題  要扣除還是賣家要另外支付
	if len(orderId) == 0 {
		resp.Fail(errorMessage.GetMessageByCode(1001002)).Send()
		return
	}
	if err := model.HandleOrderShip(orderId); err != nil {
		resp.Fail(errorMessage.GetMessageByCodeString(err.Error())).Send()
		return
	}
	resp.Success("OK").SetData(true).Send()
}


func MemberStore(ctx *gin.Context) {
	resp := response.New(ctx)
	account := ctx.Param("account")
	data, err := model.HandleMemberStores(account)
	if err != nil {
		resp.Fail(errorMessage.GetMessageByCodeString(err.Error())).Send()
		return
	}
	resp.Success("OK").SetData(data).Send()
}

func SendPlatformMessageAction(ctx *gin.Context) {
	resp := response.New(ctx)
	var params Request.PlatformMessageRequest
	if err := ctx.ShouldBindJSON(&params); err != nil {
		resp.Fail(errorMessage.GetMessageByCode(1001002)).Send()
		return
	}
	if err := model.HandleSendPlatformMessage(params); err != nil {
		resp.Fail(errorMessage.GetMessageByCodeString(err.Error())).Send()
		return
	}
	resp.Success("OK").SetData(true).Send()
}

//重新產生日結表
func AgainDailyReportAction(ctx *gin.Context) {
	resp := response.New(ctx)
	day := ctx.Param("day")
	d, _ := time.ParseInLocation("20060102", day, time.Local)
	start := d.Add(-24 * time.Hour).Format("2006/01/02")
	end := d.Format("2006/01/02")
	if err := model.GeneratorDayStatement(start, end); err != nil {
		log.Error("Day Statement Error", err)
		resp.Fail(errorMessage.GetMessageByCodeString(err.Error())).Send()
	}
	resp.Success("OK").SetData(true).Send()
}

func RefundPlatformFeeAction(ctx *gin.Context) {
	resp := response.New(ctx)
	orderId := ctx.Param("orderId")
	if err := model.HandleRefundPlatformFee(orderId); err != nil {
		resp.Fail(errorMessage.GetMessageByCodeString(err.Error())).Send()
	}
	resp.Success("OK").SetData(true).Send()
}

func CreateTaiwanArea(ctx *gin.Context) {
	resp := response.New(ctx)
	err := model.CreateTaiwanArea()
	if err != nil {
		resp.Fail(errorMessage.GetMessageByCodeString(err.Error())).Send()
		return
	}
	resp.Success("OK").Send()
}
//發送EDM簡訊
func SendEdmMessageAction(ctx *gin.Context) {
	resp := response.New(ctx)
	var params Request.PlatformSendEdmMessageRequest
	if err := ctx.ShouldBindJSON(&params); err != nil {
		resp.Fail(errorMessage.GetMessageByCode(1001002)).Send()
		return
	}
	count, err := Edm.SendEdmAllMember(params.Rule, params.Message, params.Test)
	if err != nil {
		resp.Fail(errorMessage.GetMessageByCodeString(err.Error())).Send()
		return
	}
	resp.Success("OK").SetData(count).Send()
}
//產生短網址
func GenerateShortAction(ctx *gin.Context) {
	resp := response.New(ctx)
	var params Request.GenerateShortRequest
	if err := ctx.ShouldBindJSON(&params); err != nil {
		resp.Fail(errorMessage.GetMessageByCode(1001002)).Send()
		return
	}
	data, err := Edm.HandleGenerateShortLink(params)
	if err != nil {
		resp.Fail(errorMessage.GetMessageByCodeString(err.Error())).Send()
		return
	}
	resp.Success("OK").SetData(data).Send()
}