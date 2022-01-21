package Customer

import (
	"api/services/dao/Customer"
	"api/services/database"
	"api/services/entity"
	"api/services/util/log"
)


func HandleCustomerQuestion() error {
	engine := database.GetMysqlEngine()
	defer engine.Close()

	data := getCustomerQuestion()
	for _, v := range data {
		count, _ := Customer.CountCustomerQuestionData(engine)
		v.Sort = count
		err := Customer.InsertCustomerQuestionData(engine, v)
		if err != nil {
			log.Error("Insert Customer Error!!")
			return err
		}
	}
	return nil
}

func getCustomerQuestion() []entity.CustomerQuestionData {
	var resp []entity.CustomerQuestionData

	res := entity.CustomerQuestionData{
		Question: "訂單相關問題",
		QuestionType: 1,
	}
	resp = append(resp, res)

	res = entity.CustomerQuestionData{
		Question: "帳號相關問題",
		QuestionType: 0,
	}
	resp = append(resp, res)

	res = entity.CustomerQuestionData{
		Question: "餘額相關問題",
		QuestionType: 0,
	}
	resp = append(resp, res)

	res = entity.CustomerQuestionData{
		Question: "賣場相關問題",
		QuestionType: 0,
	}
	resp = append(resp, res)

	res = entity.CustomerQuestionData{
		Question: "結帳相關問題",
		QuestionType: 0,
	}
	resp = append(resp, res)
	res = entity.CustomerQuestionData{
		Question: "其他問題",
		QuestionType: 0,
	}
	resp = append(resp, res)
	return resp
}



