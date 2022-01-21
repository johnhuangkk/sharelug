package model

import (
	"api/services/util/tools"
	"encoding/json"
	"fmt"
)


func generateLog(date, time, state string, obj interface{}) string {
	data, _ := json.Marshal(obj)
	return fmt.Sprintf("[%v %v]-[%v]-[%v]\n", tools.Now("Ymd"), tools.NowHHmmss(), state, string(data))
}

func generateRecord(date, time, direction, state string) string {
	return fmt.Sprintf("%v %v,%v,%v\n", tools.Now("Ymd"), tools.NowHHmmss(), direction, state)
}

