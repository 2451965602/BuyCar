package db

import (
	"buycar/pkg/constants"
	"buycar/pkg/errno"
	"context"
	"strconv"
)

func PurchaseGift(ctx context.Context, userID int64, giftID int64) error {

	var gift Gift
	err := DB.WithContext(ctx).Table(constants.GiftTable).Where("gift_id = ?", giftID).First(&gift).Error
	if err != nil {
		return errno.NewErrNo(errno.GiftNotFound, "未知的礼物ID")
	}

	var user User
	err = DB.WithContext(ctx).Table(constants.UserTableName).Where("user_id = ?", userID).First(&user).Error
	if err != nil {
		return errno.InternalServiceError
	}

	if user.Score < gift.ScoreCost {
		return errno.NewErrNo(errno.ScoreNotEnough, "积分不足，无法兑换该礼物")
	}

	user.Score -= gift.ScoreCost
	err = DB.WithContext(ctx).Table(constants.UserTableName).Save(&user).Error
	if err != nil {
		return err
	}

	scoreTransaction := ScoreTransaction{
		UserId:       userID,
		ChangeAmount: -gift.ScoreCost,
		Reason:       "Purchase Gift",
		RefId:        giftID,
		Description:  "Purchased gift with ID " + strconv.FormatInt(giftID, 10),
	}
	err = DB.WithContext(ctx).Table(constants.ScoreTable).Create(&scoreTransaction).Error
	if err != nil {
		return errno.NewErrNo(errno.InternalServiceErrorCode, "记录积分变更失败")
	}

	return nil

}

func GetAllGifts(ctx context.Context) ([]Gift, error) {
	var gifts []Gift
	err := DB.WithContext(ctx).Table(constants.GiftTable).Find(&gifts).Error
	if err != nil {
		return nil, errno.InternalServiceError
	}
	return gifts, nil
}
