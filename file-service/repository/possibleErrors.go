package repository 

import "errors"

var FilesErrorFailureToRetrieve = errors.New("Failure to retrieve file from DB")
var FilesErrorDbQuery = errors.New("Database query error")
var FilesErrorFileNotExist = errors.New("File doesn't exist")
var FilesErrorFileExist = errors.New("File with the save name already exists")
var FilesErrorInternal = errors.New("Internal error")
var FilesErrorFailureCreateFile = errors.New("Couldn't create new file")
var FilesErrorCopyFailure = errors.New("File couldn't be saved")
var FilesErrorDbSave = errors.New("Couldn't save file to the database")
var FilesErrorNoFilesForUser = errors.New("There is no files for user")
var FilesErrorDbScan = errors.New("Failed to scan database row")


var ErrorRootDirNotFound = errors.New("Root directory not found")
var ErrorDirectoryNotFound = errors.New("Directory not found")
var ErrorGetDirFailed = errors.New("Failed to get directory id")

var ErrorPathNotFound = errors.New("Directory path not found")
var ErrorQueryFailed = errors.New("Failed to query database")
var ErrorScanFailed = errors.New("Failed to scan row")
