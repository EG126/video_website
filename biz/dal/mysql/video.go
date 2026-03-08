package mysql

import (
	"context"
	"video_website/biz/dal/mysql/entity"
)

func CreateVideo(ctx context.Context, video *entity.Video) error {
	return DB.WithContext(ctx).Create(video).Error
}

func GetVideoByUserID(ctx context.Context, userID string, pageNum, pageSize int32) ([]*entity.Video, int64, error) {
	var videos []*entity.Video
	var total int64
	offset := (pageNum - 1) * pageSize
	db := DB.WithContext(ctx).Model(&entity.Video{}).Where("user_id = ?", userID)
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	err := db.Offset(int(offset)).Limit(int(pageSize)).Order("created_at desc").Find(&videos).Error
	return videos, total, err
}

func GetPopularVideos(ctx context.Context, pageNum, pageSize int32) ([]*entity.Video, error) {
	offset := (pageNum - 1) * pageSize
	var videos []*entity.Video
	err := DB.WithContext(ctx).Order("like_count desc,visit_count desc,created_at desc").Offset(int(offset)).Limit(int(pageSize)).Find(&videos).Error
	return videos, err
}

func SerachVideos(ctx context.Context, keyword string, pageNum, pageSize int32) ([]*entity.Video, int64, error) {
	var videos []*entity.Video
	var total int64
	db := DB.WithContext(ctx).Model(&entity.Video{}).Where("title LIKE ? OR description LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	offset := (pageNum - 1) * pageSize
	err := db.Offset(int(offset)).Limit(int(pageSize)).Find(&videos).Error
	return videos, total, err
}
