package dbwrap

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"time"

	"github.com/levigross/grequests"
	"gitlab.vrviu.com/inviu/backend-go-public/gopublic"
	"gitlab.vrviu.com/inviu/backend-go-public/gopublic/httpwrap"
)

type DBWrap struct {
	_table    string
	_host     string
	_callerID int
	_flowID   string
	_typ      reflect.Type
}

const (
	_Base64EncodeStd string = "3WiF4_yMOEqB-mdCjYTPG9sz2avUf6H8cbpI1N75xglDtwZXARJoknuerVhLQ0KS"
)

var _base64Encoding *base64.Encoding

func init() {
	_base64Encoding = base64.NewEncoding(_Base64EncodeStd)
}

func CreateDBWrap(host string, callerID int, table, flowID string, typ reflect.Type) *DBWrap {
	return &DBWrap{
		_table:    table,
		_host:     host,
		_callerID: callerID,
		_flowID:   flowID,
		_typ:      typ,
	}
}

func (iDBWrap *DBWrap) Query(query string, sortby interface{}, ascending interface{}) (interface{}, int, error) {
	info := reflect.New(iDBWrap._typ).Elem()
	var errcode int
	err := httpwrap.SearchSingleObjectFromDAL(
		iDBWrap._host,
		iDBWrap._table,
		query,
		sortby,
		ascending,
		info.Addr().Interface(),
		httpwrap.MCParam{
			FlowID:      iDBWrap._flowID,
			SessionID:   gopublic.GenerateRandonString(12),
			CallerID:    iDBWrap._callerID,
			CalleeID:    MDAL,
			InterfaceID: IMCDAL_Query,
		},
	)

	if err != nil {
		errcode = DBErrCodeQueryFailed
		if err != gopublic.ErrNotExist {
			err = fmt.Errorf(ErrStrQueryFailed, iDBWrap._table)
		}
	}

	return info.Addr().Interface(), errcode, err
}

func (iDBWrap *DBWrap) QueryAll(query string, sortby interface{}, ascending interface{}) (interface{}, int, error) {
	sliceTyp := reflect.SliceOf(iDBWrap._typ)
	infos := reflect.New(sliceTyp)
	infos.Elem().Set(reflect.MakeSlice(sliceTyp, 0, 0))

	var errcode int
	err := httpwrap.SearchAllObjectFromDAL(
		iDBWrap._host,
		iDBWrap._table,
		query,
		sortby,
		ascending,
		infos.Interface(),
		httpwrap.MCParam{
			FlowID:      iDBWrap._flowID,
			SessionID:   gopublic.GenerateRandonString(12),
			CallerID:    iDBWrap._callerID,
			CalleeID:    MDAL,
			InterfaceID: IMCDAL_Query,
		},
	)

	if err != nil {
		errcode = DBErrCodeQueryFailed
		if err != gopublic.ErrNotExist {
			err = fmt.Errorf(ErrStrQueryFailed, iDBWrap._table)
		}
	}

	return infos.Elem().Interface(), errcode, err
}

func (iDBWrap *DBWrap) QueryPage(query string, offset, limit int, sortby, ascending interface{}) (int, interface{}, int, error) {
	sliceTyp := reflect.SliceOf(iDBWrap._typ)
	infos := reflect.New(sliceTyp)
	infos.Elem().Set(reflect.MakeSlice(sliceTyp, 0, 0))

	var errcode int
	total, err := httpwrap.SearchPageObjectFromDAL(
		iDBWrap._host,
		iDBWrap._table,
		query,
		offset,
		limit,
		sortby,
		ascending,
		infos.Interface(),
		httpwrap.MCParam{
			FlowID:      iDBWrap._flowID,
			SessionID:   gopublic.GenerateRandonString(12),
			CallerID:    iDBWrap._callerID,
			CalleeID:    MDAL,
			InterfaceID: IMCDAL_Query,
		},
	)

	if err != nil {
		errcode = DBErrCodeQueryFailed
		if err != gopublic.ErrNotExist {
			err = fmt.Errorf(ErrStrQueryFailed, iDBWrap._table)
		}
	}

	return total, infos.Elem().Interface(), errcode, err
}

