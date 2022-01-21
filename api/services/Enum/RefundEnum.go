package Enum

const (
	TypeRefund = "REFUND"
	TypeReturn = "RETURN"
)

const (
	RefundStatusWait = "WAIT"
	RefundStatusAudit = "AUDIT"
	RefundStatusSuccess = "SUCCESS"
	ReturnStatusSuccess = "SUCCESS"
)

var RefundStatus = map[string]string{
	RefundStatusWait: "待退款",
	RefundStatusAudit: "退款處理中",
	RefundStatusSuccess: "已退款",
}

var ReturnStatus = map[string]string{
	RefundStatusWait: "待退貨",
	RefundStatusAudit: "退款處理中",
	ReturnStatusSuccess: "已退貨",
}
