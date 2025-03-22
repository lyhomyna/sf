# /register
Method: POST

Possible responces: 

    Code: StatusOK 200
    R_JSON: -
    Cookie field with "session-id" name

    Code: StatusMethodNotAllowed 405 
    R_JSON: { message: "Use method POST instead" }

    Code: StatusBadRequest 400
    R_JSON: { message: "Use correct user schema" }

    Code: StatusInternalServerError 500 
    R_JSON: { message: "Couldn't create user: %DATABASE_ERROR_MESSAGE" }

    Code: StatusInternalServerError 500 
    R_JSON: { message: "Couldn't create session: DATABASE_ERROR_MESSAGE" } // unknown message

    Code: StatusInternalServerError 500 
    R_JSON: { message: "Can't encrypt password. BCYRPT_ERROR_MESSAGE" } // unknown message

# /login
Method: POST 

Possible responses: 

    Code: StatusOK 200
    R_JSON: - 
    Cookie field with "session-id" name

    Code: StatusMethodNotAllowed 405
    R_JSON: { message: "Use method POST instead" }

    Code: StatusBadRequest 400
    R_JSON: { message: "Couldn't parse user" }

    Code: StatusBadRequest 400
    R_JSON: { message: "User data can't be blank line" }

    Code: StatusBadRequest 400 
    R_JSON: { message: "Password should be at least 6 chars length" }

    Code: StatusBadRequest 400 
    R_JSON: { message: "Password shouldn't contain ' or \"" }

    Code: StatusNotFound 404 
    R_JSON: { message: "User not found" }

    Code: StatusBadRequest 403 
    R_JSON: { message: "Passwords don't match" }

    Code: StatusInternalServerError 500 
    R_JSON: { message: "Internal server error" }

    Code: StatusInternalServerError 500 
    R_JSON: { message: "Couldn't create session: DATABASE_ERROR_MESSAGE" } // unknown message


# /logout
Method: GET

Possible responses:

    Code: StatusOK 200

    Code: StatusUnauthrorized 401
    R_JSON: { message: "Message from Go's http library" } // unknown message

    Code: statusinternalservererror 500
    R_JSON: { message: "Couldn't delete session: %MESSAGE_FROM_DB" } // unknown message

    Code: statusinternalservererror 500
    R_JSON: { message: "Nothing was deleted from sessions table" }

# /check-auth
Method:  GET

Possible responses:

    Code: StatusOK 200

    Code: StatusUnauthorized 401
