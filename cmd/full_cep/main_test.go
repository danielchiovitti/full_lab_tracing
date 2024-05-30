package main

import (
	full_lab_cep "full_cycle_cep"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestValidCEP(t *testing.T) {
	req, err := http.NewRequest("GET", "/cep/05330011", nil)
	if err != nil {
		t.Fatal(err)
	}

	h := full_lab_cep.InitializeHandlers()
	rr := httptest.NewRecorder()
	r := chi.NewRouter()
	r.Get("/cep/{cep}", h.HnCreateCepRoute.Get)
	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestNotFoundCEP(t *testing.T) {
	req, err := http.NewRequest("GET", "/cep/05330999", nil)
	if err != nil {
		t.Fatal(err)
	}

	h := full_lab_cep.InitializeHandlers()
	rr := httptest.NewRecorder()
	r := chi.NewRouter()
	r.Get("/cep/{cep}", h.HnCreateCepRoute.Get)
	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusNotFound, rr.Code)
}

func TestInvalidCEPFormat(t *testing.T) {
	req, err := http.NewRequest("GET", "/cep/053309999", nil)
	if err != nil {
		t.Fatal(err)
	}

	h := full_lab_cep.InitializeHandlers()
	rr := httptest.NewRecorder()
	r := chi.NewRouter()
	r.With(h.HnCreateCepRoute.RtCepValidationMiddleware.Validate).Get("/cep/{cep}", h.HnCreateCepRoute.Get)
	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusUnprocessableEntity, rr.Code)
}
