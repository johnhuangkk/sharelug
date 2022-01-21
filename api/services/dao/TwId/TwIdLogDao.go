package TwId

import (
	"api/services/VO/Response"
	"api/services/database"
	"api/services/entity"
	"api/services/util/log"
	"time"
)


func InsertTwIdLogData(engine *database.MysqlSession, data entity.TwIdVerifyData) (entity.TwIdVerifyData, error) {
	data.CreateTime = time.Now()
	data.UpdateTime = time.Now()
	_, err := engine.Session.Table(entity.TwIdVerifyData{}).Insert(&data)
	if err != nil {
		log.Error("Database Error", err)
		return data, err
	}
	return data, nil
}

//更新商品資訊
func UpdateTwIdLogData(engine *database.MysqlSession, Id int, Data entity.TwIdVerifyData, resp Response.TWIDVerify) error {
	Data.HttpResponseCode = resp.HttpCode
	Data.HttpResponseMsg = resp.HttpMsg
	Data.Response = resp.Response.CheckIDCardApply
	Data.ResponseCode = resp.RespCode
	Data.ResponseMsg = resp.RespMsg

	Data.UpdateTime = time.Now()
	_, err := engine.Session.Table("tw_id_verify_data").ID(Id).Update(Data)
	if err != nil {
		return err
	}
	return nil
}

func GetTwIdLogDataByUserId(engine *database.MysqlSession, UserId string) (entity.TwIdVerifyData, error) {
	var data entity.TwIdVerifyData
	if _, err := engine.Engine.Table(entity.TwIdVerifyData{}).Select("*").
		Where("user_id = ? AND response = ?", UserId, 1).Desc("create_time").Get(&data); err != nil {
		return data, err
	}
	return data, nil
}