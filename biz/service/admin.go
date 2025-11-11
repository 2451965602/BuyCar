package service

import (
	"buycar/biz/dal/db"
	"buycar/biz/model/admin"
	"buycar/biz/model/module"
	"buycar/pkg/errno"
	"context"
	"errors"
)

type AdminService struct{}

// QueryAllConsults 查询所有咨询记录
func (s *AdminService) QueryAllConsults(ctx context.Context) *admin.QueryAllConsultsResp {
	// 调用数据库层获取所有咨询记录
	consults, err := db.ListAllConsults(ctx)
	if err != nil {
		var Err *errno.ErrNo
		errors.As(err, &Err)
		return &admin.QueryAllConsultsResp{
			BaseResponse: &module.BaseResp{
				Code:    int32(Err.ErrorCode),
				Message: Err.ErrorMsg,
			},
		}
	}

	// 转换为模型层结构体
	modelConsults := make([]*module.Consult, 0, len(consults))
	for _, consult := range consults {
		modelConsults = append(modelConsults, consult.ToModuleStruct())
	}

	return &admin.QueryAllConsultsResp{
		BaseResponse: &module.BaseResp{
			Code:    0,
			Message: "success",
		},
		Consults: modelConsults,
	}
}

// AdminAddUser 添加用户
func (s *AdminService) AdminAddUser(ctx context.Context, req *admin.AdminAddUserReq) *admin.AdminAddUserResp {
	// 调用数据库层创建用户
	err := db.CreateUserByUserID(ctx, req.UserID, req.Password)
	if err != nil {
		var Err *errno.ErrNo
		errors.As(err, &Err)
		return &admin.AdminAddUserResp{
			BaseResponse: &module.BaseResp{
				Code:    int32(Err.ErrorCode),
				Message: Err.ErrorMsg,
			},
		}
	}

	return &admin.AdminAddUserResp{
		BaseResponse: &module.BaseResp{
			Code:    0,
			Message: "success",
		},
	}
}

// AdminDeleteUser 删除用户
func (s *AdminService) AdminDeleteUser(ctx context.Context, req *admin.AdminDeleteUserReq) *admin.AdminDeleteUserResp {
	// 调用数据库层删除用户
	err := db.DeleteUserByUserID(ctx, req.UserID)
	if err != nil {
		var Err *errno.ErrNo
		errors.As(err, &Err)
		return &admin.AdminDeleteUserResp{
			BaseResponse: &module.BaseResp{
				Code:    int32(Err.ErrorCode),
				Message: Err.ErrorMsg,
			},
		}
	}

	return &admin.AdminDeleteUserResp{
		BaseResponse: &module.BaseResp{
			Code:    0,
			Message: "success",
		},
	}
}

// QueryFeedbackAnalysis 查询用户反馈分析
func (s *AdminService) QueryFeedbackAnalysis(ctx context.Context) *admin.QueryFeedbackAnalysisResp {
	// 调用数据库层获取所有反馈记录
	feedbacks, err := db.ListAllFeedbacks(ctx)
	if err != nil {
		var Err *errno.ErrNo
		errors.As(err, &Err)
		return &admin.QueryFeedbackAnalysisResp{
			BaseResponse: &module.BaseResp{
				Code:    int32(Err.ErrorCode),
				Message: Err.ErrorMsg,
			},
		}
	}

	// 转换为模型层结构体
	modelFeedbacks := make([]*module.Feedback, 0, len(feedbacks))
	for _, feedback := range feedbacks {
		modelFeedbacks = append(modelFeedbacks, feedback.ToModuleStruct())
	}

	return &admin.QueryFeedbackAnalysisResp{
		BaseResponse: &module.BaseResp{
			Code:    0,
			Message: "success",
		},
		Feedbacks: modelFeedbacks,
	}
}
