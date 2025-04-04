package main

import (
	"fmt"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(r.Header.Get("X-User-Id"))
		w.Write([]byte("Hello"))
		w.WriteHeader(http.StatusOK)
	})

	http.ListenAndServe(":8090", mux)
}
