package data

import (
	"context"
	"github.com/hl540/my-realworld/internal/biz"
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
	err := u.db.WithContext(ctx).Model(User{}).Create(data).Error
	if err != nil {
		return nil, err
	}
	user.ID = data.ID
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
	err := u.db.WithContext(ctx).Model(User{}).Where("id = ?", user.ID).Updates(data).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (u *userRepo) GetByUsername(ctx context.Context, username string) (*biz.User, error) {
	data := &User{}
	err := u.db.WithContext(ctx).Model(User{}).Where("username = ?", username).First(data).Error
	if err != nil {
		return nil, err
	}
	return &biz.User{
		Username: data.Username,
		Password: data.Password,
		Email:    data.Email,
		Image:    data.Image,
		Bio:      data.Bio,
		ID:       data.ID,
	}, nil
}

func (u *userRepo) GetByEmail(ctx context.Context, email string) (*biz.User, error) {
	data := &User{}
	err := u.db.WithContext(ctx).Model(User{}).Where("email = ?", email).First(data).Error
	if err != nil {
		return nil, err
	}
	return &biz.User{
		Username: data.Username,
		Password: data.Password,
		Email:    data.Email,
		Image:    data.Image,
		Bio:      data.Bio,
		ID:       data.ID,
	}, nil
}

func (u *userRepo) GetByID(ctx context.Context, id interface{}) (*biz.User, error) {
	data := &User{}
	err := u.db.WithContext(ctx).Model(User{}).Where("id = ?", id).First(data).Error
	if err != nil {
		return nil, err
	}
	return &biz.User{
		Username: data.Username,
		Password: data.Password,
		Email:    data.Email,
		Image:    data.Image,
		Bio:      data.Bio,
		ID:       data.ID,
	}, nil
}

func (u *userRepo) AdditionalToArticle(ctx context.Context, articles []*biz.Article) error {
	ids := make([]uint, 0)
	for _, article := range articles {
		ids = append(ids, article.Author.ID)
	}
	// 查询author
	users := make([]*User, 0)
	err := u.db.WithContext(ctx).Model(User{}).Where("id IN (?)", ids).Find(&users).Error
	if err != nil {
		return err
	}
	// 关联user到article上
	userMap := make(map[uint]*User)
	for _, user := range users {
		userMap[user.ID] = user
	}
	for _, article := range articles {
		if _, ok := userMap[article.Author.ID]; !ok {
			continue
		}
		article.Author = &biz.Author{
			ID:        userMap[article.Author.ID].ID,
			Username:  userMap[article.Author.ID].Username,
			Image:     userMap[article.Author.ID].Image,
			Bio:       userMap[article.Author.ID].Bio,
			Following: false,
		}
	}
	return nil
}

func (u *userRepo) AddFollow(ctx context.Context, user *biz.User, targetUser *biz.User) error {
	// 查询关注是否存在
	var count int64
	query := u.db.WithContext(ctx).Model(Follow{}).Where("user_id = ? AND target_id = ?", user.ID, targetUser.ID)
	if err := query.Count(&count).Error; err != nil {
		return err
	}
	// 已经关注
	if count == 1 {
		return nil
	}
	// 添加关注
	return u.db.WithContext(ctx).Model(Follow{}).Create(&Follow{
		UserID:   user.ID,
		TargetID: targetUser.ID,
	}).Error
}

func (u *userRepo) DelFollow(ctx context.Context, user *biz.User, targetUser *biz.User) error {
	// 查询关注是否存在
	var count int64
	query := u.db.WithContext(ctx).Model(Follow{}).Where("user_id = ? AND target_id = ?", user.ID, targetUser.ID)
	if err := query.Count(&count).Error; err != nil {
		return err
	}
	// 已经取消关注
	if count == 0 {
		return nil
	}
	// 删除关注
	return u.db.WithContext(ctx).Where("user_id = ? AND target_id = ?", user.ID, targetUser.ID).Delete(&Follow{}).Error
}

func (u *userRepo) GetFollowStatus(ctx context.Context, user *biz.User, currentUserID uint) (bool, error) {
	// 查询关注信息
	var count int64
	err := u.db.WithContext(ctx).Model(Follow{}).Where("user_id = ? AND target_id = ?", currentUserID, user.ID).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (u *userRepo) AdditionalFollowToArticle(ctx context.Context, articles []*biz.Article, currentUserID uint) error {
	return nil
}
