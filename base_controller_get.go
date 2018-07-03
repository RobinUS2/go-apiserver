package apiserver

func (controller *BaseController) Get(r *BaseRequest) {
	id := r.GetID()
	if len(id) < 1 {
		// fetch list
		values := controller.GetList(r)
		r.SetValue(controller.Name()+"s", values)
	} else {
		// fetch single
		value := controller.GetSingle(r)
		if value == nil {
			return
		}
		r.SetValue(controller.Name(), value)
	}
}
