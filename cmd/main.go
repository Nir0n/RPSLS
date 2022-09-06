package main

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Nir0n/RPSLS/internals"
)

func main() {
	log.Println("Start app")

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("Port is not configurated")
	}

	shutdown_ctx, shutdown_cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer shutdown_cancel()

	request_ctx, request_cancel := context.WithTimeout(context.Background(), 50*time.Second)
	defer request_cancel()

	client := http.Client{}

	signal_ch := make(chan os.Signal, 1)
	signal.Notify(signal_ch, os.Interrupt, syscall.SIGTERM)

	http.HandleFunc("/choices", func(w http.ResponseWriter, r *http.Request) {
		choices := map[string][]interface{}{
			"id":   {1, 2, 3, 4, 5},
			"name": {"rock", "paper", "scissors", "lizard", "spock"},
		}
		json_choices, err := json.Marshal(choices)
		if err != nil {
			log.Fatalf("Error while  marshal json %s", err)
		}
		w.WriteHeader(http.StatusOK)
		w.Write(json_choices)
	})
	http.HandleFunc("/choice", func(w http.ResponseWriter, r *http.Request) {
		req, _ := http.NewRequestWithContext(request_ctx, http.MethodGet, "https://codechallenge.boohma.com/random", nil)
		response, ok := client.Do(req)
		if ok != nil {
			log.Fatalf("Error while requesting random number %s", ok)
		}
		defer response.Body.Close()
		body, err := io.ReadAll(response.Body)
		if err != nil {
			log.Fatalf("Error while extracting body of request %s", err)
		}
		var jsonRes map[string]interface{} // declaring a map for key names as string and values as interface
		_ = json.Unmarshal(body, &jsonRes)
		choice_id := internals.GenerateRandomChoice(int(jsonRes["random_number"].(float64)))
		res := map[string]interface{}{
			"id":   choice_id,
			"name": internals.Choices[choice_id],
		}
		json_res, err := json.Marshal(res)
		if err != nil {
			log.Fatalf("Error while  marshal json %s", err)
		}
		w.WriteHeader(http.StatusOK)
		w.Write(json_res)
	})
	http.HandleFunc("/play", func(w http.ResponseWriter, r *http.Request) {

		req, _ := http.NewRequestWithContext(request_ctx, http.MethodGet, "https://codechallenge.boohma.com/random", nil)
		response, ok := client.Do(req)
		if ok != nil {
			log.Fatalf("Error while requesting random number %s", ok)
		}
		defer response.Body.Close()
		body, err := io.ReadAll(response.Body)
		if err != nil {
			log.Fatalf("Error while extracting body of request %s", err)
		}
		// choice := internals.GenerateRandomChoice()
		w.WriteHeader(http.StatusOK)
		w.Write(body)
	})

	server := http.Server{Addr: ":" + port}
	go server.ListenAndServe()

	<-signal_ch
	shutdown_err := server.Shutdown(shutdown_ctx)
	if shutdown_err != nil {
		log.Fatalf("Something went wrong during shutdown: %s", shutdown_err)
	}
	log.Println("Finish app")
}
