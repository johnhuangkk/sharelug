package SCSBank

import (
	"api"
	"fmt"
	"reflect"
	"testing"
)

func init() {
	api.NewDevelopment()
}

func Test_numberStringToIntArray(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name    string
		args    args
		wantArr [12]int
	}{
		{name: "基本測試",args: args{str: "198765432198"}, wantArr: [12]int{1,9,8,7,6,5,4,3,2,1,9,8}},
		{name: "基本測試",args: args{str: "073312345678901"}, wantArr: [12]int{3,1,2,3,4,5,6,7,8,9,0,1}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotArr := numberStringToIntArray(tt.args.str); !reflect.DeepEqual(gotArr, tt.wantArr) {
				t.Errorf("numberStringToIntArray() = %v, want %v", gotArr, tt.wantArr)
			}
		})
	}
}

func Test_numberToIntArray(t *testing.T) {
	type args struct {
		num int
	}
	tests := []struct {
		name    string
		args    args
		wantArr [12]int
	}{
		{name: "基本測試",args: args{num: 198765432198}, wantArr: [12]int{1,9,8,7,6,5,4,3,2,1,9,8}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotArr := numberToIntArray(tt.args.num); !reflect.DeepEqual(gotArr, tt.wantArr) {
				t.Errorf("numberToIntArray() = %v, want %v", gotArr, tt.wantArr)
			}
		})
	}
}

func Test_checkSum(t *testing.T) {
	type args struct {
		money   [12]int
		account [12]int
	}
	tests := []struct {
		name    string
		args    args
		wantSum string
	}{
		{name: "基本測試",args: args{money: [12]int{0,0,0,0,0,0,0,0,0,0,0,1}, account: [12]int{2,0,2,9,6,0,0,0,0,0,2,2,}}, wantSum: "4"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotSum := checkSum(tt.args.money, tt.args.account); gotSum != tt.wantSum {
				t.Errorf("checkSum() = %v, want %v", gotSum, tt.wantSum)
			}
		})
	}
}

func Test_createAtmVirtualAccount(t *testing.T) {
	type args struct {
		prefix  string
		seqStr  string
		dateStr string
		money   int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{name: "基本測試",args: args{prefix: "0722",dateStr: "0296",seqStr: "0000022", money: 1}, want: "0722029600000224"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := createAtmVirtualAccount(tt.args.prefix, tt.args.seqStr, tt.args.dateStr, tt.args.money); got != tt.want {
				t.Errorf("createAtmVirtualAccount() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGenerateAtmVirtualAccount1(t *testing.T) {
	type args struct {
		money         int
		afterDayCount int
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{name: "測試帳號1",   args: args{money: 1,     afterDayCount: 2}},
		{name: "測試帳號2",   args: args{money: 2,     afterDayCount: 3}},
		{name: "測試帳號3",   args: args{money: 3,     afterDayCount: 4}},
		{name: "測試帳號4",   args: args{money: 4,     afterDayCount: 5}},
		{name: "測試帳號5",   args: args{money: 5,     afterDayCount: 6}},
		{name: "測試帳號6",   args: args{money: 10,     afterDayCount: 7}},
		{name: "測試帳號7",   args: args{money: 50,     afterDayCount: 8}},
		{name: "測試帳號8",   args: args{money: 100,     afterDayCount: 9}},
		{name: "測試帳號9",   args: args{money: 500,     afterDayCount: 10}},
		{name: "測試帳號10",  args: args{money: 1000,    afterDayCount: 11}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GenerateAtmVirtualAccount(tt.args.money, tt.args.afterDayCount)
			fmt.Println(tt.name," :",got, " Money:",tt.args.money, " AfterDay:",tt.args.afterDayCount)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateAtmVirtualAccount() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GenerateAtmVirtualAccount() got = %v, want %v", got, tt.want)
			}
		})
	}
}