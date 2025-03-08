package main

import (
	"log"
	"net/http"

	"github.com/lyhomyna/sf/file-service/api"
)

func main() {
    mux := http.NewServeMux()

    mux.HandleFunc("/save", api.HandleSave)           // POST
    mux.HandleFunc("/delete/", api.HandleDelete)      // DELETE
    mux.HandleFunc("/download/", api.HandleDownload)  // GET
    mux.HandleFunc("/filenames", api.HandleFilenames) // GET

    mux.Handle("/favion.ico", http.NotFoundHandler())

    handler := api.OptionsMiddleware(mux)

    log.Println("Port 8082. Server is running...")
    if err := http.ListenAndServe(":8082", handler); err != nil {
	log.Println("Failed to start server:", err)
    }
}
