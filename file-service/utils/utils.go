package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/lyhomyna/sf/file-service/models"
)

var sessionCookieName = "session-id"
var authServiceBaseUrl = "http://auth-service:8081"

func WriteResponse(w http.ResponseWriter, data any, code int) {
    w.Header().Set("Content-Type", "application/json")
    w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")

    w.WriteHeader(code)
    d := struct {
	Data any `json:"data"`
    } {
	Data: data,
    }

    response, err := json.Marshal(d)
    if err != nil {
	panic(err)
    }
    w.Write(response)
}

// returns user id or error
func CheckAuth(req *http.Request) (string, *models.HttpError) {
    sessionCookie, err := req.Cookie(sessionCookieName)
    if err != nil {
        return "", &models.HttpError{
	    Code: http.StatusUnauthorized,
	    Message: "Session cookie missing",
	}
    }

    uid, err := verifySession(sessionCookie); 
    if err != nil {
        return "", &models.HttpError{
	    Code: http.StatusUnauthorized,
	    Message: err.Error(),
	}
    }

    return uid, nil
}

// verifySession returns user id of error
func verifySession(sessionCookie *http.Cookie) (string, error) {
    reqUrl := fmt.Sprintf("%s/check-auth", authServiceBaseUrl)
    
    req, _ := http.NewRequest("GET", reqUrl, nil)

    req.AddCookie(sessionCookie)

    client := &http.Client{}

    resp, err := client.Do(req)
    if err != nil {
	log.Println("Unable to verify session:", err)
	return "", errors.New("Unable to verify session")
    }

    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
	log.Println("Unable to verify session: Status", resp.StatusCode)
	return "", errors.New("Unable to verify session")
    }

    var res struct {
	Id string `json:"userId"`
    }
    json.NewDecoder(resp.Body).Decode(&res)

    return res.Id, nil
}
