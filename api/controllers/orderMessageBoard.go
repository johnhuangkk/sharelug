package controllers

import (
	"api/config/middleware"
	"api/services/Enum"
	"api/services/Service/OrderMessageBoard"
	"api/services/VO/OrderMessageBoardVo"
	"api/services/VO/Request"
	"api/services/errorMessage"
	"api/services/util/log"
	"api/services/util/response"
	"github.com/gin-gonic/gin"
	"github.com/smokezl/govalidators"
)

// 增加訂單留言板訊息
func AddOrderMessageBoardAction(ctx *gin.Context) {
	resp := response.New(ctx)
	UserData := middleware.GetUserData(ctx)
	StoreData := middleware.GetStoreData(ctx)
	params := OrderMessageBoardVo.OrderMessage{}
	if len(UserData.Uid) == 0 {
		resp.Fail(1001001, "尚未登入").Send()
		return
	}
	if err := ctx.ShouldBindJSON(&params); err != nil {
		log.Debug("AddOrderMessageBoardAction => %s [%s]", err, params)
		resp.Fail(errorMessage.GetMessageByCode(1001002)).Send()
		return
	}
	validator := govalidators.New()
	if err := validator.LazyValidate(params); err != nil {
		resp.Fail(errorMessage.GetMessageByCode(1001002)).Send()
		return
	}
	// 新增留言
	if err := OrderMessageBoard.AddMessage(params, UserData, StoreData); err != nil {
		resp.Fail(errorMessage.GetMessageByCode(1001001)).Send()
		return
	}
	resp.Success("成功").Send()
}

// 取得買方和收銀機資訊
func GetBuyerAndStoreDataAction(ctx *gin.Context) {
	resp := response.New(ctx)
	UserData := middleware.GetUserData(ctx)
	if len(UserData.Uid) == 0 {
		resp.Fail(1001001, "尚未登入").Send()
		return
	}
	orderId := ctx.Param("orderId")
	storeData := middleware.GetStoreData(ctx)
	data, err := OrderMessageBoard.GetBuyerAndStoreData(orderId, UserData, storeData)
	if err != nil {
		resp.Fail(errorMessage.GetMessageByCode(1001002)).Send()
	}
	resp.Success("成功").SetData(data).Send()
}

// 取得訂單留言板訊息
func GetOrderMessageBoardAction(ctx *gin.Context) {
	resp := response.New(ctx)
	orderId := ctx.Param("orderId")
	UserData := middleware.GetUserData(ctx)
	storeData := middleware.GetStoreData(ctx)
	message, err := OrderMessageBoard.GetMessage(UserData, storeData, orderId)
	if err != nil {
		log.Debug("GetOrderMessageBoardAction orderId : %s", orderId)
		resp.Fail(errorMessage.GetMessageByCodeString(err.Error())).Send()
		return
	}
	resp.Success("成功").SetData(message).Send()
}

// 取得賣家訂單留言板列表
func GetSellerOrderMessageBoardListAction(ctx *gin.Context) {
	resp := response.New(ctx)
	params := Request.OrderMessageBoardRequest{}
	if err := ctx.BindQuery(&params); err != nil {
		log.Error("get order message board params", err)
		resp.Fail(errorMessage.GetMessageByCode(1001002)).Send()
		return
	}
	storeData := middleware.GetStoreData(ctx)
	userData := middleware.GetUserData(ctx)
	data, err := OrderMessageBoard.GetMessageBoardList(userData, storeData, params, Enum.MemberSeller)
	if err != nil {
		resp.Fail(1001001, err.Error()).Send()
		return
	}
	resp.Success("成功").SetData(data).Send()
}

// 取得買家訂單留言板列表
func GetBuyerOrderMessageBoardListAction(ctx *gin.Context) {
	resp := response.New(ctx)
	params := Request.OrderMessageBoardRequest{}
	if err := ctx.BindQuery(&params); err != nil {
		log.Error("get order message board params", err)
		resp.Fail(errorMessage.GetMessageByCode(1001002)).Send()
		return
	}
	storeData := middleware.GetStoreData(ctx)
	userData := middleware.GetUserData(ctx)
	data, err := OrderMessageBoard.GetMessageBoardList(userData, storeData, params, Enum.MemberBuyer)
	if err != nil {
		resp.Fail(1001001, err.Error()).Send()
		return
	}
	resp.Success("成功").SetData(data).Send()
}
