package main

import (
	"encoding/json"
	"net/http"
)

func main() {
	http.HandleFunc("/users/", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]string{
			"id":   "123",
			"name": "João Mauricio",
		})
	})

	http.ListenAndServe(":8001", nil)
}
