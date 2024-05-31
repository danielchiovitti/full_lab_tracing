package main

import (
	"context"
	"fmt"
	full_lab_cep "full_cycle_cep"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	p := full_lab_cep.InitializeProvider()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	tracer, err := p.GetProvider(ctx, "cep-service", "otel_collector:4317", "micro")
	if err != nil {
		log.Fatal(err)
	}

	ctx = context.Background()
	ctx, span := tracer.Start(ctx, "iniciando")
	span.End()

	go func() {
		h := full_lab_cep.InitializeHandlers()
		r := h.GetRoutes()
		fmt.Println("Starting")
		http.ListenAndServe(":3500", r)
	}()

	//go func() {
	select {
	case <-sigCh:
		log.Println("shutting down gracefully")
	case <-ctx.Done():
		log.Println("shutting down due to other reason")
	}
	//}()

	_, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	//tracer := otel.Tracer("microservice-tracer")

	//templateData := TemplateDat

}
