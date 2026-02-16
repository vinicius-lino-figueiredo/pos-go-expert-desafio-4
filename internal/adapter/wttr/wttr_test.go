package wttr_test

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/vinicius-lino-figueiredo/pos-go-expert-desafio-4/internal/adapter/wttr"
)

type roundTripperFunc func(*http.Request) (*http.Response, error)

func (f roundTripperFunc) RoundTrip(r *http.Request) (*http.Response, error) {
	return f(r)
}

type WttrSuite struct {
	suite.Suite
}

func TestWttrSuite(t *testing.T) {
	suite.Run(t, new(WttrSuite))
}

func (s *WttrSuite) newGetter(rt roundTripperFunc) *http.Client {
	return &http.Client{Transport: rt}
}

func (s *WttrSuite) TestFailedRequest() {
	expectedErr := errors.New("request failed")
	rt := roundTripperFunc(func(r *http.Request) (*http.Response, error) {
		return nil, expectedErr
	})

	tg := wttr.NewTemperatureGetter(s.newGetter(rt))
	temp, err := tg.GetTemperature(context.Background(), "São Paulo")

	s.ErrorIs(err, expectedErr)
	s.Zero(temp)
}

func (s *WttrSuite) TestNonOKStatusCode() {
	rt := roundTripperFunc(func(r *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusNotFound,
			Body:       io.NopCloser(strings.NewReader("")),
		}, nil
	})

	tg := wttr.NewTemperatureGetter(s.newGetter(rt))
	temp, err := tg.GetTemperature(context.Background(), "São Paulo")

	s.ErrorAs(err, &wttr.ErrStatusCode{})
	s.Zero(temp)
}

func (s *WttrSuite) TestInvalidJSON() {
	rt := roundTripperFunc(func(r *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader("{]")),
		}, nil
	})

	tg := wttr.NewTemperatureGetter(s.newGetter(rt))
	temp, err := tg.GetTemperature(context.Background(), "São Paulo")

	s.Error(err)
	s.Zero(temp)
}

func (s *WttrSuite) TestNoConditionFound() {
	rt := roundTripperFunc(func(r *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader(`{"current_condition":[]}`)),
		}, nil
	})

	tg := wttr.NewTemperatureGetter(s.newGetter(rt))
	temp, err := tg.GetTemperature(context.Background(), "São Paulo")

	s.ErrorIs(err, wttr.ErrNoConditionFound)
	s.Zero(temp)
}

func (s *WttrSuite) TestSuccessfulTemperature() {
	rt := roundTripperFunc(func(r *http.Request) (*http.Response, error) {
		body := `{"current_condition":[{"temp_C":"25"}]}`
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader(body)),
		}, nil
	})

	tg := wttr.NewTemperatureGetter(s.newGetter(rt))
	temp, err := tg.GetTemperature(context.Background(), "São Paulo")

	s.NoError(err)
	s.Equal(25.0, temp)
}
