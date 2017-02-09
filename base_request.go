package apiserver

import (
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

// Incoming request to the API server

type BaseRequest struct {
	Response http.ResponseWriter
	Request  *http.Request
	Params   httprouter.Params

	responseObj *ApiResponse
	body        []byte
	bodyJsonMap map[string]interface{}
}

func (request *BaseRequest) GetParam(key string) string {
	var lowerKey = strings.ToLower(key)
	var res string = ""

	// Json?
	var contentType string = request.Request.Header.Get("Content-Type")
	if contentType == "application/json" {
		// Lazy init
		if request.bodyJsonMap == nil {
			if len(request.body) == 0 {
				body, err := ioutil.ReadAll(request.Request.Body)
				if err == nil {
					request.body = body
				} else {
					log.Printf("Failed to read request body %s", err)
				}
			}
			jsonErr := json.Unmarshal(request.body, &request.bodyJsonMap)
			if jsonErr != nil {
				log.Printf("Failed to read json body %s", jsonErr)
			}
		}

		// Fetch
		if request.bodyJsonMap != nil {
			// Fetch
			if request.bodyJsonMap[key] != nil {
				res = fmt.Sprintf("%s", request.bodyJsonMap[key])
				if len(res) > 0 {
					return res
				}
			}
			// Fetch lowercase
			if request.bodyJsonMap[lowerKey] != nil {
				res = fmt.Sprintf("%s", request.bodyJsonMap[lowerKey])
				if len(res) > 0 {
					return res
				}
			}
		}
	}

	// Primary from router params, fallback to query
	res = request.Params.ByName(key)
	if len(res) > 0 {
		return res
	}

	// Form values (tries body and query)
	res = request.Request.FormValue(key)
	if len(res) > 0 {
		return res
	}

	// Form values: lowercase
	res = request.Request.FormValue(lowerKey)
	if len(res) > 0 {
		return res
	}

	// Query params
	res = request.Request.URL.Query().Get(key)
	if len(res) > 0 {
		return res
	}
	// Query params: lowercase
	res = request.Request.URL.Query().Get(lowerKey)
	return res
}

func (request *BaseRequest) GetID() string {
	return strings.TrimSpace(request.GetParam("id"))
}

func (request *BaseRequest) GetFilterVars() []interface{} {
	// @todo support
	return nil
}

func (request *BaseRequest) GetFilterQuery() interface{} {
	if len(request.GetID()) > 0 {
		return "id = ?"
	}

	for k, v := range request.Request.URL.Query() {
		if strings.Index(k, "filter_like_") == 0 {
			// like filter
			log.Printf("filter %s LIKE %s", k, v[0])
		} else if strings.Index(k, "filter_") == 0 {
			// strict filter
			log.Printf("filter %s = %s", k, v[0])
		}
	}
	// @todo support
	return nil
}

func (request *BaseRequest) init() {
	// @todo move filter query + vars + id logic to here and only do it once
}

// Set response error
func (request *BaseRequest) SetError(value interface{}) {
	valueStr := fmt.Sprintf("%s", value)
	request.responseObj.errored = true
	request.responseObj.Error(valueStr)
}

// Set response value
func (request *BaseRequest) SetValue(key string, value interface{}) {
	request.responseObj.Set(key, value)
}

func newBaseRequest(w http.ResponseWriter, r *http.Request, p httprouter.Params) *BaseRequest {
	br := &BaseRequest{
		Response:    w,
		Request:     r,
		Params:      p,
		responseObj: NewResponse(),
	}
	br.init()
	return br
}
