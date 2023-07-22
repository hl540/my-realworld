package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/hl540/my-realworld/internal/src/errors"
)

type User struct {
	Username string
	PassWord string
	Email    string
	Image    string
}

type UserRepo interface {
	// AddUser 新增用户
	AddUser(ctx context.Context, user *User) (*User, error)
	// SaveUser 保存用户信息
	SaveUser(ctx context.Context, user *User) (*User, error)
	// GetUser 获取用户信息
	GetUser(ctx context.Context, userName string) (*User, error)
}

type UserUseCase struct {
	userRepo UserRepo
	log      *log.Helper
}

func NewUserUseCase(userRepo UserRepo, logger log.Logger) *UserUseCase {
	return &UserUseCase{userRepo: userRepo, log: log.NewHelper(logger)}
}

// CreateUser 创建用户
func (uu *UserUseCase) CreateUser(ctx context.Context, user *User) (*User, error) {
	user, err := uu.userRepo.AddUser(ctx, user)
	if err != nil {
		return nil, errors.NewHTTPError(500, "body", err.Error())
	}
	return user, nil
}

// UpdateUser 更新用户信息
func (uu *UserUseCase) UpdateUser(ctx context.Context, user *User) (*User, error) {
	user, err := uu.userRepo.SaveUser(ctx, user)
	if err != nil {
		return nil, errors.NewHTTPError(500, "body", err.Error())
	}
	return user, nil
}

// GetUser 获取用户信息
func (uu *UserUseCase) GetUser(ctx context.Context, userName string) (*User, error) {
	user, err := uu.userRepo.GetUser(ctx, userName)
	if err != nil {
		return nil, errors.NewHTTPError(500, "body", err.Error())
	}
	return user, nil
}
