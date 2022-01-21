package controllers

import (
	"api/config/middleware"
	"api/services/VO/Request"
	"api/services/entity"
	"api/services/errorMessage"
	"api/services/model"
	"api/services/util/log"
	"api/services/util/response"
	"github.com/gin-gonic/gin"
)
//建立反向帳單
func PostBillOrderAction(ctx *gin.Context) {
	resp := response.New(ctx)
	params := entity.BillOrderParams{}
	if err := ctx.ShouldBindJSON(&params); err != nil {
		log.Error("post bill order params", err)
		resp.Fail(errorMessage.GetMessageByCode(1001002)).Send()
		return
	}
	resp.Fail(errorMessage.GetMessageByCode(1001001)).Send()
	//關閉上架
	//userData := middleware.GetUserData(ctx)
	//data, err := model.HandleNewBillOrder(ctx.Request.Host, userData, params)
	//if err != nil {
	//	log.Error("HandleCreateOrder", err)
	//	resp.Fail(errorMessage.GetMessageByCodeString(err.Error())).Send()
	//	return
	//}
	//resp.Success("OK").SetData(data).Send()
}
//買家查看反向帳單
func ReviewBillOrderAction(ctx *gin.Context) {
	resp := response.New(ctx)
	billId := ctx.Param("billId")
	if len(billId) == 0 {
		resp.Fail(errorMessage.GetMessageByCode(1001002)).Send()
		return
	}
	userData := middleware.GetUserData(ctx)
	data, err := model.HandleGetReviewBillOrder(userData, billId)
	if err != nil {
		log.Error("Handle Get Bill Order Error", err)
		resp.Fail(errorMessage.GetMessageByCodeString(err.Error())).Send()
		return
	}
	resp.Success("OK").SetData(data).Send()
}
//賣家查看反向帳單
func GetBillOrderAction(ctx *gin.Context) {
	resp := response.New(ctx)
	billId := ctx.Param("billId")
	if len(billId) == 0 {
		resp.Fail(errorMessage.GetMessageByCode(1001002)).Send()
		return
	}
	data, err := model.HandleGetBillOrder(billId)
	if err != nil {
		log.Error("Handle Get Bill Order Error", err)
		resp.Fail(errorMessage.GetMessageByCodeString(err.Error())).Send()
		return
	}
	resp.Success("OK").SetData(data).Send()
}
//賣家確認接受
func BillConfirmAction(ctx *gin.Context) {
	resp := response.New(ctx)
	params := Request.BillConfirmRequest{}
	if err := ctx.ShouldBindJSON(&params); err != nil {
		log.Error("Put bill Confirm params", err)
		resp.Fail(errorMessage.GetMessageByCode(1001002)).Send()
		return
	}
	storeData := middleware.GetStoreData(ctx)
	data, err := model.HandleBillConfirm(storeData, params)
	if err != nil {
		log.Error("Handle Bill Confirm Order Error", err)
		resp.Fail(errorMessage.GetMessageByCodeString(err.Error())).Send()
		return
	}
	resp.Success("OK").SetData(data).Send()
}
//買家帳單列表
func GetBillListsAction(ctx *gin.Context) {
	resp := response.New(ctx)
	params := Request.BillListRequest{}
	if err := ctx.BindQuery(&params); err != nil {
		log.Error("Get bill List params", err)
		resp.Fail(errorMessage.GetMessageByCode(1001002)).Send()
		return
	}
	userData := middleware.GetUserData(ctx)
	data, err := model.HandleBillList(userData, params)
	if err != nil {
		log.Error("Handle Bill Order Error", err)
		resp.Fail(errorMessage.GetMessageByCodeString(err.Error())).Send()
		return
	}
	resp.Success("OK").SetData(data).Send()
}
//買家帳單延期
func BillExtensionAction(ctx *gin.Context) {
	resp := response.New(ctx)
	params := Request.BillConfirmRequest{}
	if err := ctx.ShouldBindJSON(&params); err != nil {
		log.Error("Get bill List params", err)
		resp.Fail(errorMessage.GetMessageByCode(1001002)).Send()
		return
	}
	userData := middleware.GetUserData(ctx)
	if err := model.HandleBillExtension(userData, params); err != nil {
		log.Error("Handle Bill Order Error", err)
		resp.Fail(errorMessage.GetMessageByCodeString(err.Error())).Send()
		return
	}
	resp.Success("OK").SetData(true).Send()
}
//買家帳單取消
func BillCancelAction(ctx *gin.Context) {
	resp := response.New(ctx)
	params := Request.BillConfirmRequest{}
	if err := ctx.ShouldBindJSON(&params); err != nil {
		log.Error("Get bill List params", err)
		resp.Fail(errorMessage.GetMessageByCode(1001002)).Send()
		return
	}
	userData := middleware.GetUserData(ctx)
	if err := model.HandleBillCancel(userData, params); err != nil {
		log.Error("Handle Bill Order Error", err)
		resp.Fail(errorMessage.GetMessageByCodeString(err.Error())).Send()
		return
	}
	resp.Success("OK").SetData(true).Send()
}

func GetAllBillListsAction(ctx *gin.Context) {
	resp := response.New(ctx)
	params := Request.BuyerBillListRequest{}
	if err := ctx.BindQuery(&params); err != nil {
		log.Error("Get bill List params", err)
		resp.Fail(errorMessage.GetMessageByCode(1001002)).Send()
		return
	}
	userData := middleware.GetUserData(ctx)
	data, err := model.HandleAllBillList(userData, params)
	if err != nil {
		log.Error("Handle Bill Order Error", err)
		resp.Fail(errorMessage.GetMessageByCodeString(err.Error())).Send()
		return
	}
	resp.Success("OK").SetData(data).Send()
}