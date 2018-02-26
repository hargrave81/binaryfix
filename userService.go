package main

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

// User is the struct for user information
type User struct {
	Name       string
	PwdHash    string
	LockedOn   time.Time
	FixAccount string
	FixPwd     string
}

// UserTable contains the users and their password hashes
var UserTable map[string]*User

func init() {
	UserTable = make(map[string]*User)
	UserTable["chris"] = &User{Name: "chris", PwdHash: "", LockedOn: time.Time{}, FixAccount: "", FixPwd: ""}
	UserTable["chris"].SetPassword("test")
}

// IsLocked returns true or false locked state of the user
func (u *User) IsLocked() bool {
	nilTime := time.Time{}
	if u.LockedOn != nilTime {
		return true
	}
	return false
}

// SetPassword will hash a new password
func (u *User) SetPassword(password string) error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return err
	}
	u.PwdHash = string(bytes)
	return nil
}

// CheckPassword will check against the password hash
func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.PwdHash), []byte(password))
	return err == nil
}
