package model

import (
	"api/services/Enum"
	"api/services/Service/Carts"
	"api/services/Service/Product"
	"api/services/VO/CartsVo"
	"api/services/VO/Request"
	"api/services/VO/Response"
	"api/services/dao/Promotion"
	"api/services/dao/product"
	"api/services/database"
	"api/services/entity"
	"api/services/util/log"
	"encoding/json"
	"fmt"
	"time"
)

type Pay struct {
	Product []Products `json:"Product"`
	Other	Other	   `json:"Other"`
}

type Products struct {
	ProductId 		string 	`json:"ProductId"`
	ProductName		string 	`json:"ProductName"`
	Spec			string 	`json:"Spec"`
	Price			int  	`json:"Price"`
	Quantity		int		`json:"Quantity"`
	ChooseAmount	int 	`json:"ChooseAmount"`
}

type Other struct {
	ChooseShip  string	`json:"ChooseShip"`
	ShipList	string	`json:"ShipList"`
	PaywayList	string	`json:"PaywayList"`
}

func AddProductToBill(cookie string, params *Request.AddCartParams) (CartsVo.Carts, error) {
	var resp CartsVo.Carts
	engine := database.GetMysqlEngine()
	defer engine.Close()
	if err := Carts.DeleteRedisCarts(cookie, Enum.StyleBill); err != nil {
		return resp, err
	}
	carts, err := Carts.GetRedisCarts(cookie, Enum.StyleBill)
	if err != nil {
		return resp, fmt.Errorf("1001001")
	}
	productData, err := checkProductStock(engine, carts, params)
	if err != nil {
		return resp, err
	}
	//購物車內的商品數量是否為0 或 購物車內的商品是帳單
	carts, err = Carts.CheckCartsAndProduct(cookie, productData, carts)
	if err != nil {
		return resp, fmt.Errorf("1001001")
	}
	//新增的商品放入購物車
	Products := Carts.ComposeCartsProducts(params, carts)
	var selected string
	if params.Shipping == "" {
		selected = carts.Shipping
	} else {
		selected = params.Shipping
	}
	data := new(CartsVo.Carts)
	data.Products = Products
	data.Shipping = selected
	data.Realtime = productData.IsRealtime
	data.StoreId = productData.StoreId
	if err := Carts.NewRedisCarts(data, cookie, Enum.StyleBill); err != nil {
		return resp, fmt.Errorf("1001001")
	}
	resp, err = Carts.GetRedisCarts(cookie, Enum.StyleBill)
	if err != nil {
		return resp, fmt.Errorf("1001001")
	}
	return resp, nil
}
//商品加入購物車
func CheckProductAndCartSameStore(cookie string, params *Request.AddCartParams) error {
	engine := database.GetMysqlEngine()
	defer engine.Close()
	carts, err := Carts.GetRedisCarts(cookie, Enum.StyleCarts)
	if err != nil {
		return fmt.Errorf("1001001")
	}
	if len(carts.StoreId) == 0 {
		return nil
	}
	data, err := product.GetProductByProductSpecId(engine, params.ProductSpecId)
	if err != nil {
		return fmt.Errorf("1001001")
	}
	//購物車內的商品數量是否為0 或 購物車內的商品是帳單
	if data.StoreId != carts.StoreId {
		return fmt.Errorf("1001001")
	}
	return nil
}

//商品加入購物車
func AddProductToCart(cookie string, params *Request.AddCartParams) error {
	engine := database.GetMysqlEngine()
	defer engine.Close()
	carts, err := Carts.GetRedisCarts(cookie, Enum.StyleCarts)
	if err != nil {
		return fmt.Errorf("1001001")
	}
	productData, err := checkProductStock(engine, carts, params)
	if err != nil {
		return err
	}
	//購物車內的商品數量是否為0 或 購物車內的商品是帳單
	carts, err = Carts.CheckCartsAndProduct(cookie, productData, carts)
	if err != nil {
		return fmt.Errorf("1001001")
	}
	//新增的商品放入購物車
	products := Carts.ComposeCartsProducts(params, carts)
	log.Debug("Products", products)
	var selected string
	if params.Shipping == "" {
		selected = carts.Shipping
	} else {
		selected = params.Shipping
	}
	data := new(CartsVo.Carts)
	data.Products = products
	data.Shipping = selected
	data.Realtime = productData.IsRealtime
	data.StoreId = productData.StoreId
	data.CouponNumber = carts.CouponNumber
	data.Coupon = carts.Coupon
	if err = Carts.NewRedisCarts(data, cookie, Enum.StyleCarts); err != nil {
		return fmt.Errorf("1001001")
	}
	if err := Carts.DeleteRedisCarts(cookie, Enum.StyleBill); err != nil {
		return fmt.Errorf("1001001")
	}
	return nil
}
//檢查庫存及限量
func checkProductStock(engine *database.MysqlSession, cares CartsVo.Carts, params *Request.AddCartParams) (entity.ProductsData, error) {
	qty := params.Quantity
	for _, v := range cares.Products {
		if v.ProductSpecId == params.ProductSpecId {
			qty = v.Quantity + params.Quantity
		}
	}
	log.Debug("qty", qty)
	//取出商品資料
	data, err := Product.CheckProductSpecQty(engine, params.ProductSpecId, qty)
	if err != nil {
		return data, err
	}
	//檢查商品是否為限購商品
	switch data.LimitKey {
		case Enum.ProductLimitLeast:
			if qty < int(data.LimitQty) {
				return data, fmt.Errorf("1009001")
			}
		case Enum.ProductLimitMost:
			if qty > int(data.LimitQty) {
				return data, fmt.Errorf("1009002")
			}
	}
	return data, nil
}

