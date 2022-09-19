package iftest

import (
	"context"
	"database/sql"

	"github.com/takekazuomi/sqlboiler01/poc3/gen/models"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

type Accessor[K comparable, V any] interface {
	Get(context.Context, *sql.DB, K) (*V, error)
	Set(context.Context, *sql.DB, *V) error
}

type UserAccessor struct{}

var _ Accessor[int, models.User] = UserAccessor{}

func (u UserAccessor) Get(ctx context.Context, db *sql.DB, key int) (*models.User, error) {
	return models.FindUser(ctx, db, key)
}

func (u UserAccessor) Set(ctx context.Context, db *sql.DB, user *models.User) error {
	_, err := user.Update(ctx, db, boil.Infer())
	return err
}

// generic accessor 的なのができるかと思ったけどできなかった。残念。
// func (u Accessor[K, V]) Get(ctx context.Context, db *sql.DB, key K) (*V, error) {
// 	return models.FindUser(ctx, db, key)
// }
