package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/sirupsen/logrus"
)

func main() {
	logrus.Info("Starting server...")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello from Chat Service")
	})

	http.HandleFunc("/live", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "Chat Service is alive")
	})

	http.HandleFunc("/ready", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "Chat Service is ready")
	})

	port, ok := os.LookupEnv("CHAT_PORT")
	if !ok {
		port = "8083"
	}

	http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
}
