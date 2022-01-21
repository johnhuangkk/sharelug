package controllers

import (
	"api/config/middleware"
	"api/services/Enum"
	"api/services/Service/Shipment"
	"api/services/VO/ShipmentVO"
	"api/services/errorMessage"
	"api/services/util/log"
	"api/services/util/response"

	"github.com/gin-gonic/gin"
)

// 取號
func GetShipNumberAction(ctx *gin.Context) {
	resp := response.New(ctx)
	orders := ShipmentVO.Orders{}
	if err := ctx.ShouldBind(&orders); err != nil {
		log.Error("GetShipNumberAction order params error", err)
		resp.Fail(123456, "欄位未填寫完整").Send()
		return
	}
	userData := middleware.GetUserData(ctx)
	storeData := middleware.GetStoreData(ctx)
	err, code := Shipment.GetShipNumber(orders, userData.Uid, storeData.StoreId)
	if err != nil {
		resp.Fail(code, err.Error()).Send()
	} else {
		resp.Success("OK").Send()
	}
}

// 查詢貨態資料
func GetShipStatusAction(ctx *gin.Context) {
	resp := response.New(ctx)
	params := ShipmentVO.Order{}
	if err := ctx.ShouldBindJSON(&params); err != nil {
		log.Debug("getIPostMailShipStatus => ", err)
		resp.Fail(errorMessage.GetMessageByCode(1001002)).Send()
		return
	}
	log.Debug("params", params)
	userData := middleware.GetUserData(ctx)
	storeData := middleware.GetStoreData(ctx)
	data, err := Shipment.SearchSippingStatus(params, userData.Uid, storeData.StoreId)
	if err != nil {
		resp.Fail(12332, err.Error()).Send()
	}
	log.Debug("GetShipStatusAction", data)
	resp.Success("OK").SetData(data).Send()
}

// 取托運單
func GetConsignmentAction(ctx *gin.Context) {
	resp := response.New(ctx)
	orders := ShipmentVO.Orders{}
	var ship string
	var data interface{}
	var err error
	if err := ctx.ShouldBind(&orders); err != nil {
		log.Error("GetShipNumberAction order params error", err)
		resp.Fail(200, "欄位未填寫完整").Send()
		return
	}
	storeData := middleware.GetStoreData(ctx)
	ship, data, err = Shipment.GetConsignment(orders, storeData)
	log.Info("SHIP [%s]", ship)
	if err != nil {
		resp.Fail(1111, err.Error()).Send()
		return
	}
	switch ship {
	case Enum.CVS_7_ELEVEN:
		ctx.Data(200, "application/html; charset=utf-8", data.([]byte))
	case Enum.CVS_FAMILY:
		ctx.Data(200, "application/gif; charset=utf-8", data.([]byte))
	case Enum.CVS_HI_LIFE, Enum.CVS_OK_MART:
		ctx.Data(200, "application/pdf; charset=utf-8", data.([]byte))
	default:
		resp.Success("OK").SetData(data).Send()
	}
}

// 取得超商托運資訊
func GetShipDataAction(ctx *gin.Context) {
	resp := response.New(ctx)
	var params ShipmentVO.Order
	if err := ctx.BindQuery(&params); err != nil {
		log.Error("GetShipDataAction Error [%s]", err)
		resp.Fail(200, "欄位未填寫完整").Send()
		return
	}
	log.Info(`GetShipDataAction params: [%s]`, params)
	userData := middleware.GetUserData(ctx)
	storeData := middleware.GetStoreData(ctx)
	data, err := Shipment.GetCvsShippingData(params, userData.Uid, storeData.StoreId)
	if err != nil {
		resp.Fail(1111, err.Error()).Send()
		return
	}
	resp.Success("OK").SetData(data).Send()
}

// 閉轉店變更
func PutCvsShipDataAction(ctx *gin.Context) {
	resp := response.New(ctx)
	var params = ShipmentVO.SwitchOrder{}
	if err := ctx.ShouldBindJSON(&params); err != nil {
		log.Debug("PutCvsSwitchOrderAction => ", err)
		resp.Fail(errorMessage.GetMessageByCode(1001002)).Send()
		return
	}
	userData := middleware.GetUserData(ctx)
	storeData := middleware.GetStoreData(ctx)
	if err := Shipment.SwitchOrder(params, userData.Uid, storeData.StoreId); err != nil {
		resp.Fail(12323213, err.Error()).Send()
		return
	}
	resp.Success("OK").SetData(true).Send()
}
