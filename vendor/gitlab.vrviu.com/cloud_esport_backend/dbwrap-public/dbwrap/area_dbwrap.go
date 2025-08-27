package dbwrap

import (
	"errors"
	"fmt"
	"reflect"
	"time"

	"github.com/levigross/grequests"
	"gitlab.vrviu.com/inviu/backend-go-public/gopublic"
	"gitlab.vrviu.com/inviu/backend-go-public/gopublic/httpwrap"
)

type AreaDBWrap struct {
	_table    string
	_host     string
	_callerID int
	_flowID   string
	_typ      reflect.Type
}

func CreateAreaDBWrap(host string, callerID int, table, flowID string, typ reflect.Type) *AreaDBWrap {
	return &AreaDBWrap{
		_table:    table,
		_host:     host,
		_callerID: callerID,
		_flowID:   flowID,
		_typ:      typ,
	}
}

func (iDBWrap *AreaDBWrap) areaTable(area int, table string) string {
	return fmt.Sprintf("%d/%s", area, table)
}

func (iDBWrap *AreaDBWrap) Query(area int, query string, sortby interface{}, ascending interface{}) (interface{}, int, error) {
	info := reflect.New(iDBWrap._typ).Elem()
	var errcode int
	err := httpwrap.SearchSingleObjectFromDAL(
		iDBWrap._host,
		iDBWrap.areaTable(area, iDBWrap._table),
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

func (iDBWrap *AreaDBWrap) QueryAll(area int, query string, sortby interface{}, ascending interface{}) (interface{}, int, error) {
	sliceTyp := reflect.SliceOf(iDBWrap._typ)
	infos := reflect.New(sliceTyp)
	infos.Elem().Set(reflect.MakeSlice(sliceTyp, 0, 0))

	var errcode int
	err := httpwrap.SearchAllObjectFromDAL(
		iDBWrap._host,
		iDBWrap.areaTable(area, iDBWrap._table),
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

func (iDBWrap *AreaDBWrap) QueryPage(area int, query string, offset, limit int, sortby, ascending interface{}) (int, interface{}, int, error) {
	sliceTyp := reflect.SliceOf(iDBWrap._typ)
	infos := reflect.New(sliceTyp)
	infos.Elem().Set(reflect.MakeSlice(sliceTyp, 0, 0))

	var errcode int
	total, err := httpwrap.SearchPageObjectFromDAL(
		iDBWrap._host,
		iDBWrap.areaTable(area, iDBWrap._table),
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

func (iDBWrap *AreaDBWrap) Insert(area int, info interface{}) (interface{}, int, error) {
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

func (iDBWrap *AreaDBWrap) Update(area int, key, info interface{}) (interface{}, int, error) {
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
		fmt.Sprintf("%s/v1/%s/%v", iDBWrap._host, iDBWrap.areaTable(area, iDBWrap._table), key),
		nil,
		rsp.Addr().Interface())

	if err != nil {
		errcode = DBErrCodeUpdateFailed
	}
	return rsp.Addr().Interface(), errcode, err
}

func (iDBWrap *AreaDBWrap) Delete(area int, key interface{}) (int, error) {
	grsp, err := grequests.Delete(
		fmt.Sprintf("%s/v1/%s/%v", iDBWrap._host, iDBWrap.areaTable(area, iDBWrap._table), key),
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
