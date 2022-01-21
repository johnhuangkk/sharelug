package Enum

const (
	StoreRankMaster = "MASTER"
	StoreRankSlave = "SLAVE"

	StoreStatusSuccess = "SUCCESS"
	StoreStatusPending = "PENDING"
	StoreStatusClose = "CLOSE"
	StoreStatusSuspend = "SUSPEND"
	StoreStatusEnd = "END"

	StoreRankSuccess = "SUCCESS"
	StoreRankInit = "INIT"
	StoreRankDelete = "DELETE"
	StoreRankPending = "PENDING"
	StoreRankSuspend = "SUSPEND"

	UpgradeTypeRenew = "RENEW"
	UpgradeTypeSuspend = "SUSPEND"
)

var StoreRank = map[string]string {
	StoreRankMaster: "主帳號",
	StoreRankSlave: "管理帳號",
}

var StoreRankStatus = map[string]string {
	StoreRankSuccess: "啟用",
	StoreRankInit: "尚未啟用",
	StoreRankDelete: "刪除",
	StoreRankPending: "暫停中",
}

var StoreStatus = map[string]string {
	StoreStatusSuccess: "開啟中",
	//StoreStatusPending: "暫停中",
	StoreStatusClose: "暫停中",
	StoreStatusSuspend: "暫停中",
	StoreStatusEnd: "中止",
}

const (
	StoreSocialMediaTypeFacebook = "Facebook"
	StoreSocialMediaTypeInstagram = "Instagram"
	StoreSocialMediaTypeLine = "Line"
	StoreSocialMediaTypeTelegram = "Telegram"
)
