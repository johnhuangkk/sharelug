package Orders

import (
	"api/services/Enum"
	"api/services/database"
	"api/services/entity"
	"api/services/util/log"
	"fmt"
	"strings"
	"time"
)

// 新增留言資訊
func InsertOrderMessageBoardData(engine *database.MysqlSession, data entity.OrderMessageBoardData) (entity.OrderMessageBoardData, error) {
	data.CreateTime = time.Now()
	data.UpdateTime = time.Now()
	_, err := engine.Session.Table(entity.OrderMessageBoardData{}).Insert(&data)
	if err != nil {
		log.Error("InsertOrderMessageBoardData Error", err)
		return data, err
	}
	return data, nil
}

//更新回覆
func UpdateOrderMessageReply(engine *database.MysqlSession, OrderId string) error {
	sql := fmt.Sprintf("UPDATE order_message_board_data SET reply = ?, update_time = ? WHERE order_id = ?")
	_, err := engine.Session.Exec(sql, 1, time.Now(), OrderId)
	if err != nil {
		log.Error("UpdateOrderMessageBoardData Error", err)
		return err
	}
	return nil
}

// 取得留言資訊
func GetOrderMessageBoardData(engine *database.MysqlSession, orderId string) ([]entity.OrderMessageBoardData, error) {
	var data []entity.OrderMessageBoardData
	err := engine.Engine.Table(entity.OrderMessageBoardData{}).Select("*").
		Where("order_id = ? ", orderId).Asc("create_time").Find(&data)
	if err != nil {
		log.Error("GetOrderMessageBoardData Error", err)
		return data, err
	}
	return data, nil
}

func CountBuyerOrderMessageBoardNotReply(engine *database.MysqlSession, UserId string) (int64, error) {
	sql := fmt.Sprintf("SELECT count(*) FROM order_message_board_data WHERE buyer_id = ? AND message_role = ? AND reply = ?")
	result, err := engine.Engine.SQL(sql, UserId, Enum.MemberSeller, 0).Count()
	if err != nil {
		log.Error("CountOrderMessageBoardData Error", err)
		return 0, err
	}
	return result, nil
}


// 計算收銀機未回覆留言數
func CountSellerOrderMessageBoardNotReply(engine *database.MysqlSession, storeId string) (int64, error) {
	sql := fmt.Sprintf("SELECT count(*) FROM order_message_board_data WHERE store_id = ? AND message_role = ? AND reply = ?")
	result, err := engine.Engine.SQL(sql, storeId, Enum.MemberBuyer, 0).Count()
	if err != nil {
		log.Error("CountOrderMessageBoardData Error", err)
		return 0, err
	}
	return result, nil
}

func GetOrderMessageBoardList(engine *database.MysqlSession, where []string, bind []interface{}, orderBy string, limit int, start int) ([]entity.OrderMessageBoardByOrderData, error) {
	var data []entity.OrderMessageBoardByOrderData
	start = (start - 1) * limit
	sql := fmt.Sprintf("SELECT * FROM order_message_board_data m LEFT JOIN order_data o ON m.order_id = o.order_id  " +
		"WHERE %s ORDER BY %s LIMIT %v OFFSET %v", strings.Join(where, " AND "), orderBy, limit, start)
	err := engine.Engine.SQL(sql, bind...).Find(&data)
	if err != nil {
		log.Error("CountOrderMessageBoardData Error", err)
		return data, err
	}
	return data, nil
}

func CountOrderMessageBoard(engine *database.MysqlSession, where []string, bind []interface{}) (int64, error) {
	sql := fmt.Sprintf("SELECT count(*) FROM order_message_board_data m LEFT JOIN order_data o ON m.order_id = o.order_id WHERE %s", strings.Join(where, " AND "))
	result, err := engine.Engine.SQL(sql, bind...).Count()
	if err != nil {
		log.Error("Count Order Message Board Data Error", err)
		return 0, err
	}
	return result, nil
}

func CountBuyerOrderMessageBoardReply(engine *database.MysqlSession, UserId string) (int64, error) {
	sql := fmt.Sprintf("SELECT COUNT(*) FROM (" +
		"SELECT max(id), order_id FROM order_message_board_data WHERE buyer_id = ? AND message_role = ? GROUP BY order_id) a")
	result, err := engine.Engine.SQL(sql, UserId, Enum.MemberBuyer).Count()
	if err != nil {
		log.Error("CountOrderMessageBoardReply Error", err)
		return 0, err
	}
	return result, nil
}


func CountSellerOrderMessageBoardReply(engine *database.MysqlSession, storeId string) (int64, error) {
	sql := fmt.Sprintf("SELECT COUNT(*) FROM (" +
		"SELECT max(id), order_id FROM sharelug.order_message_board_data WHERE store_id = ? AND message_role = ? GROUP BY order_id) a")
	result, err := engine.Engine.SQL(sql, storeId, Enum.MemberSeller).Count()
	if err != nil {
		log.Error("CountOrderMessageBoardReply Error", err)
		return 0, err
	}
	return result, nil
}
