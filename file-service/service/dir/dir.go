package dir

import (
	"errors"
	"log"
	"net/http"
	"path"
	"strings"

	"github.com/lyhomyna/sf/file-service/config"
	"github.com/lyhomyna/sf/file-service/models"

	"github.com/lyhomyna/sf/file-service/repository"
)

type DirService struct {
    repository repository.DirRepository
    cfg *config.Config
}

func NewDirService(dirRepository repository.DirRepository) *DirService {
    return &DirService{
	repository: dirRepository,
    }
}

// Dir comes from request URL path
func (ds *DirService) ListDir(path string, userId string) ([]models.DirEntry, *models.HttpError) {
    entries, err := ds.repository.ListDir(path, userId)
    if err != nil {
	log.Println(err)
	
	if errors.Is(repository.ErrorDirectoryNotFound, err) {
	    return nil, &models.HttpError{
		Message: "Directory not found",
		Code: http.StatusNotFound,
	    }
	}

	return nil, &models.HttpError{
	    Message: "Something went wrong. Try again",
	    Code: http.StatusInternalServerError,
	}
    }

    return entries, nil
}

func (ds *DirService) CreateDirectory(userId, parentDirPath, dirName string) (*models.Dir, *models.HttpError) {
    dirId, httpErr := ds.createDir(userId, parentDirPath, dirName)
    if httpErr != nil {
	return nil, httpErr
    }

    return &models.Dir {
	Id: dirId,
	Name: dirName,
	FullPath: path.Clean("/" + strings.Trim(parentDirPath, "/") + "/" + dirName),
    }, nil
}

func (ds *DirService) createDir(userId, parentDirPath, dirName string) (string, *models.HttpError) {
	parentDirId, err := ds.repository.GetDirIdByPath(parentDirPath, userId)
	if err != nil {
	    log.Println(err)

	    // TODO: combine these two errors 
	    if errors.Is(repository.ErrorRootDirNotFound, err) || errors.Is(repository.ErrorDirectoryNotFound, err) {
		return "", &models.HttpError {
		    Message: "Parent directory not found",
		    Code: http.StatusNoContent,
		}
	    }

	    return "", &models.HttpError {
		Message: "Something went wrong. Try again",
		Code: http.StatusInternalServerError,
	    }
	}

	log.Printf(">>> TRYING TO CREATE DIR WITH PARENT DIR '%s' WITH ID '%s'", parentDirPath, parentDirId)
	dirId, err := ds.repository.CreateDir(userId, parentDirId, parentDirPath, dirName)
	if err != nil {
	    log.Println(err)
	
	    if errors.Is(repository.ErrorDirectoryAlreadyExist, err) {
		return "", &models.HttpError {
		    Message: "Directory already exist",
		    Code: http.StatusConflict,
		}
	    }

	    return "", &models.HttpError {
		Message: "Couldn't create folder. Try again",
		Code: http.StatusInternalServerError,
	    }
	}

	return dirId, nil
}

func (ds *DirService) DeleteDirectory(path, userId string) *models.HttpError {
    dirId := extractIdFromPath(path)
    err := ds.repository.DeleteDir(userId, dirId)
    if err != nil {
	log.Println(err)

	switch {
	    case errors.Is(err, repository.ErrorDelete):
		return &models.HttpError{ 
		    Message: "Could not delete directory. Try again later",
		    Code: http.StatusInternalServerError,
		}

	    default:
		return &models.HttpError{ 
		    Message: "Internal server error",
		    Code:  http.StatusInternalServerError,
		}
	}
    }

    return nil
}

func extractIdFromPath(path string) string {
    urlChanks := strings.Split(path, "/")
    return urlChanks[len(urlChanks) - 1]   
}

func (ds *DirService) CreateRootDirectory(userId string) (string, *models.HttpError) {
    dirId, httpErr := ds.createRootDir(userId)
    if httpErr != nil {
	return "", &models.HttpError{
	    Message: httpErr.Message,
	    Code: httpErr.Code,
	}
    }

    return dirId, nil
}


func (ds *DirService) createRootDir(userId string) (string, *models.HttpError) {
	dirId, err := ds.repository.CreateRootDir(userId)
	if err != nil {
	    log.Println(err)
	
	    if errors.Is(repository.ErrorDirectoryAlreadyExist, err) {
		return "", &models.HttpError {
		    Message: "Directory already exist",
		    Code: http.StatusConflict,
		}
	    }

	    return "", &models.HttpError {
		Message: "Couldn't create folder. Try again",
		Code: http.StatusInternalServerError,
	    }
	}

	return dirId, nil
}

func (ds *DirService) GetDirIdByPath(path, userId string) (string, *models.HttpError) {
    dirId, err := ds.repository.GetDirIdByPath(path, userId)
    if err != nil {
	log.Println(err)

	switch {
	    case errors.Is(err, repository.ErrorRootDirNotFound) || errors.Is(err, repository.ErrorDirectoryAlreadyExist): 
		return "", &models.HttpError {
		    Message: "Directory not found",
		    Code: http.StatusNotFound,
		}
	    case errors.Is(err, repository.ErrorGetDirFailed):
		return "", &models.HttpError {
		    Message: "Internal server error",
		    Code: http.StatusInternalServerError,
		}
	    default:
		return "", &models.HttpError {
		    Message: "WTF error occured",
		    Code: http.StatusInternalServerError,
		}
	}
    }

    return dirId, nil
}
