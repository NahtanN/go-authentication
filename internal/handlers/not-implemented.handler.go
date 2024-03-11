package handlers

import "net/http"

func NotImplementedHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Handler not implemented", http.StatusNotImplemented)
}
