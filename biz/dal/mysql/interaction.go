package mysql

import (
	"context"
	"time"
	"video_website/biz/dal/mysql/entity"
	"video_website/pkg/utils"

	pkgErrors "github.com/pkg/errors"
)

func LikeAction(ctx context.Context, userID, videoID string, actionType int32) error {
	switch actionType {
	case 1:
		like := &entity.Like{
			ID:        utils.GenerateID(),
			UserID:    userID,
			VideoID:   videoID,
			CreatedAt: time.Now(),
		}
		err := DB.WithContext(ctx).FirstOrCreate(like, "user_id = ? AND video_id = ?", userID, videoID).Error
		if err != nil {
			return pkgErrors.Wrapf(err, "LikeAction create failed, userID=%s, videoID=%s", userID, videoID)
		}
		return nil
	case 2:
		err := DB.WithContext(ctx).Where("user_id = ? AND video_id = ?", userID, videoID).Delete(&entity.Like{}).Error
		if err != nil {
			return pkgErrors.Wrapf(err, "LikeAction delete failed, userID=%s, videoID=%s", userID, videoID)
		}
		return nil
	}
	return nil
}

func CommentLikeAction(ctx context.Context, userID, commentID string, actionType int32) error {
	switch actionType {
	case 1:
		like := &entity.CommentLike{
			ID:        utils.GenerateID(),
			UserID:    userID,
			CommentID: commentID,
			CreatedAt: time.Now(),
		}
		err := DB.WithContext(ctx).FirstOrCreate(like, "user_id = ? AND comment_id = ?", userID, commentID).Error
		if err != nil {
			return pkgErrors.Wrapf(err, "CommentLikeAction create failed, userID=%s, commentID=%s", userID, commentID)
		}
		return nil
	case 2:
		err := DB.WithContext(ctx).Where("user_id = ? AND comment_id = ?", userID, commentID).Delete(&entity.CommentLike{}).Error
		if err != nil {
			return pkgErrors.Wrapf(err, "CommentLikeAction delete failed, userID=%s, commentID=%s", userID, commentID)
		}
		return nil
	}
	return nil
}

func GetUserLikedVideoIDs(ctx context.Context, userID string, pageNum, pageSize int32) ([]string, int64, error) {
	var likes []*entity.Like
	var total int64
	db := DB.WithContext(ctx).Model(&entity.Like{}).Where("user_id = ?", userID)
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, pkgErrors.Wrapf(err, "GetUserLikedVideoIDs count failed, userID=%s", userID)
	}
	offset := (pageNum - 1) * pageSize
	err := db.Offset(int(offset)).Limit(int(pageSize)).Order("created_at desc").Find(&likes).Error
	if err != nil {
		return nil, 0, pkgErrors.Wrapf(err, "GetUserLikedVideoIDs query failed, userID=%s", userID)
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
	if err != nil {
		return nil, pkgErrors.Wrapf(err, "GetVideosByIDs query failed")
	}
	return videos, nil
}

func CreateComment(ctx context.Context, _, _ string, comment *entity.Comment) error {
	err := DB.WithContext(ctx).Create(comment).Error
	if err != nil {
		return pkgErrors.Wrapf(err, "CreateComment failed, commentID=%s", comment.ID)
	}
	return nil
}

func CreateCommentReply(ctx context.Context, comment *entity.Comment) error {
	err := DB.WithContext(ctx).Create(comment).Error
	if err != nil {
		return pkgErrors.Wrapf(err, "CreateCommentReply failed, commentID=%s", comment.ID)
	}
	return nil
}

func GetCommentsByVideoID(ctx context.Context, videoID string, pageNum, pageSize int32) ([]*entity.Comment, int64, error) {
	var comments []*entity.Comment
	var total int64
	db := DB.WithContext(ctx).Model(&entity.Comment{}).Where("video_id = ?", videoID)
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, pkgErrors.Wrapf(err, "GetCommentsByVideoID count failed, videoID=%s", videoID)
	}
	offset := (pageNum - 1) * pageSize
	err := db.Offset(int(offset)).Limit(int(pageSize)).Order("created_at desc").Find(&comments).Error
	if err != nil {
		return nil, 0, pkgErrors.Wrapf(err, "GetCommentsByVideoID query failed, videoID=%s", videoID)
	}
	return comments, total, nil
}

func DeleteComment(ctx context.Context, commentID, userID string) error {
	err := DB.WithContext(ctx).Where("id = ? AND user_id = ?", commentID, userID).Delete(&entity.Comment{}).Error
	if err != nil {
		return pkgErrors.Wrapf(err, "DeleteComment failed, commentID=%s, userID=%s", commentID, userID)
	}
	return nil
}
