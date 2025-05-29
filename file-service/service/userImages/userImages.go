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
)

type UserImagesService struct {
    repository repository.UserImageRepository
}

// directory where uploaded files will be saved
var userImagesDirectoryPath = filepath.Join("userImages")

func NewUserImagesService(userImagesRepository repository.UserImageRepository) *UserImagesService {
    return &UserImagesService{
	repository: userImagesRepository,
    }
}

func (uis *UserImagesService) SaveUserImageHandler(userId string, req *http.Request) (string, *models.HttpError) {
    // get the avatar from form field by "avatar" name
    avaFile, avaFileHeader, err := req.FormFile("image")
    if err != nil {
	return "", &models.HttpError{
	    Message: "Couldn't get image from request",
	    Code: http.StatusBadRequest,
	}
    }
    defer avaFile.Close()

    // validate MIME type
    userImageActualExt, isImage := validateImageFile(avaFile)
    if !isImage {
	return "", &models.HttpError{
	    Message: "Unsupported image type. Possible types is either PNG or JPEG",
	    Code: http.StatusUnsupportedMediaType,
	}
    }
    avaFile.Seek(0, io.SeekStart)

    if err := initUserImageDir(); err != nil {
	log.Println("For some reason, couldn't initialize user images folder:", err.Error())
	return "", &models.HttpError{
	    Message: "Internal server error",
	    Code: http.StatusInternalServerError,
	}
    }

    // construct new avatar name (uuid + .ext) 
    userImageFilename := constructFilename(avaFileHeader.Filename, userImageActualExt)
    internalImageFilepath := filepath.Join(userImagesDirectoryPath, userImageFilename) 
    outFile, err := os.Create(internalImageFilepath)
    if err != nil {
	log.Println("WTF! Couldn't create a file to store userImage:", err.Error())

	return "", &models.HttpError{
	    Message: "Internal server error",
	    Code: http.StatusInternalServerError,
	}
    }
    defer outFile.Close()

    writtenBytes, err := io.Copy(outFile, avaFile)
    if err != nil {
	log.Println(err)

	return "", &models.HttpError{
	    Message: "Internal server error",
	    Code: http.StatusInternalServerError,
	}
    }
    log.Printf("User image '%s' has been stored. Written %d bytes\n", userImageFilename, writtenBytes)

    externalImageFilepath := fmt.Sprintf("image/%s", userImageFilename)
    err = uis.repository.SaveUserImage(userId, externalImageFilepath)
    if err != nil {
	removeFile(internalImageFilepath)

	return "", &models.HttpError{
	    Message: err.Error(),
	    Code: http.StatusInternalServerError,
	}
    }
    log.Printf("User image '%s' metadata has been stored.\n", userImageFilename)

    return externalImageFilepath, nil
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
	log.Printf("Urgent remove of the file with id '%s'", path)
    }
}


func (uid *UserImagesService) GetUserImageHandler(path string) (*models.ImageData, *models.HttpError) {
    pathChunks := strings.Split(path, "/")

    if len(pathChunks) != 2 {
	return nil, &models.HttpError{
	    Message: "Incorrect image path", 
	    Code: http.StatusBadRequest,
	}
    } else if (pathChunks[0] != "image") {
	return nil, &models.HttpError{
	    Message: "Incorrect image path", 
	    Code: http.StatusBadRequest,
	}
    }

    imageFile, contentTypeChunk, err := readImage(filepath.Join(userImagesDirectoryPath, pathChunks[1]))
    if err != nil {
	log.Println(err)

	return nil, &models.HttpError{
	    Message: errors.Unwrap(err).Error(), 
	    Code: http.StatusInternalServerError,
	}
    }

    return &models.ImageData {
	ImageFile: imageFile,
	ContentTypeChunk: contentTypeChunk,
    }, nil 
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
