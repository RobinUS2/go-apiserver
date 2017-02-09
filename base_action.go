package apiserver

import (
	"bytes"
)

// Sub-action on a controller (e.g. /account/movies "movies" is the action)

type BaseAction struct {
	Handle func(r *BaseRequest)
	Route  string
	Method string
}

func (action *BaseAction) FullRoute(controller BaseControllerI, single bool) string {
	var buffer bytes.Buffer
	buffer.WriteString("/") // always starts with a forward slash
	buffer.WriteString(controller.Name())
	if action.Route != "" {
		buffer.WriteString("/" + action.Route)
	}
	if single {
		buffer.WriteString("/:id")
	}
	return buffer.String()
}

func NewAction(method string, handle func(r *BaseRequest), route string) *BaseAction {
	return &BaseAction{
		Handle: handle,
		Method: method,
		Route:  route,
	}
}
