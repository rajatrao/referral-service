package domain

// Member corresponds to the members table
type Member struct {
	ID           string `json:"id,omitempty" db:"id"`
	FirstName    string `json:"first_name,omitempty" db:"first_name"`
	LastName     string `json:"last_name,omitempty" db:"last_name"`
	Email        string `json:"email,omitempty" db:"email"`
	ProgramId    string `json:"program_id,omitempty" db:"program_id"`
	ReferralCode string `json:"referral_code,omitempty" db:"referral_code"`
	IsActive     bool   `json:"is_active,omitempty" db:"is_active"`
	CreatedAt    int64  `json:"created_at,omitempty"  db:"created_at"`
	UpdatedAt    int64  `json:"updated_at,omitempty"  db:"updated_at"`
}
