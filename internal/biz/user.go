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
	Password string
	Email    string
	Image    string
	Bio      string
	ID       uint
}

type UserRepo interface {
	// Add 新增用户
	Add(ctx context.Context, user *User) (*User, error)
	// Save 保存用户信息
	Save(ctx context.Context, user *User) (*User, error)
	// GetByUsername 获取用户信息
	GetByUsername(ctx context.Context, username string) (*User, error)
	// GetByEmail 获取用户信息
	GetByEmail(ctx context.Context, email string) (*User, error)
	// GetByID 获取用户信息
	GetByID(ctx context.Context, id interface{}) (*User, error)
	// AdditionalToArticle 将user附加到article
	AdditionalToArticle(ctx context.Context, articles []*Article) error
	// AddFollow 关注用户
	AddFollow(ctx context.Context, user *User, targetUser *User) error
	// DelFollow 取消关注
	DelFollow(ctx context.Context, user *User, targetUser *User) error
	// GetFollowStatus 获取用户关注状态
	GetFollowStatus(ctx context.Context, users *User, currentUserID uint) (bool, error)
	// AdditionalFollowToArticle 附加关注信息到文章作者
	AdditionalFollowToArticle(ctx context.Context, articles []*Article, currentUserID uint) error
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

// Register 创建用户
func (uu *UserUseCase) Register(ctx context.Context, user *User) (*User, string, error) {
	// 对密码进行加密
	user.Password = util.MD5(user.Password + uu.conf.Password.GetSecretKey())
	// 调用data层持久化
	user, err := uu.userRepo.Add(ctx, user)
	if err != nil {
		return nil, "", errors.NewHTTPError(500, "body", err.Error())
	}
	// 生成token
	token, err := util.NewJwtByData(uu.conf.Jwt.GetSecretKey(), map[string]interface{}{
		util.UserID:    user.ID,
		util.UserName:  user.Username,
		util.UserEmail: user.Email,
	}).Token()
	if err != nil {
		return nil, "", errors.NewHTTPError(500, "body", err.Error())
	}
	return user, token, nil
}

// UpdateUser 更新用户信息
func (uu *UserUseCase) UpdateUser(ctx context.Context, newUser *User) (*User, error) {
	// 获取当前用户
	userInfo := util.GetUserInfo(ctx)
	if userInfo == nil {
		return nil, errors.NewHTTPError(401, "body", "there is no jwt token")
	}
	newUser.ID = userInfo.UserID

	// 对密码进行加密
	if newUser.Password != "" {
		newUser.Password = util.MD5(newUser.Password + uu.conf.Password.GetSecretKey())
	}

	// 保存用户信息
	user, err := uu.userRepo.Save(ctx, newUser)
	if err != nil {
		return nil, errors.NewHTTPError(500, "body", err.Error())
	}
	return user, nil
}

// GetUserByUsername 获取用户信息
func (uu *UserUseCase) GetUserByUsername(ctx context.Context, username string) (*User, bool, error) {
	user, err := uu.userRepo.GetByUsername(ctx, username)
	if err != nil {
		return nil, false, errors.NewHTTPError(500, "body", err.Error())
	}
	// 获取关注状态
	currentUser := util.GetUserInfo(ctx)
	following, err := uu.userRepo.GetFollowStatus(ctx, user, currentUser.UserID)
	return user, following, nil
}

// Login 用户登陆
func (uu *UserUseCase) Login(ctx context.Context, email, password string) (*User, string, error) {
	// 查询当前要登陆的用户信息
	user, err := uu.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return nil, "", errors.NewHTTPError(500, "body", err.Error())
	}

	// 检查密码
	pass := util.MD5(password + uu.conf.Password.GetSecretKey())
	if pass != user.Password {
		return nil, "", errors.NewHTTPError(500, "body", "Password error")
	}

	// 生成token
	token, err := util.NewJwtByData(uu.conf.Jwt.GetSecretKey(), map[string]interface{}{
		util.UserID:    user.ID,
		util.UserName:  user.Username,
		util.UserEmail: user.Email,
	}).Token()
	if err != nil {
		return nil, "", errors.NewHTTPError(500, "body", err.Error())
	}
	return user, token, nil
}

// CurrentUser 获取当前登陆用户
func (uu *UserUseCase) CurrentUser(ctx context.Context) (*User, error) {
	// 获取当前用户
	userInfo := util.GetUserInfo(ctx)
	if userInfo == nil {
		return nil, errors.NewHTTPError(401, "body", "there is no jwt token")
	}
	user, err := uu.userRepo.GetByID(ctx, userInfo.UserID)
	if err != nil {
		return nil, errors.NewHTTPError(500, "body", err.Error())
	}
	return user, nil
}

// FollowUser 关注用户
func (uu *UserUseCase) FollowUser(ctx context.Context, username string) (*User, error) {
	// 获取当前用户
	userInfo := util.GetUserInfo(ctx)
	if userInfo == nil {
		return nil, errors.NewHTTPError(401, "body", "there is no jwt token")
	}
	// 获取目标用户
	targetUser, err := uu.userRepo.GetByUsername(ctx, username)
	if err != nil {
		return nil, errors.NewHTTPError(500, "body", err.Error())
	}
	// 添加关注
	err = uu.userRepo.AddFollow(ctx, &User{ID: userInfo.UserID}, targetUser)
	if err != nil {
		return nil, errors.NewHTTPError(500, "body", err.Error())
	}
	return targetUser, nil
}

// UnfollowUser 关注用户
func (uu *UserUseCase) UnfollowUser(ctx context.Context, username string) (*User, error) {
	// 获取当前用户
	userInfo := util.GetUserInfo(ctx)
	if userInfo == nil {
		return nil, errors.NewHTTPError(401, "body", "there is no jwt token")
	}
	// 获取目标用户
	targetUser, err := uu.userRepo.GetByUsername(ctx, username)
	if err != nil {
		return nil, errors.NewHTTPError(500, "body", err.Error())
	}
	// 添加关注
	err = uu.userRepo.DelFollow(ctx, &User{ID: userInfo.UserID}, targetUser)
	if err != nil {
		return nil, errors.NewHTTPError(500, "body", err.Error())
	}
	return targetUser, nil
}
