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
}

func main() {
	var port int
	flag.IntVar(&port, "port", 8000, "Port to bind")
	flag.Parse()
	port_str := fmt.Sprintf(":%d", port)
	http.HandleFunc("/api/v1/check", checkHandler)
	fmt.Printf("Service listening on port: %v\n", port_str)
	err := http.ListenAndServe(port_str, nil)
	if err != nil {
		fmt.Printf("%v\n", err)
	}
}
