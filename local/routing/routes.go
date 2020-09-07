package routing

import (
	"net/http"

	"github.com/gorilla/mux"
)

// Init initializes the router middleware
func Init() *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/unfollows", func(w http.ResponseWriter, r *http.Request) {

	})

	return router
}
