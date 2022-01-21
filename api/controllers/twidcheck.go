package controllers

import (
	"api/config/middleware"
	"api/services/VO/Request"
	"api/services/errorMessage"
	"api/services/model"
	"api/services/util/response"
	"github.com/gin-gonic/gin"
	"github.com/smokezl/govalidators"
)

//送出身份證驗證
func PostTWIDCheckAction(ctx *gin.Context) {
	resp := response.New(ctx)
	params := &Request.TWIDParams{}
	UserData := middleware.GetUserData(ctx)
	if UserData.Uid == "" {
		resp.Fail(errorMessage.GetMessageByCode(1001000)).Send()
		return
	}
	if err := ctx.ShouldBindJSON(&params); err != nil {
		resp.Fail(errorMessage.GetMessageByCode(1001002)).Send()
		return
	}
	validator := govalidators.New()
	if err := validator.LazyValidate(params); err != nil {
		resp.Fail(errorMessage.GetMessageByCode(1001002)).Send()
		return
	}
	data, err := model.HandleTwIdVerify(params, UserData.Uid)
	if err != nil {
		resp.Fail(errorMessage.GetMessageByCodeString(err.Error())).Send()
		return
	}
	resp.Success("OK").SetData(data).Send()
}


//送出身份證驗證
func ErpTWIDCheckAction(ctx *gin.Context) {
	resp := response.New(ctx)
	params := &Request.TWIDParams{}
	if err := ctx.ShouldBindJSON(&params); err != nil {
		resp.Fail(errorMessage.GetMessageByCode(1001002)).Send()
		return
	}
	validator := govalidators.New()
	if err := validator.LazyValidate(params); err != nil {
		resp.Fail(errorMessage.GetMessageByCode(1001002)).Send()
		return
	}
	data, err := model.HandleErpTwIdVerify(params)
	if err != nil {
		resp.Fail(errorMessage.GetMessageByCodeString(err.Error())).Send()
		return
	}
	resp.Success("OK").SetData(data).Send()
}
//發送EMAIL驗證碼
func VerifyEmailAction(ctx *gin.Context) {
	resp := response.New(ctx)
	params := Request.VerifyEmailParams{}
	UserData := middleware.GetUserData(ctx)
	if UserData.Uid == "" {
		resp.Fail(errorMessage.GetMessageByCode(1001000)).Send()
		return
	}
	if err := ctx.ShouldBindJSON(&params); err != nil {
		resp.Fail(errorMessage.GetMessageByCode(1001002)).Send()
		return
	}
	if err := model.HandleVerifyEmail(params, UserData); err != nil {
		resp.Fail(errorMessage.GetMessageByCodeString(err.Error())).Send()
		return
	}
	resp.Success("OK").SetData(true).Send()
}
//EMAIL驗證 是否完成檢查
func VerifyCheckEmailAction(ctx *gin.Context) {
	resp := response.New(ctx)
	UserData := middleware.GetUserData(ctx)
	if UserData.Uid == "" {
		resp.Fail(errorMessage.GetMessageByCode(1001000)).Send()
		return
	}
	if err := model.HandleCheckVerifyEmail(UserData); err != nil {
		resp.Success("OK").SetData(false).Send()
		return
	}
	resp.Success("OK").SetData(true).Send()
}

func GetIndustryListAction(ctx *gin.Context) {
	resp := response.New(ctx)
	data, err := model.HandleGetIndustryList()
	if err != nil {
		resp.Fail(errorMessage.GetMessageByCodeString(err.Error())).Send()
		return
	}
	resp.Success("OK").SetData(data).Send()
}

func SetIndustryAction(ctx *gin.Context) {
	resp := response.New(ctx)
	params := Request.StoreInfoRequest{}

	UserData := middleware.GetUserData(ctx)
	if UserData.Uid == "" {
		resp.Fail(errorMessage.GetMessageByCode(1001000)).Send()
		return
	}
	if err := ctx.ShouldBindJSON(&params); err != nil {
		resp.Fail(errorMessage.GetMessageByCode(1001002)).Send()
		return
	}
	if err := model.HandleSetStoreIndustry(UserData, params); err != nil {
		resp.Fail(errorMessage.GetMessageByCodeString(err.Error())).Send()
		return
	}
	resp.Success("OK").SetData(true).Send()
}
