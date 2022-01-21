package model

import (
	"api/services/Enum"
	"api/services/Service/Mail"
	"api/services/Service/Notification"
	"api/services/VO/Request"
	"api/services/VO/Response"
	"api/services/dao/Credit"
	"api/services/dao/Customer"
	"api/services/dao/Orders"
	"api/services/database"
	"api/services/entity"
	"api/services/util/log"
	"api/services/util/tools"
	"fmt"
	"time"
)

//取出問題選單
func HandleCustomerQuestion() ([]Response.CustomerResponse, error) {
	engine := database.GetMysqlEngine()
	defer engine.Close()
	var resp []Response.CustomerResponse
	data, err := Customer.GetCustomerQuestionData(engine)
	if err != nil {
		log.Error("Get Customer Question Error")
		return nil, err
	}
	for _, v := range data {
		rep := Response.CustomerResponse{
			Question: v.Question,
			Type:     v.QuestionType,
		}
		resp = append(resp, rep)
	}
	return resp, nil
}

func HandleSendCustomer(userData entity.MemberData, params Request.CustomerRequest) error {
	engine := database.GetMysqlEngine()
	defer engine.Close()
	customerId := tools.GenerateCustomerId()
	data := entity.CustomerData{
		QuestionId: customerId,
		UserId:     userData.Uid,
		Question:   params.Question,
		OrderId:    params.OrderId,
		Contents:   params.Contents,
		Status:     `none`,
	}
	if params.RelationId == "" {
		data.RelatedId = data.QuestionId
	} else {
		data.RelatedId = params.RelationId
	}
	err := Customer.InsertCustomerData(engine, data)
	if err != nil {
		log.Error("Insert Customer Error", err)
		return err
	}
	//Todo remove when erp ready
	if err := Notification.SendUserCustomerMessage(engine, userData.Uid, params.OrderId, data.Question, customerId); err != nil {
		log.Error("Send Custmer Request Email Error", err)
	}
	return nil
	// if err := Mail.SendHostCustmerRequestEmail(userData.Email, userData.Mphone, userData.Username, data); err != nil {
	// 	log.Error("Send Custmer Request Email Error", err)
	// }
	// return nil
}

//建立聯絡我們
func HandleContactData(params Request.ContactRequest) error {
	engine := database.GetMysqlEngine()
	defer engine.Close()
	data := entity.ContactData{}
	data.Email = params.Email
	data.UserName = params.UserName
	data.Company = params.Company
	data.Telephone = params.Telephone
	data.Contents = params.Contents
	data.CreateTime = time.Now()
	_ = Mail.SendHostCustmerContactEmail(params.Company, params.Email, params.Telephone, params.UserName, params.Contents)
	err := Customer.InsertContactData(engine, data)
	if err != nil {
		log.Error("Insert Contact Error", err)
		return err
	}
	return nil
}

//訂單客服
func GetErpOrderCustomerData(engine *database.MysqlSession, OrderId string) (Response.ErpOrderCustomerResponse, error) {
	var resp Response.ErpOrderCustomerResponse
	data, err := Customer.FindCustomerByOrderId(engine, OrderId)
	if err != nil {
		log.Error("get Customer Data Error", err)
		return resp, err
	}
	for _, v := range data {
		resp.Contents += fmt.Sprintf("%s \n %s \n %s \n %s \n",
			v.CreateTime.Format("2006/01/02 15:04"), v.Contents,
			v.ReplyTime.Format("2006/01/02 15:04"), v.Reply)
		resp.Remark += fmt.Sprintf("%s by %s \n %s \n",
			v.RemarkTime.Format("2006/01/02 15:04"), v.RemarkStaff, v.Remark)
	}
	msg, err := Orders.GetOrderMessageBoardData(engine, OrderId)
	if err != nil {
		log.Error("get Message Data Error", err)
		return resp, err
	}
	for _, v := range msg {
		resp.Messages += fmt.Sprintf("%s by %s \n %s \n",
			v.CreateTime.Format("2006/01/02 15:04"), Enum.MessageRole[v.MessageRole], v.Message)
	}
	auth, err := Credit.GetGwCreditByOrderId(engine, OrderId)
	if err != nil {
		log.Error("get Message Data Error", err)
		return resp, err
	}
	resp.AuditRemark = auth.Memo

	return resp, nil
}

