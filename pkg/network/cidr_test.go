package network

import (
	"reflect"
	"testing"
)

func TestCidr_GetCidrIpRange(t *testing.T) {
	CIDR, err := NewCIDR("10.6.2.2/20")
	if err != nil {
		t.Error(err)
	}
	min, max := CIDR.GetCIDRIPRange()

	if min == "" {
		t.Error("Min error")
	}

	if max == "" {
		t.Error("Max error")
	}
}

func TestCidr_GetCidrHostNum(t *testing.T) {
	CIDR, err := NewCIDR("10.6.2.2/20")
	if err != nil {
		t.Error(err)
	}
	count := CIDR.GetCIDRHostNum()

	if count == 0 {
		t.Error("Count is nil")
	}
}

func TestCidr_GetCidrIpMask(t *testing.T) {
	CIDR, err := NewCIDR("10.6.2.2/20")
	if err != nil {
		t.Error(err)
	}
	data := CIDR.GetCIDRIPMask()

	if data == "" {
		t.Error("Netmask is nil")
	}
}

func TestNewCidr1(t *testing.T) {
	type args struct {
		ipRange string
	}

	want, _ := NewCIDR("10.6.2.2/20")
	tests := []struct {
		name    string
		args    args
		want    *CIDR
		wantErr bool
	}{
		{
			name:    "10.6.2.2/20",
			args:    args{ipRange: "10.6.2.2/20"},
			want:    want,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewCIDR(tt.args.ipRange)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewCidr() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewCIDR() got = %v, want %v", got, tt.want)
			}
		})
	}
}
