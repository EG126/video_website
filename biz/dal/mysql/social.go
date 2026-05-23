package mysql

import (
	"context"
	"video_website/biz/dal/mysql/entity"
	"video_website/pkg/utils"

	pkgErrors "github.com/pkg/errors"
)

func FollowAction(ctx context.Context, followerID, followingID string, actionType int32) error {
	switch actionType {
	case 1:
		follow := &entity.Follow{
			ID:          utils.GenerateID(),
			FollowerID:  followerID,
			FollowingID: followingID,
		}
		err := DB.WithContext(ctx).FirstOrCreate(follow, "follower_id = ? AND following_id = ?", followerID, followingID).Error
		if err != nil {
			return pkgErrors.Wrapf(err, "FollowAction create failed, follower=%s, following=%s", followerID, followingID)
		}
		return nil
	case 2:
		err := DB.WithContext(ctx).Where("follower_id = ? AND following_id = ?", followerID, followingID).Delete(&entity.Follow{}).Error
		if err != nil {
			return pkgErrors.Wrapf(err, "FollowAction delete failed, follower=%s, following=%s", followerID, followingID)
		}
		return nil
	}
	return nil
}

func GetFollowing(ctx context.Context, userID string, pageNum, pageSize int32) ([]string, int64, error) {
	var follows []*entity.Follow
	var total int64
	db := DB.WithContext(ctx).Model(&entity.Follow{}).Where("follower_id = ?", userID)
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, pkgErrors.Wrapf(err, "GetFollowing count failed, userID=%s", userID)
	}
	offset := (pageNum - 1) * pageSize
	err := db.Offset(int(offset)).Limit(int(pageSize)).Order("created_at desc").Find(&follows).Error
	if err != nil {
		return nil, 0, pkgErrors.Wrapf(err, "GetFollowing query failed, userID=%s", userID)
	}
	var ids []string
	for _, f := range follows {
		ids = append(ids, f.FollowingID)
	}
	return ids, total, nil
}

func GetFollowers(ctx context.Context, userID string, pageNum, pageSize int32) ([]string, int64, error) {
	var follows []*entity.Follow
	var total int64
	db := DB.WithContext(ctx).Model(&entity.Follow{}).Where("following_id = ?", userID)
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, pkgErrors.Wrapf(err, "GetFollowers count failed, userID=%s", userID)
	}
	offset := (pageNum - 1) * pageSize
	err := db.Offset(int(offset)).Limit(int(pageSize)).Order("created_at").Find(&follows).Error
	if err != nil {
		return nil, 0, pkgErrors.Wrapf(err, "GetFollowers query failed, userID=%s", userID)
	}
	var ids []string
	for _, f := range follows {
		ids = append(ids, f.FollowerID)
	}
	return ids, total, nil
}

func GetFriends(ctx context.Context, userID string, pageNum, pageSize int32) ([]string, int64, error) {
	subQuery := DB.WithContext(ctx).Model(&entity.Follow{}).Where("follower_id = ?", userID).Select("following_id")
	var total int64
	var friends []string
	err := DB.WithContext(ctx).Model(&entity.Follow{}).Where("follower_id IN (?) AND following_id = ?", subQuery, userID).Count(&total).Error
	if err != nil {
		return nil, 0, pkgErrors.Wrapf(err, "GetFriends count failed, userID=%s", userID)
	}
	offset := (pageNum - 1) * pageSize
	err = DB.WithContext(ctx).Model(&entity.Follow{}).Where("follower_id IN (?) AND following_id = ?", subQuery, userID).Offset(int(offset)).Limit(int(pageSize)).Pluck("follower_id", &friends).Error
	if err != nil {
		return nil, 0, pkgErrors.Wrapf(err, "GetFriends query failed, userID=%s", userID)
	}
	return friends, total, nil
}
