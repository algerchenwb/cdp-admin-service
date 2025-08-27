package modcallwrap

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/Shopify/sarama"
	"github.com/astaxie/beego"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	"gitlab.vrviu.com/inviu/backend-go-public/gopublic/kafkawrap"
	"gitlab.vrviu.com/inviu/backend-go-public/gopublic/otelwrap"
	cgsemconv "gitlab.vrviu.com/inviu/backend-go-public/gopublic/otelwrap/semconv"
	vlog "gitlab.vrviu.com/inviu/backend-go-public/gopublic/vlog/tracing"
)

var (
	// 模调异步producer名称
	_AsyncProducerName = "vrviu-modcall-async-producer"
	// 模调上报kafka topic名称
	_ModCallReportTopic = "mod_monitor_flows"
	// 初始化模调上报producer
	_OnceInitModCallProducer sync.Once
)

// VModCallIDUnique -
type VModCallIDUnique struct {
	CallerID    int    `json:"caller_id" gorm:"unique_index:idx_vmodcallid_host"`
	CaleeID     int    `json:"calee_id"  gorm:"unique_index:idx_vmodcallid_host"`
	CalleeID    int    `json:"callee_id"  gorm:"unique_index:idx_vmodcallid_host"`
	InterfaceID int    `json:"interface_id"  gorm:"unique_index:idx_vmodcallid_host"`
	CallerHost  string `json:"caller_host"  gorm:"unique_index:idx_vmodcallid_host"`
	CalleeHost  string `json:"callee_host"  gorm:"unique_index:idx_vmodcallid_host"`
}

// VModCallID -
type VModCallID struct {
	CallerID    int    `json:"caller_id"`
	CaleeID     int    `json:"calee_id"`
	CalleeID    int    `json:"callee_id"`
	InterfaceID int    `json:"interface_id"`
	CallerHost  string `json:"caller_host"`
	CalleeHost  string `json:"callee_host"`
}

func (s VModCallID) String() string {
	bs, _ := json.Marshal(s)
	return string(bs)
}

// VModCall -
type VModCall struct {
	Ctx  context.Context `json:"-"`
	Time time.Time       `json:"time"`
	VModCallID
	FlowID     string `json:"flow_id"`
	SessionID  string `json:"session_id"`
	RequestID  string `json:"request_id"`
	RetCode    int    `json:"ret_code"`
	RetMsg     string `json:"ret_msg"`
	TimeCost   int    `json:"time_cost"`
	NetSegment string `json:"-" gorm:"-"`
	Vmid       int    `json:"vmid"`
	AreaType   int    `json:"area_type"`
}

func (r *VModCall) Context() context.Context {
	if r.Ctx == nil {
		r.Ctx = otelwrap.NewSkipTraceCtx("VModCall_ctx_nil")
		vlog.Errorf(r.Ctx, "VModCall.Ctx(). ctx is null")
	}
	return r.Ctx
}

func (r *VModCall) SetContext(ctx context.Context) {
	r.Ctx = ctx
}

func (s VModCall) String() string {
	bs, _ := json.Marshal(s)
	return string(bs)
}

