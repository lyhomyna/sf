package main

import (
	"log"

	"github.com/lyhomyna/sf/file-service/config"
	"github.com/lyhomyna/sf/file-service/database"
	"github.com/lyhomyna/sf/file-service/handler"
	"github.com/lyhomyna/sf/file-service/repository"
	"github.com/lyhomyna/sf/file-service/service"

	dirRepository "github.com/lyhomyna/sf/file-service/repository/dir"
	fileRepository "github.com/lyhomyna/sf/file-service/repository/file"
	userImageRepository "github.com/lyhomyna/sf/file-service/repository/userImage"

	dirService "github.com/lyhomyna/sf/file-service/service/dir"
	fileService "github.com/lyhomyna/sf/file-service/service/file"
	userImageService "github.com/lyhomyna/sf/file-service/service/userImages"
)

var (
    cfg = config.NewConfig().WithPostgres()
    pgDb = database.GetPostgresDb(cfg)

    fr repository.FileRepository = fileRepository.NewFileRepository(pgDb)
    uir repository.UserImageRepository = userImageRepository.NewUserImageRepository(pgDb)
    dr repository.DirRepository = dirRepository.NewDirRepository(pgDb)

    fs service.FileService = fileService.NewFileService(fr)
    uis service.UserImageService = userImageService.NewUserImagesService(uir)
    ds service.DirService = dirService.NewDirService(dr)
)

func main() {
    handler := handler.NewHandler(service.Services{ 
	FileService: fs,
	UserImageService: uis,
	DirService: ds,
    }, cfg)

    if err := handler.Run(cfg); err != nil {
	log.Println(err)
    }

    log.Println("Bye-bye")
}
