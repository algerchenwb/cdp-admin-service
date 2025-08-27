package httpwrap

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"time"

	"github.com/astaxie/beego"
	"github.com/levigross/grequests"
	"gitlab.vrviu.com/inviu/backend-go-public/gopublic"
	"gitlab.vrviu.com/inviu/backend-go-public/gopublic/otelwrap"
	vlog "gitlab.vrviu.com/inviu/backend-go-public/gopublic/vlog/tracing"
)

// DALResponsePageBody DAL分页响应的Body结构
type DALResponsePageBody struct {
	Total int           `json:"total"`
	List  []interface{} `json:"list"`
}

func GetDalTimeoutDuration() (rst time.Duration) {
	return time.Duration(beego.AppConfig.DefaultInt("dal_timeout_ms", 3000)) * time.Millisecond
}

// QueryObjectFromDAL 请求DAL查询信息
// @param address: dal地址，格式为http://xxx.xxx.xxx.xxx:xxxxx
// @param   table: 查询表名
// @param     key: 查询主键值
// @pram      obj: 存放返回数据的指针
// @param mcparam: 模调参数
func QueryObjectFromDAL(address, table string, key interface{}, obj interface{}, mcparam ...MCParam) error {
	return QueryObjectFromDALWithCtx(otelwrap.NewSkipTraceCtx("QueryObjectFromDAL"), address, table, key, obj, mcparam...)
}

func QueryObjectFromDALWithCtx(ctx context.Context, address, table string, key interface{}, obj interface{}, mcparam ...MCParam) error {
	if obj == nil {
		return nil
	}

	if reflect.ValueOf(obj).Kind() != reflect.Ptr {
		return errors.New("obj is not pointer")
	}

	var head HTTPCommonHead

	_, err := HTTPRetry(
		DefaultRC,
		1,
		&RequestOptions{
			RequestName: fmt.Sprintf("%s(DALQuery)", table),
			RequestOptions: grequests.RequestOptions{
				RequestTimeout: GetDalTimeoutDuration(),
			},
			MCParam: func() MCParam {
				if len(mcparam) > 0 {
					mcparam[0].Ctx = ctx
					return mcparam[0]
				}
				return MCParam{Ctx: ctx}
			}(),
		},
		HTTPGet,
		fmt.Sprintf("%s/v1/%s/%v", address, table, key),
		&head,
		obj)

	if head.Code == -3 && head.Msg == "not exist" {
		return gopublic.ErrNotExist
	}

	return err
}

// SearchSingleObjectFromDAL 请求DAL条件查询信息
// @param   address: dal地址，格式为http://xxx.xxx.xxx.xxx:xxxxx
// @param     table: 查询表名
// @param     query: 查询条件字符串
// @param    sortby: 排序字段：nil表示不排序；多个字段排序以逗号隔开，若field1,field2
// @param ascending: 排序方向：true-升序，false-降序（默认）；多个字段排序时，按sortby指定的字段顺序传入由多个desc或asc拼接而成字符串，以逗号隔开，如desc,asc
// @param       obj: 存放返回数据的指针
// @param
func SearchSingleObjectFromDAL(address, table, query string, sortby interface{}, ascending interface{}, obj interface{}, mcparam ...MCParam) error {
	return SearchSingleObjectFromDALWithCtx(otelwrap.NewSkipTraceCtx("SearchSingleObjectFromDAL"), address, table, query, sortby, ascending, obj, mcparam...)
}

