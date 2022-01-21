package Product

import (
	"api/services/Enum"
	"api/services/Service/Carts"
	"api/services/Service/StoreService"
	"api/services/VO/CartsVo"
	"api/services/VO/Response"
	"api/services/dao/Short"
	"api/services/dao/product"
	"api/services/dao/sequence"
	"api/services/database"
	"api/services/entity"
	"api/services/util/images"
	"api/services/util/log"
	"api/services/util/tools"
	"encoding/base64"
	"encoding/json"
	"fmt"
)

func GetProductDataByProductSpecId(engine *database.MysqlSession, productSpecId string) (entity.ProductsData, error) {
	productData, err := product.GetProductByProductSpecId(engine, productSpecId)
	if err != nil {
		return entity.ProductsData{}, err
	}
	return productData, nil
}

//檢查庫存數
func CheckProductSpecQty(engine *database.MysqlSession, productSpecId string, Quantity int) (entity.ProductsData, error) {
	data, err := product.GetProductByProductSpecId(engine, productSpecId)
	if err != nil {
		return data, fmt.Errorf("1001010")
	}
	if int64(Quantity) > data.Quantity {
		return data, fmt.Errorf("1009003")
	}
	return data, nil
}

func ModifyProducts(engine *database.MysqlSession, cookie string, carts CartsVo.Carts) (Response.CartsProductModel, string, []string, []string, float64, []int, bool, error) {
	merge := make([]int, 0)
	var shipList []string
	var payWayList []string
	var subtotal float64
	var storeId string
	var isCoupon bool
	storeId = carts.StoreId
	var resp Response.CartsProductModel
	//取出免運設定
	setting, err := StoreService.GetStoreFreeShipping(engine, carts.StoreId)
	if err != nil {
		return resp, storeId, shipList, payWayList, subtotal, merge, isCoupon, err
	}
	data, err := takeCartsProduct(engine, setting, cookie, carts)
	if err != nil {
		return data, storeId, shipList, payWayList, subtotal, merge, isCoupon, err
	}
	//不可合拼計算
	for _, v := range data.Merge {
		subtotal += float64(v.Price * v.Quantity)
		merge = append(merge, v.ShipMerge)
		shipList = append(shipList, v.ShipList)
		payWayList = append(payWayList, v.PayWayList)
		resp.MergeFee += v.ShipFee
	}
	isCoupon = setting.IsCoupon
	resp.Merge = append(resp.Merge, data.Merge...)
	//免運計算
	amount := float64(0)
	qty := 0
	sub := float64(0)
	for _, v := range data.Free {
		sub += float64(v.Price * v.Quantity)
		merge = append(merge, v.ShipMerge)
		shipList = append(shipList, v.ShipList)
		payWayList = append(payWayList, v.PayWayList)
		amount += float64(v.Price * v.Quantity)
		qty += v.Quantity
	}
	//外送免運 選擇外送才有
	if setting.SelfDelivery && carts.Shipping == Enum.SELF_DELIVERY {
		//判斷是否達免運條件
		if setting.SelfDeliveryKey == Enum.FreeShipAmount && int64(amount) >= setting.SelfDeliveryFree {
			subtotal += sub
			resp.Free = append(resp.Free, data.Free...)
		} else if setting.SelfDeliveryKey == Enum.FreeShipQuantity && int64(qty) >= setting.SelfDeliveryFree {
			subtotal += sub
			resp.Free = append(resp.Free, data.Free...)
		} else {
			//未達免運丟進可合拼
			data.General = append(data.General, data.Free...)
		}
	} else {
		//判斷是否達免運條件
		if setting.FreeShipKey == Enum.FreeShipAmount && int64(amount) >= setting.FreeShip {
			subtotal += sub
			resp.Free = append(resp.Free, data.Free...)
		} else if setting.FreeShipKey == Enum.FreeShipQuantity && int64(qty) >= setting.FreeShip {
			subtotal += sub
			resp.Free = append(resp.Free, data.Free...)
		} else {
			//未達免運丟進可合拼
			data.General = append(data.General, data.Free...)
		}
	}
	//可合拼計算
	for _, v := range data.General {
		subtotal += float64(v.Price * v.Quantity)
		merge = append(merge, v.ShipMerge)
		shipList = append(shipList, v.ShipList)
		payWayList = append(payWayList, v.PayWayList)
		if v.ShipFee > resp.GeneralFee {
			resp.GeneralFee = v.ShipFee
		}
	}
	resp.General = append(resp.General, data.General...)
	return resp, storeId, shipList, payWayList, subtotal, merge, isCoupon, err
}

