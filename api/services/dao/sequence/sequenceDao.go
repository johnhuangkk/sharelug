package sequence

import (
	"api/services/database"
	"api/services/util/log"
)

// 悲觀鎖定 SequenceTable current_value
func pessimisticLockingSequenceTable(field string) (string, error) {
	engine := database.GetMysqlEngine()
	defer engine.Close()
	sql := "SELECT current_value FROM sequence WHERE name = ? FOR update"
	result, err := engine.Session.Query(sql, field)
	if err != nil {
		log.Error("Sequence Database Error", err)
		return "", err
	}
	var data string
	for _, record := range result {
		for _, v := range record {
			data = string(v)
		}
	}
	sql = "UPDATE sequence SET current_value = current_value + 1 WHERE name = ?"
	if _, err = engine.Session.Exec(sql, field); err != nil {
		log.Error("Sequence Database Error", err)
		return "", err
	}
	return data, nil
}

// 凱基虛擬帳號
func GetKgiVirtualSeq() (string, error) {
	currentValue, err := pessimisticLockingSequenceTable("KgiAtmBank")
	return currentValue, err
}

// ？？？？
func GetTinyUrlVirtualSeq() (string, error) {
	currentValue, err := pessimisticLockingSequenceTable("qrcode")
	return currentValue, err
}

// 取得郵局共用郵件編號
func GetIPostMailNoSeq() (string, error) {
	currentValue, err := pessimisticLockingSequenceTable("iPost")
	return currentValue, err
}

// 取得訂單共用編號
func GetOrderIdSeq() (string, error) {
	currentValue, err := pessimisticLockingSequenceTable("order")
	return currentValue, err
}

// 取得郵局便利包共用流水號
func GetPostBagSeq() (string, error) {
	currentValue, err := pessimisticLockingSequenceTable("postBag")
	return currentValue, err
}

func GetCustomerQuestionSeq() (string, error) {
	currentValue, err := pessimisticLockingSequenceTable("customerQuestion")
	return currentValue, err
}

func GetTaiwanAreaSeq() (string, error) {
	currentValue, err := pessimisticLockingSequenceTable("taiwanArea")
	return currentValue, err
}
func GetCouponActionSeq() (string, error) {
	currentValue, err := pessimisticLockingSequenceTable("couponAction")
	return currentValue, err
}
