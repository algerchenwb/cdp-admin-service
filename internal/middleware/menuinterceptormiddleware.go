package middleware

import (
	"cdp-admin-service/internal/helper"
	"cdp-admin-service/internal/helper/cache"
	"cdp-admin-service/internal/model/errorx"
	"fmt"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"
)

type MenuInterceptorMiddleware struct {
	cache *cache.Cache
}

func NewMenuInterceptorMiddleware(cache *cache.Cache) *MenuInterceptorMiddleware {
	return &MenuInterceptorMiddleware{
		cache: cache,
	}
}

func (m *MenuInterceptorMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		path := helper.PermStandardBase64(r.URL.Path)
		ctx := r.Context()
		isAdmin := helper.GetIsAdmin(ctx)
		roleId := helper.GetRoleId(ctx)
		userId := helper.GetUserId(ctx)
		sessionId := helper.GetSessionId(ctx)
		if isAdmin {
			next(w, r)
			return
		}
		isPrivate, err := m.cache.Sismember(ctx, sessionId, cache.PermIsPrivateKey(), path)
		if err != nil || isPrivate {
			logx.WithContext(ctx).Errorf("[%s] user[%d] not have perm[%s] err:%+v", sessionId, userId, path, err)
			http.Error(w, errorx.UnAccessError(sessionId), http.StatusOK)
			return
		}

		roleIsAdmin, err := m.cache.Sismember(ctx, sessionId, cache.RoleAdminKey(), fmt.Sprint(roleId))
		if err != nil {
			logx.WithContext(ctx).Errorf("[%s] user[%d] not have perm[%s] err:%+v", sessionId, userId, path, err)
			http.Error(w, errorx.UnAccessError(sessionId), http.StatusOK)
			return
		}
		if roleIsAdmin {
			next(w, r)
			return
		}

		is, err := m.cache.Sismember(ctx, sessionId, cache.UserPermKey(userId), path)
		if err != nil || !is {
			logx.WithContext(ctx).Errorf("[%s] user[%d] not have perm[%s] err:%+v", sessionId, userId, path, err)
			http.Error(w, errorx.UnAccessError(sessionId), http.StatusOK)
			return
		}
		next(w, r)

	}
}
