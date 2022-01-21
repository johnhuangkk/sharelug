package controllers

import (
	"api/config/middleware"
	"api/services/VO/TokenVo"
	"api/services/errorMessage"
	"api/services/model"
	"api/services/util/log"
	"api/services/util/response"
	"api/services/util/tools"
	"encoding/json"
	"github.com/gin-gonic/gin"
)

//todo 取TOKEN 判斷是否登入 TOKEN是否已過期
func PostTokenAction(ctx *gin.Context) {
	resp := response.New(ctx)
	params := TokenVo.TokenParams{}
	if err := ctx.ShouldBindJSON(&params); err != nil {
		resp.Fail(errorMessage.GetMessageByCode(1001002)).Send()
		return
	}
	jsonString, _ := json.Marshal(params)
	log.Debug("post token params", string(jsonString))
	data, err := model.HandleChangeToken(ctx, params, ctx.ClientIP())
	if err != nil {
		log.Error("Generate UUID Error", err)
	}
	resp.Success("OK").SetData(data).Send()
}
//更換UUID
func PostUuidAction(ctx *gin.Context) {
	resp := response.New(ctx)
	params := TokenVo.TokenParams{}
	if err := ctx.ShouldBindJSON(&params); err != nil {
		resp.Fail(errorMessage.GetMessageByCode(1001002)).Send()
		return
	}
	jsonString, _ := json.Marshal(params)
	log.Debug("post token params", string(jsonString))
	data, err := model.HandleChangeUuid(ctx, params, ctx.ClientIP())
	if err != nil {
		log.Error("Generate UUID Error", err)
	}
	resp.Success("OK").SetData(data).Send()
}
//取小鈴噹數
func GetNoticeAction(ctx *gin.Context) {
	resp := response.New(ctx)
	userData := middleware.GetUserData(ctx)
	storeData := middleware.GetStoreData(ctx)
	data, err := model.HandleNotice(userData, storeData)
	if err != nil {
		log.Error("Get Notice Error", err)
	}
	resp.Success("OK").SetData(data).Send()
}
//檢查UUID是否已登入
func TokenCheckAction(ctx *gin.Context) {
	resp := response.New(ctx)
	params := TokenVo.CheckTokenParams{}
	if err := ctx.ShouldBindJSON(&params); err != nil {
		resp.Fail(errorMessage.GetMessageByCode(1001002)).Send()
		return
	}
	s, _ := tools.JsonEncode(params)
	log.Debug("post token params", s)
	err := model.HandleCheckToken(params)
	if err != nil {
		log.Error("Generate UUID Error", err)
		resp.Success("OK").SetData(false).Send()
		return
	}
	resp.Success("OK").SetData(true).Send()
}