func SearchSingleObjectFromDALWithCtx(ctx context.Context, address, table, query string, sortby interface{}, ascending interface{}, obj interface{}, mcparam ...MCParam) error {
	var head HTTPCommonHead
	var body DALResponsePageBody

	if _, err := HTTPRetry(
		DefaultRC,
		1,
		&RequestOptions{
			RequestName: fmt.Sprintf("%s(DALSearchSingle)", table),
			RequestOptions: grequests.RequestOptions{
				RequestTimeout: GetDalTimeoutDuration(),
				Params: func() map[string]string {
					m := map[string]string{
						"query":  query,
						"offset": "0",
						"limit":  "1",
					}

					if sortby != nil {
						m["sortby"] = sortby.(string)
					}

					if ascending != nil {
						switch ascending := ascending.(type) {
						case bool:
							if ascending {
								m["order"] = "asc"
							} else {
								m["order"] = "desc"
							}
						case string:
							m["order"] = ascending
						}
					}

					return m
				}(),
			},
			MCParam: func() MCParam {
				if len(mcparam) > 0 {
					mcparam[0].Ctx = ctx
					return mcparam[0]
				}
				return MCParam{Ctx: ctx}
			}(),
		},
		HTTPGet,
		fmt.Sprintf("%s/v1/%s", address, table),
		&head,
		&body); err != nil && (head.Code != -3 || head.Msg != "not exist") {
		return err
	} else if head.Code == -3 && head.Msg == "not exist" {
		return gopublic.ErrNotExist
	}

	if obj == nil {
		return nil
	}

	if reflect.ValueOf(obj).Kind() != reflect.Ptr {
		return errors.New("obj is not pointer")
	}

	json.Unmarshal([]byte(gopublic.ToJSON(body.List[0])), obj)
	return nil
}

// SearchPageObjectFromDAL 请求DAL分页查询信息
// @param   address: dal地址，格式为http://xxx.xxx.xxx.xxx:xxxxx
// @param     table: 查询表名
// @param     query: 查询条件字符串
// @param    offset: 分页偏移
// @param     limit: 返回数量
// @param    sortby: 排序字段：nil表示不排序；多个字段排序以逗号隔开，若field1,field2
// @param ascending: 排序方向：true-升序，false-降序（默认）；多个字段排序时，按sortby指定的字段顺序传入由多个desc或asc拼接而成字符串，以逗号隔开，如desc,asc
// @param      objs: 存放返回数据的指针，指向一个数组
// @ -
// @return total: 总数
func SearchPageObjectFromDAL(address, table, query string, offset, limit int, sortby interface{}, ascending interface{}, objs interface{}, mcparam ...MCParam) (total int, err error) {
	return SearchPageObjectFromDALWithCtx(otelwrap.NewSkipTraceCtx("SearchPageObjectFromDAL"), address, table, query, offset, limit, sortby, ascending, objs, mcparam...)
}

func SearchPageObjectFromDALWithCtx(ctx context.Context, address, table, query string, offset, limit int, sortby interface{}, ascending interface{}, objs interface{}, mcparam ...MCParam) (total int, err error) {
	var head HTTPCommonHead
	var body DALResponsePageBody

	if _, err = HTTPRetry(
		DefaultRC,
		1,
		&RequestOptions{
			RequestName: fmt.Sprintf("%s(DALSearchPage)", table),
			RequestOptions: grequests.RequestOptions{
				RequestTimeout: GetDalTimeoutDuration(),
				Params: func() map[string]string {
					m := map[string]string{
						"query":  query,
						"offset": strconv.Itoa(offset),
						"limit":  strconv.Itoa(limit),
					}

					if sortby != nil {
						m["sortby"] = sortby.(string)
					}

					if ascending != nil {
						switch ascending := ascending.(type) {
						case bool:
							if ascending {
								m["order"] = "asc"
							} else {
								m["order"] = "desc"
							}
						case string:
							m["order"] = ascending
						}
					}

					return m
				}(),
			},
			MCParam: func() MCParam {
				if len(mcparam) > 0 {
					mcparam[0].Ctx = ctx
					return mcparam[0]
				}
				return MCParam{Ctx: ctx}
			}(),
		},
		HTTPGet,
		fmt.Sprintf("%s/v1/%s", address, table),
		&head,
		&body); err != nil && (head.Code != -3 || head.Msg != "not exist") {
		return 0, err
	} else if head.Code == -3 && head.Msg == "not exist" {
		return 0, nil
	}

	total = body.Total

	if objs == nil {
		return
	}

	if reflect.ValueOf(objs).Kind() != reflect.Ptr {
		return 0, errors.New("objs is not pointer")
	}

	objArrVal := reflect.ValueOf(objs).Elem()
	ObjType := reflect.TypeOf(objs).Elem().Elem()

	for _, v := range body.List {
		obj := reflect.New(ObjType)
		if err := json.Unmarshal([]byte(gopublic.ToJSON(v)), obj.Interface()); err != nil {
			continue
		} else {
			objArrVal.Set(reflect.Append(objArrVal, obj.Elem()))
		}
	}

	return
}

