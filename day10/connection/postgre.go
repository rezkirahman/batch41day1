package connection

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
)


var Conn *pgx.Conn

func DatabaseConnect() {
	var err error

	// urlExample := "postgres://username:password@localhost:5432/database_name"
	databaseUrl := "postgres://postgres:default@localhost:5432/personalDB"
	fmt.Println(err)

	Conn, err = pgx.Connect(context.Background(), databaseUrl)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Success connect to database")
}