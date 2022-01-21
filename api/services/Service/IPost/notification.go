package IPost

import (
	"api/services/VO/IPOSTVO"
	"api/services/model"
	"api/services/util/log"
	"api/services/util/tools"
	"encoding/json"
	"fmt"
	"strings"
)

// 驗證 CheckMacValue
func verifyNotifyCheckMacValue(params IPOSTVO.ShipStatusNotify) bool {
	// 篩選掉 CheckMacValue 只處理其他參數
	verify := params.VerifyParams()

	verify.IBoxAddress = strings.ToLower(tools.Utf8ToUnicode(params.IBoxAddress))
	verify.IBoxName = strings.ToLower(tools.Utf8ToUnicode(params.IBoxName))
	verify.PostOfficeName = strings.ToLower(tools.Utf8ToUnicode(params.PostOfficeName))
	jsonData, _ := json.Marshal(verify)

	escape := strings.ReplaceAll(string(jsonData), "\\\\", "\\")


	log.Debug("壓碼字串", escape)
	log.Debug("CheckMacValue: %s", params.CheckMacValue)
	log.Debug("壓碼字串結果: %s",strings.ToUpper(tools.SHA256(escape)))
	// 檢核碼驗證
	return params.CheckMacValue == strings.ToUpper(tools.SHA256(escape))
}

// 處理郵局通知訊息
func HandleNotification(params IPOSTVO.ShipStatusNotify) error {

	// 驗證 CheckMacValue
	if verifyNotifyCheckMacValue(params) == false {
		return fmt.Errorf("%s", "系統錯誤 檢核碼不一致")
	}

	// 通知訊息寫入貨態DB
	if err := model.NotificationInsertPostShippingStatus(params); err != nil {
		return err
	}

	return nil
}