func (vmc *VModCall) Attributes() []attribute.KeyValue {
	attrs := []attribute.KeyValue{}

	if vmc.CallerID != 0 {
		attrs = append(attrs, CallerIDKey.Int(vmc.CallerID))
	}

	if vmc.CalleeID != 0 {
		attrs = append(attrs, CalleeIDKey.Int(vmc.CalleeID))
	}

	if vmc.InterfaceID != 0 {
		attrs = append(attrs, InterfaceIDKey.Int(vmc.InterfaceID))
	}

	if len(vmc.CallerHost) != 0 {
		attrs = append(attrs, CallerHostKey.String(vmc.CallerHost))
	}

	if len(vmc.CalleeHost) != 0 {
		attrs = append(attrs, CalleeHostKey.String(vmc.CalleeHost))
	}

	if len(vmc.FlowID) != 0 {
		attrs = append(attrs, cgsemconv.FlowID(vmc.FlowID))
	}

	if len(vmc.SessionID) != 0 {
		attrs = append(attrs, SessionIDKey.String(vmc.SessionID))
	}

	if len(vmc.RequestID) != 0 {
		attrs = append(attrs, RequestIDKey.String(vmc.RequestID))
	}

	if vmc.RetCode != 0 {
		attrs = append(attrs, RetCodeKey.Int(vmc.RetCode))
	}

	if len(vmc.RetMsg) != 0 {
		attrs = append(attrs, RetMsgKey.String(vmc.RetMsg))
	}

	if vmc.TimeCost != 0 {
		attrs = append(attrs, TimeCostKey.Int(vmc.TimeCost))
	}

	if len(vmc.NetSegment) != 0 {
		attrs = append(attrs, NetSegmentKey.String(vmc.NetSegment))
	}

	if vmc.Vmid != 0 {
		attrs = append(attrs, cgsemconv.VMID(vmc.Vmid))
	}

	if vmc.AreaType != 0 {
		attrs = append(attrs, cgsemconv.AreaType(vmc.AreaType))
	}

	return attrs
}

// NewModCallHeader 生成模调的Http头
func (s VModCall) NewModCallHeader() map[string]string {
	return map[string]string{
		"vrviu-mc-flow-id":      s.FlowID,
		"vrviu-mc-session-id":   s.SessionID,
		"vrviu-mc-caller-id":    strconv.Itoa(int(s.CallerID)),
		"vrviu-mc-callee-id":    strconv.Itoa(int(s.CalleeID)),
		"vrviu-mc-interface-id": strconv.Itoa(int(s.InterfaceID)),
		"vrviu-mc-start-time":   s.Time.Format("20060102_150405.000000000"),
	}
}

// GetModcallKafkaHosts 新增一个 section modcall ，这样就可以直接append到配置文件末尾, 批量修改配置文件
func GetModcallKafkaHosts() []string {
	rst := beego.AppConfig.DefaultString("modcall::kafka_hosts", "")
	if rst != "" {
		return strings.Split(rst, ",")
	}

	return strings.Split(beego.AppConfig.DefaultString("kafka::hosts", ""), ",")
}

// Report -
func (s VModCall) Report(calleeHost string, code int, msg, requestID string, vmid int) {
	ReportModCallWithStartTimeAndVmidWithCtx(s.Context(), calleeHost, s.CallerID, s.CalleeID, s.InterfaceID, s.FlowID, s.SessionID, requestID, code, msg, s.Time, vmid)
}

func (s VModCall) ReportWithAreaType(calleeHost string, code int, msg, requestID string, vmid, areaType int) {
	ReportModCallWithAreaTypeWithCtx(s.Context(), calleeHost, s.CallerID, s.CalleeID, s.InterfaceID, s.FlowID, s.SessionID, requestID, code, msg, s.Time, vmid, areaType)
}

// ReportModCall -
func ReportModCall(vmc VModCall) {
	if !beego.AppConfig.DefaultBool("open_modcall", true) {
		return
	}

	if !otelwrap.IsSkip(vmc.Context()) {
		span := trace.SpanFromContext(vmc.Context())
		span.AddEvent("ReportModCall", trace.WithAttributes(vmc.Attributes()...))
	}

	if !kafkawrap.IsExistProducer(_AsyncProducerName) {
		_OnceInitModCallProducer.Do(func() {
			config := kafkawrap.NewDefaultProducerConfig()
			config.Producer.RequiredAcks = sarama.NoResponse
			config.Producer.Return.Successes = false
			config.Producer.Return.Errors = false
			ctx := otelwrap.NewSkipTraceCtx("CreateExplicitKafkaProducer")
			if beego.AppConfig.DefaultBool("modcall::enable_otel", false) {
				ctx = otelwrap.WithoutOtelFlag(ctx, otelwrap.OtelFlagSkip)
			}
			kafkawrap.CreateExplicitKafkaProducerWithCtx(ctx, _AsyncProducerName, GetModcallKafkaHosts(), config)
		})
	}

	if err := kafkawrap.WriteExplicitKafkaJSONWithCtx(vmc.Context(), _AsyncProducerName, _ModCallReportTopic, vmc); err != nil {
		vlog.Errorf(vmc.Context(), "[modcall]. report mod call fail. caller(%d) callee(%d) interface(%d) for: %v", vmc.CallerID, vmc.CalleeID, vmc.InterfaceID, err)
	} else {
		vlog.Debugf(vmc.Context(), "[modcall]. report mod call success. caller(%d) callee(%d) interface(%d) timecost(%d) flow_id(%s) session_id(%s) request_id(%s) errcode(%d) errmsg(%s) vmid(%v) area_type(%v)",
			vmc.CallerID, vmc.CalleeID, vmc.InterfaceID, vmc.TimeCost, vmc.FlowID, vmc.SessionID, vmc.RequestID, vmc.RetCode, vmc.RetMsg, vmc.Vmid, vmc.AreaType)
	}
}

