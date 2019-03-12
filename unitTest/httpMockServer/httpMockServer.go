package httpMockServer

import (
	"encoding/json"
	"net/http"
)

func Routes()  {
	http.HandleFunc("/sendjson", SendJSON)
}

func SendJSON(rw http.ResponseWriter, r *http.Request) {
	u := struct {
		Name string
		Age int
	}{
		Name:"张三",
		Age: 28,
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	json.NewEncoder(rw).Encode(u)
}

