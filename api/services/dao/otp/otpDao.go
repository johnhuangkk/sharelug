package otp

import (
	"api/services/database"
	"api/services/entity"
	"fmt"
	"time"
)

var otpData entity.OtpData

func InsertOtp(engine *database.MysqlSession, phone, code, uid, email string) (entity.OtpData, error) {
	now := time.Now()
	date, _ := time.ParseDuration("15m")
	var data = entity.OtpData{
		Uid:        uid,
		Phone:      phone,
		Email: 		email,
		OtpNumber:  code,
		ExpireTime: now.Add(date),
		CreateTime: now,
		SendFreq: 1,
	}
	_, err := engine.Session.Table("otp_data").Insert(&data)
	if err != nil {
		return data, err
	}
	return data, nil
}

//取出OTP DATA
func GetOtpByPhone(engine *database.MysqlSession, phone string) (entity.OtpData, error) {
	var data entity.OtpData
	sql := fmt.Sprintf("SELECT * FROM otp_data WHERE phone = ? AND otp_use = ? ORDER BY create_time DESC LIMIT 0,1")
	var _, err = engine.Engine.SQL(sql, phone, 0).Get(&data)
	if err != nil {
		return data, err
	}
	return data, nil
}

//更新Order Data
func UpdateOtpData(engine *database.MysqlSession, data entity.OtpData) error {
	_, err := engine.Session.Table("otp_data").ID(data.Id).Update(data)
	if err != nil {
		return err
	}
	return nil
}
