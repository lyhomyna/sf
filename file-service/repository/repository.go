package repository

import (
	"errors"

	"github.com/lyhomyna/sf/file-service/models"
)

var FilesErrorFailureToRetrieve = errors.New("Failure to retrieve file from DB")
var FilesErrorDbQuery = errors.New("Database query error")
var FilesErrorFileNotExist = errors.New("File doesn't exist")

type UserImagesRepository interface {
    SaveUserImage(userId string, imageUrl string) error
    GetUserImageUrl(userId string) (string, error)
}

type FilesRepository interface {
    SaveFile(userId string, filename string, fileBytes []byte) (*models.UserFile, *models.HttpError) 
    DeleteFile(userId string, fileId string) error 
    GetFile(userId string, fileId string) (*models.DbUserFile, error)
    GetFiles(userId string) ([]*models.DbUserFile, *models.HttpError)
}

