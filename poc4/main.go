package main

//go:generate sqlboiler --wipe mysql

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/takekazuomi/sqlboiler01/poc4/gen/models"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

func main() {
	boil.DebugMode = true

	dsn, ok := os.LookupEnv("DSN")
	dieFalse(ok, "env DSN not found")

	db, err := sql.Open("mysql", dsn)
	dieIf(err)

	err = db.Ping()
	dieIf(err)

	fmt.Println("connected")

	t := &models.Table1{
		Status: models.Table1StatusOrange,
	}

	err = t.Insert(context.Background(), db, boil.Infer())
	dieIf(err)

	t1 := &models.Table1{
		Status: models.Table1StatusOrange,
	}

	err = t1.Insert(context.Background(), db, boil.Infer())
	dieIf(err)

	n, err := t.Delete(context.Background(), db, false)
	dieIf(err)
	fmt.Println("deleted:", n)
}

func dieFalse(ok bool, msg string) {
	if !ok {
		panic(msg)
	}
}

func dieIf(err error) {
	if err != nil {
		panic(err)
	}
}
