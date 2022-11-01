package main

//go:generate sqlboiler --wipe mysql

import (
	"context"
	"database/sql"
	"fmt"
	"math/rand"
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
	num := rand.Int31()

	fmt.Printf("ulid: %s, num:%d\n", ulid, num)

	t := &models.Table1{
		ID:  ulid,
		Num: int(num),
		F1:  "memo",
	}

	err = t.Insert(ctx, db, boil.Infer())
	dieIf(err)

	t1, err := models.Table1s(models.Table1Where.ID.EQ(ulid)).One(ctx, db)
	dieIf(err)

	fmt.Printf("t1: %v\n", t1)

	ulid2 := newULID()
	fmt.Printf("ulid: %s\n", ulid2)
	t2 := &models.Table2{
		ID:       ulid2,
		Table1ID: ulid,
		F2:       "table2 memo",
	}
	err = t2.Insert(ctx, db, boil.Infer())
	dieIf(err)

	fmt.Printf("t2: %v\n", t2)

	// num でjoinして持ってくる

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
