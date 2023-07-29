package data

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	"github.com/hl540/my-realworld/internal/conf"
	_ "github.com/mattn/go-sqlite3"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewData, NewGormClient, NewUserRepo, NewArticleRepo, NewTagRepo)

// Data .
type Data struct {
	db  *gorm.DB
	log *log.Helper
}

func NewGormClient(conf *conf.Data, logger log.Logger) *gorm.DB {
	log := log.NewHelper(logger)
	client, err := gorm.Open(
		sqlite.Open(conf.Database.Driver),
		&gorm.Config{
			Logger: gormlogger.Default.LogMode(gormlogger.Info),
		},
	)
	if err != nil {
		log.Fatalf("failed opening connection to db: %v", err)
	}

	err = client.AutoMigrate(
		&User{},
		&Follow{},
		&Article{},
		&Tag{},
		&Favorite{},
	)
	if err != nil {
		log.Fatalf("failed creating schema resources: %v", err)
	}
	return client
}

// NewData .
func NewData(gormClient *gorm.DB, logger log.Logger) (*Data, func(), error) {
	cleanup := func() {
		log.NewHelper(logger).Info("closing the data resources")
	}
	return &Data{
		db:  gormClient,
		log: log.NewHelper(logger),
	}, cleanup, nil
}
