package controllers

import "C"
import (
	"api/config/middleware"
	"api/services/Service/Carts"
	"api/services/VO/Request"
	"api/services/errorMessage"
	"api/services/model"
	"api/services/util/log"
	"api/services/util/response"
	"github.com/gin-gonic/gin"
	"github.com/smokezl/govalidators"
)

//購物車結帳
func PostPayAction(ctx *gin.Context) {
	resp := response.New(ctx)
	params := &Request.PayParams{}
	if err := ctx.ShouldBindJSON(&params); err != nil {
		log.Error("post pay params", err)
		resp.Fail(errorMessage.GetMessageByCode(1001002)).Send()
		return
	}
	validator := govalidators.New()
	if err := validator.LazyValidate(params); err != nil {
		resp.Fail(errorMessage.GetMessageByCode(1001002)).Send()
		return
	}
	cookie := middleware.GetSessionValue(ctx)
	carts, err := Carts.GetCarts(cookie)
	if err != nil || carts.Products == nil {
		log.Error("GetRedisCarts", err, carts)
		resp.Fail(errorMessage.GetMessageByCode(1001001)).Send()
		return
	}
	// 建立訂單
	UserData := middleware.GetUserData(ctx)
	Response, err := model.HandleCreateOrder(UserData, params, carts)
	if err != nil {
		log.Error("HandleCreateOrder", err)
		resp.Fail(errorMessage.GetMessageByCodeString(err.Error())).Send()
		return
	}
	if len(Response.RtnURL) == 0 {
		carts, _ := Carts.GetCarts(cookie)
		// 刪除redis購物車
		if err = Carts.DeleteRedisCarts(cookie, carts.Style); err != nil {
			log.Error("DeleteRedisCarts", err)
			resp.Fail(1005002, err.Error()).Send()
			return
		}
	}
	resp.Success("OK").SetData(Response).Send()
}

//取出訂單資料
func GetOrderDataAction(ctx *gin.Context) {
	resp := response.New(ctx)
	orderId := ctx.Param("orderId")
	userData := middleware.GetUserData(ctx)
	storeData := middleware.GetStoreData(ctx)
	orderResp, err := model.GetOrderData(userData, storeData, orderId)
	if err != nil {
		resp.Fail(errorMessage.GetMessageByCode(1001001)).Send()
		return
	}
	log.Info("order", orderResp)
	resp.Success("OK").SetData(orderResp).Send()
}

//檢查訂單是否已付款
func GetOrderPaymentCheckAction(ctx *gin.Context) {
	resp := response.New(ctx)
	orderId := ctx.Param("orderId")

	err := model.CheckOrderPayment(orderId)
	if err != nil {
		resp.Success("OK").SetData(false).Send()
		return
	}
	resp.Success("OK").SetData(true).Send()
}
