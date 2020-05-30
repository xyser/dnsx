package network

import (
	"reflect"
	"testing"
)

func TestCidr_GetCidrIpRange(t *testing.T) {
	CIDR, err := NewCidr("10.6.2.2/20")
	if err != nil {
		t.Error(err)
	}
	min, max := CIDR.GetCidrIPRange()

	if min == "" {
		t.Error("Min error")
	} else {
		t.Log("Min: ", min)
	}

	if max == "" {
		t.Error("Max error")
	} else {
		t.Log("Max: ", max)
	}
}

func TestCidr_GetCidrHostNum(t *testing.T) {
	CIDR, err := NewCidr("10.6.2.2/20")
	if err != nil {
		t.Error(err)
	}
	count := CIDR.GetCidrHostNum()

	if count == 0 {
		t.Error("Count is nil")
	} else {
		t.Log("Host Num: ", count)
	}
}

func TestCidr_GetCidrIpMask(t *testing.T) {
	CIDR, err := NewCidr("10.6.2.2/20")
	if err != nil {
		t.Error(err)
	}
	data := CIDR.GetCidrIPMask()

	if data == "" {
		t.Error("Netmask is nil")
	} else {
		t.Log(data)
	}
}

func TestNewCidr1(t *testing.T) {
	type args struct {
		ipRange string
	}
	tests := []struct {
		name    string
		args    args
		want    *CIDR
		wantErr bool
	}{
		{
			name:    "a",
			args:    args{ipRange: "10.6.2.2/20"},
			want:    &CIDR{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewCidr(tt.args.ipRange)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewCidr() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewCidr() got = %v, want %v", got, tt.want)
			}
		})
	}
}
