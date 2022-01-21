package model

import (
	"api/services/Service/MemberService"
	"api/services/Service/TwIdVerify"
	"api/services/VO/Request"
	"api/services/VO/Response"
	"api/services/dao/TwId"
	"api/services/database"
	"api/services/entity"
	"api/services/util/log"
	"fmt"
	"github.com/spf13/viper"
)

//執行身份證驗證
func HandleTwIdVerify(params *Request.TWIDParams, userId string) (Response.TWIDResponse, error) {
	engine := database.GetMysqlEngine()
	defer engine.Close()

	var response Response.TWIDResponse

	year := params.IssueDate.IssueYear
	month := fmt.Sprintf("%02s", params.IssueDate.IssueMonth)
	day := fmt.Sprintf("%02s", params.IssueDate.IssueDay)
	applyTime := fmt.Sprintf("%s%s%s", year, month, day)

	//縣市轉換
	//if applyTime >= "0991225" {
	//	if params.IssueCounties == "10001" {
	//		params.IssueCounties = "65000"
	//	} else if params.IssueCounties == "10021" {
	//		params.IssueCounties = "67000"
	//	} else if params.IssueCounties == "10019" {
	//		params.IssueCounties = "66000"
	//	} else if params.IssueCounties == "10012" {
	//		params.IssueCounties = "64000"
	//	}
	//}
	//if applyTime >= "1031225" {
	//	if params.IssueCounties == "10003" {
	//		params.IssueCounties = "68000"
	//	}
	//}

	//查詢入容 寫入資料庫
	var Entity entity.TwIdVerifyData
	Entity.UserId = userId
	Entity.IdentityId = params.IdentityId
	Entity.IdentityName = params.IdentityName
	Entity.IdentityType = params.IssueType
	Entity.IdentityDate = applyTime
	Entity.IdentityCounties = params.IssueCounties
	data, err := TwId.InsertTwIdLogData(engine, Entity)
	if err != nil {
		log.Error("insert TwId Log error", err)
		return response, fmt.Errorf("1001001")
	}
	resp, err := TwIdVerify.SendIdentityVerify(params, userId, applyTime)
	if err != nil {
		log.Error("curl post xml error", err)
		return response, fmt.Errorf("1001001")
	}
	//結果 寫入資料庫
	if err := TwId.UpdateTwIdLogData(engine, data.Id, data, resp); err != nil {
		log.Error("Update TwId data error", err)
		return response, fmt.Errorf("1001001")
	}
	if resp.Response.CheckIDCardApply == "1" {
		if err := MemberService.AlterMemberIdentity(engine, userId, params.IdentityId, params.IdentityName); err != nil {
			return response, err
		}
		response.Message = "國民身分證資料與檔存資料相符"
		response.Status = "Success"
		return response, nil
	} else if resp.Response.CheckIDCardApply == "2" ||
		resp.Response.CheckIDCardApply == "3" ||
		resp.Response.CheckIDCardApply == "4" ||
		resp.Response.CheckIDCardApply == "6" ||
		resp.Response.CheckIDCardApply == "7" ||
		resp.Response.CheckIDCardApply == "8" {
		return response, fmt.Errorf("1003005")
	} else if resp.Response.CheckIDCardApply == "5" {
		return response, fmt.Errorf("1003006")
	}
	return response, fmt.Errorf("1003004")
}


//執行身份證驗證
func HandleErpTwIdVerify(params *Request.TWIDParams) (Response.TWIDResponse, error) {
	var response Response.TWIDResponse
	engine := database.GetMysqlEngine()
	defer engine.Close()
	userId := viper.GetString("PLATFORM.USERID")
	year := params.IssueDate.IssueYear
	month := fmt.Sprintf("%02s", params.IssueDate.IssueMonth)
	day := fmt.Sprintf("%02s", params.IssueDate.IssueDay)
	applyTime := fmt.Sprintf("%s%s%s", year, month, day)

	//查詢入容 寫入資料庫
	var Entity entity.TwIdVerifyData
	Entity.UserId = userId
	Entity.IdentityId = params.IdentityId
	Entity.IdentityName = params.IdentityName
	Entity.IdentityType = params.IssueType
	Entity.IdentityDate = applyTime
	Entity.IdentityCounties = params.IssueCounties

	data, err := TwId.InsertTwIdLogData(engine, Entity)
	if err != nil {
		log.Error("insert TwId Log error", err)
		return response, fmt.Errorf("1001001")
	}

	resp, err := TwIdVerify.SendIdentityVerify(params, userId, applyTime)
	if err != nil {
		log.Error("curl post xml error", err)
		return response, fmt.Errorf("1001001")
	}
	//結果 寫入資料庫
	if err := TwId.UpdateTwIdLogData(engine, data.Id, data, resp); err != nil {
		log.Error("Update TwId data error", err)
		return response, fmt.Errorf("1001001")
	}

	if resp.Response.CheckIDCardApply == "1" {
		response.Message = "國民身分證資料與檔存資料相符"
		response.Status = "Success"
		return response, nil
	} else if resp.Response.CheckIDCardApply == "2" || resp.Response.CheckIDCardApply == "3" || resp.Response.CheckIDCardApply == "4" {
		return response, fmt.Errorf("1003005")
	} else if resp.Response.CheckIDCardApply == "6" || resp.Response.CheckIDCardApply == "7" || resp.Response.CheckIDCardApply == "8" {
		return response, fmt.Errorf("1003007")
	} else if resp.Response.CheckIDCardApply == "5" {
		return response, fmt.Errorf("1003006")
	}
	return response, fmt.Errorf("1003004")
}