func takeCartsProduct(engine *database.MysqlSession, setting Response.StoreFreeShipResponse, cookie string, carts CartsVo.Carts) (Response.CartsProductModel, error) {
	var resp Response.CartsProductModel
	for k, v := range carts.Products {
		data, err := product.GetProductByProductSpecId(engine, v.ProductSpecId)
		if err != nil {
			return resp, err
		}
		image, err := product.GetProductFirstImageByProductId(engine, data.ProductId)
		if err != nil {
			return resp, err
		}
		var rep Response.CartsProductList
		rep.ProductSpecId = data.SpecId
		rep.ProductId = data.ProductId
		rep.ProductName = data.ProductName
		rep.ProductSpec = data.SpecName
		rep.ProductImage = images.GetImageUrl(image.Image)
		rep.Price = int(data.Price)
		rep.ShipMerge = data.ShipMerge
		rep.ShipFee = takeProductShipFee(data.Shipping, carts.Shipping)
		rep.ShipMode = rearrangeShipping(data.Shipping, setting)
		rep.PayWayMode = rearrangePayWay(data.PayWay)
		rep.FormUrl = data.FormUrl
		rep.ShipList = rearrangeShipList(data.Shipping, setting)
		rep.PayWayList = data.PayWay
		rep.LimitKey = data.LimitKey
		rep.LimitQty = int(data.LimitQty)
		//判斷庫存數
		if data.Quantity >= int64(v.Quantity) {
			rep.Quantity = v.Quantity
		} else {
			rep.Quantity = int(data.Quantity)
			carts.Products[k].Quantity = int(data.Quantity)
			err := Carts.SetRedisCarts(carts, cookie, carts.Style)
			if err != nil {
				return resp, err
			}
		}
		//外送免運設定
		if setting.FreeShipKey != Enum.FreeShipNone && data.IsFreeShip == true && data.ShipMerge == 1 || setting.SelfDelivery && data.ShipMerge == 1{
			if setting.SelfDelivery {
				for _, ship := range rep.ShipMode {
					if ship.Type == Enum.SELF_DELIVERY {
						rep.SelfDeliveryKey = setting.SelfDeliveryKey
						rep.SelfDeliveryFree = setting.SelfDeliveryFree
					}
				}
			} else {
				rep.SelfDeliveryKey = Enum.FreeShipNone
				rep.FreeShip = 0
			}
			rep.FreeShipKey = setting.FreeShipKey
			rep.FreeShip = setting.FreeShip
			resp.Free = append(resp.Free, rep)
		} else if data.ShipMerge == 0 {
			resp.Merge = append(resp.Merge, rep)
		} else {
			resp.General = append(resp.General, rep)
		}
	}
	return resp, nil
}

//重整付款資訊
func rearrangePayWay(payWayList string) []Response.PayWayMode {
	var resp []Response.PayWayMode
	if err := json.Unmarshal([]byte(payWayList), &resp); err != nil {
		log.Error("json Unmarshal Error", err)
	}
	return resp
}

//重整運費資訊
func rearrangeShipping(shipList string, setting Response.StoreFreeShipResponse) []Response.ShipMode {
	var ship []Response.ShippingMode
	if err := json.Unmarshal([]byte(shipList), &ship); err != nil {
		log.Error("json Unmarshal Error", err)
	}
	var resp []Response.ShipMode
	for _, v := range ship {
		if v.Type == Enum.SELF_DELIVERY && setting.SelfDelivery || v.Type != Enum.SELF_DELIVERY {
			res := Response.ShipMode{
				Type: v.Type,
				Text: v.Text,
			}
			resp = append(resp, res)
		}
	}
	return resp
}

func rearrangeShipList(shipList string, setting Response.StoreFreeShipResponse) string {
	var ship []Response.ShippingMode
	if err := json.Unmarshal([]byte(shipList), &ship); err != nil {
		log.Error("json Unmarshal Error", err)
	}
	var resp []Response.ShippingMode
	for _, v := range ship {
		if v.Type == Enum.SELF_DELIVERY && setting.SelfDelivery || v.Type != Enum.SELF_DELIVERY {
			resp = append(resp, v)
		}
	}
	s, _ := tools.JsonEncode(resp)
	return s
}
// 取得商品運費
func takeProductShipFee(shipList string, choose string) int {
	var fee int
	var ship []Response.ShippingMode
	if err := tools.JsonDecode([]byte(shipList), &ship); err != nil {
		log.Error("json decode error", err)
	}
	for _, value := range ship {
		if value.Type == choose {
			fee = value.Price
		}
	}
	return fee
}
//產生短網址
func GeneratorShortUrl(engine *database.MysqlSession, url string) (string, error) {
	number, _ := sequence.GetTinyUrlVirtualSeq()
	number = tools.StringPadLeft(number, 6)
	tiny := base64.StdEncoding.EncodeToString([]byte(number))
	//寫入資料庫
	var data entity.ShortUrlData
	data.Short = tiny
	data.Url = url
	if err := Short.InsertShortUrlData(engine, data); err != nil {
		return tiny, err
	}
	return tiny, nil
}

func ModifyShipList(shipList string, setting Response.StoreFreeShipResponse) []Response.ShippingMode {
	var ship []Response.ShippingMode
	if err := json.Unmarshal([]byte(shipList), &ship); err != nil {
		log.Error("json Unmarshal Error", err)
	}
	var resp []Response.ShippingMode
	for _, v := range ship {
		if v.Type == Enum.SELF_DELIVERY && setting.SelfDelivery || v.Type != Enum.SELF_DELIVERY {
			resp = append(resp, v)
		}
	}
	return resp
}
