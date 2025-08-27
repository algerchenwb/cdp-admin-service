package utils

// import (
// 	"context"
// 	"crypto/md5"
// 	"encoding/json"
// 	"fmt"
// 	"math/rand"
// 	"time"

// 	"cdp-admin-serviceee//common/globalkey"
// 	"cdp-admin-servicee//model"

// 	"github.com/golang-jwt/jwt"
// 	"github.com/zeromicro/go-zero/core/logx"
// )

// func init() {
// 	rand.Seed(time.Now().UnixNano())
// }

// func MD5(str string) string {
// 	data := []byte(str)
// 	has := md5.Sum(data)
// 	return fmt.Sprintf("%x", has)
// }

// func GetUserId(ctx context.Context) int64 {
// 	var uid int64
// 	if jsonUid, ok := ctx.Value(globalkey.SysJwtUserId).(json.Number); ok {
// 		if int64Uid, err := jsonUid.Int64(); err == nil {
// 			uid = int64Uid
// 		} else {
// 			logx.WithContext(ctx).Errorf("GetUidFromCtx err : %+v", err)
// 		}
// 	}

// 	return uid
// }

// func GetSystemHost(ctx context.Context) string {
// 	var systemHost string
// 	if jsonSystemHost, ok := ctx.Value(globalkey.SysJwtSystemHost).(string); ok {
// 		systemHost = jsonSystemHost
// 	}
// 	return systemHost
// }
// func ArrayUniqueValue[T any](arr []T) []T {
// 	size := len(arr)
// 	result := make([]T, 0, size)
// 	temp := map[any]struct{}{}
// 	for i := 0; i < size; i++ {
// 		if _, ok := temp[arr[i]]; ok != true {
// 			temp[arr[i]] = struct{}{}
// 			result = append(result, arr[i])
// 		}
// 	}

// 	return result
// }

// func ArrayContainValue(arr []int64, search int64) bool {
// 	for _, v := range arr {
// 		if v == search {
// 			return true
// 		}
// 	}

// 	return false
// }

// func Intersect(slice1 []int64, slice2 []int64) []int64 {
// 	m := make(map[int64]int64)
// 	n := make([]int64, 0)
// 	for _, v := range slice1 {
// 		m[v]++
// 	}

// 	for _, v := range slice2 {
// 		times, _ := m[v]
// 		if times == 1 {
// 			n = append(n, v)
// 		}
// 	}

// 	return n
// }

// func Difference(slice1 []int64, slice2 []int64) []int64 {
// 	m := make(map[int64]int)
// 	n := make([]int64, 0)
// 	inter := Intersect(slice1, slice2)
// 	for _, v := range inter {
// 		m[v]++
// 	}

// 	for _, v := range slice1 {
// 		times, _ := m[v]
// 		if times == 0 {
// 			n = append(n, v)
// 		}
// 	}

// 	return n
// }

// // LetterRunes 随机字符串字符池
// var LetterRunes = []rune("abcdefghijklmnopqrstuvwxyz1234567890")

// func GenerateRandonString(length int) string {
// 	b := make([]rune, length)
// 	for i := range b {
// 		b[i] = LetterRunes[rand.Intn(len(LetterRunes))]
// 	}
// 	return string(b)
// }

// func CreatRequestId() string {

// 	return GenerateRandonString(10)
// }

// func GetJwtToken(sysUser *model.SysUser, AccessExpire int64, AccessSecret string) (string, error) {
// 	iat := time.Now().Unix()
// 	claims := make(jwt.MapClaims)
// 	claims["exp"] = iat + AccessExpire
// 	claims["iat"] = iat
// 	claims[globalkey.SysJwtUserId] = sysUser.Id
// 	claims[globalkey.SysJwtAccount] = sysUser.Account
// 	token := jwt.New(jwt.SigningMethodHS256)
// 	token.Claims = claims
// 	return token.SignedString([]byte(AccessSecret))
// }

// func ArrayToString[T int | int8 | int16 | int32 | int64](arr []T) string {
// 	var str string
// 	for i, v := range arr {
// 		str += fmt.Sprintf("%d", v)
// 		if i < len(arr)-1 {
// 			str += ","
// 		}
// 	}
// 	return str
// }

// func ArrayToAnyArray[T any](arr []T) []interface{} {
// 	var anyArr []interface{}
// 	for _, v := range arr {
// 		anyArr = append(anyArr, v)
// 	}
// 	return anyArr
// }
