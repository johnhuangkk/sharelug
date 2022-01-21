package controllers

import (
	"api/config/middleware"
	"api/services/VO/Request"
	"api/services/errorMessage"
	"api/services/model"
	"api/services/util/response"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/smokezl/govalidators"
)

//取得買家訂單列表
func GetBuyerOrderListAction(ctx *gin.Context) {
	resp := response.New(ctx)
	var OrderSearch Request.OrderSearch
	if err := ctx.BindQuery(&OrderSearch); err != nil {
		resp.Fail(200, fmt.Sprintf("%v", err)).Send()
		return
	}
	//取訂單列表
	userData := middleware.GetUserData(ctx)
	data, err := model.GetSearchBuyerOrderData(userData, OrderSearch)
	if err != nil {
		resp.Fail(200, fmt.Sprintf("%v", err)).Send()
		return
	}
	resp.Success("OK").SetData(data).Send()
}

//取得賣家訂單列表
func GetSellerOrderListAction(ctx *gin.Context) {
	resp := response.New(ctx)
	var OrderSearch Request.OrderSearch
	if err := ctx.BindQuery(&OrderSearch); err != nil {
		resp.Fail(200, fmt.Sprintf("%v", err)).Send()
		return
	}
	//取訂單列表
	storeData := middleware.GetStoreData(ctx)
	data, err := model.GetSearchSellerOrderData(storeData, OrderSearch)
	if err != nil {
		resp.Fail(200, fmt.Sprintf("%v", err)).Send()
		return
	}
	resp.Success("OK").SetData(data).Send()
}

//訂單已讀
func PutOrderReadAction(ctx *gin.Context) {
	resp := response.New(ctx)
	var params Request.OrderReadParams
	err := ctx.ShouldBindJSON(&params)
	if err != nil {
		resp.Fail(1001002, fmt.Sprintf("%v", err)).Send()
		return
	}
	validator := govalidators.New()
	if err := validator.LazyValidate(params); err != nil {
		resp.Fail(errorMessage.GetMessageByCode(1001002)).Send()
		return
	}
	UserData := middleware.GetUserData(ctx)
	if err := model.HandleUnread(UserData, params); err != nil {
		resp.Fail(errorMessage.GetMessageByCode(1001002)).Send()
		return
	}
	resp.Success("OK").Send()
}

//設定訂單備註
func SetOrderMemo(ctx *gin.Context) {
	resp := response.New(ctx)
	var params Request.OrderMemoParams
	if err := ctx.ShouldBindJSON(&params); err != nil {
		resp.Fail(1001002, fmt.Sprintf("%v", err)).Send()
		return
	}
	validator := govalidators.New()
	if err := validator.LazyValidate(params); err != nil {
		resp.Fail(errorMessage.GetMessageByCode(1001002)).Send()
		return
	}
	if err := model.HandleOrderMemo(params); err != nil {
		resp.Fail(errorMessage.GetMessageByCode(1001002)).Send()
		return
	}
	resp.Success("OK").Send()
}

//todo 取得退貨退款列表
func GetOrderRefundListAction(ctx *gin.Context) {
	resp := response.New(ctx)
	StoreData := middleware.GetStoreData(ctx)
	if StoreData.StoreId == "" {
		resp.Fail(errorMessage.GetMessageByCode(1001000)).Send()
		return
	}
	var params Request.OrderRefundSearch
	if err := ctx.BindQuery(&params); err != nil {
		resp.Fail(200, fmt.Sprintf("%v", err)).Send()
		return
	}
	params.StoreId = StoreData.StoreId
	data, err := model.SearchRefundList(params)
	if err != nil {
		resp.Fail(1001001, err.Error()).Send()
		return
	}
	resp.Success("OK").SetData(data).Send()
}

//todo 設定提前付款
func SetAdvancePaymentAction(ctx *gin.Context) {
	resp := response.New(ctx)
	var params Request.SetPaymentParams
	if err := ctx.ShouldBindJSON(&params); err != nil {
		resp.Fail(1001002, fmt.Sprintf("%v", err)).Send()
		return
	}
	validator := govalidators.New()
	if err := validator.LazyValidate(params); err != nil {
		resp.Fail(errorMessage.GetMessageByCode(1001002)).Send()
		return
	}
	//todo 出貨狀態 只要以出貨就可提前撥付(排除面交自取及 CVS_PAY 貨到付款)
	userData := middleware.GetUserData(ctx)
	if err := model.HandleAdvancePayment(userData, params); err != nil {
		resp.Fail(1001001, err.Error()).Send()
		return
	}
	resp.Success("OK").SetData(true).Send()
}

//訂單延長撥付
func SetExtensionPaymentAction(ctx *gin.Context) {
	resp := response.New(ctx)
	var params Request.SetPaymentParams
	if err := ctx.ShouldBindJSON(&params); err != nil {
		resp.Fail(1001002, fmt.Sprintf("%v", err)).Send()
		return
	}
	validator := govalidators.New()
	if err := validator.LazyValidate(params); err != nil {
		resp.Fail(errorMessage.GetMessageByCode(1001002)).Send()
		return
	}
	userData := middleware.GetUserData(ctx)
	if err := model.HandleExtensionPayment(userData, params); err != nil {
		resp.Fail(1001001, err.Error()).Send()
		return
	}
	resp.Success("OK").SetData(true).Send()
}

//todo 賣家執行完成交易 fixme
func SetConfirmPaymentAction(ctx *gin.Context) {
	resp := response.New(ctx)
	var params Request.SetPaymentParams
	if err := ctx.ShouldBindJSON(&params); err != nil {
		resp.Fail(1001002, fmt.Sprintf("%v", err)).Send()
		return
	}
	storeData := middleware.GetStoreData(ctx)
	if err := model.HandleCompleteTransaction(storeData, params); err != nil {
		resp.Fail(1001001, err.Error()).Send()
		return
	}
	resp.Success("OK").SetData(true).Send()
}

//執行訂單取消交易
func SetCancelOrderAction(ctx *gin.Context) {
	resp := response.New(ctx)
	var params Request.SetPaymentParams
	if err := ctx.ShouldBindJSON(&params); err != nil {
		resp.Fail(1001002, fmt.Sprintf("%s", err.Error())).Send()
		return
	}
	//fixme 運費問題  要扣除還是賣家要另外支付
	storeData := middleware.GetStoreData(ctx)
	if err := model.HandleCancelOrder(storeData, params); err != nil {
		resp.Fail(errorMessage.GetMessageByCodeString(err.Error())).Send()
		return
	}
	resp.Success("OK").SetData(true).Send()
}
