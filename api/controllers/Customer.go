package controllers

import (
	"api/config/middleware"
	"api/services/Enum"
	"api/services/VO/Request"
	"api/services/errorMessage"
	"api/services/model"
	"api/services/util/response"

	"github.com/gin-gonic/gin"
)

//取出聯絡客服問題選單
func GetCustomerAction(ctx *gin.Context) {
	resp := response.New(ctx)
	data, err := model.HandleCustomerQuestion()
	if err != nil {
		resp.Fail(errorMessage.GetMessageByCode(1001001)).Send()
		return
	}
	resp.Success("OK").SetData(data).Send()
}

//發送聯絡客服
func PostCustomerAction(ctx *gin.Context) {
	resp := response.New(ctx)
	params := Request.CustomerRequest{}
	if err := ctx.ShouldBindJSON(&params); err != nil {
		resp.Fail(errorMessage.GetMessageByCode(1001002)).Send()
		return
	}
	userData := middleware.GetUserData(ctx)
	if err := model.HandleSendCustomer(userData, params); err != nil {
		resp.Fail(1001001, err.Error()).Send()
		return
	}
	resp.Success("OK").SetData(true).Send()
}

//發送聯絡我們
func PostContactAction(ctx *gin.Context) {
	resp := response.New(ctx)
	params := Request.ContactRequest{}
	if err := ctx.ShouldBindJSON(&params); err != nil {
		resp.Fail(errorMessage.GetMessageByCode(1001002)).Send()
		return
	}
	err := model.HandleContactData(params)
	if err != nil {
		resp.Fail(1001001, err.Error()).Send()
		return
	}
	resp.Success("OK").SetData(true).Send()
}

func GetCustomerQuestion(ctx *gin.Context) {
	resp := response.New(ctx)
	questionId := ctx.Param("questionId")
	types := ctx.Query("type")
	if len(questionId) == 0 || len(types) == 0 {
		if types != Enum.NotifyTypePlatformUser || types != Enum.NotifyTypePlatformService {
			resp.Fail(errorMessage.GetMessageByCode(1001002)).Send()
			return
		}

	}
	data, err := model.GetCustomerQuestionByMsgType(questionId, types)
	if err != nil {
		resp.Fail(errorMessage.GetMessageByCode(1001002)).Send()
		return
	}
	resp.Success("OK").SetData(data).Send()
}
