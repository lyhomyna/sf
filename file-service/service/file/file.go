package file

import (
	"errors"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/lyhomyna/sf/file-service/models"

	"github.com/lyhomyna/sf/file-service/repository"
)

type FileService struct {
    repository repository.FileRepository
}

// Directory, where uploaded files will be saved
var filesDirectory = filepath.Join("files")

func NewFileService(fileRepository repository.FileRepository) *FileService {
    return &FileService{
	repository: fileRepository,
    }
}

// SaveHandler is the handler for saving file. File should be form value with key 'file'.
func (fs *FileService) SaveFile(userId, dirId, dirPath string, req *http.Request) (*models.UserFile, *models.HttpError) {

    // Read file
    f, fh, err := req.FormFile("file")
    if err != nil {
	log.Println("Failed to read form file:", err)
	return nil, &models.HttpError{
	    Message: "Provide a file",
	    Code: http.StatusBadRequest,
	}
    }
    defer f.Close()

    if fh.Filename == "" {
	return nil, &models.HttpError{
	    Message: "Incorrect filename",
	    Code: http.StatusBadRequest,
	}
    }

    userFile, err := fs.repository.SaveFile(userId, fh.Filename, dirPath, dirId, f)
    if err != nil {
	log.Println(err)

	switch {
	    case errors.Is(err, repository.FilesErrorFileExist):
		return nil, &models.HttpError{
		    Message: "File with the same name already exists",
		    Code: http.StatusConflict,
		}

	    case errors.Is(err, repository.FilesErrorInternal):
		return nil, &models.HttpError{
		    Message: "Internal server error",
		    Code: http.StatusInternalServerError,
		}

	    case errors.Is(err, repository.FilesErrorFailureCreateFile) || errors.Is(err, repository.FilesErrorCopyFailure):
		return nil, &models.HttpError{
		    Message: "Couldn't save file",
		    Code: http.StatusInternalServerError,
		}

	    case errors.Is(err, repository.FilesErrorDbSave):
		return nil, &models.HttpError{
		    Message: "Couldn't save file's metadata to the database",
		    Code: http.StatusInternalServerError,
		}

	    default:
		return nil, &models.HttpError{
		    Message: "WTF error occurred",
		    Code: http.StatusInternalServerError,
		}
	}
    }

    return userFile, nil
}

func (fs *FileService) DeleteFile(userId, fileId string) *models.HttpError {
    err := fs.repository.DeleteFile(userId, fileId)
    if err != nil {
	log.Println(err.Error())

	if errors.Is(repository.FilesErrorFailureToRetrieve, err) || errors.Is(repository.FilesErrorFileNotExist, err) {
	    return &models.HttpError{
		Message: "There is no file to delete",
		Code: http.StatusNotFound,
	    }
	}

	if errors.Is(repository.FilesErrorDbQuery, err) {
	    return &models.HttpError{
		Message: "There is no file to delete",
		Code: http.StatusInternalServerError,
	    }
	}
    }

    return nil
}

func (fs *FileService) GetFileEntry(userId, fileId string) (*models.FileEntry, *models.HttpError) {
    file, err := fs.repository.GetFile(userId, fileId)
    if err != nil {
	log.Println(err.Error())
	if errors.Is(repository.FilesErrorFailureToRetrieve, err) {
	    return nil, &models.HttpError{
		Message: "Couldn't retrieve file",
		Code: http.StatusNotFound,
	    }
	}
    }

    fileFullPath := filepath.Join(filesDirectory, file.Filepath)
    _, err = os.Stat(fileFullPath)
    if err != nil {
	return nil, &models.HttpError{
	    Message: "File not found",
	    Code: http.StatusNotFound,
	}
    }

    return &models.FileEntry{
	Name: file.Filename,
	Path: fileFullPath,
    }, nil
}
