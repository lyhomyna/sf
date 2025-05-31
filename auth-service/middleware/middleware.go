package middleware

import (
	"log"
	"net/http"
	"os"
)

func OptionsMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
	logConnection(req)

	if req.Method == http.MethodOptions {
	    w.Header().Set("Access-Control-Allow-Origin", "http://localhost")
	    w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS")
	    w.WriteHeader(http.StatusNoContent)
	    return
	}

	if req.Method == http.MethodPost && req.URL.Path == "/register" {
	    log.Println(">> REGISTER REQUEST RECEIVED. CHECKING CREDENTIALS...")
	    username, password, ok := req.BasicAuth()
	    if !ok || username != os.Getenv("ADMIN_USERNAME") || password != os.Getenv("ADMIN_PASSWORD") {
		log.Println(">> REGISTER CREDENTIALS INCORRECT. DENY")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	    }
	    log.Println(">> REGISTER CREDENTIALS IS CORRECT. PROCEED")
	}

	next.ServeHTTP(w, req)
    })
}

// logConnection is a 'logger' of each request
func logConnection(req *http.Request) {
    log.Printf("| %s %s | %s\n", req.Method, req.URL.Path, req.RemoteAddr)
}
