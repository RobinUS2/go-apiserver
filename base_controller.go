package apiserver

import (
	"github.com/RobinUS2/go-orm"
	"log"
	"reflect"
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
	SetModel(model interface{})
	GetModel() orm.ModelI
	Orm() *orm.Orm
	InitController()
}

func (controller *BaseController) SetModel(model interface{}) {
	controller.model = model
	controller._parseEditableFields()
}

func (controller *BaseController) GetModel() orm.ModelI {
	return controller.model.(orm.ModelI)
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

func (controller *BaseController) _parseEditableFields() {
	// Reflect fields
	t := controller.GetModel().GetStructType()
	fieldNames := make([]string, 0)
	for i := 0; i < t.NumField(); i++ {
		var field reflect.StructField = t.Field(i)
		// Can not be anonymous
		if field.Anonymous {
			continue
		}
		var tags reflect.StructTag = field.Tag
		// Needs to be tagged for gorm
		if len(tags.Get("gorm")) < 1 {
			continue
		}
		// log.Printf("%d %#v", i, field)
		fieldNames = append(fieldNames, field.Name)
	}
	controller.editableFields = fieldNames
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

func (controller *BaseController) GetSingle(r *BaseRequest) interface{} {
	// fetch single
	id := r.GetID()
	if len(id) < 1 {
		r.SetError("Please provide an ID")
		return nil
	}
	value := controller.GetModel().First(controller.Orm(), r.GetFilterQuery(), id).Value()
	if value == nil {
		r.SetError("Object not found")
		return nil
	}
	return value
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
