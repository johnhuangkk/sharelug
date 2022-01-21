package cron

import (
	"api/services/Service/FamilyXml"
	"api/services/Service/HiLifeXml"
	"api/services/Service/IPost"
	"api/services/Service/OKXml"
	"api/services/Service/PostBag"
	"api/services/Service/SevenMyshipApi"
	"api/services/Task"
	"api/services/util/log"

	"github.com/robfig/cron/v3"
	"github.com/spf13/viper"
)

func Run() {
	c := cron.New(
		cron.WithSeconds(),
		// 提供秒字段的parser，如无该option秒字段不解析
		// cron.WithChain(DemoWrapper1, DemoWrapper2,cron.Recover(new(CronLogger))),
		// 使用自定义的全局cron执行链，在chain_cron.go定义
		// cron.WithLogger(new(CronLogger)),
		// 使用自定义的cron执行日志，在logger_cron.go定义
	)
	//			  ┌─────────────── 秒     (0 - 59)
	//            | ┌───────────── 分鐘   (0 - 59)
	//            | │ ┌─────────── 小時   (0 - 23)
	//            | │ │ ┌───────── 日     (1 - 31)
	//            | │ │ │ ┌─────── 月     (1 - 12)
	//            | │ │ │ │ ┌───── 星期幾 (0 - 6，0 是週日，6 是週六，7 也是週日)
	//			  │ │ │ │ │ │
	// c.AddFunc("5 * * * * *", func())

	ctx := map[string]error{}
	if viper.GetString("ENV") == "prod" {
		_, ctx["CloseAllProduct"] = c.AddFunc("00 00 17 30 9 *", Task.HandleCloseAllProductTask)
		// 每小時執行全家貨態更新
		_, ctx["Family"] = c.AddFunc("00 40 * * * *", FamilyXml.MartFamilyFetchShipping)
		// 每天凌晨執行全家店鋪資料更新
		_, ctx["FamilyStore"] = c.AddFunc("00 5 */8 * * *", FamilyXml.MartFamilyFetchStoreList)
		// 每小時執行萊爾富貨態更新
		_, ctx["HiLife"] = c.AddFunc("00 45 * * * *", HiLifeXml.MartHiLifeFetchShipping)
		// 每天凌晨執行萊爾富店鋪資料更新
		_, ctx["HiLifeStore"] = c.AddFunc("00 10 */8 * * *", HiLifeXml.MartHiLifeFetchStoreList)
		// Seven 每天下午三天更新店鋪
		_, ctx["SevenShop"] = c.AddFunc("00 10 15 * * *", SevenMyshipApi.FetchDailyShopStatus)
		// 每小時執行OK貨態更新
		_, ctx["OK"] = c.AddFunc("00 50 * * * *", OKXml.MartOkFetchShipping)
		// 每小時執行Seven貨態更新
		_, ctx["Seven"] = c.AddFunc("00 50 * * * *", SevenMyshipApi.FetchShipment)
		// 每天凌晨執行OK店鋪資料更新
		_, ctx["OKStore"] = c.AddFunc("00 15 */8 * * *", OKXml.MartOKFetchStoreList)
		// 更新上海銀行虛擬帳號入進資料
		//_, ctx["Scsb"] = c.AddFunc("00 */30 * * * *", model.FetchAccountingRecord)
		//ftp 便利包
		_, ctx["IPostFtpShip"] = c.AddFunc("00 */13 * * * *", IPost.HandleIPostCVS)
		_, ctx["PostBag"] = c.AddFunc("00 */31 * * * *", PostBag.UpdateStatus)
		_, ctx["PostBag"] = c.AddFunc("00 */5 * * * *", PostBag.BuildXml)
		_, ctx["PostBag"] = c.AddFunc("00 */11 * * * *", PostBag.UploadShipOrderFile)
		_, ctx["PostBag"] = c.AddFunc("00 */21 * * * *", PostBag.CheckFileUpdate)
		// 凌晨兩點三十分跑 寫入超商帳務
		_, ctx["CvsAccountingTask"] = c.AddFunc("00 30 02 * * *", Task.FetchCvsFtpAccountingTask)
		// 凌晨三點十五分跑 超商帳務核帳
		_, ctx["CvsAccountingTask"] = c.AddFunc("00 15 03 * * *", Task.CvsAccountingTask)
		// 更新i 郵箱資訊
		_, ctx["IPostBoxData"] = c.AddFunc("00 34 */3 * * *", IPost.UpdateIPostDataTask)
		// 撥款
		_, ctx["Balance"] = c.AddFunc("00 30 9 * * *", Task.HandleAppropriationTask)
		// 逾期未寄
		_, ctx["ShipExpire"] = c.AddFunc("00 00 */1 * * *", Task.HandleShipExpireTask)
		// 轉帳
		_, ctx["c2cTransfer"] = c.AddFunc("00 00 */1 * * *", Task.HandleC2cTransferTask)
		_, ctx["b2cTransfer"] = c.AddFunc("00 00 */1 * * *", Task.HandleB2cTransferTask)
		// 逾期未轉帳
		_, ctx["TransferExpire"] = c.AddFunc("00 00 */6 * * *", Task.HandleTransferExpireTask)
		// 逾期帳單關閉
		_, ctx["RealtimeExpire"] = c.AddFunc("00 00 */1 * * *", Task.HandleRealtimeExpireTask)
		// 信用卡檢核
		_, ctx["CheckCredit"] = c.AddFunc("00 */30 * * * *", Task.HandleCreditCheckTask)
		// 信用卡請款結果檔
		_, ctx["Credit"] = c.AddFunc("00 00 11 * * *", Task.HandleCreditReplyTask)
		// 信用卡請款
		_, ctx["CreditCapture"] = c.AddFunc("00 00 11 * * *", Task.HandleCreditCaptureTask)
		// 排程掃描正向到店第四天未領取 發出小鈴鐺通知
		_, ctx["Shop4Days"] = c.AddFunc("00 00 09 * * *", Task.ScanOrderShipStatusIsShopEqual4Days)
		// 排程掃描取號後兩天未寄出 發出小鈴鐺通知
		_, ctx["Take2Days"] = c.AddFunc("00 00 10 * * *", Task.ScanOrderShipStatusIsInitEqual2Days)
		// OK 順向到貨通知 到店 四天 七天
		_, ctx["OKShipIsShop"] = c.AddFunc("00 00 10 * * *", Task.OkOrderShipIsShopSmsNotification)
		// 加值服務到期付款及產生帳單和通知
		//_, ctx["UpgradeExpire"] = c.AddFunc("00 00 */1 * * *", Task.HandleUpgradeExpireTask)
		// 每1小時開立發票
		_, ctx["OpenInvoice"] = c.AddFunc("00 00 */2 * * *", Task.HandleInvoiceTask)
		// 產生日結表
		_, ctx["DayStatement"] = c.AddFunc("00 10 00 * * *", Task.HandleDayStatementTask)
		// 買家帳單到期
		_, ctx["BillExpire"] = c.AddFunc("00 00 6 * * *", Task.HandleBillExpireTask)

		_, ctx["KgiSpecialStoreRecord"] = c.AddFunc("00 00 00 * * *", Task.UploadKgiSpecialStoreRecord)
	}

	if len(ctx) != 0 {
		log.Debug("cron run error:", ctx)
	}
	if len(c.Entries()) > 0 {
		c.Start()
		log.Info("The Cron Jobs Running")
	} else {
		log.Info("No Cron Entries")
	}
}
//// 加值服務到期 中止服務
//_, ctx["UpgradeExpireStop"] = c.AddFunc("00 00 */1 * * *", Task.HandleUpgradeExpireStopTask)
		