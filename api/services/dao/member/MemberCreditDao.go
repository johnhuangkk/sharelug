package member

import (
	"api/services/Enum"
	"api/services/database"
	"api/services/entity"
	"api/services/util/log"
	"api/services/util/tools"
	"fmt"
	"time"
)

func InsertMemberCreditData(engine *database.MysqlSession, data entity.MemberCardData) (entity.MemberCardData, error) {
	data.CreateTime = time.Now()
	data.UpdateTime = time.Now()
	_, err := engine.Session.Table(entity.MemberCardData{}).Insert(&data)
	if err != nil {
		log.Error("Database Error", err)
		return data, err
	}
	return data, nil
}

func GetMemberCreditDataByUserId(engine *database.MysqlSession, userId string) ([]entity.MemberCardData, error) {
	var data []entity.MemberCardData
	err := engine.Engine.Table(entity.MemberCardData{}).
		Select("*").Where("member_id = ?", userId).And("status = ?", Enum.CreditStatusSuccess).Find(&data)
	if err != nil {
		return data, err
	}
	return data, nil
}

func GetMemberCreditDataByCardIdAndUserId(engine *database.MysqlSession, UserId, CardId string) (entity.MemberCardData, error) {
	var data entity.MemberCardData
	_, err := engine.Engine.Table(entity.MemberCardData{}).
		Select("*").Where("card_id = ?", CardId).And("member_id = ?", UserId).Get(&data)
	if err != nil {
		return data, err
	}
	return data, nil
}

func GetMemberCreditDataByCardId(engine *database.MysqlSession, CardId string) (entity.MemberCardData, error) {
	var data entity.MemberCardData
	sql := fmt.Sprintf("SELECT * FROM member_card_data WHERE card_id = ?")
	_, err := engine.Engine.SQL(sql, CardId).Get(&data)
	if err != nil {
		return data, err
	}
	return data, nil
}

func UpdateMemberCardDataDefault(engine *database.MysqlSession, userId string) error {
	sql := "UPDATE member_card_data SET default_card = ? WHERE default_card != ? AND member_id = ?"
	_, err := engine.Session.Exec(sql, 0, 2, userId)
	if err != nil {
		log.Error("Update Member Card Data Error", err)
		return err
	}
	return nil
}

//更新會員信用卡資訊
func UpdateMemberCardData(engine *database.MysqlSession, data entity.MemberCardData) error {
	data.UpdateTime = time.Now()
	_, err := engine.Session.Table(entity.MemberCardData{}).ID(data.CardId).Update(data)
	if err != nil {
		return err
	}
	return nil
}


func GetMemberCreditByUserIdAndNumber(engine *database.MysqlSession, number, expire, userId string) (entity.MemberCardData, error) {
	_, card := tools.ParseCredit(number)
	var data entity.MemberCardData
	sql := fmt.Sprintf("SELECT * FROM member_card_data WHERE first4_digits = ? AND last4_digits = ? AND expiry_date = ? AND member_id = ?")
	_, err := engine.Engine.SQL(sql, card[0], card[3], expire, userId).Get(&data)
	if err != nil {
		return data, err
	}
	return data, nil
}

//取出所有信用卡
func GetMemberCreditData(engine *database.MysqlSession) ([]entity.MemberCardData, error) {
	var data []entity.MemberCardData
	if err := engine.Engine.Table(entity.MemberCardData{}).Select("*").Find(&data); err != nil {
		return data, err
	}
	return data, nil
}