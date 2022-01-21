package Enum

//餘額交易狀態
const (
	BalanceTypePayment  = "PAYMENT"
	BalanceTypeRefund = "REFUND"
	BalanceTypePlatformRefund = "PLATREFUND"
	BalanceTypeAdjustment = "ADJUSTMENT"
	BalanceTypeRelieve = "RELIEVE"
	BalanceTypeBalancePay = "BALANCE"
	BalanceTypePlatform = "PLATFORM"
	BalanceTypeService = "SERVICE"
	BalanceTypeWithdraw = "WITHDRAW"
	BalanceTypeBankFee  = "BANKFEE"
	BalanceTypeRetain = "RETAIN"
	BalanceTypeWdFailed = "WDFAILED"

	BalanceTypeDeposit = "DEPOSIT"
	BalanceTypeWithdrawal = "WITHDRAWAL"
	BalanceTypeCreditWait = "CREDITWAIT"
	BalanceTypeBillFail = "BILLFAIL"
	BalanceTypeDetain = "DETAIN"
)

var BalanceTrans = map[string]string {
	BalanceTypePayment: "撥付交易款項",
	BalanceTypeRefund: "訂單退款",
	BalanceTypePlatformRefund: "平台退款",
	BalanceTypeAdjustment: "調整款項",
	BalanceTypeRelieve: "解除扣留",
	BalanceTypeBalancePay: "餘額付款",
	BalanceTypePlatform: "平台費用",
	BalanceTypeService: "加值服務費",
	BalanceTypeWithdraw: "提領",
	BalanceTypeBankFee: "轉帳手續費",
	BalanceTypeRetain: "保留款項",
	BalanceTypeWdFailed: "轉帳失敗",

	BalanceTypeDeposit: "存入",
	BalanceTypeWithdrawal: "撥款",
	BalanceTypeCreditWait: "信用卡退款暫存",
	BalanceTypeBillFail: "訂購單失效",
	BalanceTypeDetain: "扣留款項",
}





