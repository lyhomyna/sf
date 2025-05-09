package models

import "time"

type DbUser struct {
    Id 		string 	   `json:"id"`
    Email 	string 	   `json:"email"`
    Password 	string 	   `json:"password"`
    CreatedAt   time.Time  `json:"createdAt"`
    ImageUrl    string     `json:"imageUrl"`
}

type User struct {
    Email 	string 	`json:"email"`
    Password 	string 	`json:"password"`
    ImageUrl	string	`json:"imageUrl"`
}
