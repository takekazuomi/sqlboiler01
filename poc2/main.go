package main

//go:generate sqlboiler --wipe mysql

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/takekazuomi/sqlboiler01/poc2/gen/models"
	"github.com/volatiletech/null/v8"
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

	// upsert, ID指定無し
	// INSERT INTO `users` (`name`,`memo`) VALUES ('taro','') ON DUPLICATE KEY UPDATE `name` = VALUES(`name`),`memo` = VALUES(`memo`)
	u := &models.User{
		Name: "taro",
		Memo: null.StringFrom(""),
	}

	err = u.Upsert(context.Background(), db, boil.Infer(), boil.Infer())
	dieIf(err)
	fmt.Printf("user ID:%v, Name:%q, Memo:%q\n", u.ID, u.Name, u.Memo.String)

	// メモ更新、同じuserのインスタンスを使って、memoを更新
	// INSERT INTO `users` (`id`,`name`,`memo`) VALUES (3,'taro','hello world') ON DUPLICATE KEY UPDATE `name` = VALUES(`name`),`memo` = VALUES(`memo`)
	u.Memo = null.StringFrom("hello world")
	err = u.Upsert(context.Background(), db, boil.Infer(), boil.Infer())
	dieIf(err)
	fmt.Printf("user ID:%v, Name:%q, Memo:%q\n", u.ID, u.Name, u.Memo.String)

	// 別のuserのインスタンスを使って(ID無しで）Upsert。IDはauto incで振られて、taroはuniqなのでupdateされて、jiroになる
	// INSERT INTO `users` (`name`,`memo`) VALUES (?,?) ON DUPLICATE KEY UPDATE `name` = VALUES(`name`),`memo` = VALUES(`memo`)
	u1 := &models.User{
		Name: "taro",
		Memo: null.StringFrom("jiro"),
	}
	err = u1.Upsert(context.Background(), db, boil.Infer(), boil.Infer())
	dieIf(err)
	fmt.Printf("user ID:%v, Name:%q, Memo:%q\n", u1.ID, u1.Name, u1.Memo.String)

	// 同様に、yamada で新規、InsertになってIDが着く
	// INSERT INTO `users` (`name`,`memo`) VALUES (?,?) ON DUPLICATE KEY UPDATE `name` = VALUES(`name`),`memo` = VALUES(`memo`)
	// SELECT `id` FROM `users` WHERE `name`='yamada'
	// SELECTが実行されるのはなぜか？
	u2 := &models.User{
		Name: "yamada",
		Memo: null.StringFrom("jiro"),
	}
	err = u2.Upsert(context.Background(), db, boil.Infer(), boil.Infer())
	dieIf(err)
	fmt.Printf("user ID:%v, Name:%q, Memo:%q\n", u2.ID, u2.Name, u2.Memo.String)

	// yamadaのIDで、Nameをsatoにする。ID重複でinsert出来ないので、updateになって、yamadaが、sato hankoになる。
	// INSERT INTO `users` (`id`,`name`,`memo`) VALUES (5,'sato','hanako') ON DUPLICATE KEY UPDATE `name` = VALUES(`name`),`memo` = VALUES(`memo`)
	// 2回め動かすと、下記エラーになる
	// models: unable to upsert for users: Error 1062: Duplicate entry 'sato' for key 'users.unique_name'
	// yamadaをsatoに更新しようとするは、satoは既にいるので更新できない。
	// INSERT INTO `users` (`id`,`name`,`memo`) VALUES (36,'sato','hanako') ON DUPLICATE KEY UPDATE `name` = VALUES(`name`),`memo` = VALUES(`memo`)
	// ここでは、IDと、Nameを指定しているが、このデータは重複で作れないのでエラー
	u3 := &models.User{
		ID:   u2.ID,
		Name: "sato",
		Memo: null.StringFrom("hanako"),
	}
	err = u3.Upsert(context.Background(), db, boil.Infer(), boil.Infer())
	dieIf(err)
	fmt.Printf("user ID:%v, Name:%q, Memo:%q\n", u3.ID, u3.Name, u3.Memo.String)

	// ID無しで、Nameをsatoにする。satoは重複なので、insertはエラーになり、hanako が、taroに更新されることを期待する
	// INSERT INTO `users` (`name`,`memo`) VALUES ('sato','taro') ON DUPLICATE KEY UPDATE `name` = VALUES(`name`),`memo` = VALUES(`memo`)
	u4 := &models.User{
		Name: "sato",
		Memo: null.StringFrom("taro"),
	}
	err = u4.Upsert(context.Background(), db, boil.Infer(), boil.Infer())
	dieIf(err)
	fmt.Printf("user ID:%v, Name:%q, Memo:%q\n", u4.ID, u4.Name, u4.Memo.String)

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
