package domain

// Referral Program corresponds to the program table
type Program struct {
	ID        string `json:"id,omitempty" db:"id"`
	Name      string `json:"name,omitempty" db:"name"`
	Title     string `json:"title,omitempty" db:"title"`
	IsActive  bool   `json:"is_active,omitempty" db:"is_active"`
	CreatedAt int64  `json:"created_at,omitempty"  db:"created_at"`
	UpdatedAt int64  `json:"updated_at,omitempty"  db:"updated_at"`
}
