# /save
Method: POST
***multipart/form-data***. FormData should include variable with name **file**

Possible responces: 

    Code: StatusOK 200
    R_JSON: { data: "File ${filename} saved" }

    Code: StatusBadRequest 400 
    R_JSON: { data: "Use POST method instead" }

    Code: StatusBadRequest 400
    R_JSON: { data: "Provide a file" }

    Code: StatusBadRequest 400
    R_JSON: { data: "File with the same name already exist" }

    Code: StatusUnauthorized 401
    R_JSON: { data: "Session cookie missing" }

    Code: StatusUnsupportedMediaType 415
    R_JSON: { data: "Expected multipart/form-data"}

    Code: StatusInternalServerError 500
    R_JSON: { data: "Something went wrong while a server was creating file" }

    Code: StatusInternalServerError 500
    R_JSON: { data: "Something went wrong when the server was copying content from your file" }


# /delete/{filename}
Method: DELETE

Possible responses: 

    Code: StatusOK 200
    R_JSON: { data: "File deleted" }

    Code: StatusNoContent 204
    R_JSON: { data: "File doesn't exist" }

    Code: StatusBadRequest 400
    R_JSON: { data: "Use DELETE method instead" }

    Code: StatusInternalServerError 500
    R_JSON: { data: "Error removing file" }


# /download
Method: GET

Possible responses:

    Code: StatusOK 200
    Data: selected file in HTTP body 

    Code: StatusBadRequest 400
    R_JSON: { data: "Use GET method instead" }

    Code: StatusBadRequest 400
    R_JSON: { data: "File not specified" }
    
    Code: StatusNotFound 404
    R_JSON: { data: "File not found" }


# /filenames
Method: GET

Possible responses:

    Code: StatusOK 200
    R_JSON: { data: "["filename1", "filename2"]" }

    Code: StatusBadRequest 400
    R_JSON: { data: "Use GET method instead" }

    Code: StatusInternalServerError 500
    R_JSON: { data: "Cannot read from user directory" }

    Code: StatusUnauthorized 401 
    R_JSON: { data: "Session cookie missing" }

# /health
Method: GET

Possible response:
    Code: StatusOK 200
