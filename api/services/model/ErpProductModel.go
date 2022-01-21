package model

import (
	"api/services/Enum"
	"api/services/VO/Request"
	"api/services/VO/Response"
	"api/services/dao/Orders"
	"api/services/dao/Store"
	"api/services/dao/member"
	"api/services/dao/product"
	"api/services/database"
	"api/services/util/log"
	"api/services/util/tools"
	"fmt"
	"strings"
)

func HandleProductList(params Request.ErpSearchProductRequest)  ([]Response.SearchProductResponse, error) {
	engine := database.GetMysqlEngine()
	defer engine.Close()
	var resp []Response.SearchProductResponse
	switch params.Tab {
		case "Product":
			data, err := product.SearchProducts(engine, params)
			if err != nil {
				log.Debug("Get Order Error", err)
				return resp, fmt.Errorf("1001001")
			}

			for _, v := range data {
				store, err := Store.GetStoreDataByStoreId(engine, v.StoreId)
				if err != nil {
					log.Debug("Get Store data Error", err)
				}
				seller, err := member.GetMemberDataByUid(engine, store.SellerId)
				if err != nil {
					log.Debug("Get Member data Error", err)
				}
				res := v.GetSearchProduct(seller, store)
				if v.IsSpec != 0 {
					res.ProductSpec = getProductSpecToString(engine, v.ProductId)
				}
				res.ShipMode = getProductShipToString(v.Shipping)
				res.PayWayMode = getProductPayWayToString(v.PayWay)
				resp = append(resp, res)
			}
		case "Realtime":
			data, err := product.SearchProducts(engine, params)
			if err != nil {
				log.Debug("Get Order Error", err)
				return resp, fmt.Errorf("1001001")
			}
			for _, v := range data {
				store, err := Store.GetStoreDataByStoreId(engine, v.StoreId)
				if err != nil {
					log.Debug("Get Store data Error", err)
				}
				seller, err := member.GetMemberDataByUid(engine, store.SellerId)
				if err != nil {
					log.Debug("Get Member data Error", err)
				}
				res := v.GetRealtimeProduct(seller, store)
				res.ProductSpec = getProductSpecToString(engine, v.ProductId)
				res.ShipMode = getProductShipToString(v.Shipping)
				res.PayWayMode = getProductPayWayToString(v.PayWay)
				resp = append(resp, res)
			}
		case "Bill":
			data, err := Orders.SearchBills(engine, params)
			if err != nil {
				log.Debug("Get Order Error", err)
				return resp, fmt.Errorf("1001001")
			}
			for _, v := range data {
				user, err := member.GetMemberDataByUid(engine, v.BuyerId)
				if err != nil {
					log.Debug("Get Member data Error", err)
				}
				res := v.GetSearchBill(user)
				res.ReceiverAddress = GetShipAddress(engine, v.ShipType, v.ReceiverAddress)
				if !tools.InArray([]string{Enum.I_POST, Enum.CVS_FAMILY, Enum.CVS_HI_LIFE, Enum.CVS_OK_MART, Enum.CVS_7_ELEVEN}, v.ShipType) {
					res.ReceiverAddress = tools.MaskerAddress(res.ReceiverAddress)
				}
				resp = append(resp, res)
			}
	}

	return resp, nil
}

func getProductSpecToString(engine *database.MysqlSession, productId string) string {
	var resp []string
	spec, err := product.GetProductSpecByProductId(engine, productId)
	if err != nil {
		log.Debug("Get Product spec data Error", err)
	}
	for _, v := range spec {
		resp = append(resp, v.SpecName)
	}
	return strings.Join(resp, ",")
}

func getProductShipToString(shipping string) string {
	var resp []string
	var productShipList []Response.ShippingMode
	if err := tools.JsonDecode([]byte(shipping), &productShipList); err != nil {
		log.Error("Json Decode Error", err)
	}
	for _, v := range productShipList {
		resp = append(resp, v.Text)
	}
	return strings.Join(resp, ",")
}

func getProductPayWayToString(payWay string) string {
	var resp []string
	var productPayWayList []Response.PayWayMode
	if err := tools.JsonDecode([]byte(payWay), &productPayWayList); err != nil {
		log.Error("Json Decode Error", err)
	}
	for _, v := range productPayWayList {
		resp = append(resp, v.Text)
	}
	return strings.Join(resp, ",")
}
