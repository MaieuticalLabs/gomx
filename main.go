package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net"
	"net/http"
	"strings"
)

type CheckDomain struct {
	Domain string `json:"domain"`
}

type CheckResponse struct {
	Status bool `json:"status"`
}

func checkHandler(response http.ResponseWriter, request *http.Request) {
	var status bool
	var domain CheckDomain

	switch request.Method {
	case http.MethodPost:
		decoder := json.NewDecoder(request.Body)
		decoder.DisallowUnknownFields()
		err := decoder.Decode(&domain)
		if err != nil {
			http.Error(response, "Invalid JSON payload", http.StatusBadRequest)
			return
		}
		if strings.Index(domain.Domain, ".") < 1 {
			http.Error(response, "Invalid Domain", http.StatusBadRequest)
			return
		}
		fmt.Printf("Got Domain: %v\n", domain.Domain)
		mxrecords, err := net.LookupMX(domain.Domain)
		if err != nil {
			fmt.Printf("No MX record found for %v: %v\n", domain.Domain, err)
			status = false
		} else {
			status = len(mxrecords) > 0
		}
		encoder := json.NewEncoder(response)
		r := CheckResponse{status}
		err = encoder.Encode(r)
		if err != nil {
			http.Error(response, "Oops", http.StatusInternalServerError)
		}
	default:
		http.Error(response, "405 Method not allowed", http.StatusMethodNotAllowed)
	}
}

func main() {
	var address string
	flag.StringVar(&address, "address", "127.0.0.1:8000", "Address to bind to")
	flag.Parse()
	http.HandleFunc("/api/v1/check", checkHandler)
	fmt.Printf("Service listening on: %v\n", address)
	err := http.ListenAndServe(address, nil)
	if err != nil {
		fmt.Printf("%v\n", err)
	}
}
