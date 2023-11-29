package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/form", func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.Body)
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "beb", http.StatusBadRequest)
		}
		fmt.Println(body)
		bodyString := string(body)
		fmt.Println(bodyString)
		fmt.Println("BODY STRING FINISH, KEY VALUE BEGIN")
		r.ParseMultipartForm(10000)
		for key, value := range r.Form {
			log.Println(key, value)
		}
		fmt.Println("KEY VALUE FINISH, HEADERS BEGIN")
		for key, values := range r.Header {
			for _, value := range values {
				fmt.Printf("%s: %s\n", key, value)
			}
		}
	})

	http.ListenAndServe(":8080", mux)
}
