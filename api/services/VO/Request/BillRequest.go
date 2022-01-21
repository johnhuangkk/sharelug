package Request

type BillConfirmRequest struct {
	BillId string `json:"BillId"`
}

type BillListRequest struct {
	Tab   string `form:"Tab" json:"Tab"`
	Limit int64  `form:"Limit" json:"Limit"`
	Start int64  `form:"Start" json:"Start"`
}

type BuyerBillListRequest struct {
	Limit int64  `form:"Limit" json:"Limit"`
	Start int64  `form:"Start" json:"Start"`
}