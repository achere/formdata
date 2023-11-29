package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/form", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("START REQUEST")
		fmt.Printf("Method %v\n", r.Method)
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "beb", http.StatusBadRequest)
		}
		bodyString := string(body)
		fmt.Println(bodyString)
		r.Body.Close()
		r.Body = io.NopCloser(bytes.NewReader(body))
		fmt.Println("BODY STRING FINISH, KEY VALUE BEGIN")
		r.ParseMultipartForm(10 << 20)
		for key, values := range r.PostForm {
			for _, value := range values {
				fmt.Printf("%s=%s\n", key, value)
			}
		}

		fmt.Println("KEY VALUE FINISH, HEADERS BEGIN")
		for key, values := range r.Header {
			for _, value := range values {
				fmt.Printf("%s: %s\n", key, value)
			}
		}
		fmt.Println("END REQUEST")
	})

	http.ListenAndServe(":8080", mux)
}
