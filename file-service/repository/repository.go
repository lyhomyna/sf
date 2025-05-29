package repository

import (
	"io"

	"github.com/lyhomyna/sf/file-service/models"
)

type UserImageRepository interface {
    SaveUserImage(userId string, imageUrl string) error
    GetUserImageUrl(userId string) (string, error)
}

type FileRepository interface {
    SaveFile(userId, filename, dirPath, dirId string, file io.Reader) (*models.UserFile, error) 
    DeleteFile(userId string, fileId string) error 
    GetFile(userId string, fileId string) (*models.DbUserFile, error)
}

type DirRepository interface {
    ListDir(path, userId string) ([]models.DirEntry, error)
    CreateDir(userId, parentDirId, parentDirPath, name string) (string, error)
    DeleteDir(userId, dirId string) (error)
    CreateRootDir(userId string) (string, error)
    GetDirIdByPath(path, userId string) (string, error)
}
