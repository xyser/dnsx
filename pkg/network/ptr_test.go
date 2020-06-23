package network

import (
	"testing"
)

func TestPTRToIP(t *testing.T) {
	type args struct {
		name []byte
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{name: "1.0.0.127.in-addr.arpa.", args: args{name: []byte("1.0.0.127.in-addr.arpa.")}, want: "127.0.0.1"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := PTRToIP(tt.args.name); got != tt.want {
				t.Errorf("PTRToIP() = %v, want %v", got, tt.want)
			}
		})
	}
}