// SearchAllObjectFromDAL 请求DAL查询全部信息
// @param   address: dal地址，格式为http://xxx.xxx.xxx.xxx:xxxxx
// @param     table: 查询表名
// @param     query: 查询条件字符串
// @param    sortby: 排序字段：nil表示不排序；多个字段排序以逗号隔开，若field1,field2
// @param ascending: 排序方向：true-升序，false-降序（默认）；多个字段排序时，按sortby指定的字段顺序传入由多个desc或asc拼接而成字符串，以逗号隔开，如desc,asc
// @param      objs: 存放返回数据的指针，指向一个数组
func SearchAllObjectFromDAL(address, table, query string, sortby interface{}, ascending interface{}, objs interface{}, mcparam ...MCParam) error {
	return SearchAllObjectFromDALWithCtx(otelwrap.NewSkipTraceCtx("SearchAllObjectFromDAL"), address, table, query, sortby, ascending, objs, mcparam...)
}

func SearchAllObjectFromDALWithCtx(ctx context.Context, address, table, query string, sortby interface{}, ascending interface{}, objs interface{}, mcparam ...MCParam) error {
	objArrVal := reflect.ValueOf(objs).Elem()
	ObjType := reflect.TypeOf(objs).Elem().Elem()

	offset := 0
	limit := 500
	total := 500

	for offset < total {
		var head HTTPCommonHead
		var body DALResponsePageBody

		if _, err := HTTPRetry(
			DefaultRC,
			1,
			&RequestOptions{
				RequestName: fmt.Sprintf("%s(DALSearchAll)", table),
				RequestOptions: grequests.RequestOptions{
					RequestTimeout: GetDalTimeoutDuration(),
					Params: func() map[string]string {
						m := map[string]string{
							"query":  query,
							"offset": strconv.Itoa(offset),
							"limit":  strconv.Itoa(limit),
						}

						if sortby != nil {
							m["sortby"] = sortby.(string)
						}

						if ascending != nil {
							switch ascending := ascending.(type) {
							case bool:
								if ascending {
									m["order"] = "asc"
								} else {
									m["order"] = "desc"
								}
							case string:
								m["order"] = ascending
							}
						}

						return m
					}(),
				},
				MCParam: func() MCParam {
					if len(mcparam) > 0 {
						mcparam[0].Ctx = ctx
						return mcparam[0]
					}
					return MCParam{Ctx: ctx}
				}(),
			},
			HTTPGet,
			fmt.Sprintf("%s/v1/%s", address, table),
			&head,
			&body); err != nil && (head.Code != -3 || head.Msg != "not exist") {
			return err
		} else if head.Code == -3 && head.Msg == "not exist" {
			return nil
		}

		if objs == nil {
			return nil
		}

		if reflect.ValueOf(objs).Kind() != reflect.Ptr {
			return errors.New("objs is not pointer")
		}

		for _, v := range body.List {
			obj := reflect.New(ObjType)
			if err := json.Unmarshal([]byte(gopublic.ToJSON(v)), obj.Interface()); err != nil {
				continue
			} else {
				objArrVal.Set(reflect.Append(objArrVal, obj.Elem()))
			}
		}

		offset += len(body.List)
		total = body.Total
	}

	return nil
}

