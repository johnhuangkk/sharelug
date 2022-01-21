package model

import (
	"api/services/VO/Response"
	"api/services/dao/Area"
	"api/services/util/log"
	"encoding/json"
	"io/ioutil"
)

type Zipcode struct {
	Version  int64    `json:"Version"`
	CityList []Cities `json:"CityList"`
}

type Cities struct {
	Name         string   `json:"Name"`
	DistrictList []Region `json:"DistrictList"`
}

type Region struct {
	Name string `json:"Name"`
	Code string `json:"Code"`
}

var zipcode Zipcode

/**
 * 取ZIPCODE
 */
func GetZipcode() (Zipcode, error) {
	file, err := openFile("./data/zipcode/zipcode.json")
	if err != nil {
		log.Error("OpenFile Error", err)
		return zipcode, err
	}
	json.Unmarshal(file, &zipcode)
	return zipcode, nil
}

/**
 * 開啟檔案
 */
func openFile(fileName string) ([]byte, error) {
	file, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Error("Open File Error", err)
		return nil, err
	}
	return file, nil
}

func GetTaiwanArea(cityCode string) (Response.CityWithArea, error) {
	var resp Response.CityWithArea
	data, err := Area.GetTaiwanArea(cityCode)
	if err != nil {
		return resp, err
	}
	if len(data) > 0 {
		resp.CityCode = cityCode
		resp.CityName = data[0].CityName
		for _, ta := range data {
			var area Response.AreaFlagWithZipCode
			area.Name = ta.AreaName
			area.ZipCode = ta.ZipCode
			resp.Area = append(resp.Area, area)
		}
	}
	return resp, nil
}

func GetTaiwanCities() ([]Response.CityWithCode, error) {
	var resp []Response.CityWithCode
	data, err := Area.GetTaiwanCities()
	if err != nil {
		return resp, err
	}
	if len(data) > 0 {
		for _, ct := range data {
			var rcc Response.CityWithCode
			rcc.CityCode = ct.CityCode
			rcc.CityName = ct.CityName
			resp = append(resp, rcc)
		}
	}
	return resp, nil
}
func CreateTaiwanArea() error {
	data, err := GetZipcode()
	if err != nil {
		return err
	}
	for _, city := range data.CityList {
		for _, dis := range city.DistrictList {
			err := Area.InsertTaiwanArea(city.Name, dis.Code, dis.Name)
			if err != nil {
				log.Error(err.Error())
				continue
			}
		}
	}
	return nil
}
