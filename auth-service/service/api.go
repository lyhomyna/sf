package service

import (
	"context"
	"fmt"

	"github.com/lyhomyna/sf/auth-service/handlers"
)

type SiglogServer struct {
   httpServer *handlers.HttpServer 
}

func (s *SiglogServer) Run(ctx context.Context) error {
    s.httpServer = &handlers.HttpServer{} 
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
