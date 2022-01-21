package MartHiLife

import (
	"api/services/VO/HiLifeMart"
	"api/services/util/tools"
	"fmt"
	"github.com/spf13/viper"
	"net/url"
	"strings"
)

// 出貨檔 API 檢查機制
func GenerateCheckSum1(parentId, eshopId, ecDcNo, ecCvs, orderNo string, hKey, hIV string) string {
	body := "&ParentId=" + parentId + "&EshopId=" + eshopId
	body += "&EcDcNo=" + ecDcNo + "&EcCvs=" + ecCvs + "&VdrOrderNo=" + orderNo
	return privateGenerateCheckSum(body, hKey, hIV)
}

// 閉轉通知 Api 檢查checkMac
func SwitchCheckSum(s HiLifeMart.SwitchBody) error {
	hKey, hIV := viper.GetString("MartHiLife.HashKey"), viper.GetString("MartHiLife.HashIV")
	body := "&ParentId=" + s.ParentId + "&EshopId=" + s.EshopId + "&OrderNo=" + s.OrderNo
	body += "&EcOrderNo=" + s.EcOrderNo + "&OriginStoreId=" + s.OriginStoreId + "&StoreType=" + s.StoreType

	if s.ChkMac == privateGenerateCheckSum(body, hKey, hIV) {
		return nil
	}

	return fmt.Errorf("檢核失敗")
}

// 關轉店通知 API 檢查機制
func GenerateCheckSum2(parentId, eshopId, ecOrderNo, storeType, originStoreId, orderNo string, hKey, hIV string) string {
	body := "&ParentId=" + parentId + "&EshopId=" + eshopId + "&OrderNo=" + orderNo
	body += "&EcOrderNo=" + ecOrderNo + "&OriginStoreId=" + originStoreId + "&StoreType=" + storeType
	return privateGenerateCheckSum(body, hKey, hIV)
}

// 接收關轉通知 API 檢查機制
func GenerateCheckSum3(parentId, eshopId, ecDcNo, ecCvs, orderNo string, hKey, hIV string) string {
	body := "&ParentId=" + parentId + "&EshopId=" + eshopId
	body += "&EcDcNo=" + ecDcNo + "&EcCvs=" + ecCvs + "&OrderNo=" + orderNo
	return privateGenerateCheckSum(body, hKey, hIV)
}

func privateGenerateCheckSum(body, hKey, hIV string) string {
	head := "HashKey=" + hKey
	foot := "&HashIv=" + hIV
	dataStr := url.QueryEscape(head + body + foot)
	loDataStr := strings.ToLower(dataStr)
	hashDataStr := tools.MD5(loDataStr)
	UpHashDataStr := strings.ToUpper(hashDataStr)

	return UpHashDataStr
}
