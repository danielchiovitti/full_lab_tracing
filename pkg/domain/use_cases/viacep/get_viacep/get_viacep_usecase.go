package get_viacep

import (
	"full_cycle_cep/pkg/domain/contracts/viacep"
	"full_cycle_cep/pkg/domain/contracts/weatherapi"
	"full_cycle_cep/pkg/domain/models"
	"full_cycle_cep/pkg/shared/business_error"
	"full_cycle_cep/pkg/shared/constants"
	"full_cycle_cep/pkg/shared/log"
	"sync"
)

var lock sync.Mutex
var getViaCepUseCaseInstance *GetViaCepUseCase

type GetViaCepUseCase struct {
	logger             log.LoggerManagerInterface
	viaCepContract     viacep.ViaCepContractInterface
	weatherApiContract weatherapi.WeatherApiContractInterface
}

func NewGetViaCepUseCase(
	logger log.LoggerManagerInterface,
	viaCepContract viacep.ViaCepContractInterface,
	weatherApiContract weatherapi.WeatherApiContractInterface,
) *GetViaCepUseCase {
	if getViaCepUseCaseInstance == nil {
		lock.Lock()
		defer lock.Unlock()
		if getViaCepUseCaseInstance == nil {
			getViaCepUseCaseInstance = &GetViaCepUseCase{
				logger:             logger,
				viaCepContract:     viaCepContract,
				weatherApiContract: weatherApiContract,
			}
		}
	}
	return getViaCepUseCaseInstance
}

func (g *GetViaCepUseCase) Execute(cep string) (*models.CepGetResponse, error) {
	resCep, err := g.viaCepContract.GetCep(cep)

	if err != nil {
		return nil, err
	}

	if resCep.Cep == "" {
		e := &business_error.BusinessError{
			Message: string(constants.Cep_not_found),
		}
		return nil, e
	}

	resWeather, err := g.weatherApiContract.GetWeatherInfo(resCep.Localidade)
	if err != nil {
		return nil, err
	}

	return &models.CepGetResponse{
		Temp_C: resWeather.Current.TempC,
		Temp_F: (resWeather.Current.TempC * 1.8) + 32,
		Temp_K: resWeather.Current.TempC + 273,
	}, nil
}
