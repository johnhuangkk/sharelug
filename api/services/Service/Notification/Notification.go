package Notification

import (
	"api/services/Enum"
	"api/services/VO/Redis"
	"api/services/dao/UserAddressData"
	"api/services/database"
	"api/services/entity"
	"api/services/util/log"
	"api/services/util/redis"
	"encoding/json"
	"fmt"
)

//發送系統通知
func SendSystemNotify(engine *database.MysqlSession, UserId, Message, MsgType, OrderId string) error {
	data := entity.OnlineNotifyData{}
	data.Type = Enum.NotifyTypeSystem
	data.UserId = UserId
	data.Messages = Message
	data.MsgType = MsgType
	data.OrderId = OrderId
	err := UserAddressData.InsertOnlineNotifyData(engine, data)
	if err != nil {
		log.Error("Insert online Notify Message Error", err)
		return err
	}
	user := []string{UserId}
	_, err = GetNotifyCount(engine, user)
	if err != nil {
		return err
	}
	return nil
}

//發送平台通知
func SendPlatformNotify(engine *database.MysqlSession, UserId, Message, MsgType, QuestionId string) error {
	data := entity.OnlineNotifyData{}
	data.Type = Enum.NotifyTypePlatform
	data.UserId = UserId
	data.Messages = Message
	data.MsgType = MsgType
	data.OrderId = QuestionId
	err := UserAddressData.InsertOnlineNotifyData(engine, data)
	if err != nil {
		log.Error("Insert online Notify Message Error", err)
		return err
	}
	return nil
}

//記算訊息未讀數
func GetNotifyCount(engine *database.MysqlSession, userId []string) (int64, error) {
	count, err := UserAddressData.CountOnlineNotifyUnreadByUserId(engine, userId)
	if err != nil {
		log.Error("Count online Notify Message Error", err)
		return 0, err
	}
	for _, v := range userId {
		if err = setRedisNotify(v, count); err != nil {
			return 0, err
		}
	}
	return count, nil
}

//Redis 加一筆記錄
func setRedisNotify(UserId string, count int64) error {
	data := Redis.Notify{}
	data.NotifyCount = count
	jsonData, _ := json.Marshal(data)
	err := redis.New().SetHashRedis(UserId, "notify", string(jsonData))
	if err != nil {
		return err
	}
	return nil
}

//取得redis內的Notify的筆數
func GetRedisNotify(UserId string) (int64, error) {
	data := Redis.Notify{}
	val, err := redis.New().GetHashRedis(UserId, "notify")
	if err != nil {
		log.Debug("Get hash redis Error ", err)
		return 0, err
	}
	if err = json.Unmarshal([]byte(val), &data); err != nil {
		log.Debug("json hash redis Error ", err)
		return 0, err
	}
	return data.NotifyCount, nil
}

//更新通知訊息已讀
func SetNotifyRead(engine *database.MysqlSession, Id int64) error {
	data, err := UserAddressData.GetOnlineNotifyData(engine, Id)
	if err != nil {
		log.Debug("Get Online Notify Error ", err)
		return err
	}
	if len(data.Messages) == 0 {
		return fmt.Errorf("無此訊息")
	}
	if err = UserAddressData.UpdateOnlineNotifyRead(engine, data.Id); err != nil {
		log.Debug("Update Online Notify Read Error ", err)
		return err
	}
	return nil
}

func CountNotify()  {
	
}