package httpwrap

import (
	"context"
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"encoding/base64"
	"errors"
	"net/url"
	"sort"
	"strings"
	"time"

	vlog "gitlab.vrviu.com/inviu/backend-go-public/gopublic/vlog/tracing"
)

var (
	VrviuPermInfoCtxCacheKey = "vrviu-perm-info"

	SignatureVersionV1_0 = "1.0" // 签名算法版本1.0
	SignatureVersionV1_1 = "1.1" // 签名算法版本1.1: 在1.0的基础上包含url path参数
	SignatureVersionV1_2 = "1.2" // 签名算法版本1.2: 在1.1的基础上支持使用ticket代替AccessKeySecret
)

// AuthPublicParam 鉴权公共参数
type AuthPublicParam struct {
	AccessKey        string `json:"AccessKey" queryarg:"AccessKey,required"`
	Version          string `json:"Version" queryarg:"Version,required"`
	Timestamp        int64  `json:"Timestamp" queryarg:"Timestamp,required"`
	Format           string `json:"Format" queryarg:"Format,required"`
	SignatureVersion string `json:"SignatureVersion" queryarg:"SignatureVersion,required"`
	SignatureNonce   string `json:"SignatureNonce" queryarg:"SignatureNonce,required"`
	SignatureMethod  string `json:"SignatureMethod" queryarg:"SignatureMethod,required"`
	Signature        string `json:"Signature" queryarg:"Signature,required"`
	BizId            string `json:"BizId" queryarg:"BizId"` // 合作方子业务ID
	Ticket           string `json:"Ticket" queryarg:"Ticket"`
	TicketNonce      string `json:"TicketNonce" queryarg:"TicketNonce"`
}

// PermInfo 第三方鉴权key信息
type PermInfo struct {
	AppId           string    `json:"AppId"`           // 合作方ID
	BizId           string    `json:"BizId"`           // 合作方子业务ID
	AccessKeyId     string    `json:"AccessKeyId"`     // 业务密钥ID
	AccessKeySecret string    `json:"AccessKeySecret"` // 业务密钥
	Permission      int       `json:"Permission"`      // 权限
	BizType         int       `json:"BizType"`         // 业务类型
	IgnoreSignature int       `json:"IgnoreSignature"` // 忽略签名校验
	CheckExpireTime int       `json:"CheckExpireTime"` // 是否校验鉴权有效期
	ExpireTime      time.Time `json:"ExpireTime"`      // 鉴权有效期
}

// SecretHelper 获取秘钥
type SecretHelper interface {
	GetSecret(*AuthPublicParam) (*PermInfo, error)
}

// Authenticate 接口参数鉴权
// @param    method: 请求方法
// @param    secret: 鉴权秘钥
// @param signature: 请求携带的签名
// @param   urlPath: 请求的url路径（已经做了url_encode）
// @param  authargs: 鉴权公共参数
// @param queryargs: 请求查询参数
// @param      body: 请求body
func Authenticate(ctx context.Context, method, secret, signature string, urlPath string, queryargs map[string]string, body []byte) error {
	URLEncode := func(k string, v string) string {
		u := url.Values{}
		u.Set(k, v)
		return u.Encode()
	}

	var keys []string
	for k := range queryargs {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// 格式化查询字符串
	var pairs []string
	for _, key := range keys {
		value := queryargs[key]
		// golang中urlencode将空格escape为+，为了一鉴权服务保持一致，这里手动讲+号替换为%20
		pairs = append(pairs, strings.Replace(URLEncode(key, value), "+", "%20", -1))
	}
	canonicalizedQueryString := strings.Join(pairs, "&")

	// 计算BodyMd5
	hash := md5.New()
	hash.Write(body)
	bodyMd5 := base64.StdEncoding.EncodeToString(hash.Sum(nil))

	// 计算签名
	stringToSign := method + "&" + url.QueryEscape(urlPath) + "&" + url.QueryEscape(bodyMd5) + "&" + url.QueryEscape(canonicalizedQueryString)
	hmac := hmac.New(sha1.New, []byte(secret+"&"))
	hmac.Write([]byte(stringToSign))
	lsignature := base64.StdEncoding.EncodeToString(hmac.Sum(nil))

	vlog.Infof(ctx, "Authenticate(). AccessKeySecret(%s) CanonicalizedQueryString(%s) BodyMd5(%s) StringToSign(%s) Signature(%s) RequestSignature(%s)",
		secret, canonicalizedQueryString, bodyMd5, stringToSign, lsignature, signature)

	if lsignature != signature {
		return errors.New("authenticate failed")
	}

	return nil
}
