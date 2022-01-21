package SCSBank

import (
	"reflect"
	"testing"
)

func TestNewIncomeRecord(t *testing.T) {
	type args struct {
		content string
	}
	tests := []struct {
		name       string
		args       args
		wantRecord IncomeRecord
	}{
		{ name: "基本測試",
			args: args{content: "24102000077777  :20201023:00002 :20201023:100200:0722030000000441:+0000000002   :C:FT              :011123456789XXXX   :20201023:6317"}, wantRecord: IncomeRecord{
			RAccount: "24102000077777",
			SAccount: "011123456789XXXX",
			BTXDate:  "20201023",
			BUSDate:  "20201023",
			SeqNo:    "00002",
			TXDate:   "20201023",
			TXTime:   "100200",
			VAccount: "0722030000000441",
			Amount:   2,
			DC:       "C",
			Notes:    "FT",
			Valid:    true,
			Checksum: "6317",
			Raw:      "24102000077777  :20201023:00002 :20201023:100200:0722030000000441:+0000000002   :C:FT              :011123456789XXXX   :20201023:6317",
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRecord := NewIncomeRecord(tt.args.content); !reflect.DeepEqual(gotRecord, tt.wantRecord) {
				t.Errorf("NewIncomeRecord() = %v, want %v", gotRecord, tt.wantRecord)
			}
		})
	}
}
