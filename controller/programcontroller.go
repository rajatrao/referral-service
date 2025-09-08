package controller

import (
	"context"

	"referral-service/domain"
	"referral-service/repository"

	"go.uber.org/fx"
	"go.uber.org/zap"
)

// Contract for referral programs
type ProgramController interface {
	AddProgram(ctx context.Context, name string, title string, active bool) (string, error)
	UpdateProgram(ctx context.Context, id string, name *string, title *string, active *bool) (*domain.Program, error)
	GetProgram(ctx context.Context, id string) (*domain.Program, error)
	GetPrograms(ctx context.Context, page int, size int) ([]domain.Program, error)
}

type programCon struct {
	log *zap.Logger
	db  repository.Repository
}

type ProgramParams struct {
	fx.In

	Log *zap.Logger
	Db  repository.Repository
}

func ProgramNew(p ProgramParams) ProgramController {
	newController := &programCon{
		log: p.Log,
		db:  p.Db,
	}

	return newController
}

func (c *programCon) GetProgram(ctx context.Context, id string) (*domain.Program, error) {
	program, err := c.db.GetProgram(ctx, id)
	if err != nil {
		return nil, err
	}
	return &program, nil
}

func (c *programCon) GetPrograms(ctx context.Context, page int, size int) ([]domain.Program, error) {
	programs, err := c.db.GetPrograms(ctx, page, size)
	if err != nil {
		return nil, err
	}
	return programs, nil
}

func (c *programCon) AddProgram(ctx context.Context, name string, title string, active bool) (string, error) {
	programId, err := c.db.AddProgram(ctx, name, title, active)
	return programId, err
}

func (c *programCon) UpdateProgram(ctx context.Context, id string, name *string, title *string, active *bool) (*domain.Program, error) {
	err := c.db.UpdateProgram(ctx, id, name, title, active)
	if err != nil {
		return nil, err
	}
	program, getErr := c.db.GetProgram(ctx, id)
	return &program, getErr
}
