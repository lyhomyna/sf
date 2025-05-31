package dir

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/lyhomyna/sf/file-service/database"
	"github.com/lyhomyna/sf/file-service/models"
	"github.com/lyhomyna/sf/file-service/repository"
)

type DirsRepository struct {
    db *database.Postgres
}

func NewDirRepository(db *database.Postgres) *DirsRepository {
    return &DirsRepository{
	db: db,
    }
}

// Directory, where uploaded files will be saved
var filesDirectory = filepath.Join("files")

func (dr *DirsRepository) ListDir(path, userId string) ([]models.DirEntry, error) {

    dirId, err := dr.GetDirIdByPath(path, userId)
    if err != nil {
	return nil, fmt.Errorf("%w: %v", repository.ErrorPathNotFound, err)
    }

    // Directories
    dirs, err := dr.fetchSubdirectories(userId, dirId, path)
    if err != nil {
	return nil, err
    }

    // Files
    files, err := dr.fetchFiles(userId, dirId, path)
    if err != nil {
	return nil, err
    }
    
    entries := append(dirs, files...)

    return entries, nil
}

// path should be like this: /inner/directory
func (dr *DirsRepository) GetDirIdByPath(path, userId string) (string, error) {
    if path == "/" {
	return dr.getRootDirId(userId)
    }
    return dr.getNestedDirId(userId, path)
}

func (dr *DirsRepository) getRootDirId(userId string) (string, error) {
    sql := "SELECT id FROM DIRECTORIES WHERE user_id=$1 AND parent_id IS NULL"

    var rootId string
    err := dr.db.Pool.QueryRow(context.Background(), sql, userId).Scan(&rootId)
    if err != nil {
	if errors.Is(err, pgx.ErrNoRows) {
	    return "", repository.ErrorRootDirNotFound
	}
	return "", fmt.Errorf("%w: %v", repository.ErrorGetDirFailed, err)
    }

    return rootId, nil
}

func (dr *DirsRepository) getNestedDirId(userId, path string) (string, error) {
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
    depth := len(segments) + 1 // because of the above trim
    dirName := segments[len(segments)-1]

    var dirId string
    err := dr.db.Pool.QueryRow(context.Background(), recursiveSql, userId, dirName, depth).Scan(&dirId)
    if err != nil {
	if errors.Is(err, pgx.ErrNoRows) {
	    return "", repository.ErrorDirectoryNotFound
	}
	return "", fmt.Errorf("%w: %v", repository.ErrorGetDirFailed, err)
    }

    return dirId, nil
}


