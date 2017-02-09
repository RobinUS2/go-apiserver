package apiserver

func (controller *BaseController) Delete(r *BaseRequest) {
	// Fetch
	value := controller.GetSingle(r)
	if value == nil {
		return
	}

	// Delete
	controller.GetModel().Delete(controller.Orm(), value)
	r.SetValue("deleted", true)
}