func GetCustomerQuestionByMsgType(questionId string, msgType string) (Response.CustomerHistoryQuestionResponse, error) {
	engine := database.GetMysqlEngine()
	defer engine.Close()
	var resp Response.CustomerHistoryQuestionResponse
	data, err := Customer.GetCustomerQuestionByQuestionId(engine, questionId)
	if err != nil || (data == entity.CustomerData{}) {
		log.Error("Get Customer Question Error")
		return resp, err
	}
	resp.Content = data.Contents
	resp.CreateTime = data.CreateTime.Format("2006/01/02 15:04")
	resp.Question = data.Question
	resp.QuestionId = data.QuestionId
	resp.RelatedId = data.RelatedId
	resp.OrderId = data.OrderId
	if msgType != Enum.NotifyTypePlatformUser {
		resp.ReplyContent = data.Reply
		resp.ReplyTime = data.ReplyTime.Format("2006/01/02 15:04")
	}
	if data.QuestionId != data.RelatedId {
		var startTime time.Time
		if msgType == Enum.NotifyTypePlatformUser {
			startTime = data.CreateTime
		} else {
			startTime = data.ReplyTime
		}
		var relatedDatas []Response.CustomerQuestionResponse
		relatedQuestions, err := Customer.GetCustomerQuestionByRelatedIdWithTime(engine, data.RelatedId, startTime, msgType)
		if err != nil {
			log.Error("Get Customer Question Error")
			return resp, err
		}
		if len(relatedQuestions) > 0 {
			for _, question := range relatedQuestions {
				data := Response.CustomerQuestionResponse{
					OrderId:      question.OrderId,
					Question:     question.Question,
					QuestionId:   question.QuestionId,
					RelatedId:    question.RelatedId,
					Content:      question.Contents,
					ReplyContent: question.Reply,
					CreateTime:   question.CreateTime.Format("2006/01/02 15:04"),
					ReplyTime:    question.ReplyTime.Format("2006/01/02 15:04"),
				}
				relatedDatas = append(relatedDatas, data)
			}
			resp.RelatedQuestions = relatedDatas
		}

	}
	return resp, nil
}

func GetCustomerByRelatedId(questionId string) (Response.CustomerHistoryQuestionResponse, error) {
	engine := database.GetMysqlEngine()
	defer engine.Close()
	var resp Response.CustomerHistoryQuestionResponse
	data, err := Customer.GetCustomerQuestionByQuestionId(engine, questionId)
	if err != nil || (data == entity.CustomerData{}) {
		log.Error("Get Customer Question Error")
		return resp, err
	}
	resp.Content = data.Contents
	resp.CreateTime = data.CreateTime.Format("2006/01/02 15:04")
	resp.Question = data.Question
	resp.QuestionId = data.QuestionId
	resp.RelatedId = data.RelatedId
	resp.OrderId = data.OrderId
	resp.ReplyContent = data.Reply
	resp.ReplyTime = data.ReplyTime.Format("2006/01/02 15:04")
	if data.QuestionId != data.RelatedId {
		startTime := data.CreateTime
		var relatedDatas []Response.CustomerQuestionResponse
		relatedQuestions, err := Customer.GetCustomerQuestionByRelatedIdWithoutType(engine, data.RelatedId, startTime)
		if err != nil {
			log.Error("Get Customer Question Error")
			return resp, err
		}
		if len(relatedQuestions) > 0 {
			for _, question := range relatedQuestions {
				data := Response.CustomerQuestionResponse{
					OrderId:      question.OrderId,
					Question:     question.Question,
					QuestionId:   question.QuestionId,
					RelatedId:    question.RelatedId,
					Content:      question.Contents,
					ReplyContent: question.Reply,
					CreateTime:   question.CreateTime.Format("2006/01/02 15:04"),
					ReplyTime:    question.ReplyTime.Format("2006/01/02 15:04"),
				}
				relatedDatas = append(relatedDatas, data)
			}
			resp.RelatedQuestions = relatedDatas
		}

	}
	return resp, nil
}
