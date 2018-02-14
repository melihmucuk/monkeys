![monkeys team](http://i64.tinypic.com/24ewfwh.jpg)

# monkeys
tiny rest framework for golang

## Install

```
go get -u github.com/melihmucuk/monkeys
```
## Usage

```go
package main

import (
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/melihmucuk/monkeys"
)

type ItemController struct{}

func (ic ItemController) Get() *monkeys.Response {
	items := []string{"item1", "item2"}
	data := map[string][]string{"items": items}
	return &monkeys.Response{StatusCode: 200, Data: data}
}

func (ic ItemController) GetByID(ID string) *monkeys.Response {
	return &monkeys.Response{StatusCode: 200, Data: ID}
}

func (ic ItemController) Post(body monkeys.Request) *monkeys.Response {
	return &monkeys.Response{StatusCode: 200, Data: body}
}


type PingController struct{}

func (pc PingController) Get() *monkeys.Response {
	data := make(map[string]interface{})
	data["respond"] = "pong"
	return &monkeys.Response{StatusCode: 200, Data: data}
}

// HostMiddleware custom middleware
func HostMiddleware(next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		log.Printf("%s", r.Host)
		if r.Host == "localhost:5000" {
			next(w, r, p)
			return
		}
		monkeys.ErrorResponse(http.StatusForbidden, 4, "Host not allowed!", w)
	}
}

func main() {
	api := monkeys.NewAPI()
	api.Use(monkeys.Logger)
	api.Use(HostMiddleware)
	api.NewEndpointGroup("/item", ItemController{})
	api.NewEndpoint("GET", "/ping", PingController{})
	api.Start(5000)
}

```

## TODO

- [ ] Integrate all http methods (PUT, PATCH etc.)
- [ ] Create documentation
- [ ] Write tests
- [ ] Add benchmark result
- [ ] Create built-in middleware for CORS
