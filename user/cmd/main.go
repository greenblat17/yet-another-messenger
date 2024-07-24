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
		fmt.Fprintln(w, "Hello from User Service")
	})

	http.HandleFunc("/live", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "User Service is alive")
	})

	http.HandleFunc("/ready", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "User Service is ready")
	})

	http.HandleFunc("/start-up", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "User Service is start up")
	})

	port, ok := os.LookupEnv("USER_PORT")
	if !ok {
		port = "8081"
	}

	http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
}
