package subnet

import (
	"net"
	"net/http"

	"github.com/yury-kuznetsov/shortener/cmd/config"
)

// Handle is a function that takes a http.HandlerFunc as an argument and returns a modified http.HandlerFunc.
// The modified handler checks whether the request is trusted or not using the checkTrust function.
// If the request is not trusted, it returns a "Access denied" error with HTTP status code 403.
// Otherwise, it calls the original handler.
func Handle(handler http.HandlerFunc) http.HandlerFunc {
	handlerFunc := func(res http.ResponseWriter, req *http.Request) {
		if !checkTrust(req) {
			http.Error(res, "Access denied", http.StatusForbidden)
			return
		}

		handler(res, req)
	}

	return handlerFunc
}

func checkTrust(req *http.Request) bool {
	ip := net.ParseIP(req.Header.Get("X-Real-IP"))
	if ip == nil {
		return false
	}

	if config.Options.TrustedNet == "" {
		return false
	}

	_, trustNet, _ := net.ParseCIDR(config.Options.TrustedNet)

	return trustNet.Contains(ip)
}
