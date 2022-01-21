package member

import (
	"api/services/database"
	"api/services/entity"
	"api/services/util/log"
	"fmt"
	"time"
)

func InsertAccountActivity(engine *database.MysqlSession, userId, action, message string) error {
	var ent entity.AccountActivityData
	ent.UserId = userId
	ent.Action = action
	ent.Message = message
	ent.CreateTime = time.Now()
	_, err := engine.Engine.Table(entity.AccountActivityData{}).Insert(&ent)
	if err != nil {
		log.Error("Insert Account Activity Database Error", err)
		return err
	}
	return nil
}

func GetAccountActivityByUserId(engine *database.MysqlSession, userId, startTime, endTime string) ([]entity.AccountActivityData, error) {
	var ent []entity.AccountActivityData

	sql := fmt.Sprintf("SELECT * FROM account_activity_data WHERE user_id = ? AND create_time BETWEEN ? and ? ORDER BY create_time DESC LIMIT 10")
	err := engine.Engine.SQL(sql, userId, startTime, endTime).Find(&ent)
	if err != nil {
		log.Error("Get Account Activity Database Error", err)
		return ent, err
	}
	return ent, nil
}
