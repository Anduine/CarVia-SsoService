package domain

type RegisterRequest struct {
	UserID      int    `json:"-"`
	Login       string `json:"Login"`
	Password    string `json:"Password"`
	Email       string `json:"Email"`
	Role        string `json:"Role"`
	FirstName   string `json:"FirstName"`
	LastName    string `json:"LastName"`
	Phonenumber string `json:"Phonenumber"`
	Address     string `json:"Address"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserUpdateRequest struct {
	UserID       int    `json:"-"`
	Login        string `json:"Login"`
	Password     string `json:"Password"`
	HashPassword string `json:"-"`
	Email        string `json:"Email"`
	FirstName    string `json:"FirstName"`
	LastName     string `json:"LastName"`
	Phonenumber  string `json:"Phonenumber"`
	Address      string `json:"Address"`
	AvatarPath   string `json:"AvatarPath"`
}
