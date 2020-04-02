package api

import (
	"net/http"
	"os"

	"github.com/gorilla/handlers"
)

func LogRequest(h http.Handler) http.Handler {
	return handlers.CombinedLoggingHandler(os.Stdout, h)
}
