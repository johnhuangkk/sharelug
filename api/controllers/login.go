package controllers

import (
	"api/config/middleware"
	"api/services/Enum"
	"api/services/Service/Carts"
	"api/services/VO/Request"
	"api/services/errorMessage"
	"api/services/model"
	"api/services/util/log"
	"api/services/util/response"
	"api/services/util/validate"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/smokezl/govalidators"
)

//付款頁手機OTP
func SendPayOtpAction(ctx *gin.Context) {
	resp := response.New(ctx)
	params := &Request.PaySendOtpParams{}
	if err := ctx.ShouldBind(&params); err != nil {
		resp.Fail(errorMessage.GetMessageByCode(1001002)).Send()
		return
	}
	//驗證手機號碼格式
	if !validate.VerifyMobileFormat(params.Phone) {
		log.Debug("手機號碼格式錯誤 => ", params.Phone)
		resp.Fail(errorMessage.GetMessageByCode(1002002)).Send()
		return
	}
	//處理會員登入
	data, err := model.HandleLogin(params.Phone)
	if err != nil {
		log.Error("Member Login Error", err)
		resp.Fail(errorMessage.GetMessageByCodeString(err.Error())).Send()
		return
	}
	resp.Success("OK").SetData(data).Send()
}

//發送OTP認證
func SendOtpAction(ctx *gin.Context) {
	resp := response.New(ctx)
	params := &Request.SendOtpParams{}
	if err := ctx.ShouldBind(&params); err != nil {
		resp.Fail(errorMessage.GetMessageByCode(1001002)).Send()
		return
	}
	validator := govalidators.New()
	if err := validator.LazyValidate(params); err != nil {
		log.Debug("login params ", params)
		resp.Fail(1001000, fmt.Sprintf("%v", err)).Send()
		return
	}
	if !validate.VerifyMobileFormat(params.Phone) {
		resp.Fail(errorMessage.GetMessageByCode(1002002)).Send()
		return
	}
	//處理會員登入
	data, err := model.HandleLogin(params.Phone)
	if err != nil {
		log.Debug("get Member Data Error", err)
		resp.Fail(1001001, err.Error()).Send()
		return
	}
	resp.Success("發送OTP成功").SetData(data).Send()
}

//認證OTP
func ValidateOtpAction(ctx *gin.Context) {
	resp := response.New(ctx)
	params := &Request.ValidateOtpParams{}
	if err := ctx.ShouldBind(&params); err != nil {
		resp.Fail(errorMessage.GetMessageByCode(1001002)).Send()
		return
	}
	//validator := govalidators.New()
	//if err := validator.LazyValidate(params); err != nil {
	//	log.Debug("Validate otp params ", params)
	//	resp.Fail(errorMessage.GetMessageFormatCode(1001003, err)).Send()
	//	return
	//}
	if len(params.Phone) == 0 {
		resp.Fail(errorMessage.GetMessageByCode(1002009)).Send()
		return
	}
	if len(params.Code) == 0 {
		resp.Fail(errorMessage.GetMessageByCode(1002008)).Send()
		return
	}

	SessionValue := middleware.GetSessionValue(ctx)
	data, err := model.HandleValidateOtp(SessionValue, params.Code, params.Phone, ctx.ClientIP())
	if err != nil {
		resp.Fail(errorMessage.GetMessageByCodeString(err.Error())).Send()
		return
	}
	resp.Success("OTP驗證成功").SetData(data).Send()
}

//todo 切換收銀機
func ExchangeStoreAction(ctx *gin.Context) {
	resp := response.New(ctx)
	params := &Request.ExchangeStoreParams{}
	if err := ctx.ShouldBindJSON(&params); err != nil {
		resp.Fail(errorMessage.GetMessageByCode(1001002)).Send()
		return
	}
	validator := govalidators.New()
	if err := validator.LazyValidate(params); err != nil {
		log.Debug("Validate otp params ", params)
		resp.Fail(errorMessage.GetMessageFormatCode(1001002, err)).Send()
		return
	}
	data, err := model.HandleExchangeStore(ctx, params)
	if err != nil {
		resp.Fail(errorMessage.GetMessageByCodeString(err.Error())).Send()
		return
	}
	resp.Success("OK").SetData(data).Send()
}

// 登出
func LogoutAction(ctx *gin.Context) {
	resp := response.New(ctx)
	cookie := middleware.GetSessionValue(ctx)
	if err := Carts.DestroyUserLogin(cookie); err != nil {
		resp.Fail(1001001, err.Error()).Send()
		return
	}
	carts, _ := Carts.GetCarts(cookie)
	if err := Carts.DeleteRedisCarts(cookie, carts.Style); err != nil {
		resp.Fail(1001001, err.Error()).Send()
		return
	}
	userData := middleware.GetUserData(ctx)
	log.Info("會員登出", userData.Uid, middleware.GetClientIP(), Enum.SyslogSuccess)
	resp.Success("OK").SetData(true).Send()
}
