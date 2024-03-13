package models

import (
	"fmt"
	"github.com/youth-service/auth/pkg/setting"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
)

var db *gorm.DB

type Model struct {
	ID         int `gorm:"primary_key" json:"id"`
	CreatedOn  int `json:"created_on"`
	ModifiedOn int `json:"modified_on"`
	DeletedOn  int `json:"deleted_on"`
}

func Setup() {
	dns := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local", setting.DatabaseSetting.User,
		setting.DatabaseSetting.Password, setting.DatabaseSetting.Host, setting.DatabaseSetting.Name)

	myconfig := mysql.New(mysql.Config{
		DSN:                       dns,   // data source name
		DefaultStringSize:         256,   // default size for string fields
		DisableDatetimePrecision:  true,  // disable datetime precision, which not supported before MySQL 5.6
		DontSupportRenameIndex:    true,  // drop & create when rename index, rename index not supported before MySQL 5.7, MariaDB
		DontSupportRenameColumn:   true,  // `change` when rename column, rename column not supported before MySQL 8, MariaDB
		SkipInitializeWithVersion: false, // auto configure based on currently MySQL version
	})

	db, err := gorm.Open(myconfig, &gorm.Config{})
	if err != nil {
		log.Fatalf("mariadb.Setup err: %v", err)
	}

	if err = db.AutoMigrate(User{}); err != nil {
		log.Fatalf("mariadb.Setup err: %v", err)
	}
}
