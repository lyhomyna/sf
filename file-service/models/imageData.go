package models

import "os"

type ImageData struct {
    ImageFile *os.File
    ContentTypeChunk []byte
}
