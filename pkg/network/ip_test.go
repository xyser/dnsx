package network

import (
	"net"
	"reflect"
	"testing"
)

func TestToIP(t *testing.T) {
	type args struct {
		i int64
	}
	tests := []struct {
		name string
		args args
		want net.IP
	}{
		{name: "test/toIP", args: args{i: 3232286465}, want: net.IPv4(192, 168, 199, 1)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ToIP(tt.args.i); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ToIP() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToInt(t *testing.T) {
	type args struct {
		ip net.IP
	}
	tests := []struct {
		name string
		args args
		want int64
	}{
		{name: "test/toInt", args: args{ip: net.IPv4(192, 168, 199, 1)}, want: 3232286465},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ToInt(tt.args.ip); got != tt.want {
				t.Errorf("ToInt() = %v, want %v", got, tt.want)
			}
		})
	}
}
