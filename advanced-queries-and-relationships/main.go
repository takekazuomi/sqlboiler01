package main

//go:generate sqlboiler --wipe mysql

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/takekazuomi/sqlboiler01/advanced-queries-and-relationships/gen/models"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
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

	doseIt, err := models.Videos(qm.Where("id = ?", 1)).Exists(context.Background(), db)
	dieIf(err)

	fmt.Println("dose it?", doseIt)

	user, err := models.Users().One(context.Background(), db)
	dieIf(err)

	fmt.Println("user:", user.ID, user.Name)

	nVideos, err := user.Videos().Count(context.Background(), db)
	dieIf(err)

	fmt.Println("  nVideos:", nVideos)

	videos, err := user.Videos(qm.Where("id < 5"), qm.And("id < 4"), qm.Limit(3)).All(context.Background(), db)
	dieIf(err)

	fmt.Print("videos:")

	for i, v := range videos {
		if i != 0 {
			fmt.Print(", ")
		}
		fmt.Print(v.Name)
	}
	fmt.Println()

	vid1, vid2 := &models.Video{Name: "a"}, &models.Video{Name: "b"}
	err = user.AddVideos(context.Background(), db, true, vid1, vid2)
	dieIf(err)

	fmt.Println("user", user.Name)
	nVideos, err = user.Videos().Count(context.Background(), db)
	dieIf(err)

	fmt.Println("  nVideos:", nVideos)

	tags, err := models.Tags().All(context.Background(), db)
	dieIf(err)

	fmt.Println("nTags:", len(tags))

	err = vid1.AddTags(context.Background(), db, false, tags...)
	dieIf(err)

	vid1Tags := func() int {
		nTags, err := vid1.Tags().Count(context.Background(), db)
		dieIf(err)
		return int(nTags)
	}

	fmt.Println("videos:", vid1.ID, vid1.Name)
	fmt.Println("  nTags:", vid1Tags())

	err = vid1.SetTags(context.Background(), db, false, tags[0])
	dieIf(err)

	fmt.Println("videos:", vid1.ID, vid1.Name)
	fmt.Println("  nTags:", vid1Tags())

	err = vid1.RemoveTags(context.Background(), db, tags...)
	dieIf(err)

	fmt.Println("videos:", vid1.ID, vid1.Name)
	fmt.Println("  nTags:", vid1Tags())

	// eager loading = n+1 problem

	user, err = models.Users(qm.Load("Videos.Tags")).One(context.Background(), db)
	dieIf(err)

	fmt.Println("user:", user.ID, user.Name)
	fmt.Println("  nVideos:", len(user.R.Videos))
	fmt.Println("  R.Videos[0]:", user.R.Videos[0].Name)
	fmt.Println("  R.Videos[0].R.User.Name:", user.R.Videos[0].R.User.Name)

	user, err = models.Users(
		qm.Load("Videos", qm.Limit(1)),
		qm.Load("Videos.Tags", qm.Limit(1)),
	).One(context.Background(), db)
	dieIf(err)

	fmt.Println("user:", user.ID, user.Name)
	fmt.Println("  nVideos:", len(user.R.Videos))
	fmt.Println("    Videos[0].nTags:", len(user.R.Videos[0].R.Tags))

	// raq query

	var users []*models.User

	rows, err := db.QueryContext(context.Background(), "select * from users")
	dieIf(err)

	err = queries.Bind(rows, &users)
	dieIf(err)
	defer rows.Close()

	fmt.Print("raw users:")
	for i, u := range users {
		if i != 0 {
			fmt.Print(", ")
		}
		fmt.Print(u.Name)
	}
	fmt.Println()

	// little bit higher level raw query

	users = nil
	err = models.Users(qm.SQL("select * from users")).Bind(context.Background(), db, &users)
	dieIf(err)

	fmt.Println("nUsers:", len(users))

	// join struct

	type joinStruct struct {
		User  models.User  `boil:"users,bind"`
		Video models.Video `boil:"videos,bind"`
	}

	var joins []*joinStruct

	rows, err = db.QueryContext(context.Background(),
		`select users.id as "users.id", users.name as "users.name", videos.id as "videos.id", videos.name as "videos.name" 
	from users inner join videos on users.id = videos.user_id`,
	)
	dieIf(err)

	err = queries.Bind(rows, &joins)
	dieIf(err)
	defer rows.Close()

	fmt.Println("joins:", len(joins))
	fmt.Println(joins[0].User.Name, joins[0].Video.Name)
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
