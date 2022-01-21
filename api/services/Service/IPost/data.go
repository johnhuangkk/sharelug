package IPost

import (
	"api/services/VO/IPOSTVO"
	"api/services/entity"
	"api/services/model"
	"api/services/util/curl"
	"api/services/util/export"
	"api/services/util/log"
	"encoding/json"
	"github.com/spf13/viper"
	"os"
	"time"
)

// 呼叫郵局 i郵箱所有資訊
func UpdateIPostDataTask() {
	IPostBoxArray := entity.IPostBoxArray{}
	path := viper.GetString("IPOST.iPostInfoPath")

	log.Info("UpdateIPostDataTask [start]")
	body, _ := curl.Post(path, "")

	_ = json.Unmarshal(body, &IPostBoxArray)

	// 寫入資料
	err := model.SetupPostBoxData(IPostBoxArray)

	if err != nil {
		log.Error("i郵箱資料更新失敗")
		return
	}
	// 處理匯出DB i郵箱資料
	exportPostData()
	log.Info("UpdateIPostDataTask [End]")
}

// 處理匯出DB i郵箱資料
func exportPostData() {
	log.Info("ExportPostData", "輸出中....")
	var PostBoxDataArray = model.GetPostBoxStatusY()
	if len(PostBoxDataArray) == 0 {
		renameYesterdayFile()
		log.Error("func %s ", " i郵箱資料異常")
		return
	}

	var posts []IPOSTVO.IPostJson
	var post = &IPOSTVO.IPostJson{}
	var city = &IPOSTVO.City{}
	var citys []IPOSTVO.City
	var adr = &IPOSTVO.Address{}

	for _, v := range PostBoxDataArray {
		adr.Location = v.AdmLocation
		adr.Name = v.Address
		adr.Id = v.AdmId
		adr.Alias = v.AdmAlias

		if city.Name !=  v.City {
			if len(city.Name) != 0 {
				citys = append(citys, *city)
				city = &IPOSTVO.City{}
			}
		}

		city.Country = v.Country
		city.Name = v.City
		city.Zip = v.Zip
		city.Address = append(city.Address, *adr)
	}

	// 最後縣市不會變更須在此插入尾端
	citys = append(citys, *city)

	for _, c := range citys {
		if post.Country !=  c.Country {
			if len(post.Country) != 0 {
				posts = append(posts, *post)
				post = &IPOSTVO.IPostJson{}
			}
		}

		post.Country = c.Country
		post.City = append(post.City, c)
	}

	// 最後縣市不會變更須在此插入尾端
	posts = append(posts, *post)

	fileName := "IPOST_" + time.Now().Format("20060102") + ".json"
	//
	err := export.JsonToData(posts, fileName)

	// i郵箱產生失敗將昨日資訊改為今日資訊
	if err != nil {
		log.Error("JsonToData", err)
		renameYesterdayFile()
		return
	}
}

// 將前一天的資料改為今日資料
func renameYesterdayFile() {
	var config = viper.GetStringMapString("Data")
	var todayFileName = "IPOST_" + time.Now().Format("20060102") + ".json"
	var yesterdayFileName = "IPOST_" + time.Now().Add(-(time.Hour * time.Duration(24))).Format("20060102") + ".json"

	var oldSource = config["gopath"] + yesterdayFileName
	var source = config["gopath"] + todayFileName

	var oldDestination = config["wwwpath"] + yesterdayFileName
	var destination = config["wwwpath"] + todayFileName

	err := os.Rename(oldSource, source)
	if err != nil {
		log.Error("Rename source fail", err)
	}
	err = os.Rename(oldDestination, destination)
	if err != nil {
		log.Error("Rename destination fail", err)
	}
}
