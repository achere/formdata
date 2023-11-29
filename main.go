package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/form", func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.Body)
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "beb", http.StatusBadRequest)
		}
		fmt.Println(body)
		bodyString := string(body)
		fmt.Println(bodyString)
		r.ParseMultipartForm(10000)
		for key, value := range r.Form {
			log.Println(key, value)
		}
		log.Println(r.Header)
	})

	http.ListenAndServe(":8080", mux)
}
