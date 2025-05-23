package service

import "net/http"

type FileService interface {
    SaveHandler(w http.ResponseWriter, req *http.Request)
    DeleteHandler(w http.ResponseWriter, req *http.Request)
    DownloadHandler(w http.ResponseWriter, req *http.Request)
    ListDirHandler(w http.ResponseWriter, req *http.Request)
    CreateDirectoryHandler(w http.ResponseWriter, req *http.Request)
    CreateRootDirectoryHandler(w http.ResponseWriter, req *http.Request)
}

type UserImagesService interface {
    GetUserImageHandler(w http.ResponseWriter, req *http.Request)
    SaveUserImageHandler(w http.ResponseWriter, req *http.Request)
}
