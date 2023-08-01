package data

import (
	"context"
	"github.com/hl540/my-realworld/internal/biz"
	"github.com/hl540/my-realworld/internal/src/util"
	"xorm.io/builder"
	"xorm.io/xorm"
)

type articleRepo struct {
	*Data
}

func NewArticleRepo(data *Data) biz.ArticleRepo {
	return &articleRepo{Data: data}
}

func (a *articleRepo) Add(ctx context.Context, article *biz.Article) error {
	data := &Article{
		Slug:        article.Slug,
		Title:       article.Title,
		Description: article.Description,
		Body:        article.Body,
		AuthorId:    article.Author.Id,
	}
	// 事务
	_, err := a.db.Transaction(func(session *xorm.Session) (interface{}, error) {
		if _, err := session.Context(ctx).Insert(data); err != nil {
			return nil, err
		}
		tags := make([]*Tag, 0)
		for _, tag := range article.TagList {
			tags = append(tags, &Tag{ArticleId: data.Id, Name: tag})
		}
		if _, err := session.Context(ctx).Insert(tags); err != nil {
			return nil, err
		}
		return true, nil
	})
	return err
}

func (a *articleRepo) List(ctx context.Context, tagName, favoriter, author string, limit, offset int) ([]*biz.Article, int64, error) {
	query := a.db.Table(Article{})
	query = query.Join("LEFT", "tag", "tag.article_id = article.id")
	query = query.Join("LEFT", "favorite", "favorite.article_id = article.id")
	// 按tag搜索
	if tagName != "" {
		query = query.Where("tag.name = ?", tagName)
	}
	// 按收藏人搜索
	if favoriter != "" {
		subQuery := builder.Select("id").From("user").Where(builder.Eq{"username": favoriter})
		query = query.In("favorite.user_id", subQuery)
	}
	// 按作者搜索
	if author != "" {
		subQuery := builder.Select("id").From("user").Where(builder.Eq{"username": author})
		query = query.In("author_id", subQuery)
	}

	// 查询数量
	count, err := query.Context(ctx).Count()
	if err != nil {
		return nil, 0, err
	}
	if count == 0 {
		return make([]*biz.Article, 0), 0, nil
	}

	// 查询分页结果
	query = query.Limit(util.IntDefault(limit, 15), offset)
	articles := make([]*Article, 0)
	err = query.Context(ctx).Find(&articles)
	if err != nil {
		return nil, 0, nil
	}

	// 转换模型
	result := make([]*biz.Article, 0)
	for _, article := range articles {
		tArticle := &biz.Article{
			Id:          article.Id,
			Slug:        article.Slug,
			Title:       article.Title,
			Description: article.Description,
			Body:        article.Body,
			Author:      &biz.Author{Id: article.AuthorId},
			CreatedAt:   article.CreatedAt.String(),
			UpdatedAt:   article.UpdatedAt.String(),
		}
		result = append(result, tArticle)
	}
	// 附加tag信息
	if err := a.additionalTag(ctx, result); err != nil {
		return nil, 0, err
	}
	// 附加作者信息
	if err := a.additionalAuthor(ctx, result); err != nil {
		return nil, 0, err
	}
	return result, count, nil
}

// 附加tag信息
func (a *articleRepo) additionalTag(ctx context.Context, articles []*biz.Article) error {
	ids := make([]int64, 0)
	for _, article := range articles {
		ids = append(ids, article.Id)
	}
	// 查询tag
	tags := make([]*Tag, 0)
	err := a.db.Context(ctx).In("article_id", ids).Find(&tags)
	if err != nil {
		return err
	}
	tagMap := make(map[int64][]string)
	for _, tag := range tags {
		tagMap[tag.ArticleId] = append(tagMap[tag.ArticleId], tag.Name)
	}
	for _, article := range articles {
		if _, ok := tagMap[article.Id]; ok {
			article.TagList = tagMap[article.Id]
		}
	}
	return nil
}

// 附加Author信息
func (a *articleRepo) additionalAuthor(ctx context.Context, articles []*biz.Article) error {
	ids := make([]int64, 0)
	for _, article := range articles {
		ids = append(ids, article.Author.Id)
	}
	// 查询作者
	users := make([]*User, 0)
	err := a.db.Context(ctx).In("id", ids).Find(&users)
	if err != nil {
		return err
	}
	userMap := make(map[int64]*User)
	for _, user := range users {
		userMap[user.Id] = user
	}
	// 查询作者关注信息
	follows := make([]*Follow, 0)
	currentUserID := util.GetUserInfo(ctx).UserID
	err = a.db.Context(ctx).Where("user_id = ?", currentUserID).In("target_id", ids).Find(&follows)
	if err != nil {
		return err
	}
	followMap := make(map[int64]bool)
	for _, follow := range follows {
		followMap[follow.TargetId] = true
	}

	for _, article := range articles {
		if author, ok := userMap[article.Author.Id]; ok {
			article.Author = &biz.Author{
				Id:        author.Id,
				Username:  author.Username,
				Image:     author.Image,
				Bio:       author.Bio,
				Following: followMap[author.Id],
			}
		}
	}
	return nil
}

