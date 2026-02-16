package viacep_test

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/vinicius-lino-figueiredo/pos-go-expert-desafio-4/internal/adapter/viacep"
)

type roundTripperFunc func(*http.Request) (*http.Response, error)

func (f roundTripperFunc) RoundTrip(r *http.Request) (*http.Response, error) {
	return f(r)
}

type ViaCEPSuite struct {
	suite.Suite
}

func TestViaCEPSuite(t *testing.T) {
	suite.Run(t, new(ViaCEPSuite))
}

func (s *ViaCEPSuite) newGetter(rt roundTripperFunc) *http.Client {
	return &http.Client{Transport: rt}
}

func (s *ViaCEPSuite) TestFailedRequest() {
	expectedErr := errors.New("request failed")
	rt := roundTripperFunc(func(r *http.Request) (*http.Response, error) {
		return nil, expectedErr
	})

	ag := viacep.NewAddressGetter(s.newGetter(rt))
	addr, err := ag.GetAddress(context.Background(), "01001000")

	s.ErrorIs(err, expectedErr)
	s.Empty(addr)
}

func (s *ViaCEPSuite) TestNonOKStatusCode() {
	rt := roundTripperFunc(func(r *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusNotFound,
			Body:       io.NopCloser(strings.NewReader("")),
		}, nil
	})

	ag := viacep.NewAddressGetter(s.newGetter(rt))
	addr, err := ag.GetAddress(context.Background(), "01001000")

	s.ErrorAs(err, &viacep.ErrStatusCode{})
	s.Empty(addr)
}

func (s *ViaCEPSuite) TestInvalidJSON() {
	rt := roundTripperFunc(func(r *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader("{]")),
		}, nil
	})

	ag := viacep.NewAddressGetter(s.newGetter(rt))
	addr, err := ag.GetAddress(context.Background(), "01001000")

	s.Error(err)
	s.Empty(addr)
}

func (s *ViaCEPSuite) TestSuccessfulAddress() {
	rt := roundTripperFunc(func(r *http.Request) (*http.Response, error) {
		body := `{"localidade":"São Paulo"}`
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader(body)),
		}, nil
	})

	ag := viacep.NewAddressGetter(s.newGetter(rt))
	addr, err := ag.GetAddress(context.Background(), "01001000")

	s.NoError(err)
	s.Equal("São Paulo", addr)
}
