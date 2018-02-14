package monkeys

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// Available Methods
const (
	GET    = "GET"
	POST   = "POST"
	PUT    = "PUT"
	DELETE = "DELETE"
	HEAD   = "HEAD"
	PATCH  = "PATCH"
)

// GetController TO-DO
type GetController interface {
	Get() *Response
}

// GetByIDController TO-DO
type GetByIDController interface {
	GetByID(ID string) *Response
}

// PostController TO-DO
type PostController interface {
	Post(values Request) *Response
}

// Request TO-DO
type Request map[string]interface{}

// Response TO-DO
type Response struct {
	StatusCode   int         `json:"-"`
	Meta         interface{} `json:"meta,omitempty"`
	Data         interface{} `json:"data,omitempty"`
	ErrorCode    int         `json:"error_code,omitempty"`
	ErrorMessage string      `json:"error_message,omitempty"`
}

// API TO-DO
type API struct {
	router            *httprouter.Router
	routerInitialized bool
	middlewares       []func(next httprouter.Handle) httprouter.Handle
}

// NewAPI TO-DO
func NewAPI() *API {
	return &API{}
}

// Router TO-DO
func (api *API) Router() *httprouter.Router {
	if api.routerInitialized {
		return api.router
	}
	api.router = httprouter.New()
	api.routerInitialized = true
	return api.router
}

// ErrorResponse TO-DO
func ErrorResponse(statusCode, errorCode int, message string, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(statusCode)
	meta := make(map[string]interface{})
	meta["error_message"] = message
	json.NewEncoder(w).Encode(&Response{ErrorCode: errorCode, ErrorMessage: message})
}

// SuccessResponse TO-DO
func SuccessResponse(response *Response, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		ErrorResponse(http.StatusInternalServerError, 2, err.Error(), w)
	}
}

func (api *API) requestHandler(resource interface{}) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		data := new(Response)
		switch r.Method {
		case GET:
			ID := p.ByName("ID")
			if ID != "" {
				res, ok := resource.(GetByIDController)
				if !ok {
					ErrorResponse(http.StatusMethodNotAllowed, 1, "Method is not implemented!", w)
					return
				}

				data = res.GetByID(ID)
			} else {
				res, ok := resource.(GetController)
				if !ok {
					ErrorResponse(http.StatusMethodNotAllowed, 1, "Method is not implemented!", w)
					return
				}

				data = res.Get()
			}
		case POST:
			res, ok := resource.(PostController)
			if !ok {
				ErrorResponse(http.StatusMethodNotAllowed, 1, "Method is not implemented!", w)
				return
			}

			rBody := make(Request)
			err := json.NewDecoder(r.Body).Decode(&rBody)
			if err != nil {
				ErrorResponse(http.StatusBadRequest, 3, err.Error(), w)
				return
			}

			data = res.Post(rBody)
		}

		SuccessResponse(data, w)
	}
}

// NewEndpoint TO-DO
func (api *API) NewEndpoint(method, path string, resource interface{}) {
	handler := api.requestHandler(resource)
	for _, middleware := range api.middlewares {
		handler = middleware(handler)
	}
	api.Router().Handle(method, path, handler)
}

// NewEndpointGroup TO-DO
func (api *API) NewEndpointGroup(path string, resource interface{}) {
	api.NewEndpoint("GET", path, resource)
	api.NewEndpoint("GET", path+"/:ID", resource)
	api.NewEndpoint("POST", path, resource)
}

// Use TO-DO
func (api *API) Use(middleware func(next httprouter.Handle) httprouter.Handle) {
	api.middlewares = append(api.middlewares, middleware)
}

// Start TO-DO
func (api *API) Start(port int) error {
	if !api.routerInitialized {
		return errors.New("You must add at least one resource to this API")
	}
	portString := fmt.Sprintf(":%d", port)
	log.Printf("http server running on %s \n", portString)
	return http.ListenAndServe(portString, api.Router())
}
