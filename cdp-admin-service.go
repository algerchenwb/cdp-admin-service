package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"

	"cdp-admin-service/internal/config"
	"cdp-admin-service/internal/handler"
	"cdp-admin-service/internal/helper"
	table "cdp-admin-service/internal/helper/dal"
	"cdp-admin-service/internal/model/errorx"
	"cdp-admin-service/internal/svc"
	"cdp-admin-service/internal/types"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/rest/httpx"
)

var configFile = flag.String("f", "etc/cdp-admin-service-api.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)

	table.InitDBWrap(c)

	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	handler.RegisterHandlers(server, ctx)

	// 自定义错误
	httpx.SetErrorHandlerCtx(func(ctx context.Context, err error) (int, interface{}) {
		switch e := err.(type) {
		case *errorx.CodeError:
			return http.StatusOK, e.Data(helper.GetSessionId(ctx))
		default:
			return http.StatusOK, types.CommonRet{
				Ret: types.Ret{
					Code:      -1,
					Msg:       e.Error(),
					RequestId: helper.GetSessionId(ctx),
				},
				Data: nil,
			}
		}
	})
	// 自定义响应
	httpx.SetOkHandler(func(ctx context.Context, v interface{}) interface{} {
		return types.CommonRet{
			Ret: types.Ret{
				Code:      0,
				Msg:       "success",
				RequestId: helper.GetSessionId(ctx),
			},
			Data: v,
		}
	})

	if c.Mode == "dev" {
		logx.DisableStat()
	}

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
