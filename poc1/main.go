package main

//go:generate sqlboiler --wipe mysql

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/oklog/ulid/v2"
	"github.com/takekazuomi/sqlboiler01/poc1/gen/models"
	pb "github.com/takekazuomi/sqlboiler01/poc1/pkg/apis/message/v1"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"google.golang.org/protobuf/types/known/timestamppb"
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

	tx, err := db.BeginTx(context.Background(), nil)
	dieIf(err)
	defer tx.Rollback()

	now := time.Now()
	id := newULID()
	r := &models.Resource{
		ID:        id,
		Status:    models.ResourcesStatusSTATUS_ACTIVE,
		CreatedAt: now,
		UpdatedAt: now,
		DeletedAt: null.Time{},
	}
	err = r.Insert(context.Background(), tx, boil.Infer())
	dieIf(err)

	p := models.Property{
		ResourceID: id,
		FirstName:  "taro",
		LastName:   "yamada",
	}
	err = p.Insert(context.Background(), tx, boil.Infer())
	dieIf(err)
	tx.Commit()

	fmt.Println("insert resources:", id)

	id = newULID()
	now = time.Now()
	pbr := &pb.Resource{
		Name: id.String(),
		Properties: &pb.Properties{
			FirstName: "hanako",
			LastName:  "yamada",
		},
		Status:     pb.Status_STATUS_ACTIVE,
		CreateTime: timestamppb.New(now),
		UpdateTime: timestamppb.New(now),
		DeleteTime: nil,
	}

	mr := NewModelResourceFromPb(pbr)
	mp := NewModelPropertyFromPb(id, pbr.Properties)

	tx, err = db.BeginTx(context.Background(), nil)
	dieIf(err)
	defer tx.Rollback()

	err = mr.Insert(context.Background(), tx, boil.Infer())
	dieIf(err)

	err = mp.Insert(context.Background(), tx, boil.Infer())
	dieIf(err)
	tx.Commit()

	fmt.Println("insert resources:", id)

}

func newULID() ulid.ULID {
	return ulid.MustNew(ulid.Now(), ulid.DefaultEntropy())
}

func newTimestamppbFromNullTime(t null.Time) *timestamppb.Timestamp {
	if t.Valid {
		return timestamppb.New(t.Time)
	}
	return nil
}

func newNullTimeFromTimestamppb(t *timestamppb.Timestamp) null.Time {
	if t == nil {
		return null.NewTime(time.Time{}, false)
	}
	return null.NewTime(t.AsTime(), true)
}

func NewPbPropertyFromModel(source *models.Property) *pb.Properties {
	return &pb.Properties{
		FirstName: source.FirstName,
		LastName:  source.LastName,
	}
}

func NewModelPropertyFromPb(parent ulid.ULID, source *pb.Properties) *models.Property {
	return &models.Property{
		ID:         0,
		ResourceID: parent,
		FirstName:  source.FirstName,
		LastName:   source.LastName,
	}
}

func NewModelResourceFromPb(source *pb.Resource) *models.Resource {
	return &models.Resource{
		ID:        ulid.MustParse(source.Name),
		Status:    models.ResourcesStatus(pb.Status_name[int32(source.Status)]),
		CreatedAt: source.CreateTime.AsTime(),
		UpdatedAt: source.UpdateTime.AsTime(),
		DeletedAt: newNullTimeFromTimestamppb(source.DeleteTime),
	}
}

func NewPbResourceFromModel(source *models.Resource) *pb.Resource {
	return &pb.Resource{
		Name: source.ID.String(),
		Properties: &pb.Properties{
			FirstName: source.R.Properties[0].FirstName,
			LastName:  source.R.Properties[0].LastName,
		},
		Status:     pb.Status(pb.Status_value[source.Status.String()]),
		CreateTime: timestamppb.New(source.CreatedAt),
		UpdateTime: timestamppb.New(source.UpdatedAt),
		DeleteTime: newTimestamppbFromNullTime(source.DeletedAt),
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
