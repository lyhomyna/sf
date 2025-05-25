package user

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/lyhomyna/sf/auth-service/models"
	"golang.org/x/crypto/bcrypt"
)

func CreateUser(siglog *models.Siglog, req *http.Request) (string, *models.HTTPError) {
    defer req.Body.Close()

    var user models.User
    if err := decodeFromTo(req.Body, &user); err != nil {
	fmt.Println("Could't decode user.")
	return "", &models.HTTPError {
	    Code: http.StatusBadRequest,
	    Message: "Use correct user schema",
	} 
    }

    dbUser, errHttp := constructDbUser(&user);   
    if errHttp != nil {
	fmt.Println(errHttp.Message)
	return "", *&errHttp
    }

    userId, err := siglog.Users.CreateUser(dbUser)
    if err != nil {
	log.Println(err.Error())
	return "", &models.HTTPError {
	    Code: http.StatusInternalServerError,
	    Message: err.Error(),
	}
    }
    
    log.Println("User created.")
    return userId, nil
}

func decodeFromTo(rc io.ReadCloser, target any) error {
    decoder := json.NewDecoder(rc)
    if err := decoder.Decode(target); err != nil {
        return errors.New(fmt.Sprintf("Decode failure. %s", err))
    }
    return nil
}

func constructDbUser(user *models.User) (*models.DbUser, *models.HTTPError) {
    encryptedPassword, err := encryptPassword(user.Password)
    if err != nil {
	return nil, &models.HTTPError {
	    Code: http.StatusInternalServerError,
	    Message: err.Error(),
	}
    }

    dbUser := &models.DbUser{
	Id: uuid.NewString(),
	Email: user.Email,
	Password: encryptedPassword,
    }

    return dbUser, nil
}

func encryptPassword(password string) (string, error) {
    bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        return "", errors.New(fmt.Sprintf("Can't encrypt password. %s", err))
    }
    return string(bytes), nil
}


// danger function
func DeleteUser(userId string, siglog *models.Siglog) *models.HTTPError {
    if err := siglog.Users.DeleteUser(userId); err != nil {
	log.Println(fmt.Sprintf("Couldn't delete user '%s'", userId))
	return &models.HTTPError {
	    Code: http.StatusInternalServerError,
	    Message: "Couldn't delete user.",
	}
    }

    log.Println("User and session has been deleted.")
    return nil
}

func GetUserByEmail(email string, siglog *models.Siglog) (*models.DbUser, *models.HTTPError) {
    dbUser, err := siglog.Users.GetUserByEmail(email)
    if err != nil {
	log.Println(err.Error())
	return nil, &models.HTTPError{
	    Code: http.StatusInternalServerError,
	    Message: "Internal server error",
	}
    }
    if dbUser == nil {
	log.Println("Couldn't find user by email")
	return nil, &models.HTTPError {
	    Code: http.StatusNotFound,
	    Message: "User not found",
	}
    }
    return dbUser, nil
}

func ComparePasswords(passwordHash string, possiblePassword string) error {
    err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(possiblePassword))
    return err
}

func GetById(userId string, siglog *models.Siglog) (*models.User, *models.HTTPError) {
    user, err := siglog.Users.ReadUserById(userId)
    if err != nil {
	return nil, &models.HTTPError {
	    Code: http.StatusNotFound, 
	    Message: "User does not exist", 
	}
    }
    return user, nil
}



