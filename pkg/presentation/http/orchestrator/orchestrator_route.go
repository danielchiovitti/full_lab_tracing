package orchestrator

import (
	"full_cycle_cep/pkg/shared/log"
	"sync"
)

var lock sync.Mutex
var createOrchestratorInstance *CreateOrchestratorRoute

type CreateOrchestratorRoute struct {
	logger log.LoggerManagerInterface
}

func NewOrchestratorRoute(
	logger log.LoggerManagerInterface,
) *CreateOrchestratorRoute {
	if createOrchestratorInstance == nil {
		lock.Lock()
		defer lock.Unlock()
		if createOrchestratorInstance == nil {
			createOrchestratorInstance = &CreateOrchestratorRoute{
				logger: logger,
			}
		}
	}
	return createOrchestratorInstance
}
