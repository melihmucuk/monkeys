package monkeys

import (
	"log"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
)

// Logger is a middleware that logs request method, request uri and processing time.
func Logger(next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		start := time.Now()
		next(w, r, p)
		log.Printf("%s %s %v", r.Method, r.RequestURI, time.Since(start))
	}
}
