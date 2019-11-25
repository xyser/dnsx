package ip

import "testing"

func TestCidr_GetCidrIpRange(t *testing.T) {
	min, max := NewCidr("10.6.2.2/20").GetCidrIpRange()

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
	count := NewCidr("10.6.2.2/20").GetCidrHostNum()

	if count == 0 {
		t.Error("Count is nil")
	} else {
		t.Log("Host Num: ", count)
	}
}

func TestCidr_GetCidrIpMask(t *testing.T) {
	data := NewCidr("10.6.2.2/20").GetCidrIpMask()

	if data.Netmask == "" {
		t.Error("Netmask is nil")
	} else {
		t.Log(data.Netmask)
	}
}

func TestNewCidr(t *testing.T) {
	min, max := NewCidr("10.6.2.2/20").GetCidrIpRange()

	if min == "" {
		t.Error("Min error")
	} else {
		t.Log("Min ", min)
	}

	if max == "" {
		t.Error("Max error")
	} else {
		t.Log("Max ", max)
	}
}
