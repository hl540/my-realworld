package service

import (
	"context"
	"github.com/hl540/my-realworld/internal/conf"
	"github.com/hl540/my-realworld/internal/src/util"

	"github.com/hl540/my-realworld/internal/src/errors"

	pb "github.com/hl540/my-realworld/api/my_realworld/v1"
	"github.com/hl540/my-realworld/internal/biz"
)

type MyRealworldService struct {
	pb.UnimplementedMyRealworldServer

	uc   *biz.UserUseCase
	conf *conf.Server
}

func NewMyRealworldService(conf *conf.Server, uc *biz.UserUseCase) *MyRealworldService {
	return &MyRealworldService{
		uc:   uc,
		conf: conf,
	}
}

func (s *MyRealworldService) Auth(ctx context.Context, req *pb.AuthReq) (*pb.AuthRsp, error) {
	return &pb.AuthRsp{}, nil
}

func (s *MyRealworldService) Register(ctx context.Context, req *pb.RegisterReq) (*pb.RegisterRsp, error) {
	// 保存用户信息
	user, err := s.uc.CreateUser(ctx, &biz.User{
		Username: req.User.Username,
		PassWord: req.User.Password,
		Email:    req.User.Email,
		Image:    "",
	})
	if err != nil {
		return nil, err
	}
	// 生成jwt
	token, err := util.MakeJwtString(map[string]interface{}{
		util.UserName: user.Username,
	}, s.conf.Jwt.GetSecretKey())
	if err != nil {
		return nil, errors.NewHTTPError(500, "body", err.Error())
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
	userName := util.GetUserNameFromContext(ctx)
	if userName == "" {
		return nil, errors.NewHTTPError(401, "body", "there is no jwt token")
	}
	user, err := s.uc.GetUser(ctx, userName)
	if err != nil {
		return nil, err
	}
	return &pb.CurrentUserRsp{User: &pb.User{
		Email:    user.Email,
		Token:    util.ParseJwtToken(ctx),
		Username: user.Username,
		Image:    user.Image,
	}}, nil
}

func (s *MyRealworldService) UpdateUser(ctx context.Context, req *pb.UpdateUserReq) (*pb.UpdateUserRsp, error) {
	return &pb.UpdateUserRsp{}, nil
}

func (s *MyRealworldService) GetUser(ctx context.Context, req *pb.GetUserReq) (*pb.GetUserRsp, error) {
	return &pb.GetUserRsp{}, nil
}

func (s *MyRealworldService) FollowUser(ctx context.Context, req *pb.FollowUserReq) (*pb.FollowUserRsp, error) {
	return &pb.FollowUserRsp{}, nil
}

func (s *MyRealworldService) UnfollowUser(ctx context.Context, req *pb.FollowUserReq) (*pb.FollowUserRsp, error) {
	return &pb.FollowUserRsp{}, nil
}

func (s *MyRealworldService) ArticleList(ctx context.Context, req *pb.ArticleListReq) (*pb.ArticleListRsp, error) {
	return &pb.ArticleListRsp{}, nil
}

func (s *MyRealworldService) ArticleFeed(ctx context.Context, req *pb.ArticleFeedReq) (*pb.ArticleFeedRsp, error) {
	return &pb.ArticleFeedRsp{}, nil
}

func (s *MyRealworldService) GetArticle(ctx context.Context, req *pb.GetArticleReq) (*pb.GetArticleRsp, error) {
	return &pb.GetArticleRsp{}, nil
}

func (s *MyRealworldService) CreateArticle(ctx context.Context, req *pb.CreateArticleReq) (*pb.CreateArticleRsp, error) {
	return &pb.CreateArticleRsp{}, nil
}

func (s *MyRealworldService) UpdateArticle(ctx context.Context, req *pb.UpdateArticleReq) (*pb.UpdateArticleRsp, error) {
	return &pb.UpdateArticleRsp{}, nil
}

func (s *MyRealworldService) DeleteArticle(ctx context.Context, req *pb.DeleteArticleReq) (*pb.DeleteArticleRsp, error) {
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
	return nil, errors.NewHTTPError(501, "body", "xxxxxxxxxx")
	return &pb.GetTagsRsp{}, nil
}
