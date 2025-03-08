package main

import (
	"context"
	"log"

	"github.com/lyhomyna/sf/auth-service/service"
)

// var envFilePath = filepath.Join(".env")
// 
// func init() { 
//     if err := godotenv.Load(envFilePath); err != nil {
// 	panic(err)
//     }
// }

func main() {
    siglogServer := service.SiglogServer{}
    ctx := context.Background()

    errCh := make(chan error, 1)

    go func() {
	errCh <- siglogServer.Run(ctx)
    }()

    err := <- errCh
    if err != nil {
	log.Println("Server terminated by error:", err)
    }
}
