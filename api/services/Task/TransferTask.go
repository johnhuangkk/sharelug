package Task

import (
	"api/services/database"
	"api/services/model"
	"api/services/util/log"
	"time"
)

//銀行轉帳定時取回結果
func HandleC2cTransferTask() {
	log.Info("HandleC2cTransferTask [start]")
	engine := database.GetMysqlEngine()
	defer engine.Close()
	now := time.Now().AddDate(0, 0, -1).Format("20060102")
	temp := ""
	for true {
		smx, data, err := model.QueryC2cTransDateTransfer(now, now, temp)
		if err != nil {
			log.Error("Query TransDate Transfer", err)
		}
		for _, v := range data {
			log.Debug("Trans Date", v)
			if err = model.ProcessTransfer(engine, v); err != nil {
				log.Error("Process Transfer Error", v)
			}
		}
		if len(smx.SvcRs.TEMPDATA) == 0 {
			break
		} else {
			temp = smx.SvcRs.TEMPDATA
		}
	}
	log.Info("HandleC2cTransferTask [end]")
}

//銀行轉帳定時取回結果
func HandleB2cTransferTask() {
	log.Info("HandleB2cTransferTask [start]")
	engine := database.GetMysqlEngine()
	defer engine.Close()
	now := time.Now().AddDate(0, 0, -1).Format("20060102")
	temp := ""
	for true {
		smx, data, err := model.QueryB2cTransDateTransfer(now, now, temp)
		if err != nil {
			log.Error("Query TransDate Transfer", err, data)
		}
		for _, v := range data {
			log.Debug("Trans Date", v)
			if err = model.ProcessTransfer(engine, v); err != nil {
				log.Error("Process Transfer Error", v)
			}
		}
		if len(smx.SvcRs.TEMPDATA) == 0 {
			break
		} else {
			temp = smx.SvcRs.TEMPDATA
		}
	}
	log.Info("HandleB2cTransferTask [End]")
}

//取出轉帳逾期的訂單
func HandleTransferExpireTask() {
	log.Info("HandleTransferExpireTask [start]")
	engine := database.GetMysqlEngine()
	defer engine.Close()
	if err := model.ProcessTransferExpire(engine); err != nil {
		log.Error("Process Transfer Expire Error", err)
		return
	}
	log.Info("HandleTransferExpireTask [End]")
}

func HandleTransferQuery()  {
	log.Info("HandleTransferQueryTask [start]")
	engine := database.GetMysqlEngine()
	defer engine.Close()
	if err := model.ProcessTransferQuery(engine); err != nil {
		log.Error("Process Transfer Expire Error", err)
		return
	}
	log.Info("HandleTransferQueryTask [End]")
}