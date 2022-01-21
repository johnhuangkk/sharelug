package Enum

const (
	SysLogLogin = "LOGIN"
	SysLogGetCard = "GETCARD"
	SysLogSetCard = "SETCARD"



	SysLogFail = "FAIL"
	SyslogSuccess = "SUCCESS"
)

const (
	ActivityMgt = "MANAGEMENT"
	ActivityStore = "STORE"
	ActivityAppn = "APPROPRIATION"
	ActivityWithdraw = "WITHDRAW"
	ActivitySetBank = "BANK"
	ActivitySetCredit = "CREDIT"
	ActivityShutdown = "SHUTDOWN"
	ActivityCancelManager = "CANCEL"
	ActivityChangePhone = "PHONE"
	ActivityCancelEmail = "EMAIL"
)

var ActivityStatus = map[string]string {
	ActivityMgt : "指派管理帳號",
	ActivityStore : "開設收銀機",
	ActivityAppn : "款項撥付",
	ActivityWithdraw : "提領",
	ActivitySetBank: "設定銀行帳號",
	ActivitySetCredit: "設定信用卡",
	ActivityShutdown: "關閉賣場",
	ActivityCancelManager: "取消管理帳號",
	ActivityChangePhone: "變更手機號碼",
	ActivityCancelEmail: "變更email",
}