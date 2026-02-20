package domain

import "context"

type User struct {
	UserID       int    `json:"UserID"`
	Login        string `json:"Login"`
	Role         string `json:"Role"`
	FirstName    string `json:"FirstName"`
	LastName     string `json:"LastName"`
	Email        string `json:"Email"`
	HashPassword string `json:"-"`
	Address      string `json:"Address"`
	Phonenumber  string `json:"Phonenumber"`
	AvatarPath   string `json:"AvatarPath"`
}

type UserRepository interface {
	CreateUser(ctx context.Context, user User) (int, error)
	GetByEmail(ctx context.Context, email string) (User, error)
	ExistsByEmail(ctx context.Context, email string) (bool, error)
	GetByUsername(ctx context.Context, username string) (User, error)
	ExistsByUsername(ctx context.Context, username string) (bool, error)

	UpdateUserProfile(ctx context.Context, userData UserUpdateRequest) error
}
