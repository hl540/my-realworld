package data

import (
	"context"
	"github.com/hl540/my-realworld/internal/biz"
	"xorm.io/xorm"
)

type userRepo struct {
	*Data
}

func NewUserRepo(data *Data) biz.UserRepo {
	return &userRepo{Data: data}
}

func (u *userRepo) Add(ctx context.Context, user *biz.User) (*biz.User, error) {
	data := &User{
		Username: user.Username,
		Password: user.Password,
		Email:    user.Email,
		Image:    user.Image,
		Bio:      user.Bio,
	}
	_, err := u.db.Context(ctx).Insert(data)
	if err != nil {
		return nil, err
	}
	user.Id = data.Id
	return user, nil
}

func (u *userRepo) Save(ctx context.Context, user *biz.User) (*biz.User, error) {
	data := &User{
		Image: user.Image,
		Bio:   user.Bio,
	}
	if user.Username != "" {
		data.Username = user.Username
	}
	if user.Email != "" {
		data.Email = user.Email
	}
	if user.Password != "" {
		data.Password = user.Password
	}
	_, err := u.db.Context(ctx).ID(user.Id).Update(data)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (u *userRepo) GetByUsername(ctx context.Context, username string) (*biz.User, error) {
	data := &User{}
	has, err := u.db.Context(ctx).Where("username = ?", username).Get(data)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, xorm.ErrNotExist
	}
	return &biz.User{
		Username: data.Username,
		Password: data.Password,
		Email:    data.Email,
		Image:    data.Image,
		Bio:      data.Bio,
		Id:       data.Id,
	}, nil
}

func (u *userRepo) GetByEmail(ctx context.Context, email string) (*biz.User, error) {
	data := &User{}
	has, err := u.db.Context(ctx).Where("email = ?", email).Get(data)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, xorm.ErrNotExist
	}
	return &biz.User{
		Username: data.Username,
		Password: data.Password,
		Email:    data.Email,
		Image:    data.Image,
		Bio:      data.Bio,
		Id:       data.Id,
	}, nil
}

func (u *userRepo) GetByID(ctx context.Context, id interface{}) (*biz.User, error) {
	data := &User{}
	has, err := u.db.Context(ctx).ID(id).Get(data)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, xorm.ErrNotExist
	}
	return &biz.User{
		Username: data.Username,
		Password: data.Password,
		Email:    data.Email,
		Image:    data.Image,
		Bio:      data.Bio,
		Id:       data.Id,
	}, nil
}

func (u *userRepo) AddFollow(ctx context.Context, user *biz.User, targetUser *biz.User) error {
	// 查询关注是否存在
	count, err := u.db.Context(ctx).Where("user_id = ? AND target_id = ?", user.Id, targetUser.Id).Count(Follow{})
	if err != nil {
		return err
	}
	// 已经关注
	if count == 1 {
		return nil
	}
	// 添加关注
	_, err = u.db.Context(ctx).Insert(&Follow{
		UserId:   user.Id,
		TargetId: targetUser.Id,
	})
	return err
}

func (u *userRepo) DelFollow(ctx context.Context, user *biz.User, targetUser *biz.User) error {
	// 查询关注是否存在
	count, err := u.db.Context(ctx).Where("user_id = ? AND target_id = ?", user.Id, targetUser.Id).Count(Follow{})
	if err != nil {
		return err
	}
	// 已经取消关注
	if count == 0 {
		return nil
	}
	// 删除关注
	_, err = u.db.Context(ctx).Where("user_id = ? AND target_id = ?", user.Id, targetUser.Id).Delete(&Follow{})
	return err
}

func (u *userRepo) GetFollowStatus(ctx context.Context, user *biz.User, currentUserID int64) (bool, error) {
	// 查询关注信息
	count, err := u.db.Context(ctx).Where("user_id = ? AND target_id = ?", currentUserID, user.Id).Count(Follow{})
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
