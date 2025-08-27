package cloudbox

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"cdp-admin-service/internal/helper"
	"cdp-admin-service/internal/helper/diskless"
	"cdp-admin-service/internal/helper/saas"
	"cdp-admin-service/internal/model/errorx"
	instance_types "cdp-admin-service/internal/proto/instance_service/types"
	"cdp-admin-service/internal/svc"
	"cdp-admin-service/internal/types"

	ut "github.com/go-playground/universal-translator"
	translations "github.com/go-playground/validator/v10/translations/zh"

	"github.com/go-playground/locales/zh"
	"github.com/go-playground/validator/v10"
	"github.com/gocarina/gocsv"
	"github.com/zeromicro/go-zero/core/logx"
)

type CloudBoxImportLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCloudBoxImportLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CloudBoxImportLogic {
	return &CloudBoxImportLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CloudBoxImportLogic) CloudBoxImport(r *http.Request) (resp *types.CloudBoxImportResp, err error) {
	sessionId := helper.GetSessionId(l.ctx)

	resp = new(types.CloudBoxImportResp)

	data := r.FormValue("shopInfo")
	var shopInfo = new(types.ShopInfoImportExt)
	if err = json.Unmarshal([]byte(data), shopInfo); err != nil {
		return nil, errorx.NewDefaultError(errorx.CloudBoxImportUnmarshalErrorCode)
	}

	logx.WithContext(l.ctx).Infof("MIME shopInfo: %+v", shopInfo)

	_ = r.ParseMultipartForm(types.MaxFileSize)
	file, handler, err := r.FormFile("myFile")
	if err != nil {
		return nil, errorx.NewDefaultError(errorx.CloudBoxImportFormFileErrorCode)
	}

	defer file.Close()

	logx.WithContext(l.ctx).Infof("[%s] Uploaded File: %s,Size:[%d] MIME Header:%+v", sessionId, handler.Filename, handler.Size, handler.Header)

	if handler.Size > types.MaxFileSize {
		return nil, errorx.NewDefaultError(errorx.CloudBoxImportFileSizeLimitErrorCode)
	}

	buf := bytes.NewBuffer(nil)
	if _, err = io.Copy(buf, file); err != nil {
		logx.WithContext(l.ctx).Errorf("[%s] io.Copy shopInfo:%+v err: %+v", sessionId, shopInfo, err)
		return nil, errorx.NewDefaultError(errorx.CloudBoxImportFileReadErrorCode)
	}

	cloudboxInfo := []*types.CloudBoxImport{}
	if err = gocsv.UnmarshalBytes(buf.Bytes(), &cloudboxInfo); err != nil {
		logx.WithContext(l.ctx).Errorf("[%s] gocsv.UnmarshalBytes shopInfo:%+v err: %+v", sessionId, shopInfo, err)
		return nil, errorx.NewDefaultError(errorx.CloudBoxImportUnmarshalBytesErrorCode)
	}

	if len(cloudboxInfo) > types.FIXED_MAX_CSV_RECORD_NUM {
		return nil, errorx.NewDefaultError(errorx.CloudBoxImportCountLimitErrorCode)
	}

	if len(cloudboxInfo) == 0 {
		return nil, errorx.NewDefaultError(errorx.CloudBoxImportCountEmptyErrorCode)

	}

	if !helper.CheckDevicesMacUnique(cloudboxInfo) {
		return nil, errorx.NewDefaultError(errorx.CloudBoxImportMacDuplicateErrorCode)
	}
	for _, info := range cloudboxInfo {

		if err := l.CheckcloudboxInfo(info); err != nil {
			return nil, errorx.NewHandlerError(errorx.ServerErrorCode, err.Error())
		}
	}

	mapPoolInfo, err := saas.NewSaasServerService(l.svcCtx.Config.Saas.ServerHost, l.svcCtx.Config.Saas.Timeout).GetBizPoolInfo(l.ctx, l.svcCtx, sessionId, shopInfo.BizId)
	if err != nil {
		return nil, errorx.NewHandlerError(errorx.ServerErrorCode, "查询合约信息失败")
	}

	for index, info := range cloudboxInfo {
		if poolInfo, ok := mapPoolInfo[int(info.SpecInstId)]; ok {
			cloudboxInfo[index].ConfigId = poolInfo.ConfigID
		}

		logx.WithContext(l.ctx).Infof("[%s] cloudboxInfo %s", sessionId, helper.ToJSON(info))
	}

	if err := l.BatchAddCloudBox(sessionId, shopInfo, cloudboxInfo); err != nil {
		return nil, errorx.NewHandlerError(errorx.ServerErrorCode, err.Error())
	}

	return resp, nil

}

