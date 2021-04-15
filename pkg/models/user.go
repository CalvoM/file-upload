package models

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/CalvoM/file-upload/pkg/db"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID         uint64
	UserName   string
	Password   string
	Email      string
	Registered string
}

//IsUserInDb check if user exists in Db
func IsUserInDb(user *User) (isUserPresent bool) {
	rows := db.DbClient.QueryRow(`select user_id from reg_users where email=$1;`, user.Email)
	isUserPresent = false
	var userID uint64
	switch err := rows.Scan(&userID); err {
	case sql.ErrNoRows:
		log.Warn("User not found")
	case nil:
		isUserPresent = true
	default:
		log.Error("Error getting user ", err)
		isUserPresent = true
	}
	return

}
func (user *User) AddUser() (userID uint64, err error) {
	if IsUserInDb(user) {
		err = errors.New("User in DB")
		return
	}
	stmt, err := db.DbClient.Prepare("insert into reg_users(username,password,email) values($1,$2,$3) returning user_id")
	if err != nil {
		log.Error(err)
		return
	}
	defer stmt.Close()
	user.HashUserPassword()
	err = stmt.QueryRow(user.UserName, user.Password, user.Email).Scan(&userID)
	if err != nil {
		log.Error(err)
		return
	}
	user.ID = userID
	return
}

//GetUser get the user when given the email and password
func (user *User) GetUser(email, password string) (err error) {
	stmt, err := db.DbClient.Prepare("select * from reg_users where email=$1")
	if err != nil {
		log.Error(err)
		return
	}
	err = stmt.QueryRow(email).Scan(&user.ID, &user.UserName, &user.Password,
		&user.Email, &user.Registered)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer stmt.Close()
	if userPresent := user.CheckUserPassword(password); !userPresent {
		return errors.New("Passwords do not match!")
	}
	return
}

//HashUserPassword hashes the user plaintext password
func (user *User) HashUserPassword() {
	if b, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost); err != nil {
		log.Error(err)
	} else {
		user.Password = string(b)
	}
}

//CheckUserPassword check the user hashed password against the one supplied
func (user *User) CheckUserPassword(suppliedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(suppliedPassword))
	return err == nil
}
