package main

//go:generate sqlboiler --wipe mysql

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/oklog/ulid/v2"
	"github.com/takekazuomi/sqlboiler01/poc5/gen/models"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

func main() {
	boil.DebugMode = true
	ctx := context.Background()

	dsn, ok := os.LookupEnv("DSN")
	dieFalse(ok, "env DSN not found")

	db, err := sql.Open("mysql", dsn)
	dieIf(err)

	err = db.Ping()
	dieIf(err)

	fmt.Println("connected")

	ulid := newULID()
	fmt.Printf("ulid: %s\n", ulid)

	t := &models.Table1{
		ID: ulid,
	}

	// データ作る
	err = t.Insert(ctx, db, boil.Infer())
	dieIf(err)

	t1, err := models.Table1s(models.Table1Where.ID.EQ(ulid)).One(ctx, db)
	dieIf(err)
	fmt.Printf("t1: %#v\n", t1)

	// 存在しないデータをOneする
	ulid2 := newULID()
	fmt.Printf("ulid2: %s\n", ulid)

	t2, err := models.Table1s(models.Table1Where.ID.EQ(ulid2)).One(ctx, db)
	fmt.Printf("t2: %#v. err:%v\n", t2, err)

	// 存在しないデータをAllする
	t3, err := models.Table1s(models.Table1Where.ID.EQ(ulid2)).All(ctx, db)
	fmt.Printf("t3: %#v. err:%v\n", t3, err)

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

func newULID() ulid.ULID {
	return ulid.MustNew(ulid.Now(), ulid.DefaultEntropy())
}
