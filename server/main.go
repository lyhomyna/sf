package main

import (
	"log"
	"net/http"

	"github.com/lyhomyna/sf/server/api"
)

func main() {
    http.HandleFunc("/save", api.HandleSave)           // POST
    http.HandleFunc("/delete/", api.HandleDelete)      // DELETE
    http.HandleFunc("/download/", api.HandleDownload)  // GET
    http.HandleFunc("/filenames", api.HandleFilenames) // GET

    http.Handle("/favion.ico", http.NotFoundHandler())
    log.Println("Port 8080. Server is running...")
    if err := http.ListenAndServe(":8080", nil); err != nil {
	log.Println("Failed to start server:", err)
    }
}
