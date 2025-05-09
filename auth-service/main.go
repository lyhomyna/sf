package main

import (
	"context"
	"log"

	"github.com/lyhomyna/sf/auth-service/service"
)

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
