package cep

import (
	"encoding/json"
	"full_lab_tracing/pkg/shared/log"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
	"sync"
)

var lock sync.Mutex
var createCepRouteInstance CreateCepRoute

type CreateCepRoute struct {
	logger                    log.LoggerManagerInterface
	getViaCepUseCase          get_viacep.GetViaCepUseCaseInterface
	RtCepValidationMiddleware middleware.CepValidationMiddleware
}

func NewCreateCepRoute(
	logger log.LoggerManagerInterface,
	getViaCepUseCase get_viacep.GetViaCepUseCaseInterface,
) CreateCepRoute {
	if createCepRouteInstance == (CreateCepRoute{}) {
		lock.Lock()
		defer lock.Unlock()
		if createCepRouteInstance == (CreateCepRoute{}) {
			createCepRouteInstance = CreateCepRoute{
				logger:           logger,
				getViaCepUseCase: getViaCepUseCase,
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
