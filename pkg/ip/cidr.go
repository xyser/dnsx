package ip

import (
	"encoding/binary"
	"errors"
	"fmt"
	"net"
	"strconv"
	"strings"
)

// ErrNotSupportIPV6 不支持 IP V6
var ErrNotSupportIPV6 = errors.New("does not support IPv6 addresses")

type Cidr struct {
	CidrIpRange string

	maskLen int
	ipBase  string

	IP    net.IP
	IPNet *net.IPNet
}

func NewCidr(ipRange string) (*Cidr, error) {
	var cidr Cidr
	var err error

	if cidr.IP, cidr.IPNet, err = net.ParseCIDR(ipRange); err != nil {
		return nil, err
	}

	cidr.maskLen, _ = cidr.IPNet.Mask.Size()
	return &cidr, nil
}

// Contains 判断 CIDR 是否包含 某个IP
func (c *Cidr) Contains(ip string) (ok bool, err error) {
	oip := net.ParseIP(ip)
	return c.IPNet.Contains(oip), nil
}

// LastAddr 获取一个 CIDR 的广播地址
func (c *Cidr) LastAddr() (net.IP, error) { // works when the n is a prefix, otherwise...
	if c.IPNet.IP.To4() == nil {
		return net.IP{}, ErrNotSupportIPV6
	}
	ip := make(net.IP, len(c.IPNet.IP.To4()))
	binary.BigEndian.PutUint32(ip, binary.BigEndian.Uint32(c.IPNet.IP.To4())|^binary.BigEndian.Uint32(net.IP(c.IPNet.Mask).To4()))
	return ip, nil
}

// GetCidrIpRange 获取最大主机IP和最小主机IP
func (c *Cidr) GetCidrIpRange() (min, max string) {
	ip := strings.Split(c.CidrIpRange, "/")[0]
	ipSeg := strings.Split(ip, ".")
	maskLen := c.GetMaskLen()
	seg3MinIp, seg3MaxIp := getIpSeg3Range(ipSeg, maskLen)
	seg4MinIp, seg4MaxIp := getIpSeg4Range(ipSeg, maskLen)
	ipPrefix := ipSeg[0] + "." + ipSeg[1] + "."

	min = ipPrefix + strconv.Itoa(seg3MinIp) + "." + strconv.Itoa(seg4MinIp)
	max = ipPrefix + strconv.Itoa(seg3MaxIp) + "." + strconv.Itoa(seg4MaxIp)
	return min, max
}

// GetCidrHostNum CIDR地址 范围内主机数量
func (c *Cidr) GetCidrHostNum() uint {
	cidrIpNum := uint(0)
	var i = uint(32 - c.GetMaskLen() - 1)
	for ; i >= 1; i-- {
		cidrIpNum += 1 << i
	}
	return cidrIpNum
}

// GetMaskLen CIDR地址 掩码长度
func (c *Cidr) GetMaskLen() int {
	return c.maskLen
}

// GetCidrIpMask 获取CIDR掩码
func (c *Cidr) GetCidrIpMask() string {
	// ^uint32(0)二进制为32个比特1，通过向左位移，得到CIDR掩码的二进制
	cidrMask := ^uint32(0) << uint(32-c.GetMaskLen())
	// 计算CIDR掩码的四个片段，将想要得到的片段移动到内存最低8位后，将其强转为8位整型，从而得到
	cidrMaskSeg1 := uint8(cidrMask >> 24)
	cidrMaskSeg2 := uint8(cidrMask >> 16)
	cidrMaskSeg3 := uint8(cidrMask >> 8)
	cidrMaskSeg4 := uint8(cidrMask & uint32(255))

	return fmt.Sprint(cidrMaskSeg1) + "." + fmt.Sprint(cidrMaskSeg2) + "." + fmt.Sprint(cidrMaskSeg3) + "." + fmt.Sprint(cidrMaskSeg4)
}

// getIpSeg3Range 得到第三段IP的区间（第一片段.第二片段.第三片段.第四片段）
func getIpSeg3Range(ipSeg []string, maskLen int) (int, int) {
	if maskLen > 24 {
		segIp, _ := strconv.Atoi(ipSeg[2])
		return segIp, segIp
	}
	seg, _ := strconv.Atoi(ipSeg[2])
	return getIpSegRange(uint8(seg), uint8(24-maskLen))
}

// getIpSeg4Range 得到第四段IP的区间（第一片段.第二片段.第三片段.第四片段）
func getIpSeg4Range(ipSeg []string, maskLen int) (int, int) {
	seg, _ := strconv.Atoi(ipSeg[3])
	segMinIp, segMaxIp := getIpSegRange(uint8(seg), uint8(32-maskLen))
	return segMinIp + 1, segMaxIp
}

// getIpSegRange 根据用户输入的基础IP地址和CIDR掩码计算一个IP片段的区间
func getIpSegRange(userSegIp, offset uint8) (int, int) {
	var ipSegMax uint8 = 255
	netSegIp := ipSegMax << offset
	segMinIp := netSegIp & userSegIp
	segMaxIp := userSegIp&(255<<offset) | ^(255 << offset)
	return int(segMinIp), int(segMaxIp)
}
