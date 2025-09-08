package repository

import (
	"context"
	"referral-service/domain"
)

type Repository interface {
	// Program
	AddProgram(ctx context.Context,
		name string,
		title string,
		active bool) (string, error)
	UpdateProgram(ctx context.Context,
		id string,
		name *string,
		title *string,
		active *bool) error
	GetPrograms(ctx context.Context, page int, size int) ([]domain.Program, error)
	GetProgram(ctx context.Context, programId string) (domain.Program, error)
	// Member
	AddMember(ctx context.Context,
		first_name string,
		last_name *string,
		email string,
		program_id string,
		referral_code *string,
		is_active *bool) (string, error)
	GetMembers(ctx context.Context, page int, size int) ([]domain.Member, error)
	// Referral
	AddReferral(ctx context.Context,
		first_name *string,
		last_name *string,
		email *string,
		phone *string,
		referral_code string) (string, error)
	GetReferrals(ctx context.Context, page int, size int) ([]domain.Referral, error)
}
