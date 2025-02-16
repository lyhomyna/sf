package api

import "net/http"

func NewHttpServer() http.Handler {
    mux := http.NewServeMux()
    
    mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte("All's OK."))
    })

    return mux
}
