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
		port = "8080"
	}

	shutdown_ctx, shutdown_cancel := context.WithTimeout(context.Background(), 150*time.Second)
	defer shutdown_cancel()
	request_ctx, request_cancel := context.WithTimeout(context.Background(), 300*time.Second)
	defer request_cancel()

	client := http.Client{}

	signal_ch := make(chan os.Signal, 1)
	signal.Notify(signal_ch, os.Interrupt, syscall.SIGTERM)

	server := http.Server{Addr: ":" + port}
	go server.ListenAndServe()

	http.HandleFunc("/choices", func(w http.ResponseWriter, r *http.Request) {
		choices := map[string][]interface{}{
			"id":   {1, 2, 3, 4, 5},
			"name": {"rock", "paper", "scissors", "lizard", "spock"},
		}
		json_choices, err := json.Marshal(choices)
		if err != nil {
			log.Printf("Error while  marshal json %s", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write(json_choices)
	})

	http.HandleFunc("/choice", func(w http.ResponseWriter, r *http.Request) {
		var choices = map[int]string{1: "rock", 2: "paper", 3: "scissors", 4: "lizard", 5: "spock"}
		req, _ := http.NewRequestWithContext(request_ctx, http.MethodGet, "https://codechallenge.boohma.com/random", nil)
		response, ok := client.Do(req)
		if ok != nil {
			log.Printf("Error while requesting random number %s", ok)
			w.WriteHeader(http.StatusFailedDependency)
			return
		}
		defer response.Body.Close()
		body, err := io.ReadAll(response.Body)
		if err != nil {
			log.Printf("Error while extracting body of request %s", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		var jsonRes map[string]interface{}
		json.Unmarshal(body, &jsonRes)
		choice_id := internals.AnnounceRandomChoice(int(jsonRes["random_number"].(float64)))
		res := map[string]interface{}{
			"id":   choice_id,
			"name": choices[choice_id],
		}
		json_res, err := json.Marshal(res)
		if err != nil {
			log.Printf("Error while  marshal json %s", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write(json_res)
	})

	http.HandleFunc("/play", func(w http.ResponseWriter, r *http.Request) {

		initial_body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Printf("Error while extracting body of request %s", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		var json_r map[string]interface{}
		json.Unmarshal(initial_body, &json_r)
		player_choice_id := int(json_r["player"].(float64))

		req, _ := http.NewRequestWithContext(request_ctx, http.MethodGet, "https://codechallenge.boohma.com/random", nil)
		response, ok := client.Do(req)
		if ok != nil {
			log.Printf("Error while requesting random number %s", ok)
			w.WriteHeader(http.StatusFailedDependency)
			return
		}
		defer response.Body.Close()
		body, err := io.ReadAll(response.Body)
		if err != nil {
			log.Printf("Error while extracting body of request %s", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		var json_response map[string]interface{}
		_ = json.Unmarshal(body, &json_response)
		computer_choice_id := internals.AnnounceRandomChoice(int(json_response["random_number"].(float64)))
		outcome := internals.ResultCalculator(player_choice_id, computer_choice_id)
		res := map[string]interface{}{
			"results":  outcome,
			"player":   player_choice_id,
			"computer": computer_choice_id,
		}
		json_res, err := json.Marshal(res)
		if err != nil {
			log.Printf("Error while  marshal json %s", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write(json_res)
	})

	<-signal_ch
	shutdown_err := server.Shutdown(shutdown_ctx)
	if shutdown_err != nil {
		log.Fatalf("Something went wrong during shutdown: %s", shutdown_err)
	}
	log.Println("Finish app")
}
