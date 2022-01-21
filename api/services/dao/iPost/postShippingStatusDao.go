package iPost

import (
	"api/services/database"
	"api/services/entity"
	"api/services/util/log"
	"time"
)

/**

 */
func InsertPostShippingStatus(engine *database.MysqlSession, data entity.PostShippingStatus) error {
	_, err := engine.Session.Table(entity.PostShippingStatus{}).Insert(data)
	if err != nil {
		log.Error("Database Error", err)
		return err
	}
	return nil
}

/**
尋找物流配送狀態
*/
func QueryPostShippingStatus(engine *database.MysqlSession, shipNumber string, createTime time.Time) ([]entity.PostShippingStatus, error) {
	var data []entity.PostShippingStatus

	err := engine.Engine.Table(entity.PostShippingStatus{}).Select("*").
		Where(
			"mail_no = ? And ( handle_time < ?  And  handle_time > ?)",
			shipNumber,
			createTime.AddDate(0, 3, 0).Format("2006-01-02 15:04"),
			createTime.Format("2006-01-02 15:04")).
		Desc("handle_time").
		Find(&data)
	if err != nil {
		log.Error("QueryPostShippingStatus Error", err)
		return data, err
	}

	return data, nil
}

func InsertPostBagShippingStatus(engine *database.MysqlSession, data entity.PostShippingStatus) error {
	var og []entity.PostShippingStatus
	err := engine.Session.Table(entity.PostShippingStatus{}).Where(`status_code =?`, data.StatusCode).Where(`mail_no =?`, data.MailNo).Find(&og)
	if err != nil {
		log.Error("Database Error", err)
		return err
	}
	if len(og) == 0 {
		_, err = engine.Session.Table(entity.PostShippingStatus{}).Insert(data)
		if err != nil {
			log.Error("Database Error", err)
			return err
		}
	}

	return nil
}
