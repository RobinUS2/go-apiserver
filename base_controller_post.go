package apiserver

func (controller *BaseController) Post(r *BaseRequest) {
	res := controller.GetModel().Create(controller.Orm(), controller.GetEntityParams(r))
	if res.Error() != nil {
		r.SetError(res.Error())
	}
	value := res.Value()
	r.SetValue(controller.GetModel().GetName(), value)
}
