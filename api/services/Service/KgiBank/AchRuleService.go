package KgiBank

import "api/services/util/tools"

func NewKgiAchHeaderRule() *tools.FieldFormats {

	fieldFormats := &tools.FieldFormats{}
	tools.SetRules("BOF", "S", 3, "", fieldFormats) //首錄別
	tools.SetRules("CDATA", "S", 6, "", fieldFormats) //資料代號
	tools.SetRules("TDATE", "S", 8, "", fieldFormats) //處理日期
	tools.SetRules("TTIME", "S", 6, "", fieldFormats) //處理時間
	tools.SetRules("SORG", "S", 7, "", fieldFormats) //發送單位代號
	tools.SetRules("RORG", "S", 7, "", fieldFormats) //接收單位代號
	tools.SetRules("VERNO", "S", 3, "", fieldFormats) //版次
	tools.SetRules("FILLER", "S", 210, "", fieldFormats) //備用

	return fieldFormats
}

func NewKgiAchBodyRule() *tools.FieldFormats {

	fieldFormats := &tools.FieldFormats{}
	tools.SetRules("TYPE", "S", 1, "", fieldFormats) //交易型態
	tools.SetRules("TXTYPE", "S", 2, "", fieldFormats) //交易類別
	tools.SetRules("TXID", "S", 3, "", fieldFormats) //交易代號
	tools.SetRules("SEQ", "N", 8, "", fieldFormats) //交易序號
	tools.SetRules("PBANK", "S", 7, "", fieldFormats) //提出行代號
	tools.SetRules("PCLNO", "N", 16, "", fieldFormats) //發動者帳號
	tools.SetRules("RBANK", "S", 7, "", fieldFormats) //提回行代號
	tools.SetRules("RCLNO", "N", 16, "", fieldFormats) //收受者帳號
	tools.SetRules("AMT", "N", 10, "", fieldFormats) //金額
	tools.SetRules("RCODE", "N", 2, "", fieldFormats) //退件理由代號
	tools.SetRules("SCHD", "S", 1, "", fieldFormats) //提示交換次序
	tools.SetRules("CID", "S", 10, "", fieldFormats) //發動者統一編號
	tools.SetRules("PID", "S", 10, "", fieldFormats) //收受者統一編號
	tools.SetRules("SID", "S", 6, "", fieldFormats) //上市上櫃公司代號
	tools.SetRules("PDATE", "N", 8, "", fieldFormats) //原提示交易日期
	tools.SetRules("PSEQ", "N", 8, "", fieldFormats) //原提示交易序號
	tools.SetRules("PSCHD", "S", 1, "", fieldFormats) //原提示交換次序
	tools.SetRules("CNO", "S", 20, "", fieldFormats) //用戶號碼
	tools.SetRules("NOTE", "S", 40, "", fieldFormats) //發動者專用區
	tools.SetRules("MEMO", "S", 10, "", fieldFormats) //存摺摘要
	tools.SetRules("CFEE", "N", 5, "", fieldFormats) //客戶支付手續費
	tools.SetRules("NOTEB", "S", 20, "", fieldFormats) //發動行專用區
	tools.SetRules("FILLER", "S", 39, "", fieldFormats) //備用

	return fieldFormats
}

func NewKgiAchFooterRule() *tools.FieldFormats {

	fieldFormats := &tools.FieldFormats{}
	tools.SetRules("EOF", "S", 3, "", fieldFormats) //尾錄別
	tools.SetRules("CDATA", "S", 6, "", fieldFormats) //資料代號
	tools.SetRules("TDATE", "S", 8, "", fieldFormats) //處理日期
	tools.SetRules("SORG", "S", 7, "", fieldFormats) //發送單位代號
	tools.SetRules("RORG", "S", 7, "", fieldFormats) //尾錄別
	tools.SetRules("TCOUNT", "N", 8, "", fieldFormats) //總筆數
	tools.SetRules("TAMT", "N", 16, "", fieldFormats) //尾錄別
	tools.SetRules("YDATE", "S", 8, "", fieldFormats) //尾錄別
	tools.SetRules("FILLER", "S", 187, "", fieldFormats) //尾錄別

	return fieldFormats
}