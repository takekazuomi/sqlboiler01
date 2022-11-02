package main

//go:generate sqlboiler --wipe mysql

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/takekazuomi/sqlboiler01/poc7/gen/models"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

func main() {
	boil.DebugMode = true
	boil.SetLocation(time.Local)

	dsn, ok := os.LookupEnv("DSN")
	dieFalse(ok, "env DSN not found")

	db, err := sql.Open("mysql", dsn)
	dieIf(err)

	err = db.Ping()
	dieIf(err)

	fmt.Println("connected")

	ctx := context.Background()

	// upsert, ID指定無し
	// INSERT INTO `users` (`name`,`memo`,`created_at`,`updated_at`)
	// VALUES (?,?,?,?) ON DUPLICATE KEY UPDATE
	// `name` = VALUES(`name`),`memo` = VALUES(`memo`),`created_at` = VALUES(`created_at`),`updated_at` = VALUES(`updated_at`)
	// [taro  2022-11-01 23:44:22.878110492 +0000 UTC 2022-11-01 23:44:22.878110492 +0000 UTC]
	u := &models.User{
		Name: "taro jiro",
		Memo: "",
	}

	err = u.Upsert(ctx, db, boil.Blacklist(models.UserColumns.ID, models.UserColumns.CreatedAt), boil.Infer())
	dieIf(err)
	fmt.Printf("1 upsert\tuser ID:%v, Name:%q, Memo:%q, %v, %v\n", u.ID, u.Name, u.Memo, u.CreatedAt, u.UpdatedAt)

	// 読み直し
	//	select * from `users` where `id`=?
	err = u.Reload(ctx, db)
	dieIf(err)
	fmt.Printf("2 reload\tuser ID:%v, Name:%q, Memo:%q, %v, %v\n", u.ID, u.Name, u.Memo, u.CreatedAt, u.UpdatedAt)

	// select * from `users` where `id`=?
	// [1]
	u2, err := models.FindUser(context.Background(), db, u.ID)
	dieIf(err)
	fmt.Printf("3 find\tuser ID:%v, Name:%q, Memo:%q, %v, %v\n", u2.ID, u2.Name, u2.Memo, u2.CreatedAt, u2.UpdatedAt)

	// upsert, update側
	// INSERT INTO `users` (`id`,`name`,`memo`,`created_at`,`updated_at`) VALUES (?,?,?,?,?) ON DUPLICATE KEY UPDATE `name` = VALUES(`name`),`memo` = VALUES(`memo`),`created_at` = VALUES(`created_at`),`updated_at` = VALUES(`updated_at`)
	// [1 taro update memo 2022-11-02 08:44:23 +0900 JST 2022-11-01 23:44:22.887576297 +0000 UTC]
	u.Memo = "update memo 2"

	// updateの場合は、IDとCreateAtを外す
	err = u.Upsert(ctx, db, boil.Blacklist(models.UserColumns.ID, models.UserColumns.CreatedAt), boil.Infer())
	dieIf(err)
	fmt.Printf("4 upsert\tuser ID:%v, Name:%q, Memo:%q, %v, %v\n", u.ID, u.Name, u.Memo, u.CreatedAt, u.UpdatedAt)

	// 読み直し
	err = u.Reload(ctx, db)
	dieIf(err)
	fmt.Printf("5 reload\tuser ID:%v, Name:%q, Memo:%q, %v, %v\n", u.ID, u.Name, u.Memo, u.CreatedAt, u.UpdatedAt)

	// read all

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
