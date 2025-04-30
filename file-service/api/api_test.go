package api

import (
	"bytes"
	"net/http"
	"os"
	"path/filepath"
	"testing"
)

func TestDeleteFile(t *testing.T) {
    tUserId := "asdgwtqwegsdS33"
    tFilename := "testingfile.txt"
    var tBuf bytes.Buffer
    tBuf.WriteString("TESTING!")

    saveUploadedFile(tUserId, tFilename, &tBuf) // should be OK 
    
    errMsg, code := deleteFile(tUserId, tFilename)
    if errMsg != "" && code != -1 {
	t.Fatal(errMsg)
    }
}

func TestDeleteNonExistentFile(t *testing.T) {
    tUserId := "asdgwtqwegsdS33"
    tFilename := "GSLkGSGSoiyygybaoioBYSOOEktjTkenGdkjOehwlWQLhwqltnLNglsdkfjGLhl"
    errMsg, code := deleteFile(tUserId, tFilename)
    if errMsg != "File don't exist." && code != http.StatusNoContent {
	t.Fatal()
    }
}


func TestCorrectFileSave(t *testing.T) {
    tUserId := "asdgwtqwegsdS33"
    tFilename := "testingfile.txt"
    tFile, err := os.Create(tFilename)
    if err != nil {
 	t.Fatal("Couldn't create a file.")
    }
    defer tFile.Close()
    defer os.Remove(tFilename) 

    errMsg, statusCode := saveUploadedFile(tUserId, tFilename, tFile)
    if errMsg != "" && statusCode != -1 {
 	t.Fatal()
    }
    
    deleteFile(tUserId, tFilename)
}

func TestSameFileSave(t *testing.T) {
    tUserId := "asdgwtqwegsdS33"
    tFilename := "testingfile.txt"
    tFile, err := os.Create(tFilename)
    if err != nil {
	t.Fatal("Couldn't create a file.")
    }
    defer tFile.Close()
    defer os.Remove(tFilename)

    errMsg, statusCode := saveUploadedFile(tUserId, tFilename, tFile)
    if errMsg != "" && statusCode != -1 {
	t.Fatal()
    }
    
    errMsg, statusCode = saveUploadedFile(tUserId, tFilename, tFile)
    if errMsg != "File with the same name already exist" && statusCode != http.StatusBadRequest {
	t.Fatal()
    }

    deleteFile(tUserId, tFilename)
}

func TestValidImageValidation(t *testing.T) {
    testImagePath := filepath.Join("/", "home", "qqweq", "d", "testdata", "hello.png")

    testImage, err := os.Open(testImagePath)
    if err != nil {
	t.Fatal("Couldn't open test image:", err.Error())
    }
    defer testImage.Close()

    _, isImage := validateImageFile(testImage)
    if !isImage {
	t.Fatal("Wrong return")
    }
}

func TestNotValidImageValidation(t *testing.T) {
    testNotImagePath := filepath.Join("/", "home", "qqweq", "d", "testdata", "notimage")

    testFile, err := os.Open(testNotImagePath)
    if err != nil {
	t.Fatal("Couldn't open test image:", err.Error())
    }
    defer testFile.Close()

    _, isImage := validateImageFile(testFile)
    if isImage {
	t.Fatal("Wrong return")
    }
}