func (iDBWrap *DBWrap) Insert(info interface{}) (interface{}, int, error) {
	rsp := reflect.New(iDBWrap._typ).Elem()
	errcode, err := httpwrap.HTTPRetry(
		httpwrap.DefaultRC,
		1,
		&httpwrap.RequestOptions{
			RequestOptions: grequests.RequestOptions{
				RequestTimeout: 3000 * time.Millisecond,
				JSON:           info,
			},
			MCParam: httpwrap.MCParam{
				FlowID:      iDBWrap._flowID,
				SessionID:   gopublic.GenerateRandonString(12),
				CallerID:    iDBWrap._callerID,
				CalleeID:    MDAL,
				InterfaceID: IMCDAL_Insert,
			},
		},
		httpwrap.HTTPPost,
		fmt.Sprintf("%s/v1/%s", iDBWrap._host, iDBWrap._table),
		nil,
		rsp.Addr().Interface())

	if err != nil {
		errcode = DBErrCodeInsertFailed
	}
	return rsp.Addr().Interface(), errcode, err
}

func (iDBWrap *DBWrap) Update(key, info interface{}) (interface{}, int, error) {
	rsp := reflect.New(iDBWrap._typ).Elem()
	errcode, err := httpwrap.HTTPRetry(
		httpwrap.DefaultRC,
		1,
		&httpwrap.RequestOptions{
			RequestOptions: grequests.RequestOptions{
				RequestTimeout: 3000 * time.Millisecond,
				JSON:           info,
			},
			MCParam: httpwrap.MCParam{
				FlowID:      iDBWrap._flowID,
				SessionID:   gopublic.GenerateRandonString(12),
				CallerID:    iDBWrap._callerID,
				CalleeID:    MDAL,
				InterfaceID: IMCDAL_Update,
			},
		},
		httpwrap.HTTPPut,
		fmt.Sprintf("%s/v1/%s/%v", iDBWrap._host, iDBWrap._table, key),
		nil,
		rsp.Addr().Interface())

	if err != nil {
		errcode = DBErrCodeUpdateFailed
	}
	return rsp.Addr().Interface(), errcode, err
}

func (iDBWrap *DBWrap) Delete(key interface{}) (int, error) {
	grsp, err := grequests.Delete(
		fmt.Sprintf("%s/v1/%s/%v", iDBWrap._host, iDBWrap._table, key),
		&grequests.RequestOptions{
			RequestTimeout: 3000 * time.Millisecond,
		},
	)
	if err != nil {
		return DBErrCodeDeleteFailed, err
	}

	var rsp httpwrap.HTTPResponse
	err = grsp.JSON(&rsp)
	if err != nil {
		return DBErrCodeDeleteFailed, gopublic.ErrUnmarshal
	}

	// 不存在认为成功
	if rsp.Head.Code == -3 && rsp.Head.Msg == "not exist" {
		return 0, nil
	} else if rsp.Head.Code != 0 {
		return DBErrCodeDeleteFailed, errors.New(rsp.Head.Msg)
	}

	return 0, nil
}

