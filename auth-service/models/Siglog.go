package models

type Siglog struct {
    Users UserDao
    Sessions SessionDao
}

type UserDao interface {
    CreateUser(user *DbUser) (string, error)
    DeleteUser(userId string) error
    ReadUserById(id string) (*User, error)
    FindUser(user *User) (string, error)
    GetUserByEmail(email string) (*DbUser, error)
}

type SessionDao interface {
    CreateSession(userId string) (string, error)
    DeleteSession(sessionId string) error
    UserIdFromSessionId(sessionId string) (string, error)
}
