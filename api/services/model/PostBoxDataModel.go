package model

import (
	"api/services/VO/IPOSTVO"
	"api/services/dao/iPost"
	"api/services/database"
	"api/services/entity"
	"api/services/util/log"
)

func insertIPOSTBoxData(IPOSTBoxData entity.IPOSTBoxData) entity.PostBoxData {
	var ent entity.PostBoxData

	ent.AdmId = IPOSTBoxData.ADMId
	ent.AdmName = IPOSTBoxData.ADMName
	ent.AdmAlias = IPOSTBoxData.ADMName
	ent.AdmLocation = IPOSTBoxData.ADMLocation
	ent.Country = IPOSTBoxData.Country
	ent.Zip = IPOSTBoxData.Zip
	ent.City = IPOSTBoxData.City
	ent.Address = IPOSTBoxData.Address
	ent.Longitude = IPOSTBoxData.Longitude
	ent.Latitude = IPOSTBoxData.Latitude
	ent.GovNo = IPOSTBoxData.GovNo
	ent.BoxStatus = "Y"

	return ent
}

func update(data entity.PostBoxData, IPostBox entity.IPOSTBoxData) entity.PostBoxData {

	data.AdmId = IPostBox.ADMId
	data.AdmName = IPostBox.ADMName
	data.AdmAlias = IPostBox.ADMName
	data.AdmLocation = IPostBox.ADMLocation
	data.Country = IPostBox.Country
	data.Zip = IPostBox.Zip
	data.City = IPostBox.City
	data.Address = IPostBox.Address
	data.Longitude = IPostBox.Longitude
	data.Latitude = IPostBox.Latitude
	data.GovNo = IPostBox.GovNo
	data.BoxStatus = "Y"

	return data
}

// 使用櫃體編號Id 取得地址
func GetPostBoxAddressById(engine *database.MysqlSession, admId string) IPOSTVO.IPostZipAddress {
	ip, _ := iPost.SelectPostBoxByAdmId(engine, admId)
	return IPOSTVO.IPostZipAddress{
		Id: ip.AdmId,
		Alias: ip.AdmAlias,
		Zip:     ip.Zip,
		Address: ip.AdmAlias + "," + ip.Zip + ip.Country + ip.City + ip.Address + " [" + ip.AdmLocation + "]",
		Status: ip.BoxStatus,
	}
}

/**
寫入 i郵箱資料
*/
func SetupPostBoxData(IPostBoxArray entity.IPostBoxArray) error {
	engine := database.GetMysqlEngine()
	defer engine.Close()

	err := iPost.SetPostBoxStatusClose(engine)
	if err != nil {
		return err
	}

	for _, IPOSTData := range IPostBoxArray {
		data, _ := iPost.SelectPostBoxByAdmId(engine, IPOSTData.ADMId)

		if len(data.AdmId) == 0 {
			ent := insertIPOSTBoxData(IPOSTData)
			_, err := iPost.InsertPostBox(engine, ent)
			if err != nil {
				log.Error("IPOSTBoxData IPOSTData %v", IPOSTData)
				log.Error("IPOSTBoxData select %v", data)
				log.Error("IPOSTBoxData insert %v", ent)
				//return fmt.Errorf("IPOSTBoxData Insert Fail")
			}
		} else {
			ent := update(data, IPOSTData)
			_, err := iPost.UpdatePostBox(engine, ent.AdmId, ent)

			if err != nil {
				log.Error("IPOSTBoxData update %v", data)
				//return fmt.Errorf("IPOSTBoxData Update Fail")
			}
		}
	}

	return nil
}

/*
取得 i郵箱資料資料
*/
func GetPostBoxStatusY() []entity.PostBoxData {
	engine := database.GetMysqlEngine()
	defer engine.Close()
	const STATUS = "Y"
	var dataArray []entity.PostBoxData
	var data []entity.PostBoxData

	city, _ := iPost.SelectDistinctCountry(engine)
	for _, c := range city {
		data, _ = iPost.SelectPostBoxByBoxByStatusAndCountry(engine, STATUS, c["country"])
		dataArray = append(dataArray, data...)
	}

	return dataArray
}
