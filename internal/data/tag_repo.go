package data

import (
	"context"
	"github.com/hl540/my-realworld/internal/biz"
)

type tagRepo struct {
	*Data
}

func NewTagRepo(data *Data) biz.TagRepo {
	return &tagRepo{Data: data}
}

func (t *tagRepo) AdditionalToArticle(ctx context.Context, articles []*biz.Article) error {
	ids := make([]uint, 0)
	for _, article := range articles {
		ids = append(ids, article.ID)
	}
	// 查询tag
	tags := make([]*Tag, 0)
	err := t.db.WithContext(ctx).Model(Tag{}).Where("article_id IN (?)", ids).Find(&tags).Error
	if err != nil {
		return err
	}
	// 关联tag到article上
	tagMap := make(map[uint][]string)
	for _, tag := range tags {
		tagMap[tag.ArticleID] = append(tagMap[tag.ArticleID], tag.Name)
	}
	for _, article := range articles {
		if _, ok := tagMap[article.ID]; ok {
			article.TagList = tagMap[article.ID]
		}
	}
	return nil
}
