package helper

import (
	"cdp-admin-service/internal/model/errorx"
	"context"
	"fmt"
	"reflect"
	"strings"

	"github.com/zeromicro/go-zero/core/logx"
)

func CheckCommQueryParam(ctx context.Context, sessionId string, condList []string, orders string, sorts string, jsonObj interface{}) (qry string, err error) {
	for _, v := range condList {
		pos := strings.Index(v, ":")
		if pos == 0 {
			logx.WithContext(ctx).Errorf("[%s] unknown condition statement[%v]", sessionId, v)
			return "", errorx.NewDefaultCodeError(fmt.Sprintf("unknown condition statement[%v]", v))
		} else if pos == -1 {
			if !strings.Contains(v, "__isnull") && !strings.Contains(v, "__isnotnull") {
				logx.WithContext(ctx).Errorf("[%s] unknown condition statement[%v]  __isnull", sessionId, v)
				return "", errorx.NewDefaultCodeError(fmt.Sprintf("unknown condition statement[%v] __isnull", v))
			}
			pos = len(v)
		}
		if !CheckJsonTagExist(v[0:pos], jsonObj) {
			logx.WithContext(ctx).Errorf("[%s] unknown condition statement[%v]  CheckJsonTagExist", sessionId, v)
			return "", errorx.NewDefaultCodeError(fmt.Sprintf("unknown condition statement[%v] CheckJsonTagExist", v))
		} else {
			qry = qry + "$" + v
		}
	}
	if qry != "" {
		qry = qry[1:]
	}

	if sorts != "" && orders != "" {
		sortList := strings.Split(sorts, ",")
		for _, v := range sortList {
			if v == "" || !CheckJsonTagExist(v, jsonObj) {
				logx.WithContext(ctx).Errorf("[%s] bad param, sort[%v] ", sessionId, v)
				return "", errorx.NewDefaultCodeError(fmt.Sprintf("bad param, sort[%v]", v))
			}
		}
		orderList := strings.Split(orders, ",")
		for _, v := range orderList {
			if v == "asc" || v == "desc" {
				continue
			} else {
				logx.WithContext(ctx).Errorf("[%s] bad param, sort[%v] ", sessionId, v)
				return "", errorx.NewDefaultCodeError(fmt.Sprintf("bad param, sort[%v]", v))
			}
		}

		if len(sortList) != len(orderList) {
			logx.WithContext(ctx).Errorf("[%s] bad param, len(sorts) != len(orders)", sessionId)
			return "", errorx.NewDefaultCodeError("bad param, len(sorts) != len(orders)")
		}
	} else if sorts != "" || orders != "" {
		logx.WithContext(ctx).Errorf("[%s] bad param, len(sorts) != len(orders)", sessionId)
		return "", errorx.NewDefaultCodeError("bad param, len(sorts) != len(orders)")
	}

	return qry, nil
}

func CheckJsonTagExist(jsonTag string, val interface{}) bool {
	if strings.Contains(jsonTag, "___") {
		return false
	}
	if pos := strings.Index(jsonTag, "__"); pos > 0 {
		jsonTag = jsonTag[0:pos]
	}
	item := reflect.ValueOf(val)
	for j := 0; j < item.Type().NumField(); j++ {
		tagName := item.Type().Field(j).Tag.Get("json")
		if tagName == jsonTag {
			return true
		}
	}
	return false
}
