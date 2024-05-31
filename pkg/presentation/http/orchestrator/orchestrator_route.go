package orchestrator

import (
	"fmt"
	"full_cycle_cep/pkg/core/middleware"
	"full_cycle_cep/pkg/domain/contracts/init_provider"
	"full_cycle_cep/pkg/shared/log"
	"github.com/go-chi/chi/v5"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"io/ioutil"
	"net/http"
	"sync"
)

var lock sync.Mutex
var createOrchestratorInstance CreateOrchestratorRoute

type CreateOrchestratorRoute struct {
	logger                    log.LoggerManagerInterface
	RtCepValidationMiddleware middleware.CepValidationMiddleware
	initProvider              *init_provider.InitProvider
}

func NewOrchestratorRoute(
	logger log.LoggerManagerInterface,
	initProvider *init_provider.InitProvider,
) CreateOrchestratorRoute {
	if createOrchestratorInstance == (CreateOrchestratorRoute{}) {
		lock.Lock()
		defer lock.Unlock()
		if createOrchestratorInstance == (CreateOrchestratorRoute{}) {
			createOrchestratorInstance = CreateOrchestratorRoute{
				logger:       logger,
				initProvider: initProvider,
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
	carrier := propagation.HeaderCarrier(r.Header)
	ctx := r.Context()
	ctx = otel.GetTextMapPropagator().Extract(ctx, carrier)
	tracer, _ := c.initProvider.GetProvider(ctx, "orchestrator-service", "otel_collector:4317", "ms")
	ctx, span := tracer.Start(ctx, "ini orchestration")
	defer span.End()

	url := "http://cep_weather:3500/cep/05330011" // Substitua pela URL da sua API

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		fmt.Printf("Erro ao criar a requisição: %v\n", err)
		return
	}

	otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(req.Header))

	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		fmt.Printf("Erro ao fazer a requisição: %v\n", err)
		return
	}
	defer response.Body.Close()

	respBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Printf("Erro ao ler o corpo da resposta: %v\n", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(respBody)
}
