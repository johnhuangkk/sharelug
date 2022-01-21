package controllers

import (
	"api/config/middleware"
	"api/services/VO/Request"
	"api/services/errorMessage"
	"api/services/model"
	"api/services/util/response"
	"github.com/gin-gonic/gin"
)

//提領頁 fixme 增加預設帳號
func WithdrawBankCodeAction(ctx *gin.Context) {
	resp := response.New(ctx)
	userData := middleware.GetUserData(ctx)
	data, err := model.GetBankCode(userData)
	if err != nil {
		resp.Fail(1001001, err.Error()).Send()
		return
	}
	resp.Success("OK").SetData(data).Send()
}

//提領
func WithdrawAction(ctx *gin.Context) {
	resp := response.New(ctx)
	userData := middleware.GetUserData(ctx)
	var params Request.WithdrawRequest
	if err := ctx.ShouldBindJSON(&params); err != nil {
		resp.Fail(errorMessage.GetMessageByCode(1001002)).Send()
		return
	}
	err := model.HandleWithdraw(params, userData)
	if err != nil {
		resp.Fail(errorMessage.GetMessageByCodeString(err.Error())).Send()
		return
	}
	resp.Success("OK").SetData(true).Send()
}

//刪除提領帳號
func WithdrawDeleteAction(ctx *gin.Context) {
	resp := response.New(ctx)
	var params Request.EditWithdrawRequest
	if err := ctx.ShouldBindJSON(&params); err != nil {
		resp.Fail(errorMessage.GetMessageByCode(1001002)).Send()
		return
	}
	UserData := middleware.GetUserData(ctx)
	if err := model.HandleDeleteWithdraw(UserData, params); err != nil {
		resp.Fail(1001001, err.Error()).Send()
		return
	}
	resp.Success("OK").SetData(true).Send()
}

func WithdrawChangeDefaultAction(ctx *gin.Context) {
	resp := response.New(ctx)
	var params Request.EditWithdrawRequest
	if err := ctx.ShouldBindJSON(&params); err != nil {
		resp.Fail(errorMessage.GetMessageByCode(1001002)).Send()
		return
	}
	UserData := middleware.GetUserData(ctx)
	if err := model.HandleChangeDefaultWithdraw(UserData, params); err != nil {
		resp.Fail(1001001, err.Error()).Send()
		return
	}
	resp.Success("OK").SetData(true).Send()
}