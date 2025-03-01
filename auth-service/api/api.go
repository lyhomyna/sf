package api

import (
	"context"
	"fmt"
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


