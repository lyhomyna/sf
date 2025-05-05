package database

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/lyhomyna/sf/file-service/config"
)

type Postgres struct {
    Pool *pgxpool.Pool 
}

var db *Postgres

func GetPostgresDb(cfg *config.Config) (*Postgres) {
    ctx := context.Background()
    if db == nil {
	connPool, err := pgxpool.New(ctx, cfg.PostgresConfig.GetConnString())
	if err != nil {
	    log.Printf("Failed to connect to DB.\n%s", err)
	    panic(err)
	}

	db = &Postgres {
	    Pool: connPool,
	}
    }

    if err := db.Pool.Ping(ctx); err != nil {
	log.Printf("Connection to DB established, but hadn't got ping response.\n%s", err)
	panic(err)
    }

    log.Println("Connection to Postgres Db established")
    return db 
}
