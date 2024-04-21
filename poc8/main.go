package main

//go:generate sqlboiler --wipe mysql

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/takekazuomi/sqlboiler01/poc8/gen/models"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

func main() {
	//boil.DebugMode = true
	ctx := context.Background()

	dsn, ok := os.LookupEnv("DSN")
	dieFalse(ok, "env DSN not found")

	db, err := sql.Open("mysql", dsn)
	dieIf(err)

	err = db.Ping()
	dieIf(err)

	fmt.Println("connected")

	u := &models.User{
		Name: "foo",
	}

	err = u.Insert(ctx, db, boil.Infer())
	dieIf(err)

	fmt.Printf("%v, %v, %v, %v, %v\n", u.ID, u.Name, u.CreatedAt.UnixNano(), u.UpdatedAt.UnixNano(), nullTimeUnixNano(u.DeletedAt))

	err = u.Reload(ctx, db)
	dieIf(err)

	fmt.Printf("Nano:\t%v, %v, %v, %v, %v\n", u.ID, u.Name, u.CreatedAt.UnixNano(), u.UpdatedAt.UnixNano(), nullTimeUnixNano(u.DeletedAt))
	fmt.Printf("Micro:\t%v, %v, %v, %v, %v\n", u.ID, u.Name, u.CreatedAt.UnixMicro(), u.UpdatedAt.UnixMicro(), nullTimeUnixNano(u.DeletedAt))
	fmt.Printf("Unix:\t%v, %v, %v, %v, %v\n", u.ID, u.Name, u.CreatedAt.Unix(), u.UpdatedAt.Unix(), nullTimeUnixNano(u.DeletedAt))

	// roundの確認
	// 1668122405569040000
	// 1668122405569040
	// 1668122405
	// 169040000
	// 100
	sec := int64(1668122405)
	nsec := int64(100)

	for i, j := 0, int64(0); i < 11; i++ {
		u := &models.User{
			Name:      "foo " + strconv.Itoa(i),
			CreatedAt: time.Unix(sec, j),
		}
		err = u.Insert(ctx, db, boil.Infer())
		dieIf(err)
		err = u.Reload(ctx, db)
		dieIf(err)

		// MySQLでは、datetime(6)としたときに、micro secでroundされる
		// https://dev.mysql.com/doc/refman/8.0/en/fractional-seconds.html#:~:text=Inserting%20a%20TIME%2C%20DATE%2C%20or%20TIMESTAMP%20value%20with%20a%20fractional%20seconds%20part%20into%20a%20column%20of%20the%20same%20type%20but%20having%20fewer%20fractional%20digits%20results%20in%20rounding.%20Consider%20a%20table%20created%20and%20populated%20as%20follows%3A
		fmt.Printf("%v: %v, %v, %v, %v, %v\n", i, j, u.ID, u.CreatedAt, u.CreatedAt.UnixNano(), u.CreatedAt.Round(time.Microsecond).UnixMicro())
		// i=5 の時に、繰り上げられて、1 micro secになる
		j += nsec
	}
}

func nullTimeUnixNano(t null.Time) int64 {
	if !t.Valid {
		return 0
	}
	return t.Time.UnixNano()
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
