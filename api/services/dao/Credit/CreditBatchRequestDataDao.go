package Credit

import (
	"api/services/database"
	"api/services/entity"
	"api/services/util/log"
	"time"
)

// 新增
func InsertCreditBatchRequestData(engine *database.MysqlSession, data entity.CreditBatchRequestData) (entity.CreditBatchRequestData, error) {
	data.CreateTime = time.Now()
	data.UpdateTime = time.Now()
	_, err := engine.Session.Table("credit_batch_request_data").Insert(&data)
	if err != nil {
		log.Error("Database Error", err)
		return data, err
	}
	return data, nil
}

// 更新
func UpdateCreditBatchRequestData(engine *database.MysqlSession, BatchId string, data entity.CreditBatchRequestData) (int64, error) {
	affected, err := engine.Session.Table("credit_batch_request_data").Where("batch_id = ?", BatchId).Update(data)
	if err != nil {
		return affected, err
	}
	return affected, nil
}