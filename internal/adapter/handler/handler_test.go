package handler_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/vinicius-lino-figueiredo/pos-go-expert-desafio-4/internal/adapter/handler"
)

type mockAddressGetter struct {
	address string
	err     error
}

func (m *mockAddressGetter) GetAddress(_ context.Context, _ string) (string, error) {
	return m.address, m.err
}

type mockTemperatureGetter struct {
	temp float64
	err  error
}

func (m *mockTemperatureGetter) GetTemperature(_ context.Context, _ string) (float64, error) {
	return m.temp, m.err
}

type HandlerSuite struct {
	suite.Suite
}

func TestHandlerSuite(t *testing.T) {
	suite.Run(t, new(HandlerSuite))
}

func (s *HandlerSuite) TestInvalidPostalCodeTooShort() {
	h := handler.NewHandler(&mockAddressGetter{}, &mockTemperatureGetter{})
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/temperature/123", nil)

	h.ServeHTTP(rec, req)

	s.Equal(http.StatusBadRequest, rec.Code)
}

func (s *HandlerSuite) TestInvalidPostalCodeNonNumeric() {
	h := handler.NewHandler(&mockAddressGetter{}, &mockTemperatureGetter{})
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/temperature/abcdefgh", nil)

	h.ServeHTTP(rec, req)

	s.Equal(http.StatusBadRequest, rec.Code)
}

func (s *HandlerSuite) TestAddressGetterError() {
	ag := &mockAddressGetter{err: errors.New("not found")}
	h := handler.NewHandler(ag, &mockTemperatureGetter{})
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/temperature/01001000", nil)

	h.ServeHTTP(rec, req)

	s.Equal(http.StatusInternalServerError, rec.Code)
}

func (s *HandlerSuite) TestTemperatureGetterError() {
	ag := &mockAddressGetter{address: "São Paulo"}
	tg := &mockTemperatureGetter{err: errors.New("service unavailable")}
	h := handler.NewHandler(ag, tg)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/temperature/01001000", nil)

	h.ServeHTTP(rec, req)

	s.Equal(http.StatusInternalServerError, rec.Code)
}

func (s *HandlerSuite) TestSuccessfulResponse() {
	ag := &mockAddressGetter{address: "São Paulo"}
	tg := &mockTemperatureGetter{temp: 25.0}
	h := handler.NewHandler(ag, tg)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/temperature/01001000", nil)

	h.ServeHTTP(rec, req)

	s.Equal(http.StatusOK, rec.Code)

	var resp handler.Response
	err := json.NewDecoder(rec.Body).Decode(&resp)
	s.NoError(err)
	s.Equal(25.0, resp.TempC)
	s.Equal(25.0*1.8+32, resp.TempF)
	s.Equal(25.0+273, resp.TempK)
}
