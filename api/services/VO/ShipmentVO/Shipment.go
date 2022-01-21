package ShipmentVO

/**
	取號
*/
type Orders struct {
	OrderId []string `form:"orderId" json:"orderId" valid:"required"`
}

type Order struct {
	OrderId string `form:"orderId" json:"orderId" valid:"required"`
}

type SwitchOrder struct {
	OrderId string `form:"orderId" json:"orderId" valid:"required"`
	StoreId string `form:"storeId" json:"storeId" valid:"StoreId"`
}


// 郵箱取號設定賣家地址
type SellerSenderAddress struct {
	Id string
	Zip string
	Address string
}