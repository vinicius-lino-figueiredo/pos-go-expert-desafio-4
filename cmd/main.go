// Package main TODO
package main

import (
	"context"
	"crypto/tls"
	"log"
	"net/http"
	"os/signal"
	"syscall"

	"github.com/vinicius-lino-figueiredo/pos-go-expert-desafio-4/internal/adapter/handler"
	"github.com/vinicius-lino-figueiredo/pos-go-expert-desafio-4/internal/adapter/viacep"
	"github.com/vinicius-lino-figueiredo/pos-go-expert-desafio-4/internal/adapter/wttr"
)

const addr = ":8000"

func main() {
	ag := viacep.NewAddressGetter(http.DefaultClient)
	tg := wttr.NewTemperatureGetter(http.DefaultClient)

	http.DefaultClient.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	hr := handler.NewHandler(ag, tg)

	ctx, cancel := context.WithCancelCause(context.Background())

	server := http.Server{Addr: addr, Handler: hr}

	ctx, stop := signal.NotifyContext(ctx, syscall.SIGINT)
	defer stop()

	go func() {
		if err := server.ListenAndServe(); err != nil {
			cancel(err)
		}
	}()

	<-ctx.Done()

	log.Println(context.Cause(ctx))

}
