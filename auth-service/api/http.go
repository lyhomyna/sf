package api

import (
	"net/http"

	"github.com/lyhomyna/sf/auth-service/database/models"
)

func NewHttpServer(dao *models.Siglog) http.Handler {
    mux := http.NewServeMux()

    // TODO: Use controllers to manipulate endpoints
    panic("Not yet implemented.")

    return mux
}
