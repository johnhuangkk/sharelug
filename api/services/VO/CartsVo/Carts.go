package CartsVo

import (
	"api/services/Enum"
	"api/services/VO/Request"
	"api/services/entity"
	"api/services/util/tools"
)

type Carts struct {
	Products     []CartProduct `json:"Products"`
	Shipping     string        `json:"Shipping"`
	SubTotal     float64       `json:"SubTotal"`
	StoreId      string        `json:"StoreId"`
	BeforeTotal  float64       `json:"BeforeTotal"`
	Total        float64       `json:"Total"`
	ShipFee      float64       `json:"ShipFee"`
	Realtime     int           `json:"Realtime"`
	Coupon       float64       `json:"Coupon"`
	CouponNumber string        `json:"CouponNumber"`
	Style        string
}

type CartProduct struct {
	ProductSpecId string `json:"productSpecId"`
	Quantity      int    `json:"quantity"`
}

func (c Carts) GenerateOrder(orderId, buyerId string, storeData entity.StoreData, params *Request.PayParams) entity.OrderData {
	var data entity.OrderData
	data.OrderId = orderId
	data.SellerId = storeData.SellerId
	data.StoreId = storeData.StoreId
	data.OrderStatus = Enum.OrderInit
	data.RefundStatus = Enum.OrderRefundInit
	data.ShipStatus = Enum.OrderShipInit
	data.CaptureStatus = Enum.OrderCaptureInit
	data.CaptureApply = Enum.OrderCaptureInit
	data.PayWay = params.PayWay
	data.ShipType = c.Shipping
	data.SubTotal = c.SubTotal
	data.ShipFee = c.ShipFee
	data.TotalAmount = c.Total
	data.BuyerId = buyerId
	data.BuyerName = params.BuyerName
	data.BuyerPhone = params.BuyerPhone
	data.BuyerNotes = params.BuyerNotes
	if params.ReceiverName != "" {
		data.ReceiverName = params.ReceiverName
		data.ReceiverPhone = params.ReceiverPhone
	} else {
		data.ReceiverName = params.BuyerName
		data.ReceiverPhone = params.BuyerPhone
	}
	data.CsvCheck = 1
	//超商取貨付款 需將超取付入帳 設為0
	if params.PayWay == Enum.CvsPay {
		data.CsvCheck = 0
	}
	data.Coupon = c.Coupon
	data.CouponNumber = c.CouponNumber
	data.BeforeTotal = c.BeforeTotal
	//判斷無需配送和面交，不需要寫入地址
	if !tools.InArray([]string{Enum.F2F, Enum.NONE}, c.Shipping) {
		data.ReceiverAddress = params.Address
	}
	return data
}
