package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"time"
)

type Config struct {
	Port string
}

type GreetingService interface {
	Greet(ctx context.Context, firstName, lastName string) (string, error)
	Goodbye(ctx context.Context, firstName, lastName string) (string, error)
}

type EnglishGreetingService struct{}

func (EnglishGreetingService) Greet(ctx context.Context, firstName, lastName string) (string, error) {
	if lastName == "" {
		return fmt.Sprintf("Hello, %s!", firstName), nil
	}
	return fmt.Sprintf("Hello, %s %s!", firstName, lastName), nil
}

func (EnglishGreetingService) Goodbye(ctx context.Context, firstName, lastName string) (string, error) {
	if lastName == "" {
		return fmt.Sprintf("Goodbye, %s!", firstName), nil
	}
	return fmt.Sprintf("Goodbye, %s %s!", firstName, lastName), nil
}

type Handler struct {
	service GreetingService
}

func (h *Handler) Greet(w http.ResponseWriter, r *http.Request) {
	firstName := r.URL.Query().Get("first_name")
	lastName := r.URL.Query().Get("last_name")

	if firstName == "" {
		firstName = "World"
	}

	msg, err := h.service.Greet(r.Context(), firstName, lastName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintln(w, msg)
}

func (h *Handler) Goodbye(w http.ResponseWriter, r *http.Request) {
	firstName := r.URL.Query().Get("first_name")
	lastName := r.URL.Query().Get("last_name")

	if firstName == "" {
		firstName = "World"
	}

	msg, err := h.service.Goodbye(r.Context(), firstName, lastName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintln(w, msg)
}

func loadConfig() Config {
	var cfg Config

	flag.StringVar(&cfg.Port, "port", "8080", "HTTP server port")
	flag.Parse()

	return cfg
}

func main() {
	cfg := loadConfig()

	go func() {
		for {
			fmt.Println("ok")
			time.Sleep(5 * time.Second)
		}
	}()

	handler := &Handler{
		service: EnglishGreetingService{},
	}

	http.HandleFunc("/greet", handler.Greet)
	http.HandleFunc("/goodbye", handler.Goodbye)

	addr := ":" + cfg.Port

	fmt.Printf("Listening on %s\n", addr)

	if err := http.ListenAndServe(addr, nil); err != nil {
		panic(err)
	}
}