package userImage

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/jackc/pgx/v5"
	"github.com/lyhomyna/sf/file-service/database"
)

type UserImagesRepository struct {
    db *database.Postgres
}

func NewUserImageRepository(db *database.Postgres) *UserImagesRepository {
    return &UserImagesRepository{
	db: db,
    }
}

// directory where uploaded files will be saved
var userImagesDirectoryPath = filepath.Join("userImages")

func (r *UserImagesRepository) GetUserImageDirectoryPath() string {
    return userImagesDirectoryPath
}

func (r *UserImagesRepository) SaveUserImage(userId string, imageUrl string) error {
    ctx := context.Background()

    sql := "UPDATE users SET image_url=$1 WHERE id=$2"

    _, err := r.db.Pool.Exec(ctx, sql, imageUrl, userId)
    if err != nil {
	log.Printf("Failed to update image for user %s: %v", userId, err)
	return errors.New("Failed to update image")
    }

    return nil
}

func (r *UserImagesRepository) GetUserImageUrl(userId string) (string, error) {
    ctx := context.Background()
    sql := "SELECT image_url FROM users WHERE id=$1"

    var imageUrl string
    err := r.db.Pool.QueryRow(ctx, sql, userId).Scan(&imageUrl)
    if err != nil {
	if errors.Is(err, pgx.ErrNoRows) {
	    log.Printf("No image for user %s", userId)
	    return "", errors.New("No image for user") 
	}

	log.Printf("Failed to scan row to get image for user %s: %v", userId, err)
	return "", errors.New("Failed to get image")
    }

    return imageUrl, nil
}

func (r *UserImagesRepository) InitUserImageDir() error {
    err := os.MkdirAll(userImagesDirectoryPath, os.ModePerm)
    return err
}

func (r *UserImagesRepository) RemoveImage(path string) error {
    err := os.Remove(path)
    if err == nil {
	err = fmt.Errorf("Urgent remove of the image '%s'. %w", path, err)
    }
    return err
}

func (r *UserImagesRepository) ReadImage(imagePath string) (*os.File, []byte, error) {
    image, err := os.Open(imagePath)
    if err != nil {
	return nil, []byte{}, fmt.Errorf("%w: %v", errors.New("User picture missing"), err)
    }

    buf := make([]byte, 512)
    _, err = image.Read(buf)
    if err != nil {
	return nil, []byte{}, fmt.Errorf("%w: %v", errors.New("Couldn't read user image"), err)
    }
    image.Seek(0, io.SeekStart)

    return image, buf, nil
}
