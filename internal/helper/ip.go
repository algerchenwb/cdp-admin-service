package helper

import (
	"encoding/binary"
	"fmt"
	"net"
)

func IpToNum(ipStr string) (uint32, error) {
	ip := net.ParseIP(ipStr).To4()
	if ip == nil {
		return 0, fmt.Errorf("invalid IPv4 address: %s", ipStr)
	}
	return binary.BigEndian.Uint32(ip), nil
}

func NumToIP(ipInt uint32) string {
	ip := make(net.IP, 4)
	binary.BigEndian.PutUint32(ip, ipInt)
	return ip.String()
}