func GetProducts(engine *database.MysqlSession, cares CartsVo.Carts, params *Request.AddCartParams) (entity.ProductsData, error) {
	qty := params.Quantity
	for _, v := range cares.Products {
		if v.ProductSpecId == params.ProductSpecId {
			qty = v.Quantity + params.Quantity
		}
	}
	log.Debug("qty", qty)
	//取出商品資料
	data, err := Product.CheckProductSpecQty(engine, params.ProductSpecId, qty)
	if err != nil {
		return data, err
	}
	//檢查商品是否為限購商品
	switch data.LimitKey {
	case Enum.ProductLimitLeast:
		if qty < int(data.LimitQty) {
			return data, fmt.Errorf("1009001")
		}
	case Enum.ProductLimitMost:
		if qty > int(data.LimitQty) {
			return data, fmt.Errorf("1009002")
		}
	}
	return data, nil
}


/**
 * 取得購物車內的數量
 */
func GetCartsCount(cookie string) (int, error) {
	cart, err := Carts.GetCarts(cookie)
	if err != nil {
		return 0, fmt.Errorf("1009005")
	}
	return len(cart.Products), nil
}


//檢查運送方式 (20210329拿掉)
//func checkShipping(carts CartsVo.Carts,params *Request.AddCartParams) bool {
//	if carts.Shipping != "" && params.Shipping != "" {
//		if carts.Shipping != params.Shipping {
//			return true
//		}
//	}
//	return false
//}

func GetCartsData(cookie string, params Request.AddCartParams) (Response.CartPayResponse, error) {
	var resp Response.CartPayResponse
	var carts CartsVo.Carts
	if len(params.ProductSpecId) != 0 {
		cartData, _ := Carts.GetRedisCarts(cookie, Enum.StyleCarts)
		if cartData.Products == nil {
			data, err := AddProductToBill(cookie, &params)
			if err != nil {
				return resp, fmt.Errorf("1009005")
			}
			carts = data
		} else {
			if err := AddProductToCart(cookie, &params); err != nil {
				return resp, err
			}
			data, err := Carts.GetRedisCarts(cookie, Enum.StyleCarts)
			if err != nil {
				return resp, fmt.Errorf("1009005")
			}
			carts = data
		}
	} else {
		data, err := Carts.GetRedisCarts(cookie, Enum.StyleCarts)
		if err != nil {
			return resp, fmt.Errorf("1009005")
		}
		carts = data
	}
	resp, err := GetUserCartData(cookie, carts)
	if err != nil {
		return resp, err
	}
	return resp, nil
}

//取得使用者購物車內的資料
func GetUserCartData(cookie string,carts CartsVo.Carts) (Response.CartPayResponse, error) {
	engine := database.GetMysqlEngine()
	defer engine.Close()
	var cartResp Response.CartPayResponse

	if len(carts.Products) == 0 {
		log.Error("get product Error")
		return cartResp, fmt.Errorf("1009005")
	}
	product, storeId, shipList, payWayList, subtotal, merge, isCoupon, err := Product.ModifyProducts(engine, cookie, carts)
	if err != nil {
		log.Error("get product Error", err)
		return cartResp, fmt.Errorf("1001001")
	}
	shippingMode := modifyShip(shipList, merge)
	payWayMode := modifyPayWay(payWayList)
	shipFee := float64(product.GeneralFee + product.MergeFee)
	beforeTotal := subtotal + float64(product.GeneralFee + product.MergeFee)
	total := subtotal + float64(product.GeneralFee + product.MergeFee) - carts.Coupon
	cartResp = Response.CartPayResponse{
		Product: product,
		Subtotal: subtotal,
		Shipping: carts.Shipping,
		ShipFee: shipFee,
		Coupon: carts.Coupon,
		BeforeTotal: beforeTotal,
		Total: total,
		ShipList: shippingMode,
		PayWayList: payWayMode,
		StoreId: storeId,
		CouponNumber: carts.CouponNumber,
		IsCoupon: isCoupon,
	}
	if err := Carts.UpdateRedisCarts(carts, cookie, subtotal, shipFee, beforeTotal, total); err != nil {
		log.Error("update redis carts data error", err)
		return cartResp, fmt.Errorf("1001001")
	}
	return cartResp, nil
}

