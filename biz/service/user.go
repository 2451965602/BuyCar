package service

import (
	"buycar/biz/dal/db"
	"buycar/biz/model/module"
	"buycar/biz/model/user"
	"buycar/pkg/utils"
	"context"

	"github.com/cloudwego/hertz/pkg/app"
)

type UserService struct {
	ctx context.Context
	c   *app.RequestContext
}

func NewUserService(ctx context.Context, c *app.RequestContext) *UserService {
	return &UserService{ctx: ctx, c: c}
}

func (s *UserService) Register(req *user.RegisterReq) error {

	passwordHash, err := utils.EncryptPassword(req.Password)
	if err != nil {
		return err
	}

	req.Password = passwordHash

	err = db.CreateUser(s.ctx, req.Username, req.Password, req.Email, req.Name, req.Telephone)
	if err != nil {
		return err
	}
	return nil
}

func (s *UserService) LoginIn(req *user.LoginReq) (*module.User, error) {
	userInfo, err := db.GetUserByUserName(s.ctx, req.Username)
	if err != nil {
		return nil, err
	}

	err = utils.ComparePassword(userInfo.PasswordHash, req.Password)
	if err != nil {
		return nil, err
	}

	return userInfo.ToModuleStruct(), nil
}

func (s *UserService) CreateFeedback(req *user.FeedbackReq) error {
	userID := GetUidFormContext(s.c)
	err := db.CreateFeedback(s.ctx, userID, req.ConsultID, req.Content)
	if err != nil {
		return err
	}
	return nil
}
