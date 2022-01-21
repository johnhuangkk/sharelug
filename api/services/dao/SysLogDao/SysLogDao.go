package SysLogDao

import (
	"api/services/database"
	"api/services/entity"
	"api/services/util/log"
)

func InsertSysLog(engine *database.MysqlSession, data entity.SystemLog) error {
	if _, err := engine.Session.Table(entity.SystemLog{}).Insert(&data); err != nil {
		log.Error("Insert Status history Log Database Error", err)
		return err
	}
	return nil
}

func GetSystemLogByUserId(engine *database.MysqlSession, userId, start, end string) ([]entity.SystemLog, error) {
	var data []entity.SystemLog
	if err := engine.Engine.Table(entity.SystemLog{}).Select("*").
		Where("user_id = ? AND create_time BETWEEN ? AND ?", userId, start, end).
		Desc("create_time").Find(&data); err != nil {
		return data, err
	}
	return data, nil
}