// 變更運送方式
func ChangeUserCartsShipping(cookie, choose string) error {
	carts, err := Carts.GetCarts(cookie)
	if err != nil {
		return fmt.Errorf("1009005")
	}
	carts.Shipping = choose
	if err := Carts.SetRedisCarts(carts, cookie, carts.Style); err != nil {
		return fmt.Errorf("1009006")
	}
	return nil
}

// 變更商品數量
func ChangeUserCartProductQuantity(cookie, productSpecId, Type string) error {
	engine := database.GetMysqlEngine()
	defer engine.Close()
	carts, err := Carts.GetCarts(cookie)
	if err != nil {
		return fmt.Errorf("1009005")
	}
	for k, v := range carts.Products {
		var quantity int
		if v.ProductSpecId == productSpecId {
			if Type == "add" {
				quantity = v.Quantity + 1
			} else {
				quantity = v.Quantity - 1
			}
			data, err := product.GetProductByProductSpecId(engine, v.ProductSpecId)
			if err != nil {
				return fmt.Errorf("1001001")
			}
			switch data.LimitKey {
				case Enum.ProductLimitLeast:
					if quantity < int(data.LimitQty) {
						return fmt.Errorf("1009001")
					}
				case Enum.ProductLimitMost:
					if quantity > int(data.LimitQty) {
						return fmt.Errorf("1009002")
					}
			}
			if data.Quantity >= int64(quantity) {
				carts.Products[k].Quantity = quantity
			} else {
				return fmt.Errorf("1009003")
			}
		}
	}
	if err := Carts.SetRedisCarts(carts, cookie, carts.Style); err != nil {
		return fmt.Errorf("1009004")
	}
	return nil
}
//撿查優惠卷及優惠卷放入購物車
func CheckCartCoupon(cookie, number string) error {
	engine := database.GetMysqlEngine()
	defer engine.Close()
	if len(number) != 8 {
		return fmt.Errorf("1009007")
	}
	carts, err := Carts.GetCarts(cookie)
	if err != nil {
		return fmt.Errorf("1009005")
	}
	//取出 購物車內的賣場
	storeId := carts.StoreId
	// Coupon Number 是否存在 StoreId 和 number
	data, err := Promotion.GetPromotionCodeByCodeAndStoreId(engine, storeId, number)
	if err != nil {
		return fmt.Errorf("1001001")
	}
	if len(data.PromotionCode.Code) == 0 {
		return fmt.Errorf("1009007")
	}
	if data.PromotionCode.IsUsed != false {
		return fmt.Errorf("1009009")
	}
	//判斷是否到期
	if time.Now().After(data.Promotion.CloseTime) || !data.Promotion.StopTime.IsZero() && time.Now().After(data.Promotion.StopTime) {
		return fmt.Errorf("1009009")
	}
	cart, err := GetUserCartData(cookie, carts)
	if err != nil {
		return fmt.Errorf("1001001")
	}
	//判斷 折價後金額是否小於等於0
	if cart.BeforeTotal - data.Promotion.Value <= 0 {
		return fmt.Errorf("1009008")
	}
	//如果存在 優惠金額是多少
	carts.Coupon = data.Promotion.Value
	carts.CouponNumber = number
	if err := Carts.SetRedisCarts(carts, cookie, carts.Style); err != nil {
		return fmt.Errorf("1009004")
	}

	return nil
}
//刪除購物車優惠卷資料
func DeleteCartCoupon(cookie string) error {
	engine := database.GetMysqlEngine()
	defer engine.Close()
	carts, err := Carts.GetCarts(cookie)
	if err != nil {
		return fmt.Errorf("1009005")
	}
	//清空優惠卷的相關資料
	carts.Coupon = 0
	carts.CouponNumber = ""
	if err := Carts.SetRedisCarts(carts, cookie, carts.Style); err != nil {
		return fmt.Errorf("1009004")
	}
	return nil
}
/**
 * 刪除購物車內的商品
 */
