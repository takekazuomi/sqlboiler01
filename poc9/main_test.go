package main

import (
	"context"
	"log/slog"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/takekazuomi/sqlboiler01/poc5/gen/models"
)

func Benchmark_createData(b *testing.B) {
	//b.Skip("Skip Benchmark_createData()")

	slog.Info("Benchmark_createData()", slog.Int("b.N", b.N))

	db := connectDB()
	defer db.Close()

	createData(db, b.N)

	c, err := models.Table1s().Count(context.Background(), db)
	dieIf(err)

	slog.Info("Table1s", slog.Int64("count", c))

}

func Benchmark_allTable1(b *testing.B) {
	slog.Info("Benchmark_allTable1()", slog.Int("b.N", b.N))

	db := connectDB()
	defer db.Close()

	c, err := models.Table1s().Count(context.Background(), db)
	dieIf(err)

	slog.Info("Table1s", slog.Int64("count", c))

	allTable1(db, c)

}
