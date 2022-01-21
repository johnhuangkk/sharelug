package OrderMessageBoardVo

import (
	"api/services/entity"
)

type OrderMessage struct {
	OrderId   string `form:"OrderId" json:"OrderId" validate:"required"`
	Message   string `form:"Message" json:"Message" validate:"required"`
}

func (od *OrderMessage) GetOrderMessageBoardEnt() entity.OrderMessageBoardData {
	var data = entity.OrderMessageBoardData{}
	data.OrderId = od.OrderId
	data.Message = od.Message

	return data
}
