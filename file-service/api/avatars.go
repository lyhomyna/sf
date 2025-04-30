package api

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
)

// directory where uploaded files will be saved
var avatarsDirectoryPath = filepath.Join("avatars")

func HandleAvatarSave(w http.ResponseWriter, req *http.Request) {
    if req.Method != http.MethodPost {
	http.Error(w, "Use POST method instead", http.StatusBadRequest)
	return
    }
    _, httpErr := checkAuth(req)
    if httpErr != nil {
	writeResponse(w, httpErr.Message, httpErr.Code)
	return
    }

    // get the avatar from form field by "avatar" name
    avaFile, _, err := req.FormFile("avatar")
    if err != nil {
	writeResponse(w, "Couldn't get avatar from request", http.StatusBadRequest)
	return
    }
    defer avaFile.Close()

    // validate MIME type
    avatarExt, isImage := validateImageFile(avaFile)
    if !isImage {
	writeResponse(w, "Unsupported image type. Possible types is either PNG or JPEG", http.StatusUnsupportedMediaType)
	return
    }
    avaFile.Seek(0, io.SeekStart)

    if err := initAvatarDir(); err != nil {
	log.Println("For some reason, couldn't initialize avatars folder:", err.Error())
	writeResponse(w, "Internal server error X_X", http.StatusInternalServerError)
	return
    }

    // construct new avatar name (uuid + .ext) 
    avatarFilepath := filepath.Join(avatarsDirectoryPath, constructFilename(avatarExt)) 
    outFile, err := os.Create(avatarFilepath)
    if err != nil {
	log.Println("WTF! Couldn't create a file to store avatar:", err.Error())
	writeResponse(w, "Internal server error X_X", http.StatusInternalServerError)
	return
    }
    defer outFile.Close()

    writtenBytes, err := io.Copy(outFile, avaFile)
    if err != nil {
	log.Println("Couldn't copy file:", err.Error())
	writeResponse(w, "Internal server error X_X", http.StatusInternalServerError)
	return
    }

    responseMgs := fmt.Sprintf("Avatar has been stored. Written %d bytes", writtenBytes)
    writeResponse(w, responseMgs, http.StatusOK)
    // save avatar to database (OMG X_X)
    // 
    // side notes
    // avatar path should be stored in the database
}

// validateImageFile returns extension and isImage indicator
func validateImageFile(possibleImage io.Reader) (string, bool) {
    allowed := map[string] bool {
	"image/jpeg": true,
	"image/png": true,
    }

    buffer := make([]byte, 512)
    _, err := possibleImage.Read(buffer)
    if err != nil {
	log.Println("Error reading avatar file from request:", err.Error())
	return "", false
    }
    
    contentType := http.DetectContentType(buffer)
    log.Println("ContentType of avatar is:", contentType)
    if !allowed[contentType] {
	return "", false
    }

    log.Println("File is an image")
    ext := strings.Split(contentType, "/")[1]
    return ext, true
}

func initAvatarDir() error {
    err := os.MkdirAll(avatarsDirectoryPath, os.ModePerm)
    log.Println("initAvatarDir() err is:", err)
    return err
}

func constructFilename(extension string) string {
    newFilename := fmt.Sprintf("%s.%s", uuid.NewString(), extension)
    return newFilename 
}

func HandleAvatarGet(w http.ResponseWriter, req *http.Request) {
    panic("Not yet implemented")
}
