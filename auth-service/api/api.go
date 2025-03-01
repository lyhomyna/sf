package api

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/lyhomyna/sf/auth-service/database"
)

type SiglogServer struct {
   httpServer *httpServer 
}

func (s *SiglogServer) Run(ctx context.Context) error {
    s.httpServer = &httpServer{} 
    ctx, cancel := context.WithCancel(ctx)

    errCh := make(chan error, 1)
    go func() {
	err := s.httpServer.Run(ctx)
	if err != nil {
	    err = fmt.Errorf("HTTP server error. %w", err)
	}
	errCh <- err
    } ()

    err := <-errCh

    cancel()
    return err
}

type httpServer struct {
    http *http.Server
}

func (s *httpServer) Run(ctx context.Context) error {
    handler := NewHttpServer() 

    s.http = &http.Server {
	Addr: ":8080",
	Handler: handler,
	ReadHeaderTimeout: 5 * time.Second, // mitigate risk of Slowloris Attack
    }

    siglog := database.GetDao()
    if siglog == nil {
	return errors.New("Couldn't get database dao.")
    }

    log.Println("HTTP server is running on port 8080")
    if err := s.http.ListenAndServe(); err != nil {
	return err
    }

    return nil 
}
