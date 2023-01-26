package main

import (
	"testing"
)

func TestStringParser(t *testing.T) {

	type args struct {
		str []string
	}

	tests := []struct {
		name         string
		args         args
		wantCategory string
		wantSum      string
		wantCur      string
	}{
		{
			name: "Test1",
			args: args{
				str: []string{"Продукты", "100Rub"},
			},
			wantCategory: "Продукты",
			wantSum:      "100",
			wantCur:      "Rub",
		},
		{
			name: "Test2",
			args: args{
				str: []string{"Услуги", "100USD"},
			},
			wantCategory: "Услуги",
			wantSum:      "100",
			wantCur:      "USD",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			category, sum, cur := stringParser(tt.args.str)
			if category != tt.wantCategory || sum != tt.wantSum || cur != tt.wantCur {
				t.Errorf("%s: Провал \n", tt.name)
			}
		})
	}
}
