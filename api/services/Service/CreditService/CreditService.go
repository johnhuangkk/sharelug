package CreditService

import (
	"api/config/middleware"
	"api/services/Enum"
	"api/services/dao/Credit"
	"api/services/database"
	"api/services/entity"
	"api/services/util/log"
	"github.com/spf13/viper"
)

//建立GW刷卡資料
func NewAuthGwData(engine *database.MysqlSession, params entity.AuthRequest) (entity.GwCreditAuthData, error) {
	var MerchantId string
	var TerminalId string
	var err error
	//是否使用 3D 取受權
	if params.IsFirst {
		log.Debug("Credit Auth 3D")
		//交易模式 C2C 或 B2C
		if params.Type == Enum.OrderTransC2c {
			//是否使用次特店代碼
			if params.IsMerchant {
				MerchantId, TerminalId, err = getSeller3DMerchantId(engine, params.SellerId)
				if err != nil {
					return entity.GwCreditAuthData{}, err
				}
			} else {
				MerchantId = viper.GetString("KgiCredit.C2C.3D.MerchantID")
				TerminalId = viper.GetString("KgiCredit.C2C.3D.TerminalID")
			}
		} else {
			MerchantId = viper.GetString("KgiCredit.B2C.3D.MerchantID")
			TerminalId = viper.GetString("KgiCredit.B2C.3D.TerminalID")
		}
	} else {
		log.Debug("Credit Auth N3D")
		//交易模式 C2C 或 B2C
		if params.Type == Enum.OrderTransC2c {
			//是否使用次特店代碼
			if params.IsMerchant {
				MerchantId, TerminalId, err = getSellerN3DMerchantId(engine, params.SellerId)
				if err != nil {
					return entity.GwCreditAuthData{}, err
				}
			} else {
				MerchantId = viper.GetString("KgiCredit.C2C.N3D.MerchantID")
				TerminalId = viper.GetString("KgiCredit.C2C.N3D.TerminalID")
			}
		} else {
			MerchantId = viper.GetString("KgiCredit.B2C.N3D.MerchantID")
			TerminalId = viper.GetString("KgiCredit.B2C.N3D.TerminalID")
		}
	}
	result := params.GenerateGwCreditAuthData(middleware.ClientIP, MerchantId, TerminalId)
	data, err := Credit.InsertGwCreditData(engine, result)
	if err != nil {
		return data, err
	}
	return data, nil
}
//取得3D特店代碼
func getSeller3DMerchantId(engine *database.MysqlSession, sellerId string) (string, string, error) {
	var merchantId string
	var terminalId string
	//取出賣家的次特店代碼
	data, err := Credit.GetSellerMerchantIdBySellerId(engine, sellerId)
	if err != nil {
		return merchantId, terminalId, err
	}
	if len(data.MerchantId) != 0 {
		//有次特店代碼
		merchantId = data.MerchantId
		terminalId = data.Terminal3dId
	} else {
		//無次特店代碼
		//取出未使用次特店代碼
		merchant, err := Credit.GetSellerMerchantIdUnused(engine)
		if err != nil {
			return merchantId, terminalId, err
		}
		if len(merchant.MerchantId) == 0 {
			//已無未使用次特店代碼 代入平台特店代碼
			merchantId = viper.GetString("KgiCredit.C2C.3D.MerchantID")
			terminalId = viper.GetString("KgiCredit.C2C.3D.TerminalID")
		} else {
			//有未使用次特店代碼及更新次特店代碼
			merchant.UserId = sellerId
			merchant.IsUsed = true
			if err := Credit.UpdateSellerMerchantId(engine, merchant); err != nil {
				return merchantId, terminalId, err
			}
			merchantId = merchant.MerchantId
			terminalId = merchant.Terminal3dId
		}
	}
	return merchantId, terminalId, nil
}
//取得N3D特店代碼
func getSellerN3DMerchantId(engine *database.MysqlSession, sellerId string) (string, string, error) {
	var merchantId string
	var terminalId string

	//取出賣家的次特店代碼
	data, err := Credit.GetSellerMerchantIdBySellerId(engine, sellerId)
	if err != nil {
		return merchantId, terminalId, err
	}
	if len(data.MerchantId) != 0 {
		//有次特店代碼
		merchantId = data.MerchantId
		terminalId = data.Terminaln3dId
	} else {
		//無次特店代碼
		//取出未使用次特店代碼
		merchant, err := Credit.GetSellerMerchantIdUnused(engine)
		if err != nil {
			return merchantId, terminalId, err
		}
		if len(merchant.MerchantId) == 0 {
			//已無未使用次特店代碼 代入平台特店代碼
			merchantId = viper.GetString("KgiCredit.C2C.N3D.MerchantID")
			terminalId = viper.GetString("KgiCredit.C2C.N3D.TerminalID")
		} else {
			//有未使用次特店代碼及更新次特店代碼
			merchant.UserId = sellerId
			merchant.IsUsed = true
			if err := Credit.UpdateSellerMerchantId(engine, merchant); err != nil {
				return merchantId, terminalId, err
			}
			merchantId = merchant.MerchantId
			terminalId = merchant.Terminaln3dId
		}
	}
	return merchantId, terminalId,  nil
}
