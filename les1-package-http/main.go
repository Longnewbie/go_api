package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/demo", demoFunc)

	log.Println("Server is starting...")
	err := http.ListenAndServe(":8080", nil) // init server - localhost:8080

	if err != nil {
		log.Fatal("Server error:", err)
	}
}

func demoFunc(w http.ResponseWriter, r *http.Request) {
	log.Printf("%+v", r)

	// check method
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	res := map[string]string{
		"message": "Chào mừng Bro đến với Go API",
		"info":    "Golang API",
	}

	w.Header().Set("Content-Type", "application/json") // default json
	w.Header().Set("X-Course", "Golang")

	// Cách 1 sử dụng Marshal
	// data, err := json.Marshal(res)
	// if err != nil {
	// 	http.Error(w, "Internal server error", http.StatusInternalServerError)
	// 	return
	// }

	// w.Write(data)

	// Cách 2 sử dụng Encoder trên ResponseWriter (code ngắn hơn)
	json.NewEncoder(w).Encode(res)
}