func (dr *DirsRepository) fetchSubdirectories(userId, dirId, path string) ([]models.DirEntry, error) {
    const dirsSQL = `
	SELECT id, name
	FROM directories
	WHERE parent_id=$1 AND user_id=$2;
    `
    rows, err := dr.db.Pool.Query(context.Background(), dirsSQL, dirId, userId)
    if err != nil {
	return nil, fmt.Errorf("%w: %v", repository.ErrorQueryFailed, err)
    }
    defer rows.Close()

    result := []models.DirEntry{}
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

func (dr *DirsRepository) fetchFiles(userId, dirId, path string) ([]models.DirEntry, error) {
    const filesSql = `
	SELECT id, filename
	FROM files
	WHERE directory_id=$1 AND user_id=$2;
    `
    rows, err := dr.db.Pool.Query(context.Background(), filesSql, dirId, userId)
    if err != nil {
	return nil, fmt.Errorf("%w: %v", repository.ErrorQueryFailed, err)
    }
    defer rows.Close()

    result := []models.DirEntry{}
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

// :::::::::::::::::::::::::::::::::::::::::::::::::::::::::

func (dr *DirsRepository) CreateDir(userId, parentDirId, parentDirPath, name string) (string, error) {
    ctx := context.Background()

    // check if parent directory exist and belongs to user
    var tmp string
    checkSql:= `SELECT id FROM directories WHERE id=$1 AND user_id=$2`
    err := dr.db.Pool.QueryRow(ctx, checkSql, parentDirId, userId).Scan(&tmp)
    if err == pgx.ErrNoRows {
	return "", repository.ErrorParentDirNotFound
    } else if err != nil {
	return "", fmt.Errorf("Check parent dir error: %w", err)
    }

    // check if directory exist
    dupCheckSql := "SELECT id FROM directories WHERE name=$1 AND parent_id=$2 AND user_id=$3"
    err = dr.db.Pool.QueryRow(ctx, dupCheckSql, name, parentDirId, userId).Scan(&tmp)
    if err == nil {
	return "", repository.ErrorDirectoryAlreadyExist
    } else if err != pgx.ErrNoRows {
	return "", fmt.Errorf("Check duplicate dir error: %w", err)
    }

    newFilePath := filepath.Join(filesDirectory, userId, parentDirPath, name)
    if err := os.MkdirAll(newFilePath, os.ModePerm); err != nil {
	return "", fmt.Errorf("%w: %v", repository.FilesErrorInternal, err) 
    }

    newDirId := uuid.NewString()
    insertSql := `
	INSERT INTO directories (id, name, parent_id, user_id)
	VALUES ($1, $2, $3, $4)` 
    _, err = dr.db.Pool.Exec(ctx, insertSql, newDirId, name, parentDirId, userId)
    if err != nil {
	return "", fmt.Errorf("Create dir insert error: %w", err)
    }

    return newDirId, nil
}

// :::::::::::::::::::::::::::::::::::::::::::::::::::::::::

func (dr *DirsRepository) DeleteDir(userId, dirId string) (error) {
    ctx := context.Background()

    dirPath, err := dr.getDirPathById(userId, dirId)
    if err != nil {
	return err
    }

    tx, err := dr.db.Pool.Begin(ctx)
    if err != nil {
	return fmt.Errorf("%w: %v", repository.ErrorBeginTransaction, err)
    }
    defer tx.Rollback(ctx)

    recursiveSql := `
	WITH RECURSIVE dirs_to_delete AS (
	    SELECT id FROM directories WHERE id=$1 AND user_id=$2
	    UNION
	    SELECT d.id
	    FROM directories d
	    INNER JOIN dirs_to_delete dt ON d.parent_id=dt.id
	),
	deleted_files AS (
	    DELETE FROM files
	    WHERE directory_id IN (
		SELECT id FROM dirs_to_delete
	    )
	)
	DELETE FROM directories
	WHERE id IN (SELECT id FROM dirs_to_delete); `

    _, err = dr.db.Pool.Exec(ctx, recursiveSql, dirId, userId)
    if err != nil {
	return fmt.Errorf("%w: %v", repository.ErrorDelete, err)
    }

    err = os.RemoveAll(dirPath)
    if err != nil {
	return fmt.Errorf("%w: %v", repository.ErrorDeleteFolder, err)
    }

    if err := tx.Commit(ctx); err != nil {
	return fmt.Errorf("%w: %v", repository.ErrorCommitTransaction, err)
    }

    return nil
}

func (dr *DirsRepository) getDirPathById(userId, dirId string) (string, error) {
    ctx := context.Background()

    recursiveSql := `
	WITH RECURSIVE dir_path AS (
	    SELECT id, name, parent_id 
	    FROM directories
	    WHERE user_id=$1
	    AND id=$2
	    
	    UNION ALL

	    SELECT d.id, d.name, d.parent_id
	    FROM directories d
	    INNER JOIN dir_path dp ON d.id=dp.parent_id
	)
	SELECT name, parent_id 
	FROM dir_path; `

	rows, err := dr.db.Pool.Query(ctx, recursiveSql, userId, dirId)
	if err != nil {
	    if errors.Is(err, pgx.ErrNoRows) {
		return "", fmt.Errorf("%w: %v", repository.ErrorNoRows, err)
	    }

	    return "", fmt.Errorf("%w: %v", repository.FilesErrorDbQuery, err)
	}

	var dirNames []string
	for rows.Next() {
	    var dirName string
	    var parentId sql.NullString 

	    err := rows.Scan(&dirName, &parentId)
	    if err != nil {
		return "", fmt.Errorf("%w: %v", repository.ErrorScanFailed, err)
	    }

	    // root folder has "root" name, but in filepath that name is't relevant
	    if !parentId.Valid {
		continue
	    }

	    dirNames = append(dirNames, dirName)
	}

	pathParts := []string{"files", userId}
	for i := len(dirNames)-1; i >= 0; i-- {
	    pathParts = append(pathParts, dirNames[i])
	}
	return filepath.Join(pathParts...), nil
}

// :::::::::::::::::::::::::::::::::::::::::::::::::::::::::

func (dr *DirsRepository) CreateRootDir(userId string) (string, error) {
    ctx := context.Background()

    // check if directory exist
    var tmp string
    dupCheckSql := "SELECT id FROM directories WHERE user_id=$1 AND parent_id IS NULL"
    err := dr.db.Pool.QueryRow(ctx, dupCheckSql, userId).Scan(&tmp)
    if err == nil {
	return "", repository.ErrorDirectoryAlreadyExist
    } else if err != pgx.ErrNoRows {
	return "", fmt.Errorf("Check duplicate dir error: %w", err)
    }

    newFilePath := filepath.Join(filesDirectory, userId)
    if err := os.MkdirAll(newFilePath, os.ModePerm); err != nil {
	return "", fmt.Errorf("%w: %v", repository.FilesErrorInternal, err) 
    }

    newDirId := uuid.NewString()
    insertSql := `
	INSERT INTO directories (id, user_id, name, parent_id)
	VALUES ($1, $2, $3, NULL)` 
    _, err = dr.db.Pool.Exec(ctx, insertSql, newDirId, userId, "root")
    if err != nil {
	return "", fmt.Errorf("Create dir insert error: %w", err)
    }

    return newDirId, nil
}
