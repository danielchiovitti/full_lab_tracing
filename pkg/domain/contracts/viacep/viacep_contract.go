package viacep

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"full_cycle_cep/pkg/domain/models"
	"full_cycle_cep/pkg/shared/log"
	"io"
	"net/http"
	"sync"
)

var lock sync.Mutex
var viaCepContractInstance *ViaCepContract

type ViaCepContract struct {
	logger log.LoggerManagerInterface
}

func NewViaCepContract(logger log.LoggerManagerInterface) *ViaCepContract {
	if viaCepContractInstance == nil {
		lock.Lock()
		defer lock.Unlock()
		if viaCepContractInstance == nil {
			viaCepContractInstance = &ViaCepContract{
				logger: logger,
			}
		}
	}
	return viaCepContractInstance
}

func (v *ViaCepContract) GetCep(cep string) (*models.ViaCepContractModel, error) {
	var cModel models.ViaCepContractModel
	url := fmt.Sprintf("https://viacep.com.br/ws/%s/json/", cep)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

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

	err = json.Unmarshal(body, &cModel)
	if err != nil {
		fmt.Println("Error unmarshal:", err)
		return nil, err
	}

	return &cModel, nil
}
