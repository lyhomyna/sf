package models

import "time"

type DbUser struct {
    Id 		string 	   `json:"id"`
    Email 	string 	   `json:"email"`
    Password 	string 	   `json:"password"`
    CreatedAt   time.Time  `json:"created_at"`
    ImageUrl    string     `json:"image_url"`
}

type User struct {
    Email 	string 	`json:"email"`
    Password 	string 	`json:"password"`
}
