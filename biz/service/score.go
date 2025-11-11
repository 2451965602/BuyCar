package service

import (
	"buycar/biz/dal/db"
	"buycar/biz/model/module"
	"context"

	"github.com/cloudwego/hertz/pkg/app"
)

type ScoreService struct {
	ctx context.Context
	c   *app.RequestContext
}

func NewScoreService(ctx context.Context, c *app.RequestContext) *ScoreService {
	return &ScoreService{ctx: ctx, c: c}
}

func (s *ScoreService) GetUserScore() (int64, error) {
	userId := GetUidFormContext(s.c)

	user, err := db.GetUserByUserId(s.ctx, userId)
	if err != nil {
		return 0, err
	}
	return user.Score, nil
}

func (s *ScoreService) PurchaseGift(giftID int64) error {
	userId := GetUidFormContext(s.c)

	err := db.PurchaseGift(s.ctx, userId, giftID)
	if err != nil {
		return err
	}
	return nil
}

func (s *ScoreService) QueryGiftList() ([]*module.Gift, error) {
	gifts, err := db.GetAllGifts(s.ctx)
	if err != nil {
		return nil, err
	}
	var giftList []*module.Gift
	for _, gift := range gifts {
		giftList = append(giftList, gift.ToModuleStruct())
	}
	return giftList, nil
}
