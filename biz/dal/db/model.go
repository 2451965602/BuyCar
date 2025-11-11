package db

import (
	"buycar/biz/model/module"
	"time"
)

type User struct {
	UserId       int64
	Username     string
	PasswordHash string
	IsAdmin      bool
	Score        int64
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func (u User) ToModuleStruct() *module.User {
	return &module.User{
		UserID:    u.UserId,
		Username:  u.Username,
		IsAdmin:   u.IsAdmin,
		Score:     u.Score,
		CreatedAt: u.CreatedAt.String(),
		UpdatedAt: u.UpdatedAt.String(),
	}
}

type Feedback struct {
	Id        int64
	UserId    int64
	ConsultId int64
	Content   string
	CreatedAt time.Time
}

func (f Feedback) ToModuleStruct() *module.Feedback {
	return &module.Feedback{
		ID:        f.Id,
		UserID:    f.UserId,
		ConsultID: f.ConsultId,
		Content:   f.Content,
		CreatedAt: f.CreatedAt.String(),
	}
}
