package db

import (
	"buycar/biz/model/module"
	"fmt"
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

type Consult struct {
	ConsultId       int64
	UserId          *int64
	Title           *string
	BudgetRange     *string
	PreferredType   *string
	UseCase         *string
	FuelType        *string
	BrandPreference *string
	LlmModel        *string
	LlmPrompt       *string
	LlmResponse     *string
	Recommendations *string
	Status          string
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

func (c Consult) ToModuleStruct() *module.Consult {
	var userID *string
	if c.UserId != nil {
		uidStr := fmt.Sprintf("%d", *c.UserId)
		userID = &uidStr
	}

	return &module.Consult{
		ConsultID:       fmt.Sprintf("%d", c.ConsultId),
		UserID:          userID,
		BudgetRange:     c.BudgetRange,
		PreferredType:   c.PreferredType,
		UseCase:         c.UseCase,
		FuelType:        c.FuelType,
		BrandPreference: c.BrandPreference,
		LlmModel:        c.LlmModel,
		LlmPrompt:       c.LlmPrompt,
		LlmResponse:     c.LlmResponse,
		Recommendations: c.Recommendations,
		CreatedAt:       c.CreatedAt.Unix(),
		UpdatedAt:       c.UpdatedAt.Unix(),
	}
}
