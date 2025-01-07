package main

import (
	"log"
	"net/http"
	"path/filepath"
)

var filesDirectory = filepath.Join("..", "files")

func main() {
    fs := http.FileServer(http.Dir(filesDirectory))
    http.Handle("/", fs)
    
    if err := http.ListenAndServe(":8080", nil); err != nil {
	log.Println("Failed to start server:", err)
    }
}
