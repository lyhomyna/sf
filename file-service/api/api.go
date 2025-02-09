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

// directory where uploaded files will be saved
var filesDirectory = filepath.Join("..", "..", "files")

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

// HandleSave is the handler for saving file. File should be form value with key 'file'.
func HandleSave(w http.ResponseWriter, req *http.Request) {
    if req.Method != http.MethodPost {
	http.Error(w, "Use POST method instead", http.StatusBadRequest)
	return
    }

    f, fh, err := req.FormFile("file")
    if err != nil {
	log.Println("Failed to read form value", err)
	writeResponse(w, "Provide a file", http.StatusBadRequest)
	return
    }
    defer f.Close()
    
    if errMsg, statusCode := saveUploadedFile(fh.Filename, f); errMsg != "" && statusCode != -1 {
	writeResponse(w, errMsg, statusCode)
	return
    }

    writeResponse(w, "File saved", http.StatusOK)
}

// saveupLoadedFile saves file into server's forlder and returns error message as string and http status code as int. If error message == "" and status code == -1 there is not errors and file uploaded successfully.
func saveUploadedFile(filename string, uploadedFile multipart.File) (string, int) {
    newFilepath := filepath.Join(filesDirectory, filename)

    if _, err := os.Stat(newFilepath); err == nil {
	log.Println("File with the same name already exists")
	return "File with the same name already exists", http.StatusBadRequest 
    }

    nf, err := os.Create(newFilepath)
    if err != nil {
	log.Println("Failed to create a file:", err)
	return "Something went wrong while a server was creating file", http.StatusInternalServerError
    }
    defer nf.Close()

    // Copy content
    _, err = io.Copy(nf, uploadedFile)
    if err != nil {
	os.Remove(newFilepath)
	log.Println("Failed to save file:", err)
	return "Something went wrong when the server was copying content from your file", http.StatusInternalServerError
    }

    return "", -1
}

func HandleDelete(w http.ResponseWriter, req *http.Request) {
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

    writeResponse(w, "File deleted", http.StatusOK)
}

// deleteFile deletes file from server and returns error message as string and http response code as int. If error message == "" and status code == -1 there is not errors and file uploaded successfully.
func deleteFile(filename string) (string, int) {
    fullFilepath := filepath.Join(filesDirectory, filename)

    if _, err := os.Stat(fullFilepath); err != nil {
	return "File doesn't exist", http.StatusNoContent 
    }

    if err := os.Remove(fullFilepath); err != nil {
	return "Error removing file", http.StatusInternalServerError 
    }

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
    w.Header().Set("Content-Type", "application/octet-stream")
    w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")

    http.ServeFile(w,req, filepath.Join(filesDirectory, filename))
}

func HandleFilenames(w http.ResponseWriter, req *http.Request) {
    if req.Method != http.MethodGet {
	http.Error(w, "Use GET method instead", http.StatusBadRequest)
	return
    }

    if filenames, err := getFilenames(); err != nil {
	writeResponse(w, err.Message, err.Code)
    } else {
	writeResponse(w, filenames, http.StatusOK)
    }
}

type responseError struct {
    Code int
    Message string
}

func getFilenames() ([]string, *responseError) {
    entries, err := os.ReadDir(filesDirectory)
    if err != nil {
	return nil, &responseError {
	    Code: http.StatusInternalServerError,
	    Message: "Cannot read from user directory",
	}
    }

    filenames := []string {}
    for _, entry := range entries {
	filenames = append(filenames, entry.Name())
    }

    return filenames, nil
}


func writeResponse(w http.ResponseWriter, data any, code int) {
    w.Header().Set("Content-Type", "application/json")
    w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")

    w.WriteHeader(code)
    d := struct {
	Data any `json:"data"`
    } {
	Data: data,
    }

    response, err := json.Marshal(d)
    if err != nil {
	panic(err)
    }
    w.Write(response)
}

// logConnection is a 'logger' of each request
func logConnection(req *http.Request) {
    log.Printf("%s | %s %s\n", req.RemoteAddr, req.Method,  req.URL.Path)
}
