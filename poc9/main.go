package main

//go:generate sqlboiler --wipe mysql

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/takekazuomi/sqlboiler01/poc9/gen/models"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

func createData(db *sql.DB, count int) {
	// start tx
	tx, err := db.Begin()
	dieIf(err)

	defer func() {
		_ = tx.Rollback()
	}()

	// prepare data template
	table1 := models.Table1{}

	nameFmt := dummyName(1024)

	// insert a row
	for i := 0; i < count; i++ {
		table1.Name = fmt.Sprintf(nameFmt, i)
		//slog.Info("before", slog.Any("table1", table1))
		// ID is auto increment, so we need to exclude it from insert
		err = table1.Insert(context.Background(), tx, boil.Blacklist(models.Table1Columns.ID))
		//slog.Info("inserted", slog.Any("table1", table1))
		dieIf(err)
	}
	err = tx.Commit()
	dieIf(err)
}

const dummy = "dummy data"

func dummyName(n int) string {
	var name strings.Builder

	j := (n - 10 + 1) / len(dummy)
	name.Grow(n)
	name.WriteString("%10d:")
	for i := 0; i < j; i++ {
		name.WriteString(dummy)
	}
	return name.String()
}

func connectDB() *sql.DB {
	dsn, ok := os.LookupEnv("DSN")
	dieFalse(ok, "env DSN not found")

	db, err := sql.Open("mysql", dsn)
	dieIf(err)

	err = db.Ping()
	dieIf(err)

	return db
}

func allTable1(db *sql.DB, limit int64) int {
	table1s, err := models.Table1s(
		qm.Limit(int(limit)),
	).All(context.Background(), db)
	dieIf(err)

	for _, v := range table1s {
		_ = v
	}
	return len(table1s)
}

func main() {
	boil.DebugMode = true
	_ = context.Background()

	db := connectDB()

	fmt.Println("connected")

	createData(db, 10)
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
