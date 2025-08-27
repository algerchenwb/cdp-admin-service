package dbwrap

// interface for DBWrap
type IDBWrap interface {
	Query(query string, sortby interface{}, ascending interface{}) (interface{}, int, error)
	QueryAll(query string, sortby interface{}, ascending interface{}) (interface{}, int, error)
	QueryPage(query string, offset, limit int, sortby, ascending interface{}) (int, interface{}, int, error)
	Insert(info interface{}) (interface{}, int, error)
	Update(key, info interface{}) (interface{}, int, error)
	Delete(key interface{}) (int, error)
}

type IAreaDBWrap interface {
	Query(area int, query string, sortby interface{}, ascending interface{}) (interface{}, int, error)
	QueryAll(area int, query string, sortby interface{}, ascending interface{}) (interface{}, int, error)
	QueryPage(area int, query string, offset, limit int, sortby, ascending interface{}) (int, interface{}, int, error)
	Insert(area int, info interface{}) (interface{}, int, error)
	Update(area int, key, info interface{}) (interface{}, int, error)
	Delete(area int, key interface{}) (int, error)
}
