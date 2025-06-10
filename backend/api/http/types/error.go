package types

import (
	"encoding/json"
	"net/http"
)

func ProcessError(w http.ResponseWriter, err error, resp any, successResponse int) {
	if err != nil {
		http.Error(w, "Internal Error", http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(successResponse)
	if resp != nil {
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			http.Error(w, "Internal Error", http.StatusInternalServerError)
			return
		}
	}
}
