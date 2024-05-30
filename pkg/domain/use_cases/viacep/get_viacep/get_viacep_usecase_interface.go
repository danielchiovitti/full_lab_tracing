package get_viacep

import "full_cycle_cep/pkg/domain/models"

type GetViaCepUseCaseInterface interface {
	Execute(cep string) (*models.CepGetResponse, error)
}
