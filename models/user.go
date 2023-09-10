// models/user.go

package models

import "time"

type User struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Password  string    `json:"-"`
	Phone     string    `json:"phone"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `json:"deleted_at"`
}

// NewUser creates a new User instance.
func NewUser(name, email, password, phone, status string) *User {
	return &User{
		Name:      name,
		Email:     email,
		Password:  password,
		Phone:     phone,
		Status:    status,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		DeletedAt: time.Time{}, // Assuming initially deleted_at is empty
	}
}
