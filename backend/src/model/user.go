package model

type User struct {
	Username  string `json:"username" bson:"username"`
	Password  string `json:"password" bson:"password"`
	Email     string `json:"email" bson:"email"`
	CreatedAt string `json:"created" bson:"created"`
	UpdatedAt string `json:"updated" bson:"updated"`
}
