package helper

import (
	"cdp-admin-service/internal/types"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net"
	"regexp"
	"strconv"
	"strings"
)

func SliceToString[T any](slice []T) string {
	var result string
	for _, v := range slice {
		result += fmt.Sprintf("%v,", v)
	}
	if len(result) > 0 {
		result = result[:len(result)-1]
	}
	return result
}

func StringToInt64Slice(str string) []int64 {
	var result []int64
	for _, v := range strings.Split(str, ",") {
		vv, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			continue
		}
		result = append(result, vv)
	}
	return result
}

func ToJSON(object interface{}) string {
	bytes, _ := json.Marshal(object)
	return string(bytes)
}

// 检查 IP 地址类型和私有性
func CheckIP(ipStr string) (isValid, isIPv4, isPrivate bool) {
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return false, false, false
	}

	// 区分 IPv4/IPv6
	isIPv4 = ip.To4() != nil

	// 检查私有地址（兼容 Go 1.17+ 的 ip.IsPrivate()）
	if isIPv4 {
		// IPv4 私有地址范围：10.0.0.0/8, 172.16.0.0/12, 192.168.0.0/16
		ip4 := ip.To4()
		switch {
		case ip4[0] == 10:
			isPrivate = true
		case ip4[0] == 172 && ip4[1] >= 16 && ip4[1] <= 31:
			isPrivate = true
		case ip4[0] == 192 && ip4[1] == 168:
			isPrivate = true
		}
	} else {
		// IPv6 私有地址范围：fc00::/7 (ULA)
		isPrivate = len(ip) == net.IPv6len && (ip[0]&0xfe == 0xfc)
	}

	return true, isIPv4, isPrivate

}

func RemoveFromStr(str string, value string) string {
	var result string
	for _, v := range strings.Split(str, ",") {
		if v == value {
			continue
		}
		result += fmt.Sprintf("%s,", v)
	}
	if len(result) > 0 {
		result = result[:len(result)-1]
	}
	return result

}

// LetterRunes 随机字符串字符池
var LetterRunes = []rune("abcdefghijklmnopqrstuvwxyz1234567890")

var LargeLetterRunes = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ")

func GenerateRandonString(length int) string {
	b := make([]rune, length)
	for i := range b {
		b[i] = LetterRunes[rand.Intn(len(LetterRunes))]
	}
	return string(b)
}

func GenerateRandonLargeString(length int) string {
	b := make([]rune, length)
	for i := range b {
		b[i] = LargeLetterRunes[rand.Intn(len(LargeLetterRunes))]
	}
	return string(b)
}

func CreatRequestId() string {

	return GenerateRandonString(10)
}

func CreatCloudBoxName() string {

	return GenerateRandonLargeString(4)
}

func FormatMacAddress(mac string) (string, error) {
	// 检查MAC地址长度是否正确
	if len(mac) != 12 {
		return "", errors.New("Invalid MAC address length")
	}

	// 使用strings.Builder来构建新的MAC地址格式
	var builder strings.Builder
	for i, char := range mac {
		// 每两个字符后插入冒号
		if i%2 == 0 && i != 0 {
			builder.WriteString(":")
		}
		builder.WriteRune(char)
	}
	formattedMac := strings.ToLower(builder.String())
	return formattedMac, nil
}

// mac地址转换  145F01C189E4 =>14:5f:01:c1:89:e4
func ConvertMacAddress(mac string) string {
	// 移除 MAC 地址中的冒号或短横线
	mac = strings.ReplaceAll(mac, ":", "")
	mac = strings.ReplaceAll(mac, "-", "")

	// 将 MAC 地址转换为大写
	mac = strings.ToUpper(mac)

	return mac
}

// mac地址转换  14:5f:01:c1:89:e4 => 145F01C189E4
func ConvertMacAddressToString(condList []string) []string {
	for index, cond := range condList {
		if strings.Contains(cond, "mac_address__contains") {

			re := regexp.MustCompile(`(?i)mac_address__contains:([a-f0-9:]+)`)
			matches := re.FindStringSubmatch(cond)

			if len(matches) > 1 {

				macPart := matches[1]                              // 提取匹配到的 MAC 地址部分
				macPart = strings.ReplaceAll(macPart, ":", "")     // 移除 MAC 地址部分中的冒号
				remainingPart := strings.ToUpper(macPart)          // 将剩余的部分转换为大写
				output := "mac_address__contains:" + remainingPart // 重新组合字符串
				condList[index] = output
			}
		}
	}
	return condList
}

func CheckDevicesMacUnique(d []*types.CloudBoxImport) bool {
	inResult := make(map[string]bool)
	for idx := range d {
		if _, ok := inResult[d[idx].Mac]; !ok {
			inResult[d[idx].Mac] = true
		} else {
			return false
		}
	}
	return true
}

func StringToInt64(s string) int64 {
	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0
	}
	return i
}
