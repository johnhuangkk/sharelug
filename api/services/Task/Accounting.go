package Task

import (
	"api/services/Enum"
	"api/services/Service/CvsAccounting"
	"api/services/Service/FamilyXml"
	"api/services/Service/HiLifeXml"
	"api/services/Service/OKXml"
	"api/services/database"
	"api/services/util/log"
	"time"
)

// 寫入超商核帳資料
func FetchCvsFtpAccountingTask() {
	FamilyXml.MartFamilyFetchAccounting()
	HiLifeXml.MartHiLifeAccounting()
	OKXml.MarkOkAccounting()
}

// 超商核帳
func CvsAccountingTask() {
	engine := database.GetMysqlEngine()
	defer engine.Close()

	{
		limitTime := time.Now().Add(-(time.Hour * time.Duration(24*13))).Format(`2006-01-02`)
		log.Info(`CvsAccountingTask Start`)
		{
			CvsAccounting.Send(engine, Enum.CVS_FAMILY, limitTime)
			CvsAccounting.Receive(engine, Enum.CVS_FAMILY, limitTime)
			log.Info(`CvsAccountingTask %s End`, Enum.CVS_FAMILY)
		}
		{
			CvsAccounting.Send(engine, Enum.CVS_HI_LIFE, limitTime)
			CvsAccounting.Receive(engine, Enum.CVS_HI_LIFE, limitTime)
			log.Info(`CvsAccountingTask %s End`, Enum.CVS_HI_LIFE)
		}
		{
			CvsAccounting.Send(engine, Enum.CVS_OK_MART, limitTime)
			CvsAccounting.Receive(engine, Enum.CVS_OK_MART, limitTime)
			log.Info(`CvsAccountingTask %s End`, Enum.CVS_OK_MART)
		}
		{
			CvsAccounting.Send(engine, Enum.CVS_7_ELEVEN, limitTime)
			CvsAccounting.Receive(engine, Enum.CVS_7_ELEVEN, limitTime)
			log.Info(`CvsAccountingTask %s End`, Enum.CVS_7_ELEVEN)
		}
		log.Info(`CvsAccountingTask End`)
	}
}
