package repository

import (
	"errors"
	"io"

	"github.com/lyhomyna/sf/file-service/models"
)

var FilesErrorFailureToRetrieve = errors.New("Failure to retrieve file from DB")
var FilesErrorDbQuery = errors.New("Database query error")
var FilesErrorFileNotExist = errors.New("File doesn't exist")
var FilesErrorFileExist = errors.New("File with the save name already exists")
var FilesErrorInternal = errors.New("Internal error")
var FilesErrorFailureCreateFile = errors.New("Couldn't create new file")
var FilesErrorCopyFailure = errors.New("File couldn't be saved")
var FilesErrorDbSave = errors.New("Couldn't save file to the database")

type UserImagesRepository interface {
    SaveUserImage(userId string, imageUrl string) error
    GetUserImageUrl(userId string) (string, error)
}

type FilesRepository interface {
    SaveFile(userId string, filename string, file io.Reader) (*models.UserFile, error) 
    DeleteFile(userId string, fileId string) error 
    GetFile(userId string, fileId string) (*models.DbUserFile, error)
    GetFiles(userId string) ([]*models.DbUserFile, *models.HttpError)
}

