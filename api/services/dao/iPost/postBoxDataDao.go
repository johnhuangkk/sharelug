package iPost

import (
	"api/services/database"
	"api/services/entity"
	"api/services/util/log"
	"time"
)

func UpdatePostBox(engine *database.MysqlSession, ADMId string, data entity.PostBoxData) (entity.PostBoxData, error) {
	data.UpdateTime = time.Now()

	_, err := engine.Session.Table(entity.PostBoxData{}).ID(ADMId).Update(&data)

	if err != nil {
		log.Error("Database Error", err)
		return data, err
	}
	return data, nil
}


func InsertPostBox(engine *database.MysqlSession, data entity.PostBoxData) (entity.PostBoxData, error) {
	data.UpdateTime = time.Now()
	_, err := engine.Session.Table(entity.PostBoxData{}).Insert(&data)
	if err != nil {
		log.Error("Database Error", err)
		return data, err
	}
	return data, nil
}

/*
	透過 AdmId 找到單筆資料
*/
func SelectPostBoxByAdmId(engine *database.MysqlSession, AdmId string) (entity.PostBoxData, error) {
	var data entity.PostBoxData
	_, err := engine.Engine.Table("post_box_data").Select("*").Where("adm_id = ?", AdmId).Get(&data)
	if err != nil {
		return data, err
	}
	return data, nil
}

/*
	取得縣市唯一
 */
func SelectDistinctCountry(engine *database.MysqlSession) ([]map[string]string, error) {
	data, err := engine.Engine.QueryString("select distinct (country) from post_box_data ")
	if err != nil {
		log.Error("SelectDistinctCountry Error", err)
		return data, err
	}
	return data, nil
}

/*
透過狀態以及城市找尋 i郵箱資訊
 */
func SelectPostBoxByBoxByStatusAndCountry(engine *database.MysqlSession, status string, country string) ([]entity.PostBoxData, error) {
	var data []entity.PostBoxData
	query := map[string]interface{}{}
	query["box_status"] = status
	query["country"] = country

	err := engine.Engine.Table(entity.PostBoxData{}).Select("*").Where(query).OrderBy("zip asc, city asc").Find(&data)

	if err != nil {
		return data, err
	}

	return data, nil
}

/*
關閉全部收件i郵箱位置
 */
func SetPostBoxStatusClose(engine *database.MysqlSession) error {
	var sql = "UPDATE post_box_data set box_status = 'N' "
	_, err := engine.Session.Exec(sql)
	if err != nil {
		log.Error("SetPostBoxStatusClose Error", err)
		return err
	}
	return nil
}
