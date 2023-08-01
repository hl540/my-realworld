package service

import (
	"context"
	"github.com/hl540/my-realworld/internal/src/util"

	"github.com/hl540/my-realworld/internal/src/errors"

	pb "github.com/hl540/my-realworld/api/my_realworld/v1"
	"github.com/hl540/my-realworld/internal/biz"
)

type MyRealworldService struct {
	pb.UnimplementedMyRealworldServer

	uu *biz.UserUseCase
	au *biz.ArticleUseCase
}

func NewMyRealworldService(uu *biz.UserUseCase, au *biz.ArticleUseCase) *MyRealworldService {
	return &MyRealworldService{
		uu: uu,
		au: au,
	}
}

func (s *MyRealworldService) Auth(ctx context.Context, req *pb.AuthReq) (*pb.AuthRsp, error) {
	// 校验
	if req.User.Email == "" || req.User.Password == "" {
		return nil, errors.NewHTTPError(500, "body", "The username and password are required")
	}
	user, token, err := s.uu.Login(ctx, req.User.Email, req.User.Password)
	if err != nil {
		return nil, err
	}
	return &pb.AuthRsp{
		User: &pb.User{
			Email:    user.Email,
			Token:    token,
			Username: user.Username,
			Image:    user.Image,
		},
	}, nil
}

func (s *MyRealworldService) Register(ctx context.Context, req *pb.RegisterReq) (*pb.RegisterRsp, error) {
	// 保存用户信息
	user, token, err := s.uu.Register(ctx, &biz.User{
		Username: req.User.Username,
		Password: req.User.Password,
		Email:    req.User.Email,
		Image:    "",
	})
	if err != nil {
		return nil, err
	}
	return &pb.RegisterRsp{
		User: &pb.User{
			Email:    user.Email,
			Token:    token,
			Username: user.Username,
			Image:    user.Image,
		},
	}, nil
}

func (s *MyRealworldService) CurrentUser(ctx context.Context, req *pb.CurrentUserReq) (*pb.CurrentUserRsp, error) {
	user, err := s.uu.CurrentUser(ctx)
	if err != nil {
		return nil, err
	}
	return &pb.CurrentUserRsp{User: &pb.User{
		Email:    user.Email,
		Token:    util.ParseTokenStr(ctx),
		Username: user.Username,
		Image:    user.Image,
	}}, nil
}

func (s *MyRealworldService) UpdateUser(ctx context.Context, req *pb.UpdateUserReq) (*pb.UpdateUserRsp, error) {
	user, err := s.uu.UpdateUser(ctx, &biz.User{
		Username: req.User.Username,
		Password: req.User.Password,
		Email:    req.User.Email,
		Image:    req.User.Image,
		Bio:      req.User.Bio,
	})
	if err != nil {
		return nil, err
	}
	return &pb.UpdateUserRsp{User: &pb.User{
		Email:    user.Email,
		Token:    util.ParseTokenStr(ctx),
		Username: user.Username,
		Image:    user.Image,
		Bio:      user.Bio,
	}}, nil
}

func (s *MyRealworldService) GetUser(ctx context.Context, req *pb.GetUserReq) (*pb.GetUserRsp, error) {
	// 校验
	if req.Username == "" {
		return nil, errors.NewHTTPError(500, "body", "The username are required")
	}
	// 查询指定用户信息
	user, following, err := s.uu.GetUserByUsername(ctx, req.Username)
	if err != nil {
		return nil, err
	}
	return &pb.GetUserRsp{Profile: &pb.Author{
		Username:  user.Username,
		Bio:       user.Bio,
		Image:     user.Image,
		Following: following,
	}}, nil
}

func (s *MyRealworldService) FollowUser(ctx context.Context, req *pb.FollowUserReq) (*pb.FollowUserRsp, error) {
	if req.Username == "" {
		return nil, errors.NewHTTPError(500, "body", "The username are required")
	}
	user, err := s.uu.FollowUser(ctx, req.Username)
	if err != nil {
		return nil, err
	}
	return &pb.FollowUserRsp{Profile: &pb.Author{
		Username:  user.Username,
		Bio:       user.Bio,
		Image:     user.Image,
		Following: true, // 固定值
	}}, nil
}

func (s *MyRealworldService) UnfollowUser(ctx context.Context, req *pb.FollowUserReq) (*pb.FollowUserRsp, error) {
	if req.Username == "" {
		return nil, errors.NewHTTPError(500, "body", "The username are required")
	}
	user, err := s.uu.UnfollowUser(ctx, req.Username)
	if err != nil {
		return nil, err
	}
	return &pb.FollowUserRsp{Profile: &pb.Author{
		Username:  user.Username,
		Bio:       user.Bio,
		Image:     user.Image,
		Following: false, // 固定值
	}}, nil
}

func (s *MyRealworldService) ArticleList(ctx context.Context, req *pb.ArticleListReq) (*pb.ArticleListRsp, error) {
	articles, count, err := s.au.ArticleList(ctx, req.Tag, req.Author, req.Favorited, req.Limit, req.Offset)
	if err != nil {
		return nil, err
	}
	result := &pb.ArticleListRsp{
		Articles:      make([]*pb.Article, 0),
		ArticlesCount: count,
	}
	for _, article := range articles {
		result.Articles = append(result.Articles, &pb.Article{
			Slug:           article.Slug,
			Title:          article.Title,
			Description:    article.Description,
			Body:           article.Body,
			TagList:        article.TagList,
			CreatedAt:      article.CreatedAt,
			UpdatedAt:      article.UpdatedAt,
			Favorited:      article.Favorited,
			FavoritesCount: uint64(article.FavoritesCount),
			Author: &pb.Author{
				Username:  article.Author.Username,
				Bio:       article.Author.Bio,
				Image:     article.Author.Image,
				Following: article.Author.Following,
			},
		})
	}
	return result, nil
}

