package viacep

import "full_cycle_cep/pkg/domain/models"

type ViaCepContractInterface interface {
	GetCep(cep string) (*models.ViaCepContractModel, error)
}
