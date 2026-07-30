package http

import (
	"delivery/internal/generated/servers"
	"encoding/json"
)

func WrapError(code int, err error) []byte {
	e := servers.Error{
		Code:    int32(code),
		Message: err.Error(),
	}
	var bytes []byte
	bytes, _ = json.MarshalIndent(e, "", "    ")

	return bytes
}
