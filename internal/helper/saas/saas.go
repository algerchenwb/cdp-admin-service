package saas

import "sync"

type Request struct {
	Header string `header:"X-System"`
}

type SaasCommonResp struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}
type SaasServerService struct {
	Host    string
	Timeout int
}

var instance *SaasServerService
var once sync.Once

func NewSaasServerService(host string, timeout int) *SaasServerService {
	once.Do(func() {
		instance = &SaasServerService{
			Host:    host,
			Timeout: timeout,
		}
	})
	return instance
}
