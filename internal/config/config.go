package config

import (
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	rest.RestConf
	Salt    string
	JwtAuth struct {
		AccessSecret string
		AccessExpire int64
	}

	OutSide struct {
		SaasHost               string
		DisklessHost           string
		DisklessCloudImageHost string
		VmMgrHost              string
		DisklessEdge1          string
		Timeout                int64
	}

	EsportDal struct {
		EsportDalHost    string
		EsportDalTimeout int64
	}

	UserRpc        zrpc.RpcClientConf
	AES            AESConfig
	Saas           SaasConfig
	DisklessConfig DisklessConfig
	Account        AccountConfig
	Schema         SchemaConfig
	Instance       InstanceConfig
	Menu           MenuConfig
}

type MenuConfig struct {
	SuanliRootId  int64
	ShigongRootId int64
}

type AccountConfig struct {
	DefaultPassword string
}

type RpcConfig struct {
	EtcdHosts []string
	Key       string
	EndPoints []string
}

type SpecificationConfig struct {
	DefaultTotalInstances           int64
	DefaultTimeResourcePoolNum      int64
	DefaultFrequencyResourcePoolNum int64
	DefaultValidityPeriodYears      int
}

type AESConfig struct {
	Key  string
	Salt string
}

type SaasConfig struct {
	ServerHost    string
	Timeout       int
	SpecIdList    []int64
	Specification SpecificationConfig
}

type DisklessConfig struct {
	AggregatorHost string
	Timeout        int
}

type SchemaConfig struct {
	DefaultSchemaConfig      string
	DefaultResetSchemaConfig string
	DefaultOsMntType         int
	DefaultStorageType       int
	DefaultWrbackType        int
	DefaultBootPnpPath       string
}

type InstanceConfig struct {
	InvaildVersion int64
	DefaultPoolId  int64
}
