package helper

import (
	"errors"
	"strings"

	"google.golang.org/grpc/status"
)

func ParseGRPCError(err error) error {
	if err == nil {
		return nil
	}

	st, ok := status.FromError(err)
	errMsg := ""
	if ok {
		errMsg = st.Message()
	} else {
		errMsg = err.Error()
	}
	errMsg = strings.ReplaceAll(errMsg, "租户", "合约")
	return errors.New(errMsg)
}
