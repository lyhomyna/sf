package files

import (
	"context"
	"crypto/sha256"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/lyhomyna/sf/file-service/database"
	"github.com/lyhomyna/sf/file-service/models"
)

type FilesRepository struct {
    db *database.Postgres
}

// Directory, where uploaded files will be saved
var filesDirectory = filepath.Join("files")

func NewFilesRepository(db *database.Postgres) *FilesRepository {
    return &FilesRepository{
	db: db,
    }
}

// SaveupLoadedFile saves file into server's forlder and returns HttpError object if there is an error. 
func (pr *FilesRepository) SaveFile(userId string, filename string, fileBytes []byte) *models.HttpError {
    newFilePath := filepath.Join(filesDirectory, userId, filename)

    if _, err := os.Stat(newFilePath); err == nil {
	log.Println("File with the same name already exists")
	return &models.HttpError{
	    Message: "File with the same name already exists",
	    Code: http.StatusBadRequest,
	}
    }

    if err := os.MkdirAll(filepath.Dir(newFilePath), os.ModePerm); err != nil {
	return &models.HttpError{
	    Message: "Something went wrong",
	    Code: http.StatusInternalServerError,
	}
    }

    nf, err := os.Create(newFilePath)
    if err != nil {
	log.Println("Failed to create a file:", err)
	return &models.HttpError{
	    Message: "Something went wrong while a server was creating file",
	    Code: http.StatusInternalServerError,
	}
    }
    defer nf.Close()

    // Copy content
    _, err = nf.Write(fileBytes)
    if err != nil {
	os.Remove(newFilePath)
	log.Println("Failed to save file:", err)
	return &models.HttpError{
	    Message: "Something went wrong when the server was copying content from your file",
	    Code: http.StatusInternalServerError,
	}
    }

    hash, err := calculateSHA256(nf)
    if err != nil {
	fmt.Println("Error calculating hash for file:", err.Error())
	return &models.HttpError{ 
	    Message: "Something went wrong",
	    Code: http.StatusInternalServerError,
	}
    }
    fmt.Println("Hash for file:", hash)

    err = pr.saveFileToDb(userId, filename, newFilePath, len(fileBytes), hash)
    if err != nil {
	os.Remove(newFilePath)
	
	return &models.HttpError{
	    Message: "Couldn't save file",
	    Code: http.StatusInternalServerError,
	}
    }

    return nil
}

func (pr *FilesRepository) saveFileToDb(userId string, filename string, filepath string, size int, hash string) error {
    log.Println("TRYING TO SAVE FILE...")

    ctx := context.Background()
    fileId := uuid.NewString()
    sql := "INSERT INTO files (id, user_id, filename, filepath, size, hash) VALUES ($1, $2, $3, $4, $5, $6)"

    _, err := pr.db.Pool.Exec(ctx, sql, fileId, userId, filename, filepath, size, hash)
    if err != nil {
	log.Println("Failed to save file to db:", err.Error())
	return err 
    }

    log.Println("FILE SAVED SUCCESSFULY!")
    return nil
}

func calculateSHA256(file io.Reader) (string, error) {
    hasher := sha256.New()
    if _, err := io.Copy(hasher, file); err != nil {
	return "", err
    }

    hash := hasher.Sum(nil)
    return fmt.Sprintf("%x", hash), nil
}

// DeleteFile deletes file from server and returns error if there is an error
func (pr *FilesRepository) DeleteFile(userId string, filename string) *models.HttpError {
    


    fullFilepath := filepath.Join(filesDirectory, userId, filename)
    if _, err := os.Stat(fullFilepath); err != nil {
	log.Println("File doesn't exist:", err.Error())
	return &models.HttpError{
	    Message: "File doesn't exist",
	    Code: http.StatusNoContent,
	}
    }

    if err := os.Remove(fullFilepath); err != nil {
	log.Println("File doesn't exist:", err.Error())
	return &models.HttpError{
	    Message: "Error removing file",
	    Code: http.StatusInternalServerError,
	}
    }

    // Somehow delete file from db (maybe by HASH :))) )

    return nil 
}

// GetFilenames returns list of user filenames or an error
func (pr *FilesRepository) GetFilenames(userId string) ([]string, *models.HttpError) {
    filesPath := filepath.Join(filesDirectory, userId)
    entries, err := os.ReadDir(filesPath)
    if err != nil {
	return nil, &models.HttpError { Code: http.StatusInternalServerError, Message: "Cannot read from user directory",
	}
    }

    filenames := []string {}
    for _, entry := range entries {
	filenames = append(filenames, entry.Name())
    }

    return filenames, nil
}
