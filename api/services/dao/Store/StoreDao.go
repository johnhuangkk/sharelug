package Store

import (
	"api/services/Enum"
	"api/services/database"
	"api/services/entity"
	"api/services/util/log"
	"api/services/util/tools"
	"fmt"
	"strings"
	"time"
)

// 取出 GetStoreData ...
func GetStoresByUid(engine *database.MysqlSession, Uid string) ([]entity.StoreDataResp, error) {
	var store []entity.StoreDataResp
	sql := fmt.Sprintf("SELECT * FROM store_rank_data a LEFT JOIN store_data b ON a.store_id = b.store_id" +
		" WHERE seller_id = ? AND rank = ?")
	err := engine.Engine.SQL(sql, Uid, Enum.StoreRankMaster).Find(&store)
	if err != nil {
		return store, err
	}
	return store, nil
}

func GetStoreDefaultDataByUid(engine *database.MysqlSession, Uid string) (entity.StoreData, error) {
	var store entity.StoreData
	if _, err := engine.Engine.Table(entity.StoreData{}).Select("*").
		Where("seller_id = ?", Uid).And("store_default = ?", 1).Get(&store); err != nil {
		return store, err
	}
	return store, nil
}

// 取出 GetStoreData ...
func GetStoreDataByStoreId(engine *database.MysqlSession, Sid string) (entity.StoreData, error) {
	var store entity.StoreData
	if _, err := engine.Engine.Table(entity.StoreData{}).Select("*").Where("store_id = ?", Sid).Get(&store); err != nil {
		return store, err
	}
	return store, nil
}

func GetStoreDataByUserIdAndStoreId(engine *database.MysqlSession, UserId, StoreId string) (entity.StoreDataResp, error) {
	var store entity.StoreDataResp
	sql := fmt.Sprintf("SELECT * FROM store_data s LEFT JOIN store_rank_data r ON s.store_id = r.store_id" +
		" WHERE s.store_id = ? AND r.user_id = ? AND r.rank_status = ? ORDER BY rank ASC, store_default DESC")
	_, err := engine.Engine.SQL(sql, StoreId, UserId, Enum.StoreRankSuccess).Get(&store)
	if err != nil {
		return store, err
	}
	return store, nil
}

func GetStoreByUserIdAndStoreId(engine *database.MysqlSession, UserId, StoreId string) (entity.StoreDataResp, error) {
	var store entity.StoreDataResp
	sql := fmt.Sprintf("SELECT * FROM store_data s LEFT JOIN store_rank_data r ON s.store_id = r.store_id" +
		" WHERE s.store_id = ? AND r.user_id = ? ORDER BY rank ASC, store_default DESC")
	_, err := engine.Engine.SQL(sql, StoreId, UserId).Get(&store)
	if err != nil {
		return store, err
	}
	return store, nil
}

//更新Order Data
func UpdateStoreData(engine *database.MysqlSession, storeId string, storeData entity.StoreData) error {
	storeData.UpdateTime = time.Now()
	_, err := engine.Session.Table(entity.StoreData{}).ID(storeId).Update(storeData)
	if err != nil {
		return err
	}
	return nil
}

// 新增 Insert Member
func InsertStoreData(engine *database.MysqlSession, uid, storeName, picture, expire string) (entity.StoreData, error) {
	var storeData entity.StoreData
	storeData.StoreId = tools.GeneratorStoreId()
	storeData.SellerId = uid
	if len(expire) == 0 {
		storeData.StoreDefault = 1
	}
	storeData.StoreName = storeName
	storeData.StoreStatus = Enum.StoreStatusSuccess
	storeData.StorePicture = picture
	storeData.ExpireTime = expire
	storeData.CreateTime = time.Now()
	storeData.UpdateTime = time.Now()
	_, err := engine.Session.Table("store_data").Insert(&storeData)
	if err != nil {
		log.Error("Database Error", err)
		return storeData, err
	}
	return storeData, nil
}

func CountStoreByUserId(engine *database.MysqlSession, UserId string) int64 {
	sql := fmt.Sprintf("SELECT count(*) FROM store_data WHERE seller_id = ?")
	result, err := engine.Engine.SQL(sql, UserId).Count()
	if err != nil {
		log.Error("Count Store Manager Error", err)
		return 0
	}
	return result
}

func CountStoreManager(engine *database.MysqlSession, StoreId string) int64 {
	sql := fmt.Sprintf("SELECT count(*) FROM store_rank_data WHERE store_id = ? AND rank = ? AND rank_status != ? AND rank_status != ?")
	result, err := engine.Engine.SQL(sql, StoreId, Enum.StoreRankSlave, Enum.StoreRankDelete, Enum.StoreRankSuspend).Count()
	if err != nil {
		log.Error("Count Store Manager Error", err)
		return 0
	}
	return result
}

