package Enum

const (
	InvoiceType07 = "07" //一般稅額計算
	InvoiceType08 = "08" //特種稅額計算

	InvoiceTaxType01 = "1" //應稅
	InvoiceTaxType02 = "2" //零稅率
	InvoiceTaxType03 = "3" //免稅
	InvoiceTaxType04 = "4" //應稅(特種稅率)
	InvoiceTaxType09 = "9" //混合應稅與免稅或零稅率

	InvoiceTaxRate = 1.05

	InvoiceAssignStatusEnable    = "ENABLED"   //啟用
	InvoiceAssignStatusDeEnabled = "DEENABLED" //停用

	AllowanceTypeBuyer    = "1" //買方開立折讓證明單
	AllowanceTypeSeller   = "2" //賣方折讓證明通知單
	AllowanceStatusInit   = "INIT"
	AllowanceStatusCancel = "CANCEL"

	InvoiceStatusCancel = "CANCEL" //
	InvoiceStatusNot    = "NOT"    //未開獎
	InvoiceStatusLose   = "LOSE"   //未中獎
	InvoiceStatusWin    = "WIN"    //已中獎
	InvoiceStatusDonate = "DONATE" //已捐贈
)

var InvoiceStatus = map[string]string{
	InvoiceStatusNot:    "未開獎",
	InvoiceStatusLose:   "未中獎",
	InvoiceStatusWin:    "已中獎",
	InvoiceStatusDonate: "已捐贈",
}

