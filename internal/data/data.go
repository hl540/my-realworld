package data

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	"github.com/hl540/my-realworld/internal/conf"
	_ "github.com/mattn/go-sqlite3"
	"xorm.io/xorm"
	xormlog "xorm.io/xorm/log"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewData, NewXormClient, NewUserRepo, NewArticleRepo)

// Data .
type Data struct {
	db  *xorm.Engine
	log *log.Helper
}

func NewXormClient(conf *conf.Data, logger log.Logger) *xorm.Engine {
	log := log.NewHelper(logger)
	client, err := xorm.NewEngine(conf.GetDatabase().GetDriver(), conf.GetDatabase().GetSource())
	if err != nil {
		log.Fatalf("failed opening connection to db: %v", err)
	}
	err = client.Sync(&User{}, &Follow{}, &Article{}, &Tag{}, &Favorite{})
	if err != nil {
		log.Fatalf("failed opening connection to db: %v", err)
	}
	client.ShowSQL(true)
	client.Logger().SetLevel(xormlog.LOG_DEBUG)
	return client
}

// NewData .
func NewData(xormClient *xorm.Engine, logger log.Logger) (*Data, func(), error) {
	cleanup := func() {
		log.NewHelper(logger).Info("closing the data resources")
	}
	return &Data{
		db:  xormClient,
		log: log.NewHelper(logger),
	}, cleanup, nil
}
