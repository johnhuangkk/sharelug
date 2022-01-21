package Area

import (
	"api/services/database"
	"api/services/entity"
	"api/services/util/log"
	"api/services/util/tools"
)

func InsertTaiwanArea(cityName string, zipcode string, areaName string) error {
	engine := database.GetMysqlEngine()
	defer engine.Close()
	var existData entity.TaiwanArea
	var newData entity.TaiwanArea
	newData.AreaName = areaName
	newData.CityName = cityName
	newData.ZipCode = zipcode
	flag, err := engine.Engine.Table(entity.TaiwanArea{}).Where("city_name = ?", cityName).Get(&existData)
	if err != nil {
		log.Error(err.Error())
		return err
	}
	if flag {
		newData.CityCode = existData.CityCode
	} else {
		newData.CityCode = tools.GenerateTaiwanAreaId()
	}
	_, err = engine.Session.Table(entity.TaiwanArea{}).Insert(&newData)
	if err != nil {
		log.Error(err.Error())
		return err
	}
	return nil
}

func GetTaiwanArea(cityCode string) ([]entity.TaiwanArea, error) {
	engine := database.GetMysqlEngine()
	defer engine.Close()
	var data []entity.TaiwanArea
	err := engine.Engine.Table(entity.TaiwanArea{}).Where("city_code=?", cityCode).Find(&data)
	if err != nil {
		log.Error(err.Error())
		return data, err
	}
	return data, nil
}

func GetTaiwanCities() ([]entity.TaiwanArea, error) {
	engine := database.GetMysqlEngine()
	defer engine.Close()
	var data []entity.TaiwanArea
	err := engine.Engine.Table(entity.TaiwanArea{}).Distinct("city_code,city_name").Find(&data)
	if err != nil {
		log.Error(err.Error())
		return data, err
	}
	return data, nil
}
