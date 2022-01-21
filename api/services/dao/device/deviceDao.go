package device

import (
	"api/services/VO/TokenVo"
	"api/services/database"
	"api/services/entity"
	"api/services/util/log"
	"github.com/google/uuid"
	"time"
)

func InsertDevice(engine *database.MysqlSession, params TokenVo.TokenParams, ip string) (entity.DeviceData, error) {
	uuid := uuid.New()
	var data entity.DeviceData
	data.DeviceUuid = uuid.String()
	data.Platform = params.Platform
	data.PlatformVersion = params.PlatformVersion
	data.PlatformDevice = params.PlatformDevice
	data.PlatformIP = ip
	data.FcmToken = params.FcmToken
	data.CreateTime = time.Now()
	data.UpdateTime = time.Now()
	_, err := engine.Session.Table("device_data").Insert(&data)
	if err != nil {
		log.Error("Insert Device Data Database Error", err)
		return data, err
	}
	return data, nil
}