func (s *MyRealworldService) ArticleFeed(ctx context.Context, req *pb.ArticleFeedReq) (*pb.ArticleFeedRsp, error) {
	return &pb.ArticleFeedRsp{}, nil
}

func (s *MyRealworldService) GetArticle(ctx context.Context, req *pb.GetArticleReq) (*pb.GetArticleRsp, error) {
	article, err := s.au.GetArticle(ctx, req.GetSlug())
	if err != nil {
		return nil, err
	}
	return &pb.GetArticleRsp{Article: &pb.Article{
		Slug:           article.Slug,
		Title:          article.Title,
		Description:    article.Description,
		Body:           article.Body,
		TagList:        article.TagList,
		CreatedAt:      article.CreatedAt,
		UpdatedAt:      article.UpdatedAt,
		Favorited:      article.Favorited,
		FavoritesCount: uint64(article.FavoritesCount),
		Author: &pb.Author{
			Username:  article.Author.Username,
			Bio:       article.Author.Bio,
			Image:     article.Author.Image,
			Following: article.Author.Following,
		},
	}}, nil
}

func (s *MyRealworldService) CreateArticle(ctx context.Context, req *pb.CreateArticleReq) (*pb.CreateArticleRsp, error) {
	article, err := s.au.CreateArticle(ctx, &biz.Article{
		Title:       req.Article.Title,
		Description: req.Article.Description,
		Body:        req.Article.Body,
		TagList:     req.Article.TagList,
	})
	if err != nil {
		return nil, err
	}
	return &pb.CreateArticleRsp{Article: &pb.Article{
		Slug:           "",
		Title:          article.Title,
		Description:    article.Description,
		Body:           article.Body,
		TagList:        article.TagList,
		CreatedAt:      article.CreatedAt,
		UpdatedAt:      article.UpdatedAt,
		Favorited:      article.Favorited,
		FavoritesCount: uint64(article.FavoritesCount),
		Author: &pb.Author{
			Username:  article.Author.Username,
			Bio:       article.Author.Bio,
			Image:     article.Author.Bio,
			Following: article.Author.Following,
		},
	}}, nil
}

func (s *MyRealworldService) UpdateArticle(ctx context.Context, req *pb.UpdateArticleReq) (*pb.UpdateArticleRsp, error) {
	if req.Slug == "" {
		return nil, errors.NewHTTPError(500, "body", "The slug are required")
	}
	if req.Article.Title == "" && req.Article.Description == "" && req.Article.Body == "" {
		return nil, errors.NewHTTPError(500, "body", "Updated content is a must")
	}
	article, err := s.au.UpdateArticle(ctx, &biz.Article{
		Slug:        req.Slug,
		Title:       req.Article.Title,
		Description: req.Article.Description,
		Body:        req.Article.Body,
	})
	if err != nil {
		return nil, err
	}
	return &pb.UpdateArticleRsp{Article: &pb.Article{
		Slug:           article.Slug,
		Title:          article.Title,
		Description:    article.Description,
		Body:           article.Body,
		TagList:        article.TagList,
		CreatedAt:      article.CreatedAt,
		UpdatedAt:      article.UpdatedAt,
		Favorited:      article.Favorited,
		FavoritesCount: uint64(article.FavoritesCount),
		Author: &pb.Author{
			Username:  article.Author.Username,
			Bio:       article.Author.Bio,
			Image:     article.Author.Image,
			Following: article.Author.Following,
		},
	}}, nil
}

func (s *MyRealworldService) DeleteArticle(ctx context.Context, req *pb.DeleteArticleReq) (*pb.DeleteArticleRsp, error) {
	if req.Slug == "" {
		return nil, errors.NewHTTPError(500, "body", "The slug are required")
	}
	err := s.au.DeleteArticle(ctx, req.Slug)
	if err != nil {
		return nil, err
	}
	return &pb.DeleteArticleRsp{}, nil
}

func (s *MyRealworldService) CommentsArticle(ctx context.Context, req *pb.CommentsArticleReq) (*pb.CommentsArticleRsp, error) {
	return &pb.CommentsArticleRsp{}, nil
}

func (s *MyRealworldService) GetComments(ctx context.Context, req *pb.GetCommentsReq) (*pb.GetCommentsRsp, error) {
	return &pb.GetCommentsRsp{}, nil
}

func (s *MyRealworldService) DeleteComments(ctx context.Context, req *pb.DeleteCommentsReq) (*pb.DeleteCommentsRsp, error) {
	return &pb.DeleteCommentsRsp{}, nil
}

func (s *MyRealworldService) FavoriteArticle(ctx context.Context, req *pb.FavoriteArticleReq) (*pb.CommentsArticleRsp, error) {
	return &pb.CommentsArticleRsp{}, nil
}

func (s *MyRealworldService) UnfavoriteArticle(ctx context.Context, req *pb.FavoriteArticleReq) (*pb.FavoriteArticleReq, error) {
	return &pb.FavoriteArticleReq{}, nil
}

func (s *MyRealworldService) GetTags(ctx context.Context, req *pb.GetTagsReq) (*pb.GetTagsRsp, error) {
	tags, err := s.au.GetTags(ctx)
	if err != nil {
		return nil, err
	}
	return &pb.GetTagsRsp{Tags: tags}, nil
}
