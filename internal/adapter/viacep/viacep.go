// Package viacep TODO
package viacep

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/vinicius-lino-figueiredo/pos-go-expert-desafio-4/domain"
)

const baseURL = "http://viacep.com.br/ws/"

// ErrStatusCode TODO
type ErrStatusCode struct {
	Status int
}

// Error implements [error].
func (e ErrStatusCode) Error() string {
	return fmt.Sprintf("unexpected status code %d", e.Status)
}

// ViaCEP TODO
type ViaCEP struct {
	cl *http.Client
}

// NewAddressGetter returns a new implementation of [domain.AddressGetter].
func NewAddressGetter(cl *http.Client) domain.AddressGetter {
	return &ViaCEP{cl: cl}
}

// GetAddress implements [domain.AddressGetter].
func (a *ViaCEP) GetAddress(ctx context.Context, postalCode string) (string, error) {
	u, err := a.getURL(postalCode)
	if err != nil {
		return "", fmt.Errorf("mounting url: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return "", fmt.Errorf("creating request: %w", err)
	}

	res, err := a.cl.Do(req)
	if err != nil {
		return "", fmt.Errorf("doing request: %w", err)
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return "", ErrStatusCode{Status: res.StatusCode}
	}

	var body viaCEP
	if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
		return "", fmt.Errorf("decoding response: %w", err)
	}

	return body.Localidade, nil
}

func (a *ViaCEP) getURL(postalCode string) (string, error) {
	return url.JoinPath(baseURL, postalCode, "json")
}

type viaCEP struct {
	Cep         string `json:"cep"`
	Logradouro  string `json:"logradouro"`
	Complemento string `json:"complemento"`
	Unidade     string `json:"unidade"`
	Bairro      string `json:"bairro"`
	Localidade  string `json:"localidade"`
	Uf          string `json:"uf"`
	Estado      string `json:"estado"`
	Regiao      string `json:"regiao"`
	Ibge        string `json:"ibge"`
	Gia         string `json:"gia"`
	Ddd         string `json:"ddd"`
	Siafi       string `json:"siafi"`
}
