package History

import (
	"api/services/Enum"
	"api/services/VO/Request"
	"api/services/VO/Response"
	"api/services/dao/History"
	"api/services/database"
	"api/services/entity"
	"api/services/util/log"
	"api/services/util/tools"
	"encoding/json"
	"strconv"
	"time"
)

func GenerateStatusLog(table string, field string, dataId string, oldValue string, newValue string, userName string) error {
	engine := database.GetMysqlEngine()
	defer engine.Close()
	var UserName string
	if len(userName) == 0 {
		UserName = "SYSTEM"
	} else {
		UserName = userName
	}

	var Entity entity.StatusHistoryLog
	Entity.Table = table
	Entity.Field = field
	Entity.DataId = dataId
	Entity.OldValue = oldValue
	Entity.NewValue = newValue
	Entity.OperateUserId = UserName
	Entity.CreateTime = time.Now()

	_, err := History.InsertHistoryLog(engine, Entity)
	if err != nil {
		log.Error("Insert Status history Log Database Error", err)
		return err
	}
	return nil
}

func GenerateNewProductLog(newData entity.ProductData, userName string) {
	engine := database.GetMysqlEngine()
	defer engine.Close()
	var UserName string
	if len(userName) == 0 {
		UserName = "SYSTEM"
	} else {
		UserName = userName
	}
	var productLog entity.ProductHistoryLog
	productLog.NewValue = newData
	productLog.Action = "create"
	productLog.OperateUserId = UserName
	productLog.ProductId = newData.ProductId

	err := History.InsertProductLog(engine, productLog)
	if err != nil {
		log.Error("Insert Status history Log Database Error", err)
	}
}
func GenerateEditProductLog(params *Request.EditProductParams, oldData entity.ProductData, userName string) {
	engine := database.GetMysqlEngine()
	defer engine.Close()
	var UserName string

	if len(userName) == 0 {
		UserName = "SYSTEM"
	} else {
		UserName = userName
	}

	newProduct := buildProduct(params, oldData.StoreId, oldData.ProductId, oldData.ProductStatus)
	var productLog entity.ProductHistoryLog
	productLog.NewValue = newProduct
	productLog.OperateUserId = UserName
	productLog.ProductId = oldData.ProductId
	productLog.OldValue = oldData
	if params.StatusDown == 1 {
		productLog.Action = "down"
	} else if params.StatusDelete == 1 {
		productLog.Action = "delete"
	} else {
		if oldData.ProductStatus == Enum.ProductStatusDown && params.StatusDown == 0 {
			productLog.Action = "reopen"
		} else {
			productLog.Action = "edit"
		}
	}
	err := History.InsertProductLog(engine, productLog)
	if err != nil {
		log.Error("Insert Status history Log Database Error", err)
	}
}
func GenerateBatchDownProductLog(engine *database.MysqlSession, old, new entity.ProductData, userName string) {
	var productLog entity.ProductHistoryLog
	productLog.NewValue = new
	productLog.OperateUserId = userName
	productLog.ProductId = old.ProductId
	productLog.OldValue = old
	if new.ProductStatus == Enum.ProductStatusSuccess {
		productLog.Action = "batch-open"
	} else if new.ProductStatus == Enum.ProductStatusDown {
		productLog.Action = "batch-down"
	} else {
		productLog.Action = "batch-pending"
	}

	err := History.InsertProductLog(engine, productLog)
	if err != nil {
		log.Error("Insert Status history Log Database Error", err)
	}

}
func GenerateBatchPayWayProductLog(engine *database.MysqlSession, old, new entity.ProductData, userName string) {
	var productLog entity.ProductHistoryLog
	productLog.NewValue = new
	productLog.OperateUserId = userName
	productLog.ProductId = old.ProductId
	productLog.OldValue = old
	productLog.Action = "batch-payway"

	err := History.InsertProductLog(engine, productLog)
	if err != nil {
		log.Error("Insert Status history Log Database Error", err)
	}

}
func GenerateBatchShipWayProductLog(engine *database.MysqlSession, old, new entity.ProductData, userName string) {
	var productLog entity.ProductHistoryLog
	productLog.NewValue = new
	productLog.OperateUserId = userName
	productLog.ProductId = old.ProductId
	productLog.OldValue = old
	productLog.Action = "batch-shipway"

	err := History.InsertProductLog(engine, productLog)
	if err != nil {
		log.Error("Insert Status history Log Database Error", err)
	}

}
func buildProduct(params *Request.EditProductParams, storeId, productId, status string) entity.ProductData {
	p, _ := strconv.Atoi(params.Price)
	payWay, _ := json.Marshal(rearrangesPayWayMode(params.PayWayList, params.ShippingList))
	shipping, _ := json.Marshal(rearrangeShipMode(params.ShippingList))
	newData := entity.ProductData{
		ProductId:   productId,
		Stock:       int64(params.ProductQty),
		Price:       int64(p),
		FormUrl:     params.FormUrl,
		ShipMerge:   params.ShipMerge,
		IsSpec:      params.IsSpec,
		IsRealtime:  params.IsRealtime,
		ProductName: params.ProductName,
		LimitKey:    params.LimitKey,
		LimitQty:    int64(params.LimitQty),
		PayWay:      string(payWay),
		Shipping:    string(shipping),
		StoreId:     storeId,
	}
	if params.StatusDown == 1 {
		newData.ProductStatus = Enum.ProductStatusDown
	} else if params.StatusDelete == 1 {
		newData.ProductStatus = Enum.ProductStatusDelete
	} else {
		newData.ProductStatus = status
	}
	return newData
}
func rearrangeShipMode(shipping []Request.NewShipping) []entity.NewShippingMode {
	var ShippingMode []entity.NewShippingMode
	for _, val := range shipping {
		var data entity.NewShippingMode
		data.Type = val.ShipType
		data.Text = Enum.Shipping[val.ShipType]
		data.Price = val.ShipFee
		data.Remark = val.ShipRemark
		ShippingMode = append(ShippingMode, data)
	}
	return ShippingMode
}
func rearrangesPayWayMode(payway []string, shipping []Request.NewShipping) []Response.PayWayMode {
	var paywayMode []Response.PayWayMode
	for _, val := range payway {
		if val == Enum.CvsPay {
			if checkShipIsCvs(shipping) {
				mode := Response.PayWayMode{
					Type: val,
					Text: Enum.PayWay[val],
				}
				paywayMode = append(paywayMode, mode)
			}
		} else {
			mode := Response.PayWayMode{
				Type: val,
				Text: Enum.PayWay[val],
			}
			paywayMode = append(paywayMode, mode)
		}
	}
	return paywayMode
}

func checkShipIsCvs(shipping []Request.NewShipping) bool {
	for _, v := range shipping {
		if tools.InArray([]string{Enum.CVS_HI_LIFE, Enum.CVS_FAMILY, Enum.CVS_OK_MART, Enum.CVS_7_ELEVEN}, v.ShipType) {
			return true
		}
	}
	return false
}
