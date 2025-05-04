package files

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/lyhomyna/sf/file-service/database"
	"github.com/lyhomyna/sf/file-service/models"
	"github.com/lyhomyna/sf/file-service/repository"
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
func (pr *FilesRepository) SaveFile(userId string, filename string, fileBytes []byte) (*models.UserFile, *models.HttpError) {
    newFilePath := filepath.Join(filesDirectory, userId, filename)

    if _, err := os.Stat(newFilePath); err == nil {
	log.Println("File with the same name already exists")
	return nil, &models.HttpError{
	    Message: "File with the same name already exists",
	    Code: http.StatusBadRequest,
	}
    }

    if err := os.MkdirAll(filepath.Dir(newFilePath), os.ModePerm); err != nil {
	return nil, &models.HttpError{
	    Message: "Something went wrong",
	    Code: http.StatusInternalServerError,
	}
    }

    nf, err := os.Create(newFilePath)
    if err != nil {
	log.Println("Failed to create a file:", err)
	return nil, &models.HttpError{
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
	return nil, &models.HttpError{
	    Message: "Something went wrong when the server was copying content from your file",
	    Code: http.StatusInternalServerError,
	}
    }

    hash, err := calculateSHA256(nf)
    if err != nil {
	fmt.Println("Error calculating hash for file:", err.Error())
	return nil, &models.HttpError{ 
	    Message: "Something went wrong",
	    Code: http.StatusInternalServerError,
	}
    }
    fmt.Println("Hash for file:", hash)

    fileId, err := pr.saveFileToDb(userId, filename, newFilePath, len(fileBytes), hash)
    if err != nil {
	os.Remove(newFilePath)
	
	return nil, &models.HttpError{
	    Message: "Couldn't save file",
	    Code: http.StatusInternalServerError,
	}
    }

    return &models.UserFile {
	Id: fileId,
	Filename: filename,
    }, nil
}

func (pr *FilesRepository) saveFileToDb(userId string, filename string, filepath string, size int, hash string) (string, error) {
    log.Println("TRYING TO SAVE FILE...")

    ctx := context.Background()
    fileId := uuid.NewString()
    sql := "INSERT INTO files (id, user_id, filename, filepath, size, hash) VALUES ($1, $2, $3, $4, $5, $6)"

    _, err := pr.db.Pool.Exec(ctx, sql, fileId, userId, filename, filepath, size, hash)
    if err != nil {
	log.Println("Failed to save file to db:", err.Error())
	return "", err 
    }

    log.Println("FILE SAVED SUCCESSFULLY!")
    return fileId, nil
}

func calculateSHA256(file io.Reader) (string, error) {
    hasher := sha256.New()
    if _, err := io.Copy(hasher, file); err != nil {
	return "", err
    }

    hash := hasher.Sum(nil)
    return fmt.Sprintf("%x", hash), nil
}

var FilesErrorFileNotExist = errors.New("File doesn't exist")
// DeleteFile deletes file from server and returns error if there is an error
func (pr *FilesRepository) DeleteFile(userId string, fileId string) error {
    uf, err := pr.getFile(userId, fileId)
    if err != nil {
	return err
    }

    if _, err := os.Stat(uf.Filepath); err != nil {
	pr.removeFileFromDb(fileId)
	return fmt.Errorf("%w: %v", repository.FilesErrorFileNotExist, err)
    }

    if err := pr.removeFileFromDb(fileId); err != nil {
	return err
    }

    os.Remove(uf.Filepath)

    return nil 
}

// Can return FilesErrorFailureToRetrieve
func (pr *FilesRepository) getFile(userId string, fileId string) (*models.DbUserFile, error) { 
    ctx := context.Background()
    sql := "SELECT * FROM files WHERE user_id=$1 AND id=$2"

    var uf models.DbUserFile
    if err := pr.db.Pool.QueryRow(ctx, sql, userId, fileId).Scan(&uf.Id, &uf.UserId, &uf.Filename, &uf.Filepath, &uf.Size, &uf.Hash, &uf.LastAccessed); err != nil {
	return nil, fmt.Errorf("%w: %v", repository.FilesErrorFailureToRetrieve, err)
    }

    return &uf, nil
}

// Can return FilesErrorDbQuery
func (pr *FilesRepository) removeFileFromDb(fileId string) error {
    ctx := context.Background()
    sql := "DELETE FROM files WHERE id=$1"

    _, err := pr.db.Pool.Exec(ctx, sql, fileId)
    if err != nil {
	return fmt.Errorf("%w:%v", repository.FilesErrorDbQuery, err)	 
    }

    return nil
}

// GetFilenames returns list of user filenames or an error
func (pr *FilesRepository) GetFiles(userId string) ([]*models.DbUserFile, *models.HttpError) {
    ctx := context.Background()
    sql := "SELECT * FROM files WHERE user_id=$1";
    userFiles := []*models.DbUserFile{}

    rows, err := pr.db.Pool.Query(ctx, sql, userId)
    if err != nil {
	if errors.Is(err, pgx.ErrNoRows) {
	    return userFiles, nil	
	}

	log.Println("Couldn't read user files")
	return nil, &models.HttpError{
	    Message: "Couldn't read user files",
	    Code: http.StatusInternalServerError,
	}
    }

    for rows.Next() {
	var uf models.DbUserFile
	if err := rows.Scan(&uf.Id, &uf.UserId, &uf.Filename, &uf.Filepath, &uf.Size, &uf.Hash, &uf.LastAccessed); err != nil {
	    log.Println("Error retrieving user file:", err.Error())
	    return nil, &models.HttpError{
		Message: "Error retrieving user file",
		Code: http.StatusInternalServerError,
	    }
	}
	userFiles = append(userFiles, &uf)
    }

    return userFiles, nil
}