func GetUserAllStoreIdStore(engine *database.MysqlSession, UserId string) ([]entity.StoreData, error) {
	var data []entity.StoreData
	sql := fmt.Sprintf("SELECT * FROM store_data WHERE seller_id = ?")
	err := engine.Engine.SQL(sql, UserId).Find(&data)
	if err != nil {
		return data, err
	}
	return data, nil
}

func GetStoreManagerListByStoreId(engine *database.MysqlSession, StoreId string) ([]entity.StoreRankResp, error) {
	var data []entity.StoreRankResp
	sql := fmt.Sprintf("SELECT * FROM store_rank_data s LEFT JOIN member_data m ON s.user_id = m.uid" +
		" WHERE s.store_id = ? AND s.rank = ? AND rank_status != ? AND rank_status != ? ORDER BY rank_id ASC")
	err := engine.Engine.SQL(sql, StoreId, Enum.StoreRankSlave, Enum.StoreRankDelete, Enum.StoreRankSuspend).Find(&data)
	if err != nil {
		return data, err
	}
	return data, nil
}

func GetStoreManagerAndMemberByManagerId(engine *database.MysqlSession, ManagerId int) (entity.StoreRankResp, error) {
	var data entity.StoreRankResp
	sql := fmt.Sprintf("SELECT * FROM store_rank_data s LEFT JOIN member_data m ON s.user_id = m.uid" +
		" WHERE rank_id = ?")
	_, err := engine.Engine.SQL(sql, ManagerId).Get(&data)
	if err != nil {
		return data, err
	}
	return data, nil
}

//更新
func UpdateStoreManagerData(engine *database.MysqlSession, data entity.StoreRankData) error {
	data.UpdateTime = time.Now()
	_, err := engine.Session.Table("store_rank_data").ID(data.RankId).Update(data)
	if err != nil {
		return err
	}
	return nil
}

func UpdateStoreVerifyIdentity(engine *database.MysqlSession, uid string, is int64) error {
	now := time.Now()
	sql := fmt.Sprintf("UPDATE store_data SET verify_identity = ?, update_time = ? WHERE seller_id = ?")
	_, err := engine.Session.Exec(sql, is, now, uid)
	if err != nil {
		log.Error("UpdateOrderMessageBoardData Error", err)
		return err
	}
	return nil
}

func GetStoreDataBySellerIdAndStoreId(engine *database.MysqlSession, UserId, StoreId string) (entity.StoreUserData, error) {
	var store entity.StoreUserData
	sql := fmt.Sprintf("SELECT * FROM store_data s LEFT JOIN member_data m ON s.seller_id = m.uid" +
		" WHERE s.store_id = ? AND s.seller_id = ?")
	_, err := engine.Engine.SQL(sql, StoreId, UserId).Get(&store)
	if err != nil {
		return store, err
	}
	return store, nil
}

func GetStoreBySellerId(engine *database.MysqlSession, UserId string) ([]entity.StoreData, error) {
	var data []entity.StoreData
	if err := engine.Engine.Table(entity.StoreData{}).Select("*").
		Where("seller_id = ?", UserId).Find(&data); err != nil {
		return data, err
	}
	return data, nil
}

func GetStoreFreeShipByStoreId(engine *database.MysqlSession, storeId string) (entity.StoreData, error) {
	var store entity.StoreData
	_, err := engine.Engine.Table(entity.StoreData{}).Select("*").Where("store_id = ?", storeId).Get(&store)
	if err != nil {
		return store, err
	}
	return store, nil
}

func GetStoresByStoreName(engine *database.MysqlSession, storeName string) ([]entity.StoreData, error) {
	var store []entity.StoreData
	err := engine.Engine.Table(entity.StoreData{}).Select("*").
		Where("store_name = ?", storeName).Find(&store)
	if err != nil {
		return store, err
	}
	return store, nil
}

func GetStoreManagersByStoreId(engine *database.MysqlSession, StoreId string) ([]entity.StoreRankData, error) {
	var data []entity.StoreRankData
	if err := engine.Engine.Table(entity.StoreRankData{}).Select("*").Where("store_id=? AND rank =?", StoreId, Enum.StoreRankSlave).Find(&data); err != nil {
		return data, err
	}
	return data, nil
}
func GetStoreSlaveManagerByUserId(engine *database.MysqlSession, stores []string, userId string) ([]entity.StoreWithRank, error) {
	var d []entity.StoreWithRank
	if err := engine.Engine.Table(entity.StoreRankData{}).Select("*").
		Join("LEFT", entity.StoreData{}, "store_rank_data.store_id = store_data.store_id").
		NotIn("store_rank_data.store_id", stores).Where("store_rank_data.user_id = ?", userId).Find(&d); err != nil {
		return d, err
	}
	return d, nil
}

