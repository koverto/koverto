// Package health defines health checking and status reporting for the Koverto
// service and all services it connects to.
package health

import (
	"context"
	"encoding/json"
	"net/http"
	"sync"

	"github.com/koverto/koverto/internal/pkg/resolver"

	debug "github.com/micro/go-micro/v2/debug/service/proto"
)

// Handler implements http.Handler to respond to health check requests over
// HTTP.
type Handler struct {
	*resolver.Resolver
}

type healthResponse struct {
	*debug.HealthResponse
	Error   error
	Service string
}

func (r *healthResponse) GetStatus() interface{} {
	if r.Error != nil {
		return r.Error
	}

	return r.HealthResponse.GetStatus()
}

// NewHandler returns a new Handler.
func NewHandler(r *resolver.Resolver) http.Handler {
	return &Handler{r}
}

func (h *Handler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	wg := &sync.WaitGroup{}
	wg.Add(h.Length())

	statuses := make(chan *healthResponse, h.Length())

	for _, name := range h.ClientSet.Keys() {
		go func(name string) {
			defer wg.Done()

			req := h.Service.Client().NewRequest(name, "Debug.Health", &debug.HealthRequest{})
			rsp := &debug.HealthResponse{}

			err := h.Service.Client().Call(context.Background(), req, rsp)
			statuses <- &healthResponse{rsp, err, name}
		}(name)
	}

	wg.Wait()
	close(statuses)

	ok := true
	response := make(map[string]interface{}, 3)

	for status := range statuses {
		if status.GetStatus() != "ok" {
			ok = false
		}

		response[status.Service] = status.GetStatus()
	}

	if ok {
		rw.WriteHeader(http.StatusOK)
	} else {
		rw.WriteHeader(http.StatusInternalServerError)
	}

	j, _ := json.Marshal(response)
	_, _ = rw.Write(j)
}
