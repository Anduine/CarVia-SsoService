package domain

import "context"

type RegisterRequest struct {
	UserID      int     `json:"user_id"`
	Login    	  string  `json:"login"`
	Password 	  string  `json:"password"`
	Email    	  string  `json:"email"`
	Role 	 	    string  `json:"role"`
	FirstName	  string	`json:"first_name"`
	LastName	  string	`json:"last_name"`
	Phonenumber string	`json:"phonenumber"`
	Address		  string	`json:"address"`
}

type LoginRequest struct {
	Login    	 string  `json:"login"`
	Email    	 string  `json:"email"`
	Password 	 string  `json:"password"`
}

type UserUpdateRequest struct {
	UserID    	 int
	Login      	 string
	Password   	 string
	HashPassword string 
	FirstName 	 string
	LastName  	 string
	Email    	 	 string
	Phone    	   string
	Address  	   string
	AvatarPath  *string
}

type User struct {
	UserID       int	  `json:"user_id"`
	Login		     string	`json:"login"`
	Role         string	`json:"role"`
	FirstName	   string	`json:"first_name"`
	LastName	   string	`json:"last_name"`
	Email        string	`json:"email"`
	HashPassword string	`json:"-"`
	Address		   string	`json:"address"`
	Phone        string	`json:"phonenumber"`
	AvatarPath   string	`json:"avatar"`
}

type UserRepository interface {
	CreateUser(ctx context.Context, user User) (int, error)
	GetByUsername(ctx context.Context, email string) (User, error)
	ExistsByEmail(ctx context.Context, email string) (bool, error)
	ExistsByUsername(ctx context.Context, username string) (bool, error)

	UpdateUserProfile(ctx context.Context, userData UserUpdateRequest) error
}
