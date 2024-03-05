package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/log", func(w http.ResponseWriter, r *http.Request) {
		body := make([]byte, 1<<20)
		log.Println("Request header: ", r.Header)
		if _, err := r.Body.Read(body); err == nil {
			fmt.Fprintf(w, "not ok %v", err)
		} else {
			log.Println(string(body))
			fmt.Fprintf(w, "ok")
		}
	})
	log.Println("now listening on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
