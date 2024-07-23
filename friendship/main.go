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
		fmt.Fprintln(w, "Hello from Friendship Service")
	})

	http.HandleFunc("/live", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "Friendship Service is alive")
	})

	http.HandleFunc("/ready", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "Friendship Service is ready")
	})

	port, ok := os.LookupEnv("FRIENDSHIP_PORT")
	if !ok {
		port = "8082"
	}

	http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
}
