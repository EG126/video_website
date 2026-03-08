package mysql

import (
	"video_website/biz/dal/mysql/entity"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Init(dsn string) error {
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}
	DB = db
	return DB.AutoMigrate(&entity.User{}, &entity.Video{}, &entity.Like{}, &entity.Comment{}, &entity.Follow{})
}
