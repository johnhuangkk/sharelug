package Enum

const (
	CVS_FAMILY           = "CVS_FAMILY"
	CVS_HI_LIFE          = "CVS_HI_LIFE"
	CVS_OK_MART          = "CVS_OK_MART"
	CVS_7_ELEVEN         = "CVS_7_ELEVEN"
	DELIVERY_POST_BAG1   = "DELIVERY_POST_BAG1"
	DELIVERY_POST_BAG2   = "DELIVERY_POST_BAG2"
	DELIVERY_POST_BAG3   = "DELIVERY_POST_BAG3"
	DELIVERY_I_POST_BAG1 = "DELIVERY_I_POST_BAG1"
	DELIVERY_I_POST_BAG2 = "DELIVERY_I_POST_BAG2"
	DELIVERY_I_POST_BAG3 = "DELIVERY_I_POST_BAG3"
	DELIVERY_POST        = "DELIVERY_POST"
	DELIVERY_T_CAT       = "DELIVERY_T_CAT"
	DELIVERY_E_CAN       = "DELIVERY_E_CAN"
	DELIVERY_OTHER       = "DELIVERY_OTHER"
	I_POST               = "I_POST"
	F2F                  = "F2F"
	NONE                 = "NONE"
	SELF_DELIVERY        = "SELF_DELIVERY"
)

var Shipping = map[string]string{
	CVS_7_ELEVEN:         "7-11超商",
	CVS_FAMILY:           "全家超商",
	CVS_OK_MART:          "OK超商",
	CVS_HI_LIFE:          "萊爾富超商",
	DELIVERY_POST_BAG1:   "郵局便利包",
	DELIVERY_POST_BAG2:   "郵局便利包",
	DELIVERY_POST_BAG3:   "郵局便利包",
	DELIVERY_I_POST_BAG1: "i 郵箱交寄，宅配到府",
	DELIVERY_I_POST_BAG2: "i 郵箱交寄，宅配到府",
	DELIVERY_I_POST_BAG3: "i 郵箱交寄，宅配到府",
	DELIVERY_POST:        "中華郵政",
	DELIVERY_T_CAT:       "黑貓宅急便",
	DELIVERY_E_CAN:       "宅配通",
	DELIVERY_OTHER:       "宅配到府",
	I_POST:               "i郵箱",
	F2F:                  "面交/自取",
	NONE:                 "無須配送",
	SELF_DELIVERY: 		  "外送",
}

const (
	Credit    = "CREDIT"
	Transfer  = "TRANSFER"
	Balance   = "BALANCE"
	CvsPay    = "CVS_PAY"
	TaiwanPay = "TW_PAY"
)

var PayWay = map[string]string{
	Credit:    "信用卡",
	Transfer:  "ATM轉帳",
	Balance:   "Check'Ne餘額",
	CvsPay:    "貨到付款",
	TaiwanPay: "台灣PAY",
}

var PayWayReport = map[string]string{
	Credit:    "信用卡",
	Transfer:  "虛擬帳號",
	Balance:   "Check'Ne餘額",
	CvsPay:    "貨到付款",
	TaiwanPay: "台灣PAY",
}

const (
	ProductStatusSuccess = "SUCCESS"
	ProductStatusCancel  = "CANCEL"
	ProductStatusDown    = "DOWN"    //下架
	ProductStatusPending = "PENDING" //賣場關閉商品下架
	ProductStatusDelete  = "DELETE"
	ProductStatusOverdue = "OVERDUE"
	ProductStatusPaid    = "PAID"
	ProductForcedRemove  = "REMOVE" //官方強制下架
)

var RealtimeStatus = map[string]string{
	ProductStatusSuccess: "帳單成立",
	ProductStatusCancel:  "已取消",
	ProductStatusOverdue: "逾期未繳",
	ProductStatusPaid:    "已付款",
}

const (
	ProductLimitNone  = "NONE"
	ProductLimitLeast = "LEAST"
	ProductLimitMost  = "MOST"

	FreeShipNone     = "NONE"
	FreeShipAmount   = "AMOUNT"
	FreeShipQuantity = "QUANTITY"
)