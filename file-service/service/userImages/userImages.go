package userImages

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/lyhomyna/sf/file-service/models"
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

func (uis *UserImagesService) SaveUserImageHandler(w http.ResponseWriter, req *http.Request) {
    if req.Method != http.MethodPost {
	http.Error(w, "Use POST method instead", http.StatusBadRequest)
	return
    }
    userId, httpErr := utils.CheckAuth(req)
    if httpErr != nil {
	utils.WriteResponse(w, httpErr.Message, httpErr.Code)
	return
    }

    // get the avatar from form field by "avatar" name
    avaFile,avaFileHeader, err := req.FormFile("image")
    if err != nil {
	utils.WriteResponse(w, "Couldn't get image from request", http.StatusBadRequest)
	return
    }
    defer avaFile.Close()

    // validate MIME type
    userImageActualExt, isImage := validateImageFile(avaFile)
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
    userImageFilename := constructFilename(avaFileHeader.Filename, userImageActualExt)
    internalImageFilepath := filepath.Join(userImagesDirectoryPath, userImageFilename) 
    outFile, err := os.Create(internalImageFilepath)
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
    log.Printf("User image '%s' has been stored. Written %d bytes\n", userImageFilename, writtenBytes)

    externalImageFilepath := fmt.Sprintf("image/%s", userImageFilename)
    err = uis.repository.SaveUserImage(userId, externalImageFilepath)
    if err != nil {
	removeFile(internalImageFilepath)
	utils.WriteResponse(w, err.Error(), http.StatusInternalServerError)
    }
    log.Printf("User image '%s' metadata has been stored.\n", userImageFilename)

    utils.WriteResponseV2(w, &models.ImageJson { ImageUrl: externalImageFilepath }, http.StatusOK)
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
	log.Println("Error reading user image file from request. Reason:", err)
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

func constructFilename(filename string, actualExtension string) string {
    hash := sha1.New()
    hash.Write([]byte(filename))

    newFilename := fmt.Sprintf("%s.%s", fmt.Sprintf("%x", hash.Sum(nil)), actualExtension)
    return newFilename 
}

func removeFile(path string) {
    err := os.Remove(path)
    if err == nil {
	log.Printf("File '%s' properly removed", path)
    }
}


func (uid *UserImagesService) GetUserImageHandler(w http.ResponseWriter, req *http.Request) {
    _, httpErr := utils.CheckAuth(req)
    if httpErr != nil {
	utils.WriteResponse(w, httpErr.Message, httpErr.Code)
    }

    imagePath := strings.TrimPrefix(req.URL.Path, "/")
    imagePathChunks := strings.Split(imagePath, "/")

    if len(imagePathChunks) != 2 {
	utils.WriteResponse(w, "Incorrect image path", http.StatusBadRequest)
	return
    } else if (imagePathChunks[0] != "image") {
	utils.WriteResponse(w, "Incorrect image path", http.StatusBadRequest)
	return
    }

    image, contentType, err := readImage(filepath.Join(userImagesDirectoryPath, imagePathChunks[1]))
    if err != nil {
	log.Println(err)
	utils.WriteResponse(w, errors.Unwrap(err).Error(), http.StatusInternalServerError)
	return
    }
    defer image.Close()

    log.Printf("Sending image '%s' back to the user...\n", imagePathChunks[1])

    w.Header().Set("Content-Type", http.DetectContentType(contentType))
    w.WriteHeader(http.StatusOK)
    if _, err = io.Copy(w, image); err != nil {
	log.Printf("Image '%s' wasn't sent to the user. Reason:%v", imagePathChunks[1], err)
	return
    }

    log.Printf("Image '%s' sent", imagePathChunks[1])
}

func readImage(imagePath string) (*os.File, []byte, error) {
    image, err := os.Open(imagePath)
    if err != nil {
	return nil, []byte{}, fmt.Errorf("%w: %v", errors.New("User picture missing"), err)
    }

    buf := make([]byte, 512)
    _, err = image.Read(buf)
    if err != nil {
	return nil, []byte{}, fmt.Errorf("%w: %v", errors.New("Couldn't read user image"), err)
    }
    image.Seek(0, io.SeekStart)

    return image, buf, nil
}
