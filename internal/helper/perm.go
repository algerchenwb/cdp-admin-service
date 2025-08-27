package helper

import (
	"encoding/base64"
	"strings"
)

func PermStandardBase64(perm string) string {
	var standardPerm string
	if strings.HasPrefix(perm, "/") {
		standardPerm = perm
	} else {
		standardPerm = "/" + perm
	}
	return base64.StdEncoding.EncodeToString([]byte(standardPerm))
}

func PermStandard(perm string) string {
	var standardPerm string
	if strings.HasPrefix(perm, "/") {
		standardPerm = perm
	} else {
		standardPerm = "/" + perm
	}
	return standardPerm
}
