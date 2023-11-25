package model

type User struct {
	Username string `json:"username" bson:"username"`
	Password string `json:"password" bson:"password"`
	Name     string `json:"name" bson:"name"`
	Email    string `json:"email" bson:"email"`
	City     string `json:"city" bson:"city"`
}
