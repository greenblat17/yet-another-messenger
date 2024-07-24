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
		fmt.Fprintln(w, "Hello from Notification Service")
	})

	http.HandleFunc("/live", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "Notification Service is alive")
	})

	http.HandleFunc("/ready", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "Notification Service is ready")
	})

	port, ok := os.LookupEnv("NOTIFICATION_PORT")
	if !ok {
		port = "8084"
	}

	http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
}
