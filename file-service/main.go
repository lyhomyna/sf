package main

import (
	"log"
	"net/http"

	"github.com/lyhomyna/sf/file-service/config"
	"github.com/lyhomyna/sf/file-service/database"
	"github.com/lyhomyna/sf/file-service/middleware"
	"github.com/lyhomyna/sf/file-service/service"

	filesRepository "github.com/lyhomyna/sf/file-service/repository/files"
	userImagesRepository "github.com/lyhomyna/sf/file-service/repository/userImages"

	filesService "github.com/lyhomyna/sf/file-service/service/files"
	userImagesService "github.com/lyhomyna/sf/file-service/service/userImages"
)

var (
    cfg = config.NewConfig().WithPostgres()
    pgDb = database.GetPostgresDb(cfg)

    fr = filesRepository.NewFilesRepository(pgDb)
    uir = userImagesRepository.NewUserImagesRepository(pgDb)

    fs service.FileService = filesService.NewFilesService(fr, cfg)
    uis service.UserImagesService = userImagesService.NewUserImagesService(uir)
)

func main() {
    mux := http.NewServeMux()

    mux.HandleFunc("/", fs.ListDirHandler)     // GET
    mux.HandleFunc("/create-directory", fs.CreateDirectoryHandler)
    mux.HandleFunc("/create-root", fs.CreateRootDirectoryHandler)

    mux.HandleFunc("/save", fs.SaveHandler)               // POST
    mux.HandleFunc("/delete/", fs.DeleteHandler)          // DELETE
    mux.HandleFunc("/download/", fs.DownloadHandler)      // GET

    mux.HandleFunc("/image/", uis.GetUserImageHandler)        // GET
    mux.HandleFunc("/save-image", uis.SaveUserImageHandler)  // POST

    mux.HandleFunc("/health", func(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusOK)
    })
    mux.Handle("/favion.ico", http.NotFoundHandler())

    handler := middleware.OptionsMiddleware(mux)

    log.Printf("File service. Server is running on port %s", cfg.ServerPort)
    if err := http.ListenAndServe(cfg.ServerPort, handler); err != nil {
	log.Println("Failed to start server:", err)
    }

    log.Println("Bye-bye")
}
