package cep

import (
	"encoding/json"
	"full_cycle_cep/pkg/core/middleware"
	"full_cycle_cep/pkg/domain/contracts/init_provider"
	"full_cycle_cep/pkg/domain/use_cases/viacep/get_viacep"
	"full_cycle_cep/pkg/shared/business_error"
	"full_cycle_cep/pkg/shared/constants"
	"full_cycle_cep/pkg/shared/log"
	"github.com/go-chi/chi/v5"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"net/http"
	"sync"
)

var lock sync.Mutex
var createCepRouteInstance CreateCepRoute

type CreateCepRoute struct {
	logger                    log.LoggerManagerInterface
	getViaCepUseCase          get_viacep.GetViaCepUseCaseInterface
	RtCepValidationMiddleware middleware.CepValidationMiddleware
	initProvider              *init_provider.InitProvider
}

func NewCreateCepRoute(
	logger log.LoggerManagerInterface,
	getViaCepUseCase get_viacep.GetViaCepUseCaseInterface,
	initProvider *init_provider.InitProvider,
) CreateCepRoute {
	if createCepRouteInstance == (CreateCepRoute{}) {
		lock.Lock()
		defer lock.Unlock()
		if createCepRouteInstance == (CreateCepRoute{}) {
			createCepRouteInstance = CreateCepRoute{
				logger:           logger,
				getViaCepUseCase: getViaCepUseCase,
				initProvider:     initProvider,
			}
		}
	}
	return createCepRouteInstance
}

func (c *CreateCepRoute) GetCepRoute() *chi.Mux {
	r := chi.NewRouter()
	r.With(c.RtCepValidationMiddleware.Validate).Get("/{cep}", c.Get)
	return r
}

func (c *CreateCepRoute) Get(w http.ResponseWriter, r *http.Request) {
	carrier := propagation.HeaderCarrier(r.Header)
	ctx := r.Context()
	ctx = otel.GetTextMapPropagator().Extract(ctx, carrier)
	tracer, _ := c.initProvider.GetProvider(ctx, "cep-service", "otel_collector:4317", "ms")
	ctx, span := tracer.Start(ctx, "ini cep")
	defer span.End()

	cep := chi.URLParam(r, "cep")
	res, err := c.getViaCepUseCase.Execute(cep)

	if err != nil {
		if cepNotFoundErr, ok := err.(*business_error.BusinessError); ok {
			if cepNotFoundErr.Message == string(constants.Cep_not_found) {
				http.Error(w, "can not find zipcode", http.StatusNotFound)
				return
			}

			if cepNotFoundErr.Message == string(constants.Invalid_Cep) {
				http.Error(w, "invalid zipcode", http.StatusUnprocessableEntity)
				return
			}
		}
		http.Error(w, "Error", http.StatusInternalServerError)
		return
	}

	jsonResponse, err := json.Marshal(res)
	if err != nil {
		http.Error(w, "can not find zipcode", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}
