package models

type Dir struct {
    Id       string 	`json:"id"`
    Name     string 	`json:"name"`
    FullPath string 	`json:"fullPath"`
}

type DirEntry struct {
    Id           string `json:"id"`
    Type         string `json:"type"` // "dir" or "file"
    Name         string `json:"name"`
    FullFilepath string `json:"path"`
}
