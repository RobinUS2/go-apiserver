package apiserver

import (
	"github.com/RobinUS2/go-orm"
	"log"
	"strings"
	"strconv"
	"time"
	"fmt"
	"database/sql"
)

type BaseController struct {
	name           string
	customActions  []*BaseAction
	model          interface{}
	orm            *orm.Orm
	editableFields []string

	ModelFn func() interface{}
	SpecializeRowFn func(rawRow interface{}, data map[string]interface{}) interface{}
}

type BaseControllerI interface {
	Name() string
	CustomActions() []*BaseAction
	Orm() *orm.Orm
	InitController()
	Model() interface{}
}

func (controller *BaseController) GetWhere(request *BaseRequest) *Where {
	return &Where{
		//Query: "customerId = ?",
		//Vars: []interface{}{
		//	1,
		//},
	}
}

func (controller *BaseController) GetSingle(request *BaseRequest) interface{} {
	elm := controller.Model()
	where := controller.GetWhere(request)
	db := controller.Orm().Where(where.Query, where.Vars...).First(elm, request.GetID())
	if db.Error != nil {
		log.Printf("%v", db.Error)
		return nil
	}
	return elm
}

func (controller *BaseController) GetList(request *BaseRequest) interface{} {
	var modelRows = make([]interface{}, 0)
	var rows *sql.Rows
	var err error
	// needs to be directly on db to work
	clone := controller.Orm()
	where := controller.GetWhere(request)
	rows, err = clone.Model(controller.Model()).Where(where.Query, where.Vars...).Rows()
	if err != nil {
		log.Printf("%#v", err)
		return nil
	}
	columns, _ := rows.Columns()
	count := len(columns)
	values := make([]interface{}, count)
	valuePtrs := make([]interface{}, count)
	defer rows.Close()
	for rows.Next() {
		for i := range columns {
			valuePtrs[i] = &values[i]
		}
		rows.Scan(valuePtrs...)
		rowData := make(map[string]interface{})

		for i, col := range columns {
			rowData[col] = values[i]
		}

		// update row res
		rowRes := controller.SpecializeRow(controller.Model(), rowData)
		modelRows = append(modelRows, rowRes)

	}
	return modelRows
}

func ParseID(val interface{}) uint {
	if val == nil {
		return 0
	}
	if f, ok := val.(int64); ok {
		return uint(f)
	}
	u, _ := strconv.ParseUint(fmt.Sprintf("%s", val.([]byte)), 10, 64)
	return uint(u)
}

func ParseTime(val interface{}) *time.Time {
	if val == nil {
		return nil
	}
	if f, ok := val.(time.Time); ok {
		return &f
	}
	if f, ok := val.(*time.Time); ok {
		return f
	}
	return nil
}

func ParseString(val interface{}) string {
	if val == nil {
		return ""
	}
	if f, ok := val.(string); ok {
		return f
	}
	if f, ok := val.([]byte); ok {
		return string(f)
	}
	return ""
}

type Where struct {
	Query string
	Vars  []interface{}
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

func (controller *BaseController) Model() interface{} {
	return controller.ModelFn()
}

func (controller *BaseController) SpecializeRow(rawRow interface{}, data map[string]interface{}) interface{} {
	return controller.SpecializeRowFn(rawRow, data)
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
