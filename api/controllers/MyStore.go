package controllers

import (
	"api/config/middleware"
	"api/services/VO/Request"
	"api/services/errorMessage"
	"api/services/model"
	"api/services/util/log"
	"api/services/util/response"
	"github.com/gin-gonic/gin"
	"io"
	"os"
	"time"
)

//取我的收銀機
func GetMyStoreAction(ctx *gin.Context) {
	resp := response.New(ctx)
	userData := middleware.GetUserData(ctx)
	storeData := middleware.GetStoreData(ctx)
	res, err := model.GetMyStore(storeData, userData)
	if err != nil {
		resp.Fail(1001001, err.Error()).Send()
		return
	}
	resp.Success("OK").SetData(res).Send()
}

//取我的收銀機列表
func GetMyStoreListAction(ctx *gin.Context) {
	resp := response.New(ctx)
	userData := middleware.GetUserData(ctx)
	data, err := model.GetStoreList(userData)
	if err != nil {
		resp.Fail(1001001, err.Error()).Send()
		return
	}
	resp.Success("OK").SetData(data).Send()
}

//取收銀機銷售統計
func GetSalesReportAction(ctx *gin.Context) {
	resp := response.New(ctx)
	params := &Request.SalesReportRequest{}
	if err := ctx.BindQuery(&params); err != nil {
		log.Error("get sales report params", err)
		resp.Fail(errorMessage.GetMessageByCode(1001002)).Send()
		return
	}
	storeData := middleware.GetStoreData(ctx)
	data, err := model.GetSalesStatisticsReport(storeData, params)
	if err != nil {
		resp.Fail(1001001, err.Error()).Send()
		return
	}
	resp.Success("OK").SetData(data).Send()
}

//報表明細下載
func DownLoadSalesReportAction(ctx *gin.Context) {
	resp := response.New(ctx)
	params := &Request.SalesReportRequest{}
	if err := ctx.BindQuery(&params); err != nil {
		log.Error("get sales report params", err)
		resp.Fail(errorMessage.GetMessageByCode(1001002)).Send()
		return
	}
	storeData := middleware.GetStoreData(ctx)
	filename, err := model.DownLoadSalesStatisticsReport(storeData, params)
	if err != nil {
		resp.Fail(1001001, err.Error()).Send()
		return
	}
	file, err := os.Open(filename) //Create a file
	if err != nil {
		resp.Fail(errorMessage.GetMessageByCode(1001001)).Send()
		return
	}
	defer file.Close()
	ctx.Writer.Header().Add("Content-Disposition", "attachment; filename="+time.Now().Format("20060102")+".xlsx")
	ctx.Writer.Header().Add("Content-type", "application/octet-stream")
	_, err = io.Copy(ctx.Writer, file)
	if err != nil {
		resp.Fail(errorMessage.GetMessageByCode(1001001)).Send()
		return
	}
}
//使用者操作記錄
func UserOperateRecordAction(ctx *gin.Context) {
	resp := response.New(ctx)
	userData := middleware.GetUserData(ctx)
	data, err := model.GetUserOperateRecord(userData)
	if err != nil {
		resp.Fail(errorMessage.GetMessageByCodeString(err.Error())).Send()
	}
	resp.Success("OK").SetData(data).Send()
}
