// 判断host是IP或域名
package helper

import (
	"net"
)

func IsIP(host string) bool {
	return net.ParseIP(host) != nil
}
