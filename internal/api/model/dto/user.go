package dto

type UserRegister struct {
	DisplayName     *string `validate:"omitempty,min=3,max=52" json:"display_name"`
	Username        string  `validate:"required,username" json:"username"`
	Email           string  `validate:"required,email" json:"email"`
	Password        string  `validate:"required,password" json:"password"`
	PasswordConfirm string  `validate:"required" json:"password_confirm"`
}

type UserLogin struct {
	Login    string `validate:"required" json:"login"`
	Password string `validate:"required" json:"password"`
}

func (u *UserRegister) PasswordCheck() bool {
	return u.Password == u.PasswordConfirm
}
