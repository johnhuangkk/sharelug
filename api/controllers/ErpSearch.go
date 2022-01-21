package controllers

import (
	"api/services/Service/Balance"
	"api/services/VO/Request"
	"api/services/errorMessage"
	"api/services/model"
	"api/services/util/response"
	"github.com/gin-gonic/gin"
)

func SearchUserBalanceAction(ctx *gin.Context) {
	resp := response.New(ctx)
	params := Request.SearchBalanceRequest{}
	if err := ctx.ShouldBindJSON(&params); err != nil {
		resp.Fail(errorMessage.GetMessageByCode(1001002)).Send()
		return
	}
	data, err := model.HandleSearchBalance(params)
	if err != nil {
		resp.Fail(errorMessage.GetMessageByCodeString(err.Error())).Send()
		return
	}
	resp.Success("OK").SetData(data).Send()
}

func RecalculateBalanceRetainAction(ctx *gin.Context) {
	resp := response.New(ctx)
	uid := ctx.Param("uid")
	if err := Balance.RecalculateBalanceRetain(uid); err != nil {
		resp.Fail(errorMessage.GetMessageByCodeString(err.Error())).Send()
		return
	}
	resp.Success("OK").SetData(true).Send()
}

func RecalculateBalanceAction(ctx *gin.Context) {
	resp := response.New(ctx)
	uid := ctx.Param("uid")
	if err := Balance.RecalculateBalance(uid); err != nil {
		resp.Fail(errorMessage.GetMessageByCodeString(err.Error())).Send()
		return
	}
	resp.Success("OK").SetData(true).Send()
}
