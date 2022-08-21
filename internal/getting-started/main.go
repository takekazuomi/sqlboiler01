package main

//go:generate sqlboiler --wipe mysql

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/takekazuomi/sqlboiler01/internal/getting-started/gen/models"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

func main() {
	dsn, ok := os.LookupEnv("DSN")
	dieFalse(ok, "env DSN not found")

	db, err := sql.Open("mysql", dsn)
	dieIf(err)

	err = db.Ping()
	dieIf(err)

	fmt.Println("connected")

	u := &models.User{Name: "john"}
	err = u.Insert(context.Background(), db, boil.Infer())
	dieIf(err)

	fmt.Println("user id", u.ID)

	got, err := models.Users().One(context.Background(), db)
	dieIf(err)

	fmt.Println("got user:", got.ID)

	found, err := models.FindUser(context.Background(), db, u.ID)
	dieIf(err)

	fmt.Println("found user:", found.ID)

	found.Name = "jane"
	rows, err := found.Update(context.Background(), db, boil.Whitelist(models.UserColumns.Name))
	dieIf(err)

	fmt.Println("updated", rows, "users")

	foundAgain, err := models.Users(qm.Where("name = ?", found.Name)).One(context.Background(), db)
	dieIf(err)
	fmt.Println("found again:", foundAgain.ID, foundAgain.Name)

	exists, err := models.UserExists(context.Background(), db, foundAgain.ID)
	dieIf(err)

	fmt.Println("user", foundAgain.ID, "exists?", exists)

	count, err := models.Users().Count(context.Background(), db)
	dieIf(err)

	fmt.Println("there are", count, "users")

	// 10:19
	videos, err := got.Videos().All(context.Background(), db)
	dieIf(err)

	for _, v := range videos {
		fmt.Println("  got video:", v.ID)
	}
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
