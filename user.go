package userchat

import (
	"errors"
	"strings"
)

type User struct {
	Id       int
	Username string
	Password string
}

var (
	ErrUserNotFound = errors.New("user not found")
)

// In a real app this would be stored in a database
var DataUsers = []*User{
	{
		Id:       23,
		Username: "Wolverine",
		Password: "password",
	},
	{
		Id:       48,
		Username: "Storm",
		Password: "password",
	},
}

func GetUsers() []*User {
	return DataUsers
}

func GetUserById(id int) (*User, error) {
	for _, v := range DataUsers {
		if v.Id == id {
			return v, nil
		}
	}
	return nil, ErrUserNotFound
}

func GetUserByUsername(username string) (*User, error) {
	username = strings.ToLower(username)
	for _, v := range DataUsers {
		if strings.ToLower(v.Username) == username {
			return v, nil
		}
	}
	return nil, ErrUserNotFound
}

func ComparePassword(hashed, given string) bool {
	return hashed == given
}
