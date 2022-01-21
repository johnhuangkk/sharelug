package Balance

import (
	"api/services/Enum"
	"api/services/dao/Balance"
	"api/services/database"
	"api/services/entity"
	"api/services/util/log"
	"fmt"
	"time"
)

//取餘額
func GetBalanceByUid(engine *database.MysqlSession, UserId string) float64 {
	data, err := Balance.GetBalanceAccountLastByUserId(engine, UserId)
	if err != nil {
		log.Error("Get Balance Account Database Error")
		return 0
	}
	return data.Balance
}
//保留餘額
func GetBalanceRetainsByUid(engine *database.MysqlSession, UserId string) float64 {
	data, err := Balance.GetBalanceRetainAccountLastByUserId(engine, UserId)
	if err != nil {
		log.Error("Get Balance Account Database Error")
		return 0
	}
	return data.Balance
}

//扣除餘額
func WithdrawDeduction(data entity.WithdrawData) error {
	comment := fmt.Sprintf("提領日期：%s<br>%s/****-%v",
		data.WithdrawTime.Format("2006/01/02"), data.BankName, data.BankAccount[len(data.BankAccount)-4:])
	err := Withdrawal(data.UserId, data.WithdrawId, float64(data.WithdrawAmt), Enum.BalanceTypeWithdraw, comment)
	if err != nil {
		return err
	}
	if data.WithdrawFee != 0 {
		Comment := fmt.Sprintf("提領日期：%s", time.Now().Format("2006/01/02"))
		err = Withdrawal(data.UserId, data.WithdrawId, float64(data.WithdrawFee), Enum.BalanceTypeBankFee, Comment)
		if err != nil {
			return err
		}
	}
	return nil
}

//入帳
func Deposit(UserId, dataId string, amount float64, transType, comment string) error {
	engine := database.GetMysqlEngine()
	defer engine.Close()
	data, err := Balance.GetBalanceAccountLastByUserId(engine, UserId)
	if err != nil {
		log.Error("Get Balance Account Database Error")
		return err
	}
	var ent entity.BalanceAccountData
	ent.UserId = UserId
	ent.DataId = dataId
	ent.TransType = transType
	ent.In = amount
	ent.Out = 0
	ent.Balance = data.Balance + amount
	ent.Comment = comment
	ent.CreateTime = time.Now()
	err = Balance.InsertBalanceAccountData(engine, ent)
	if err != nil {
		log.Error("Insert Balance Account Database Error")
		return err
	}
	return nil
}

//出帳
func Withdrawal(UserId string, dataId string, amount float64, transType string, comment string) error {
	engine := database.GetMysqlEngine()
	defer engine.Close()
	data, err := Balance.GetBalanceAccountLastByUserId(engine, UserId)
	if err != nil {
		log.Error("Get Balance Account Database Error")
		return err
	}
	var ent entity.BalanceAccountData
	ent.UserId = UserId
	ent.DataId = dataId
	ent.TransType = transType
	ent.In = 0
	ent.Out = amount
	ent.Balance = data.Balance - amount
	ent.Comment = comment
	ent.CreateTime = time.Now()
	err = Balance.InsertBalanceAccountData(engine, ent)
	if err != nil {
		log.Error("Insert Balance Account Database Error")
		return err
	}
	return nil
}

//保留餘額存入
func RetainDeposit(SellerId, dataId string, amount float64, transType, comment string) error {
	engine := database.GetMysqlEngine()
	defer engine.Close()
	data, err := Balance.GetBalanceRetainAccountLastByUserId(engine, SellerId)
	if err != nil {
		log.Error("Get Balance Account Database Error")
		return err
	}
	var ent entity.BalanceRetainAccountData
	ent.UserId = SellerId
	ent.DataId = dataId
	ent.TransType = transType
	ent.In = amount
	ent.Out = 0
	ent.Balance = data.Balance + amount
	ent.Comment = comment
	ent.CreateTime = time.Now()
	if err := Balance.InsertBalanceRetainAccountData(engine, ent); err != nil {
		log.Error("Insert Balance Account Database Error")
		return err
	}
	return nil
}

//保留餘額支出
func RetainWithdrawal(SellerId, dataId string, amount float64, transType, comment string) error {
	engine := database.GetMysqlEngine()
	defer engine.Close()
	data, err := Balance.GetBalanceRetainAccountLastByUserId(engine, SellerId)
	if err != nil {
		log.Error("Get Balance Account Database Error")
		return err
	}
	var ent entity.BalanceRetainAccountData
	ent.UserId = SellerId
	ent.DataId = dataId
	ent.TransType = transType
	ent.In = 0
	ent.Out = amount
	ent.Balance = data.Balance - amount
	ent.Comment = comment
	ent.CreateTime = time.Now()
	if err := Balance.InsertBalanceRetainAccountData(engine, ent); err != nil {
		log.Error("Insert Balance Account Database Error")
		return err
	}
	return nil
}

func CountBalanceByUserIdAndOrderId(engine *database.MysqlSession, UserId, OrderId string) (int64, error) {
	count, err := Balance.CountBalanceByUserIdAndOrderId(engine, UserId, OrderId)
	if err != nil {
		return count, err
	}
	return count, nil
}

func RecalculateBalance(uid string) error {
	engine := database.GetMysqlEngine()
	defer engine.Close()
	data, err := Balance.GetBalancesByUserId(engine, uid)
	if err != nil {
		log.Error("Get Balances Error", err)
		return err
	}
	for k, v := range data {
		if k != 0 {
			balance, _ := Balance.GetBalanceById(engine, data[k-1].Id)
			log.Debug("Balance 1", balance)
			log.Debug("Balance 2", int64(balance.Balance))
			v.Balance = float64(int64(balance.Balance) + int64(v.In) - int64(v.Out))
			log.Debug("Balance 3", int64(v.Balance))
			if err := Balance.UpdateBalancesData(engine, v); err != nil {
				log.Debug("Update Balances Error")
				return err
			}
		} else {
			v.Balance = v.In - v.Out
			if err := Balance.UpdateBalancesData(engine, v); err != nil {
				log.Debug("Update Balances Error")
				return err
			}
		}
	}
	return nil
}

func RecalculateBalanceRetain(uid string) error {
	engine := database.GetMysqlEngine()
	defer engine.Close()
	if err := RecalculateRetain(engine, uid); err != nil {
		log.Error("Recalculate Retain Error", err)
		return fmt.Errorf("1001001")
	}
	return nil
}

func RecalculateRetain(engine *database.MysqlSession, uid string) error {
	data, err := Balance.GetBalanceRetainsByUserId(engine, uid)
	if err != nil {
		log.Error("Get Balances Error", err)
		return err
	}
	for k, v := range data {
		if k != 0 {
			balance, _ := Balance.GetBalanceRetainById(engine, data[k-1].Id)
			log.Debug("Balance 1", balance)
			log.Debug("Balance 2", int64(balance.Balance))
			v.Balance = float64(int64(balance.Balance) + int64(v.In) - int64(v.Out))
			log.Debug("Balance 3", int64(v.Balance))
		} else {
			v.Balance = v.In - v.Out
		}
		if err := Balance.UpdateBalanceRetainData(engine, v); err != nil {
			log.Debug("Update Balances Error")
			return err
		}
	}
	return nil
}