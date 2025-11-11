package db

import (
	"buycar/pkg/constants"
	"buycar/pkg/errno"
	"context"
	"errors"

	"gorm.io/gorm"
)

func CreateUser(ctx context.Context, username, password string) error {
	user := &User{
		Username:     username,
		PasswordHash: password,
	}

	err := DB.WithContext(ctx).Table(constants.UserTableName).Create(user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errno.NewErrNo(errno.ServiceUserExist, " 用户已存在")
		}
		return errno.NewErrNo(errno.InternalDatabaseErrorCode, "创建用户失败: "+err.Error())
	}
	return nil

}

func GetUserByUserName(ctx context.Context, username string) (*User, error) {
	var user User
	err := DB.WithContext(ctx).Table(constants.UserTableName).Where("username = ?", username).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errno.NewErrNo(errno.ServiceUserNotExist, "用户不存在")
		}
		return nil, errno.NewErrNo(errno.InternalDatabaseErrorCode, "查询用户失败: "+err.Error())
	}
	return &user, nil
}

func CreateFeedback(ctx context.Context, userId, ConsultId int64, content string) error {
	feedback := &Feedback{
		UserId:    userId,
		ConsultId: ConsultId,
		Content:   content,
	}

	err := DB.WithContext(ctx).Table(constants.UserTableName).Create(feedback).Error
	if err != nil {
		return errno.NewErrNo(errno.InternalDatabaseErrorCode, "创建反馈失败: "+err.Error())
	}
	return nil
}

func GetUserByUserId(ctx context.Context, userId int64) (*User, error) {
	var user User
	err := DB.WithContext(ctx).Where("user_id = ?", userId).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errno.NewErrNo(errno.ServiceUserNotExist, "用户不存在")
		}
		return nil, errno.NewErrNo(errno.InternalDatabaseErrorCode, "查询用户失败: "+err.Error())
	}
	return &user, nil
}
