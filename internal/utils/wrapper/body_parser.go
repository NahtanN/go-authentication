package wrapper_utils

import (
	"encoding/json"
	"net/http"
)

func BodyParser[R interface{}](request *R, req *http.Request) error {
	return json.NewDecoder(req.Body).Decode(&request)
}
