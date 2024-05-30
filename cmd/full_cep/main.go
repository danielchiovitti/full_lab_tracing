package main

import (
	"context"
	"fmt"
	full_lab_cep "full_cycle_cep"
	"go.opentelemetry.io/otel"
	"log"
	"net/http"
	"os"
	"os/signal"
)

func main() {
	p := full_lab_cep.InitializeProvider()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	shutdown, err := p.GetProvider("", "")
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := shutdown(ctx); err != nil {
			log.Fatal("failed to shutdown TracerProvider: %w", err)
		}
	}()

	tracer := otel.Tracer("microservice-tracer")

	//templateData := TemplateDat

	h := full_lab_cep.InitializeHandlers()
	r := h.GetRoutes()
	fmt.Println("Starting")
	http.ListenAndServe(":3500", r)
}
