package models

type Siglog struct {
    Users UserDao
    Sessions SessionDao
}

type UserDao interface {
    CreateUser(user *User)
    DeleteUser(user *User) error
    ReadUserById(id string) (*User, error)
}

type SessionDao interface {
    CreateSession(userId string) (string, error)
    DeleteSession(sessionId string) error
    UserIdFromSessionId(sessionId string) (string, error)
}
