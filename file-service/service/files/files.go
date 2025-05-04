package files

import (
	"errors"
	"io"
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
	http.Error(w, "Use POST method instead", http.StatusBadRequest)
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
	log.Println("Failed to read form value", err)
	utils.WriteResponse(w, "Provide a file", http.StatusBadRequest)
	return
    }
    defer f.Close()

    fileBytes, err := io.ReadAll(f)
    if err != nil {
	log.Println("Failed to read file", err)
	utils.WriteResponse(w, "Failed to read file", http.StatusBadRequest)
    }


    userFile, httpErr := fs.repository.SaveFile(userId, fh.Filename, fileBytes)
    if httpErr != nil {
	utils.WriteResponse(w, httpErr.Message, httpErr.Code)
	return
    }

    utils.WriteResponse(w, userFile, http.StatusOK)
}

func (fs *FilesService) DeleteHandler(w http.ResponseWriter, req *http.Request) {
    if req.Method != http.MethodDelete {
	http.Error(w, "Use DELETE method instead", http.StatusBadRequest)
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
	    utils.WriteResponse(w, "There is not file to delete", http.StatusNotFound)
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
	http.Error(w, "Use GET method instead", http.StatusBadRequest)
	return
    }
    userId, httpErr := utils.CheckAuth(req)
    if httpErr != nil {
	utils.WriteResponse(w, httpErr.Message, httpErr.Code)
	return
    }

    filename := strings.TrimPrefix(req.URL.Path, "/download/")   
    if filename == "" {
	utils.WriteResponse(w, "File not specified", http.StatusBadRequest)
	return
    }
    filepathToDownload := filepath.Join(filesDirectory, userId, filename) 

    _, err := os.Stat(filepathToDownload)
    if err != nil {
	utils.WriteResponse(w, "File not found", http.StatusNotFound)
	return
    }

    w.Header().Set("Content-Disposition", "attachment; filename="+filename)
    w.Header().Set("Content-Type", "application/octet-stream")

    http.ServeFile(w, req, filepathToDownload)
}

func (fs *FilesService) FilesHanlder(w http.ResponseWriter, req *http.Request) {
    if req.Method != http.MethodGet {
	http.Error(w, "Use GET method instead", http.StatusBadRequest)
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
