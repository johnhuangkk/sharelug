package Email

import (
	"api/services/database"
	"api/services/entity"
	"api/services/util/log"
	"time"
)

func InsertEmailVerifyData(engine *database.MysqlSession, data entity.EmailVerifyData) error {
	data.CreateTime = time.Now()
	data.UpdateTime = time.Now()
	_, err := engine.Session.Table("email_verify_data").Insert(&data)
	if err != nil {
		log.Error("Mail Verify Database Error", err)
		return err
	}
	return nil
}

func GetEmailVerifyDataByCode(engine *database.MysqlSession, code string) (entity.EmailVerifyData, error) {
	var data entity.EmailVerifyData
	_, err := engine.Engine.Table(entity.EmailVerifyData{}).
		Select("*").Where("verify_code = ?", code).Get(&data)
	if err != nil {
		return data, err
	}
	return data, nil
}

func UpdateEmailVerifyDataByCode(engine *database.MysqlSession, status string, id int64) error {
	sql := "UPDATE email_verify_data SET verify_status = ? WHERE id = ?"
	if _, err := engine.Session.Exec(sql, status, id); err != nil {
		log.Error("Update Member Mail verify code Data Error", err)
		return err
	}
	return nil
}