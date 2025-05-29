package handler

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"path"
	"strings"

	"github.com/lyhomyna/sf/file-service/config"
	"github.com/lyhomyna/sf/file-service/middleware"
	"github.com/lyhomyna/sf/file-service/models"
	"github.com/lyhomyna/sf/file-service/service"
	"github.com/lyhomyna/sf/file-service/utils"
)

type Handler struct {
	services service.Services
	cfg      *config.Config
}

func NewHandler(services service.Services, cfg *config.Config) *Handler {
	return &Handler{
		services: services,
		cfg:      cfg,
	}
}

func (h *Handler) Run() error {
	mux := http.NewServeMux()

	mux.HandleFunc("/", h.listDirHandler)                    // GET
	mux.HandleFunc("/create-directory", h.createDirHandler)  // POST
	mux.HandleFunc("/delete-directory/", h.deleteDirHandler) // POST
	mux.HandleFunc("/create-root", h.createRootDirHandler)   // POST

	mux.HandleFunc("/save", h.saveFileHandler)          // POST
	mux.HandleFunc("/delete/", h.deleteFileHandler)     // DELETE
	mux.HandleFunc("/download/", h.downloadFileHanlder) // GET

	mux.HandleFunc("/image/", h.getUserImageHandler)      // GET
	mux.HandleFunc("/save-image", h.saveUserImageHandler) // POST

	mux.HandleFunc("/health", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	mux.Handle("/favion.ico", http.NotFoundHandler())

	handler := middleware.OptionsMiddleware(mux)

	log.Printf("File service. Server is running on port %s", h.cfg.ServerPort)
	if err := http.ListenAndServe(h.cfg.ServerPort, handler); err != nil {
		return fmt.Errorf("Failed to start server: %w", err)
	}

	return nil
}

func (h *Handler) listDirHandler(w http.ResponseWriter, req *http.Request) {
	userId, httpErr := authorizeRequest(req, http.MethodGet)
	if httpErr != nil {
		http.Error(w, httpErr.Message, httpErr.Code)
		return
	}

	cleanPath := path.Clean("/" + strings.Trim(req.URL.Path, "/"))
	entries, listDirHttpErr := h.services.DirService.ListDir(cleanPath, userId)
	if listDirHttpErr != nil {
		http.Error(w, listDirHttpErr.Message, listDirHttpErr.Code)
		return
	}

	utils.WriteResponseV2(w, entries, http.StatusOK)
}

func (h *Handler) createDirHandler(w http.ResponseWriter, req *http.Request) {
	userId, httpErr := authorizeRequest(req, http.MethodPost)
	if httpErr != nil {
		http.Error(w, httpErr.Message, httpErr.Code)
		return
	}

	dirName := req.FormValue("name")
	parentDirPath := req.FormValue("curr_dir")

	dir, createDirHttpError := h.services.DirService.CreateDirectory(userId, parentDirPath, dirName)
	if createDirHttpError != nil {
		http.Error(w, createDirHttpError.Message, createDirHttpError.Code)
		return
	}

	utils.WriteResponseV2(w, dir, http.StatusOK)
}

func (h *Handler) deleteDirHandler(w http.ResponseWriter, req *http.Request) {
	userId, httpErr := authorizeRequest(req, http.MethodDelete)
	if httpErr != nil {
		http.Error(w, httpErr.Message, httpErr.Code)
		return
	}

	// >>  TEST THIS (do that need "/" in the end of a path?)
	cleanPath := path.Clean("/" + strings.Trim(req.URL.Path, "/"))
	if delDirHttpError := h.services.DirService.DeleteDirectory(cleanPath, userId); delDirHttpError != nil {
		http.Error(w, delDirHttpError.Message, delDirHttpError.Code)
		return
	}

	utils.WriteResponse(w, "Directory deleted successfully", http.StatusOK)
}

func (h *Handler) createRootDirHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(w, "Use POST method instead", http.StatusMethodNotAllowed)
		return
	}

	if !h.hasRootCreationAccess(req) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	userId := req.FormValue("userId")

	dirId, httpErr := h.services.DirService.CreateRootDirectory(userId)
	if httpErr != nil {
		http.Error(w, httpErr.Message, httpErr.Code)
		return
	}

	response := struct {
		Id string `json:"id"`
	}{
		Id: dirId,
	}

	utils.WriteResponseV2(w, response, http.StatusOK)
}

func (h *Handler) hasRootCreationAccess(req *http.Request) bool {
	authHeader := req.Header.Get("Authorization")
	return authHeader == "Bearer "+h.cfg.AuthToFileToken
}

