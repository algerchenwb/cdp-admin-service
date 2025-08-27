package helper

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/levigross/grequests"
	"github.com/zeromicro/go-zero/core/logx"
)

type HttpCommonResponse struct {
	Head struct {
		Code      int    `json:"code"`
		Msg       string `json:"msg"`
		RequestID string `json:"request_id"`
	} `json:"ret"`
	Body interface{} `json:"body"`
}

func parseHttpResponse(resp *grequests.Response, body interface{}) (int, error) {
	if int(resp.StatusCode/100) != 2 {
		return -3001, fmt.Errorf("non-2xx status code[%v]", resp.StatusCode)
	}

	if body != nil {
		if err := json.Unmarshal(resp.Bytes(), body); err != nil {
			return -3001, fmt.Errorf("json Unmarshal failure, err[%v]", err)
		}
	}

	return 0, nil
}

func HttpGet(ctx context.Context, sessionId string, url string, timeoutInMS int, rspBody interface{}) (errCode int, err error) {
	logx.WithContext(ctx).Debug(fmt.Sprintf("[%v] url[%v]", sessionId, url))

	// 请求参数
	ro := grequests.RequestOptions{
		RequestTimeout: time.Duration(timeoutInMS) * time.Millisecond,
		Headers: map[string]string{
			"X-System": "diskless-aggregator",
		},
	}

	// 发送请求
	grsp, e := grequests.Get(url, &ro)
	if e != nil {
		err = e
		errCode = -3001
		return
	}

	logx.WithContext(ctx).Debug(fmt.Sprintf("[%v] response[%+v]", sessionId, grsp))
	return parseHttpResponse(grsp, rspBody)
}

func HttpPut(ctx context.Context, sessionId string, url string, timeoutInMS int, reqBody interface{}, rspBody interface{}) (errCode int, err error) {
	logx.WithContext(ctx).Debug(fmt.Sprintf("[%v] url[%v], request[%+v]", sessionId, url, reqBody))
	// 请求参数
	ro := grequests.RequestOptions{
		RequestTimeout: time.Duration(timeoutInMS) * time.Millisecond,
		JSON:           reqBody,
	}

	// 发送请求
	grsp, e := grequests.Put(url, &ro)
	if e != nil {
		err = e
		errCode = -3001
		return
	}

	logx.WithContext(ctx).Debug(fmt.Sprintf("[%v] response[%+v]", sessionId, grsp))
	return parseHttpResponse(grsp, rspBody)
}

func HttpPost(ctx context.Context, sessionId string, url string, timeoutInMS int, reqBody interface{}, rspBody interface{}) (errCode int, err error) {
	logx.WithContext(ctx).Debug(fmt.Sprintf("[%v] url[%v], request[%+v]", sessionId, url, reqBody))
	// fmt.Printf("[%v] url[%v], request[%+v]\n", sessionId, url, reqBody)
	// 请求参数
	ro := grequests.RequestOptions{
		RequestTimeout: time.Duration(timeoutInMS) * time.Millisecond,
		JSON:           reqBody,
		Headers: map[string]string{
			"X-System": "diskless-aggregator",
		},
	}

	// 发送请求
	grsp, e := grequests.Post(url, &ro)
	if e != nil {
		err = e
		errCode = -3001
		return
	}

	logx.WithContext(ctx).Debugf("[%v] response[%+v]", sessionId, grsp)
	// fmt.Printf("[%v] response[%+v]\n", sessionId, grsp)
	return parseHttpResponse(grsp, rspBody)
}

func HttpPostWithHeaders(ctx context.Context, sessionId string, url string, timeoutInMS int, reqBody interface{}, rspBody interface{}, headers map[string]string) (errCode int, err error) {
	logx.WithContext(ctx).Debug(fmt.Sprintf("[%v] url[%v], request[%+v]", sessionId, url, reqBody))
	// fmt.Printf("[%v] url[%v], request[%+v]\n", sessionId, url, reqBody)
	// 请求参数
	ro := grequests.RequestOptions{
		RequestTimeout: time.Duration(timeoutInMS) * time.Millisecond,
		JSON:           reqBody,
		Headers:        headers,
	}

	// 发送请求
	grsp, e := grequests.Post(url, &ro)
	if e != nil {
		err = e
		errCode = -3001
		return
	}

	logx.WithContext(ctx).Debug(fmt.Sprintf("[%v] response[%+v]", sessionId, grsp))
	// fmt.Printf("[%v] response[%+v]\n", sessionId, grsp)
	return parseHttpResponse(grsp, rspBody)
}
