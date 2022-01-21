package OrderMessageBoard

import (
	"api/services/Enum"
	"api/services/Service/Notification"
	"api/services/VO/OrderMessageBoardVo"
	"api/services/VO/Request"
	"api/services/VO/Response"
	"api/services/dao/Orders"
	"api/services/dao/Store"
	"api/services/dao/member"
	"api/services/database"
	"api/services/entity"
	"api/services/util/log"
	"api/services/util/tools"
	"fmt"
)

// 新增留言
func AddMessage(params OrderMessageBoardVo.OrderMessage, UserData entity.MemberData, store entity.StoreDataResp) error {
	engine := database.GetMysqlEngine()
	defer engine.Close()

	orderData, _ := Orders.GetOrderByOrderId(engine, params.OrderId)

	var orderMessageData = entity.OrderMessageBoardData{
		OrderId: params.OrderId,
		Message: params.Message,
	}
	log.Info("member", UserData.Uid, orderData.BuyerId)
	log.Info("store", orderData.StoreId, store.StoreId, store)
	// 判斷新增訂單留言板使用者是否為買方 或是 店家
	if orderData.BuyerId != UserData.Uid && orderData.StoreId != store.StoreId {
		log.Error("AddMessage [%s] %s", UserData.Uid, orderData)
		return fmt.Errorf("系統錯誤")
	}
	// 買方新增訊息
	if orderData.BuyerId == UserData.Uid {
		orderMessageData.MessageRole = Enum.MemberBuyer
		orderMessageData.GetBuyerMessageData(UserData)
		StoreData, err := Store.GetStoreDataByStoreId(engine, orderData.StoreId)
		if err != nil {
			return fmt.Errorf("系統錯誤")
		}
		orderMessageData.StoreId = StoreData.StoreId
		orderMessageData.StoreName = StoreData.StoreName
		if _, err := Orders.InsertOrderMessageBoardData(engine, orderMessageData); err != nil {
			log.Error("AddOrderMessageBoard InsertOrderMessageBoardData Error : %s", err)
			return fmt.Errorf("AddOrderMessageBoard InsertOrderMessageBoardData Error : %s", err)
		}
		if err := Notification.SendOrderCustomerMessage(engine, orderData.OrderId, 1, orderMessageData); err != nil {
			return fmt.Errorf("系統錯誤")
		}
	} else {
		// 店家新增訊息
		orderMessageData.MessageRole = Enum.MemberSeller
		orderMessageData.GetStoreMessageData(store)
		buyerData, err := member.GetMemberDataByUid(engine, orderData.BuyerId)
		if err != nil {
			return err
		}
		orderMessageData.GetBuyerMessageData(buyerData)
		// 回寫回覆
		if err := Orders.UpdateOrderMessageReply(engine, orderData.OrderId); err != nil {
			log.Error("UpdateOrderMessageBoard UpdateOrderMessageBoardData Error : %s", err)
			return fmt.Errorf("系統錯誤")
		}
		if _, err := Orders.InsertOrderMessageBoardData(engine, orderMessageData); err != nil {
			log.Error("AddOrderMessageBoard InsertOrderMessageBoardData Error : %s", err)
			return fmt.Errorf("AddOrderMessageBoard InsertOrderMessageBoardData Error : %s", err)
		}
		if err := Notification.SendReplyOrderCustomerMessage(engine, orderData); err != nil {
			return fmt.Errorf("系統錯誤")
		}
	}
	return nil
}

// 取得買方和賣家店舖資訊
func GetBuyerAndStoreData(orderId string, userData entity.MemberData, storeData entity.StoreDataResp) (OrderMessageBoardVo.BuyerStorePictureData, error) {
	var p = OrderMessageBoardVo.BuyerStorePictureData{}
	engine := database.GetMysqlEngine()
	defer engine.Close()
	orderData, _ := Orders.GetOrderByOrderId(engine, orderId)
	if userData.Uid != orderData.BuyerId && storeData.StoreId != orderData.StoreId {
		log.Error("此訂單非你所屬 [Uid: %s, Sid: %s]", userData.Uid, storeData.StoreId)
		return p, fmt.Errorf("此訂單非你所屬")
	}
	store, _ := Store.GetStoreDataByStoreId(engine, orderData.StoreId)
	buyer, _ := member.GetMemberDataByUid(engine, orderData.BuyerId)
	p.OrderId = orderId
	p.IsBuyer = orderData.BuyerId == userData.Uid
	p.Buyer.Id = buyer.Uid
	p.Buyer.Name = buyer.Username
	p.Buyer.Picture = buyer.Picture
	p.Store.Id = store.StoreId
	p.Store.Name = store.StoreName
	p.Store.Picture = store.StorePicture
	p.OrderInfo = takeOrderInfo(engine, orderData)
	return p, nil
}

