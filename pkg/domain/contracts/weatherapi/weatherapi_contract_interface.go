package weatherapi

import "full_cycle_cep/pkg/domain/models"

type WeatherApiContractInterface interface {
	GetWeatherInfo(city string) (*models.WeatherApiModel, error)
}
