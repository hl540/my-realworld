package data

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/hl540/my-realworld/internal/biz"
	entuser "github.com/hl540/my-realworld/internal/data/ent/user"
)

type userRepo struct {
	data *Data
	log  *log.Helper
}

func NewUserRepo(data *Data, logger log.Logger) biz.UserRepo {
	return &userRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (u *userRepo) AddUser(ctx context.Context, user *biz.User) (*biz.User, error) {
	po, err := u.data.db.User.Create().
		SetUsername(user.Username).
		SetPassword(user.PassWord).
		SetEmail(user.Email).
		SetImage(user.Image).Save(ctx)
	if err != nil {
		return nil, err
	}
	return &biz.User{
		Username: po.Username,
		PassWord: po.Password,
		Email:    po.Email,
		Image:    po.Image,
	}, nil
}

func (u *userRepo) SaveUser(ctx context.Context, user *biz.User) (*biz.User, error) {
	po, err := u.data.db.User.Query().Where(entuser.Username(user.Username)).First(ctx)
	if err != nil {
		return nil, err
	}
	_, err = u.data.db.User.Update().
		SetPassword(user.PassWord).
		SetPassword(user.Email).
		SetPassword(user.Image).
		Where(entuser.ID(po.ID)).Save(ctx)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (u *userRepo) GetUser(ctx context.Context, userName string) (*biz.User, error) {
	po, err := u.data.db.User.Query().Where(entuser.Username(userName)).First(ctx)
	if err != nil {
		return nil, err
	}
	return &biz.User{
		Username: po.Username,
		PassWord: po.Password,
		Email:    po.Email,
		Image:    po.Image,
	}, nil
}
