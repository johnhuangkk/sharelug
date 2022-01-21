package Enum

//轉帳狀態
const (
	TransferInit = "Init"
	TransferSuccess = "Success"
	TransferDuplicate = "Duplicate"
)

//信用卡狀態
const (
	CreditTransStatusInit = "INIT"
	CreditTransStatusSuccess = "SUCCESS"
	CreditTransStatusFail = "FAIL"
	CreditTransStatusCancel = "CANCEL"

	CreditTransTypeAuth = "AUTH"
	CreditTransTypeVoid = "VOID"
	CreditTransTypeRefund = "REFUND"

	CreditAuditInit = "INIT"
	CreditAuditWait = "WAIT"
	CreditAuditNote = "NOTE"
	CreditAuditPending = "PENDING"
	CreditAuditRefused = "REFUSED"
	CreditAuditRelease = "RELEASE"

	CreditCaptureInit = "INIT"
	CreditCaptureFail = "FAIL"
	CreditCaptureSuccess = "SUCCESS"
)

var CreditCaptureStatus = map[string]string{
	CreditCaptureInit: "待請款",
	CreditCaptureSuccess: "已請款",
}

var AuthReport = map[string]string {
	CreditTransTypeAuth: "代收交易款項",
	CreditTransTypeRefund: "交易退款款項",
}

const (
	ServiceTypeUpgrade = "Upgrade"
	ServiceTypeShop = "Shop"
)