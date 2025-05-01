package userImages

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"

	"github.com/lyhomyna/sf/file-service/repository"
	"github.com/lyhomyna/sf/file-service/utils"
)
type UserImagesService struct {
    repository repository.UserImagesRepository
}

// directory where uploaded files will be saved
var userImagesDirectoryPath = filepath.Join("userImages")

func NewUserImagesService(userImagesRepository repository.UserImagesRepository) *UserImagesService {
    return &UserImagesService{
	repository: userImagesRepository,
    }
}

func (uis *UserImagesService) GetUserImageHandler(w http.ResponseWriter, req *http.Request) {
    if req.Method != http.MethodPost {
	http.Error(w, "Use POST method instead", http.StatusBadRequest)
	return
    }
    _, httpErr := utils.CheckAuth(req)
    if httpErr != nil {
	utils.WriteResponse(w, httpErr.Message, httpErr.Code)
	return
    }

    // get the avatar from form field by "avatar" name
    avaFile, _, err := req.FormFile("avatar")
    if err != nil {
	utils.WriteResponse(w, "Couldn't get image from request", http.StatusBadRequest)
	return
    }
    defer avaFile.Close()

    // validate MIME type
    userImageExt, isImage := validateImageFile(avaFile)
    if !isImage {
	utils.WriteResponse(w, "Unsupported image type. Possible types is either PNG or JPEG", http.StatusUnsupportedMediaType)
	return
    }
    avaFile.Seek(0, io.SeekStart)

    if err := initUserImageDir(); err != nil {
	log.Println("For some reason, couldn't initialize user images folder:", err.Error())
	utils.WriteResponse(w, "Internal server error X_X", http.StatusInternalServerError)
	return
    }

    // construct new avatar name (uuid + .ext) 
    userImageFilepath := filepath.Join(userImagesDirectoryPath, constructFilename(userImageExt)) 
    outFile, err := os.Create(userImageFilepath)
    if err != nil {
	log.Println("WTF! Couldn't create a file to store userImage:", err.Error())
	utils.WriteResponse(w, "Internal server error X_X", http.StatusInternalServerError)
	return
    }
    defer outFile.Close()

    writtenBytes, err := io.Copy(outFile, avaFile)
    if err != nil {
	utils.WriteResponse(w, "Internal server error X_X", http.StatusInternalServerError)
	return
    }

    responseMgs := fmt.Sprintf("User image has been stored. Written %d bytes", writtenBytes)
    utils.WriteResponse(w, responseMgs, http.StatusOK)
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
	log.Println("Error reading user image file from request:", err.Error())
	return "", false
    }
    
    contentType := http.DetectContentType(buffer)
    if !allowed[contentType] {
	return "", false
    }

    ext := strings.Split(contentType, "/")[1]
    return ext, true
}

func initUserImageDir() error {
    err := os.MkdirAll(userImagesDirectoryPath, os.ModePerm)
    return err
}

func constructFilename(extension string) string {
    newFilename := fmt.Sprintf("%s.%s", uuid.NewString(), extension)
    return newFilename 
}
func (uid *UserImagesService) SaveUserImageHandler(w http.ResponseWriter, req *http.Request) {

}
