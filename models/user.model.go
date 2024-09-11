package models

type User struct {
	Username string `bson:"username" json:"username" validate:"required"`
	Email    string `bson:"email" json:"email" validate:"required,email"`
	Password string `bson:"password" json:"password" validate:"required,min=8"`
}

type LoginData struct {
	Email    string `bson:"email" json:"email" validate:"required,email"`
	Password string `bson:"password" json:"password" validate:"required,min=8"`
}
