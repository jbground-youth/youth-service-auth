package models

type User struct {
	Model

	ID       uint64 `gorm:"ID" json:"id"`
	UserId   string `gorm:"USERID" json:"username"`
	Password string `gorm:"PASSWORD" json:"password"`
	Phone    string `gorm:"PHONE" json:"phone"`
}

func (u User) FindUser(userid string) *User {
	var user = User{}

	db.Where("id = ?", userid).First(&user)

	return &user
}
