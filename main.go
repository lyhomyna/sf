package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

type ResponseMessage struct {
    Message string `json:"message"`
}

var filesDirectory = filepath.Join("..", "files")

func main() {
    http.HandleFunc("/save", saveFile) 

    fs := http.FileServer(http.Dir(filesDirectory))
    http.Handle("/", fs)
    
    log.Println("Server is running...")
    if err := http.ListenAndServe(":8080", nil); err != nil {
	log.Println("Failed to start server:", err)
    }
}


func saveFile(w http.ResponseWriter, req *http.Request) {
    // Get file from form
    f, fh, err := req.FormFile("file")
    if err != nil {
	log.Fatalln("Failed to read form value.", err)
    }
    defer f.Close()
    
    // TODO: encrypt filename

    // Create new file
    newFilepath := filepath.Join(filesDirectory, fh.Filename)
    nf, err := os.Create(newFilepath)
    if err != nil {
	f.Close()
	log.Fatalln("Failed to create a file:", err)
    }
    defer nf.Close()

    // Copy content
    _, err = io.Copy(nf, f)
    if err != nil {
	f.Close(); nf.Close()
	os.Remove(newFilepath)
	log.Fatalln("Failed to save file:", err)
    }
    log.Printf("File %s saved.\n", newFilepath)

    // Write response
    w.Header().Set("Content-Type", "application/json")
    w.Header().Set("Date", time.Now().UTC().String())

    response, _ := json.Marshal(ResponseMessage{Message: "file saved"})
    w.Write(response)
}

func downloadFile() {

}

func deleteFile() {

}