// 附加文章收藏信息信息
func (a *articleRepo) additionalFollowing(ctx context.Context, articles []*biz.Article) error {
	ids := make([]int64, 0)
	for _, article := range articles {
		ids = append(ids, article.Id)
	}
	// 查询当前用户是否收藏
	favorites := make([]*Favorite, 0)
	currentUserID := util.GetUserInfo(ctx).UserID
	err := a.db.Context(ctx).Where("user_id = ?", currentUserID).In("target_id", ids).Find(&favorites)
	if err != nil {
		return err
	}
	favoriteMap := make(map[int64]bool)
	for _, favorite := range favorites {
		favoriteMap[favorite.ArticleId] = true
	}
	// 查询关注总量
	favoritesCount := make([]struct {
		Count     int64
		ArticleId int64
	}, 0)
	query := a.db.Context(ctx).Table(Favorite{}).Select("COUNT(*) as count, article_id")
	err = query.In("article_id", ids).GroupBy("article_id").Find(&favoritesCount)
	if err != nil {
		return err
	}
	favoritesCountMap := make(map[int64]int64)
	for _, favorites := range favoritesCount {
		favoritesCountMap[favorites.ArticleId] = favorites.Count
	}
	for _, article := range articles {
		article.FavoritesCount = favoritesCountMap[article.Id]
	}
	return nil
}

func (a *articleRepo) AllTag(ctx context.Context) ([]string, error) {
	tags := make([]*Tag, 0)
	err := a.db.Context(ctx).GroupBy("name").Find(&tags)
	if err != nil {
		return nil, err
	}
	result := make([]string, 0)
	for _, tag := range tags {
		result = append(result, tag.Name)
	}
	return result, nil
}

func (a *articleRepo) GetBySlug(ctx context.Context, slug string) (*biz.Article, error) {
	var article = &Article{}
	ex, err := a.db.Context(ctx).Table(Article{}).Where("slug = ?", slug).Get(article)
	if err != nil {
		return nil, err
	}
	if !ex {
		return nil, xorm.ErrNotExist
	}
	result := make([]*biz.Article, 0)
	result = append(result, &biz.Article{
		Id:          article.Id,
		Slug:        article.Slug,
		Title:       article.Title,
		Description: article.Description,
		Body:        article.Body,
		Author:      &biz.Author{Id: article.AuthorId},
		CreatedAt:   article.CreatedAt.String(),
		UpdatedAt:   article.UpdatedAt.String(),
	})
	// 附加tag信息
	if err := a.additionalTag(ctx, result); err != nil {
		return nil, err
	}
	// 附加作者信息
	if err := a.additionalAuthor(ctx, result); err != nil {
		return nil, err
	}
	return result[0], nil
}

func (a *articleRepo) Save(ctx context.Context, article *biz.Article) error {
	updata := make(map[string]interface{})
	if article.Title != "" {
		updata["title"] = article.Title
	}
	if article.Description != "" {
		updata["description"] = article.Description
	}
	if article.Body != "" {
		updata["body"] = article.Body
	}
	_, err := a.db.Context(ctx).Table(Article{}).ID(article.Id).Update(updata)
	return err
}

func (a *articleRepo) Delete(ctx context.Context, article *biz.Article) error {
	// 事务
	_, err := a.db.Transaction(func(session *xorm.Session) (interface{}, error) {
		// 删除文章
		_, err := session.Context(ctx).Table(Article{}).ID(article.Id).Delete(Article{})
		if err != nil {
			return nil, err
		}
		// 删除文章的标签
		_, err = session.Context(ctx).Table(Tag{}).Where("article_id = ?", article.Id).Delete(Tag{})
		if err != nil {
			return nil, err
		}
		// 删除文章的收藏
		_, err = session.Context(ctx).Table(Favorite{}).Where("article_id = ?", article.Id).Delete(Favorite{})
		if err != nil {
			return nil, err
		}
		return nil, nil
	})
	return err
}
