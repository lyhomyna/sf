package database

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

var connPool *pgxpool.Pool

func ConnectToDb(ctx context.Context) (*pgxpool.Pool) {
    if connPool == nil {
	connString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", os.Getenv("POSTGRES_USER"), os.Getenv("POSTGRES_PASSWORD"), os.Getenv("POSTGRES_HOST"), os.Getenv("POSTGRES_PORT"), os.Getenv("POSTGRES_NAME"))

	var err error
	connPool, err = pgxpool.New(ctx, connString)
	if err != nil {
	    log.Printf("Failed connect to DB.\n%s", err)
	    return nil
	}

	if err := connPool.Ping(ctx); err != nil {
	    log.Printf("Connection to DB established, but haven't got ping response.\n%s", err)
	    return nil
	}

	log.Println("Connection to postgres DB established.")
    }

    return connPool 
}
