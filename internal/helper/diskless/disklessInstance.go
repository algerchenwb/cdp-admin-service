package diskless

import instance_types "cdp-admin-service/internal/proto/instance_service/types"

type HTTPCommonHead struct {
	Code   int    `json:"code"`
	Msg    string `json:"msg,omitempty"`
	Detail string `json:"detail,omitempty"`
}

type GetInstanceList struct {
	Head HTTPCommonHead                       `json:"ret"`
	Body instance_types.ListInstancesResponse `json:"body,omitempty"`
}
