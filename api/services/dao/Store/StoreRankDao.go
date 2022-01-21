package Store

import (
	"api/services/Enum"
	"api/services/database"
	"api/services/entity"
	"api/services/util/log"
	"time"
)

//建立賣家權限
func InsertStoreRankData(engine *database.MysqlSession, uid string, storeId string, rank string, status string, Email string) (entity.StoreRankData, error) {
	storeData := entity.StoreRankData{
		StoreId:    storeId,
		UserId:     uid,
		Rank:       rank,
		RankStatus: status,
		Email: Email,
		CreateTime: time.Now(),
		UpdateTime: time.Now(),
	}
	_, err := engine.Session.Table(entity.StoreRankData{}).Insert(&storeData)
	if err != nil {
		log.Error("Store Rank Data Database Error", err)
		return storeData, err
	}
	return storeData, nil
}
// 取出 GetStoreRankData ...
func GetStoreRankListByUid(engine *database.MysqlSession, UserId string) ([]entity.StoreDataResp, error) {
	var StoreRank []entity.StoreDataResp
	err := engine.Engine.Table(entity.StoreData{}).Select("*").
		Join("LEFT", entity.StoreRankData{}, "store_data.store_id = store_rank_data.store_id").
		Where("store_rank_data.user_id = ?", UserId).
		And("store_data.store_status != ?", Enum.StoreStatusEnd).
		And("store_rank_data.rank_status != ?", Enum.StoreRankSuspend).
		And("store_rank_data.rank_status != ?", Enum.StoreRankDelete).
		Asc("store_rank_data.rank").Find(&StoreRank)
	if err != nil {
		log.Error("Store Join Store Rank Database Error", err)
		return StoreRank, err
	}
	return StoreRank, nil
}
//計算管理者的賣場數量
func CountStoreRankByUid(engine *database.MysqlSession, Uid string) (int64, error) {
	result, err := engine.Engine.Table(entity.StoreRankData{}).
		And("user_id = ? AND rank_status != ? AND rank_status != ?", Uid, Enum.StoreRankDelete, Enum.StoreStatusSuspend).Count()
	if err != nil {
		log.Error("Count Store Rank Database Error", err)
		return result, err
	}
	return result, nil
}
//取會員賣場權限
func GetStoreManagerByStoreIdAndUserId(engine *database.MysqlSession, StoreId, UserId string) (entity.StoreRankData, error) {
	var data entity.StoreRankData
	if _, err := engine.Engine.Table(entity.StoreRankData{}).Select("*").
		Where("store_id = ?", StoreId).And("user_id = ?", UserId).Get(&data); err != nil {
		log.Error("get store rank data Database Error", err)
		return data, err
	}
	return data, nil
}
//取出收銀機下的管理者
func GetStoreManagerByStoreId(engine *database.MysqlSession, StoreId string) ([]entity.StoreRankData, error) {
	var data []entity.StoreRankData
	if err := engine.Engine.Table(entity.StoreRankData{}).Select("*").
		Where("store_id = ? AND rank = ? AND rank_status != ?", StoreId, Enum.StoreRankSlave, Enum.StoreRankDelete).Find(&data); err != nil {
		return data, err
	}
	return data, nil
}

//func GetStoreManagerByStoreId(engine *database.MysqlSession, StoreId string) ([]entity.StoreRankData, error) {
//	var StoreRank []entity.StoreRankData
//	sql := fmt.Sprintf("SELECT * FROM store_rank_data WHERE store_id = ? AND rank = ? AND rank_status = ?")
//	err := engine.Engine.SQL(sql, StoreId, Enum.StoreRankSlave, Enum.StoreRankSuccess).Find(&StoreRank)
//	if err != nil {
//		log.Error("get store rank data Database Error", err)
//		return StoreRank, err
//	}
//	return StoreRank, nil
//}
//取出收銀機下的所有主帳號及管理者
func GetStoreByStoreId(engine *database.MysqlSession, StoreId string) ([]entity.StoreRankData, error) {
	var StoreRank []entity.StoreRankData
	if err := engine.Engine.Table(entity.StoreRankData{}).Select("*").Where("store_id = ?", StoreId).
		And("rank_status = ?", Enum.StoreRankSuccess).Find(&StoreRank); err != nil {
		log.Error("get store rank data Database Error", err)
		return StoreRank, err
	}
	return StoreRank, nil
}
//更新Order Data
func UpdateStoreRankData(engine *database.MysqlSession, data entity.StoreRankData) (int64, error) {
	data.UpdateTime = time.Now()
	affected, err := engine.Session.Table(entity.StoreRankData{}).ID(data.RankId).AllCols().Update(data)
	if err != nil {
		return affected, err
	}
	return affected, nil
}

func CountStoreManagerBySellerId(engine *database.MysqlSession, SellerId string) (int64, error) {
	count, err := engine.Engine.Table(entity.StoreData{}).Select("count(*)").
		Join("LEFT", entity.StoreRankData{}, "store_data.store_id = store_rank_data.store_id").
		Where("store_data.seller_id = ?", SellerId).
		And("store_data.store_status = ?", Enum.StoreStatusSuccess).
		And("store_rank_data.rank = ?", Enum.StoreRankSlave).
		And("store_rank_data.rank_status != ?", Enum.StoreRankDelete).Count()
	if err != nil {
		return count, err
	}
	return count, nil
}
