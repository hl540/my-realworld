package data

import (
	"context"
	"github.com/hl540/my-realworld/internal/conf"
	"github.com/hl540/my-realworld/internal/data/ent"
	"github.com/hl540/my-realworld/internal/data/ent/migrate"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	_ "github.com/mattn/go-sqlite3"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewData, NewEntClient, NewUserRepo)

// Data .
type Data struct {
	db  *ent.Client
	log *log.Helper
}

func NewEntClient(conf *conf.Data, logger log.Logger) *ent.Client {
	log := log.NewHelper(logger)
	client, err := ent.Open(conf.Database.Driver, conf.Database.Source)
	if err != nil {
		log.Fatalf("failed opening connection to db: %v", err)
	}
	err = client.Schema.Create(context.Background(), migrate.WithForeignKeys(false))
	if err != nil {
		log.Fatalf("failed creating schema resources: %v", err)
	}
	return client
}

// NewData .
func NewData(endClient *ent.Client, logger log.Logger) (*Data, func(), error) {
	cleanup := func() {
		log.NewHelper(logger).Info("closing the data resources")
	}
	return &Data{
		db:  endClient,
		log: log.NewHelper(logger),
	}, cleanup, nil
}
