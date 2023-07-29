package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/hl540/my-realworld/internal/src/errors"
	"github.com/hl540/my-realworld/internal/src/util"
	"time"
)

type Article struct {
	ID             uint
	Slug           string
	Title          string
	Description    string
	Body           string
	TagList        []string // 标签
	Author         *Author  // 作者
	Favorited      bool     // 是否收藏
	FavoritesCount uint64   // 收藏次数
	CreatedAt      string
	UpdatedAt      string
}

type ArticleRepo interface {
	// Add 新增文章
	Add(ctx context.Context, article *Article) error
	// List 获取文章列表
	List(ctx context.Context, tagName, favoriter, author string, limit, offset int) ([]*Article, int64, error)
}

type TagRepo interface {
	// AdditionalToArticle 将tag附加到article上
	AdditionalToArticle(ctx context.Context, articles []*Article) error
}

type ArticleUseCase struct {
	articleRepo ArticleRepo
	userRepo    UserRepo
	tagRepo     TagRepo
	log         *log.Helper
}

func NewArticleUseCase(articleRepo ArticleRepo, userRepo UserRepo, tagRepo TagRepo, logger log.Logger) *ArticleUseCase {
	return &ArticleUseCase{
		articleRepo: articleRepo,
		userRepo:    userRepo,
		tagRepo:     tagRepo,
		log:         log.NewHelper(logger),
	}
}

// CreateArticle 创建文章
func (au *ArticleUseCase) CreateArticle(ctx context.Context, article *Article) (*Article, error) {
	// 获取当前用户
	userInfo := util.GetUserInfo(ctx)
	if userInfo == nil {
		return nil, errors.NewHTTPError(401, "body", "there is no jwt token")
	}
	user, err := au.userRepo.GetByID(ctx, userInfo.UserID)
	if err != nil {
		return nil, errors.NewHTTPError(500, "body", err.Error())
	}
	// 文章作者信息
	article.Author = &Author{
		ID:       user.Id,
		Username: user.Username,
		Image:    user.Image,
		Bio:      user.Bio,
	}
	// 创建文章
	article.Slug = util.MD5(article.Title + time.Now().String())
	if err := au.articleRepo.Add(ctx, article); err != nil {
		return nil, errors.NewHTTPError(500, "body", err.Error())
	}
	return article, nil
}

// ArticleList 获取文章列表
func (au *ArticleUseCase) ArticleList(ctx context.Context, tagName, author, favorited string, limit, offset uint64) ([]*Article, uint64, error) {
	// 查询文章
	articles, count, err := au.articleRepo.List(ctx, tagName, favorited, author, int(limit), int(offset))
	if err != nil {
		return nil, 0, errors.NewHTTPError(500, "body", err.Error())
	}
	// 附加tag信息
	if au.tagRepo.AdditionalToArticle(ctx, articles); err != nil {
		return nil, 0, errors.NewHTTPError(500, "body", err.Error())
	}
	// 附加author信息
	if au.userRepo.AdditionalToArticle(ctx, articles); err != nil {
		return nil, 0, errors.NewHTTPError(500, "body", err.Error())
	}
	return articles, uint64(count), nil
}
