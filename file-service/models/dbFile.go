package models

import "time"

type DirEntry struct {
    Id           string `json:"id"`
    Type         string `json:"type"` // "dir" or "file"
    Name         string `json:"name"`
    FullFilepath string `json:"fullFilepath"`
}

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
