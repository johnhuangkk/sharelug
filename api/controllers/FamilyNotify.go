package controllers

import (
	"api/services/Service/FamilyNotificationService"
	"api/services/VO/FamilyMart"
	"api/services/util/log"
	"api/services/util/response"
	"api/services/util/tools"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

// 超商表示不要驗證 Signature
func checkSignature(params FamilyMart.FamilyParams) bool {

	text := fmt.Sprintf(`ApiKey=%s&Data=%s&TimeStamp=%s`, params.ApiKey, params.DataUrlEncode(), params.TimeStamp)
	priKey := viper.GetString(`MartFamily901.PriKey`)
	hash := tools.SHA512Mac(text, priKey)

	log.Error("簽章 Hash[%s] Signature[%s]", hash, params.Signature)
	return hash == params.Signature
}

// 處理參數 全家表示不用檢查
func handleParams(ctx *gin.Context) (FamilyMart.FamilyParams, error) {
	var params FamilyMart.FamilyParams

	if err := ctx.ShouldBind(&params); err != nil {
		log.Error("格式數量不符", params)
		return params, failResponse(ctx, 003, "格式數量不符")
	}

	//if err := FamilyNotificationService.checkEmptyFields(params); err != nil {
	//	log.Error("檔案格式錯誤", params)
	//	return params, failResponse(ctx, 004, "檔案內容錯誤")
	//}

	// 檢查簽章
	//if !checkSignature(params) {
	//	return params, failResponse(ctx, 999, "簽章不一致")
	//}

	return params, nil
}

func failResponse(ctx *gin.Context, code int, err string) error {
	ctx.String(code, err)
	return fmt.Errorf(err)
}

// 到店寄件
func FamilySendAction(ctx *gin.Context) {

	var params FamilyMart.FamilyParams
	var err error

	rsp := response.New(ctx)
	// Url處理參數
	if params, err = handleParams(ctx); err != nil {
		return
	}

	rsp.XML(FamilyNotificationService.HandleNotify(params, FamilyNotificationService.Sent))
}

// 寄件離店
func FamilyLeaveAction(ctx *gin.Context) {
	var params FamilyMart.FamilyParams
	var err error

	rsp := response.New(ctx)
	// Url處理參數
	if params, err = handleParams(ctx); err != nil {
		return
	}

	rsp.XML(FamilyNotificationService.HandleNotify(params, FamilyNotificationService.SentLeave))
}

// 寄件到店
func FamilyEnterAction(ctx *gin.Context) {
	var params FamilyMart.FamilyParams
	var err error

	rsp := response.New(ctx)
	// Url處理參數
	if params, err = handleParams(ctx); err != nil {
		return
	}

	rsp.XML(FamilyNotificationService.HandleNotify(params, FamilyNotificationService.Enter))
}

// 寄件取貨
func FamilyPickupAction(ctx *gin.Context) {
	var params FamilyMart.FamilyParams
	var err error

	rsp := response.New(ctx)
	// Url處理參數
	if params, err = handleParams(ctx); err != nil {
		return
	}

	rsp.XML(FamilyNotificationService.HandleNotify(params, FamilyNotificationService.PickUp))
}

// 寄件閉轉店
func FamilySwitchAction(ctx *gin.Context) {
	var params FamilyMart.FamilyParams
	var err error

	rsp := response.New(ctx)
	// Url處理參數
	if params, err = handleParams(ctx); err != nil {
		return
	}

	rsp.XML(FamilyNotificationService.HandleNotify(params, FamilyNotificationService.Switch))
}
