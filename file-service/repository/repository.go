package repository

import (
	"io"

	"github.com/lyhomyna/sf/file-service/models"
)

type UserImagesRepository interface {
    SaveUserImage(userId string, imageUrl string) error
    GetUserImageUrl(userId string) (string, error)
}

type FilesRepository interface {
    SaveFile(userId string, filename string, file io.Reader, dir string) (*models.UserFile, error) 
    DeleteFile(userId string, fileId string) error 
    GetFile(userId string, fileId string) (*models.DbUserFile, error)
    GetFiles(userId string) ([]*models.DbUserFile, error)
    ListDir(path, userId string) ([]models.DirEntry, error)
    CreateDir(userId, parentDirId, name string) (string, error)
    CreateRootDir(userId string) (string, error)
    GetDirIdByPath(userId, path string) (string, error)
}
