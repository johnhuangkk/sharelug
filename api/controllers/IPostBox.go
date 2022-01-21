package controllers

import (
	"api/services/Service/IPost"
	"api/services/VO/IPOSTVO"
	"api/services/errorMessage"
	"api/services/util/log"
	"api/services/util/response"
	"encoding/json"
	"github.com/gin-gonic/gin"
)

// 郵局呼叫api
func IPostShipStatusUpdateNotification(ctx *gin.Context) {
	resp := response.New(ctx)
	params := IPOSTVO.ShipStatusNotify{}
	if err := ctx.ShouldBindJSON(&params); err != nil {
		log.Debug("IPostShipStatusUpdateNotification => ", err)
		resp.Fail(errorMessage.GetMessageByCode(1001002)).Send()
		return
	}

	jsonData, _ := json.Marshal(params)
	log.Debug("auth Notify params", string(jsonData))

	log.Info("get Notification %s", params)
	// 處理郵局通知訊息
	err := IPost.HandleNotification(params)

	if err != nil {
		resp.Conflict(err.Error()).Send()
		return
	}

	resp.Success("成功").Send()
}
