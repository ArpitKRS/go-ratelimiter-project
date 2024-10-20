package main

import (
	"log"

	tollbooth "github.com/didip/tollbooth/v7"

	"encoding/json"
	"net/http"
)

type Message struct {
	Status string `json:"status"`
	Body   string `json:"body"`
}

func endpointHandler(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)

	message := Message{
		Status: "Successfull",
		Body:   "Hi! You've reached the API.",
	}

	err := json.NewEncoder(writer).Encode(&message)
	if err != nil {
		return
	}
}

func main() {
	message := Message{
		Status: "Request Failed",
		Body:   "The API is at capacity, try again later.",
	}

	jsonMessage, _ := json.Marshal(message)
	tlbthLimiter := tollbooth.NewLimiter(1, nil)
	tlbthLimiter.SetMessageContentType("application/json")
	tlbthLimiter.SetMessage(string(jsonMessage))
	http.Handle("/ping", tollbooth.LimitFuncHandler(tlbthLimiter, endpointHandler))

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Println("Error while listening to the port : 8080", err)
	}
}

// Command to test: Invoke-WebRequest -Uri http://localhost:8080/ping
// Commnad for too many API calls: for ($i = 1; $i -le 6; $i++) { curl.exe http://localhost:8080/ping }
