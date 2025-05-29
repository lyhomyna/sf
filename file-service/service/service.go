package service

import (
	"net/http"

	"github.com/lyhomyna/sf/file-service/models"
)

type Services struct {
    FileService FileService
    UserImageService UserImageService
    DirService DirService
}

type FileService interface {
    SaveFile(userId, dirId, dirPath string, req *http.Request) (*models.UserFile, *models.HttpError)
    DeleteFile(userId, fileId string) *models.HttpError
    GetFileEntry(userId, fileId string) (*models.FileEntry, *models.HttpError)
}

type UserImageService interface {
    SaveUserImageHandler(userId string, req *http.Request) (string, *models.HttpError)
    GetUserImageHandler(path string) (*models.ImageData, *models.HttpError)
}

type DirService interface {
    ListDir(path, userId string) ([]models.DirEntry, *models.HttpError)
    CreateDirectory(userId, parentDirPath, dirName string) (*models.Dir, *models.HttpError)
    DeleteDirectory(path, userId string) *models.HttpError
    CreateRootDirectory(userId string) (string, *models.HttpError)
    GetDirIdByPath(path, userId string) (string, *models.HttpError)
}
