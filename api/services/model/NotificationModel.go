package model

import (
	"api/services/Enum"
	"api/services/Service/Notification"
	"api/services/VO/Request"
	"api/services/VO/Response"
	"api/services/dao/UserAddressData"
	"api/services/database"
	"api/services/entity"
	"api/services/util/log"
	"fmt"
	"strings"
)
//訊息通知列表
func GetNotificationList(userData entity.MemberData, storeData entity.StoreDataResp, params Request.NotificationRequest) (Response.NotificationResponse, error) {
	engine := database.GetMysqlEngine()
	defer engine.Close()
	var resp Response.NotificationResponse
	var userId []string
	if storeData.Rank == Enum.StoreRankMaster {
		userId = []string{storeData.StoreId, userData.Uid}
	} else {
		userId = []string{storeData.StoreId}
	}
	count, err := UserAddressData.CountOnlineNotifyByUserId(engine, userId, strings.ToUpper(params.Tab))
	if err != nil {
		return resp, fmt.Errorf("1001001")
	}
	resp.NotifyCount = count
	data, err := UserAddressData.GetOnlineNotifyList(engine, userId, params.Tab, params.Start, params.Limit)
	if err != nil {
		return resp, fmt.Errorf("1001001")
	}
	for _, value := range data {
		var res Response.NotifyResponse
		res.MessageId = value.Id
		res.Message = value.Messages
		res.MsgType = value.MsgType
		res.OrderId = value.OrderId
		res.CreateTime = value.CreateTime.Format("2006/01/02 15:04")
		res.Unread = false
		if value.Unread == 1 {
			res.Unread = true
		}
		resp.NotifyMessage = append(resp.NotifyMessage, res)
	}
	platform, err := UserAddressData.CountOnlineNotifyUnreadByTypeAndUserId(engine, userId, Enum.NotifyTypePlatform)
	if err != nil {
		log.Error("Count Online Notify Data Error", err)
		return resp, fmt.Errorf("1001001")
	}
	system, err := UserAddressData.CountOnlineNotifyUnreadByTypeAndUserId(engine, userId, Enum.NotifyTypeSystem)
	if err != nil {
		log.Error("Count Online Notify Data Error", err)
		return resp, fmt.Errorf("1001001")
	}
	resp.Tabs.Platform = platform
	resp.Tabs.System = system

	return resp, nil
}

//訊息通知已讀
func NotifyRead(userData entity.MemberData, params Request.NotificationReadRequest) (int64, error) {
	engine := database.GetMysqlEngine()
	defer engine.Close()
	if err := Notification.SetNotifyRead(engine, params.MessageId); err != nil {
		return 0, fmt.Errorf("1001001")
	}
	count, err := UserAddressData.CountOnlineNotifyUnreadByTypeAndUserId(engine, []string{userData.Uid}, Enum.NotifyTypePlatform)
	if err != nil {
		return 0, fmt.Errorf("1001001")
	}
	return count, nil
}
