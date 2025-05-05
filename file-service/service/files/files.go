package files

import (
	"errors"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/lyhomyna/sf/file-service/models"
	"github.com/lyhomyna/sf/file-service/utils"

	"github.com/lyhomyna/sf/file-service/repository"
)

type FilesService struct {
    repository repository.FilesRepository
}

// Directory, where uploaded files will be saved
var filesDirectory = filepath.Join("files")

func NewFilesService(filesRepository repository.FilesRepository) *FilesService {
    return &FilesService{
	repository: filesRepository,
    }
}

// SaveHandler is the handler for saving file. File should be form value with key 'file'.
func (fs *FilesService)SaveHandler(w http.ResponseWriter, req *http.Request) {
    if req.Method != http.MethodPost {
	http.Error(w, "Use POST method instead", http.StatusMethodNotAllowed)
	return
    }

    userId, httpErr := utils.CheckAuth(req)
    if httpErr != nil {
	utils.WriteResponse(w, httpErr.Message, httpErr.Code)
	return
    }

    if !strings.Contains(req.Header.Get("Content-Type"), "multipart/form-data") {
	utils.WriteResponse(w, "Expected multipart/form-data", http.StatusUnsupportedMediaType)
	return
    }

    // Read file
    f, fh, err := req.FormFile("file")
    if err != nil {
	log.Println("Failed to read form file:", err)
	utils.WriteResponse(w, "Provide a file", http.StatusBadRequest)
	return
    }
    defer f.Close()

    //fileBytes, err := io.ReadAll(f)
    //if err != nil {
    //	log.Println("Failed to read from file:", err)
    //	utils.WriteResponse(w, "Failed to read from file", http.StatusBadRequest)
    //	return
    //}

    if fh.Filename == "" {
	    utils.WriteResponse(w, "Filename is required", http.StatusBadRequest)
	    return
    }

    userFile, err := fs.repository.SaveFile(userId, fh.Filename, f)
    if err != nil {
	log.Println(err)

	switch {
	    case errors.Is(repository.FilesErrorFileExist, err):
		utils.WriteResponse(w, "File with the same name already exists", http.StatusConflict)

	    case errors.Is(repository.FilesErrorInternal, err):
		utils.WriteResponse(w, "Internal server error", http.StatusInternalServerError)

	    case errors.Is(repository.FilesErrorFailureCreateFile, err):
		utils.WriteResponse(w, "Couldn't save file", http.StatusInternalServerError)

	    case errors.Is(repository.FilesErrorCopyFailure, err):
		utils.WriteResponse(w, "Couldn't save file", http.StatusInternalServerError)

	    case errors.Is(repository.FilesErrorDbSave, err):
		utils.WriteResponse(w, "Couldn't save your file to the database", http.StatusInternalServerError)

	    default:
		utils.WriteResponse(w, "Unknown error occurred", http.StatusInternalServerError)
	}

	return
    }

    utils.WriteResponse(w, userFile, http.StatusOK)
}

func (fs *FilesService) DeleteHandler(w http.ResponseWriter, req *http.Request) {
    if req.Method != http.MethodDelete {
	http.Error(w, "Use DELETE method instead", http.StatusMethodNotAllowed)
	return
    }
    userId, httpErr := utils.CheckAuth(req)
    if httpErr != nil {
	utils.WriteResponse(w, httpErr.Message, httpErr.Code)
	return
    }

    fileId := strings.TrimPrefix(req.URL.Path, "/delete/")

    err := fs.repository.DeleteFile(userId, fileId)
    if err != nil {
	log.Println(err.Error())

	if errors.Is(repository.FilesErrorFailureToRetrieve, err) || errors.Is(repository.FilesErrorFileNotExist, err) {
	    utils.WriteResponse(w, "There is no file to delete", http.StatusNotFound)
	    return
	}

	if errors.Is(repository.FilesErrorDbQuery, err) {
	    utils.WriteResponse(w, "Internal server error", http.StatusInternalServerError)
	    return
	}
    }

    utils.WriteResponse(w, "File deleted", http.StatusOK)
}

// DownloadHandler downloads file into client downloads folder
func (fs *FilesService) DownloadHandler(w http.ResponseWriter, req *http.Request) {
    if req.Method != http.MethodGet {
	http.Error(w, "Use GET method instead", http.StatusMethodNotAllowed)
	return
    }
    userId, httpErr := utils.CheckAuth(req)
    if httpErr != nil {
	utils.WriteResponse(w, httpErr.Message, httpErr.Code)
	return
    }

    fileId := strings.TrimPrefix(req.URL.Path, "/download/")   
    if fileId == "" {
	utils.WriteResponse(w, "File not specified", http.StatusBadRequest)
	return
    }

    file, err := fs.repository.GetFile(userId, fileId)
    if err != nil {
	log.Println(err.Error())
	if errors.Is(repository.FilesErrorFailureToRetrieve, err) {
	    utils.WriteResponse(w, "Couldn't retrieve file", http.StatusNotFound)
	    return
	}
    }

    _, err = os.Stat(file.Filepath)
    if err != nil {
	utils.WriteResponse(w, "File not found", http.StatusNotFound)
	return
    }

    w.Header().Set("Content-Disposition", "attachment; filename="+file.Filename)
    w.Header().Set("Content-Type", "application/octet-stream")

    http.ServeFile(w, req, file.Filepath)
}

func (fs *FilesService) FilesHanlder(w http.ResponseWriter, req *http.Request) {
    if req.Method != http.MethodGet {
	http.Error(w, "Use GET method instead", http.StatusMethodNotAllowed)
	return
    }
    userId, httpErr := utils.CheckAuth(req)
    if httpErr != nil {
	utils.WriteResponse(w, httpErr.Message, httpErr.Code)
	return
    }

    userFiles, httpErr := fs.repository.GetFiles(userId)
    if httpErr != nil {
	utils.WriteResponse(w, httpErr.Message, httpErr.Code)
	return
    }

    responseUserFiles := []*models.UserFile{}

    for _, userFile := range userFiles {
	responseUserFiles = append(responseUserFiles, &models.UserFile{
	    Id: userFile.Id,
	    Filename: userFile.Filename,
	    LastAccessed: userFile.LastAccessed.Unix(),
	})
    }

    utils.WriteResponse(w, responseUserFiles, http.StatusOK)
}
