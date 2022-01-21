package controllers

import (
	"api/config/middleware"
	"api/services/Service/Erp"
	"api/services/Task"
	"api/services/VO/Request"
	"api/services/errorMessage"
	"api/services/model"
	"api/services/util/log"
	"api/services/util/response"
	"io"
	"os"

	"github.com/gin-gonic/gin"
)

// 取得登入 使用者資訊＆店舖資訊
func GetMemberDataAction(ctx *gin.Context) {
	resp := response.New(ctx)
	userData := middleware.GetUserData(ctx)
	storeData := middleware.GetStoreData(ctx)
	data, err := model.HandleGetMemberInfo(userData, storeData)
	if err != nil {
		resp.Fail(errorMessage.GetMessageByCodeString(err.Error())).Send()
		return
	}
	resp.Success("成功").SetData(data).Send()
}

// 驗證使用者EMAIL
func CheckMemberEmailVerifyAction(ctx *gin.Context) {
	resp := response.New(ctx)
	userData := middleware.GetUserData(ctx)
	if len(userData.Email) != 0 {
		resp.Success("OK").SetData(true).Send()
	} else {
		resp.Success("OK").SetData(false).Send()
	}
}

// 取得驗證使用者EMAIL
func GetEmailVerifyAction(ctx *gin.Context) {
	resp := response.New(ctx)
	var params Request.MemberEmailVerifyRequest
	err := ctx.BindQuery(&params)
	if err != nil {
		log.Error("post params error", err)
		resp.Fail(errorMessage.GetMessageByCode(1001002)).Send()
		return
	}
	data, err := model.HandleGetEmailVerify(params)
	if err != nil {
		resp.Fail(errorMessage.GetMessageByCodeString(err.Error())).Send()
		return
	}
	resp.Success("OK").SetData(data).Send()
}

//驗證會員EMAIL
func MemberEmailVerifyAction(ctx *gin.Context) {
	resp := response.New(ctx)
	var params Request.MemberEmailVerifyRequest
	err := ctx.ShouldBindJSON(&params)
	if err != nil {
		log.Error("post params error", err)
		resp.Fail(errorMessage.GetMessageByCode(1001002)).Send()
		return
	}
	userData := middleware.GetUserData(ctx)
	data, err := model.HandleMemberEmailVerify(userData, params)
	if err != nil {
		resp.Fail(errorMessage.GetMessageByCodeString(err.Error())).Send()
		return
	}
	resp.Success("OK").SetData(data).Send()
}

func MemberCompanyVerify(ctx *gin.Context) {
	resp := response.New(ctx)
	var params Request.MemberCompanyVerifyRequest
	err := ctx.ShouldBindJSON(&params)
	if err != nil {
		log.Error("post params error", err)
		resp.Fail(errorMessage.GetMessageByCode(1001002)).Send()
		return
	}
	data, err := model.HandleMemberCompany(params)
	if err != nil {
		resp.Fail(errorMessage.GetMessageByCodeString(err.Error())).Send()
		return
	}
	resp.Success("OK").SetData(data).Send()
}

func MemberCompanyVerifyPendingList(ctx *gin.Context) {
	resp := response.New(ctx)

	data, err := model.HandleCompanyPendingList()
	if err != nil {
		resp.Fail(errorMessage.GetMessageByCodeString(err.Error())).Send()
		return
	}
	resp.Success("OK").SetData(data).Send()
}

func MemberSpecialStoreVerifyInfo(ctx *gin.Context) {
	resp := response.New(ctx)
	memberId := ctx.Param("memberId")

	data, err := model.HandleMemberCompanyInfo(memberId)
	if err != nil {
		resp.Fail(errorMessage.GetMessageByCodeString(err.Error())).Send()
		return
	}
	resp.Success("OK").SetData(data).Send()
}

func UpdateMemberCompany(ctx *gin.Context) {
	resp := response.New(ctx)
	memberId := ctx.Param("memberId")
	var params Request.MemberPersonalRequest
	if err := ctx.ShouldBindJSON(&params); err != nil {
		resp.Fail(errorMessage.GetMessageByCode(1001001)).Send()
		return
	}
	err := model.UpdateMemberCompanyInfo(memberId, params)
	if err != nil {
		resp.Fail(errorMessage.GetMessageByCode(1001001)).Send()
		return
	}
	resp.Success("OK").Send()

}

func MemberSpecialStoreVerifyList(ctx *gin.Context) {
	resp := response.New(ctx)

	data, err := model.GetMemberSpecialStoreVerifyList()
	if err != nil {
		resp.Fail(errorMessage.GetMessageByCode(1001001)).Send()
		return
	}
	resp.Success("OK").SetData(data).Send()
}

func FindMemberSpecialStoreVerifyRecord(ctx *gin.Context) {
	resp := response.New(ctx)
	var params Request.MemberPersonalSpecialStoreRequest
	if err := ctx.ShouldBindJSON(&params); err != nil {
		resp.Fail(errorMessage.GetMessageByCode(1001002)).Send()
		return
	}
	data, err := model.GetMemberSpecialStoreRecords(params)
	if err != nil {
		resp.Fail(errorMessage.GetMessageByCode(1001001)).Send()
		return
	}
	resp.Success("OK").SetData(data).Send()
}

func GetMemberSendToBank(ctx *gin.Context) {
	resp := response.New(ctx)
	Task.UploadKgiSpecialStoreRecord()
	resp.Success("OK").Send()
}

func GetMemberSendToKgiFtp(ctx *gin.Context) {
	resp := response.New(ctx)
	data, err := Erp.GetFolderList()
	if err != nil {
		resp.Fail(errorMessage.GetMessageByCode(1001001)).Send()
		return
	}
	resp.Success("OK").SetData(data).Send()
}

func GetKgiSpecialStoreExcel(ctx *gin.Context) {
	resp := response.New(ctx)
	var params Request.KgiSpecialStoreExcel
	if err := ctx.ShouldBindJSON(&params); err != nil {
		log.Error(err.Error())
		resp.Fail(errorMessage.GetMessageByCode(1001001)).Send()
		return
	}

	file, err := Erp.GetSpecialStoreFile(params.Filename)
	if err != nil {
		log.Error(err.Error())
		resp.Fail(errorMessage.GetMessageByCode(1001001)).Send()
		return
	}
	file, err = os.Open(file.Name()) //Create a file
	if err != nil {
		resp.Fail(errorMessage.GetMessageByCode(1001001)).Send()
		return
	}
	defer file.Close()
	ctx.Writer.Header().Add("Content-Disposition", "attachment; filename="+params.Filename)
	ctx.Writer.Header().Add("Content-type", "application/octet-stream")
	_, err = io.Copy(ctx.Writer, file)
	if err != nil {
		log.Error(err.Error())
		resp.Fail(errorMessage.GetMessageByCode(1001001)).Send()
		return
	}
}
