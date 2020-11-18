package utils

import (
	"encoding/json"
	"net/http"
)

const (
	//StatusOk when request is success
	StatusOk = "ok"
	// StatusNOk when request fails
	StatusNOk = "nok"
)

//Message to represent messages in response
type Message struct {
	Message string `json:"message"`
}

//Response is the struct for general api response
type Response struct {
	Status string       `json:"status"`
	Error  string       `json:"error,omitempty"`
	Result *interface{} `json:"result,omitempty"`
}

//Send function to send general api response
func Send(w http.ResponseWriter, status int, payload interface{}) {
	response := Response{
		Status: StatusOk,
		Result: &payload,
	}
	result, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(result)
}

//Fail function to send general api response in case of error
func Fail(w http.ResponseWriter, status int, details string) {
	response := &Response{
		Status: StatusNOk,
		Error:  details,
	}
	result, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(result)
}
