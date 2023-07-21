package api_helpers

import (
	"encoding/json"
	"log"
	"net/http"
)

func ResponseJSON(
	w *http.ResponseWriter,
	statusCode int,
	mp *map[string]interface{},
	message *string,
) {
	if mp == nil {
		mp = &map[string]interface{}{}
	}
	if message != nil {
		(*mp)[MessageField] = *message
	}

	(*w).Header().Set("Content-Type", "application/json")
	res, err := json.Marshal(mp)
	if err != nil {
		(*w).WriteHeader(http.StatusInternalServerError)
		log.Printf("Error happened in JSON marshal. Err: %s\n", err.Error())
		return
	}

	(*w).WriteHeader(statusCode)
	(*w).Write(res)
}
