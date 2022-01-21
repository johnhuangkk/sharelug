package model

import (
	"api/services/Service/Sms"
	"api/services/dao/member"
	"api/services/database"
	"api/services/util/log"
)

func sendPaySms(engine *database.MysqlSession, buyerId string, message string) error {

	buyerData, err := member.GetMemberDataByUid(engine, buyerId)
	if err != nil {
		log.Error("send message error", err)
		return err
	}
	err = Sms.PushMessageSms(buyerData.Mphone, message)
	if err != nil {
		log.Error("send message error", err)
		return err
	}
	return nil
}



