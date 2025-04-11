package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// directory where uploaded files will be saved
var filesDirectory = filepath.Join("files")

var sessionCookieName = "session-id"
var authServiceBaseUrl = "http://auth-service:8081"

func OptionsMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
	// logConnection(req)

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

    sessionCookie, err := req.Cookie(sessionCookieName)
    if err != nil {
	log.Println(err)
	writeResponse(w, "Session cookie missing", http.StatusUnauthorized)
	return
    }

    userId, err := verifySession(sessionCookie)
    if err != nil {
	writeResponse(w, err.Error(), http.StatusUnauthorized)
	return
    }

    f, fh, err := req.FormFile("file")
    if err != nil {
	log.Println("Failed to read form value", err)
	writeResponse(w, "Provide a file", http.StatusBadRequest)
	return
    }
    defer f.Close()
    
    if errMsg, statusCode := saveUploadedFile(userId, fh.Filename, f); errMsg != "" && statusCode != -1 {
	writeResponse(w, errMsg, statusCode)
	return
    }

    writeResponse(w, "File saved", http.StatusOK)
}

// saveupLoadedFile saves file into server's forlder and returns error message as string and http status code as int. If error message == "" and status code == -1 there is not errors and file uploaded successfully.
func saveUploadedFile(userId string, filename string, uploadedFile io.Reader) (string, int) {
    newFilepath := filepath.Join(filesDirectory, userId, filename)

    if _, err := os.Stat(newFilepath); err == nil {
	log.Println("File with the same name already exists")
	return "File with the same name already exists", http.StatusBadRequest 
    }

    if err := os.MkdirAll(filepath.Dir(newFilepath), os.ModePerm); err != nil {
	return "Something went wrong", http.StatusInternalServerError
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

    sessionCookie, err := req.Cookie(sessionCookieName)
    if err != nil {
	log.Println(err)
	writeResponse(w, "Session cookie missing", http.StatusUnauthorized)
	return
    }

    userId, err := verifySession(sessionCookie)
    if err != nil {
	writeResponse(w, err.Error(), http.StatusUnauthorized)
	return
    }

    filename := strings.TrimPrefix(req.URL.Path, "/delete/")

    errMsg, statusCode := deleteFile(userId, filename)
    if errMsg != "" && statusCode != -1 {
	writeResponse(w, errMsg, statusCode)
	return
    }

    writeResponse(w, "File deleted", http.StatusOK)
}

// deleteFile deletes file from server and returns error message as string and http response code as int. If error message == "" and status code == -1 there is not errors and file uploaded successfully.
func deleteFile(userId string, filename string) (string, int) {
    fullFilepath := filepath.Join(filesDirectory, userId, filename)

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
    if req.Method != http.MethodGet {
	http.Error(w, "Use GET method instead", http.StatusBadRequest)
	return
    }

    // get session cookie
    sessionCookie, err := req.Cookie(sessionCookieName)
    if err != nil {
	writeResponse(w, "Session cookie missing", http.StatusUnauthorized)
	return
    }
    // verify session cookie and get user id
    userId, err := verifySession(sessionCookie)
    if err != nil {
	writeResponse(w, err.Error(), http.StatusUnauthorized)
	return
    }

    filename := strings.TrimPrefix(req.URL.Path, "/download/")   
    if filename == "" {
	writeResponse(w, "File not specified", http.StatusBadRequest)
	return
    }
    filepathToDownload := filepath.Join(filesDirectory, userId, filename) 

    _, err = os.Stat(filepathToDownload)
    if err != nil {
	writeResponse(w, "File not found", http.StatusNotFound)
	return
    }

    w.Header().Set("Content-Disposition", "attachment; filename="+filename)
    w.Header().Set("Content-Type", "application/octet-stream")
    // w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")

    http.ServeFile(w,req, filepathToDownload)
}

func HandleFilenames(w http.ResponseWriter, req *http.Request) {
    if req.Method != http.MethodGet {
	http.Error(w, "Use GET method instead", http.StatusBadRequest)
	return
    }

    // get session cookie
    sessionCookie, err := req.Cookie(sessionCookieName)
    if err != nil {
	writeResponse(w, "Session cookie missing", http.StatusUnauthorized)
	return
    }
    // verify session cookie and get user id
    userId, err := verifySession(sessionCookie)
    if err != nil {
	writeResponse(w, err.Error(), http.StatusUnauthorized)
	return
    }

    if filenames, err := getFilenames(userId); err != nil {
	writeResponse(w, err.Message, err.Code)
    } else {
	writeResponse(w, filenames, http.StatusOK)
    }
}

func verifySession(sessionCookie *http.Cookie) (string, error) {
    reqUrl := fmt.Sprintf("%s/check-auth", authServiceBaseUrl)
    
    req, _ := http.NewRequest("GET", reqUrl, nil)

    req.AddCookie(sessionCookie)

    client := &http.Client{}

    resp, err := client.Do(req)
    if err != nil {
	log.Println("Unable to verify session:", err)
	return "", errors.New("Unable to verify session")
    }

    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
	log.Println("Unable to verify session: Status", resp.StatusCode)
	return "", errors.New("Unable to verify session")
    }

    var res struct {
	Id string `json:"userId"`
    }
    json.NewDecoder(resp.Body).Decode(&res)

    return res.Id, nil
}

type responseError struct {
    Code int
    Message string
}

func getFilenames(userId string) ([]string, *responseError) {
    filesPath := filepath.Join(filesDirectory, userId)
    entries, err := os.ReadDir(filesPath)
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
