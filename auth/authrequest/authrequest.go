package authrequest

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"

	"github.com/factorysh/chevillette/memory"
)

// See https://nginx.org/en/docs/http/ngx_http_auth_request_module.html

type AuthRequest struct {
	memory memory.Memory
}

func New(memory *memory.Memory) *AuthRequest {
	return &AuthRequest{
		memory: *memory,
	}
}

func (a *AuthRequest) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ip, err := getIP(r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
	}
	var status int
	log.Println("Keys :", ip, r.UserAgent())
	if a.memory.HasKey([]string{ip, r.UserAgent()}) {
		status = http.StatusOK
	} else {
		status = http.StatusForbidden
	}
	log.Printf(`%s %d "%v"`, ip, status, r.UserAgent())
	w.WriteHeader(status)
}

func (a *AuthRequest) ListenAndServe(addr string) error {
	return http.ListenAndServe(addr, a)
}

// See https://golangbyexample.com/golang-ip-address-http-request/
func getIP(r *http.Request) (string, error) {
	//Get IP from the X-REAL-IP header
	ip := r.Header.Get("X-REAL-IP")
	netIP := net.ParseIP(ip)
	if netIP != nil {
		return ip, nil
	}

	//Get IP from X-FORWARDED-FOR header
	ips := r.Header.Get("X-FORWARDED-FOR")
	splitIps := strings.Split(ips, ",")
	for _, ip := range splitIps {
		netIP := net.ParseIP(ip)
		if netIP != nil {
			return ip, nil
		}
	}

	//Get IP from RemoteAddr
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return "", err
	}
	netIP = net.ParseIP(ip)
	if netIP != nil {
		return ip, nil
	}
	return "", fmt.Errorf("No valid ip found")
}
