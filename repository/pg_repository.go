package repository

import (
	"context"
	"database/sql"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"referral-service/domain"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"go.uber.org/config"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type pgRepository struct {
	log    *zap.Logger
	db     *sqlx.DB
	txOpts *sql.TxOptions
}

type Params struct {
	fx.In

	Log *zap.Logger
	Cfg config.Provider
}

func New(p Params) (Repository, error) {

	cfg := p.Cfg

	connStr := fmt.Sprintf(
		"user=%s dbname=%s password=%s host=%s port=5432 sslmode=disable",
		cfg.Get("postgres.user").String(),
		cfg.Get("postgres.db_name").String(),
		cfg.Get("postgres.password").String(),
		cfg.Get("postgres.host").String(),
	)

	p.Log.Info(connStr)
	db, err := sqlx.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("sql Open %w", err)
	}

	return &pgRepository{
		log:    p.Log,
		db:     db,
		txOpts: &sql.TxOptions{Isolation: sql.LevelSerializable},
	}, nil
}

// program

func (r *pgRepository) AddProgram(ctx context.Context, name string, title string, active bool) (string, error) {
	// Open a new transaction to update table with data.
	tx, err := r.db.BeginTxx(ctx, r.txOpts)
	if err != nil {
		return "", fmt.Errorf("schema transaction begin %w", err)
	}
	defer tx.Rollback()

	programQuery, err := tx.PrepareNamedContext(
		ctx,
		"INSERT INTO programs (id, name, title, is_active, created_at, updated_at) VALUES (:id, :name, :title, :is_active, :created_at, :updated_at)",
	)
	if err != nil {
		return "", fmt.Errorf("PrepareNamedContext %w", err)
	}

	programId := uuid.New().String()

	_, err = programQuery.ExecContext(
		ctx,
		&domain.Program{
			ID:        programId,
			Name:      name,
			Title:     title,
			IsActive:  active,
			CreatedAt: time.Now().UTC().Unix(),
			UpdatedAt: time.Now().UTC().Unix(),
		},
	)
	if err != nil {
		return "", fmt.Errorf("program insert exec %w", err)
	}
	err = tx.Commit()
	if err != nil {
		return "", fmt.Errorf("commit transaction %w", err)
	}

	return programId, nil
}

func (r *pgRepository) UpdateProgram(ctx context.Context, id string, name *string, title *string, active *bool) error {
	tx, err := r.db.BeginTxx(ctx, r.txOpts)
	if err != nil {
		return fmt.Errorf("schema transaction begin %w", err)
	}
	defer tx.Rollback()

	query := "UPDATE programs SET "
	params := map[string]interface{}{"id": id}
	var sets []string

	if name != nil {
		sets = append(sets, "name=:name")
		params["name"] = *name
	}
	if title != nil {
		sets = append(sets, "title=:title")
		params["title"] = *title
	}
	if active != nil {
		sets = append(sets, "is_active=:is_active")
		params["is_active"] = *active
	}

	sets = append(sets, "updated_at=:updated_at")
	params["updated_at"] = time.Now().UTC().Unix()

	query += strings.Join(sets, ", ")
	query += " WHERE id=:id"

	_, err = tx.NamedExec(query, params)

	if err != nil {
		return fmt.Errorf("program update exec %w", err)
	}
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("commit transaction %w", err)
	}
	return nil
}

func (r *pgRepository) GetPrograms(ctx context.Context, page int, size int) ([]domain.Program, error) {
	programs := []domain.Program{}
	offset := (page - 1) * size
	query := "SELECT * FROM programs order by created_at LIMIT $1 OFFSET $2"

	err := r.db.Select(&programs, query, size, offset)
	return programs, err
}

func (r *pgRepository) GetProgram(ctx context.Context, programId string) (domain.Program, error) {
	program := domain.Program{}
	err := r.db.Get(&program, "SELECT * FROM programs WHERE id=$1", programId)
	return program, err
}

// member

