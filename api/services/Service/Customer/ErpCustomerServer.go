package Customer

import (
	"api/services/VO/CustomerVo"
	"api/services/dao/Customer"
	"api/services/database"
	"api/services/util/log"
	"time"
)

//取出客服內容
func GetOrderCustomer(engine *database.MysqlSession, orderId string) ([]CustomerVo.CustomerResponse, error) {
	var resp []CustomerVo.CustomerResponse
	data, err := Customer.FindCustomerByOrderId(engine, orderId)
	if err != nil {
		log.Error("Get Customer Error", err)
		return resp, err
	}
	for _, v := range data {
		var res CustomerVo.CustomerResponse
		res = v.GetCustomerData()
		resp = append(resp, res)
	}
	return resp, nil
}

//取出客服備註
func GetOrderCustomerRemark(engine *database.MysqlSession, orderId string) error {

	return nil
}

//客服回覆
func AnswerOrderCustomer(engine *database.MysqlSession, orderId, title, contents string) error {
	data, err := Customer.GetCustomerByOrderId(engine, orderId)
	if err != nil {
		log.Error("Get Customer Error", err)
		return err
	}
	if len(data.OrderId) == 0 {
		return err
	}
	data.ReplyTitle = title
	data.Reply = contents
	data.ReplyTime = time.Now()
	if err := Customer.UpdateCustomerData(engine, data); err != nil {
		log.Error("Update Customer Error", err)
		return err
	}
	return nil
}

//客服備註
func RemarkOrderCustomer(engine *database.MysqlSession, orderId, contents string) error {
	data, err := Customer.GetCustomerByOrderId(engine, orderId)
	if err != nil {
		log.Error("Get Customer Error", err)
		return err
	}
	if len(data.OrderId) == 0 {
		return err
	}
	data.Remark = contents
	data.RemarkTime = time.Now()
	data.RemarkStaff = ""
	if err := Customer.UpdateCustomerData(engine, data); err != nil {
		log.Error("Update Customer Error", err)
		return err
	}
	return nil
}
