package response

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// RespForm ...
type RespForm struct {
	Type    int
	Msg     string
	Content interface{}
}

// PrepareResponse ...
func PrepareResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	w.WriteHeader(statusCode)
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}
}

// JsonResponse ...
func JsonResponse(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET")
		w.Header().Set("Content-Type", "application/json")
		next(w, req)
	}
}
