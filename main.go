package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
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
		r.Body.Close()
		r.Body = io.NopCloser(bytes.NewReader(body))
		target, err := url.Parse("https://emea2.owndata.com/api/v1/services/1104801/gdpr/forget")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		req := &http.Request{
			Method:        r.Method,
			URL:           target,
			Proto:         r.Proto,
			ProtoMajor:    r.ProtoMajor,
			ProtoMinor:    r.ProtoMinor,
			Header:        r.Header,
			Body:          r.Body,
			ContentLength: r.ContentLength,
			Host:          target.Host,
		}
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			fmt.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()
		for key, values := range resp.Header {
			for _, value := range values {
				w.Header().Add(key, value)
			}
		}
		w.WriteHeader(resp.StatusCode)
		io.Copy(w, resp.Body)
	})
	http.ListenAndServe(":8080", mux)
}
