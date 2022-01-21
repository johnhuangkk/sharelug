package controllers

import (
	"api/config/middleware"
	"api/services/Service/Upgrade"
	"api/services/VO/Request"
	"api/services/errorMessage"
	"api/services/model"
	"api/services/util/response"
	"github.com/gin-gonic/gin"
)

//取出升級方案
func UpgradeAction(ctx *gin.Context) {
	resp := response.New(ctx)
	UserData := middleware.GetUserData(ctx)
	StoreData := middleware.GetStoreData(ctx)
	data, err := Upgrade.GetUpgradePlan(UserData, StoreData)
	if err != nil {
		resp.Fail(errorMessage.GetMessageByCode(1001002)).Send()
		return
	}
	resp.Success("OK").SetData(data).Send()
}
//賣場升級取結帳 (放入購物車)
func GetUpgradePayAction(ctx *gin.Context) {
	resp := response.New(ctx)
	params := Request.GetB2CPayRequest{}
	if err := ctx.ShouldBindJSON(&params); err != nil {
		resp.Fail(errorMessage.GetMessageByCode(1001002)).Send()
		return
	}
	StoreData := middleware.GetStoreData(ctx)
	UserData := middleware.GetUserData(ctx)
	cookie := middleware.GetSessionValue(ctx)
	if len(cookie) == 0 {
		resp.Fail(errorMessage.GetMessageByCode(1001001)).Send()
		return
	}
	data, err := model.HandleGetB2CPay(cookie, UserData, StoreData, params)
	if err != nil {
		resp.Fail(errorMessage.GetMessageByCodeString(err.Error())).Send()
		return
	}
	resp.Success("OK").SetData(data).Send()
}
//升級方案付款
func UpgradePayAction(ctx *gin.Context) {
	resp := response.New(ctx)
	params := Request.B2CPayRequest{}
	if err := ctx.ShouldBindJSON(&params); err != nil {
		resp.Fail(errorMessage.GetMessageByCode(1001002)).Send()
		return
	}
	cookie := middleware.GetSessionValue(ctx)
	if len(cookie) == 0 {
		resp.Fail(errorMessage.GetMessageByCode(1001001)).Send()
		return
	}
	data, err := model.HandleB2CPay(cookie, params)
	if err != nil {
		resp.Fail(errorMessage.GetMessageByCodeString(err.Error())).Send()
		return
	}
	resp.Success("OK").SetData(data).Send()
}

//未付款帳單數
func CountUpgradeOrderAction(ctx *gin.Context) {
	resp := response.New(ctx)
	UserData := middleware.GetUserData(ctx)
	data, err := Upgrade.HandleCountUnpaidUpgradeOrder(UserData)
	if err != nil {
		resp.Fail(1001002, err.Error()).Send()
		return
	}
	resp.Success("OK").SetData(data).Send()
}

//未付款帳單列表
func GetUpgradeOrderListAction(ctx *gin.Context) {
	resp := response.New(ctx)
	UserData := middleware.GetUserData(ctx)
	data, err := Upgrade.HandleGetUnpaidUpgradeOrder(UserData)
	if err != nil {
		resp.Fail(1001002, err.Error()).Send()
		return
	}
	resp.Success("OK").SetData(data).Send()
}

//付款完成 取升級方案訂單結果
func GetUpgradeOrderAction(ctx *gin.Context) {
	resp := response.New(ctx)
	orderId := ctx.Param("orderId")
	if len(orderId) == 0 {
		resp.Fail(errorMessage.GetMessageByCode(1001002)).Send()
		return
	}
	userdata := middleware.GetUserData(ctx)
	data, err := Upgrade.HandleGetUpgradeOrder(orderId, userdata)
	if err != nil {
		resp.Fail(errorMessage.GetMessageByCodeString(err.Error())).Send()
		return
	}
	resp.Success("OK").SetData(data).Send()
}