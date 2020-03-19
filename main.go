package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"net"
	"net/http"
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
	buf := new(bytes.Buffer)
	buf.ReadFrom(request.Body)
	json.Unmarshal(buf.Bytes(), &domain)
	fmt.Printf("Got Domain: %v\n", domain.Domain)
	_, err := net.LookupHost(domain.Domain)
	if err != nil {
		fmt.Printf("Host not valid: %v\n", err)
		status = false
	} else {
		mxrecords, err := net.LookupMX(domain.Domain)
		if err != nil {
			fmt.Printf("No MX record found for %v: %v\n", domain.Domain, err)
			status = false
		} else {
			status = len(mxrecords) > 0
		}
	}
	r := CheckResponse{status}
	data, _ := json.Marshal(r)
	buf = bytes.NewBuffer(data)
	response.Write(buf.Bytes())
}

func main() {
	var port int
	flag.IntVar(&port, "port", 8000, "Port to bind")
	flag.Parse()
	port_str := fmt.Sprintf(":%d", port)
	http.HandleFunc("/api/v1/check", checkHandler)
	fmt.Printf("Service listening on port: %v\n", port_str)
	http.ListenAndServe(port_str, nil)
}
