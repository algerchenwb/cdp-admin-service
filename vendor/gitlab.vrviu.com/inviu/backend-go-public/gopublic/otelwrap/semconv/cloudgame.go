package semconv

import (
	"fmt"

	"go.opentelemetry.io/otel/attribute"
)

const (
	CGFlowID             = attribute.Key("flow_id")
	CGSessionID          = attribute.Key("session_id")
	CGAreaType           = attribute.Key("area")
	CGIDC                = attribute.Key("idc")
	CGVMID               = attribute.Key("vmid")
	CGURCode             = attribute.Key("urcode")
	CGBizID              = attribute.Key("biz")
	CGUserID             = attribute.Key("uuid")
	CGAppUserID          = attribute.Key("app_userid")
	CGTargetUserID       = attribute.Key("target_uuid")
	CGGameID             = attribute.Key("ugid")
	CGEvent              = attribute.Key("cg.event")
	CGDeviceID           = attribute.Key("devid")
	CGDeviceType         = attribute.Key("dev_type")
	CGChargeID           = attribute.Key("cid")
	CGChargeThirdID      = attribute.Key("third_cid")
	CGCloudArchiveName   = attribute.Key("cloud_archive_name")
	CGReportEventKey     = attribute.Key("report.event")
	CGReportEventDescKey = attribute.Key("report.desc")
	CGRequestArgsKey     = attribute.Key("request.args")
)

func FlowID(flowID string) attribute.KeyValue {
	return CGFlowID.String(flowID)
}

func SessionID(sessionID string) attribute.KeyValue {
	return CGSessionID.String(sessionID)
}

func AreaType(areaType int) attribute.KeyValue {
	return CGAreaType.Int(areaType)
}

func IDC(idc int) attribute.KeyValue {
	return CGIDC.Int(idc)
}

func VMID(vmid interface{}) attribute.KeyValue {
	return CGVMID.String(fmt.Sprint(vmid))
}

func URCode(urcode interface{}) attribute.KeyValue {
	return CGURCode.String(fmt.Sprint(urcode))
}

func BizID(bizid int) attribute.KeyValue {
	return CGBizID.Int(bizid)
}

func UserID(uuid interface{}) attribute.KeyValue {
	return CGUserID.String(fmt.Sprint(uuid))
}

func AppUserID(appUserID string) attribute.KeyValue {
	return CGAppUserID.String(appUserID)
}

func TargetUserID(targetUUID interface{}) attribute.KeyValue {
	return CGTargetUserID.String(fmt.Sprint(targetUUID))
}

func GameID(gid interface{}) attribute.KeyValue {
	return CGGameID.String(fmt.Sprint(gid))
}

func Event(event string) attribute.KeyValue {
	return CGEvent.String(event)
}

func DeviceID(devid string) attribute.KeyValue {
	return CGDeviceID.String(devid)
}

func DeviceType(devType int) attribute.KeyValue {
	return CGDeviceType.Int(devType)
}

func ChargeID(cid interface{}) attribute.KeyValue {
	return CGChargeID.String(fmt.Sprint(cid))
}

func ChargeThirdID(thirdChargeID string) attribute.KeyValue {
	return CGChargeThirdID.String(thirdChargeID)
}

func ReportEvent(event int64) attribute.KeyValue {
	return CGReportEventKey.Int64(event)
}

func ReportEventDesc(desc string) attribute.KeyValue {
	return CGReportEventDescKey.String(desc)
}

func RequestArgs(args interface{}) attribute.KeyValue {
	return CGRequestArgsKey.String(fmt.Sprint(args))
}
