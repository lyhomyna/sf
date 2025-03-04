package api

import (
	"context"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/lyhomyna/sf/auth-service/handlers/user"
	"github.com/lyhomyna/sf/auth-service/repository"
)

type httpServer struct {
    http *http.Server
}

func (s *httpServer) Run(ctx context.Context) error {
    userRoutes := user.Routes()
    sessionRoutes := user.Routes()

    siglog := repository.GetSiglog()

    if siglog == nil {
	return errors.New("Couldn't get SigLog.")
    }

    handler := http.NewServeMux()
    handler.

    s.http = &http.Server {
	Addr: ":8080",
	Handler: handler,
	ReadHeaderTimeout: 5 * time.Second, // mitigate risk of Slowloris Attack
    }
    
    log.Println("HTTP server is running on port 8080")
    if err := s.http.ListenAndServe(); err != nil {
	return err
    }

    return nil 
}
