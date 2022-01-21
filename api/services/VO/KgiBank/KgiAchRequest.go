package KgiBank

type AchHeader struct {
	BOF    string //首錄別: 預設BOF
	CDATA  string //資料代號: 預設 ACHP01
	TDATE  string //處理日期: 民國YYYYMMDD
	TTIME  string //處理時間: HHMMSS
	SORG   string //發送單位代號: 預設
	RORG   string //接收單位代號:  預設9990250
	VERNO  string //版次: (預設固定為V10)
	FILLER string //備用: 預設空白
}

type AchBody struct {
	TYPE   string //交易型態: 預設 N
	TXTYPE string //交易類別: (SC:代付案件   SD:代收案件)
	TXID   string //交易代號:（預設405）
	SEQ    string //交易序號: 序號不得重複
	PBANK  string //提出行代號: （預設 8090072）
	PCLNO  string //發動者帳號: 位數不足時，右靠左補零
	RBANK  string //提回行代號:
	RCLNO  string //收受者帳號: 位數不足時，右靠左補零
	AMT    string //金額: 右靠左補零，金額不得為零。
	RCODE  string //退件理由代號: 預設空白
	SCHD   string //提示交換次序: 預設 B
	CID    string //發動者統一編號: (統一編號或身分證字號)
	PID    string //收受者統一編號: (統一編號或身分證字號)
	SID    string //上市上櫃公司代號: 預設空白
	PDATE  string //原提示交易日期: 預設空白
	PSEQ   string //原提示交易序號: 預設空白
	PSCHD  string //原提示交換次序: 預設空白
	CNO    string //用戶號碼: 預設空白
	NOTE   string //發動者專用區: 預設空白
	MEMO   string //存摺摘要: CheckNe
	CFEE   string //客戶支付手續費:
	NOTEB  string //發動行專用區: 預設空白
	FILLER string //備用: 預設空白
}

type AchFooter struct {
	EOF    string //尾錄別: 預設EOF
	CDATA  string //資料代號:  預設ACHP01
	TDATE  string //處理日期:  民國YYYYMMDD
	SORG   string //發送單位代號: 代表行代號
	RORG   string //接收單位代號: 預設 9990250
	TCOUNT string  //總筆數:
	TAMT   string  //總金額:
	YDATE  string //前一營業日日期: 預設空白
	FILLER string //備用: 預設空白
}
