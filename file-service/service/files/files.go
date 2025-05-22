package files

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
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
func (fs *FilesService) SaveHandler(w http.ResponseWriter, req *http.Request) {
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

    // Read file's dir
    dir := req.FormValue("dir");
    if dir == "" {
	utils.WriteResponse(w, "File's dire couldn't be empty.", http.StatusBadRequest)
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

    if fh.Filename == "" {
	    utils.WriteResponse(w, "Filename is required", http.StatusBadRequest)
	    return
    }

    userFile, err := fs.repository.SaveFile(userId, fh.Filename, f, dir)
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
    
    // to do not go beyond permitted access limits
    if !isUrlPathCorrect(req.URL.Path) {
	utils.WriteResponse(w,  "Access denied. You don't have permission to perform this action.", http.StatusForbidden)
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

// path must start with '/'
func isUrlPathCorrect(path string) bool {
    splitted := strings.Split(path[1:], "/") 
    if (len(splitted) < 2 || len(splitted) > 2) {
	return false
    }
    return true
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

    // to do not go beyond permitted access limits
    if !isUrlPathCorrect(req.URL.Path) {
	utils.WriteResponse(w,  "Access denied. You don't have permission to perform this action.", http.StatusForbidden)
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

// Dir comes from request URL path
func (fs *FilesService) ListDirHandler(w http.ResponseWriter, req *http.Request) {
    if req.Method != http.MethodGet {
	http.Error(w, "Use GET method instead", http.StatusMethodNotAllowed)
	return
    }
    userId, httpErr := utils.CheckAuth(req)
    if httpErr != nil {
	utils.WriteResponse(w, httpErr.Message, httpErr.Code)
	return
    }

    cleanPath := path.Clean("/" + strings.Trim(req.URL.Path, "/"))

    entries, err := fs.repository.ListDir(cleanPath, userId)
    if err != nil {
	log.Println(err)
	
	if errors.Is(repository.ErrorDirectoryNotFound, err) {
	    utils.WriteResponse(w, "Directory not found", http.StatusNotFound)
	    return
	}

	utils.WriteResponse(w, "Something went wrong. Try again", http.StatusInternalServerError)
	return
    }

    utils.WriteResponseV2(w, entries, http.StatusOK)
}

func (fs *FilesService) CreateDirectoryHandler(w http.ResponseWriter, req *http.Request) {
    if req.Method != http.MethodPost {
	http.Error(w, "Use POST method instead", http.StatusMethodNotAllowed)
	return
    }
    userId, httpErr := utils.CheckAuth(req)
    if httpErr != nil {
	utils.WriteResponse(w, httpErr.Message, httpErr.Code)
	return
    }

    dirName := req.FormValue("name")
    parentDir := req.FormValue("curr_dir")

    // What errors can GetDirIdByPath return?
    parentDirId, err := fs.repository.GetDirIdByPath(userId, parentDir)
    if err != nil {
	log.Println(err)

	// TODO: combine these two errors 
	if errors.Is(repository.ErrorRootDirNotFound, err) || errors.Is(repository.ErrorDirectoryNotFound, err) {
	    utils.WriteResponse(w, "Parent directory not found", http.StatusNoContent)
	    return
	}

	utils.WriteResponse(w, "Something went wrong. Try again", http.StatusInternalServerError)
	return
    }

    dirId, err := fs.repository.CreateDir(userId, parentDirId, dirName)
    if err != nil {
	log.Println(err)
    
	if errors.Is(repository.ErrorDirectoryAlreadyExist, err) {
	    utils.WriteResponse(w, "Directory already exist", http.StatusConflict)
	    return
	}

	utils.WriteResponse(w, "Couldn't create folder. Try again", http.StatusInternalServerError)
	return
    }

    response := struct {
	Id   string `json:"id"`
	Name string `json:"name"`
	FullPath string `json:"fullPath"`
    }{
	    Id: dirId,
	    Name: dirName,
	    FullPath: path.Clean("/" + strings.Trim(parentDir, "/") + "/" + dirName),
    }

    utils.WriteResponseV2(w, response, http.StatusOK)
}