func (iDBWrap *DBWrap) RawQuery(sql string, rettype reflect.Type) (interface{}, int, error) {
	var head httpwrap.HTTPCommonHead
	var body httpwrap.DALResponsePageBody

	errcode, err := httpwrap.HTTPRetry(
		httpwrap.DefaultRC,
		1,
		&httpwrap.RequestOptions{
			RequestOptions: grequests.RequestOptions{
				RequestTimeout: 3 * time.Second,
				Params: map[string]string{
					"sql": _base64Encoding.EncodeToString([]byte(sql)),
				},
			},
			MCParam: httpwrap.MCParam{
				FlowID:      iDBWrap._flowID,
				SessionID:   gopublic.GenerateRandonString(12),
				CallerID:    iDBWrap._callerID,
				CalleeID:    MDAL,
				InterfaceID: IMCDAL_Query,
			},
		},
		httpwrap.HTTPGet,
		fmt.Sprintf("%s/v2/raw_query", iDBWrap._host),
		&head,
		&body,
	)

	if err != nil {
		if head.Code == -3 && head.Msg == "not exist" {
			return nil, head.Code, gopublic.ErrNotExist
		} else {
			return nil, DBErrCodeRawQueryFailed, fmt.Errorf(ErrStrRawQueryFailed)
		}
	}

	sliceTyp := reflect.SliceOf(rettype)
	infos := reflect.New(sliceTyp)
	infos.Elem().Set(reflect.MakeSlice(sliceTyp, 0, 0))

	for _, v := range body.List {
		obj := reflect.New(rettype)
		if err := json.Unmarshal([]byte(gopublic.ToJSON(v)), obj.Interface()); err != nil {
			continue
		} else {
			infos.Elem().Set(reflect.Append(infos.Elem(), obj.Elem()))
		}
	}

	return infos.Elem().Interface(), errcode, err
}

func (iDBWrap *DBWrap) AreaQuery(areaType int, query string, sortby interface{}, ascending interface{}) (interface{}, int, error) {
	info := reflect.New(iDBWrap._typ).Elem()
	var errcode int
	err := httpwrap.SearchSingleObjectFromDAL(
		iDBWrap._host,
		strconv.FormatInt(int64(areaType), 10)+"/"+iDBWrap._table,
		query,
		sortby,
		ascending,
		info.Addr().Interface(),
		httpwrap.MCParam{
			FlowID:      iDBWrap._flowID,
			SessionID:   gopublic.GenerateRandonString(12),
			CallerID:    iDBWrap._callerID,
			CalleeID:    MDAL,
			InterfaceID: IMCDAL_Query,
		},
	)

	if err != nil {
		errcode = DBErrCodeQueryFailed
		if err != gopublic.ErrNotExist {
			err = fmt.Errorf(ErrStrQueryFailed, iDBWrap._table)
		}
	}

	return info.Addr().Interface(), errcode, err
}

func (iDBWrap *DBWrap) AreaQueryAll(areaType int, query string, sortby interface{}, ascending interface{}) (interface{}, int, error) {
	sliceTyp := reflect.SliceOf(iDBWrap._typ)
	infos := reflect.New(sliceTyp)
	infos.Elem().Set(reflect.MakeSlice(sliceTyp, 0, 0))

	var errcode int
	err := httpwrap.SearchAllObjectFromDAL(
		iDBWrap._host,
		strconv.FormatInt(int64(areaType), 10)+"/"+iDBWrap._table,
		query,
		sortby,
		ascending,
		infos.Interface(),
		httpwrap.MCParam{
			FlowID:      iDBWrap._flowID,
			SessionID:   gopublic.GenerateRandonString(12),
			CallerID:    iDBWrap._callerID,
			CalleeID:    MDAL,
			InterfaceID: IMCDAL_Query,
		},
	)

	if err != nil {
		errcode = DBErrCodeQueryFailed
		if err != gopublic.ErrNotExist {
			err = fmt.Errorf(ErrStrQueryFailed, iDBWrap._table)
		}
	}

	return infos.Elem().Interface(), errcode, err
}

