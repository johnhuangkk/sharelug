package Store

import (
	"api/services/database"
	"api/services/entity"
	"api/services/util/log"
	"time"
)
// 新增 Insert Store Social Media Data
func InsertStoreSocialMediaData(engine *database.MysqlSession, data entity.StoreSocialMediaData) error {
	data.UpdateTime = time.Now()
	data.CreateTime = time.Now()
	if _, err := engine.Session.Table(entity.StoreSocialMediaData{}).Insert(&data); err != nil {
		log.Error("Database Error", err)
		return err
	}
	return nil
}
// 更新StoreSocialMediaData
func UpdateStoreSocialMediaData(engine *database.MysqlSession, data entity.StoreSocialMediaData) error {
	data.UpdateTime = time.Now()
	if _, err := engine.Session.Table(entity.StoreSocialMediaData{}).ID(data.StoreId).AllCols().Update(data); err != nil {
		log.Error("Update member Error", err)
		return err
	}
	return nil
}
// 取出 GetStoreSocialMediaData
func GetStoreSocialMediaDataByStoreId(engine *database.MysqlSession, storeId string) (entity.StoreSocialMediaData, error) {
	var data entity.StoreSocialMediaData
	if _, err := engine.Engine.Table(entity.StoreSocialMediaData{}).
		Select("*").Where("store_id = ?", storeId).Get(&data); err != nil {
		return data, err
	}
	return data, nil
}