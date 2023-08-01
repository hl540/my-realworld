package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/hl540/my-realworld/internal/src/errors"
	"github.com/hl540/my-realworld/internal/src/util"
	"time"
)

type Author struct {
	Id        int64
	Username  string
	Image     string
	Bio       string
	Following bool
}

type Article struct {
	Id             int64
	Slug           string
	Title          string
	Description    string
	Body           string
	TagList        []string // 标签
	Author         *Author  // 作者
	Favorited      bool     // 是否收藏
	FavoritesCount int64    // 收藏次数
	CreatedAt      string
	UpdatedAt      string
}

type ArticleRepo interface {
	// Add 新增文章
	Add(ctx context.Context, article *Article) error
	// Save 新增文章
	Save(ctx context.Context, article *Article) error
	// Delete 删除文章
	Delete(ctx context.Context, article *Article) error
	// List 获取文章列表
	List(ctx context.Context, tagName, favoriter, author string, limit, offset int) ([]*Article, int64, error)
	// AllTag 获取全部tag
	AllTag(ctx context.Context) ([]string, error)
	// GetBySlug 获取全部文章
	GetBySlug(ctx context.Context, slug string) (*Article, error)
}

type ArticleUseCase struct {
	articleRepo ArticleRepo
	userRepo    UserRepo
	log         *log.Helper
}

func NewArticleUseCase(articleRepo ArticleRepo, userRepo UserRepo, logger log.Logger) *ArticleUseCase {
	return &ArticleUseCase{
		articleRepo: articleRepo,
		userRepo:    userRepo,
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
		Id:       user.Id,
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
	return articles, uint64(count), nil
}

// GetTags 获取所有tag
func (au *ArticleUseCase) GetTags(ctx context.Context) ([]string, error) {
	tags, err := au.articleRepo.AllTag(ctx)
	if err != nil {
		return nil, errors.NewHTTPError(500, "body", err.Error())
	}
	return tags, nil
}

// GetArticle 获取文章信息
func (au *ArticleUseCase) GetArticle(ctx context.Context, slug string) (*Article, error) {
	article, err := au.articleRepo.GetBySlug(ctx, slug)
	if err != nil {
		return nil, errors.NewHTTPError(500, "body", err.Error())
	}
	return article, nil
}

// UpdateArticle 编辑文章
func (au *ArticleUseCase) UpdateArticle(ctx context.Context, article *Article) (*Article, error) {
	// 先查询文章
	data, err := au.articleRepo.GetBySlug(ctx, article.Slug)
	if err != nil {
		return nil, errors.NewHTTPError(500, "body", err.Error())
	}
	data.Title = article.Title
	data.Description = article.Description
	data.Body = article.Body
	err = au.articleRepo.Save(ctx, data)
	if err != nil {
		return nil, errors.NewHTTPError(500, "body", err.Error())
	}
	return data, nil
}

// DeleteArticle 删除文章
func (au *ArticleUseCase) DeleteArticle(ctx context.Context, slug string) error {
	// 先查询文章
	data, err := au.articleRepo.GetBySlug(ctx, slug)
	if err != nil {
		return errors.NewHTTPError(500, "body", err.Error())
	}
	err = au.articleRepo.Delete(ctx, data)
	if err != nil {
		return errors.NewHTTPError(500, "body", err.Error())
	}
	return nil
}
