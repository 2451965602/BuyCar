package db

import (
	"buycar/pkg/errno"
	"context"
	"errors"
	"gorm.io/gorm"
)

// ListAllConsults 查询所有购车咨询记录
func ListAllConsults(ctx context.Context) ([]*Consult, error) {
	var consults []Consult
	err := DB.WithContext(ctx).Find(&consults).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errno.NewErrNo(errno.ServiceConsultNotExist, "暂无咨询记录")
		}
		return nil, errno.NewErrNo(errno.InternalDatabaseErrorCode, "查询咨询记录失败: "+err.Error())
	}

	result := make([]*Consult, len(consults))
	for i := range consults {
		result[i] = &consults[i]
	}
	return result, nil
}

// CreateUserByUserID 管理员创建用户
func CreateUserByUserID(ctx context.Context, userIDStr, password string) error {
	// 检查用户名是否已存在
	var existing User
	err := DB.WithContext(ctx).Where("username = ?", userIDStr).First(&existing).Error
	if err == nil {
		return errno.NewErrNo(errno.ServiceUserExist, "用户已存在")
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return errno.NewErrNo(errno.InternalDatabaseErrorCode, "检查用户是否存在时出错: "+err.Error())
	}

	user := &User{
		Username:     userIDStr,
		PasswordHash: password,
		IsAdmin:      false, // 管理员不能通过接口创建其他管理员
		Score:        0,
	}

	err = DB.WithContext(ctx).Create(user).Error
	if err != nil {
		return errno.NewErrNo(errno.InternalDatabaseErrorCode, "创建用户失败: "+err.Error())
	}
	return nil
}

// DeleteUserByUserID 删除用户
func DeleteUserByUserID(ctx context.Context, userIDStr string) error {
	// 先检查用户是否存在
	var user User
	err := DB.WithContext(ctx).Where("username = ?", userIDStr).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errno.NewErrNo(errno.ServiceUserNotExist, "用户不存在")
		}
		return errno.NewErrNo(errno.InternalDatabaseErrorCode, "查询用户失败: "+err.Error())
	}

	// 执行删除
	err = DB.WithContext(ctx).Where("username = ?", userIDStr).Delete(&User{}).Error
	if err != nil {
		return errno.NewErrNo(errno.InternalDatabaseErrorCode, "删除用户失败: "+err.Error())
	}
	return nil
}

// ListAllFeedbacks 查询所有用户反馈
func ListAllFeedbacks(ctx context.Context) ([]*Feedback, error) {
	var feedbacks []Feedback
	err := DB.WithContext(ctx).Find(&feedbacks).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errno.NewErrNo(errno.ServiceFeedbackNotExist, "暂无反馈记录")
		}
		return nil, errno.NewErrNo(errno.InternalDatabaseErrorCode, "查询反馈记录失败: "+err.Error())
	}

	result := make([]*Feedback, len(feedbacks))
	for i := range feedbacks {
		result[i] = &feedbacks[i]
	}
	return result, nil
}
