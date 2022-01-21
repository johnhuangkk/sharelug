package controllers

import (
	"api/config/middleware"
	"api/services/Enum"
	"api/services/Service/Balance"
	"api/services/VO/Request"
	"api/services/errorMessage"
	"api/services/model"
	"api/services/util/log"
	"api/services/util/response"
	"api/services/util/tools"
	"github.com/gin-gonic/gin"
	"net/http"
)

// 變更預設信用卡
func ChangeDefaultCreditAction(ctx *gin.Context) {
	resp := response.New(ctx)
	var params Request.EditCreditRequest
	if err := ctx.ShouldBindJSON(&params); err != nil {
		resp.Fail(errorMessage.GetMessageByCode(1001002)).Send()
		return
	}
	UserData := middleware.GetUserData(ctx)
	if err := model.HandleChangeDefaultCredit(UserData, params); err != nil {
		resp.Fail(1001001, err.Error()).Send()
		return
	}
	resp.Success("OK").SetData(true).Send()
}

// 刪除信用卡
func DeleteCreditAction(ctx *gin.Context) {
	resp := response.New(ctx)
	var params Request.EditCreditRequest
	if err := ctx.ShouldBindJSON(&params); err != nil {
		resp.Fail(errorMessage.GetMessageByCode(1001002)).Send()
		return
	}
	UserData := middleware.GetUserData(ctx)
	if err := model.HandleDeleteCredit(UserData, params); err != nil {
		resp.Fail(1001001, err.Error()).Send()
		return
	}
	resp.Success("OK").SetData(true).Send()
}

// 取信用卡資料
func GetCardAction(ctx *gin.Context) {
	resp := response.New(ctx)
	UserData := middleware.GetUserData(ctx)
	data, err := model.GetCreditData(UserData)
	if err != nil {
		log.Info("會員讀取信用卡資訊", UserData.Uid, middleware.GetClientIP(), Enum.SysLogFail)
		resp.Fail(1001001, err.Error()).Send()
		return
	}
	log.Info("會員讀取信用卡資訊", UserData.Uid, middleware.GetClientIP(), Enum.SyslogSuccess)
	resp.Success("OK").SetData(data).Send()
}

//取出最後一次使用的地址
func GetDeliveryLastAddressAction(ctx *gin.Context) {
	resp := response.New(ctx)
	var params Request.GetAddressParams
	if err := ctx.BindQuery(&params); err != nil {
		log.Error("Get Address params", err)
		return
	}
	userData := middleware.GetUserData(ctx)
	data, err := model.GetDeliveryLastAddress(userData, params)
	if err != nil {
		resp.Fail(1001001, err.Error()).Send()
		return
	}
	resp.Success("OK").SetData(data).Send()
}

// 開啟3D認證頁
func CreditCheckAction(ctx *gin.Context) {
	content := ctx.Param("content")
	log.Debug("content", tools.Base64DecodeByString(content))
	ctx.Writer.WriteHeader(http.StatusOK)
	html := tools.Base64DecodeByString(content)
	ctx.Writer.Write([]byte(html))
}

// 接收3D認證結果
func CreditConfirmAction(ctx *gin.Context) {
	PayType := ctx.Param("pay")
	params := &Request.Credit3dCheckParams{}
	if err := ctx.ShouldBind(&params); err != nil {
		log.Error("post pay params", err)
		return
	}
	s, _ := tools.JsonEncode(params)
	log.Debug("credit config", s)
	if err := Balance.Credit3DConfirm(PayType, params); err != nil {
		log.Info("Credit 3D Confirm Error", err)
		if params.OrderID[0:2] == "BM" {
			ctx.Redirect(http.StatusFound, "/payV2/fail")
		} else {
			ctx.Redirect(http.StatusFound, "/pay/fail")
		}
		return
	}
	log.Info("Credit 3D Confirm Succ", params.OrderID, PayType)
	if PayType == Enum.OrderTransC2c {
		if params.OrderID[0:2] == "BM" {
			ctx.Redirect(http.StatusFound, "/payV2/succ/"+params.OrderID)
		} else {
			ctx.Redirect(http.StatusFound, "/pay/succ/"+params.OrderID)
		}
	} else if PayType == Enum.OrderTransB2c {
		ctx.Redirect(http.StatusFound, "/store/setting/upgrade-cart/succ/"+params.OrderID)
	} else if PayType == Enum.OrderTransBill {
		ctx.Redirect(http.StatusFound, "/bill/preview/"+params.OrderID)
	}
}
