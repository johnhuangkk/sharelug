package model

import (
	"api/services/Enum"
	"api/services/Service/Excel"
	"api/services/VO/Request"
	"api/services/VO/Response"
	"api/services/dao/Balance"
	"api/services/dao/Store"
	"api/services/database"
	"api/services/entity"
	"fmt"
	"time"
)

//帳戶明細
func GetBalanceList(userData entity.MemberData, params Request.BalanceRequest) (Response.BalanceResponse, error) {
	engine := database.GetMysqlEngine()
	defer engine.Close()
	start, end := setBalanceParams(params)
	var resp Response.BalanceResponse
	count, err := Balance.CountBalanceListByUserId(engine, userData.Uid, start, end)
	if err != nil {
		return resp, fmt.Errorf("系統錯誤")
	}
	data, err := Balance.GetBalanceListByUserId(engine, userData.Uid, start, end, params.Limit, params.Page)
	if err != nil {
		return resp, fmt.Errorf("系統錯誤")
	}
	for _, v := range data {
		var res Response.BalanceAccountList
		res.Date = v.CreateTime.Format("2006/01/02")
		res.TransText = Enum.BalanceTrans[v.TransType]
		res.In = int64(v.In)
		res.Out = int64(v.Out)
		res.Balance = int64(v.Balance)
		res.Comment = v.Comment
		resp.BalanceAccountList = append(resp.BalanceAccountList, res)
	}
	resp.BalanceAccountCount = count
	resp.StartTime = start
	t, err := time.Parse("2006-01-02", end)
	if err != nil {
		return resp, fmt.Errorf("系統錯誤")
	}
	resp.EndTime = t.Add(-(time.Hour * time.Duration(24))).Format("2006-01-02")
	return resp, nil
}

func HandleBalanceExport(userData entity.MemberData, params Request.BalanceRequest) (string, error) {
	engine := database.GetMysqlEngine()
	defer engine.Close()
	start, end := setBalanceParams(params)
	data, err := Balance.GetBalancesByDateAndUserId(engine, userData.Uid, start, end)
	if err != nil {
		return "", fmt.Errorf("1001001")
	}
	var report []Response.BalanceAccountList
	for _, v := range data {
		var res Response.BalanceAccountList
		res.Date = v.CreateTime.Format("2006/01/02")
		res.TransText = Enum.BalanceTrans[v.TransType]
		res.In = int64(v.In)
		res.Out = int64(v.Out)
		res.Balance = int64(v.Balance)
		res.Comment = v.Comment
		report = append(report, res)
	}
	filename, err := Excel.BalanceNew().ToBalanceReportFile(report)
	if err != nil {
		return "", fmt.Errorf("1001001")
	}
	return filename, nil
}

func setBalanceParams(params Request.BalanceRequest) (string, string) {
	var start string
	var end string
	switch params.Tab {
		case "SevenDay":
			start = time.Now().Add(-(time.Hour * time.Duration(24) * 7)).Format("2006-01-02")
			end  = time.Now().Add(time.Hour * time.Duration(24)).Format("2006-01-02")
		case "TenDay":
			start = time.Now().Add(-(time.Hour * time.Duration(24) * 10)).Format("2006-01-02")
			end  = time.Now().Add(time.Hour * time.Duration(24)).Format("2006-01-02")
		case "TwentyDay":
			start = time.Now().Add(-(time.Hour * time.Duration(24) * 20)).Format("2006-01-02")
			end  = time.Now().Add(time.Hour * time.Duration(24)).Format("2006-01-02")
		case "Custom":
			StartTime, _ := time.Parse("2006-01-02", params.StartTime)
			start = StartTime.Format("2006-01-02")
			EndTime, _ := time.Parse("2006-01-02", params.EndTime)
			end  = EndTime.Add(time.Hour * time.Duration(24)).Format("2006-01-02")
		default:
			start = time.Now().Add(-(time.Hour * time.Duration(24) * 7)).Format("2006-01-02")
			end  = time.Now().Add(time.Hour * time.Duration(24)).Format("2006-01-02")
	}
	return start, end
}

//保留款項列表
func GetRetainList(userData entity.MemberData, params Request.RetainRequest) (Response.RetainAccountResponse, error) {
	engine := database.GetMysqlEngine()
	defer engine.Close()
	var resp Response.RetainAccountResponse
	count, err := Balance.CountBalanceRetainListByUserId(engine, userData.Uid)
	if err != nil {
		return resp, fmt.Errorf("系統錯誤")
	}
	data, err := Balance.GetBalanceRetainListByUserId(engine, userData.Uid, params.Limit, params.Page)
	if err != nil {
		return resp, fmt.Errorf("系統錯誤")
	}
	for _, v := range data {
		storeData, _ := Store.GetStoreDataByStoreId(engine, v.Order.StoreId)
		var res Response.RetainAccountList
		res.Date = v.Retain.CreateTime.Format("2006/01/02")
		res.StoreName = storeData.StoreName
		res.OrderId = v.Retain.DataId
		res.Amount = int64(v.Retain.In)
		resp.RetainAccountList = append(resp.RetainAccountList, res)
	}
	resp.RetainAccountCount = count
	return resp, nil
}


func GetMyAccount(params Request.MyAccountRequest, userData entity.MemberData) (Response.MyAccountResponse, error) {
	engine := database.GetMysqlEngine()
	defer engine.Close()
	var resp Response.MyAccountResponse
	balance, err := Balance.GetBalanceAccountLastByUserId(engine, userData.Uid)
	if err != nil {
		return resp, fmt.Errorf("系統錯誤")
	}
	resp.Balance = int64(balance.Balance)
	return resp, nil
}
