package Request

type OrderRefundSearch struct {
	StoreId string `form:"StoreId" json:"StoreId"`
	Tab     string `form:"tab"`
	Limit   int    `form:"Limit" json:"Limit"`
	Start   int    `form:"Start" json:"Start"`
}

type RefundParams struct {
	OrderId string `form:"OrderId" json:"OrderId" valid:"required"`
	Amount  int64  `form:"Amount" json:"Amount" valid:"required"`
}

type ReturnParams struct {
	OrderId    string       `form:"OrderId" json:"OrderId" valid:"required"`
	IsReturn   int          `form:"Status" json:"IsReturn" valid:"required"`
	ReturnList []ReturnList `form:"ReturnList" json:"ReturnList" valid:"required"`
}

type ReturnList struct {
	ProductSpecId string `form:"ProductSpecId" json:"ProductSpecId" valid:"required"`
	Qty           int64  `form:"Qty" json:"Qty" valid:"required"`
}

type RefundQuery struct {
	OrderId string `form:"OrderId" json:"OrderId" valid:"required"`
}

type ReturnQuery struct {
	OrderId string `form:"OrderId" json:"OrderId"`
}

type ReturnListQuery struct {
	OrderId       string `form:"OrderId" json:"OrderId" valid:"required"`
	ProductSpecId string `form:"ProductSpecId" json:"ProductSpecId" valid:"required"`
}

type ReturnConfirmParams struct {
	ReturnId    string       `form:"ReturnId" json:"ReturnId" valid:"required"`
}
