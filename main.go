package main

import (
	"encoding/json"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

type ResponseMessage struct {
    Message string `json:"message"`
}

// files directory where an uploaded files will be saved
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
	log.Println("Failed to read form value.", err)
	writeErrorResponse(w, "Provide a file.", http.StatusBadRequest)

	return
    }
    defer f.Close()
    
    if errMsg, statusCode := saveUploadedFile(fh.Filename, f); errMsg != "" && statusCode != -1 {
	writeErrorResponse(w, errMsg, statusCode)
	return
    }

    // Write response
    w.Header().Set("Content-Type", "application/json")

    response, _ := json.Marshal(ResponseMessage{Message: "file saved"})
    w.Write(response)
}

func deleteFile() {
    panic("Not yet implemented.")
}

func downloadFile() {
    panic("Not yet implemented.")
}

// writeErrorRespose created to write error responses fast
func writeErrorResponse(w http.ResponseWriter, message string, code int) {
    w.WriteHeader(code)
    response, _ := json.Marshal(ResponseMessage{Message: message})
    w.Write(response)
}

// saveupLoadedFile saves file into server's forlder and returns error message as string and staus code as int. If error message == "" and status code == -1 there is not errors and file uploaded successfully.
func saveUploadedFile(filename string, uploadedFile multipart.File) (string, int) {
    newFilepath := filepath.Join(filesDirectory, filename)

    if _, err := os.Stat(newFilepath); err == nil {
	log.Println("File with the same name already exist")
	return "File with the same name already exist.", http.StatusBadRequest 
    }

    nf, err := os.Create(newFilepath)
    if err != nil {
	log.Println("Failed to create a file:", err)
	return "Something went wrong while a server was creating file.", http.StatusInternalServerError
    }
    defer nf.Close()

    // Copy content
    _, err = io.Copy(nf, uploadedFile)
    if err != nil {
	os.Remove(newFilepath)
	log.Println("Failed to save file:", err)
	return "Something went wrong when the server was copying content from your file.", http.StatusInternalServerError
    }

    log.Printf("File %s saved.\n", newFilepath)

    return "", -1
}
