package main

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/fire", fireHandler)
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

func fireHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	token := string(body)
	resp, err := sendRequest(token)

	if err != nil {
		fmt.Println("Error performing request:", err)
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
}

func sendRequest(token string) (*http.Response, error) {
	apiURL := "https://emea2.owndata.com/api/v1/services/1104801/gdpr/forget"

	formData := map[string]string{
		"table_name": "Contact",
		"record_id":  "0033z000033pldQAAQ",
	}

	body, contentType, err := createMultipartBody(formData)
	if err != nil {
		fmt.Println("Error creating request body:", err)
		return nil, err
	}

	req, err := http.NewRequest("POST", apiURL, body)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return nil, err
	}

	req.Header.Set("Content-Type", contentType)

	req.Header.Set("Authorization", "Bearer "+token)

	// Perform the request
	client := http.Client{}
	return client.Do(req)
}

func createMultipartBody(formData map[string]string) (*bytes.Buffer, string, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	for key, value := range formData {
		err := writer.WriteField(key, value)
		if err != nil {
			return nil, "", err
		}
	}

	err := writer.Close()
	if err != nil {
		return nil, "", err
	}

	contentType := writer.FormDataContentType()

	return body, contentType, nil
}
