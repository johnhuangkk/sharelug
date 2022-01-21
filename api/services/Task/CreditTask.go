package Task

import (
	"api/services/Service/KgiBank"
	"api/services/model"
	"api/services/util/log"
	"github.com/spf13/viper"
	"time"
)

func HandleCreditCheckTask()  {
	log.Info("Handle Credit Reply Task [Start]")
	if err := model.HandleCheckCredit(); err != nil {
		log.Error("Check credit Error", err)
	}
	log.Info("Handle Credit Reply Task [End]")
}

//信用卡請款回覆檔處理
func HandleCreditReplyTask() {
	log.Info("Handle Credit Reply Task")
	var MerchantID = viper.GetString("KgiCredit.C2C.3D.MerchantID")
	now := time.Now()
	d, _ := time.ParseDuration("-24h")
	yesterday := now.Add(d)
	err := KgiBank.DownloadFile(yesterday, MerchantID)
	if err != nil {
		log.Error("Download File Error", err)
	}
}
//信用卡請款
func HandleCreditCaptureTask() {
	log.Info("[start] Handle Credit Reply Task")
	KgiBank.HandleC2C3DCapture()
	KgiBank.HandleC2CN3DCapture()
	log.Info("[end] Handle Credit Reply Task")
}
