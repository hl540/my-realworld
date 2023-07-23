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
		SetImage(user.Image).
		SetBio(user.Bio).Save(ctx)
	if err != nil {
		return nil, err
	}
	return &biz.User{
		Id:       po.ID,
		Username: po.Username,
		PassWord: po.Password,
		Email:    po.Email,
		Image:    po.Image,
	}, nil
}

func (u *userRepo) SaveUser(ctx context.Context, user *biz.User) (*biz.User, error) {
	up := u.data.db.User.Update()
	if user.PassWord != "" {
		up = up.SetPassword(user.PassWord)
	}
	if user.Username != "" {
		up = up.SetUsername(user.Username)
	}
	if user.Email != "" {
		up = up.SetEmail(user.Email)
	}
	up = up.SetImage(user.Image)
	up = up.SetBio(user.Bio)
	up = up.Where(entuser.ID(user.Id))
	_, err := up.Save(ctx)
	if err != nil {
		return nil, err
	}
	return u.GetUserByID(ctx, user.Id)
}

func (u *userRepo) GetUserByUsername(ctx context.Context, username string) (*biz.User, error) {
	po, err := u.data.db.User.Query().Where(entuser.Username(username)).First(ctx)
	if err != nil {
		return nil, err
	}
	return &biz.User{
		Id:       po.ID,
		Username: po.Username,
		PassWord: po.Password,
		Email:    po.Email,
		Image:    po.Image,
		Bio:      po.Bio,
	}, nil
}

func (u *userRepo) GetUserByEmail(ctx context.Context, email string) (*biz.User, error) {
	po, err := u.data.db.User.Query().Where(entuser.Email(email)).First(ctx)
	if err != nil {
		return nil, err
	}
	return &biz.User{
		Id:       po.ID,
		Username: po.Username,
		PassWord: po.Password,
		Email:    po.Email,
		Image:    po.Image,
		Bio:      po.Bio,
	}, nil
}

func (u *userRepo) GetUserByID(ctx context.Context, id int) (*biz.User, error) {
	po, err := u.data.db.User.Query().Where(entuser.ID(id)).First(ctx)
	if err != nil {
		return nil, err
	}
	return &biz.User{
		Id:       po.ID,
		Username: po.Username,
		PassWord: po.Password,
		Email:    po.Email,
		Image:    po.Image,
		Bio:      po.Bio,
	}, nil
}
