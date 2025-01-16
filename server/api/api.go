package api

import (
	"encoding/json"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// files directory where an uploaded files will be saved
var (
    filesDirectory = filepath.Join("..", "..", "..", "files")
    fileServer = http.FileServer(http.Dir(filesDirectory))
)

func HandleSave(w http.ResponseWriter, req *http.Request) {
    logConnection(req)

    if req.Method != http.MethodPost {
	http.Error(w, "Use POST method instead", http.StatusBadRequest)
	return
    }

    f, fh, err := req.FormFile("file")
    if err != nil {
	log.Println("Failed to read form value.", err)
	writeResponse(w, "Provide a file.", http.StatusBadRequest)
	return
    }
    defer f.Close()
    
    if errMsg, statusCode := saveUploadedFile(fh.Filename, f); errMsg != "" && statusCode != -1 {
	writeResponse(w, errMsg, statusCode)
	return
    }

    // Write response
    w.Header().Set("Content-Type", "application/json")

    writeResponse(w, "file saved", http.StatusOK)
}

// saveupLoadedFile saves file into server's forlder and returns error message as string and http status code as int. If error message == "" and status code == -1 there is not errors and file uploaded successfully.
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
func HandleDelete(w http.ResponseWriter, req *http.Request) {
    logConnection(req)

    if req.Method != http.MethodDelete {
	http.Error(w, "Use DELETE method instead", http.StatusBadRequest)
	return
    }

    filename := strings.TrimPrefix(req.URL.Path, "/delete/")

    errMsg, statusCode := deleteFile(filename)
    if errMsg != "" && statusCode != -1 {
	writeResponse(w, errMsg, statusCode)
	return
    }

    writeResponse(w, "file deleted", http.StatusOK)
}

// deleteFile deletes file from server and returns error message as string and http response code as int. If error message == "" and status code == -1 there is not errors and file uploaded successfully.
func deleteFile(filename string) (string, int) {
    fullFilepath := filepath.Join(filesDirectory, filename)

    if _, err := os.Stat(fullFilepath); err != nil {
	return "File don't exist.", http.StatusNoContent 
    }

    if err := os.Remove(fullFilepath); err != nil {
	return err.Error(), http.StatusInternalServerError 
    }
    
    log.Printf("File %s deleted successfully.", filename)

    return "", -1
}

// handleDownload downloads file into client downloads folder
func HandleDownload(w http.ResponseWriter, req *http.Request) {
    logConnection(req)

    if req.Method != http.MethodGet {
	http.Error(w, "Use GET method instead", http.StatusBadRequest)
	return
    }

    filename := strings.TrimPrefix(req.URL.Path, "/download/")   
    if filename == "" {
	writeResponse(w, "File not specified", http.StatusBadRequest)
	return
    }

    _, err := os.Stat(filepath.Join(filesDirectory, filename))
    if err != nil {
	writeResponse(w, "File not found", http.StatusNotFound)
	return
    }

    w.Header().Set("Content-Disposition", "attachment; filename="+filename)
    http.StripPrefix("/download/", fileServer).ServeHTTP(w, req)
}

type responseMessage struct {
    Message string `json:"message"`
}

func writeResponse(w http.ResponseWriter, message string, code int) {
    w.WriteHeader(code)
    response, _ := json.Marshal(responseMessage{Message: message})
    w.Write(response)
}


// logConnection is a 'logger' of each request
func logConnection(req *http.Request) {
    log.Printf("%s | %s %s\n", req.RemoteAddr, req.Method,  req.URL.Path)
}

