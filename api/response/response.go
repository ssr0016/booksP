package response

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"practice/api/errors"
)

type Response interface {
	WriteTo(w http.ResponseWriter)
}

type response struct {
	status int
	body   *bytes.Buffer
	header http.Header
	err    error
}

func (r *response) SetHeader(key, value string) *response {
	r.header.Set(key, value)
	return r
}

func (r *response) WriteTo(w http.ResponseWriter) {
	header := w.Header()
	for k, v := range r.header {
		header[k] = v
	}

	if r.err != nil {
		rw := w.(ResponseWriter)
		rw.Error(r.err)
	}

	w.WriteHeader(r.status)
	w.Write(r.body.Bytes())
}

func Respond(status int, body interface{}) *response {
	var b []byte
	switch t := body.(type) {
	case []byte:
		b = t
	case string:
		b = []byte(t)
	default:
		var err error
		if b, err = json.Marshal(body); err != nil {
			return Error(http.StatusInternalServerError, fmt.Errorf("body json marshal"))
		}
	}

	return &response{
		status: status,
		body:   bytes.NewBuffer(b),
		header: make(http.Header),
	}
}

func JSON(status int, body interface{}) *response {
	return Respond(status, body).SetHeader("Content-Type", "application/json")
}

func Text(status int, body interface{}) *response {
	return Respond(status, body).SetHeader("Content-Type", "text/plain")
}

func HTML(status int, body interface{}) *response {
	return Respond(status, body).SetHeader("Content-Type", "text/html")
}

func XML(status int, body interface{}) *response {
	return Respond(status, body).SetHeader("Content-Type", "application/xml")
}

func Excel(status int, body interface{}) *response {
	res := Respond(status, body)
	res.SetHeader("Content-Type", "application/octet-stream")
	res.SetHeader("Content-Transfer-Encoding", "binary")
	return res
}

func Success(message string) *response {
	resp := make(map[string]interface{})
	resp["message"] = message
	return JSON(http.StatusOK, resp)
}

func Error(status int, err error) *response {
	var result errors.ErrorStatus
	switch err := err.(type) {
	case errors.ErrorStatus:
		result = err
	default:
		result = errors.New(errors.ErrGeneral, err.Error())
	}

	res := JSON(status, result)
	if err != nil {
		res.err = err
	}

	return res
}
