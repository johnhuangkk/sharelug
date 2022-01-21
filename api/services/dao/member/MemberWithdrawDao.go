package member

import (
	"api/services/Enum"
	"api/services/database"
	"api/services/entity"
	"api/services/util/log"
	"time"
)

func InsertMemberWithdraw(engine *database.MysqlSession, data entity.MemberWithdrawData) error {
	data.CreateTime = time.Now()
	data.UpdateTime = time.Now()
	_, err := engine.Session.Table(entity.MemberWithdrawData{}).Insert(&data)
	if err != nil {
		log.Error("Insert Member Withdraw Database Error", err)
		return err
	}
	return nil
}

func GetMemberWithdrawByAccount(engine *database.MysqlSession, account string) (entity.MemberWithdrawData, error) {
	var data entity.MemberWithdrawData
	_, err := engine.Session.Table(entity.MemberWithdrawData{}).
		Select("*").Where("account = ?", account).Get(&data)
	if err != nil {
		return data, err
	}
	return data, nil
}

func GetMemberWithdrawByIdAndUserId(engine *database.MysqlSession, userId string, Id int64) (entity.MemberWithdrawData, error) {
	var data entity.MemberWithdrawData
	_, err := engine.Engine.Table(entity.MemberWithdrawData{}).
		Select("*").Where("id = ?", Id).And("user_id = ?", userId).Get(&data)
	if err != nil {
		return data, err
	}
	return data, nil
}

func GetMemberWithdrawListByUserId(engine *database.MysqlSession, userId string) ([]entity.MemberWithdrawData, error) {
	var data []entity.MemberWithdrawData
	err := engine.Engine.Table(entity.MemberWithdrawData{}).
		Select("*").Where("user_id = ?", userId).And("status = ?", Enum.WithdrawStatusSuccess).Find(&data)
	if err != nil {
		return data, err
	}
	return data, nil
}

func UpdateMemberWithdrawDeleteDefault(engine *database.MysqlSession, userId string) error {
	sql := "UPDATE member_withdraw_data SET act_default = ? WHERE user_id = ?"
	_, err := engine.Session.Exec(sql, 0, userId)
	if err != nil {
		log.Error("Update Member Withdraw Delete Default Error", err)
		return err
	}
	return nil
}

func UpdateMemberWithdrawSetDefault(engine *database.MysqlSession, userId, account string) error {
	sql := "UPDATE member_withdraw_data SET act_default = ? WHERE user_id = ? AND account = ?"
	_, err := engine.Session.Exec(sql, 1, userId, account)
	if err != nil {
		log.Error("Update Member Withdraw Default Error", err)
		return err
	}
	return nil
}


func UpdateMemberWithdraw(engine *database.MysqlSession, Id int64, data entity.MemberWithdrawData) error {
	data.UpdateTime = time.Now()
	_, err := engine.Session.Table(entity.MemberWithdrawData{}).ID(Id).Update(data)
	if err != nil {
		return err
	}
	return nil
}