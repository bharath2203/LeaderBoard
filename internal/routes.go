package internal

import (
	"TopKScores/internal/handlers"
	"net/http"
)

type Route struct {
	Name    string
	Method  string
	Path    string
	Handler http.HandlerFunc
}

type Routes []Route

func getRoutes(handlerService *handlers.Handler) Routes {
	var routes = Routes{
		Route{
			Name:    "AddScore of the Game",
			Method:  http.MethodPost,
			Path:    "/v1/game/add-score",
			Handler: handlerService.AddSingleScore,
		},
		Route{
			Name:    "AddScore of the Game to the queue",
			Method:  http.MethodPost,
			Path:    "/v1/game/queue/add-score",
			Handler: handlerService.AddSingleScoreToQueue,
		},
		Route{
			Name:    "Get Top K scores",
			Method:  http.MethodGet,
			Path:    "/v1/game/top-scores",
			Handler: handlerService.GetTopKScores,
		},
	}
	return routes
}
