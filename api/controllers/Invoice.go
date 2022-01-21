package controllers

import (
	"api/config/middleware"
	"api/services/Service/Invoice"
	"api/services/VO/InvoiceVo"
	"api/services/VO/Request"
	"api/services/errorMessage"
	"api/services/model"
	"api/services/util/curl"
	"api/services/util/log"
	"api/services/util/response"
	"api/services/util/tools"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/smokezl/govalidators"
	"github.com/spf13/viper"
	"net/http"
	"net/url"
)

//取出發票載具資料
func GetCarrierAction(ctx *gin.Context) {
	resp := response.New(ctx)
	userData := middleware.GetUserData(ctx)
	data, err := model.HandleGetCarrier(userData)
	if err != nil {
		resp.Fail(1001001, err.Error()).Send()
		return
	}
	resp.Success("OK").SetData(data).Send()
}

//儲存發票載具資料
func PostCarrierAction(ctx *gin.Context) {
	resp := response.New(ctx)
	params := Request.CarrierRequest{}
	if err := ctx.ShouldBindJSON(&params); err != nil {
		resp.Fail(1001002, err.Error()).Send()
		return
	}
	validator := govalidators.New()
	if err := validator.LazyValidate(params); err != nil {
		resp.Fail(200, err.Error()).Send()
		return
	}
	userData := middleware.GetUserData(ctx)
	if err := model.HandlePostCarrier(userData, params); err != nil {
		resp.Fail(1001001, err.Error()).Send()
		return
	}
	resp.Success("OK").SetData(true).Send()
}
//發票列表
func GetInvoiceListAction(ctx *gin.Context) {
	resp := response.New(ctx)
	params := Request.InvoiceListRequest{}
	if err := ctx.BindQuery(&params); err != nil {
		resp.Fail(1001002, err.Error()).Send()
		return
	}
	userData := middleware.GetUserData(ctx)
	data, err := model.HandleGetInvoiceList(userData, params)
	if err != nil {
		resp.Fail(1001001, err.Error()).Send()
		return
	}
	resp.Success("OK").SetData(data).Send()
}
//發票內容
func GetInvoiceDetailAction(ctx *gin.Context) {
	resp := response.New(ctx)
	params := InvoiceVo.InvoiceRequest{}
	if err := ctx.BindQuery(&params); err != nil {
		resp.Fail(1001002, err.Error()).Send()
		return
	}
	userData := middleware.GetUserData(ctx)
	data, err := model.HandleGetInvoiceDetail(userData, params)
	if err != nil {
		resp.Fail(1001001, err.Error()).Send()
		return
	}
	resp.Success("OK").SetData(data).Send()
}

//發票大平台綁定傳送綁定
func InvoiceBindCarrierAction(ctx *gin.Context) {
	token,_ := ctx.GetPostForm("token")
	ban,_ := ctx.GetPostForm("ban")

	log.Debug("Bind Platform Error", token, ban)
	//
	// From 大平台
	params := Request.BindPlatformRequest{}
	if err := ctx.ShouldBind(&params); err != nil {
		log.Error("Bind Platform Error", err)
		ctx.Redirect(http.StatusFound, "/error/404?title=系統訊息&message="+err.Error())
	}
	r, _ := tools.JsonEncode(params)
	log.Debug("Bind Platform Request", r)
	data, err := model.HandleBindPlatform(params)
	if err != nil {
		log.Error("Bind Platform Error", err)
		ctx.Redirect(http.StatusFound, "/error/404?title=系統訊息&message="+err.Error())
	}
	log.Debug("Bind Platform Request")
	ctx.Redirect(http.StatusSeeOther, "/gui-login?Token="+data)
}
//發票大平台綁定驗證綁定資料
func InvoiceVerifyCarrierAction(ctx *gin.Context) {
	resp := response.New(ctx)
	params := Request.BindCarrierRequest{}
	if err := ctx.ShouldBindJSON(&params); err != nil {
		log.Error("Bind Platform Error", err)
		ctx.Redirect(http.StatusFound, "/error/404?title=系統訊息&message="+err.Error())
	}
	userData := middleware.GetUserData(ctx)
	data, err := model.HandleVerifyBindCarrier(userData, params)
	if err != nil {
		resp.Fail(1001001, err.Error()).Send()
		return
	}
	resp.Success("OK").SetData(data).Send()
}
// 大平台綁定小平台
func InvoiceReCaptchaAction(ctx *gin.Context) {
	resp := response.New(ctx)
	params := Request.BindCarrierRequest{}
	if err := ctx.ShouldBindJSON(&params); err != nil {
		log.Error("InvoiceReCaptchaAction Error", err)
		return
	}

	// 組合form data
	PostValue := url.Values{}
	PostValue.Add(`secret`, viper.GetString(`GCP.reCaptchaSecret`))
	PostValue.Add(`response`, params.Token)

	// google reCaptcha response struct
	var gResp = struct {
		Success bool `json:"success"`
		ChallengeTs string `json:"challenge_ts"`
		Hostname string `json:"hostname"`
	}{}
	// google reCaptcha uri
	var reCaptchaUri = `https://www.google.com/recaptcha/api/siteverify`

	data, err := curl.Post(reCaptchaUri, PostValue.Encode())
	_ = json.Unmarshal(data, &gResp)
	log.Info(`InvoiceReCaptchaAction data`, gResp)
	if err != nil {
		log.Error(`InvoiceReCaptchaAction Error`, err.Error())
		resp.Fail(1001001, err.Error()).Send()
		return
	}
	resp.Success("OK").SetData(gResp.Success).Send()
}
//上傳發票開獎結果檔
func ImportAwardedAction(ctx *gin.Context) {
	resp := response.New(ctx)
	file, header, err := ctx.Request.FormFile("file")
	if err != nil {
		resp.Fail(errorMessage.GetMessageByCode(1001002)).Send()
		return
	}
	if _, err := tools.UploadAwardedFile(file, header); err != nil {
		resp.Fail(errorMessage.GetMessageByCode(1001002)).Send()
		return
	}
	if err := Invoice.ReadAwardedPathFiles(); err != nil {
		resp.Fail(1001001, err.Error()).Send()
		return
	}
	resp.Success("OK").SetData(true).Send()
}
//上傳發票字軌
func ImportInvoiceTrackAction(ctx *gin.Context) {
	resp := response.New(ctx)
	params := Request.InvoiceTrackRequest{}
	if err := ctx.ShouldBindJSON(&params); err != nil {
		log.Error("InvoiceReCaptchaAction Error", err)
		return
	}
	if err := Invoice.HandleNewInvoiceTrack(params); err != nil {
		resp.Fail(errorMessage.GetMessageByCodeString(err.Error())).Send()
		return
	}
	resp.Success("OK").SetData(true).Send()
}
//發票內容
func GetAwardedInvoiceAction(ctx *gin.Context) {
	resp := response.New(ctx)
	orderId := ctx.Param("orderId")
	data, err := model.HandleGetAwardedInvoice(orderId)
	if err != nil {
		resp.Fail(1001001, err.Error()).Send()
		return
	}
	resp.Success("OK").SetData(data).Send()
}

