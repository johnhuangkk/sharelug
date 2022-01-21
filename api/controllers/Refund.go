package controllers

import (
	"api/config/middleware"
	"api/services/VO/Request"
	"api/services/errorMessage"
	"api/services/model"
	"api/services/util/log"
	"api/services/util/response"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/smokezl/govalidators"
)

//取得退款列表
func GetRefundListAction(ctx *gin.Context) {
	resp := response.New(ctx)
	var params Request.RefundQuery
	err := ctx.BindQuery(&params)
	if err != nil {
		resp.Fail(200, fmt.Sprintf("%v", err)).Send()
		return
	}
	result, err := model.QueryRefundList(params)
	if err != nil {
		resp.Fail(errorMessage.GetMessageByCode(1001001)).Send()
		return
	}
	resp.Success("OK").SetData(result).Send()
}

//取得退款資料
func GetRefundAction(ctx *gin.Context) {
	resp := response.New(ctx)
	var params Request.RefundQuery
	err := ctx.BindQuery(&params)
	if err != nil {
		resp.Fail(200, fmt.Sprintf("%v", err)).Send()
		return
	}

	result, err := model.QueryRefund(params)
	if err != nil {
		resp.Fail(errorMessage.GetMessageByCode(1001001)).Send()
		return
	}
	resp.Success("OK").SetData(result).Send()
}

//執行退款
func PostRefundAction(ctx *gin.Context) {
	resp := response.New(ctx)
	params := &Request.RefundParams{}
	if err := ctx.ShouldBindJSON(&params); err != nil {
		log.Error("post refund params", err)
		resp.Fail(errorMessage.GetMessageByCode(1001002)).Send()
		return
	}
	validator := govalidators.New()
	if err := validator.LazyValidate(params); err != nil {
		resp.Fail(errorMessage.GetMessageByCode(1001002)).Send()
		return
	}
	storeData := middleware.GetStoreData(ctx)
	if err := model.HandleRefund(storeData, params); err != nil {
		resp.Fail(1001001, err.Error()).Send()
		return
	}
	resp.Success("OK").SetData(true).Send()
}

//取得退貨列表
func GetReturnListAction(ctx *gin.Context) {
	resp := response.New(ctx)
	params := &Request.ReturnListQuery{}
	if err := ctx.BindQuery(&params); err != nil {
		log.Error("post return params", err)
		resp.Fail(errorMessage.GetMessageByCode(1001002)).Send()
		return
	}
	result, err := model.QueryReturnList(params)
	if err != nil {
		resp.Fail(errorMessage.GetMessageByCode(1001001)).Send()
		return
	}
	resp.Success("OK").SetData(result).Send()
}

//取得退貨資料
func GetReturnAction(ctx *gin.Context) {
	resp := response.New(ctx)
	params := &Request.ReturnQuery{}
	if err := ctx.BindQuery(&params); err != nil {
		log.Error("post return params", err)
		resp.Fail(errorMessage.GetMessageByCode(1001002)).Send()
		return
	}
	result, err := model.QueryReturn(params)
	if err != nil {
		resp.Fail(errorMessage.GetMessageByCode(1001001)).Send()
		return
	}

	resp.Success("OK").SetData(result).Send()
}

//執行退貨
func PostReturnAction(ctx *gin.Context) {
	resp := response.New(ctx)
	params := &Request.ReturnParams{}
	if err := ctx.ShouldBindJSON(&params); err != nil {
		log.Error("post return params", err)
		resp.Fail(errorMessage.GetMessageByCode(1001002)).Send()
		return
	}
	validator := govalidators.New()
	if err := validator.LazyValidate(params); err != nil {
		resp.Fail(errorMessage.GetMessageByCode(1001002)).Send()
		return
	}
	err := model.HandleReturn(params)
	if err != nil {
		resp.Fail(1001001, err.Error()).Send()
		return
	}

	resp.Success("OK").SetData(true).Send()
}

//退貨完成確認
func PostReturnConfirmAction(ctx *gin.Context) {
	resp := response.New(ctx)
	params := &Request.ReturnConfirmParams{}
	if err := ctx.ShouldBindJSON(&params); err != nil {
		log.Error("post return params", err)
		resp.Fail(errorMessage.GetMessageByCode(1001002)).Send()
		return
	}
	err := model.HandleReturnConfirm(params)
	if err != nil {
		resp.Fail(1001001, err.Error()).Send()
		return
	}
	resp.Success("OK").SetData(true).Send()
}
