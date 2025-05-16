package files

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
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

var FileErrorFileExist = errors.New("File with the save name already exists")
// SaveupLoadedFile saves file into server's forlder and returns HttpError object if there is an error. 
func (pr *FilesRepository) SaveFile(userId string, filename string, file io.Reader) (*models.UserFile, error) {
    newFilePath := filepath.Join(filesDirectory, userId, filename)

    if _, err := os.Stat(newFilePath); err == nil {
	return nil, fmt.Errorf("%w: %v", repository.FilesErrorFileExist, err)
    }

    if err := os.MkdirAll(filepath.Dir(newFilePath), os.ModePerm); err != nil {
	return nil, fmt.Errorf("%w: %v", repository.FilesErrorInternal, err) 
    }

    nf, err := os.Create(newFilePath)
    if err != nil {
	return nil, fmt.Errorf("%w: %v", repository.FilesErrorFailureCreateFile, err)
    }
    defer nf.Close()

    // Copy content
    hash := sha256.New()
    writtenBytes, err := io.Copy(nf, io.TeeReader(file, hash))
    if err != nil {
	os.Remove(newFilePath)
	return nil, fmt.Errorf("%w: %v", repository.FilesErrorCopyFailure, err)
    }

    hashString := fmt.Sprintf("%x", hash.Sum(nil))
    fileId, err := pr.saveFileToDb(userId, filename, newFilePath, writtenBytes, hashString)
    if err != nil {
	os.Remove(newFilePath)
	return nil, fmt.Errorf("%w: %v", repository.FilesErrorDbSave, err)
    }

    return &models.UserFile {
	Id: fileId,
	Filename: filename,
    }, nil
}

func (pr *FilesRepository) saveFileToDb(userId string, filename string, filepath string, size int64, hash string) (string, error) {
    ctx := context.Background()
    fileId := uuid.NewString()
    sql := 

    `INSERT INTO files (id, user_id, filename, filepath, size, hash) 
	VALUES ($1, $2, $3, $4, $5, $6)`

    _, err := pr.db.Pool.Exec(ctx, sql, fileId, userId, filename, filepath, size, hash)
    if err != nil {
	return "", err 
    }

    return fileId, nil
}

// DeleteFile deletes file from server and returns error if there is an error
func (pr *FilesRepository) DeleteFile(userId string, fileId string) error {
    uf, err := pr.GetFile(userId, fileId)
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
func (pr *FilesRepository) GetFile(userId string, fileId string) (*models.DbUserFile, error) { 
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
func (pr *FilesRepository) GetFiles(userId string) ([]*models.DbUserFile, error) {
    ctx := context.Background()
    sql := "SELECT * FROM files WHERE user_id=$1"
    userFiles := []*models.DbUserFile{}

    rows, err := pr.db.Pool.Query(ctx, sql, userId)
    if err != nil {
	if errors.Is(err, pgx.ErrNoRows) {
	    return userFiles, nil	
	}

	return nil, fmt.Errorf("%w: %v", repository.FilesErrorDbQuery, err)
    }

    for rows.Next() {
	var uf models.DbUserFile
	if err := rows.Scan(&uf.Id, &uf.UserId, &uf.Filename, &uf.Filepath, &uf.Size, &uf.Hash, &uf.LastAccessed); err != nil {
	    return nil, fmt.Errorf("%w :%v", repository.FilesErrorDbScan, err)
	}
	userFiles = append(userFiles, &uf)
    }

    return userFiles, nil
}
