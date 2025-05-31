package models

import "os"

type ImageData struct {
    ImageFile  		*os.File
    ContentTypeChunk 	[]byte
}

type ImageJson struct {
    ImageUrl 	string 	`json:"imageUrl"`
}
