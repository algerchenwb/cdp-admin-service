package spec

import (
	"context"
	"fmt"

	"cdp-admin-service/internal/helper"
	table "cdp-admin-service/internal/helper/dal"
	"cdp-admin-service/internal/helper/saas"
	"cdp-admin-service/internal/model/errorx"
	"cdp-admin-service/internal/svc"
	"cdp-admin-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
	"gitlab.vrviu.com/inviu/backend-go-public/gopublic"
)

type SpecLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSpecLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SpecLogic {
	return &SpecLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SpecLogic) Spec(req *types.SpecListReq) (resp *types.SpecListResp, err error) {
	sessionId := helper.GetSessionId(l.ctx)
	resp = &types.SpecListResp{}

	specMap := make(map[int]types.Spec)

	if req.AreaId != 0 {
		_, _, err := table.T_TCdpAreaInfoService.Query(l.ctx, sessionId, fmt.Sprintf("area_id:%d$status:1", req.AreaId), nil, nil)
		if err != nil {
			l.Logger.Errorf("[%s] T_TCdpAreaInfoService  areaId:%d, err:%v", sessionId, req.AreaId, err)
			return nil, errorx.NewDefaultError(errorx.AreaNotFound)
		}
		_, areaConfigs, err := saas.NewSaasServerService(l.svcCtx.Config.Saas.ServerHost, l.svcCtx.Config.Saas.Timeout).
			GetAreaConfigInfoList(l.ctx, sessionId, uint64(req.AreaId))
		if err != nil {
			l.Logger.Errorf("[%s] areaId:%d,  GetAreaConfigInfoList  err:%v", sessionId, req.AreaId, err)
			return nil, errorx.NewDefaultError(errorx.AreaNotFound)
		}

		areaInfo, _, err := table.T_TCdpAreaInfoService.Query(l.ctx, sessionId, fmt.Sprintf("area_id:%d", req.AreaId), nil, nil)
		if err != nil {
			l.Logger.Errorf("[%s] Table Query err. areaId[%d] err:%+v", sessionId, req.AreaId, err)
			return nil, errorx.NewDefaultError(errorx.AreaNotFound)
		}

		l.Logger.Infof("[%s] areaId:%d, areaInfo:%s", sessionId, req.AreaId, helper.ToJSON(areaInfo))

		_, pools, err := saas.NewSaasServerService(l.svcCtx.Config.Saas.ServerHost, l.svcCtx.Config.Saas.Timeout).
			GetResourcePoolInfoList(l.ctx, sessionId, []string{fmt.Sprintf("agent_id:%d$biz_id:%d$state:1", areaInfo.AgentId, 0)})
		if err != nil {
			l.Logger.Errorf("[%s] GetResourcePoolInfoList err. AgentId[%d] err:%+v", sessionId, areaInfo.AgentId, err)
			return nil, errorx.NewDefaultError(errorx.AreaNotFound)
		}

		AgentOutSpecId := make([]uint64, 0)
		for _, pool := range pools {
			AgentOutSpecId = append(AgentOutSpecId, uint64(pool.OuterSpecID))
		}

		l.Logger.Infof("[%s] areaId:%d, agentId:%d AgentOutSpecId:%s", sessionId, req.AreaId, areaInfo.AgentId, helper.ToJSON(AgentOutSpecId))

		for _, areaConfig := range areaConfigs {
			for _, spec := range areaConfig.OuterSpecInfoList {
				if !gopublic.Uint64InArray(uint64(spec.Id), AgentOutSpecId) {
					continue
				}
				specMap[int(spec.Id)] = types.Spec{
					OuterSpecId: int(spec.Id),
					InnerSpecId: int(spec.InnerSpecID),
					Name:        spec.Name,
				}
			}
		}
	}

	if req.PrimaryId != 0 {
	}
	if req.AgentId != 0 {
	}

	for _, spec := range specMap {
		resp.List = append(resp.List, spec)
	}
	resp.Total = len(resp.List)

	return
}