func (l *CloudBoxImportLogic) BatchAddCloudBox(sessionId string, shopInfo *types.ShopInfoImportExt, cloudBoxInfos []*types.CloudBoxImport) (err error) {

	var errMsgList []string
	for _, item := range cloudBoxInfos {
		macString := helper.ConvertMacAddress(item.Mac)
		saasReq := saas.CreateDeviceInfoReq{
			MACAddress:  macString,
			UpdateState: 0,
			BizID:       shopInfo.BizId,
			DeviceName:  item.Name,
			UpdateBy:    shopInfo.AccountName,
		}
		_, err := saas.CreateEsportDeviceInfo(l.ctx, sessionId, l.svcCtx.Config.OutSide.SaasHost, saasReq)
		if err != nil {
			errMsgList = append(errMsgList, err.Error())
			continue
		}

		// 调用无盘实例的接口
		if len(item.ConfigId) == 0 {
			item.ConfigId = "0"
		}
		schemeId, _ := strconv.ParseInt(item.ConfigId, 10, 64)

		instReq := &instance_types.CreateInstanceRequest{
			FlowId:     sessionId,
			DeviceType: 2, // 0-云主机 1-本地主机 2-云盒 3-本地盒子
			Vlan:       int(shopInfo.VlanId),
			InstanceInfo: instance_types.InstanceInfo{
				BootMac:  item.Mac,
				BootType: item.Mode,
				NetInfo: instance_types.NetInfo{
					Ip:       item.IP,
					Hostname: item.Name,
				},
				DefaultConfig: instance_types.DefaultConfig{
					NetInfo: instance_types.NetInfo{
						Ip:       item.IP,
						Hostname: item.Name,
					},
				},
				State:    100,
				Remark:   item.Name,
				SchemeId: schemeId,
			},
		}

		if err := diskless.NewDisklessWebGateway(l.ctx, l.svcCtx).CreateInstance(shopInfo.AreaId, instReq); err != nil {
			errMsgList = append(errMsgList, err.Error())
		}
	}
	if len(errMsgList) != 0 {
		err = errors.New(strings.Join(errMsgList, ","))
	}
	return
}

func (l *CloudBoxImportLogic) CheckcloudboxInfo(info *types.CloudBoxImport) (err error) {

	validate := validator.New()
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := fld.Tag.Get("label")
		return name
	})

	trans, _ := ut.New(zh.New()).GetTranslator("zh")
	translations.RegisterDefaultTranslations(validate, trans)
	if validateErr := validate.StructCtx(l.ctx, info); validateErr != nil {
		for _, err1 := range validateErr.(validator.ValidationErrors) {
			err = errors.New(info.Name + ":" + err1.Translate(trans))
			logx.Error(err)
			return err
		}
	}

	// 检查字符串是否只包含ASCII字符
	if !regexp.MustCompile(`^[\x00-\x7F]+$`).MatchString(info.Name) {
		return errors.New("包含ASCII字符")
	}

	// 检查字符串是否只包含大小写字母、数字、中划线(-)、英文字符点(.)
	if !regexp.MustCompile(`^[a-zA-Z0-9\-.]+$`).MatchString(info.Name) {
		return errors.New("包含非法字符")
	}

	// 检查字符串是否只有一个字符且为.时，不满足
	if info.Name == "." {
		return errors.New("包含非法字符")
	}
	return
}
