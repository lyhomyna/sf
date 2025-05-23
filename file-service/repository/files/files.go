package files

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

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

//////////////////////////////////////////////////////////////////////////////////////////////////////////////////

var FileErrorFileExist = errors.New("File with the save name already exists")
// SaveupLoadedFile saves file into server's forlder and returns HttpError object if there is an error. 
func (pr *FilesRepository) SaveFile(userId string, filename string, file io.Reader, dirPath string) (*models.UserFile, error) {
    // get directory id
    dirId, err := pr.GetDirIdByPath(userId, dirPath)
    if err != nil {
	return nil, fmt.Errorf("%w: %v", repository.ErrorDirectoryNotFound, err)
    }	

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
    fileId, err := pr.saveFileToDb(userId, dirId, filename, newFilePath, writtenBytes, hashString)
    if err != nil {
	os.Remove(newFilePath)
	return nil, fmt.Errorf("%w: %v", repository.FilesErrorDbSave, err)
    }

    return &models.UserFile {
	Id: fileId,
	Filename: filename,
    }, nil
}

func (pr *FilesRepository) saveFileToDb(userId, dirId, filename, filepath string, size int64, hash string) (string, error) {
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

//////////////////////////////////////////////////////////////////////////////////////////////////////////////////

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

//////////////////////////////////////////////////////////////////////////////////////////////////////////////////

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

//////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// ::::::::::::::::::
// ::: DEPRECATED :::
// ::::::::::::::::::
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

//////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func (pr *FilesRepository) ListDir(path, userId string) ([]models.DirEntry, error) {
    dirId, err := pr.GetDirIdByPath(userId, path)
    if err != nil {
	return nil, fmt.Errorf("%w: %v", repository.ErrorPathNotFound, err)
    }

    // Directories
    dirs, err := pr.fetchSubdirectories(userId, dirId, path)
    if err != nil {
	return nil, err
    }

    // Files
    files, err := pr.fetchFiles(userId, dirId, path)
    if err != nil {
	return nil, err
    }

    return append(dirs, files...), nil
}

// path should look like this: /inner/directory
func (pr *FilesRepository) GetDirIdByPath(userId, path string) (string, error) {
    if path == "/" {
	return pr.getRootDirId(userId)
    }
    return pr.getNestedDirId(userId, path)
}

func (pr *FilesRepository) getRootDirId(userId string) (string, error) {
    sql := "SELECT id FROM DIRECTORIES WHERE user_id=$1 AND parent_id=NULL"

    var rootId string
    err := pr.db.Pool.QueryRow(context.Background(), sql, userId).Scan(&rootId)
    if err != nil {
	if errors.Is(err, pgx.ErrNoRows) {
	    return "", repository.ErrorRootDirNotFound
	}
	return "", fmt.Errorf("%w: %v", repository.ErrorGetDirFailed, err)
    }

    return rootId, nil
}

func (pr *FilesRepository) getNestedDirId(userId, path string) (string, error) {
    // $1 - userId, $2 - top dir name, $3 - depth
    recursiveSql := `
	WITH RECURSIVE dir_path AS (
	    SELECT id, name, parent_id, 1 AS depth
	    FROM directories
	    WHERE user_id = $1 AND parent_id IS NULL

	    UNION ALL

	    SELECT d.id, d.name, d.parent_id, dp.depth + 1
	    FROM directories d
	    JOIN dir_path dp ON d.parent_id = dp.id
	    WHERE d.user_id = $1
	)
	SELECT id
	FROM dir_path
	WHERE name = $2 
	  AND depth = $3;
    `

    segments := strings.Split(strings.Trim(path, "/"), "/")
    depth := len(segments)
    dirName := segments[depth-1]

    var dirId string
    err := pr.db.Pool.QueryRow(context.Background(), recursiveSql, userId, dirName, depth).Scan(&dirId)
    if err != nil {
	if errors.Is(err, pgx.ErrNoRows) {
	    return "", repository.ErrorDirectoryNotFound
	}
	return "", fmt.Errorf("%w: %v", repository.ErrorGetDirFailed, err)
    }

    return dirId, nil
}


