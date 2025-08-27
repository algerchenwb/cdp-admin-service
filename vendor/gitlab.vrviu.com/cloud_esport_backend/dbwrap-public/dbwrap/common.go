package dbwrap

import "reflect"

type DBWrapInterface interface {
	// 单条查询
	Query(query string, sortby, ascending interface{}) (interface{}, int, error)
	// 全部查询
	QueryAll(query string, sortby, ascending interface{}) (interface{}, int, error)
	// 分页查询
	QueryPage(query string, offset, limit int, sortby, ascending interface{}) (int, interface{}, int, error)
	// 插入
	Insert(info interface{}) (interface{}, int, error)
	// 更新
	Update(key, info interface{}) (interface{}, int, error)
	// 删除
	Delete(key interface{}) (int, error)
	// 原始SQL语句
	RawQuery(sql string, rettype reflect.Type) (interface{}, int, error)

	// 单条查询
	AreaQuery(areaType int, query string, sortby, ascending interface{}) (interface{}, int, error)
	// 全部查询
	AreaQueryAll(areaType int, query string, sortby, ascending interface{}) (interface{}, int, error)
	// 分页查询
	AreaQueryPage(areaType int, query string, offset, limit int, sortby, ascending interface{}) (int, interface{}, int, error)
	// 插入
	AreaInsert(areaType int, info interface{}) (interface{}, int, error)
	// 更新
	AreaUpdate(areaType int, key, info interface{}) (interface{}, int, error)
	// 删除
	AreaDelete(areaType int, key interface{}) (int, error)
}

var (
	_tUnionUserAreaInfo                     = "t_union_user_area_info"
	_tUnionGameMetaInfo                     = "t_union_game_meta_info"
	_tUnionUserBasicInfo                    = "t_union_user_basic_info"
	_tUnionUserThirdInfo                    = "t_union_user_third_info"
	_tUnionGameAreaInfo                     = "t_union_game_area_info"
	_tUnionAreaInfo                         = "t_union_area_info"
	_tUnionUserGameAreaStorageInfo          = "t_union_user_game_area_storage_info"
	_tUnionArchiveBakupTaskInfo             = "t_union_archive_backup_task_info"
	_tUnionArchiveImportInfo                = "t_union_archive_import_info"
	_tUnionArchiveManageInfo                = "t_union_archive_manage_info"
	_tUnionAccessPointInfo                  = "t_union_access_point_info"
	_tUnionAccessPointOuterInfo             = "t_union_access_point_outer_info"
	_tUnionAreaZoneInfo                     = "t_union_area_zone_info"
	_tAreaAccessPointManagerInfo            = "t_area_access_point_manager_info"
	_tZoneAccessPointManagerInfo            = "t_zone_access_point_manager_info"
	_tUnionRemoveTaskInfo                   = "t_union_remove_task_info"
	_tUnionCleanupTaskInfo                  = "t_union_cleanup_task_info"
	_tUnionAccessPointGeographyScheduleInfo = "t_union_access_point_geography_schedule_info"
	_tGameVirtualMachinePortInfo            = "t_game_virtual_machine_port_info"
	_tAreaAccessPointInfo                   = "t_area_access_point_info"
	_tAreaAccessPointOuterInfo              = "t_area_access_point_outer_info"
	_tAreaAccessChannelInfo                 = "t_area_access_channel_info"
	_tGameVirtualMachineInfo                = "t_game_virtual_machine_info"
)

// 模块类型枚举
const (
	MDAL int = 100000025
)

// 接口类型枚举
const (
	IMCDAL_Query  int = 300000089
	IMCDAL_Insert int = 300000090
	IMCDAL_Update int = 300000091
)

const (
	ErrCodeOK               = 0
	ESErrCodeQueryFailed    = 3200
	ESErrCodeConnFailed     = 3201
	DBErrCodeCreateFailed   = 3100
	DBErrCodeUpdateFailed   = 3101
	DBErrCodeQueryFailed    = 3102
	DBErrCodeInsertFailed   = 3103
	DBErrCodeDeleteFailed   = 3104
	DBErrCodeRawQueryFailed = 3105
)

const (
	ErrStrOK             = "ok. table:%s"
	ErrStrCreateFailed   = "create failed. table:%s"
	ErrStrUpdateFailed   = "update failed. table:%s"
	ErrStrQueryFailed    = "query failed. table:%s"
	ErrStrInsertFailed   = "insert failed. table:%s"
	ErrStrDeleteFailed   = "delete failed. table:%s"
	ErrStrRawQueryFailed = "raw sql exec failed"
)
