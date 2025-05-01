package files

import "github.com/lyhomyna/sf/file-service/database"

type ProductRepository struct {
    db *database.Postgres
}

func NewFilesRepository(db *database.Postgres) *ProductRepository {
    return &ProductRepository{
	db: db,
    }
}
