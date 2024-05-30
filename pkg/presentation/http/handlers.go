package http

import (
	"full_cycle_cep/pkg/presentation/http/cep"
	"full_cycle_cep/pkg/shared/log"
	"github.com/go-chi/chi/v5"
)

type Handlers struct {
	logger           log.LoggerManagerInterface
	HnCreateCepRoute cep.CreateCepRoute
}

func ProvideHandlers(
	logger log.LoggerManagerInterface,
	createCepRoute cep.CreateCepRoute,
) *Handlers {
	return &Handlers{
		logger:           logger,
		HnCreateCepRoute: createCepRoute,
	}
}

func (h *Handlers) GetRoutes() *chi.Mux {
	r := chi.NewRouter()
	r.Mount("/cep", h.HnCreateCepRoute.GetCepRoute())
	return r
}
