package main

//go:generate sqlboiler --wipe mysql

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	//boil.DebugMode = true
	_ = context.Background()

	dsn, ok := os.LookupEnv("DSN")
	dieFalse(ok, "env DSN not found")

	db, err := sql.Open("mysql", dsn)
	dieIf(err)

	err = db.Ping()
	dieIf(err)

	fmt.Println("connected")

}

func dieIf(err error) {
	if err != nil {
		panic(err)
	}
}

func dieFalse(ok bool, msg string) {
	if !ok {
		panic(msg)
	}
}
