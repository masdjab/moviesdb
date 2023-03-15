package database

import (
	"crypto/rand"
	"database/sql"
	"encoding/base32"
	"errors"
	"fmt"
	"time"

	"moviesdb.com/model"
)

func (db *DbModel) GetUserByToken(token string) (*model.User, error) {
	var lastLogin int64
	user := &model.User{}
	cmd := `SELECT id, display_name, email, uname, token, last_login 
					FROM users WHERE token = '%s'`
	query := fmt.Sprintf(cmd, escape(token))
	err := db.conn.QueryRow(query).Scan(&user.UserId, &user.DisplayName, &user.Email, 
		&user.UserName, &user.Token, &lastLogin)
	
	if err == sql.ErrNoRows {
		return nil, nil
	}
	
	if err != nil {
		return nil, err
	}

	if lastLogin > 0 {
		user.LastLogin = time.Unix(lastLogin, 0)
	}

	return user, nil
}

func (db *DbModel) GetUserByUserNameAndPassword(uname string, psswd string) (*model.User, error) {
	var lastLogin int64
	var pswdMatch int

	user := &model.User{}
	cmd := `SELECT id, display_name, email, uname, token, last_login, IF(MD5('aa%szz') = psswd, 1, 0) pwd_match 
					FROM users WHERE uname = '%s'`
	query := fmt.Sprintf(cmd, escape(psswd), escape(uname))
	err := db.conn.QueryRow(query).Scan(&user.UserId, &user.DisplayName, &user.Email, 
		&user.UserName, &user.Token, &lastLogin, &pswdMatch)
	if err == sql.ErrNoRows {
		return nil, errors.New("Invalid User Name")
	}

	if err != nil {
		return nil, err
	}

	if pswdMatch == 0 {
		return nil, errors.New("Incorrect password")
	}

	user.LastLogin = time.Unix(lastLogin, 0)

	return user, nil
}

func (db *DbModel) LoginUser(uname string, psswd string) (*model.User, error) {
	user, err := db.GetUserByUserNameAndPassword(uname, psswd)
	if err != nil {
		return nil, err
	}

	for {
		token := generateRandomToken()
		tmpUser, err := db.GetUserByToken(token)
		if err != nil {
			return nil, err
		}

		if tmpUser == nil {
			cmd := "UPDATE users SET token = '%s' WHERE id = %d"
			query := fmt.Sprintf(cmd, escape(token), user.UserId)
			_, err = db.conn.Query(query)
			if err != nil {
				return nil, err
			}

			user, err = db.GetUserByToken(token)
			if err != nil {
				return nil, err
			}

			return user, nil
		}
	}
}

func (db *DbModel) LogoutUser(userId int64) {
	cmd := "UPDATE users SET last_login = 0, token = '' WHERE id = %d"
	query := fmt.Sprintf(cmd, userId)
	_, _ = db.conn.Query(query)
}

func generateRandomToken() string {
	bytes := make([]byte, 32)
	_, err := rand.Read(bytes)
	if err != nil {
			panic(err)
	}
	return base32.StdEncoding.EncodeToString(bytes)
}
