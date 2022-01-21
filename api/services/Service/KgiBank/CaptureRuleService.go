package KgiBank

import (
	"api/services/util/log"
	"api/services/util/tools"
)

func NewKgiHeaderRule() *tools.FieldFormats {

	fieldFormats := &tools.FieldFormats{}
	tools.SetRules("Flag", "S", 1, "", fieldFormats)        //檔頭旗標
	tools.SetRules("MerchantId", "S", 10, "", fieldFormats) //特店代號
	tools.SetRules("SendDate", "S", 8, "", fieldFormats)    //送檔日期
	tools.SetRules("Seq", "N", 10, "", fieldFormats)        //序號
	tools.SetRules("Count", "N", 12, "", fieldFormats)      //總筆數
	tools.SetRules("Symbol", "S", 1, "", fieldFormats)      //金額正負號
	tools.SetRules("ToTal", "N", 12, "", fieldFormats)      //總金額
	tools.SetRules("Filler", "S", 216, "", fieldFormats)    //FILLER填空白

	log.Debug("ssss", fieldFormats)
	return fieldFormats
}

func NewKgiBodyRule() *tools.FieldFormats {

	fieldFormats := &tools.FieldFormats{}
	tools.SetRules("MerchantId", "S", 10, "", fieldFormats)  //特店代號
	tools.SetRules("TerminalId", "S", 8, "", fieldFormats)   //端末機代號
	tools.SetRules("OrderId", "S", 40, "", fieldFormats)     //訂單編號
	tools.SetRules("Space", "S", 19, "", fieldFormats)       //空白
	tools.SetRules("TranAmount", "N", 8, "", fieldFormats)   //交易金額
	tools.SetRules("AuthCode", "S", 8, "", fieldFormats)     //授權碼
	tools.SetRules("TranType", "N", 2, "", fieldFormats)     //交易碼
	tools.SetRules("TranDate", "S", 8, "", fieldFormats)     //交易日期
	tools.SetRules("Custom", "S", 16, "", fieldFormats)      //使用者自訂欄位
	tools.SetRules("CardInfo", "S", 40, "", fieldFormats)    //持卡人資訊
	tools.SetRules("ProcessDate", "S", 6, "", fieldFormats)  //帳單處理日期
	tools.SetRules("ResponseCode", "S", 3, "", fieldFormats) //回應碼
	tools.SetRules("ResponseMsg", "SC", 16, "", fieldFormats) //回應訊息
	tools.SetRules("BatchSeq", "N", 6, "", fieldFormats)     //Batch and seq. No.
	tools.SetRules("Mark", "S", 1, "", fieldFormats)         //分期付款或紅利積點註記
	tools.SetRules("NumberOfPay", "N", 2, "", fieldFormats)  //分期數
	tools.SetRules("FirstPayment", "N", 8, "", fieldFormats) //首期金額
	tools.SetRules("EachPayment", "N", 8, "", fieldFormats)  //每期金額
	tools.SetRules("Fees", "N", 6, "", fieldFormats)         //手續費
	tools.SetRules("DeductPoint", "N", 8, "", fieldFormats)  //本次扣底點數
	tools.SetRules("Symbol", "S", 1, "", fieldFormats)  //餘額正負號
	tools.SetRules("PointBalance", "N", 8, "", fieldFormats)  //卡人點數餘額
	tools.SetRules("Deductible", "N", 10, "", fieldFormats)  //卡人自付金額
	tools.SetRules("PaymentDate", "S", 8, "", fieldFormats)  //付款日
	tools.SetRules("Verify", "S", 1, "", fieldFormats)  //3D 認證結果
	tools.SetRules("Foreign", "S", 1, "", fieldFormats)  //國外卡
	tools.SetRules("Reserved", "S", 18, "", fieldFormats)  //預留

	return fieldFormats
}
