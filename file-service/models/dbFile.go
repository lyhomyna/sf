package models

import "time"

type DbUserFile struct {
    Id 			string
    UserId 		string
    Filename		string
    Filepath		string
    Size		int64
    Hash		string
    LastAccessed	time.Time
}

type UserFile struct {
    Id	string		`json:"id"`
    Filename string	`json:"filename"`
    LastAccessed int64	`json:"lastAccessed"`
}
