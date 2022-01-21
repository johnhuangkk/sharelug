package controllers

import (
	"api/services/Service/KgiBank"
	"api/services/Service/Notification"
	"api/services/VO/Request"
	"api/services/dao/Customer"
	"api/services/dao/UserAddressData"
	"api/services/database"
	"api/services/errorMessage"
	"api/services/model"
	"api/services/util/log"
	"api/services/util/response"

	"github.com/gin-gonic/gin"
)

func IndexAction(ctx *gin.Context) {
	//engine := database.GetMysqlEngine()
	//defer engine.Close()
	//Task.HandleInvoiceTask()
	log.Debug("aaaaa")

	c := &model.Dome{}
	s, err := c.Test()
	if  err != nil {
		log.Debug("sss", s)
	}
	log.Debug("CCCC", s)
}

func SendCustomerReply(ctx *gin.Context) {
	resp := response.New(ctx)
	var params Request.CustomerReplyRequest
	if err := ctx.ShouldBindJSON(&params); err != nil {
		resp.Fail(errorMessage.GetMessageByCode(1001002)).Send()
		return
	}
	engine := database.GetMysqlEngine()
	defer engine.Close()

	data, _ := Customer.GetCustomerById(engine, "19")
	content := "經查詢訂單編號：B210527044738 訂單繳費期限已過，無法執行轉帳。\n請您重新下單訂購。\n商品連結為：http://go81.me/MDAwNDgz\n\n如有其他疑問，請隨時來信詢問。\n謝謝。"
	err := Notification.SendCustomerReplyMessage(engine, data.UserId, data.OrderId, content, data.Question)
	if err != nil {
		log.Error("Get image Error", err)
	}
}

func CreditCaptureAction(ctx *gin.Context) {
	KgiBank.HandleC2C3DCapture()
	KgiBank.HandleC2CN3DCapture()
}

//上線通知我
func PostNotifyMessageAction(ctx *gin.Context) {
	resp := response.New(ctx)
	params := Request.NotifyRequest{}
	if err := ctx.ShouldBindJSON(&params); err != nil {
		log.Debug("PostNotifyMessageAction => ", err)
		resp.Fail(errorMessage.GetMessageByCode(1001002)).Send()
		return
	}
	err := UserAddressData.InsertNotifyMessage(params)
	if err != nil {
		resp.Fail(1001001, err.Error()).Send()
		return
	}
	resp.Success("OK").SetData(true).Send()
}

//重新讀取信用卡回應檔
func CreditAgainReadRespondAction(ctx *gin.Context) {
	resp := response.New(ctx)
	if err := KgiBank.HandleRespond(); err != nil {
		resp.Fail(1001001, err.Error()).Send()
		return
	}
	resp.Success("OK").SetData(true).Send()
}

func GetLogsAction(ctx *gin.Context) {
	resp := response.New(ctx)
	params := Request.LogsRequest{}
	err := ctx.BindQuery(&params)
	if err != nil {
		resp.Fail(200, err.Error()).Send()
		return
	}
	log.Warning("Get Logs", params)
	resp.Success("OK").SetData(true).Send()
}