// InsertObjectToDAL 请求DAL插入记录
// @param   address: dal地址，格式为http://xxx.xxx.xxx.xxx:xxxxx
// @param     table: 插入表名
// @param       obj: 插入数据
func InsertObjectToDAL(address, table string, obj interface{}, mcparam ...MCParam) error {
	return InsertObjectToDALWithCtx(otelwrap.NewSkipTraceCtx("InsertObjectToDAL"), address, table, obj, mcparam...)
}

func InsertObjectToDALWithCtx(ctx context.Context, address, table string, obj interface{}, mcparam ...MCParam) error {
	_, err := HTTPRetry(
		DefaultRC,
		1,
		&RequestOptions{
			RequestName: fmt.Sprintf("%s(DALInsert)", table),
			RequestOptions: grequests.RequestOptions{
				RequestTimeout: GetDalTimeoutDuration(),
				JSON:           obj,
			},
			MCParam: func() MCParam {
				if len(mcparam) > 0 {
					mcparam[0].Ctx = ctx
					return mcparam[0]
				}
				return MCParam{Ctx: ctx}
			}(),
		},
		HTTPPost,
		fmt.Sprintf("%s/v1/%s", address, table),
		nil,
		nil,
	)
	return err
}

// UpdateObjectToDAL 请求DAL更新记录
// @param   address: dal地址，格式为http://xxx.xxx.xxx.xxx:xxxxx
// @param     table: 更新表名
// @param       key: 更新数据key
// @param       obj: 更新数据
func UpdateObjectToDAL(address, table string, key, obj interface{}, mcparam ...MCParam) error {
	return UpdateObjectToDALWithCtx(otelwrap.NewSkipTraceCtx("UpdateObjectToDAL"), address, table, key, obj, mcparam...)
}

func UpdateObjectToDALWithCtx(ctx context.Context, address, table string, key, obj interface{}, mcparam ...MCParam) error {
	_, err := HTTPRetry(
		DefaultRC,
		1,
		&RequestOptions{
			RequestName: fmt.Sprintf("%s(DALUpdate)", table),
			RequestOptions: grequests.RequestOptions{
				RequestTimeout: GetDalTimeoutDuration(),
				JSON:           obj,
			},
			MCParam: func() MCParam {
				if len(mcparam) > 0 {
					mcparam[0].Ctx = ctx
					return mcparam[0]
				}
				return MCParam{Ctx: ctx}
			}(),
		},
		HTTPPut,
		fmt.Sprintf("%s/v1/%s/%v", address, table, key),
		nil,
		nil,
	)
	return err
}

// DeleteObjectToDAL 请求DAL更新记录
// @param   address: dal地址，格式为http://xxx.xxx.xxx.xxx:xxxxx
// @param     table: 表名
// @param        id: 数据主键id
func DeleteObjectToDAL(address, table string, id interface{}, mcparam ...MCParam) error {
	return DeleteObjectToDALWithCtx(otelwrap.NewSkipTraceCtx("DeleteObjectToDAL"), address, table, id, mcparam...)
}

func DeleteObjectToDALWithCtx(ctx context.Context, address, table string, id interface{}, mcparam ...MCParam) error {
	var head HTTPCommonHead
	_, err := HTTPRetry(
		DefaultRC,
		1,
		&RequestOptions{
			RequestName: fmt.Sprintf("%s(DALDelete)", table),
			RequestOptions: grequests.RequestOptions{
				RequestTimeout: GetDalTimeoutDuration(),
			},
			MCParam: func() MCParam {
				if len(mcparam) > 0 {
					mcparam[0].Ctx = ctx
					return mcparam[0]
				}
				return MCParam{Ctx: ctx}
			}(),
		},
		HttpDelete,
		fmt.Sprintf("%s/v1/%s/%v", address, table, id),
		&head,
		nil,
	)
	// 不存在认为成功
	if head.Code == -3 {
		vlog.Warnf(ctx, "DeleteObjectToDAL(). not exist. table(%s) id(%v)", table, id)
		return nil
	}

	return err
}
