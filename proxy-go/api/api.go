package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"proxy-go/lib"
	"strings"
)

type Response struct {
 	EncryptHost string `json:"encrypted_host"`
}

func encryptSubdomainHandler(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	token := strings.TrimPrefix(authHeader, "Bearer ")
	if token != "YOUR_ACCESS_TOKEN" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	ip := r.URL.Query().Get("ip")
	port := r.URL.Query().Get("port")
	host := fmt.Sprintf("%s:%s", ip, port)
	encryptedResult, _ := lib.EncryptSubdomain(host)
	response := Response{EncryptHost: encryptedResult}
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "JSON encoding error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}


func EncryptSubdomainApi() {

	http.HandleFunc("/encrypt/subdomain", encryptSubdomainHandler)
	
	apiServerHost := "0.0.0.0:8080"
	
	fmt.Printf("api server for encrypting host start on %s ...\n", apiServerHost)
	
	fmt.Printf("[GET] http://%s/encrypt/subdomain?ip=<ip>&port=<port> \n", apiServerHost)
	
	err := http.ListenAndServe(apiServerHost, nil)
	
	if err != nil {
		fmt.Println("Server error:", err)
	}
}
