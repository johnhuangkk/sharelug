package model

import (
	"api/services/Enum"
	"api/services/VO/Request"
	"api/services/VO/Response"
	"api/services/dao/Balance"
	"api/services/dao/member"
	"api/services/database"
	"fmt"
	"time"
)

func HandleSearchBalance(params Request.SearchBalanceRequest) (Response.BalanceResponse, error) {
	engine := database.GetMysqlEngine()
	defer engine.Close()
	var resp Response.BalanceResponse
	//取出會員資料
	userData, err := member.GetMemberDataByUid(engine, params.Account)
	if err != nil {
		return resp, fmt.Errorf("1001001")
	}
	var start, end string
	if len(params.StartTime) == 0 || len(params.EndTime) == 0 {
		start = time.Now().AddDate(0, 0, -30).Format("2006-01-02")
		end = time.Now().Add(time.Hour * time.Duration(24)).Format("2006-01-02")
	} else {
		StartTime, _ := time.Parse("2006-01-02", params.StartTime)
		start = StartTime.Format("2006-01-02")
		EndTime, _ := time.Parse("2006-01-02", params.EndTime)
		end = EndTime.Add(time.Hour * time.Duration(24)).Format("2006-01-02")
	}
	count, err := Balance.CountBalanceListByUserId(engine, userData.Uid, start, end)
	if err != nil {
		return resp, fmt.Errorf("1001001")
	}
	data, err := Balance.GetBalanceListNoLimitByUserId(engine, userData.Uid, start, end)
	if err != nil {
		return resp, fmt.Errorf("1001001")
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
		return resp, fmt.Errorf("1001001")
	}
	resp.EndTime = t.Add(-(time.Hour * time.Duration(24))).Format("2006-01-02")
	return resp, nil
}
