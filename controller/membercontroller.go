package controller

import (
	"context"

	"referral-service/domain"
	"referral-service/repository"

	"go.uber.org/fx"
	"go.uber.org/zap"
)

// Contract to handle referral program membership flows
type MemberController interface {
	AddMember(ctx context.Context,
		first_name string,
		last_name *string,
		email string,
		program_id string,
		referral_code *string,
		is_active *bool,
	) (string, error)
	GetMembers(ctx context.Context, page int, size int) ([]domain.Member, error)
}

type memberCon struct {
	log *zap.Logger
	db  repository.Repository
}

type MemberParams struct {
	fx.In

	Log *zap.Logger
	Db  repository.Repository
}

func MemberNew(p MemberParams) MemberController {
	newController := &memberCon{
		log: p.Log,
		db:  p.Db,
	}

	return newController
}

func (c *memberCon) GetMembers(ctx context.Context, page int, size int) ([]domain.Member, error) {
	members, err := c.db.GetMembers(ctx, page, size)
	if err != nil {
		return nil, err
	}
	return members, nil
}

func (c *memberCon) AddMember(ctx context.Context,
	first_name string,
	last_name *string,
	email string,
	program_id string,
	referral_code *string,
	is_active *bool) (string, error) {
	memberId, err := c.db.AddMember(ctx,
		first_name,
		last_name,
		email,
		program_id,
		referral_code,
		is_active,
	)
	return memberId, err
}
