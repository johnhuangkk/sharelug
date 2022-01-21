package MemberService

import (
	"api/config/middleware"
	"api/services/Enum"
	"api/services/Service/SysLog"
	"api/services/dao/Credit"
	"api/services/dao/member"
	"api/services/database"
	"api/services/entity"
	"api/services/util/log"
	"api/services/util/tools"
	"fmt"
	"strings"
)

//取得信用卡資料
func GetCreditData(engine *database.MysqlSession, cardId, number, expire, userId string) (entity.MemberCardData, error) {
	var resp entity.MemberCardData
	count := strings.Count(number, "*")
	if count == 0 {
		//卡號沒有*的
		expireDate := ""
		if len(expire) > 0 {
			exp := strings.Split(expire, "/")
			expireDate = fmt.Sprintf("%s/%s", exp[1], exp[0])
		}
		data, err := member.GetMemberCreditByUserIdAndNumber(engine, number, expireDate, userId)
		if err != nil {
			return resp, err
		}
		if len(data.CardId) == 0 {
			data, err = createCreditData(engine, number, expireDate, userId)
			if err != nil {
				return resp, err
			}
		}
		resp = data
	} else {
		//取信用卡資料
		data, err := member.GetMemberCreditDataByCardIdAndUserId(engine, userId, cardId)
		if err != nil {
			return resp, err
		}
		data.Frequency = resp.Frequency + 1
		if err := member.UpdateMemberCardData(engine, data); err != nil {
			log.Error("Update Member Card Error!!")
			return resp, err
		}
		log.Info("會員讀取信用卡", userId, middleware.GetClientIP(), Enum.SyslogSuccess)
		resp = data
	}
	return resp, nil
}

func createCreditData(engine *database.MysqlSession, cardNumber, expireDate, userId string) (entity.MemberCardData, error)  {
	var data entity.MemberCardData
	CardNumber, card := tools.ParseCredit(cardNumber)
	expiryDate := strings.Replace(expireDate, "/", "", -1)
	bank, err := CompareCreditBank(engine, card[0] + card[1])
	if err != nil {
		return data, err
	}
	data.CardId = tools.GeneratorCardId()
	data.BankName = bank.BankName
	data.IsDebit = bank.IsDebit
	data.CardType = bank.CardType
	data.MemberId = userId
	data.First4Digits = card[0]
	data.Last4Digits = card[3]
	data.CardNumber, _ = tools.AwsKMSEncrypt(CardNumber)
	data.ExpiryDate = expiryDate
	data.DefaultCard = "1"
	if len(bank.BankName) == 0 {
		data.IsForeign = 1
	}
	data.Status = Enum.CreditStatusSuccess
	if err := member.UpdateMemberCardDataDefault(engine, userId); err != nil {
		log.Error("Update Member Card Error!!")
		return data, err
	}
	if _, err := member.InsertMemberCreditData(engine, data); err != nil {
		log.Error("Insert Member Card Error!!")
		return data, err
	}
	if err := SysLog.SetCreditSystemLog(userId, data.BankName, data.Last4Digits); err != nil {
		log.Error("System Log Error", err)
	}
	log.Info("會員新增信用卡", userId, middleware.GetClientIP(), Enum.SyslogSuccess)
	return data, nil
}

//比對信用卡發卡行
func CompareCreditBank(engine *database.MysqlSession, number string) (entity.BankBinCode, error) {
	var resp entity.BankBinCode
	data, err := Credit.GetBankBinCodeByBinCode(engine, number[0:6])
	if err != nil {
		return resp, fmt.Errorf("系統錯誤")
	}
	for _, v := range data {
		if v.BinNumber == number[0:6] {
			resp = v
		}
		if v.BinNumber == number {
			resp = v
		}
	}
	return resp, nil
}

//變更會員信用卡已請款
func ChangeCreditCardIsRelease(engine *database.MysqlSession, cardId string) error {
	data, err := member.GetMemberCreditDataByCardId(engine, cardId)
	if err != nil {
		return err
	}
	data.IsRelease = 1
	if err := member.UpdateMemberCardData(engine, data); err != nil {
		return err
	}
	return nil
}

//修改會員信用卡銀行及卡別資料
func ChangeMemberCardData() error {
	engine := database.GetMysqlEngine()
	defer engine.Close()
	data, err := member.GetMemberCreditData(engine)
	if err != nil {
		log.Error("Get Member Card Error", err)
		return err
	}
	for _, v := range data {
		cardNumber, _ := tools.AwsKMSDecrypt(v.CardNumber)
		if len(cardNumber) != 0 {
			log.Debug("cardNumber", cardNumber[0:9])
			bank, err := CompareCreditBank(engine, cardNumber[0:9])
			if err != nil {
				log.Error("Get Bank bin code Error", err)
				return err
			}
			v.IsDebit = bank.IsDebit
			v.BankName = bank.BankName
			v.CardType = bank.CardType
			if err := member.UpdateMemberCardData(engine, v); err != nil {
				log.Debug("Update Member Card Error", err)
			}
		}
	}
	return nil
}

