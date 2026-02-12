package helper

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func ReadParams(r *http.Request) string {
	return httprouter.ParamsFromContext(r.Context()).ByName("id")
}
