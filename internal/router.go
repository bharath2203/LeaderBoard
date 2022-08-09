package internal

import (
	"TopKScores/internal/handlers"
	"github.com/gorilla/mux"
	"net/http"
)

func NewRouter(handlerService *handlers.Handler) *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range getRoutes(handlerService) {
		var handler http.Handler
		handler = route.Handler

		router.
			Methods(route.Method).
			Path(route.Path).
			Name(route.Name).
			Handler(handler)
	}
	return router
}
