package weatherapi

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"full_cycle_cep/pkg/domain/models"
	"io"
	"net/http"
	"net/url"
	"sync"
)

var lock sync.Mutex
var weatherApiContractInstance WeatherApiContract

type WeatherApiContract struct{}

func NewWeatherApiContract() *WeatherApiContract {
	if weatherApiContractInstance == (WeatherApiContract{}) {
		lock.Lock()
		defer lock.Unlock()
		if weatherApiContractInstance == (WeatherApiContract{}) {
			weatherApiContractInstance = WeatherApiContract{}
		}
	}
	return &weatherApiContractInstance
}

func (w *WeatherApiContract) GetWeatherInfo(city string) (*models.WeatherApiModel, error) {
	var wModel models.WeatherApiModel
	//url := fmt.Sprintf("http://api.weatherapi.com/v1/current.json?key=cfdb3f4abe324925b07214804242205&q=%s&aqi=no", city)
	baseUrl := "http://api.weatherapi.com/v1/current.json"
	apiKey := "cfdb3f4abe324925b07214804242205"
	location := city
	aqi := "no"

	parsedURL, err := url.Parse(baseUrl)

	query := parsedURL.Query()
	query.Set("key", apiKey)
	query.Set("q", location)
	query.Set("aqi", aqi)
	parsedURL.RawQuery = query.Encode()

	finalURL := parsedURL.String()

	if err != nil {
		fmt.Printf("Erro ao analisar a URL base: %v\n", err)
		return nil, err
	}

	req, err := http.NewRequest("GET", finalURL, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return nil, err
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{Transport: tr}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error creating client:", err)
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading body:", err)
		return nil, err
	}

	fmt.Println(string(body))

	var jsonCheck map[string]interface{}
	if err := json.Unmarshal([]byte(body), &jsonCheck); err != nil {
		return nil, fmt.Errorf("response is not valid JSON: %v", err)
	}

	err = json.Unmarshal(body, &wModel)
	if err != nil {
		fmt.Println("Error unmarshal:", err)
		return nil, err
	}

	return &wModel, nil
}
