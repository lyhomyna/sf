package main

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

    errMsg, statusCode := saveUploadedFile(tFilename, tFile)
    if errMsg != "" && statusCode != -1 {
 	t.Fatal()
    }

    os.Remove(tFilename)
    // TODO: SOMEHOW REMOVE FILE FROM SERVER
}

func TestSameFileSave(t *testing.T) {
    tFilename := "testingfile.txt"
    tFile, err := os.Create(tFilename)
    if err != nil {
	t.Fatal("Couldn't create a file.")
    }
    defer tFile.Close()

    errMsg, statusCode := saveUploadedFile(tFilename, tFile)
    if errMsg != "" && statusCode != -1 {
	t.Fatal()
    }
    
    errMsg, statusCode = saveUploadedFile(tFilename, tFile)
    if errMsg != "File with the same name already exist" && statusCode != http.StatusBadRequest {
	t.Fatal()
    }

    os.Remove(tFilename)
    // TODO: SOMEHOW REMOVE FILE FROM SERVER
}

