package controllers

import (
	"api/services/Service/Balance"
	"api/services/entity"
	"api/services/errorMessage"
	"api/services/model"
	"api/services/util/log"
	"api/services/util/response"
	"api/services/util/tools"
	"encoding/json"
	"github.com/gin-gonic/gin"
)

func TransferNotifyAction(ctx *gin.Context) {
	resp := response.New(ctx)
	params := entity.TransferParams{}
	if err := ctx.ShouldBind(&params); err != nil {
		resp.Fail(errorMessage.GetMessageByCode(1001001)).Send()
		return
	}
	//記錄LOG
	if err := model.TransferResponse(params); err != nil {
		resp.Fail(errorMessage.GetMessageByCodeString(err.Error())).Send()
		return
	}
	//資料回寫
	if err := model.HandleTransfer(params); err != nil {
		resp.Fail(errorMessage.GetMessageByCodeString(err.Error())).Send()
		return
	}
	resp.Success("OK").Send()
}
//信用卡回傳結果
func AuthNotifyAction(ctx *gin.Context) {
	resp := response.New(ctx)
	params := entity.AuthParams{}
	if err := ctx.Bind(&params); err != nil {
		resp.Fail(errorMessage.GetMessageByCode(1001001)).Send()
	}
	if err := Balance.AuthResponse(params); err != nil {
		resp.Fail(errorMessage.GetMessageByCode(1001001)).Send()
	}
	jsonData, _ := json.Marshal(params)
	log.Debug("auth Notify params", string(jsonData))
	resp.Success("OK").Send()
}
//上傳次特店代碼
func ImportSpecialStoreAction(ctx *gin.Context) {
	resp := response.New(ctx)
	file, header, err := ctx.Request.FormFile("file")
	if err != nil {
		resp.Fail(errorMessage.GetMessageByCode(1001002)).Send()
		return
	}
	filename, err := tools.UploadFile(file, header)
	if err != nil {
		resp.Fail(errorMessage.GetMessageByCode(1001002)).Send()
		return
	}
	if err := model.HandleSpecialStoreCode(filename); err != nil {
		resp.Fail(errorMessage.GetMessageByCodeString(err.Error())).Send()
		return
	}
	resp.Success("OK").SetData(true).Send()
}