// ReportModCallWithStartTime -
func ReportModCallWithStartTime(calleeHost string, modCaller, modCallee, modInterface int, FlowID, SessionID, RequestID string, iRet int, iRetMsg string, starttime time.Time) {
	ReportModCallWithStartTimeAndVmid(calleeHost, modCaller, modCallee, modInterface, FlowID, SessionID, RequestID, iRet, iRetMsg, starttime, -1)
}

func ReportModCallWithStartTimeWithCtx(ctx context.Context, calleeHost string, modCaller, modCallee, modInterface int, FlowID, SessionID, RequestID string, iRet int, iRetMsg string, starttime time.Time) {
	ReportModCallWithStartTimeAndVmidWithCtx(ctx, calleeHost, modCaller, modCallee, modInterface, FlowID, SessionID, RequestID, iRet, iRetMsg, starttime, -1)
}

// ReportModCallWithStartTimeAndVmid -
func ReportModCallWithStartTimeAndVmid(calleeHost string, modCaller, modCallee, modInterface int, FlowID, SessionID, RequestID string, iRet int, iRetMsg string, starttime time.Time, vmid int) {
	ReportModCallWithAreaType(calleeHost, modCaller, modCallee, modInterface, FlowID, SessionID, RequestID, iRet, iRetMsg, starttime, vmid, -1)
}

func ReportModCallWithStartTimeAndVmidWithCtx(ctx context.Context, calleeHost string, modCaller, modCallee, modInterface int, FlowID, SessionID, RequestID string, iRet int, iRetMsg string, starttime time.Time, vmid int) {
	ReportModCallWithAreaTypeWithCtx(ctx, calleeHost, modCaller, modCallee, modInterface, FlowID, SessionID, RequestID, iRet, iRetMsg, starttime, vmid, -1)
}

func ReportModCallWithAreaType(calleeHost string, modCaller, modCallee, modInterface int, FlowID, SessionID, RequestID string, iRet int, iRetMsg string, starttime time.Time, vmid, areaType int) {
	ReportModCallWithAreaTypeWithCtx(otelwrap.NewSkipTraceCtx("ReportModCall"), calleeHost, modCaller, modCallee, modInterface, FlowID, SessionID, RequestID, iRet, iRetMsg, starttime, vmid, areaType)
}

func ReportModCallWithAreaTypeWithCtx(ctx context.Context, calleeHost string, modCaller, modCallee, modInterface int, FlowID, SessionID, RequestID string, iRet int, iRetMsg string, starttime time.Time, vmid, areaType int) {
	iTimeCost := int(time.Since(starttime).Nanoseconds() / 1000000)
	vmcid := VModCallID{
		CallerID:    modCaller,
		CalleeID:    modCallee,
		InterfaceID: modInterface,
		CallerHost:  GetCallerHost(),
		CalleeHost:  GetCalleeHost(calleeHost),
	}

	ReportModCall(VModCall{
		Ctx:        ctx,
		VModCallID: vmcid,
		FlowID:     FlowID,
		SessionID:  SessionID,
		RequestID:  RequestID,
		RetCode:    iRet,
		RetMsg:     iRetMsg,
		TimeCost:   iTimeCost,
		Time:       starttime,
		Vmid:       vmid,
		AreaType:   areaType,
	})
}

