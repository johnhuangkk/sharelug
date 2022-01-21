package Enum

const (
	CashFlowApply = "APPLY"
	FileExport = "EXPORT"

/**
 * CFAP ->金流確認 (Cash Flow Apply)
 * CFMC ->金流主管確認 (Cash Flow Manager Confirm)
 * FEXP -> 匯出檔案 (File Export)
 * FNAL -> 款項已匯出 (Final)
 * ACMC -> 出納放行
 * FCNF -> 總經理放行
 * RETN-> 銀行檔已回傳
 * DONE ->確認結果
 */

	EmailVerifyWait = "WAIT"
	EmailVerifySuccess = "SUCCESS"
	EmailVerifyExpired = "EXPIRED"
	EmailVerifyAlready = "ALREADY"

	WithdrawStatusInit = "INIT"
	WithdrawStatusSuccess = "SUCCESS"
	WithdrawStatusPending = "PENDING"
	WithdrawStatusWait = "WAIT"
	WithdrawStatusDelete = "DELETE"
	WithdrawStatusFailed = "FAILED"

	CreditStatusSuccess = "SUCCESS"
	CreditStatusDelete = "DELETE"

	EmailVerifyTypeUser = "USER"
	EmailVerifyTypeStore = "STORE"
)

var WithdrawStatus = map[string]string {
	//未提出、已提出、提出錯誤、退件、已匯款
	WithdrawStatusInit: "未提出",
	WithdrawStatusFailed: "提出錯誤",
	WithdrawStatusSuccess: "已匯款",
}

var EmailVerify = map[string]string {
	EmailVerifyWait: "待驗證",
	EmailVerifySuccess: "驗證完成",
	EmailVerifyExpired: "已到期",
}