package routes

import (
	"fmt"
	"net/http"

	"github.com/spoonboy-io/link/internal"
)

// Routes makes the application context, logger and config availalble to the handlers
type Routes struct {
	App *internal.App
}

// Ping provides an endpoint to check the server is running and responding
func (r *Routes) Ping(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "text/plain")

	res := "Hello from Link!\n"

	r.App.Logger.Info("Served GET / request - 200 OK")
	_, _ = fmt.Fprint(w, res)
}
