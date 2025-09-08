package domain

// Referral corresponds to the referrals table
type Referral struct {
	ID           string `json:"id,omitempty" db:"id"`
	FirstName    string `json:"first_name,omitempty" db:"first_name"`
	LastName     string `json:"last_name,omitempty" db:"last_name"`
	Email        string `json:"email,omitempty" db:"email"`
	Phone        string `json:"phone,omitempty" db:"phone"`
	ReferralCode string `json:"referral_code,omitempty" db:"referral_code"`
	Status       string `json:"status,omitempty" db:"status"`
	CreatedAt    int64  `json:"created_at,omitempty"  db:"created_at"`
	UpdatedAt    int64  `json:"updated_at,omitempty"  db:"updated_at"`
	ProgramId    string `json:"program_id,omitempty" db:"program_id"`
	MemberId     string `json:"member_id,omitempty" db:"member_id"`
}
