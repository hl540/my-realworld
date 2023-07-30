package data

import (
	"context"
	"github.com/hl540/my-realworld/internal/biz"
	"github.com/hl540/my-realworld/internal/src/util"
	"gorm.io/gorm"
)

type articleRepo struct {
	*Data
}

func NewArticleRepo(data *Data) biz.ArticleRepo {
	return &articleRepo{Data: data}
}

func (a *articleRepo) Add(ctx context.Context, article *biz.Article) error {
	// 事务
	return a.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 插入文章
		poArticle := &Article{
			Slug:        article.Slug,
			Title:       article.Title,
			Description: article.Description,
			Body:        article.Body,
			AuthorID:    article.Author.ID,
		}
		if err := tx.Create(poArticle).Error; err != nil {
			return err
		}
		// 插入标签
		if len(article.TagList) > 0 {
			poTags := make([]*Tag, 0)
			for _, tag := range article.TagList {
				poTags = append(poTags, &Tag{ArticleID: poArticle.ID, Name: tag})
			}
			if err := tx.Create(poTags).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (a *articleRepo) List(ctx context.Context, tagName, favoriter, author string, limit, offset int) ([]*biz.Article, int64, error) {
	query := a.db.WithContext(ctx).Model(Article{})

	// 按tag搜索
	if tagName != "" {
		query = query.Joins("LEFT JOIN tag ON tag.article_id = article.id")
		query = query.Where("tag.name = ?", tagName)
	}

	// 按收藏人搜索
	if favoriter != "" {
		query = query.Joins("LEFT JOIN favorite ON favorite.article_id = article.id")
		query = query.Where(
			"favorite.user_id IN (?)",
			a.db.Select("id").Model(User{}).Where("username = ?", favoriter),
		)
	}

	// 按作者搜索
	if author != "" {
		query = query.Joins("LEFT JOIN user ON article.author_id = user.id")
		query = query.Where("user.username = ?", author)
	}

	// 查询数量
	var count int64
	if err := query.Count(&count).Error; err != nil {
		return nil, 0, err
	}
	if count == 0 {
		return make([]*biz.Article, 0), 0, nil
	}

	// 查询分页结果
	limit = util.IntDefault(limit, 15)
	query = query.Offset(offset).Limit(limit)
	var articles []*Article
	if err := query.Find(&articles).Error; err != nil {
		return nil, 0, nil
	}

	// 转换模型
	result := make([]*biz.Article, 0)
	for _, article := range articles {
		result = append(result, &biz.Article{
			ID:             article.ID,
			Slug:           article.Slug,
			Title:          article.Title,
			Description:    article.Description,
			Body:           article.Body,
			TagList:        make([]string, 0),
			Author:         &biz.Author{ID: article.AuthorID},
			Favorited:      false,
			FavoritesCount: 0,
			CreatedAt:      article.CreatedAt.String(),
			UpdatedAt:      article.UpdatedAt.String(),
		})
	}
	return result, count, nil
}

func (a *articleRepo) AllTag(ctx context.Context) ([]string, error) {
	tags := make([]*Tag, 0)
	err := a.db.WithContext(ctx).Model(Tag{}).Group("name").Find(&tags).Error
	if err != nil {
		return nil, err
	}
	result := make([]string, 0)
	for _, tag := range tags {
		result = append(result, tag.Name)
	}
	return result, nil
}
