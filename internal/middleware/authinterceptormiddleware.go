package middleware

import (
	"cdp-admin-service/internal/helper/cache"
	table "cdp-admin-service/internal/helper/dal"
	"cdp-admin-service/internal/model/errorx"
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"cdp-admin-service/internal/helper"

	"github.com/zeromicro/go-zero/core/logx"
)

type AuthInterceptorMiddleware struct {
	cache *cache.Cache
}

func NewAuthInterceptorMiddleware(cache *cache.Cache) *AuthInterceptorMiddleware {
	return &AuthInterceptorMiddleware{
		cache: cache,
	}
}

func (m *AuthInterceptorMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		sessionId := helper.GetSessionId(ctx)

		token := r.Header.Get("Authorization")
		claims, err := helper.ParseToken(token)
		if err != nil {
			logx.WithContext(ctx).Errorf("[%s], parse token error: %v", sessionId, err)
			http.Error(w, errorx.UnauthorizedError(sessionId), http.StatusOK)
			return
		}
		originToken, err := m.cache.Get(ctx, sessionId, cache.UserOnlineKey(claims.UserID))
		if err != nil || originToken != token {
			logx.WithContext(ctx).Errorf("[%s] 0token err  useirId[%d] originToken[%s] != token[%s]  err: %v", sessionId, claims.UserID, originToken, token, err)
			http.Error(w, errorx.UnauthorizedError(sessionId), http.StatusOK)
			return
		}

		user, _, err := table.T_TCdpSysUserService.Query(ctx, sessionId, fmt.Sprintf("id:%d$status__ex:%d", claims.UserID, table.SysUserStatusDisable), "", "")
		if err != nil {
			logx.WithContext(ctx).Errorf("[%s] query user[%d] info error: %v", sessionId, claims.UserID, err)
			http.Error(w, errorx.UnauthorizedError(sessionId), http.StatusOK)
			return
		}
		ctx = context.WithValue(ctx, helper.UserKey{}, helper.User{
			UserId:   int64(user.Id),
			UserName: user.Nickname,
			Platform: user.Platform,
			AreaIds:  strings.Join(strings.Split(user.AreaIds, ","), ","),
			RoleId:   int32(user.RoleId),
			IsAdmin:  int32(user.IsAdmin),
		})
		SetUserId(w, int64(user.Id))
		SetUserAccount(w, user.Account)
		SetSessionId(w, sessionId)
		SetUserPlatform(w, user.Platform)

		logx.WithContext(ctx).Debugf("[%s] user info: %s", sessionId, helper.ToJSON(user))
		next(w, r.WithContext(ctx))
	}
}

func SetUserId(w http.ResponseWriter, userId int64) {
	w.Header().Set("X-User-Id", fmt.Sprintf("%d", userId))
}

func SetUserAccount(w http.ResponseWriter, account string) {
	w.Header().Set("X-User-Account", account)
}

func SetSessionId(w http.ResponseWriter, sessionId string) {
	w.Header().Set("X-Session-Id", sessionId)
}

func SetUserPlatform(w http.ResponseWriter, platform int32) {
	w.Header().Set("X-User-Platform", fmt.Sprintf("%d", platform))
}

func GetUserId(w http.ResponseWriter) int64 {
	userId, err := strconv.ParseInt(w.Header().Get("X-User-Id"), 10, 64)
	if err != nil {
		return 0
	}
	return userId
}

func GetUserAccount(w http.ResponseWriter) string {
	return w.Header().Get("X-User-Account")
}

func GetSessionId(w http.ResponseWriter) string {
	return w.Header().Get("X-Session-Id")
}

func GetUserPlatform(w http.ResponseWriter) int32 {
	platform, err := strconv.ParseInt(w.Header().Get("X-User-Platform"), 10, 32)
	if err != nil {
		return 0
	}
	return int32(platform)
}
