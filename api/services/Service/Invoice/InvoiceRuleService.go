package Invoice

import "api/services/util/tools"

func NewInvoiceRule() *tools.FieldFormats {
	fieldFormats := &tools.FieldFormats{}
	tools.SetRules("TCompanyBan", "S", 10, "", fieldFormats) //賣方-營業人總機構統一編號
	tools.SetRules("InvoiceYm", "S", 5, "", fieldFormats) //發票所屬年月
	tools.SetRules("InvoiceAxle", "S", 2, "", fieldFormats) //發票號碼-字軌
	tools.SetRules("InvoiceNumber", "S", 8, "", fieldFormats) //發票號碼-號碼
	tools.SetRules("SellerName", "S", 60, "", fieldFormats) //賣方-營業人名稱
	tools.SetRules("SellerBan", "S", 10, "", fieldFormats) //賣方-營業人統一編號
	tools.SetRules("Address", "S", 100, "", fieldFormats)//賣方-營業人地址
	tools.SetRules("InvoiceDate", "S", 10, "", fieldFormats)//發票日期
	tools.SetRules("InvoiceTime", "S", 8, "", fieldFormats)//發票時間
	tools.SetRules("TotalAmount", "S", 12, "", fieldFormats)//總計
	tools.SetRules("CarrierType", "S", 6, "", fieldFormats)//載具類別號碼
	tools.SetRules("CarrierName", "S", 60, "", fieldFormats)//載具類別名稱
	tools.SetRules("CarrierIdClear", "S", 64, "", fieldFormats)//載具顯碼id
	tools.SetRules("CarrierIdHide", "S", 64, "", fieldFormats)//載具隱碼id
	tools.SetRules("RandomNumber", "S", 4, "", fieldFormats)//四位隨機碼
	tools.SetRules("PrizeType", "S", 1, "", fieldFormats)//中獎獎別
	tools.SetRules("PrizeAmt", "S", 10, "", fieldFormats)//中獎獎金
	tools.SetRules("BuyerBan", "S", 10, "", fieldFormats)//買受人-扣繳單位統一編號
	tools.SetRules("DepositMK", "S", 1, "", fieldFormats)//整合服務平台已匯款註記
	tools.SetRules("DataType", "S", 1, "", fieldFormats)//資料類別
	tools.SetRules("Code", "S", 2, "", fieldFormats)//例外代碼
	tools.SetRules("Print", "S", 2, "", fieldFormats)//列印格式
	tools.SetRules("Identifier", "S", 24, "", fieldFormats)//唯一識別碼
	return fieldFormats
}