func (pr *FilesRepository) fetchSubdirectories(userId, dirId, path string) ([]models.DirEntry, error) {
    const dirsSQL = `
	SELECT id, name
	FROM directories
	WHERE parent_id=$1 AND user_id=$2;
    `
    rows, err := pr.db.Pool.Query(context.Background(), dirsSQL, dirId, userId)
    if err != nil {
	return nil, fmt.Errorf("%w: %v", repository.ErrorQueryFailed, err)
    }
    defer rows.Close()

    var result []models.DirEntry
    for rows.Next() {
	var id, name string
	if err := rows.Scan(&id, &name); err != nil {
	    return nil, fmt.Errorf("%w: %v", repository.ErrorScanFailed, err)
	}
	result = append(result, models.DirEntry{
	    Id: id,
	    Type: "dir",
	    Name: name,
	    FullFilepath: joinPath(path, name, true),
	})
    }

    return result, nil
}

func (pr *FilesRepository) fetchFiles(userId, dirId, path string) ([]models.DirEntry, error) {
    const filesSql = `
	SELECT id, name
	FROM files
	WHERE directory_id=$1 AND user_id=$2;
    `
    rows, err := pr.db.Pool.Query(context.Background(), filesSql, dirId, userId)
    if err != nil {
	return nil, fmt.Errorf("%w: %v", repository.ErrorQueryFailed, err)
    }
    defer rows.Close()

    var result []models.DirEntry
    for rows.Next() {
	var id, name string
	if err := rows.Scan(&id, &name); err != nil {
	    return nil, fmt.Errorf("%w: %v", repository.ErrorScanFailed, err)
	}
	result = append(result, models.DirEntry{
	    Id: id,
	    Type: "file",
	    Name: name,
	    FullFilepath: joinPath(path, name, true),
	})
    }
    return result, nil
}

func joinPath(base, name string, isDir bool) string {
    if base == "/" {
	base = ""
    }
    if isDir {
	return base + "/" + name + "/"
    }
    return base + "/" + name
}

//////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func (pr *FilesRepository) CreateDir(userId, parentDirId, name string) (string, error) {
    ctx := context.Background()

    // check if parent directory exist and belongs to user
    var tmp string
    checkSql:= `SELECT id FROM directories WHERE id=$1 AND user_id=$2`
    err := pr.db.Pool.QueryRow(ctx, checkSql, parentDirId, userId).Scan(&tmp)
    if err == pgx.ErrNoRows {
	return "", repository.ErrorParentDirNotFound
    } else if err != nil {
	return "", fmt.Errorf("Check parent dir error: %w", err)
    }

    // check if directory exist
    dupCheckSql := "SELECT id FROM directories WHERE name=$1 AND parent_id=$2 AND user_id=$3"
    err = pr.db.Pool.QueryRow(ctx, dupCheckSql, name, parentDirId, userId).Scan(&tmp)
    if err == nil {
	return "", repository.ErrorDirectoryAlreadyExist
    } else if err != pgx.ErrNoRows {
	return "", fmt.Errorf("Check duplicate dir error: %w", err)
    }

    newDirId := uuid.NewString()
    insertSql := `
	INSERT INTO directories (id, name, parent_id, user_id)
	VALUES ($1, $2, $3, $4, $5)` 
    _, err = pr.db.Pool.Exec(ctx, insertSql, newDirId, name, parentDirId, userId)
    if err != nil {
	return "", fmt.Errorf("Create dir insert error: %w", err)
    }

    return newDirId, nil
}

func (pr *FilesRepository) CreateRootDir(userId string) (string, error) {
    ctx := context.Background()

    // check if directory exist
    var tmp string
    dupCheckSql := "SELECT id FROM directories WHERE user_id=$1 AND parent_id=NULL"
    err := pr.db.Pool.QueryRow(ctx, dupCheckSql, userId).Scan(&tmp)
    if err == nil {
	return "", repository.ErrorDirectoryAlreadyExist
    } else if err != pgx.ErrNoRows {
	return "", fmt.Errorf("Check duplicate dir error: %w", err)
    }

    newDirId := uuid.NewString()
    insertSql := `
	INSERT INTO directories (id, user_id, name, parent_id)
	VALUES ($1, $2, $3, NULL)` 
    _, err = pr.db.Pool.Exec(ctx, insertSql, newDirId, userId, "root")
    if err != nil {
	return "", fmt.Errorf("Create dir insert error: %w", err)
    }

    return newDirId, nil
}
