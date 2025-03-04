package session

import (
	"net/http"

	"github.com/lyhomyna/sf/auth-service/models"
)

func Routes(siglog models.Siglog) http.Handler {
    mux := http.NewServeMux()

    mux.HandleFunc("/register", func(w http.ResponseWriter, req *http.Request) {
	register(siglog, w, req)
    })
    mux.HandleFunc("/login", func(w http.ResponseWriter, req *http.Request) {
	login(siglog, w, req)
    })
    mux.HandleFunc("/logout", func(w http.ResponseWriter, req *http.Request) {
	logout(siglog, w, req)
    })

    return mux
}

func register(siglog models.Siglog, w http.ResponseWriter, req *http.Request) {
    panic("Not yet implemented.")
}

func login(siglog models.Siglog, w http.ResponseWriter, req *http.Request) {
    panic("Not yet implemented.")
}

func logout(siglog models.Siglog, w http.ResponseWriter, req *http.Request) {
    panic("Not yet implemented.")
}
