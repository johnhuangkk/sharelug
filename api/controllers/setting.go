package controllers

import (
	"api/config/middleware"
	"api/services/Enum"
	"api/services/VO/Request"
	"api/services/VO/Response"
	"api/services/errorMessage"
	"api/services/model"
	"api/services/util/response"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

//收銀機設定
func SettingStoreAction(ctx *gin.Context) {
	resp := response.New(ctx)
	storeData := middleware.GetStoreData(ctx)
	data, err := model.GetMyStoreInfo(storeData)
	if err != nil {
		resp.Fail(1001001, err.Error()).Send()
		return
	}
	resp.Success("OK").SetData(data).Send()
}

//收銀機設定、收銀機名稱、收銀機狀態
func PutSettingStoreAction(ctx *gin.Context) {
	resp := response.New(ctx)
	var params Request.SettingStoreRequest
	if err := ctx.ShouldBindJSON(&params); err != nil {
		resp.Fail(1001002, fmt.Sprintf("%v", err)).Send()
		return
	}
	storeData := middleware.GetStoreData(ctx)
	if storeData.Rank != Enum.StoreRankMaster {
		resp.Fail(errorMessage.GetMessageByCode(1001005)).Send()
		return
	}
	if err := model.SetStoreInfo(storeData.StoreId, params); err != nil {
		resp.Fail(1001001, err.Error()).Send()
		return
	}
	resp.Success("OK").SetData(true).Send()
}

//免運費設定
func SettingFreeShipAction(ctx *gin.Context) {
	resp := response.New(ctx)
	storeData := middleware.GetStoreData(ctx)
	if len(storeData.StoreId) == 0 {
		resp.Fail(errorMessage.GetMessageByCode(1001000)).Send()
		return
	}
	data, err := model.HandleGetStoreFreeShip(storeData)
	if err != nil {
		resp.Fail(errorMessage.GetMessageByCodeString(err.Error())).Send()
		return
	}
	resp.Success("OK").SetData(data).Send()
}

//免運費設定
func PutSettingFreeShipAction(ctx *gin.Context) {
	resp := response.New(ctx)
	params := Request.SettingFreeShipRequest{}
	err := ctx.ShouldBindJSON(&params)
	if err != nil {
		resp.Fail(errorMessage.GetMessageFormatCode(1001002, err)).Send()
		return
	}
	storeData := middleware.GetStoreData(ctx)
	if len(storeData.StoreId) == 0 {
		resp.Fail(errorMessage.GetMessageByCode(1001000)).Send()
		return
	}
	if err := model.HandleSettingFreeShip(storeData, params); err != nil {
		resp.Fail(errorMessage.GetMessageByCodeString(err.Error())).Send()
		return
	}
	resp.Success("OK").SetData(true).Send()
}

//會員資料設定
func SettingAccountAction(ctx *gin.Context) {
	resp := response.New(ctx)
	userData := middleware.GetUserData(ctx)
	data, err := model.GetMyUserInfo(userData.Uid)
	if err != nil {
		resp.Fail(1001001, err.Error()).Send()
		return
	}
	resp.Success("OK").SetData(data).Send()
}

//會員資料設定 儲存
func PutSettingAccountAction(ctx *gin.Context) {
	resp := response.New(ctx)
	var params Request.SettingUserRequest
	if err := ctx.ShouldBindJSON(&params); err != nil {
		resp.Fail(1001002, err.Error()).Send()
		return
	}
	userData := middleware.GetUserData(ctx)
	if err := model.SetMyUserInfo(userData, params); err != nil {
		resp.Fail(errorMessage.GetMessageByCodeString(err.Error())).Send()
		return
	}
	resp.Success("OK").SetData(true).Send()
}

//帳戶重發認證信
func InviteAccountAction(ctx *gin.Context) {
	resp := response.New(ctx)
	userData := middleware.GetUserData(ctx)
	if err := model.HandleInviteAccount(userData); err != nil {
		resp.Fail(1001001, err.Error()).Send()
		return
	}
	resp.Success("OK").SetData(true).Send()
}

//帳戶明細
func MyBalanceAction(ctx *gin.Context) {
	resp := response.New(ctx)
	var params Request.BalanceRequest
	err := ctx.BindQuery(&params)
	if err != nil {
		resp.Fail(1001002, fmt.Sprintf("%v", err)).Send()
		return
	}
	userData := middleware.GetUserData(ctx)
	data, err := model.GetBalanceList(userData, params)
	if err != nil {
		resp.Fail(1001001, err.Error()).Send()
		return
	}
	resp.Success("OK").SetData(data).Send()
}

//帳戶明細下載
func ExportMyBalanceAction(ctx *gin.Context) {
	resp := response.New(ctx)
	var params Request.BalanceRequest
	if err := ctx.BindQuery(&params); err != nil {
		resp.Fail(errorMessage.GetMessageByCode(1001002)).Send()
		return
	}
	userData := middleware.GetUserData(ctx)
	filename, err := model.HandleBalanceExport(userData, params)
	if err != nil {
		resp.Fail(errorMessage.GetMessageByCodeString(err.Error())).Send()
		return
	}
	file, err := os.Open(filename) //Create a file
	if err != nil {
		resp.Fail(errorMessage.GetMessageByCode(1001001)).Send()
		return
	}
	defer file.Close()
	ctx.Writer.Header().Add("Content-Disposition", "attachment; filename=checkne"+time.Now().Format("20060102150405")+".xlsx")
	ctx.Writer.Header().Add("Content-type", "application/octet-stream")
	_, err = io.Copy(ctx.Writer, file)
	if err != nil {
		resp.Fail(errorMessage.GetMessageByCode(1001001)).Send()
		return
	}
}

//保留款項
func MyRetainAction(ctx *gin.Context) {
	resp := response.New(ctx)
	var params Request.RetainRequest
	err := ctx.BindQuery(&params)
	if err != nil {
		resp.Fail(1001002, fmt.Sprintf("%v", err)).Send()
		return
	}
	userData := middleware.GetUserData(ctx)
	data, err := model.GetRetainList(userData, params)
	if err != nil {
		resp.Fail(1001001, err.Error()).Send()
		return
	}
	resp.Success("OK").SetData(data).Send()
}

func MyAccountAction(ctx *gin.Context) {
	resp := response.New(ctx)
	var params Request.MyAccountRequest
	err := ctx.BindQuery(&params)
	if err != nil {
		resp.Fail(1001002, fmt.Sprintf("%v", err)).Send()
		return
	}
	userData := middleware.GetUserData(ctx)
	data, err := model.GetMyAccount(params, userData)
	if err != nil {
		resp.Fail(1001001, err.Error()).Send()
		return
	}
	resp.Success("OK").SetData(data).Send()
}

//通知訊息列表
func NotificationListAction(ctx *gin.Context) {
	resp := response.New(ctx)
	params := Request.NotificationRequest{}
	err := ctx.BindQuery(&params)
	if err != nil {
		resp.Fail(errorMessage.GetMessageByCode(1001002)).Send()
		return
	}
	storeData := middleware.GetStoreData(ctx)
	userData := middleware.GetUserData(ctx)
	data, err := model.GetNotificationList(userData, storeData, params)
	if err != nil {
		resp.Fail(errorMessage.GetMessageByCodeString(err.Error())).Send()
		return
	}
	resp.Success("OK").SetData(data).Send()
}

//訊息通知已讀
func NotificationReadAction(ctx *gin.Context) {
	resp := response.New(ctx)
	params := Request.NotificationReadRequest{}
	if err := ctx.ShouldBindJSON(&params); err != nil {
		resp.Fail(errorMessage.GetMessageByCode(1001002)).Send()
		return
	}
	userData := middleware.GetUserData(ctx)
	count, err := model.NotifyRead(userData, params)
	if err != nil {
		resp.Fail(errorMessage.GetMessageByCodeString(err.Error())).Send()
		return
	}
	var data Response.NoticeResponse
	data.Message = count
	resp.Success("OK").SetData(data).Send()
}

//取出社群帳號連結
func SettingSocialMediaAction(ctx *gin.Context) {
	resp := response.New(ctx)
	storeData := middleware.GetStoreData(ctx)
	data, err := model.GetSocialMediaInfo(storeData)
	if err != nil {
		resp.Fail(errorMessage.GetMessageByCodeString(err.Error())).Send()
		return
	}
	resp.Success("OK").SetData(data).Send()
}

//設定社群帳號連結
func PutSettingSocialMediaAction(ctx *gin.Context) {
	resp := response.New(ctx)
	params := Request.StoreSocialMediaRequest{}
	if err := ctx.ShouldBindJSON(&params); err != nil {
		resp.Fail(errorMessage.GetMessageByCode(1001002)).Send()
		return
	}
	storeData := middleware.GetStoreData(ctx)
	if err := model.SetSocialMediaInfo(storeData, params); err != nil {
		resp.Fail(errorMessage.GetMessageByCodeString(err.Error())).Send()
		return
	}
	resp.Success("OK").SetData(true).Send()
}

func GetDeliveryArea(ctx *gin.Context) {
	resp := response.New(ctx)
	cityCode := ctx.Param("cityCode")
	if len(cityCode) == 0 {
		resp.Fail(errorMessage.GetMessageByCode(1001002)).Send()
	}
	data, err := model.GetTaiwanArea(cityCode)
	if err != nil {
		resp.Fail(errorMessage.GetMessageByCode(1001001)).Send()
	}
	resp.Success("OK").SetData(data).Send()
}

func GetDeliveryCities(ctx *gin.Context) {
	resp := response.New(ctx)
	data, err := model.GetTaiwanCities()
	if err != nil {
		resp.Fail(errorMessage.GetMessageByCode(1001001)).Send()
	}

	resp.Success("OK").SetData(data).Send()
}

func SetStoreSelfDelivery(ctx *gin.Context) {
	resp := response.New(ctx)
	// uid, sid := middleware.GetSession(ctx)
	params := Request.SelfDeliveryAreaRequest{}
	if err := ctx.ShouldBindJSON(&params); err != nil {
		resp.Fail(errorMessage.GetMessageByCode(1001002)).Send()
		return
	}
	if params.Enable == true && len(params.Section) == 0 {
		resp.Fail(errorMessage.GetMessageByCode(1010001)).Send()
		return
	}
	err := model.HandleSelfDeliveryArea(params)
	if err != nil {
		resp.Fail(errorMessage.GetMessageByCode(1001001)).Send()
		return
	}
	resp.Success("OK").Send()
}

func SetStoreSelfDeliveryChargeFree(ctx *gin.Context) {
	resp := response.New(ctx)

	params := Request.SelfDeliveryFeeRequest{}
	if err := ctx.ShouldBindJSON(&params); err != nil {
		resp.Fail(errorMessage.GetMessageByCode(1001002)).Send()
		return
	}
	err := model.HandleSelfDeliveryChargeFree(params)
	if err != nil {
		resp.Fail(errorMessage.GetMessageByCode(1001001)).Send()
		return
	}
	resp.Success("OK").Send()
}

func GetStoreSelfDelivery(ctx *gin.Context) {
	resp := response.New(ctx)
	storeId := ctx.Param("storeId")
	data, err := model.HandleGetStoreSelfDeliveryArea(storeId)
	if err != nil {
		resp.Fail(errorMessage.GetMessageByCode(1001001)).Send()
		return
	}
	resp.Success("OK").SetData(data).Send()
}
