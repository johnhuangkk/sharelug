package controllers

import (
	"api/config/middleware"
	"api/services/Enum"
	"api/services/VO/Request"
	"api/services/errorMessage"
	"api/services/model"
	"api/services/util/log"
	"api/services/util/response"
	"api/services/util/tools"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/smokezl/govalidators"
)

//賣家
func GetSellerShipListAction(ctx *gin.Context) {
	resp := response.New(ctx)
	var params Request.OrderSearch
	if err := ctx.BindQuery(&params); err != nil {
		resp.Fail(200, fmt.Sprintf("%v", err)).Send()
		return
	}
	//取訂單列表
	storeData := middleware.GetStoreData(ctx)
	data, err := model.GetSearchShipData(storeData, params)
	if err != nil {
		resp.Fail(200, fmt.Sprintf("%v", err)).Send()
		return
	}
	resp.Success("OK").SetData(data).Send()
}

//出貨回填 運送方式及單號 api
func SetOrderShippingAction(ctx *gin.Context) {
	resp := response.New(ctx)
	var params Request.SetShipNumberParams
	if err := ctx.ShouldBindJSON(&params); err != nil {
		log.Debug("ShouldBindJSON Error", err.Error())
		resp.Fail(errorMessage.GetMessageByCode(1001002)).Send()
		return
	}
	validator := govalidators.New()
	if err := validator.LazyValidate(params); err != nil {
		resp.Fail(errorMessage.GetMessageByCode(1001002)).Send()
		return
	}
	storeData := middleware.GetStoreData(ctx)
	if err := model.HandleSetShipNumber(storeData, params); err != nil {
		resp.Fail(1001001, err.Error()).Send()
		return
	}
	resp.Success("OK").SetData(true).Send()
}
//宅配匯出批次出貨單
func ExportOrderShippingAction(ctx *gin.Context) {
	resp := response.New(ctx)
	var params Request.ExportOrderShippingParams
	if err := ctx.ShouldBindJSON(&params); err != nil {
		resp.Fail(errorMessage.GetMessageByCode(1001002)).Send()
		return
	}
	storeData := middleware.GetStoreData(ctx)
	var filename string
	var err error
	switch params.ShipType {
		case Enum.DELIVERY_OTHER:
			filename, err = model.HandleBatchShipExport(storeData, params)
		case Enum.F2F:
			filename, err = model.HandleBatchF2fExport(storeData, params)
		case Enum.SELF_DELIVERY:
			filename, err = model.HandleBatchSelfDeliveryExport(storeData, params)
	}
	if err != nil {
		log.Error("get Export file Error")
		resp.Fail(errorMessage.GetMessageByCode(1001001)).Send()
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
//宅配匯出批次出貨單PDF
func ExportOrderShippingPdfAction(ctx *gin.Context) {
	resp := response.New(ctx)
	var params Request.ExportOrderShippingParams
	if err := ctx.ShouldBindJSON(&params); err != nil {
		resp.Fail(errorMessage.GetMessageByCode(1001002)).Send()
		return
	}
	storeData := middleware.GetStoreData(ctx)
	var data interface{}
	var err error
	switch params.ShipType {
		case Enum.DELIVERY_OTHER:
			data, err = model.HandleBatchDeliveryPdfExport(storeData, params)
		case Enum.SELF_DELIVERY:
			data, err = model.HandleBatchSelfDeliveryPdfExport(storeData, params)
	}
	if err != nil {
		resp.Fail(errorMessage.GetMessageByCodeString(err.Error())).Send()
		return
	}
	resp.Success("OK").SetData(data).Send()
}

func ExportOrderShipSendAction(ctx *gin.Context) {
	resp := response.New(ctx)
	var params Request.ExportOrderShippingParams
	if err := ctx.ShouldBindJSON(&params); err != nil {
		resp.Fail(errorMessage.GetMessageByCode(1001002)).Send()
		return
	}
	storeData := middleware.GetStoreData(ctx)
	var data interface{}
	var err error
	switch params.ShipType {
		case Enum.DELIVERY_OTHER:
			data, err = model.HandleBatchDeliverySendExport(storeData, params)
		case Enum.F2F:
			data, err = model.HandleBatchF2fSndExport(storeData, params)
		case Enum.SELF_DELIVERY:
			data, err = model.HandleBatchSelfDeliverySendExport(storeData, params)
	}
	if err != nil {
		resp.Fail(errorMessage.GetMessageByCodeString(err.Error())).Send()
		return
	}
	resp.Success("OK").SetData(data).Send()
}
//匯入批次出貨檔
func ImportBatchShippingAction(ctx *gin.Context) {
	resp := response.New(ctx)
	file, header, err := ctx.Request.FormFile("file")
	if err != nil {
		resp.Fail(errorMessage.GetMessageByCode(1001002)).Send()
		return
	}
	filename, err := tools.UploadFile(file, header)
	if err != nil {
		resp.Fail(errorMessage.GetMessageByCode(1001009)).Send()
		return
	}
	storeData := middleware.GetStoreData(ctx)
	data, err := model.HandleBatchShipFile(storeData, filename)
	if err != nil {
		resp.Fail(errorMessage.GetMessageByCodeString(err.Error())).Send()
		return
	}
	resp.Success("OK").SetData(data).Send()
}

func ProcessBatchShippingAction(ctx *gin.Context)  {
	resp := response.New(ctx)
	batchId := ctx.Param("batchId")
	storeData := middleware.GetStoreData(ctx)
	if err := model.HandleBatchShip(storeData, batchId); err != nil {
		resp.Fail(errorMessage.GetMessageByCodeString(err.Error())).Send()
		return
	}
	resp.Success("OK").SetData(true).Send()
}
//面交匯出批次出貨單
//func ExportOrderF2fAction(ctx *gin.Context) {
//	resp := response.New(ctx)
//	var params Request.ExportOrderShippingParams
//	if err := ctx.ShouldBindJSON(&params); err != nil {
//		resp.Fail(errorMessage.GetMessageByCode(1001002)).Send()
//		return
//	}
//	storeData := middleware.GetStoreData(ctx)
//	filename, err := model.HandleBatchF2fExport(storeData, params)
//	if err != nil {
//		resp.Fail(1001001, err.Error()).Send()
//		return
//	}
//	file, err := os.Open(filename) //Create a file
//	if err != nil {
//		resp.Fail(errorMessage.GetMessageByCode(1001001)).Send()
//		return
//	}
//	defer file.Close()
//	ctx.Writer.Header().Add("Content-Disposition", "attachment; filename=checkne"+time.Now().Format("20060102150405")+".xlsx")
//	ctx.Writer.Header().Add("Content-type", "application/octet-stream")
//	_, err = io.Copy(ctx.Writer, file)
//	if err != nil {
//		resp.Fail(errorMessage.GetMessageByCode(1001001)).Send()
//		return
//	}
//}
//面交匯出批次出貨單PDF
//func ExportOrderF2fPdfAction(ctx *gin.Context) {
//	resp := response.New(ctx)
//	var params Request.ExportOrderShippingParams
//	if err := ctx.ShouldBindJSON(&params); err != nil {
//		resp.Fail(errorMessage.GetMessageByCode(1001002)).Send()
//		return
//	}
//	storeData := middleware.GetStoreData(ctx)
//	data, err := model.HandleBatchF2fExportPdf(storeData, params)
//	if err != nil {
//		resp.Fail(errorMessage.GetMessageByCodeString(err.Error())).Send()
//		return
//	}
//	resp.Success("OK").SetData(data).Send()
//}

//func ExportOrderSelfDeliveryAction(ctx *gin.Context) {
//	resp := response.New(ctx)
//	var params Request.ExportOrderShippingParams
//	if err := ctx.ShouldBindJSON(&params); err != nil {
//		resp.Fail(errorMessage.GetMessageByCode(1001002)).Send()
//		return
//	}
//	storeData := middleware.GetStoreData(ctx)
//	filename, err := model.HandleBatchSelfDeliveryExport(storeData, params)
//	if err != nil {
//		resp.Fail(1001001, err.Error()).Send()
//		return
//	}
//	file, err := os.Open(filename) //Create a file
//	if err != nil {
//		resp.Fail(errorMessage.GetMessageByCode(1001001)).Send()
//		return
//	}
//	defer file.Close()
//	ctx.Writer.Header().Add("Content-Disposition", "attachment; filename=checkne"+time.Now().Format("20060102150405")+".xlsx")
//	ctx.Writer.Header().Add("Content-type", "application/octet-stream")
//	_, err = io.Copy(ctx.Writer, file)
//	if err != nil {
//		resp.Fail(errorMessage.GetMessageByCode(1001001)).Send()
//		return
//	}
//}
//面交匯出批次出貨單PDF
//func ExportOrderSelfDeliveryPdfAction(ctx *gin.Context) {
//	resp := response.New(ctx)
//	var params Request.ExportOrderShippingParams
//	if err := ctx.ShouldBindJSON(&params); err != nil {
//		resp.Fail(errorMessage.GetMessageByCode(1001002)).Send()
//		return
//	}
//	storeData := middleware.GetStoreData(ctx)
//	data, err := model.HandleBatchSelfDeliveryExportPdf(storeData, params)
//	if err != nil {
//		resp.Fail(errorMessage.GetMessageByCodeString(err.Error())).Send()
//		return
//	}
//	resp.Success("OK").SetData(data).Send()
//}
