package apiserver

import (
	"github.com/RobinUS2/golang-jresp"
)

type ApiResponse struct {
	jresp.JResp
	hasError  bool
	errorCode int
}

func (a *ApiResponse) SetErrorCode(errorCode int) {
	a.errorCode = errorCode
}

func NewResponse() *ApiResponse {
	return &ApiResponse{
		JResp:    *jresp.NewJsonResp(),
		hasError: false,
	}
}
