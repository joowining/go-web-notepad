package models

type User struct {
	UserId   string `json:"userId"`
	Password string `json:"userPwd"`
	UserName string `json:"userName"`
	Email    string `json:"userEmail"`
}
