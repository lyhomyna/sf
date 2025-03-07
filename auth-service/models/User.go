package models

type DbUser struct {
    Id 		string 	`json:"id"`
    Email 	string 	`json:"email"`
    Password 	string 	`json:"password"`
}

type User struct {
    Email 	string 	`json:"email"`
    Password 	string 	`json:"password"`
}
