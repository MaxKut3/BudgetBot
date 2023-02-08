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

func TestValidateMessageText(t *testing.T) {
	type args struct {
		wordList []string
	}

	tests := []struct {
		name  string
		args  args
		wants bool
	}{
		{
			name: "Test 1",
			args: args{
				[]string{"Магнит", "Продукты", "1000Rub"},
			},
			wants: false,
		},
		{
			name: "Test 2",
			args: args{
				[]string{"Продукты", "1000Rub"},
			},
			wants: true,
		},
		{
			name: "Test 3",
			args: args{
				[]string{"Продукты", "Rub"},
			},
			wants: false,
		},
		{
			name: "Test 4",
			args: args{
				[]string{"Продукты", "1000Rub"},
			},
			wants: true,
		},
		{
			name: "Test 5",
			args: args{
				[]string{"Продукты", "1000Руб"},
			},
			wants: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := validateMessageText(tt.args.wordList)
			if res != tt.wants {
				t.Errorf("%s: fail \n", tt.name)
			}
		})
	}
}
