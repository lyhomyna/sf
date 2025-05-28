package siglog

import (
	"context"
	"fmt"

	"github.com/lyhomyna/sf/auth-service/handlers"
	"github.com/lyhomyna/sf/auth-service/repository"
	"github.com/lyhomyna/sf/auth-service/service"
)

type SiglogServer struct {
   httpServer *handlers.HttpServer 
}

func (s *SiglogServer) Run(ctx context.Context) error {
    ctx, cancel := context.WithCancel(ctx)

    errCh := make(chan error, 1)
    go func() {
	repos := repository.GetRepos(ctx)
	if repos == nil {
	    errCh <- fmt.Errorf("Coudn't ger repos")
	    return
	}
	services := service.GetServices(ctx, repos)

	s.httpServer = &handlers.HttpServer{
	    Services: services,
	} 

	err := s.httpServer.Run(ctx)
	if err != nil {
	    errCh <- fmt.Errorf("HTTP server error. %w", err)
	}
    } ()

    err := <-errCh

    cancel()
    return err
}
