# SQLBoiler - Advanced Queries and Relationships

Example code of [SQLBoiler: Advanced Queries and Relationships](https://www.youtube.com/watch?v=iiJuM9NR8No).

schemaは、getting started と同じ

- users ユーザー
- videos ユーザーが見た英語
- tags videoのタグ
- tags_videos tagとvideoの交差エンティティ


## 関心した機能

VideoにTagを追加すると、
```go
err = vid1.AddTags(context.Background(), db, false, tags...)
```

tags_videosに行が追加される。
```sql
insert into `tags_videos` (`video_id`, `tag_id`) values (?, ?)
```

Videoのタグを数えると、
```go
nTags, err := vid1.Tags().Count(context.Background(), db)
```

tags_videosの該当video_idの行を数える。

```sql
SELECT COUNT(*) FROM `tags` INNER JOIN `tags_videos` on `tags`.`id` = `tags_videos`.`tag_id` WHERE (`tags_videos`.`video_id`=?);
```

前のほうで、Userが見たVideoを追加したとき、

```go
err = user.AddVideos(context.Background(), db, true, vid1, vid2)
```

videosの行が追加されてた。
```sql
INSERT INTO `videos` (`user_id`,`name`) VALUES (?,?)
```

試しに、VideoにTagを追加するときに、`insert bool` を、trueで実行
```go
err = vid1.AddTags(context.Background(), db, true, tags...)
```

```sql
INSERT INTO `tags` (`id`,`name`) VALUES (?,?)
[1 action]
panic: failed to insert into foreign table: models: unable to insert into tags: Error 1062: Duplicate entry '1' for key 'tags.PRIMARY'
```

既にtagsにある行をinsertしようとして、Duplicateのエラーになった。insertするかどうかをしてするようになっているので、直感的ではないが分かるような仕様。



## あとで

`maxBadConnRetries` というのがある。コネクションプールに有効なコネクションがない場合はどうなるのか？

```go:sql.go:go/src/database/sql/sql.go
// ExecContext executes a query without returning any rows.
// The args are for any placeholder parameters in the query.
func (db *DB) ExecContext(ctx context.Context, query string, args ...any) (Result, error) {
	var res Result
	var err error
	var isBadConn bool
	for i := 0; i < maxBadConnRetries; i++ {
		res, err = db.exec(ctx, query, args, cachedOrNewConn)
		isBadConn = errors.Is(err, driver.ErrBadConn)
		if !isBadConn {
			break
		}
	}
	if isBadConn {
		return db.exec(ctx, query, args, alwaysNewConn)
	}
	return res, err
}
```

