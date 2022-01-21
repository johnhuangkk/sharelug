package Customer

import (
	"api/services/Enum"
	"api/services/database"
	"api/services/entity"
	"api/services/util/log"
	"fmt"
	"time"
)

func InsertCustomerQuestionData(engine *database.MysqlSession, data entity.CustomerQuestionData) error {
	data.CreateTime = time.Now()
	_, err := engine.Session.Table(entity.CustomerQuestionData{}).Insert(&data)
	if err != nil {
		log.Error("Insert Customer Question Database Error", err)
		return err
	}
	return nil
}

func CountCustomerQuestionData(engine *database.MysqlSession) (int64, error) {
	sql := fmt.Sprintf("SELECT count(*) FROM customer_question_data")
	result, err := engine.Engine.SQL(sql).Count()
	if err != nil {
		log.Error("Count Customer Question Database Error", err)
		return 0, err
	}
	return result, nil
}

func GetCustomerQuestionData(engine *database.MysqlSession) ([]entity.CustomerQuestionData, error) {
	var data []entity.CustomerQuestionData
	sql := fmt.Sprintf("SELECT * FROM customer_question_data ORDER BY sort ASC")
	err := engine.Engine.SQL(sql).Find(&data)
	if err != nil {
		log.Error("Get Customer Question Database Error", err)
		return data, err
	}
	return data, nil
}

func InsertCustomerData(engine *database.MysqlSession, data entity.CustomerData) error {
	data.CreateTime = time.Now()
	data.UpdateTime = time.Now()
	_, err := engine.Session.Table(entity.CustomerData{}).Insert(&data)
	if err != nil {
		log.Error("Insert Customer Question Database Error", err)
		return err
	}
	return nil
}

func FindCustomerByOrderId(engine *database.MysqlSession, OrderId string) ([]entity.CustomerData, error) {
	var data []entity.CustomerData
	if err := engine.Engine.Table(entity.CustomerData{}).Select("*").Where("order_id = ?", OrderId).
		Desc("create_time").Find(&data); err != nil {
		log.Error("Get Customer Database Error", err)
		return data, err
	}
	return data, nil
}
func GetCustomerQuestionByRelatedIdWithTime(engine *database.MysqlSession, relatedId string, startTime time.Time, msgType string) ([]entity.CustomerData, error) {
	var data []entity.CustomerData

	query := engine.Engine.Table(entity.CustomerData{}).Where(`related_id = ?`, relatedId)
	switch msgType {
	case Enum.NotifyTypePlatformService:
		if err := query.Where(`customer_data.reply_time < ?`, startTime.Format(`2006-01-02 15:04:05`)).Desc("create_time").Find(&data); err != nil {
			return data, err
		}
	case Enum.NotifyTypePlatformUser:
		if err := query.Where(`customer_data.create_time < ?`, startTime.Format(`2006-01-02 15:04:05`)).Desc("create_time").Find(&data); err != nil {
			return data, err
		}

	}
	return data, nil
}
func GetCustomerQuestionByRelatedIdWithoutType(engine *database.MysqlSession, relatedId string, startTime time.Time) ([]entity.CustomerData, error) {
	var data []entity.CustomerData

	query := engine.Engine.Table(entity.CustomerData{}).Where(`related_id = ?`, relatedId)
	if err := query.Where(`customer_data.create_time < ?`, startTime.Format(`2006-01-02 15:04:05`)).Desc("create_time").Find(&data); err != nil {
		return data, err
	}
	return data, nil
}
func GetCustomerQuestionByQuestionId(engine *database.MysqlSession, questionId string) (entity.CustomerData, error) {
	var data entity.CustomerData
	if _, err := engine.Engine.Table(entity.CustomerData{}).Where(`question_id = ?`, questionId).Get(&data); err != nil {
		return data, err
	}
	return data, nil
}
func GetCustomerQuestionById(engine *database.MysqlSession, questionId string) (entity.CustomerData, error) {
	var data entity.CustomerData
	if _, err := engine.Engine.Table(entity.CustomerData{}).Where(`id = ?`, questionId).Get(&data); err != nil {
		return data, err
	}
	return data, nil
}
func GetCustomerByOrderId(engine *database.MysqlSession, OrderId string) (entity.CustomerData, error) {
	var data entity.CustomerData
	if _, err := engine.Engine.Table(entity.CustomerData{}).Select("*").Where("order_id = ?", OrderId).
		Desc("create_time").Get(&data); err != nil {
		log.Error("Get Customer Database Error", err)
		return data, err
	}
	return data, nil
}

func UpdateCustomerData(engine *database.MysqlSession, data entity.CustomerData) error {
	if _, err := engine.Session.Table(entity.CustomerData{}).Where("id = ?", data.Id).Update(data); err != nil {
		return err
	}
	return nil
}

func InsertContactData(engine *database.MysqlSession, data entity.ContactData) error {
	_, err := engine.Session.Table(entity.ContactData{}).Insert(data)
	if err != nil {
		log.Error("Insert Contact Error", err)
		return err
	}
	return nil
}

func GetCustomerById(engine *database.MysqlSession, id string) (entity.CustomerData, error) {
	var data entity.CustomerData
	if _, err := engine.Engine.Table(entity.CustomerData{}).Select("*").Where("question_id = ?", id).Get(&data); err != nil {
		log.Error("Get Customer Database Error", err)
		return data, err
	}
	return data, nil
}
