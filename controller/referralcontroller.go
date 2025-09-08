package controller

import (
	"context"

	"referral-service/domain"
	"referral-service/repository"

	"go.uber.org/fx"
	"go.uber.org/zap"
)

// Contract defines Member Referrals flows
type ReferralController interface {
	AddReferral(ctx context.Context,
		first_name *string,
		last_name *string,
		email *string,
		phone *string,
		referral_code string) (string, error)
	GetReferrals(ctx context.Context, page int, size int) ([]domain.Referral, error)
}

type referralCon struct {
	log *zap.Logger
	db  repository.Repository
}

type ReferralParams struct {
	fx.In

	Log *zap.Logger
	Db  repository.Repository
}

func ReferralNew(p ReferralParams) ReferralController {
	newController := &referralCon{
		log: p.Log,
		db:  p.Db,
	}

	return newController
}

func (c *referralCon) GetReferrals(ctx context.Context, page int, size int) ([]domain.Referral, error) {
	referrals, err := c.db.GetReferrals(ctx, page, size)
	if err != nil {
		return nil, err
	}
	return referrals, nil
}

func (c *referralCon) AddReferral(ctx context.Context,
	first_name *string,
	last_name *string,
	email *string,
	phone *string,
	referral_code string) (string, error) {
	referralId, err := c.db.AddReferral(ctx,
		first_name,
		last_name,
		email,
		phone,
		referral_code,
	)
	return referralId, err
}
