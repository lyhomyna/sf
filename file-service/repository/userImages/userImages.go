package userImages

import "github.com/lyhomyna/sf/file-service/database"

type UserImagesRepository struct {
    db *database.Postgres
}

func NewUserImagesRepository(db *database.Postgres) *UserImagesRepository {
    return &UserImagesRepository{
	db: db,
    }
}

func (r *UserImagesRepository) SaveUserImage(userId string, imageName string) error {
    panic("Not yet implemented")
}
