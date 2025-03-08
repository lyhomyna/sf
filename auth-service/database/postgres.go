package database 

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5"
)

var dbConnection *pgx.Conn

func ConnectToDb(ctx context.Context) (*pgx.Conn) {
    if dbConnection == nil {
	connString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_NAME"))

	var err error
	dbConnection, err = pgx.Connect(ctx, connString)
	if err != nil {
	    log.Printf("Failed connect to DB.\n%s", err)
	    return nil
	}

	if err := dbConnection.Ping(ctx); err != nil {
	    log.Printf("Connection to DB established, but haven't got ping response.\n%s", err)
	    return nil
	}

	log.Println("Connection to postgres DB established.")
    }

    return dbConnection
}