func DeleteCartsProduct(cookie, productSpecId string) error {
	carts, err := Carts.GetCarts(cookie)
	if err != nil {
		return fmt.Errorf("1009005")
	}
	if len(carts.Products) == 0 {
		return fmt.Errorf("1009005")
	}
	var products []CartsVo.CartProduct
	for k, v := range carts.Products {
		if v.ProductSpecId != productSpecId {
			products = append(products, carts.Products[k])
		}
	}
	cart := new(CartsVo.Carts)
	cart.Products = products
	cart.Shipping = carts.Shipping
	cart.StoreId = carts.StoreId
	cart.Style = carts.Style
	if len(cart.Products) == 0 {
		err = Carts.DeleteRedisCarts(cookie, carts.Style)
		if err != nil {
			return fmt.Errorf("1001001")
		}
	} else {
		err = Carts.NewRedisCarts(cart, cookie, cart.Style)
		if err != nil {
			return fmt.Errorf("1001001")
		}
	}
	return nil
}

/**
 * 合併付款方式
 */
func modifyPayWay(payWayList []string) []Response.PayWayMode {
	payWay := getPaywayIntersect(payWayList)
	var payWayMode []Response.PayWayMode
	for key, value := range payWay {
		payWayData := Response.PayWayMode{
			Type: key,
			Text: value,
		}
		payWayMode = append(payWayMode, payWayData)
	}
	return payWayMode
}

var shipMode map[string]int

/**
 * 合併運送方式
 */
func modifyShip(shipList []string, merge []int) []Response.ShippingMode {
	ShipList, Merge := reFormatShip(shipList, merge)
	shipMode = getShipIntersect(ShipList, Merge)
	var shippingMode []Response.ShippingMode
	for key, value := range shipMode {
		shippingMode = append(shippingMode, Response.ShippingMode{
			Type: key,
			Text: Enum.Shipping[key],
			Price: value,
		})
	}
	return shippingMode
}
//重新排序運費
func reFormatShip(shipList []string, merge []int) ([]string, []int) {
	var list1 []string
	var list2 []string
	var merge1 []int
	var merge2 []int
	for k, v := range merge {
		if v == 1 {
			list1 = append(list1, shipList[k])
			merge1 = append(merge1, v)
		} else {
			list2 = append(list2, shipList[k])
			merge2 = append(merge2, v)
		}
	}
	list1 = append(list1, list2...)
	merge1 = append(merge1, merge2...)
	return list1, merge1
}
/**
 * 取得選擇的運費
 */
func getShipFee(choose string) int {
	var fee int
	for key, value := range shipMode {
		if key == choose {
			fee = value
		}
	}
	return fee
}

/**
 * 付款方式合併取出交集
 */
func getPaywayIntersect(paywayList []string) map[string]string {
	var n1 map[string]string
	var n2 map[string]string
	for _, v := range paywayList {
		n2 = paywayToStringArray(v)
		if n1 == nil {
			n1 = n2
			n2 = nil
		} else {
			n1 = paywayIntersect(n1, n2)
			n2 = nil
		}
	}
	return n1
}

/**
 * 付款方式轉為陣列
 */
func paywayToStringArray(str string) map[string]string {
	var payway []Response.PayWayMode
	n := make(map[string]string)
	_ = json.Unmarshal([]byte(str), &payway)
	for _, val := range payway {
		n[val.Type] = val.Text
	}
	return n
}

/**
 * 付款方式取出交集
 */
func paywayIntersect(slice1, slice2 map[string]string) map[string]string {
	nn := make(map[string]string)
	for key, value := range slice2 {
		if _, ok := slice1[key]; ok  {
			nn[key] = value
		}
	}
	return nn
}

/**
 * 運送方式取出交集
 */
func getShipIntersect(shipList []string, merge []int) map[string]int{
	var n1 map[string]int
	var n2 map[string]int

	for k, v := range shipList {
		n2 = shipToStringArray(v)
		if n1 == nil {
			n1 = n2
			n2 = nil
		} else {
			n1 = shipIntersect(n1, n2, merge[k])
			n2 = nil
		}
	}
	return n1
}

/**
 * 轉換成string Array
 */
func shipToStringArray(str string) map[string]int {
	var ship []Response.ShippingMode
	n := make(map[string]int)
	_ = json.Unmarshal([]byte(str), &ship)

	for _, val := range ship {
		n[val.Type] = val.Price
	}
	return n
}

/**
 * 求交集
 */
func shipIntersect(slice1 map[string]int, slice2 map[string]int, m int) map[string]int {
	nn := make(map[string]int)
	for key, value := range slice2 {
		if _, ok := slice1[key]; ok  {
			if m != 0 {
				if slice1[key] > value {
					nn[key] = slice1[key]
				} else {
					nn[key] = value
				}
			} else {
				nn[key] = slice1[key] + value
			}
		}
	}
	return nn
}