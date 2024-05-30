package orchestrator

import (
	"full_cycle_cep/pkg/core/middleware"
	"full_cycle_cep/pkg/shared/log"
	"github.com/go-chi/chi/v5"
	"net/http"
	"sync"
)

var lock sync.Mutex
var createOrchestratorInstance CreateOrchestratorRoute

type CreateOrchestratorRoute struct {
	logger                    log.LoggerManagerInterface
	RtCepValidationMiddleware middleware.CepValidationMiddleware
}

func NewOrchestratorRoute(
	logger log.LoggerManagerInterface,
) CreateOrchestratorRoute {
	if createOrchestratorInstance == (CreateOrchestratorRoute{}) {
		lock.Lock()
		defer lock.Unlock()
		if createOrchestratorInstance == (CreateOrchestratorRoute{}) {
			createOrchestratorInstance = CreateOrchestratorRoute{
				logger: logger,
			}
		}
	}
	return createOrchestratorInstance
}

func (c *CreateOrchestratorRoute) GetOrchestratorRoute() *chi.Mux {
	r := chi.NewRouter()
	r.With(c.RtCepValidationMiddleware.Validate).Get("/{cep}", c.Get)
	return r
}

func (c *CreateOrchestratorRoute) Get(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(""))
}
