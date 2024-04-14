package utils

import (
	"encoding/json"
	"net/http"
)

// type RespError struct {
// 	Msg error `json:"error,omitempty"`
// }

// type RespOk struct {
// 	Msg any
// }

// func ResponseErr(w http.ResponseWriter, r *http.Request, message error, httpStatusCode int) {
// 	w.WriteHeader(httpStatusCode)
// 	render.JSON(w, r, RespError{Msg: message})
// }

// func ResponseOK(w http.ResponseWriter, r *http.Request, message any, httpStatusCode int) {
// 	w.WriteHeader(httpStatusCode)
// 	render.JSON(w, r, RespOk{Msg: message})
// }

const (
	Error = "error"
)

func ResponseJSON(w http.ResponseWriter, key string, message any, httpStatusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatusCode)
	data := make(map[string]any)
	data[key] = message
	jsonResp, _ := json.Marshal(data)
	w.Write(jsonResp)
}
