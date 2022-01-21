package StoreService

import (
	"api/services/Enum"
	"api/services/Service/MemberService"
	"api/services/VO/Response"
	"api/services/dao/Store"
	"api/services/dao/member"
	"api/services/database"
	"api/services/entity"
	"fmt"
)

//建立管理帳號
func CreateStoreManageData(engine *database.MysqlSession, storeId, phone, email string) (entity.MemberData, error) {
	//檢查是否有帳號
	userData, err := member.GetMemberDataByPhone(engine, phone)
	if err != nil {
		return userData, fmt.Errorf("系統錯誤")
	}
	if len(userData.Uid) == 0 {
		userData, err = MemberService.BuildNewMember(engine, phone)
		if err != nil {
			return userData, fmt.Errorf("系統錯誤")
		}
	}
	storeUser, err := Store.GetStoreManagerByStoreIdAndUserId(engine, storeId, userData.Uid)
	if err != nil {
		return userData, fmt.Errorf("系統錯誤")
	}
	if len(storeUser.UserId) != 0 && storeUser.RankStatus != Enum.StoreRankDelete {
		return userData, fmt.Errorf("此電話號碼已在管理員列表中，無法重複新增管理員。")
	}
	if len(storeUser.UserId) == 0 {
		_, err = Store.InsertStoreRankData(engine, userData.Uid, storeId, Enum.StoreRankSlave, Enum.StoreRankInit, email)
		if err != nil {
			return userData, fmt.Errorf("系統錯誤")
		}
	} else {
		storeUser.RankStatus = Enum.StoreRankInit
		storeUser.Email = email
		if _, err := Store.UpdateStoreRankData(engine, storeUser); err != nil {
			return userData, fmt.Errorf("系統錯誤")
		}
	}
	return userData, nil
}

//取出使用者的所有StoreId
func GetUserAllStore(engine *database.MysqlSession, UserId string) ([]string, error) {
	data, err := Store.GetUserAllStoreIdStore(engine, UserId)
	if err != nil {
		return nil, fmt.Errorf("系統錯誤")
	}
	var store []string
	for _, v := range data {
		store = append(store, v.StoreId)
	}
	return store, nil
}

//建立賣場資料
func CreateStoreData(engine *database.MysqlSession, uid, storeName, picture, expire string) (entity.StoreData, error) {
	storeData, err := Store.InsertStoreData(engine, uid, storeName, picture, expire)
	if err != nil {
		return storeData, err
	}
	//建立管理帳號
	_, err = CreateStoreManageDefaultData(engine, uid, storeData.StoreId)
	if err != nil {
		return storeData, err
	}
	return storeData, nil
}

//建立預設管理帳號
func CreateStoreManageDefaultData(engine *database.MysqlSession, uid, storeId string) (entity.StoreRankData, error) {
	storeRankData, err := Store.InsertStoreRankData(engine, uid, storeId, Enum.StoreRankMaster, Enum.StoreRankSuccess, "")
	if err != nil {
		return storeRankData, err
	}
	return storeRankData, nil
}

//取出免運設定
func GetStoreFreeShipping(engine *database.MysqlSession, storeId string) (Response.StoreFreeShipResponse, error) {
	var resp Response.StoreFreeShipResponse
	data, err := Store.GetStoreFreeShipByStoreId(engine, storeId)
	if err != nil {
		return resp, err
	}
	resp.FreeShipKey = data.FreeShipKey
	resp.FreeShip = data.FreeShip
	resp.SelfDelivery = data.SelfDelivery
	resp.SelfDeliveryFree = data.FreeSelfDelivery
	resp.SelfDeliveryKey = data.SelfDeliveryKey
	resp.IsCoupon = data.EnablePromo
	switch data.FreeShipKey {
	case Enum.FreeShipAmount:
		resp.FreeShipText = fmt.Sprintf("滿 $%v 免運", data.FreeShip)
	case Enum.FreeShipQuantity:
		resp.FreeShipText = fmt.Sprintf("%v 件免運", data.FreeShip)
	}
	return resp, nil
}
