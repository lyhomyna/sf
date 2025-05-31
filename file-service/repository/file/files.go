package file

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/lyhomyna/sf/file-service/database"
	"github.com/lyhomyna/sf/file-service/models"
	"github.com/lyhomyna/sf/file-service/repository"
)

type FileRepository struct {
    db *database.Postgres
}

func NewFileRepository(db *database.Postgres) *FileRepository {
    return &FileRepository{
	db: db,
    }
}

// Directory, where uploaded files will be saved
var filesDirectory = filepath.Join("files")

// :::::::::::::::::::::::::::::::::::::::::::::::::::::::::

var FileErrorFileExist = errors.New("File with the save name already exists")
// SaveupLoadedFile saves file into server's forlder and returns HttpError object if there is an error. 
func (pr *FileRepository) SaveFile(userId, filename, dirPath, dirId string, file io.Reader) (*models.UserFile, error) {
    newFilePath := filepath.Join(filesDirectory, userId, dirPath, filename)

    // check if file already exists on disk
    if _, err := os.Stat(newFilePath); err == nil {
	return nil, fmt.Errorf("%w: %v", repository.FilesErrorFileExist, err)
    }

    // create directory if not exitst (this only for new users)
    if err := os.MkdirAll(filepath.Dir(newFilePath), os.ModePerm); err != nil {
	return nil, fmt.Errorf("%w: %v", repository.FilesErrorInternal, err) 
    }

    // create file in filesystem and compute hash
    nf, err := os.Create(newFilePath)
    if err != nil {
	return nil, fmt.Errorf("%w: %v", repository.FilesErrorFailureCreateFile, err)
    }
    defer nf.Close()

    hash := sha256.New()
    writtenBytes, err := io.Copy(nf, io.TeeReader(file, hash))
    if err != nil {
	os.Remove(newFilePath)
	return nil, fmt.Errorf("%w: %v", repository.FilesErrorCopyFailure, err)
    }

    hashString := fmt.Sprintf("%x", hash.Sum(nil))
    fileId, err := pr.saveFileToDb(userId, dirId, filename, filepath.Join(userId, dirPath, filename), writtenBytes, hashString)
    if err != nil {
	os.Remove(newFilePath)
	return nil, fmt.Errorf("%w: %v", repository.FilesErrorDbSave, err)
    }

    return &models.UserFile {
	Id: fileId,
	Filename: filename,
    }, nil
}

func (pr *FileRepository) saveFileToDb(userId, dirId, filename, filepath string, size int64, hash string) (string, error) {
    ctx := context.Background()
    fileId := uuid.NewString()
    sql := 
    `INSERT INTO files (id, user_id, directory_id, filename, filepath, size, hash) 
	VALUES ($1, $2, $3, $4, $5, $6, $7)`

    _, err := pr.db.Pool.Exec(ctx, sql, fileId, userId, dirId, filename, filepath, size, hash)
    if err != nil {
	return "", err 
    }

    return fileId, nil
}

// :::::::::::::::::::::::::::::::::::::::::::::::::::::::::

// DeleteFile deletes file from server and returns error if there is an error
func (pr *FileRepository) DeleteFile(userId string, fileId string) error {
    uf, err := pr.GetFile(userId, fileId)
    if err != nil {
	return err
    }

    fullFilePath := filepath.Join(filesDirectory, uf.Filepath)

    if _, err := os.Stat(fullFilePath); err != nil {
	return fmt.Errorf("%w: %v", repository.FilesErrorFileNotExist, err)
    }

    if err := pr.removeFileFromDb(fileId); err != nil {
	return err
    }

    err = os.Remove(fullFilePath)

    return err
}

// Can return FilesErrorDbQuery
func (pr *FileRepository) removeFileFromDb(fileId string) error {
    ctx := context.Background()

    sql := "DELETE FROM files WHERE id=$1"

    _, err := pr.db.Pool.Exec(ctx, sql, fileId)
    if err != nil {
	return fmt.Errorf("%w:%v", repository.FilesErrorDbQuery, err)	 
    }

    return nil
}

// :::::::::::::::::::::::::::::::::::::::::::::::::::::::::

// Can return FilesErrorFailureToRetrieve
func (pr *FileRepository) GetFile(userId string, fileId string) (*models.DbUserFile, error) { 
    ctx := context.Background()

    sql := "SELECT id, user_id, filename, filepath, size, hash, last_accessed FROM files WHERE user_id=$1 AND id=$2"

    var uf models.DbUserFile
    if err := pr.db.Pool.QueryRow(ctx, sql, userId, fileId).Scan(&uf.Id, &uf.UserId, &uf.Filename, &uf.Filepath, &uf.Size, &uf.Hash, &uf.LastAccessed); err != nil {
	return nil, fmt.Errorf("%w: %v", repository.FilesErrorFailureToRetrieve, err)
    }

    return &uf, nil
}
