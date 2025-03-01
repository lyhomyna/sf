package api

import (
	"context"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/lyhomyna/sf/auth-service/controllers"
	"github.com/lyhomyna/sf/auth-service/database"
)

type httpServer struct {
    http *http.Server
}

func (s *httpServer) Run(ctx context.Context) error {
    siglog := database.GetSiglog()

    if siglog == nil {
	return errors.New("Couldn't get SigLog.")
    }

    sessionsController := controllers.SessionsController{
	Siglog: siglog,
    }
    usersController := controllers.UsersController{
	Siglog: siglog,
    }

    handler := NewHttpServer(usersController, sessionsController) 

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

var usersController controllers.UsersController
var sessionsController controllers.SessionsController

func NewHttpServer(uc controllers.UsersController, sc controllers.SessionsController) http.Handler {
    mux := http.NewServeMux()

    panic("Not yet implemented.")

    mux.HandleFunc("/login", nil)
    mux.HandleFunc("/register", nil)
    mux.HandleFunc("/logout", nil)

    return mux
}
