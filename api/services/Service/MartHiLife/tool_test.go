package MartHiLife

import "testing"

func Test_privateCheckSum1(t *testing.T) {
	type args struct {
		parentId string
		eshopId  string
		ecDcNo   string
		ecCvs    string
		orderNo  string
		hKey     string
		hIV      string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"基本測試" ,args{
			parentId: "124",
			eshopId:  "901",
			ecDcNo:   "D11",
			ecCvs:    "HILIFEC2C",
			orderNo:  "20201120001H",
			hKey:     "5593df4f4d10",
			hIV:      "b41752b9acbd",
		}, "4B1223532ACF6FF1E50B52D23E03A558"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GenerateCheckSum1(tt.args.parentId, tt.args.eshopId, tt.args.ecDcNo, tt.args.ecCvs, tt.args.orderNo, tt.args.hKey, tt.args.hIV); got != tt.want {
				t.Errorf("privateGenerateCheckSum() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_privateGenerateCheckSum2(t *testing.T) {
	type args struct {
		parentId      string
		eshopId       string
		ecOrderNo     string
		storeType     string
		originStoreId string
		orderNo       string
		hKey          string
		hIV           string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{ name: "基本測試", args: args{
			parentId:      "012",
			eshopId:       "001",
			ecOrderNo:     "01200001234",
			storeType:     "2",
			originStoreId: "3750",
			orderNo:       "19HAK53243004",
			hKey:          "e7e10c28d3ec",
			hIV:           "513cf55d4727",
		}, want: "4A11142870357BB72DF19D0AC7C9BCEF"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GenerateCheckSum2(tt.args.parentId, tt.args.eshopId, tt.args.ecOrderNo, tt.args.storeType, tt.args.originStoreId, tt.args.orderNo, tt.args.hKey, tt.args.hIV); got != tt.want {
				t.Errorf("privateGenerateCheckSum2() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_privateGenerateCheckSum3(t *testing.T) {
	type args struct {
		parentId string
		eshopId  string
		ecDcNo   string
		ecCvs    string
		orderNo  string
		hKey     string
		hIV      string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{ name: "基本測試", args: args{
			parentId:      "124",
			eshopId:       "901",
			ecCvs:         "HILIFEC2C",
			ecDcNo:        "D11",
			orderNo:       "20RKB40416476",
			hKey:          "5593df4f4d10",
			hIV:           "b41752b9acbd",
		}, want: "778561D74C80D429A0AB3BF5DDF863E8"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GenerateCheckSum3(tt.args.parentId, tt.args.eshopId, tt.args.ecDcNo, tt.args.ecCvs, tt.args.orderNo, tt.args.hKey, tt.args.hIV); got != tt.want {
				t.Errorf("GenerateCheckSum3() = %v, want %v", got, tt.want)
			}
		})
	}
}