//取出所有管理者
func GetStoreAllManagerListBySellerId(engine *database.MysqlSession, SellerId string) ([]entity.StoreAndStoreRankEnt, error) {
	var data []entity.StoreAndStoreRankEnt
	if err := engine.Engine.Table(entity.StoreData{}).Select("*").
		Join("LEFT", entity.StoreRankData{}, "store_data.store_id = store_rank_data.store_id").
	    Where("store_data.seller_id = ?", SellerId).
		And( "store_rank_data.rank = ?", Enum.StoreRankSlave).
		And("store_rank_data.rank_status != ?", Enum.StoreRankDelete).
		And("store_rank_data.rank_status != ?", Enum.StoreRankSuspend).Find(&data); err != nil {
		return data, err
	}
	return data, nil
}

func GetUserAllStoreData(engine *database.MysqlSession) ([]entity.StoreData, error) {
	var data []entity.StoreData
	err := engine.Engine.Table(entity.StoreData{}).Select("*").Asc("create_time").Find(&data)
	if err != nil {
		return data, err
	}
	return data, nil
}

func UpdateStoreSelfDeliveryArea(engine *database.MysqlSession, storeId string, enable bool) error {
	_, err := engine.Session.Table(entity.StoreData{}).Cols("self_delivery").Where("store_id=?", storeId).Update(entity.StoreData{
		SelfDelivery: enable,
	})
	if err != nil {
		log.Error(err.Error(), "Update store self-delivery error")
		return err
	}
	return nil
}

func InsertOrUpdateStoreSelfDeliveryArea(engine *database.MysqlSession, storeId, cityCode string, areaList []string) error {
	var data entity.StoreSelfDeliveryArea
	flag, err := engine.Engine.Table(entity.StoreSelfDeliveryArea{}).Where("store_id=? AND city_code=?", storeId, cityCode).Get(&data)
	if err != nil {
		log.Error(err.Error())
		return err
	}
	if flag {
		_, err := engine.Session.Table(entity.StoreSelfDeliveryArea{}).Where("store_id=? AND city_code=?", storeId, cityCode).Cols("area_zip").Update(entity.StoreSelfDeliveryArea{
			AreaZip: strings.Join(areaList, ","),
		})
		if err != nil {
			log.Error(err.Error())
			return err
		}
	} else {
		var data entity.StoreSelfDeliveryArea
		data.AreaZip = strings.Join(areaList, ",")
		data.CityCode = cityCode
		data.StoreId = storeId
		_, err := engine.Session.Table(entity.StoreSelfDeliveryArea{}).Insert(&data)
		if err != nil {
			log.Error(err.Error())
			return err
		}
	}
	return nil
}

func UpdateStoreSelfDeliveryChargeFree(engine *database.MysqlSession, storeId string, fee int64, key string) error {
	var data entity.StoreData
	cols := []string{"free_self_delivery"}
	if len(key) > 0 {
		cols = append(cols, "self_delivery_key")
		data.SelfDeliveryKey = key
	}
	data.FreeSelfDelivery = fee
	_, err := engine.Session.Table(entity.StoreData{}).Cols(cols...).Where("store_id=?", storeId).Update(&data)
	if err != nil {
		log.Error(err.Error())
		return err
	}
	return nil
}

func GetStoreSelfDeliveryArea(engine *database.MysqlSession, storeId string) ([]entity.StoreSelfDeliveryArea, error) {
	var data []entity.StoreSelfDeliveryArea
	err := engine.Engine.Table(entity.StoreSelfDeliveryArea{}).Where("store_id=?", storeId).Find(&data)
	if err != nil {
		log.Error(err.Error())
		return data, err
	}

	return data, nil
}

func DeleteStoreSelfDeliveryArea(engine *database.MysqlSession, storeId string) error {

	_, err := engine.Session.Table(entity.StoreSelfDeliveryArea{}).Delete(entity.StoreSelfDeliveryArea{
		StoreId: storeId,
	})
	if err != nil {
		log.Error(err.Error())
		return err
	}
	return nil
}
func CheckStoreManager(engine *database.MysqlSession, userId string, storeId string) (bool, error) {
	var data entity.StoreRankData
	flag, err := engine.Engine.Table(entity.StoreRankData{}).
		Where("store_id =? AND user_id =? AND rank_status =?", storeId, userId, Enum.StoreRankSuccess).Get(&data)
	if err != nil {
		log.Error("Find store manager error", err.Error())
		return false, err
	}
	return flag, nil
}

func UpdateStorePromoEnable(engine *database.MysqlSession, storeId string, enable bool) error {
	_, err := engine.Session.Table(entity.StoreData{}).
		Where("store_id=?", storeId).Cols("enable_promo").Update(entity.StoreData{
		EnablePromo: enable,
	})
	if err != nil {
		log.Error("Update Store Promo Enable", err.Error())
		return err
	}
	return nil
}
