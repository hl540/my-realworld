package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/hl540/my-realworld/internal/conf"
	"github.com/hl540/my-realworld/internal/src/errors"
	"github.com/hl540/my-realworld/internal/src/util"
)

type User struct {
	Username string
	PassWord string
	Email    string
	Image    string
	Bio      string
	Id       int
}

type UserRepo interface {
	// AddUser 新增用户
	AddUser(ctx context.Context, user *User) (*User, error)
	// SaveUser 保存用户信息
	SaveUser(ctx context.Context, user *User) (*User, error)
	// GetUserByUsername 获取用户信息
	GetUserByUsername(ctx context.Context, username string) (*User, error)
	// GetUserByEmail 获取用户信息
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	// GetUserByID 获取用户信息
	GetUserByID(ctx context.Context, id int) (*User, error)
}

type UserUseCase struct {
	userRepo UserRepo
	log      *log.Helper
	conf     *conf.Server
}

func NewUserUseCase(conf *conf.Server, userRepo UserRepo, logger log.Logger) *UserUseCase {
	return &UserUseCase{
		conf:     conf,
		userRepo: userRepo,
		log:      log.NewHelper(logger),
	}
}

// CreateUser 创建用户
func (uu *UserUseCase) CreateUser(ctx context.Context, user *User) (*User, error) {
	// 对密码进行加密
	user.PassWord = util.MakePassword(user.PassWord, uu.conf.Password.GetSecretKey())
	// 调用data层持久化
	user, err := uu.userRepo.AddUser(ctx, user)
	if err != nil {
		return nil, errors.NewHTTPError(500, "body", err.Error())
	}
	return user, nil
}

// UpdateUser 更新用户信息
func (uu *UserUseCase) UpdateUser(ctx context.Context, user *User) (*User, error) {
	if user.PassWord != "" {
		// 对密码进行加密
		user.PassWord = util.MakePassword(user.PassWord, uu.conf.Password.GetSecretKey())
	}
	user, err := uu.userRepo.SaveUser(ctx, user)
	if err != nil {
		return nil, errors.NewHTTPError(500, "body", err.Error())
	}
	return user, nil
}

// GetUserByUsername 获取用户信息
func (uu *UserUseCase) GetUserByUsername(ctx context.Context, username string) (*User, error) {
	user, err := uu.userRepo.GetUserByUsername(ctx, username)
	if err != nil {
		return nil, errors.NewHTTPError(500, "body", err.Error())
	}
	return user, nil
}

// GetUserByEmail 获取用户信息
func (uu *UserUseCase) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	user, err := uu.userRepo.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, errors.NewHTTPError(500, "body", err.Error())
	}
	return user, nil
}

func (uu *UserUseCase) CurrentUser(ctx context.Context, id int) (*User, error) {
	user, err := uu.userRepo.GetUserByID(ctx, id)
	if err != nil {
		return nil, errors.NewHTTPError(500, "body", err.Error())
	}
	return user, nil
}
