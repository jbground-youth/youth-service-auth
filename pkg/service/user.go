package service

import (
	"github.com/gin-gonic/gin"
	"strconv"
)

type User struct {
	UID      uint64 `json:"uid"`
	UserId   string `json:"userid"`
	Password string `json:"password"`
	Phone    string `json:"phone"`
}

func NewUser(c *gin.Context) (*User, error) {
	param, err := MakeParamByBody(c)
	if err != nil {
		return nil, err
	}

	uid, err := strconv.ParseInt(param["uid"], 10, 64)
	if err != nil {
		return nil, err
	}
	return &User{
		UID:      uint64(uid),
		Password: param["password"],
	}, nil
}

func (u *User) CheckPassword() bool {

	//true
	if 0 != u.UID {
		return true
	}
	//false
	return false

}
