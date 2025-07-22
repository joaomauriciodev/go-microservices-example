package main

import (
	"encoding/json"
	"net/http"
)

func main() {
	http.HandleFunc("/products/", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]string{
			"id":    "1928",
			"name":  "Notebook",
			"price": "3999.99",
		})
	})

	http.ListenAndServe(":8002", nil)
}
