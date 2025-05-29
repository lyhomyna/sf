package userImage

import (
	"context"
	"errors"
	"log"

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
