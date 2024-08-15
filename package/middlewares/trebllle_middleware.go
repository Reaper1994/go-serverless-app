package middlewares

import (
	"github.com/treblle/treblle-go"
	"net/http"
)

// TreblleMiddleware wraps the given handler with Treblle middlewares for Api ops.
func TreblleMiddleware(apiKey string, projectId string, next http.Handler) http.Handler {

	treblle.Configure(treblle.Configuration{
		APIKey:    apiKey,
		ProjectID: projectId,
	})

	return treblle.Middleware(next)
}
