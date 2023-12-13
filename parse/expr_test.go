package parse

import "testing"

func TestExprTrue(t *testing.T) {
	type args struct {
		expr string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{name: "yes", args: args{expr: "yes"}, want: true},
		{name: "on", args: args{expr: "on"}, want: true},
		{name: "true", args: args{expr: "true"}, want: true},
		{name: "t", args: args{expr: "t"}, want: true},
		{name: "enabled", args: args{expr: "enabled"}, want: true},
		{name: "active", args: args{expr: "active"}, want: true},
		{name: "allow", args: args{expr: "allow"}, want: true},
		{name: "1", args: args{expr: "1"}, want: true},
		// false
		{name: "0", args: args{expr: "0"}, want: false},
		{name: "deny", args: args{expr: "deny"}, want: false},
		{name: "dead", args: args{expr: "dead"}, want: false},
		{name: "disabled", args: args{expr: "disabled"}, want: false},
		{name: "f", args: args{expr: "f"}, want: false},
		{name: "false", args: args{expr: "false"}, want: false},
		{name: "off", args: args{expr: "off"}, want: false},
		{name: "no", args: args{expr: "no"}, want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ExprTrue(tt.args.expr); got != tt.want {
				t.Errorf("ExprTrue() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestExprSec(t *testing.T) {
	type args struct {
		expr string
	}
	tests := []struct {
		name    string
		args    args
		wantSec int
	}{
		{name: "sec", args: args{expr: "1231253s"}, wantSec: 1231253},
		{name: "min-sec", args: args{expr: "4m30s"}, wantSec: 270},
		{name: "hour-min-sec", args: args{expr: "4h2m22s"}, wantSec: 14542},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotSec := ExprSec(tt.args.expr); gotSec != tt.wantSec {
				t.Errorf("ExprSec() = %v, want %v", gotSec, tt.wantSec)
			}
		})
	}
}