// GetCallerHost 获取主调的ip:port
// 要求HOST环境变量记录物理机ip
// 端口采用容器内运行的端口， TODO: 这个端口可能需要做转换
func GetCallerHost() string {
	return GetCallerHostWithCtx(otelwrap.NewSkipTraceCtx("GetCallerHost"))
}

func GetCallerHostWithCtx(ctx context.Context) string {
	ip := os.Getenv("HOST")
	if ip == "" {
		ip = GetPhysicalIPWithCtx(ctx)
	}
	return fmt.Sprintf("%v:%v", ip, beego.BConfig.Listen.HTTPPort)
}

// GetCalleeHost 将 域名:port 形式的地址转化为 ip:port的形式
func GetCalleeHost(host string) string {
	if host == "" {
		return host
	}
	hostSli := strings.Split(host, ",")
	if len(hostSli) < 2 {
		return GetSingleCalleeHost(host)
	}
	first := GetSingleCalleeHost(hostSli[0])
	hostSli[0] = first
	return strings.Join(hostSli, ",")

}

// GetSingleCalleeHost -
func GetSingleCalleeHost(host string) string {
	if host == "" {
		return host
	}
	ipport := host
	if strings.Contains(ipport, "://") {
		ipport = strings.SplitAfter(ipport, "://")[1]
	}
	sli := strings.Split(ipport, ":")
	if len(sli) != 2 {
		return host
	}

	ip := GetIPByDomain(sli[0])
	return ip + ":" + sli[1]
}

// GetIPByDomain 根据域名获取ip，没找到则直接返回域名
func GetIPByDomain(name string) string {
	ipSli, err := net.LookupIP(name)
	if err != nil || len(ipSli) == 0 {
		return name
	}
	return ipSli[0].String()
}

// GetInterfaceIPMap -
func GetInterfaceIPMap() (map[string]string, error) {
	ips := make(map[string]string)

	interfaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	for _, i := range interfaces {
		if (i.Flags & net.FlagUp) == 0 {
			continue
		}
		byName, err := net.InterfaceByName(i.Name)
		if err != nil {
			return nil, err
		}
		addresses, err := byName.Addrs()
		for _, v := range addresses {
			if ipnet, ok := v.(*net.IPNet); ok {
				if ipnet.IP.To4() != nil {
					ips[byName.Name] = ipnet.IP.String()
				}
			}
		}
	}
	return ips, nil
}

// GetPhysicalIP 优先先过滤掉127.0.0.1 和 docker 的ip
func GetPhysicalIP() (ip string) {
	return GetPhysicalIPWithCtx(otelwrap.NewSkipTraceCtx("GetPhysicalIP"))
}

func GetPhysicalIPWithCtx(ctx context.Context) (ip string) {
	var name string
	// defer func() {
	// 	gopublic.PrintTraceLog("[%v] name %v ip %v", fileline.FileLineFuncMulti(1), name, ip)
	// }()
	ipMap, err := GetInterfaceIPMap()
	if err != nil {
		vlog.Errorf(ctx, "GetInterfaceIPMap failed for:%v", err)
		return ""
	}

	prefixListStr := beego.AppConfig.DefaultString("net_interface_prefix_list", "en,eth")
	prefixList := strings.Split(prefixListStr, ",")
	for name, ip = range ipMap {
		lowerName := strings.ToLower(name)
		for _, pre := range prefixList {
			if strings.HasPrefix(lowerName, pre) {
				return ip
			}
		}
	}

	blackPrefixListStr := beego.AppConfig.DefaultString("net_interface_black_prefix_list", "lo,docker,vir,br")
	prefixList = strings.Split(blackPrefixListStr, ",")
	for name, ip = range ipMap {
		lowerName := strings.ToLower(name)
		for _, pre := range prefixList {
			if strings.HasPrefix(lowerName, pre) {
				continue
			}
			return ip
		}
	}
	return ""
}
