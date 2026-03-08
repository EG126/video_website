package mysql

import (
	"context"
	"video_website/biz/dal/mysql/entity"
)

func FollowAction(ctx context.Context, followerID, followingID string, actionType int32) error {
	if actionType == 1 { // 关注
		follow := &entity.Follow{
			FollowerID:  followerID,
			FollowingID: followingID,
		}
		return DB.WithContext(ctx).FirstOrCreate(follow).Error
	} else if actionType == 2 { // 取消关注
		return DB.WithContext(ctx).Where("follower_id = ? AND following_id = ?", followerID, followingID).Delete(&entity.Follow{}).Error
	}
	return nil
}

func GetFollowing(ctx context.Context, userID string, pageNum, pageSize int32) ([]string, int64, error) {
	var follows []*entity.Follow
	var total int64
	db := DB.WithContext(ctx).Model(&entity.Follow{}).Where("following_id = ?", userID)
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	offset := (pageNum - 1) * pageSize
	err := db.Offset(int(offset)).Limit(int(pageSize)).Order("created_at desc").Find(&follows).Error
	if err != nil {
		return nil, 0, err
	}
	var ids []string
	for _, f := range follows {
		ids = append(ids, f.FollowerID)
	}
	return ids, total, nil
}

func GetFollowers(ctx context.Context, userID string, pageNum, pageSize int32) ([]string, int64, error) {
	var follows []*entity.Follow
	var total int64
	db := DB.WithContext(ctx).Where("follower_id = ?", userID)
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	offset := (pageNum - 1) * pageSize
	err := db.Offset(int(offset)).Limit(int(pageSize)).Order("created_at").Find(&follows).Error
	if err != nil {
		return nil, 0, err
	}
	var ids []string
	for _, f := range follows {
		ids = append(ids, f.FollowingID)
	}
	return ids, total, nil
}

func GetFriends(ctx context.Context, userID string, pageNum, pageSize int32) ([]string, int64, error) {
	subQuery := DB.WithContext(ctx).Model(&entity.Follow{}).Where("follower_id = ?", userID).Select("following_id")
	var total int64
	var friends []string
	err := DB.WithContext(ctx).Model(&entity.Follow{}).Where("follower_id IN (?) AND following_id = ?", subQuery, userID).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}
	offset := (pageNum - 1) * pageSize
	// Pluck:提取单列
	err = DB.WithContext(ctx).Model(&entity.Follow{}).Where("follower_id IN (?) AND following_id = ?", subQuery, userID).Offset(int(offset)).Limit(int(pageSize)).Pluck("follower_id", &friends).Error
	return friends, total, err
}
