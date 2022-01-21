package member

import (
	"api/services/Enum"
	"api/services/database"
	"api/services/entity"
	"api/services/util/log"
	"time"
)

//新增會員發票載具資訊
func InsertMemberCarrierData(engine *database.MysqlSession, data entity.MemberCarrierData) error {
	data.CreateTime = time.Now()
	data.UpdateTime = time.Now()
	_, err := engine.Session.Table(entity.MemberCarrierData{}).Insert(&data)
	if err != nil {
		log.Error("Insert Member Carrier Database Error", err)
		return err
	}
	return nil
}

//更新會員發票載具資訊
func UpdateMemberCarrierData(engine *database.MysqlSession, data entity.MemberCarrierData) error {
	data.UpdateTime = time.Now()
	_, err := engine.Session.Table(entity.MemberCarrierData{}).ID(data.MemberId).Update(data)
	if err != nil {
		log.Error("Update Member Carrier Database Error", err)
		return err
	}
	return nil
}

//取得會員發票載具資訊
func GetMemberCarrierByMemberId(engine *database.MysqlSession, memberId string) (entity.MemberCarrierData, error) {
	var data entity.MemberCarrierData
	_, err := engine.Engine.Table(entity.MemberCarrierData{}).Where("member_id = ?", memberId).Get(&data)
	if err != nil {
		log.Error("Get Member Carrier Database Error", err)
		return data, err
	}
	return data, nil
}

func InsertDonateData(data entity.DonateData) error {
	engine := database.GetMysqlEngine()
	defer engine.Close()
	_, err := engine.Session.Table(entity.DonateData{}).Insert(&data)
	if err != nil {
		log.Error("Insert Bank Code Database Error", err)
		return err
	}
	return nil
}

func GetDonateData(engine *database.MysqlSession) ([]entity.DonateData, error) {
	var data []entity.DonateData
	if err := engine.Session.Table(entity.DonateData{}).Where("donate_status = ?", Enum.OrderSuccess).Find(&data); err != nil {
		log.Error("Insert Bank Code Database Error", err)
		return data, err
	}
	return data, nil
}

func GetDonateDataByDonateCode(engine *database.MysqlSession, DonateCode string) (entity.DonateData, error) {
	var data entity.DonateData
	if _, err := engine.Session.Table(entity.DonateData{}).Where("donate_code = ?", DonateCode).Get(&data); err != nil {
		log.Error("Insert Bank Code Database Error", err)
		return data, err
	}
	return data, nil
}
