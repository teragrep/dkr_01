package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(os.Stdout, "[TestServer] Received request: ", r.URL.Path)

		for hk, hv := range r.Header {
			fmt.Fprintf(os.Stdout, "[TestServer] %s=%s\n", hk, hv)
		}
		fmt.Fprintf(w, "echo %v", r.URL.Path)
	})

	var port string
	if port = os.Getenv("HTTP_PORT"); port == "" {
		port = "9000"
	}

	log.Printf("[TestServer] Start server on %v %v", "localhost", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}
