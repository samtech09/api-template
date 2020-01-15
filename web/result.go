package web

import "net/http"

//Result is datatype that will be retured from all APIs
type Result struct {
	HTTPStatus int
	ErrCode    int
	ErrText    string
	Data       string //interface{}
	//IsJSON  bool
}

//NewResult object
func NewResult() *Result {
	r := Result{}
	r.HTTPStatus = http.StatusOK
	r.ErrCode = 1
	return &r
}

//NewResultFilled creates a new instance of Result struct with given parameters
func NewResultFilled(data string, httpstatus int, errcode int, errmsg string) *Result {
	r := Result{Data: data, ErrCode: errcode, ErrText: errmsg, HTTPStatus: httpstatus}
	return &r
}
