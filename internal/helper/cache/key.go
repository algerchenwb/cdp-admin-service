package cache

import "strconv"

const (
	UserOnlinePrefix    = "user_online_prefix:"
	UserPermPrefix      = "user_perm_prefix:"
	RoleAdminPrefix     = "role_admin"
	PermIsPrivatePrefix = "perm_is_private"
	PermIsLoggingPrefix = "perm_is_logging_prefix:"
)

func UserOnlineKey(userId int64) string {
	return UserOnlinePrefix + strconv.FormatInt(userId, 10)
}

func UserPermKey(userId int64) string {
	return UserPermPrefix + strconv.FormatInt(userId, 10)
}

func RoleAdminKey() string {
	return RoleAdminPrefix
}

func PermIsPrivateKey() string {
	return PermIsPrivatePrefix
}

func PermIsLoggingKey(perm string) string {
	return PermIsLoggingPrefix + perm
}
