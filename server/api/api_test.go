package api 

import (
	"net/http"
	"os"
	"testing"
)

func TestCorrectFileSave(t *testing.T) {
    tFilename := "testingfile.txt"
    tFile, err := os.Create(tFilename)
    if err != nil {
 	t.Fatal("Couldn't create a file.")
    }
    defer tFile.Close()
    defer os.Remove(tFilename) 

    errMsg, statusCode := saveUploadedFile(tFilename, tFile)
    if errMsg != "" && statusCode != -1 {
 	t.Fatal()
    }
    
    deleteFile(tFilename)
}

func TestSameFileSave(t *testing.T) {
    tFilename := "testingfile.txt"
    tFile, err := os.Create(tFilename)
    if err != nil {
	t.Fatal("Couldn't create a file.")
    }
    defer tFile.Close()
    defer os.Remove(tFilename)

    errMsg, statusCode := saveUploadedFile(tFilename, tFile)
    if errMsg != "" && statusCode != -1 {
	t.Fatal()
    }
    
    errMsg, statusCode = saveUploadedFile(tFilename, tFile)
    if errMsg != "File with the same name already exist" && statusCode != http.StatusBadRequest {
	t.Fatal()
    }

    deleteFile(tFilename)
}

func TestDeleteFile(t *testing.T) {
    tFilename := "testingfile.txt"
    tFile, err := os.Create(tFilename)
    if err != nil {
	t.Fatal("Couldn't create a file.")
    }
    defer tFile.Close()
    defer os.Remove(tFilename)

    saveUploadedFile(tFilename, tFile) // should be OK 
    
    errMsg, code := deleteFile(tFilename)
    if errMsg != "" && code != -1 {
	t.Fatal()
    }
}

func TestDeleteNonExistentFile(t *testing.T) {
    tFilename := "GSLkGSGSoiyygybaoioBYSOOEktjTkenGdkjOehwlWQLhwqltnLNglsdkfjGLhl"
    errMsg, code := deleteFile(tFilename)
    if errMsg != "File don't exist." && code != http.StatusNoContent {
	t.Fatal()
    }
}

