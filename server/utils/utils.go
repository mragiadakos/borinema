package utils

type ErrorMsg struct {
	VariableErrors map[string]error
	Error          error
	Status         int
}

func NewErrorMsg() *ErrorMsg {
	errMsg := ErrorMsg{
		VariableErrors: map[string]error{},
		Error:          nil,
	}
	return &errMsg
}

func (e *ErrorMsg) HasErrors() bool {
	if e.VariableErrors != nil {
		return len(e.VariableErrors) > 0 || e.Error != nil
	}
	return e.Error != nil
}

func (e *ErrorMsg) GetVariable(s string) error {
	err, ok := e.VariableErrors[s]
	if ok {
		return err
	}
	return nil
}

func (e *ErrorMsg) Json() map[string]interface{} {
	m := map[string]interface{}{}
	errs := map[string]string{}
	for k, v := range e.VariableErrors {
		errs[k] = v.Error()
	}
	m["variable_errors"] = errs
	if e.Error != nil {
		m["error"] = e.Error.Error()
	}
	return m
}