func (h *Handler) saveFileHandler(w http.ResponseWriter, req *http.Request) {
	userId, httpErr := authorizeRequest(req, http.MethodPost)
	if httpErr != nil {
		http.Error(w, httpErr.Message, httpErr.Code)
		return
	}

	if !strings.Contains(req.Header.Get("Content-Type"), "multipart/form-data") {
		http.Error(w, "Expected multipart/form-data", http.StatusUnsupportedMediaType)
		return
	}

	// Read file's dir
	dirPath := req.FormValue("dir")
	if dirPath == "" {
		http.Error(w, "The file directory cannot be empty", http.StatusBadRequest)
		return
	}

	dirId, getDirIdError := h.services.DirService.GetDirIdByPath(dirPath, userId)
	if getDirIdError != nil {
		http.Error(w, getDirIdError.Message, getDirIdError.Code)
		return
	}

	responseFile, saveHttpError := h.services.FileService.SaveFile(userId, dirId, dirPath, req)
	if saveHttpError != nil {
		http.Error(w, saveHttpError.Message, saveHttpError.Code)
		return
	}

	utils.WriteResponseV2(w, responseFile, http.StatusOK)
}

func (h *Handler) deleteFileHandler(w http.ResponseWriter, req *http.Request) {
	userId, httpErr := authorizeRequest(req, http.MethodDelete)
	if httpErr != nil {
		http.Error(w, httpErr.Message, httpErr.Code)
		return
	}

	fileId := extractIdFromPath(req.URL.Path)
	if fileId == "" {
		http.Error(w, "File not specified", http.StatusBadRequest)
		return
	}

	if deleteHttpErr := h.services.FileService.DeleteFile(userId, fileId); deleteHttpErr != nil {
		http.Error(w, deleteHttpErr.Message, deleteHttpErr.Code)
		return
	}

	// >>> TEST THIS
	utils.WriteResponseV2(w, "File successfully deleted", http.StatusOK)
}

func extractIdFromPath(path string) string {
	urlChanks := strings.Split(path, "/")
	return urlChanks[len(urlChanks)-1]
}

func (h *Handler) downloadFileHanlder(w http.ResponseWriter, req *http.Request) {
	userId, httpErr := authorizeRequest(req, http.MethodGet)
	if httpErr != nil {
		http.Error(w, httpErr.Message, httpErr.Code)
		return
	}

	fileId := extractIdFromPath(req.URL.Path)
	if fileId == "" {
		http.Error(w, "File not specified", http.StatusBadRequest)
		return
	}

	fileEntry, getFileEntryHttpError := h.services.FileService.GetFileEntry(userId, fileId)
	if getFileEntryHttpError != nil {
		http.Error(w, getFileEntryHttpError.Message, getFileEntryHttpError.Code)
		return
	}

	w.Header().Set("Content-Disposition", "attachment; filename="+fileEntry.Name)
	w.Header().Set("Content-Type", "application/octet-stream")

	http.ServeFile(w, req, fileEntry.Path)
}

func (h *Handler) saveUserImageHandler(w http.ResponseWriter, req *http.Request) {
	userId, httpErr := authorizeRequest(req, http.MethodPost)
	if httpErr != nil {
		http.Error(w, httpErr.Message, httpErr.Code)
		return
	}

	imagePath, saveHttpErr := h.services.UserImageService.SaveUserImage(userId, req)
	if saveHttpErr != nil {
		http.Error(w, saveHttpErr.Message, saveHttpErr.Code)
		return
	}

	response := &models.ImageJson{
		ImageUrl: imagePath,
	}

	utils.WriteResponseV2(w, response, http.StatusOK)
}

func (h *Handler) getUserImageHandler(w http.ResponseWriter, req *http.Request) {
	_, httpErr := authorizeRequest(req, http.MethodGet)
	if httpErr != nil {
		http.Error(w, httpErr.Message, httpErr.Code)
		return
	}

	// >>  TEST THIS (do that need "/" in the end of a path?)
	cleanPath := path.Clean("/" + strings.Trim(req.URL.Path, "/"))

	imageData, getImageHttError := h.services.UserImageService.GetUserImage(cleanPath)
	if getImageHttError != nil {
		http.Error(w, getImageHttError.Message, getImageHttError.Code)
		return
	}
	defer imageData.ImageFile.Close()

	w.Header().Set("Content-Type", http.DetectContentType(imageData.ContentTypeChunk))
	if _, err := io.Copy(w, imageData.ImageFile); err != nil {
		log.Printf("Image wasn't sent to the user. Reason:%v", err)
		return
	}
}

func authorizeRequest(req *http.Request, allowedHttpMethod string) (string, *models.HttpError) {
	if req.Method != allowedHttpMethod {
		return "", &models.HttpError{
			Message: fmt.Sprintf("Use %s method instead", allowedHttpMethod),
			Code:    http.StatusMethodNotAllowed,
		}
	}

	userId, httpErr := utils.CheckAuth(req)
	if httpErr != nil {
		return "", httpErr
	}

	return userId, nil
}
