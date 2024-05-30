package http

import (
	"full_cycle_cep/pkg/presentation/http/cep"
	"full_cycle_cep/pkg/presentation/http/orchestrator"
	"full_cycle_cep/pkg/shared/log"
	"github.com/go-chi/chi/v5"
)

type Handlers struct {
	logger                    log.LoggerManagerInterface
	HnCreateCepRoute          cep.CreateCepRoute
	HnCreateOrchestratorRoute orchestrator.CreateOrchestratorRoute
}

func ProvideHandlers(
	logger log.LoggerManagerInterface,
	createCepRoute cep.CreateCepRoute,
	createOrchestratorRoute orchestrator.CreateOrchestratorRoute,
) *Handlers {
	return &Handlers{
		logger:                    logger,
		HnCreateCepRoute:          createCepRoute,
		HnCreateOrchestratorRoute: createOrchestratorRoute,
	}
}

func (h *Handlers) GetRoutes() *chi.Mux {
	r := chi.NewRouter()
	r.Mount("/cep", h.HnCreateCepRoute.GetCepRoute())
	r.Mount("/orchestrator", h.HnCreateOrchestratorRoute.GetOrchestratorRoute())
	return r
}
