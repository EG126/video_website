package mysql

import (
	"context"
	"time"
	"video_website/biz/dal/mysql/entity"
	"video_website/pkg/utils"
)

func LikeAction(ctx context.Context, userID, videoID string, actionType int32) error {
	if actionType == 1 {
		like := &entity.Like{
			ID:        utils.GenerateID(),
			UserID:    userID,
			VideoID:   videoID,
			CreatedAt: time.Now(),
		}
		return DB.WithContext(ctx).FirstOrCreate(like, "user_id = ? AND video_id = ?", userID, videoID).Error
	} else if actionType == 2 {
		return DB.WithContext(ctx).Where("user_id = ? AND video_id = ?", userID, videoID).Delete(&entity.Like{}).Error
	}
	return nil
}

func GetUserLikedVideoIDs(ctx context.Context, userID string, pageNum, pageSize int32) ([]string, int64, error) {
	var likes []*entity.Like
	var total int64
	db := DB.WithContext(ctx).Model(&entity.Like{}).Where("user_id = ?", userID)
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	offset := (pageNum - 1) * pageSize
	err := db.Offset(int(offset)).Limit(int(pageSize)).Order("created_at desc").Find(&likes).Error
	if err != nil {
		return nil, 0, err
	}
	var ids []string
	for _, l := range likes {
		ids = append(ids, l.VideoID)
	}
	return ids, total, nil
}

func GetVideosByIDs(ctx context.Context, videoIDs []string) ([]*entity.Video, error) {
	if len(videoIDs) == 0 {
		return []*entity.Video{}, nil
	}
	var videos []*entity.Video
	err := DB.WithContext(ctx).Where("id IN ?", videoIDs).Find(&videos).Error
	return videos, err
}

func CreateComment(ctx context.Context, userID, videoID string, comment *entity.Comment) error {
	return DB.WithContext(ctx).Create(comment).Error
}

func GetCommentsByVideoID(ctx context.Context, videoID string, pageNum, pageSize int32) ([]*entity.Comment, int64, error) {
	var comments []*entity.Comment
	var total int64
	db := DB.WithContext(ctx).Model(&entity.Comment{}).Where("video_id = ?", videoID)
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	offset := (pageNum - 1) * pageSize
	err := db.Offset(int(offset)).Limit(int(pageSize)).Order("created_at desc").Find(&comments).Error
	return comments, total, err
}

func DeleteComment(ctx context.Context, commentID, userID string) error {
	return DB.WithContext(ctx).Where("id = ? AND user_id = ?", commentID, userID).Delete(&entity.Comment{}).Error
}
