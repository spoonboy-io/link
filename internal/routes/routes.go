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

// Home - there is nothing useful being done here, but could be used for ping check
func (r *Routes) Home(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "text/plain")

	res := "Link is up\n"

	r.App.Logger.Info("Served GET / request - 200 OK")
	_, _ = fmt.Fprint(w, res)
}
