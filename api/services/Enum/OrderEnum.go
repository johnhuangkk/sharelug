package Enum

const (
	OrderInit    = "INIT"
	OrderFail    = "FAIL"
	OrderExpire  = "EXPIRE"
	OrderWait    = "WAIT"
	OrderAudit   = "AUDIT"
	OrderPaid    = "PAID"
	OrderCancel  = "CANCEL"
	OrderSuccess = "SUCCESS"
)

var OrderStatus = map[string]string{
	OrderInit:    "已成立",
	OrderFail:    "失敗",
	OrderWait:    "待付款",
	OrderExpire:  "逾期未付",
	OrderAudit:   "訂單處理中",
	OrderPaid:    "已付款",
	OrderCancel:  "已取消",
	OrderSuccess: "已完成",
}

var ErpOrderStatus = map[string]string{
	OrderInit:    "已成立",
	OrderFail:    "失敗",
	OrderWait:    "待付款",
	OrderExpire:  "逾期未付",
	OrderAudit:   "訂單處理中",
	OrderPaid:    "已付款",
	OrderCancel:  "已取消",
	OrderSuccess: "已完成",
}

const (
	OrderRefundInit    = "INIT"
	OrderRefund        = "REFUND"
	OrderReturnSuccess = "SUCCESS"

	OrderAuditInit    = "INIT"
	OrderAuditNote    = "NOTE"
	OrderAuditPending = "PENDING"
	OrderAuditRefused = "REFUSED"
	OrderAuditRelease = "RELEASE"
)

var OrderAuditStatus = map[string]string{
	OrderAuditInit:    "待審",
	OrderAuditNote:    "照會",
	OrderAuditPending: "待決",
	OrderAuditRefused: "拒絕",
	OrderAuditRelease: "放行",
}

var OrderRefundStatus = map[string]string{
	OrderRefundInit:    "無",
	OrderRefund:        "退貨退款",
	OrderReturnSuccess: "退貨完成",
}

const (
	OrderShipInit                = "INIT"
	OrderShipTake                = "TAKE"
	OrderShipment                = "SHIPMENT"
	OrderShipTransit             = "TRANSIT"
	OrderShipShop                = "SHOP"
	OrderShipFail                = "FAIL"
	OrderShipOverdue             = "OVERDUE"
	OrderShipNotTaken            = "NOTTAKEN"
	OrderShipSuccess             = "SUCCESS"
	OrderShipReceiverStoreSwitch = "ShipRrStoreSwitch"
	OrderShipSenderStoreSwitch   = "ShipSrStoreSwitch"

	OrderShipSend         = "寄"
	OrderShipSenderIsSend = "ShipSenderIsSend"
	OrderShipOnShipping   = "ShipOnShipping"
	OrderShipNone         = "NONE"
)

var OrderShipStatus = map[string]string{
	OrderShipInit:         "待出貨",
	OrderShipTake:         "待出貨",
	OrderShipment:         "已出貨",
	OrderShipTransit:      "已出貨",
	OrderShipShop:         "已出貨",
	OrderShipFail:         "出貨失敗",
	OrderShipOverdue:      "待出貨",
	OrderShipNotTaken:     "逾期未取",
	OrderShipSuccess:      "已完成",
	OrderShipSenderIsSend: "賣家已到店交寄",
	OrderShipNone:         "已出貨",
}

var OrderF2fStatus = map[string]string{
	OrderShipInit:         "待取貨",
	OrderShipNone:         "已出貨",
}



const (
	OrderCaptureInit     = "INIT"
	OrderCaptureProgress = "PROGRESS"
	OrderCaptureAdvance  = "ADVANCE"
	OrderCapturePostpone = "POSTPONE"
	OrderCaptureSuccess  = "SUCCESS"
	OrderCaptureSuspend  = "SUSPEND"
)

var OrderCaptureStatus = map[string]string{
	OrderCaptureInit:     "未需撥付",
	OrderCaptureProgress: "待撥付",
	OrderCaptureAdvance:  "提前撥付",
	OrderCapturePostpone: "延後撥付",
	OrderCaptureSuccess:  "已撥付",
	OrderCaptureSuspend:  "暫停撥付",
}

const (
	OrderTransC2c  = "C2C"
	OrderTransB2c  = "B2C"
	OrderTransBill = "BILL"
	OrderTrans3D   = "3D"
	OrderTransN3D  = "N3D"
)

const (
	BillStatusInit    = "INIT"
	BillStatusOverdue = "OVERDUE"
	BillStatusClose   = "CLOSE"
	BillStatusCancel  = "CANCEL"
)

const (
	InvoiceOpenStatusNot  = "INIT"
	InvoiceOpenStatusOpen = "OPEN"
)

const (
	B2cOrderTypeUpgrade = "UPGRADE"
	B2cOrderTypeBilling = "BILLING"
)
