package repository

import (
	"github.com/lyhomyna/sf/file-service/models"
)

type UserImagesRepository interface {
    SaveUserImage(userId string, imageUrl string) error
    GetUserImageUrl(userId string) (string, error)
}

type FilesRepository interface {
    SaveFile(userId string, filename string, fileBytes []byte) *models.HttpError 
    DeleteFile(userId string, filename string) *models.HttpError 
    GetFilenames(userId string) ([]string, *models.HttpError)
}
