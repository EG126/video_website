package mysql

import (
	"context"
	"video_website/biz/dal/mysql/entity"

	"gorm.io/gorm"
)

func CreateUser(ctx context.Context, user *entity.User) error {
	return DB.WithContext(ctx).Create(user).Error
}

func GetUserByUsername(ctx context.Context, username string) (*entity.User, error) {
	var user entity.User
	err := DB.WithContext(ctx).Where("username = ?", username).First(&user).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &user, err
}

func GetUserByID(ctx context.Context, id string) (*entity.User, error) {
	var user entity.User
	err := DB.WithContext(ctx).Where("id = ?", id).First(&user).Error
	return &user, err
}

func UpdateUserAvatar(ctx context.Context, userID, avatarURL string) error {
	return DB.WithContext(ctx).Model(&entity.User{}).Where("id = ?", userID).Update("avatar_url", avatarURL).Error
}

func IsUserExist(ctx context.Context, username string) (bool, error) {
	var count int64
	err := DB.WithContext(ctx).Model(&entity.User{}).Where("username = ?", username).Count(&count).Error
	return count > 0, err
}
