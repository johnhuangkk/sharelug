package OKMart

func GenerateOkReturnToString(code string) string {
	if "T00" == code {
		return " 正常驗退"
	}
	if "T01" == code {
		return " 閉店、整修、無路線路順 T02 : 無進貨資料"
	}
	if "T03" == code {
		return " 條碼錯誤"
	}
	if "T08" == code {
		return " 超材"
	}
	if "D04" == code {
		return " 大物流包裝不良(滲漏)"
	}
	if "S06" == code {
		return " 小物流破損"
	}
	if "S07" == code {
		return " 門市反應商品包裝不良(滲漏)"
	}
	return ""
}

func GenerateOkF10ToString(code string) string {
	if "E02" == code {
		return "踢退 ECNO(ECNO) 廠商不存在"
	}
	if "E03" == code {
		return "踢退 ODNO(STNO) 必須為整數"
	}
	if "E04" == code {
		return "踢退訂單編號 11 碼全為 0 錯誤"
	}

	if "E06" == code {
		return "踢退無店舖編號"
	}
	if "E07" == code {
		return "踢退門市閉店"
	}
	if "E08" == code {
		return "踢退店編前一碼異常"
	}
	if "E09" == code {
		return "踢退代收金額有誤 AMT 必須為數值"
	}
	if "E10" == code {
		return "踢退代收金額有誤 AMT 必須為 0 或正整數"
	}
	if "E11" == code {
		return "踢退商品實際金額有誤 REALAMT 必須為數值"
	}
	if "E12" == code {
		return "踢退商品實際金額有誤 REALAMT 必須為 0 或正整數"
	}
	if "E13" == code {
		return "踢退 EC 不可使用 SERCODE"
	}
	if "E14" == code {
		return "踢退無 TRADETYPE"
	}
	if "E15" == code {
		return "踢退訂單不存在"
	}
	if "E16" == code {
		return "踢退訂單狀態異常"
	}
	if "E17" == code {
		return "踢退 TRADETYPE 錯誤"
	}

	if "000" == code {
		return "成功"
	}
	return ""
}

func GenerateOkBarCodeToShipNo(bc1,bc2 string) string {
	if len(bc1) != 9 || len(bc2) != 16  {
		return ""
	}

	prefix := bc1[3:6]
	subfix := bc2[0:8]
	return prefix + subfix
}