func takeOrderInfo(engine *database.MysqlSession, order entity.OrderData) OrderMessageBoardVo.OrderInfo {
	var resp OrderMessageBoardVo.OrderInfo
	data, err := Orders.GetOrderDetailListByOrderId(engine, order.OrderId)
	if err != nil {
		log.Debug("Get Order Detail Error", err)
	}
	for _, v := range data {
		resp.Products = append(resp.Products, v.ProductName)
	}
	resp.BuyerName = tools.MaskerName(order.BuyerName)
	resp.BuyerPhone = tools.MaskerPhoneLater(order.BuyerPhone)
	return resp
}

// 取得留言
func GetMessage(UserData entity.MemberData, storeData entity.StoreDataResp, orderId string) ([]entity.OrderMessageBoardData, error) {
	engine := database.GetMysqlEngine()
	defer engine.Close()
	var messageData []entity.OrderMessageBoardData
	if len(orderId) == 0 {
		log.Debug("GetOrderMessageBoardAction orderId : %s", orderId)
		return nil, fmt.Errorf("1001002")
	}
	if len(UserData.Uid) == 0 {
		return nil, fmt.Errorf("1001000")
	}
	orderData, err := Orders.GetOrderByOrderId(engine, orderId)
	if err != nil {
		return nil, fmt.Errorf("1001001")
	}
	if orderData.BuyerId == UserData.Uid || orderData.SellerId == storeData.SellerId {
		data, err := Orders.GetOrderMessageBoardData(engine, orderId)
		if err != nil {
			return nil, fmt.Errorf("1001001")
		}
		messageData = data
	}
	return messageData, nil
}

//取得留言列表
func GetMessageBoardList(userData entity.MemberData, storeData entity.StoreDataResp, params Request.OrderMessageBoardRequest, Type string) (Response.OrderMessageBoardResponse, error) {
	engine := database.GetMysqlEngine()
	defer engine.Close()
	var resp Response.OrderMessageBoardResponse
	per := tools.CheckIsZero(params.Limit, 30)
	page := tools.CheckIsZero(params.Start, 1)
	where, bind, orderBy := setMessageParams(params, Type, storeData, userData)
	result, err := Orders.GetOrderMessageBoardList(engine, where, bind, orderBy, per, page)
	if err != nil {
		return resp, fmt.Errorf("系統錯誤")
	}
	count, err := Orders.CountOrderMessageBoard(engine, where, bind)
	if err != nil {
		return resp, fmt.Errorf("系統錯誤")
	}
	for _, v := range result {
		var res Response.MessageList
		res.OrderId = v.Order.OrderId
		res.OrderTime = v.Order.CreateTime.Format("2006/01/02 15:04")
		res.MessageTime = v.MessageBoard.CreateTime.Format("2006/01/02 15:04")
		res.MessageContent = v.MessageBoard.Message
		res.MessageReply = "false"
		if v.MessageBoard.Reply == 1 {
			res.MessageReply = "true"
		}
		resp.MessageList = append(resp.MessageList, res)
	}
	resp.BoardCount = count
	Inbox, err := Orders.CountSellerOrderMessageBoardNotReply(engine, storeData.StoreId)
	if Type == Enum.MemberBuyer {
		Inbox, err = Orders.CountBuyerOrderMessageBoardNotReply(engine, userData.Uid)
	}
	if err != nil {
		return resp, fmt.Errorf("系統錯誤")
	}
	resp.Tabs.Inbox = Inbox
	reply, err := Orders.CountSellerOrderMessageBoardReply(engine, storeData.StoreId)
	if Type == Enum.MemberBuyer {
		reply, err = Orders.CountBuyerOrderMessageBoardReply(engine, userData.Uid)
	}
	if err != nil {
		return resp, fmt.Errorf("系統錯誤")
	}
	resp.Tabs.Replied = reply
	return resp, nil
}

func setMessageParams(params Request.OrderMessageBoardRequest, Type string, storeData entity.StoreDataResp, userData entity.MemberData) ([]string, []interface{}, string) {
	var sql 	[]string
	var bind 	[]interface{}
	if Type == Enum.MemberBuyer {
		sql = append(sql, "m.buyer_id = ?")
		bind = append(bind, userData.Uid)
	} else {
		sql = append(sql, "m.store_id = ?")
		bind = append(bind, storeData.StoreId)
	}
	switch params.Tabs {
		case "Replied":
			sql = append(sql, "m.message_role = ?")
			if Type == Enum.MemberBuyer {
				bind = append(bind, Enum.MemberSeller)
			} else {
				bind = append(bind, Enum.MemberBuyer)
			}
			sql = append(sql, "m.reply = ?")
			bind = append(bind, 1)
	}

	var orderBy string
	switch params.OrderBy {
	case "Message":
		orderBy = fmt.Sprintf("m.reply ASC, m.create_time DESC")
	case "Order":
		orderBy = fmt.Sprintf("m.reply ASC, o.create_time DESC")
	default:
		orderBy = fmt.Sprintf("m.reply ASC, m.create_time DESC")
	}
	return sql, bind, orderBy
}