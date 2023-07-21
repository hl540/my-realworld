package service

import (
	"context"
	"fmt"

	"github.com/hl540/my-realworld/internal/src/errors"

	pb "github.com/hl540/my-realworld/api/my_realworld/v1"
	"github.com/hl540/my-realworld/internal/biz"
	"github.com/hl540/my-realworld/internal/src/middleware/auth"
)

type MyRealworldService struct {
	pb.UnimplementedMyRealworldServer

	uc *biz.GreeterUsecase
}

func NewMyRealworldService(uc *biz.GreeterUsecase) *MyRealworldService {
	return &MyRealworldService{uc: uc}
}

func (s *MyRealworldService) Auth(ctx context.Context, req *pb.AuthReq) (*pb.AuthRsp, error) {
	return &pb.AuthRsp{}, nil
}
func (s *MyRealworldService) Register(ctx context.Context, req *pb.RegisterReq) (*pb.RegisterRsp, error) {
	return &pb.RegisterRsp{}, nil
}
func (s *MyRealworldService) CurrentUser(ctx context.Context, req *pb.CurrentUserReq) (*pb.CurrentUserRsp, error) {
	return &pb.CurrentUserRsp{}, nil
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
	data, _ := auth.FromContext(ctx)
	fmt.Printf("%+v", data)
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
