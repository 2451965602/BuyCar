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

type ScoreTransaction struct {
	Id           int64
	UserId       int64
	ChangeAmount int64
	Reason       string
	RefId        int64
	Description  string
	CreatedAt    time.Time
}

func (s ScoreTransaction) ToModuleStruct() *module.ScoreTransaction {
	return &module.ScoreTransaction{
		ID:           s.Id,
		UserID:       s.UserId,
		ChangeAmount: s.ChangeAmount,
		Reason:       s.Reason,
		RefID:        s.RefId,
		Description:  s.Description,
		CreatedAt:    s.CreatedAt.String(),
	}
}

type Gift struct {
	GiftId      int64
	Name        string
	Description string
	ScoreCost   int64
	Stock       int64
	Status      string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (g Gift) ToModuleStruct() *module.Gift {
	return &module.Gift{
		GiftID:      g.GiftId,
		Name:        g.Name,
		Description: g.Description,
		ScoreCost:   g.ScoreCost,
		Stock:       g.Stock,
		Status:      g.Status,
		CreatedAt:   g.CreatedAt.String(),
		UpdatedAt:   g.UpdatedAt.String(),
	}
}

type GiftPurchase struct {
	Id        int64
	UserId    int64
	GiftId    int64
	Number    int64
	ScoreCost int64
	Status    string
	CreatedAt time.Time
}
