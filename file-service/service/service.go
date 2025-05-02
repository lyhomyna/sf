package service

import "net/http"

type FileService interface {
    SaveHandler(w http.ResponseWriter, req *http.Request)
    DeleteHandler(w http.ResponseWriter, req *http.Request)
    DownloadHandler(w http.ResponseWriter, req *http.Request)
    FilesHanlder(w http.ResponseWriter, req *http.Request)
}

type UserImagesService interface {
    GetUserImageHandler(w http.ResponseWriter, req *http.Request)
    SaveUserImageHandler(w http.ResponseWriter, req *http.Request)
}
