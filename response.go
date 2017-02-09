package apiserver

import (
	"github.com/RobinUS2/golang-jresp"
)

type ApiResponse struct {
	jresp.JResp
	errored bool
}

func NewResponse() *ApiResponse {
	return &ApiResponse{
		JResp:   *jresp.NewJsonResp(),
		errored: false,
	}
}
