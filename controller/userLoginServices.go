package controller

import (
	"errors"
	"github.com/sql-queries/common"
	"github.com/sql-queries/models"
	"github.com/sql-queries/server"
	"strings"
)

type User models.User

type UserLogin struct {
	Email     string `gorm:"type:varchar(100);unique_index"`
	Password  string
	SessionId string
	IsActive  bool
}

type Deactivate struct {
	Email    string `gorm:"type:varchar(100);unique_index"`
	IsActive bool
}

func (self *UserLogin) Format() *UserLogin {
	self.Email = strings.ToLower(self.Email)
	return self
}

func (self *User) FindByEmail(mail string) (error, *User) {
	newUser := &User{}
	queryResult := server.Conn().Where(&User{Email: mail}).First(newUser)
	if queryResult.Error != nil {
		return errors.New("Error while connecting to database "), nil
	} else {
		return nil, newUser
	}
}

func (self *User) GetCurrentUserFromHeaders(SessionID string) (error, string) {
	user := &User{}

	currentUser := server.Conn().Where(&User{SessionId: SessionID}).First(user)
	if currentUser.Error != nil {
		return errors.New("Error while connecting to database "), ""
	}

	return nil, user.Email
}

func (self *UserLogin) ValidateLogin() (string, *User) {
	user := &User{}
	if self.Email != "" && self.Password != "" {
		_, user = user.FindByEmail(self.Email)

		if user == nil || user.Email == "" {
			return "Error login, user doesn't exist ", nil
		}

		if self.Email != user.Email {
			return "Error login, Wrong email or password ", nil
		}

		wrongPassword := common.CheckPasswordHash(self.Password, user.Password)

		if !wrongPassword {
			return "Error login, Wrong email or password", nil
		}

		if !user.IsActive {
			return "Please reactivate your account or call customer support", nil
		}

	} else {
		return "Please insert email and password to login", nil
	}
	user.Password = ""
	return "", user
}
