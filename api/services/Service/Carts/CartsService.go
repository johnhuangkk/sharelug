package Carts

import (
	"api/services/Enum"
	"api/services/VO/CartsVo"
	"api/services/VO/Request"
	"api/services/entity"
	"api/services/util/log"
	"api/services/util/redis"
	"api/services/util/session"
	"encoding/json"
	"time"
)

//新的購物車
func NewRedisCarts(data *CartsVo.Carts, cookie, style string) error {
	data.Style = style
	jsonData, _ := json.Marshal(data)
	val, _ := redis.New().GetHashRedis(cookie, style)
	if err := redis.New().SetHashRedis(cookie, style, string(jsonData)); err != nil {
		log.Error("b")
		return err
	}
	if len(val) == 0 {
		//設定購物車時間
		var min = 60 * time.Minute
		if err := redis.New().SetExpireRedis(cookie, min); err != nil {
			log.Error("c")
			return err
		}
	}
	return nil
}

//修改購物車入容
func SetRedisCarts(data CartsVo.Carts, cookie, style string) error {
	jsonData, _ := json.Marshal(data)
	if err := redis.New().SetHashRedis(cookie, style, string(jsonData)); err != nil {
		log.Error("Set Cart For Redis Error !!", err)
		return err
	}
	return nil
}

/**
 * 刪除購物車
 */
func DeleteRedisCarts(cookie, style string) error {
	if err := redis.New().DelHashRedis(cookie, style); err != nil {
		log.Debug("delete hash redis => ", err)
		return err
	}
	return nil
}

func GetCarts(cookie string) (CartsVo.Carts, error) {
	var resp CartsVo.Carts
	carts, err := GetRedisCarts(cookie, Enum.StyleCarts)
	if err != nil {
		log.Debug("Get redis carts Error", err)
		return resp, err
	}
	bill, err := GetRedisCarts(cookie, Enum.StyleBill)
	if err != nil {
		log.Debug("Get redis carts Error", err)
		return resp, err
	}
	if carts.Products == nil && bill.Products == nil {
		return resp, err
	} else if bill.Products == nil {
		return carts, nil
	} else {
		return bill, nil
	}
}
//取出Redis購物車內容
func GetRedisCarts(cookie, style string) (CartsVo.Carts, error) {
	var data CartsVo.Carts
	val, err := redis.New().GetHashRedis(cookie, style)
	if err != nil {
		log.Debug("Get hash redis Error ", err, cookie, style)
		if err := redis.New().SetHashRedis(cookie, style, ""); err != nil {
			return data, err
		}
		return data, nil
	}
	if len(val) != 0 {
		err := json.Unmarshal([]byte(val), &data)
		if err != nil {
			log.Debug("json hash redis Error ", err, cookie, val, style)
			return data, err
		}
	}
	return data, nil
}
// 更新Redis Carts Data
func UpdateRedisCarts(carts CartsVo.Carts, cookie string, subtotal, shipfee, beforeTotal, total float64) error {
	carts.SubTotal = subtotal
	carts.ShipFee = shipfee
	carts.BeforeTotal = beforeTotal
	carts.Total	= total
	err := SetRedisCarts(carts, cookie, carts.Style)
	if err != nil {
		return err
	}
	return nil
}

//登出刪除 Session
func DestroyUserLogin(cookie string) error {
	sess := session.Session{
		Name:"user" + cookie,
		TTL: 0,
	}
	if err := sess.Destroy(); err != nil{
		log.Error("Session Destroy Error", err)
		return err
	}
	return nil
}

//此商品是否為帳單,是否為同賣家,購物車內的是否為帳單
func CheckCartsAndProduct(cookie string, data entity.ProductsData, carts CartsVo.Carts) (CartsVo.Carts, error) {
	var resp CartsVo.Carts
	if len(carts.Products) == 0 || carts.Realtime == 1 || data.StoreId != carts.StoreId || data.IsRealtime == 1 {
		if err := DeleteRedisCarts(cookie, carts.Style); err != nil {
			return resp, err
		}
	} else {
		resp = carts
	}
	return resp, nil
}

//搜尋購物車內是否相同商品
func SearchCartsMatchItems(cartProduct []CartsVo.CartProduct, params *Request.AddCartParams) bool{
	var change = false
	for k, v := range cartProduct {
		if params.ProductSpecId == v.ProductSpecId {
			cartProduct[k].Quantity = v.Quantity + params.Quantity
			change = true
		}
	}
	return change
}

func ComposeCartsProducts(params *Request.AddCartParams, carts CartsVo.Carts) []CartsVo.CartProduct {
	var Products []CartsVo.CartProduct
	Products = carts.Products
	log.Debug("carts product", Products)
	//判斷購物車是否有物品
	if carts.Products == nil {
		var product1 CartsVo.CartProduct
		product1.ProductSpecId = params.ProductSpecId
		product1.Quantity = params.Quantity
		Products = append(Products, product1)
	} else {
		change := SearchCartsMatchItems(Products, params)
		if !change {
			var product2 CartsVo.CartProduct
			product2.ProductSpecId = params.ProductSpecId
			product2.Quantity = params.Quantity
			Products = append(Products, product2)
		}
	}
	log.Debug("carts product", Products)
	return Products
}
//B2C購物車
func NewB2cRedisCarts(data entity.B2cOrderData, cookie, style string) error {
	jsonData, _ := json.Marshal(data)
	val, _ := redis.New().GetHashRedis(cookie, style)
	if err := redis.New().SetHashRedis(cookie, style, string(jsonData)); err != nil {
		log.Error("Set Redis Error", err)
		return err
	}
	if len(val) == 0 {
		//設定購物車時間
		var min = 60 * time.Minute
		if err := redis.New().SetExpireRedis(cookie, min); err != nil {
			log.Error("Set Redis Expire Error", err)
			return err
		}
	}
	return nil
}
//取出B2c Redis購物車內容
func GetB2cRedisCarts(cookie, style string) (entity.B2cOrderData, error) {
	var data entity.B2cOrderData
	val, err := redis.New().GetHashRedis(cookie, style)
	if err != nil {
		log.Debug("Get hash redis Error ", err, cookie, style)
		if err := redis.New().SetHashRedis(cookie, style, ""); err != nil {
			return data, err
		}
		return data, nil
	}
	if len(val) != 0 {
		err := json.Unmarshal([]byte(val), &data)
		if err != nil {
			log.Debug("json hash redis Error ", err, cookie, val, style)
			return data, err
		}
	}
	return data, nil
}