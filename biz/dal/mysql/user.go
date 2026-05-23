package mysql

import (
	"context"
	"video_website/biz/dal/mysql/entity"

	pkgErrors "github.com/pkg/errors"
	gorm "gorm.io/gorm"
)

func CreateUser(ctx context.Context, user *entity.User) error {
	err := DB.WithContext(ctx).Create(user).Error
	if err != nil {
		return pkgErrors.Wrapf(err, "CreateUser failed, username=%s", user.Username)
	}
	return nil
}

func GetUserByUsername(ctx context.Context, username string) (*entity.User, error) {
	var user entity.User
	err := DB.WithContext(ctx).Where("username = ?", username).First(&user).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, pkgErrors.Wrapf(err, "GetUserByUsername failed, username=%s", username)
	}
	return &user, nil
}

func GetUserByID(ctx context.Context, id string) (*entity.User, error) {
	var user entity.User
	err := DB.WithContext(ctx).Where("id = ?", id).First(&user).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, pkgErrors.Wrapf(err, "GetUserByID failed, id=%s", id)
	}
	return &user, nil
}

func UpdateUserAvatar(ctx context.Context, userID, avatarURL string) error {
	err := DB.WithContext(ctx).Model(&entity.User{}).Where("id = ?", userID).Update("avatar_url", avatarURL).Error
	if err != nil {
		return pkgErrors.Wrapf(err, "UpdateUserAvatar failed, userID=%s", userID)
	}
	return nil
}

func IsUserExist(ctx context.Context, username string) (bool, error) {
	var count int64
	err := DB.WithContext(ctx).Model(&entity.User{}).Where("username = ?", username).Count(&count).Error
	if err != nil {
		return false, pkgErrors.Wrapf(err, "IsUserExist failed, username=%s", username)
	}
	return count > 0, nil
}

func UpdateUserMFA(ctx context.Context, userID, mfaSecret string) error {
	err := DB.WithContext(ctx).Model(&entity.User{}).Where("id = ?", userID).Update("mfa_secret", mfaSecret).Error
	if err != nil {
		return pkgErrors.Wrapf(err, "UpdateUserMFA failed, userID=%s", userID)
	}
	return nil
}

func GetUserMFASecret(ctx context.Context, userID string) (string, error) {
	var user entity.User
	err := DB.WithContext(ctx).Select("mfa_secret").Where("id = ?", userID).First(&user).Error
	if err == gorm.ErrRecordNotFound {
		return "", nil
	}
	if err != nil {
		return "", pkgErrors.Wrapf(err, "GetUserMFASecret failed, userID=%s", userID)
	}
	return user.MFASecret, nil
}
