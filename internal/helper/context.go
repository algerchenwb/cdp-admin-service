package helper

import (
	"context"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"gitlab.vrviu.com/inviu/backend-go-public/gopublic"
)

type SessionIdKey struct{}
type UserKey struct{}
type User struct {
	UserId   int64
	UserName string
	Platform int32
	AreaIds  string
	RoleId   int32
	IsAdmin  int32
}

func GetUserId(ctx context.Context) int64 {
	user, _ := ctx.Value(UserKey{}).(User)
	return user.UserId
}

func GetUserName(ctx context.Context) string {
	user, _ := ctx.Value(UserKey{}).(User)
	return user.UserName
}

func GetPlatform(ctx context.Context) int32 {
	user, _ := ctx.Value(UserKey{}).(User)
	return user.Platform
}

func GetAreaIds(ctx context.Context) string {
	user, _ := ctx.Value(UserKey{}).(User)
	return user.AreaIds
}

const (
	PlatformSuanli  = 1
	PlatformShigong = 2
)
const (
	SystemHostSuanli  = "suanli"
	SystemHostShigong = "shigong"
)

func GetSystemHost(ctx context.Context) string {
	user, _ := ctx.Value(UserKey{}).(User)
	switch user.Platform {
	case PlatformSuanli:
		return SystemHostSuanli
	case PlatformShigong:
		return SystemHostShigong

	}
	return ""
}

func GetRoleId(ctx context.Context) int32 {
	user, _ := ctx.Value(UserKey{}).(User)
	return user.RoleId
}

func GenSessionId(ctx context.Context) (string, context.Context) {
	sessionId, _ := ctx.Value(SessionIdKey{}).(string)
	if sessionId == "" {
		sessionId = fmt.Sprintf("%s%s%d", "cdp_", time.Now().Format("20060102150405"), rand.Intn(1000))
	}
	return sessionId, context.WithValue(ctx, SessionIdKey{}, sessionId)
}

func GetSessionId(ctx context.Context) string {
	sessionId, _ := ctx.Value(SessionIdKey{}).(string)
	return sessionId
}

func GetIsAdmin(ctx context.Context) bool {
	user, _ := ctx.Value(UserKey{}).(User)
	return user.IsAdmin == 1
}

func CheckAreaId(ctx context.Context, areaId string) bool {
	user, _ := ctx.Value(UserKey{}).(User)
	if user.AreaIds == "" {
		return false
	}
	areaIdList := strings.Split(user.AreaIds, ",")
	return gopublic.StringInArray(areaId, areaIdList)
}
