package errwrap

var (
	Unknown = New(-1, "unknown") // 未知错误
	Succ    = New(0, "succ")     // 返回成功

	// 1xxx 通用错误
	// 10xx 参数相关错误
	ErrInvalidRequestParam   = New(1000, "invalid request parameter") // 请求参数异常
	ErrInvalidRequestMessage = New(1001, "invalid request message")   // 请求消息异常
	ErrInvalidToken          = New(1002, "invalid token")             // token 非法
	ErrInvalidRequest        = New(1003, "invalid request")           // 非法请求
	ErrHandleTimeout         = New(1004, "handle timeout")            // 请求处理超时
	ErrNoPermission          = New(1005, "no permission")             // 权限不足

	// 11xx 其他类型
	ErrJsonMarshal        = New(1100, "json marshal failed")       // 序列化 JSON 数据失败
	ErrJsonUnmarshal      = New(1101, "json unmarshal failed")     // 反序列化 JSON 数据失败
	ErrPBMarshal          = New(1102, "protobuf marshal failed")   // 序列化 protobuf 数据失败
	ErrPBUnmarshal        = New(1103, "protobuf unmarshal failed") // 反序列化 protobuf 数据失败
	ErrItemNotExist       = New(1104, "item not exist")            // 记录不存在
	ErrItemAlreadyExist   = New(1105, "item already exist")        // 记录已存在
	ErrInvalidConfigValue = New(1106, "invalid config value")      // 配置非法
	ErrParseISP           = New(1107, "parse isp failed")          // 解析 isp 信息失败
	ErrQueryAccessInfo    = New(1108, "query access info failed")  // 查询区域接入信息失败
	ErrUrcode             = New(1109, "urcode failed")             // urcode 属性配置异常
	ErrUrcodeImage        = New(1110, "urcode image failed")       // urcode image 属性配置异常
	ErrUrcodeDevice       = New(1110, "urcode device failed")      // urcode device 属性配置异常
	ErrQueryUrcode        = New(1111, "query urcode failed")       // 请求 resource_type 查询 urcode 属性失败
	ErrEdgeUnavailable    = New(1112, "edge unavailable failed")   // edge下线
	ErrWriteData          = New(1113, "write data failed")         // 写数据失败

	// 2xxx 业务相关错误
	// 21xx 游戏相关错误
	ErrAnotherGameIsRunning              = New(2100, "another game is running")                    // 已运行另一个游戏
	ErrGameUnavailable                   = New(2101, "unavailable game")                           // 游戏不可用
	ErrGameNotLaunch                     = New(2102, "game not launch")                            // 游戏未启动
	ErrSpecifyGameNotRunning             = New(2103, "specify game not running")                   // 指定游戏未运行
	ErrInvalidGameStatus                 = New(2104, "invalid game status")                        // 游戏状态异常
	ErrGameDurationNotEnough             = New(2105, "game duration not enough")                   // 剩余游戏时长不足
	ErrGameAlreadyRunning                = New(2106, "game already running")                       // 该游戏已运行
	ErrGameInExit                        = New(2107, "game in exit")                               // 游戏退出中
	ErrGameNotAuthorized                 = New(2108, "game not authorized")                        // 游戏未授权
	ErrNoCkdeyLimit                      = New(2109, "no ckdey limit")                             // 游戏无 cdkey 控制
	ErrNoMoreCdkey                       = New(2110, "no more cdkey")                              // 无空闲 cdkey
	ErrNoAvailableArea                   = New(2111, "no available area")                          // 无可用的游戏分区
	ErrCalcResourceTypeRate              = New(2112, "calc resource type rate failed")             // 计算资源类型比例失败
	ErrExplicitAllocatingByVmid          = New(2113, "explicit allocating by vmid")                // 禁止指定 vmid 的显示分配置换资源
	ErrVmStateError                      = New(2114, "vm state error")                             // 虚拟机状态异常
	ErrVmStreamMethodError               = New(2115, "vm stream method error")                     // 虚拟机串流类型非法
	ErrFlowIdNotMatch                    = New(2116, "flow id not match")                          // union 和区域 flow_id 不匹配
	ErrInvalidStatusToStartLive          = New(2117, "invalid status to start live")               // 当前游戏状态下无法开始直播
	ErrInvalidExtraParam                 = New(2118, "invalid extra param")                        // 游戏附加启动参数非法
	ErrAssignVm                          = New(2119, "assign vm failed")                           // 请求 vm-manager 申请 vm 失败
	ErrQueryUserLatestArchiveVersion     = New(2120, "query user latest archive version failed")   // 查询用户最新存档版本失败
	ErrFlowUnmarshal                     = New(2121, "unmarshal flow failed")                      // 反序列化 flow 失败
	ErrReleaseVm                         = New(2122, "release vm failed")                          // 请求 vm-manager 释放 vm 失败
	ErrReleaseUnion                      = New(2123, "release union failed")                       // 请求union-access释放区域失败
	ErrPrelaunch                         = New(2124, "prelaunch failed")                           // 提前分配游戏资源失败
	ErrCalcApInfo                        = New(2125, "calc ap info failed")                        // 计算资源出口线路失败
	ErrQueryVmInfo                       = New(2126, "query vm info failed")                       // 查询虚拟机资源信息失败
	ErrUpdateStreamInfo                  = New(2127, "update stream info failed")                  // 请求GSM更新stream info失败
	ErrQueryStreamStatus                 = New(2128, "query stream status failed")                 // 请求区域查询串流状态失败
	ErrUpdateStreamInfoNonRecoverable    = New(2129, "update stream info failed, non-recoverable") // 请求GSM更新stream info失败且无法恢复
	ErrNoAvailableResource               = New(2130, "no available resourced")                     // 无可用资源
	ErrLockResourceSpecImageNotMatchGame = New(2131, "lock resource spec image not match game")    // 指定规格镜像和游戏不匹配
	ErrLockResourceSpecNotExist          = New(2132, "lock resource spec not exist")               // 指定规格不存在
	ErrLockResourceSpecNotBelongToBiz    = New(2133, "lock resource spec not belong to biztype")   // 指定规格不属于该业务类型
	ErrPrejoinGame                       = New(2134, "prejoin game failed")                        // 参与者提前接入同实例失败

	// 22xx 用户相关错误
	ErrUserInGame         = New(2200, "user in game")                 // 正在游戏中
	ErrUserIsOffline      = New(2201, "user is offline")              // 用户不在线
	ErrUserNotInGame      = New(2202, "user not in game")             // 用户不在游戏中
	ErrQueryUserBasicInfo = New(2203, "query user basic info failed") // 查询用户信息失败
	ErrQueryUserThirdInfo = New(2204, "query user third info failed") // 查询用户第三方信息失败
	ErrQueryUserAreaInfo  = New(2205, "query user area info failed")  // 查询用户区域映射信息失败
	ErrExpireToken        = New(2206, "expire token failed")          // 过期 token 失败
	ErrLoadUserCache      = New(2207, "load user cache failed")       // 获取用户缓存失败
	ErrQueryUtoken        = New(2208, "query utoken failed")          // 获取 utoken 失败
	ErrAcquireUserDlock   = New(2209, "acquire user dlock failed")    // 获取用户分布式锁失败

	// 23xx storage-schedule & archiver-manager 相关错误
	// ...

	// 24xx union-storage-scheduler 相关错误
	// ...

	// 3xxx 组件错误
	// 30xx 网络相关错误(HTTP/WS 等)
	ErrHttpGet                 = New(3000, "http get failed")                 // GET 请求失败
	ErrHttpPost                = New(3001, "http post failed")                // POST 请求失败
	ErrHttpPut                 = New(3002, "http put failed")                 // PUT 请求失败
	ErrHttpDelete              = New(3003, "http delete failed")              // DELETE 请求失败
	ErrHttpParseUrl            = New(3004, "http parse url failed")           // 请求地址格式非法
	ErrHttpUnsupportMethod     = New(3005, "http unsupport method")           // 不支持的方法
	ErrHttpPutError            = New(3006, "http put error")                  // PUT 请求失败
	ErrHttpDeleteError         = New(3007, "http delete error")               // DELETE 请求失败
	ErrHttpError               = New(3008, "http error")                      // 请求返回非 2xx 状态码
	ErrHttpParseResponseBody   = New(3009, "http parse response body failed") // 响应的 body 解析为 json 失败
	ErrHttpResponseBodyEmpty   = New(3010, "http response body empty failed") // 响应的 body 为空
	ErrHttpResponseCodeNotZero = New(3011, "http response code not zero")     // 响应的 body 中 code 非 0
	ErrHttpAuth                = New(3012, "http auth failed")                // 接口鉴权失败
	ErrHttpParseAuthParam      = New(3013, "http parse auth param failed")    // 解析鉴权公共参数失败

	// 31xx DB 相关错误
	ErrDalCreate = New(3100, "dal create failed") // DAL 新建失败
	ErrDalUpdate = New(3101, "dal update failed") // DAL 修改失败
	ErrDalQuery  = New(3102, "dal query failed")  // DAL 查询失败
	ErrDalDelete = New(3103, "dal delete failed") // DAL 删除失败

	// 32xx ES 相关错误
	ErrEsQuery   = New(3200, "query es failed")   // 查询 ES 失败
	ErrEsTimeout = New(3200, "es timeout failed") // 请求超时

	// 33xx Redis 相关错误
	ErrRedisQuery         = New(3300, "redis query failed")          // reids查询失败
	ErrRedisWrite         = New(3301, "redis write failed")          // redis写入失败
	ErrRedisDelete        = New(3302, "redis delete failed")         // redis删除失败
	ErrRedisOperate       = New(3303, "redis operate failed")        // redis操作失败
	ErrRedisParseResponse = New(3304, "redis parse response failed") // redis解析返回数据失败

	// 34xx Zookeeper 相关错误
	ErrZkQueryNode            = New(3400, "zk query node failed")      // 查询节点失败
	ErrZkCreateNode           = New(3401, "zk create node failed")     // 创建节点失败
	ErrZkDeleteNode           = New(3402, "zk delete node failed")     // 删除节点失败
	ErrZkSetNode              = New(3403, "zk set node failed")        // 节点写入失败
	ErrZkNodeDeleted          = New(3404, "zk node deleted")           // 节点被删除
	ErrZkWatchUnexpectedEvent = New(3405, "zk watch unexpected event") // watch未预期的事件
	ErrZkUniqueIdNoResource   = New(3430, "zk unique id no resource")  // 申请区间唯一 ID：ID 资源已用尽
	ErrZkUniqueIdReleased     = New(3431, "zk unique id is released")  // 唯一 ID 已被释放
	ErrZkUniqueIdAlloc        = New(3432, "zk alloc unique id failed") // 获取唯一 ID 失败
	ErrZkDflagWild            = New(3440, "zk wild dflag")             // 分布式标记: 未获取的标记
	ErrZkDflagAcceptException = New(3441, "zk dflag accept exception") // 分布式标记: 收到异常通知
	ErrZkAcquireDlock         = New(3442, "zk acquire dlock failed")   // 分布式锁：获取锁失败

	// 35xx Kafka 相关错误
	ErrKafkaTopicSubscription = New(3500, "kafka topic subscription failed") // 订阅 topic 失败
	ErrKafkaProduceMsg        = New(3501, "kafka produce message failed")    // 发送消息失败

	// 4xxx 内部服务错误
	// 40xx vm-manager
	ErrAllocResource    = New(4000, "alloc resource failed")    // 申请分配资源失败
	ErrReleaseResource  = New(4001, "release resource failed")  // 申请释放资源失败
	ErrQueryRank        = New(4002, "query rank failed")        // 查询申请排名失败
	ErrExchangeResource = New(4003, "exchange resource failed") // 申请更换资源失败

	// 41xx 计费相关错误
	ErrAllocChargeId = New(4100, "alloc charge id failed") // 申请计费 ID 失败
	ErrStartCharge   = New(4101, "start charge failed")    // 开始计费失败
	ErrFinishCharge  = New(4102, "finish charge failed")   // 结束计费失败

	// 42xx iaas 相关错误
	ErrInstanceNotFound        = New(4200, "instance not found")           // 实例未发现
	ErrApplyInstances          = New(4201, "apply instances failed")       // 申请实例失败
	ErrQueryInstanceCount      = New(4202, "query instance count failed")  // 查询实例数量失败
	ErrQueryInstanceStatus     = New(4203, "query instance status failed") // 查询实例状态失败
	ErrInstanceStatusException = New(4204, "instance status exception")    // 实例状态异常
	ErrCreatePackage           = New(4220, "create package failed")        // 创建数据包失败
	ErrAreaPackageTask         = New(4221, "area package task failed")     // 区域数据包任务失败
	ErrQueryAreaConfig         = New(4222, "query area config failed")     // 查询区域配置失败
	ErrQueryAreaPackage        = New(4223, "query area package failed")    // 查询区域数据包失败
	ErrUpdatePackageStatus     = New(4224, "update package status failed") // 更新区域数据包状态失败
	ErrAreaPackageUnavailable  = New(4225, "package unavailable failed")   // 区域数据包不可用

	// 网关相关错误
	ErrGatewayResponseNoRetCode = New(-99999, "response body ret.code not found") // 返回数据中没有业务返回码
)

// Code 获取错误码
func Code(err error) int {
	if err == nil {
		return Succ.Code
	}
	if e, ok := err.(*ErrorWrap); ok {
		return e.Code
	}
	return Unknown.Code
}

// Msg 获取错误信息
func Msg(err error) string {
	if err == nil {
		return Succ.Message
	}
	if e, ok := err.(*ErrorWrap); ok {
		return e.Message
	}
	return err.Error()
}