func (iDBWrap *DBWrap) AreaQueryPage(areaType int, query string, offset, limit int, sortby, ascending interface{}) (int, interface{}, int, error) {
	sliceTyp := reflect.SliceOf(iDBWrap._typ)
	infos := reflect.New(sliceTyp)
	infos.Elem().Set(reflect.MakeSlice(sliceTyp, 0, 0))

	var errcode int
	total, err := httpwrap.SearchPageObjectFromDAL(
		iDBWrap._host,
		strconv.FormatInt(int64(areaType), 10)+"/"+iDBWrap._table,
		query,
		offset,
		limit,
		sortby,
		ascending,
		infos.Interface(),
		httpwrap.MCParam{
			FlowID:      iDBWrap._flowID,
			SessionID:   gopublic.GenerateRandonString(12),
			CallerID:    iDBWrap._callerID,
			CalleeID:    MDAL,
			InterfaceID: IMCDAL_Query,
		},
	)

	if err != nil {
		errcode = DBErrCodeQueryFailed
		if err != gopublic.ErrNotExist {
			err = fmt.Errorf(ErrStrQueryFailed, iDBWrap._table)
		}
	}

	return total, infos.Elem().Interface(), errcode, err
}

func (iDBWrap *DBWrap) AreaInsert(areaType int, info interface{}) (interface{}, int, error) {
	rsp := reflect.New(iDBWrap._typ).Elem()
	errcode, err := httpwrap.HTTPRetry(
		httpwrap.DefaultRC,
		1,
		&httpwrap.RequestOptions{
			RequestOptions: grequests.RequestOptions{
				RequestTimeout: 3000 * time.Millisecond,
				JSON:           info,
			},
			MCParam: httpwrap.MCParam{
				FlowID:      iDBWrap._flowID,
				SessionID:   gopublic.GenerateRandonString(12),
				CallerID:    iDBWrap._callerID,
				CalleeID:    MDAL,
				InterfaceID: IMCDAL_Insert,
			},
		},
		httpwrap.HTTPPost,
		fmt.Sprintf("%s/v1/%d/%s", iDBWrap._host, areaType, iDBWrap._table),
		nil,
		rsp.Addr().Interface())

	if err != nil {
		errcode = DBErrCodeInsertFailed
	}
	return rsp.Addr().Interface(), errcode, err
}

func (iDBWrap *DBWrap) AreaUpdate(areaType int, key, info interface{}) (interface{}, int, error) {
	rsp := reflect.New(iDBWrap._typ).Elem()
	errcode, err := httpwrap.HTTPRetry(
		httpwrap.DefaultRC,
		1,
		&httpwrap.RequestOptions{
			RequestOptions: grequests.RequestOptions{
				RequestTimeout: 3000 * time.Millisecond,
				JSON:           info,
			},
			MCParam: httpwrap.MCParam{
				FlowID:      iDBWrap._flowID,
				SessionID:   gopublic.GenerateRandonString(12),
				CallerID:    iDBWrap._callerID,
				CalleeID:    MDAL,
				InterfaceID: IMCDAL_Update,
			},
		},
		httpwrap.HTTPPut,
		fmt.Sprintf("%s/v1/%d/%s/%v", iDBWrap._host, areaType, iDBWrap._table, key),
		nil,
		rsp.Addr().Interface())

	if err != nil {
		errcode = DBErrCodeUpdateFailed
	}
	return rsp.Addr().Interface(), errcode, err
}

func (iDBWrap *DBWrap) AreaDelete(areaType int, key interface{}) (int, error) {
	grsp, err := grequests.Delete(
		fmt.Sprintf("%s/v1/%d/%s/%v", iDBWrap._host, areaType, iDBWrap._table, key),
		&grequests.RequestOptions{
			RequestTimeout: 3000 * time.Millisecond,
		},
	)
	if err != nil {
		return DBErrCodeDeleteFailed, err
	}

	var rsp httpwrap.HTTPResponse
	err = grsp.JSON(&rsp)
	if err != nil {
		return DBErrCodeDeleteFailed, gopublic.ErrUnmarshal
	}

	// 不存在认为成功
	if rsp.Head.Code == -3 && rsp.Head.Msg == "not exist" {
		return 0, nil
	} else if rsp.Head.Code != 0 {
		return DBErrCodeDeleteFailed, errors.New(rsp.Head.Msg)
	}

	return 0, nil
}
