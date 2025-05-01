package middleware

import (
	"log"
	"net/http"
)

func OptionsMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
	logConnection(req)

	if req.Method == http.MethodOptions {
	    w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	    w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS")
	    w.WriteHeader(http.StatusNoContent)
	    return
	}
	    
	next.ServeHTTP(w, req)
    })
}

// logConnection is a 'logger' of each request
func logConnection(req *http.Request) {
    log.Printf("%s | %s %s\n", req.RemoteAddr, req.Method,  req.URL.Path)
}
