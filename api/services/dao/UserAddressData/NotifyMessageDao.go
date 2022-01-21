package UserAddressData

import (
	"api/services/VO/Request"
	"api/services/database"
	"api/services/entity"
	"api/services/util/log"
	"fmt"
	"time"
)

// 寫入一筆資訊
func InsertNotifyMessage(params Request.NotifyRequest) error  {
	engine := database.GetMysqlEngine()
	defer engine.Close()

	data := entity.NotifyMessageData{}
	data.Email = params.Email
	data.Phone = params.Tel
	data.CreateTime = time.Now()

	if _, err := engine.Session.Table(entity.NotifyMessageData{}).Insert(data); err != nil {
		log.Error("Insert Notify Message Error", err)
		return err
	}
	return nil
}

//建立線上通知訊息
func InsertOnlineNotifyData(engine *database.MysqlSession, data entity.OnlineNotifyData) error {
	data.CreateTime = time.Now()
	data.UpdateTime = time.Now()
	if _, err := engine.Session.Table(entity.OnlineNotifyData{}).Insert(data); err != nil {
		log.Error("Insert Contact Error", err)
		return err
	}
	return nil
}
//計算未讀數
func CountOnlineNotifyUnreadByUserId(engine *database.MysqlSession, UserId []string) (int64, error) {
	result, err := engine.Engine.Table(entity.OnlineNotifyData{}).Select("count(*)").
		Where("unread = ?", 0).In("user_id", UserId).Count()
	if err != nil {
		log.Error("Count Online Notify Error", err)
		return 0, err
	}
	return result, nil
}

func CountOnlineNotifyUnreadByTypeAndUserId(engine *database.MysqlSession, UserId []string, Type string) (int64, error) {
	result, err := engine.Engine.Table(entity.OnlineNotifyData{}).Select("count(*)").
		Where("type = ? AND unread = ?", Type, 0).In("user_id", UserId).Count()
	if err != nil {
		log.Error("Count Online Notify Error", err)
		return 0, err
	}
	return result, nil
}

//更新線上通知訊息已讀
func UpdateOnlineNotifyRead(engine *database.MysqlSession, Id int64) error {
	sql := fmt.Sprintf("UPDATE online_notify_data SET unread = ?, update_time = ? WHERE id = ?")
	if _, err := engine.Session.Exec(sql, 1, time.Now(), Id); err != nil {
		log.Error("Update Online Notify Read Error", err)
		return err
	}
	return nil
}

func GetOnlineNotifyData(engine *database.MysqlSession, Id int64) (entity.OnlineNotifyData, error) {
	data := entity.OnlineNotifyData{}
	if _, err := engine.Engine.Table(entity.OnlineNotifyData{}).Select("*").
		Where("id = ?", Id).Get(&data); err != nil {
		return data, err
	}
	return data, nil
}

//計算線上通知訊息筆數
func CountOnlineNotifyByUserId(engine *database.MysqlSession, UserId []string, Type string) (int64, error) {
	result, err := engine.Engine.Table(entity.OnlineNotifyData{}).Select("count(*)").
		Where("type = ?", Type).In("user_id", UserId).Count()
	if err != nil {
		log.Error("Count Online Notify Error", err)
		return 0, err
	}
	return result, nil
}

//線上通知訊息列表
func GetOnlineNotifyList(engine *database.MysqlSession, UserId []string, Type string, start, limit int64) ([]entity.OnlineNotifyData, error) {
	var resp []entity.OnlineNotifyData
	start = (start - 1) * limit
	if err := engine.Engine.Table(entity.OnlineNotifyData{}).Select("*").
		Where("type = ?", Type).In("user_id", UserId).Desc("create_time").
		Limit(int(limit), int(start)).Find(&resp); err != nil {
		log.Error("Get BalanceAccountData Database Error", err)
		return resp, err
	}
	return resp, nil
}

//計算線上通知訊息筆數
func CountOnlineNotifySystemByUserId(engine *database.MysqlSession, UserId string) int64 {
	sql := fmt.Sprintf("SELECT count(*) FROM online_notify_data WHERE user_id = ? AND messages like '%s'", "%衷心感謝你加入成為Check’Ne的會員。%")
	result, err := engine.Engine.SQL(sql, UserId).Count()
	if err != nil {
		log.Error("Update Online Notify Read Error", err)
		return 0
	}
	return result
}
