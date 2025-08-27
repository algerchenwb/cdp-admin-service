package svc

import (
	"cdp-admin-service/internal/config"
	"cdp-admin-service/internal/helper/cache"
	"cdp-admin-service/internal/helper/queue"
	"cdp-admin-service/internal/middleware"

	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
	"gitlab.vrviu.com/cloud_esport_backend/saas_user_server/pb/saas_user"
)

type ServiceContext struct {
	Config                    config.Config
	LogwayHandleMiddleware    rest.Middleware
	AuthInterceptorMiddleware rest.Middleware
	MenuInterceptorMiddleware rest.Middleware
	Cache                     *cache.Cache
	UserRpc                   saas_user.UserServiceClient
	AgentRpc                  saas_user.AgentServiceClient
	PrimaryRpc                saas_user.PrimaryServiceClient
}

func NewServiceContext(c config.Config) *ServiceContext {
	cache := cache.NewCache(c)
	q := queue.InitLogQueue()
	conn := zrpc.MustNewClient(c.UserRpc).Conn()
	return &ServiceContext{
		Config:                    c,
		Cache:                     cache,
		LogwayHandleMiddleware:    middleware.NewLogwayHandleMiddleware(cache, q).Handle,
		AuthInterceptorMiddleware: middleware.NewAuthInterceptorMiddleware(cache).Handle,
		MenuInterceptorMiddleware: middleware.NewMenuInterceptorMiddleware(cache).Handle,
		UserRpc:                   saas_user.NewUserServiceClient(conn),
		AgentRpc:                  saas_user.NewAgentServiceClient(conn),
		PrimaryRpc:                saas_user.NewPrimaryServiceClient(conn),
	}
}

func (s *ServiceContext) IsAdmin(userId int64) bool {
	return userId == 1
}
