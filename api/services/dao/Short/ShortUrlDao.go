package Short

import (
	"api/services/database"
	"api/services/entity"
	"api/services/util/log"
	"time"
)

func InsertShortUrlData(engine *database.MysqlSession, data entity.ShortUrlData) error {
	data.CreateTime = time.Now()
	_, err := engine.Session.Table(entity.ShortUrlData{}).Insert(&data)
	if err != nil {
		log.Error("Short Url Database insert Error", err)
		return err
	}
	return nil
}

func GetShortUrlDataByShort(engine *database.MysqlSession, Short string) (entity.ShortUrlData, error) {
	var data entity.ShortUrlData
	if _, err := engine.Engine.Table(entity.ShortUrlData{}).Select("*").Where("short = ?", Short).Get(&data); err != nil {
		log.Error("get RealTim Database Error", err)
		return data, err
	}
	return data, nil
}

