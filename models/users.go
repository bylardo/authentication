package models

type User struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Fullname string `json:"fullname"`
	Role     string `json:"role"`
}
