package models

type ViaCepContractModel struct {
	Cep         string
	Logradouro  string
	Complemento string
	Bairro      string
	Localidade  string
	Uf          string
	Ibge        string
	Gia         string
	Ddd         string
	Siafi       string
	Erro        bool
}
