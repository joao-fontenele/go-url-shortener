package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type statusResponse struct {
	Running bool `json:"running"`
}

func statusHandler(w http.ResponseWriter, r *http.Request) {
	data := statusResponse{Running: true}
	js, err := json.Marshal(data)
	fmt.Println("GET /status")
	if err != nil {
		fmt.Println("err")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

// Init sets up the server and it's routes
func Init() {
	http.HandleFunc("/status", statusHandler)
	log.Println("listening on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
