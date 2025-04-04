package models

type UserLoginRequest struct {
	Login    string `json:"login" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UserChangePassword struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required"`
}

type UserRegisterRequest struct {
	FIO      string   `json:"fio" binding:"required"`
	Login    string   `json:"login" binding:"required"`
	Password string   `json:"password" binding:"required"`
	Roles    []string `json:"roles" binding:"required"`
}