func (r *pgRepository) AddMember(ctx context.Context,
	first_name string,
	last_name *string,
	email string,
	program_id string,
	referral_code *string,
	is_active *bool) (string, error) {
	tx, err := r.db.BeginTxx(ctx, r.txOpts)
	if err != nil {
		return "", fmt.Errorf("schema transaction begin %w", err)
	}
	defer tx.Rollback()

	programQuery, err := tx.PrepareNamedContext(
		ctx,
		"INSERT INTO members (id, first_name, last_name, email, program_id, referral_code, is_active, created_at, updated_at) VALUES (:id, :first_name, :last_name, :email, :program_id, :referral_code, :is_active, :created_at, :updated_at)",
	)
	if err != nil {
		return "", fmt.Errorf("PrepareNamedContext %w", err)
	}

	memberId := uuid.New().String()
	var code = randomLowercaseString(5)
	if referral_code != nil {
		code = *referral_code
	}
	_, err = programQuery.ExecContext(
		ctx,
		&domain.Member{
			ID:           memberId,
			FirstName:    first_name,
			LastName:     *last_name,
			Email:        email,
			ProgramId:    program_id,
			ReferralCode: code,
			IsActive:     *is_active,
			CreatedAt:    time.Now().UTC().Unix(),
			UpdatedAt:    time.Now().UTC().Unix(),
		},
	)
	if err != nil {
		return "", fmt.Errorf("member insert exec %w", err)
	}
	err = tx.Commit()
	if err != nil {
		return "", fmt.Errorf("commit transaction %w", err)
	}

	return memberId, nil
}

func (r *pgRepository) GetMembers(ctx context.Context, page int, size int) ([]domain.Member, error) {
	members := []domain.Member{}
	offset := (page - 1) * size
	query := "SELECT * FROM members order by created_at LIMIT $1 OFFSET $2"
	err := r.db.Select(&members, query, size, offset)
	return members, err
}

func (r *pgRepository) GetMember(ctx context.Context, memberId string) (domain.Member, error) {
	member := domain.Member{}
	err := r.db.Get(&member, "SELECT * FROM programs WHERE id=$1", memberId)
	return member, err
}

// referral

func (r *pgRepository) AddReferral(ctx context.Context,
	first_name *string,
	last_name *string,
	email *string,
	phone *string,
	referral_code string) (string, error) {
	tx, err := r.db.BeginTxx(ctx, r.txOpts)
	if err != nil {
		return "", fmt.Errorf("schema transaction begin %w", err)
	}
	defer tx.Rollback()

	referralQuery, err := tx.PrepareNamedContext(
		ctx,
		"INSERT INTO referrals (id, first_name, last_name, email, phone, referral_code, status, created_at, updated_at) VALUES (:id, :first_name, :last_name, :email, :phone, :referral_code, :status, :created_at, :updated_at)",
	)
	if err != nil {
		return "", fmt.Errorf("PrepareNamedContext %w", err)
	}

	referralId := uuid.New().String()

	_, err = referralQuery.ExecContext(
		ctx,
		&domain.Referral{
			ID:           referralId,
			FirstName:    *first_name,
			LastName:     *last_name,
			Email:        *email,
			Phone:        *phone,
			ReferralCode: referral_code,
			Status:       "pending",
			CreatedAt:    time.Now().UTC().Unix(),
			UpdatedAt:    time.Now().UTC().Unix(),
		},
	)
	if err != nil {
		return "", fmt.Errorf("referral insert exec %w", err)
	}
	err = tx.Commit()
	if err != nil {
		return "", fmt.Errorf("commit transaction %w", err)
	}

	return referralId, nil
}

func (r *pgRepository) GetReferrals(ctx context.Context, page int, size int) ([]domain.Referral, error) {
	referrals := []domain.Referral{}
	offset := (page - 1) * size
	query := "SELECT r.*,m.program_id as program_id, m.id as member_id FROM referrals r join members m on r.referral_code = m.referral_code order by r.created_at LIMIT $1 OFFSET $2"
	err := r.db.Select(&referrals, query, size, offset)
	return referrals, err
}

// generates random string of specified length
// helps to generate referral code
func randomLowercaseString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyz"
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}
