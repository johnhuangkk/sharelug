package Task

import (
	"api/services/Enum"
	"api/services/Service/Upgrade"
	"api/services/dao/member"
	"api/services/database"
	"api/services/model"
	"api/services/util/log"
	"time"
)

//加值服務到期付款及產生帳單和通知
func HandleUpgradeExpireTask()  {
	log.Info("HandleUpgradeExpireTask [start]")
	engine := database.GetMysqlEngine()
	defer engine.Close()
	expire := time.Now().Add(-(time.Hour * time.Duration(24)))
	data, err  := member.GetMemberUpgradeExpire(engine, expire)
	if err != nil {
		log.Error("Get Member Upgrade Expire Error", err)
	}
	//產生帳單
	for _, v := range data {
		if v.UpgradeType == Enum.UpgradeTypeSuspend {
			if err := Upgrade.MemberDemoteLevelService(engine, v, 0); err != nil {
				log.Error("Upgrade Demote Level Error", err)
			}
		} else {
			if err := model.GeneratorWaitPaymentOrder(engine, v); err != nil {
				log.Error("Generator Upgrade Order Error", err)
			}
		}
	}
	log.Info("HandleUpgradeExpireTask [end]")
}

//加值服務到期 中止服務
func HandleUpgradeExpireStopTask()  {
	log.Info("HandleUpgradeExpireStopTask [start]")
	engine := database.GetMysqlEngine()
	defer engine.Close()
	//取出到期10天的會員
	expire := time.Now().Add(-(time.Hour * time.Duration(24) * 10))
	data, err  := member.GetMemberUpgradeExpire(engine, expire)
	if err != nil {
		log.Error("Get Member Upgrade Expire Error", err)
	}
	//關閉賣場及管理者
	for _, v := range data {
		err := Upgrade.CloseUpgradeService(engine, v)
		if err != nil {
			log.Error("Close Store Upgrade Error", err)
		}
	}
	log.Info("HandleUpgradeExpireStopTask [end]")
}