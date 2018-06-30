package apiserver

import (
	"github.com/RobinUS2/go-orm"
	"log"
	"strings"
)

type BaseController struct {
	name           string
	customActions  []*BaseAction
	model          interface{}
	orm            *orm.Orm
	editableFields []string
}

type BaseControllerI interface {
	Name() string
	CustomActions() []*BaseAction
	// Get(r *BaseRequest)
	// Put(r *BaseRequest)
	// Post(r *BaseRequest)
	// Delete(r *BaseRequest)
	Orm() *orm.Orm
	InitController()
}

// Get parameters (lower case key => value map)
func (controller *BaseController) GetEntityParams(request *BaseRequest) map[string]interface{} {
	m := make(map[string]interface{})
	for _, key := range controller.GetEditableFields() {
		val := request.GetParam(key)
		// Skip empty params
		if len(val) == 0 {
			continue
		}
		m[strings.ToLower(key)] = val
	}
	log.Printf("%v", m)
	return m
}

func (controller *BaseController) GetEditableFields() []string {
	return controller.editableFields
}

func (controller *BaseController) Orm() *orm.Orm {
	return controller.orm
}

func (controller *BaseController) CustomActions() []*BaseAction {
	return controller.customActions
}

func (controller *BaseController) Name() string {
	return controller.name
}

func (controller *BaseController) RegisterAction(action *BaseAction) {
	if controller.customActions == nil {
		controller.customActions = make([]*BaseAction, 0)
	}
	controller.customActions = append(controller.customActions, action)
}

func NewController(orm *orm.Orm, name string) BaseController {
	if orm == nil {
		panic("Controller can not be created without ORM")
	}
	c := &BaseController{
		name: name,
		orm:  orm,
	}
	return *c
}
