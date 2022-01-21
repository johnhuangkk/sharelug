package UserAddressData

import (
	"api/services/database"
	"api/services/entity"
	"api/services/util/log"
	"api/services/util/tools"
	"fmt"
	"time"
)

// 同一配送方式只能有一個寄件地址
func CheckSendAddressUniqOfShip(engine *database.MysqlSession, info entity.UserAddressData) (bool, error) {
	exist := new(entity.UserAddressData)

	query := map[string]interface{}{}
	query["uid"] = info.Uid
	query["ship"] = info.Ship
	query["status"] = info.Status
	query["type"] = info.Type

	has, err := engine.Engine.Table(entity.UserAddressData{}).Select("*").Where(query).Exist(exist)
	if err != nil {
		log.Error("CheckSendAddressUniqOfShip Error", err)
		return has, err
	}
	return has, err
}


// 檢查相同地址是否存在
func CheckExistAddress(engine *database.MysqlSession, info entity.UserAddressData) (bool, error) {
	exist := new(entity.UserAddressData)

	query := map[string]interface{}{}
	query["address"] = info.Address
	query["uid"] = info.Uid
	query["ship"] = info.Ship
	query["status"] = info.Status
	query["real_name"] = info.RealName
	query["phone"] = info.Phone

	has, err := engine.Engine.Table(entity.UserAddressData{}).Select("*").Where(query).Exist(exist)

	if err != nil {
		log.Error("CheckExistAddress Error", err)
		return has, err
	}

	return has, err
}


// 寫入一筆資訊
func InsertAddressInfo(engine *database.MysqlSession, info entity.UserAddressData) (entity.UserAddressData, error)  {
	now := time.Now()
	info.UpdateTime = now
	info.UaId = tools.MD5(now.Format("20060102150405") + info.Uid)
	_, err := engine.Session.Table("user_address_data").Insert(info)
	if err != nil {
		log.Error("InsertAddressInfo Error", err)
		return info, err
	}
	return info, err
}

// 更新資訊
func UpdateSenderInfo(engine *database.MysqlSession, info entity.UserAddressData) error {
	_, err := engine.Session.Table("user_address_data").ID(info.UaId).Update(info)
	if err != nil {
		log.Error("UpdateSenderInfo Error", err)
		return err
	}
	return nil
}

// 取得單筆地址資訊
func QueryAddressByUaId(engine *database.MysqlSession, uaId string, uid string) (entity.UserAddressData, error) {
	info := entity.UserAddressData{}
	query := map[string]interface{}{}
	query["status"] = "Y"
	query["uid"] = uid
	query["ua_id"] = uaId

	_, err := engine.Engine.Table("user_address_data").Select("*").Where(query).Get(&info)

	if err != nil {
		log.Error("get QueryAddressByUaId Error", err)
		return info, err
	}

	return info, err
}

// 配送方式取得預設寄件/收件地址
func QuerySendAddressByShip(engine *database.MysqlSession, uid string, ship string, aType string) (entity.UserAddressData, error) {
	info := entity.UserAddressData{}
	query := map[string]interface{}{}
	query["status"] = "Y"
	query["uid"] = uid
	query["ship"] = ship
	query["type"] = aType
	_, err := engine.Engine.Table("user_address_data").Select("*").Where(query).Get(&info)
	if err != nil {
		log.Error("get QueryAddressByUaId Error", err)
		return info, err
	}
	return info, err
}

// 取得地址
func QueryAddresses(engine *database.MysqlSession, uid string, aType string, ship string) ([]entity.UserAddressData, error) {
	var info  []entity.UserAddressData
	query := map[string]interface{}{}
	query["status"] = "Y"
	query["uid"] = uid
	query["type"] = aType

	if len(ship) != 0 {
		query["ship"] = ship
	}
	err := engine.Engine.Table("user_address_data").Select("*").Where(query).OrderBy("ship asc, update_time desc").Find(&info)
	if err != nil {
		log.Error("get QueryAddresses Error", err)
		return info, err
	}
	return info, nil
}


func GetExistAddress(engine *database.MysqlSession, info entity.UserAddressData) (entity.UserAddressData, error) {
	var data entity.UserAddressData
	query := map[string]interface{}{}
	query["address"] = info.Address
	query["uid"] = info.Uid
	query["ship"] = info.Ship
	query["status"] = info.Status
	_, err := engine.Engine.Table(entity.UserAddressData{}).Select("*").Where(query).Get(&data)
	if err != nil {
		log.Error("CheckExistAddress Error", err)
		return data, err
	}
	return data, err
}

func UpdateAddresses(engine *database.MysqlSession, data entity.UserAddressData) error {
	data.UpdateTime = time.Now()
	_, err := engine.Session.Table(entity.UserAddressData{}).Where("ua_id = ?", data.UaId).Update(data)
	if err != nil {
		return  err
	}
	return nil
}

func GetDeliveryLastAddress(engine *database.MysqlSession, userId, payWay string) (entity.UserAddressData, error) {
	var data entity.UserAddressData
	sql := fmt.Sprintf("SELECT * FROM user_address_data WHERE uid = ? AND ship = ? AND status = ? AND type = ? ORDER BY update_time DESC LIMIT 0,1")
	_, err := engine.Engine.SQL(sql, userId, payWay, "Y", "R").Get(&data)
	if err != nil {
		log.Error("Get Delivery Last Address Error", err)
		return data, err
	}
	return data, nil
}

func GetUserAddresses(engine *database.MysqlSession, id string) (entity.UserAddressData, error) {
	var data entity.UserAddressData
	_, err := engine.Session.Table(entity.UserAddressData{}).Where("ua_id = ?", id).Get(&data)
	if err != nil {
		return data, err
	}
	return data, nil
}
