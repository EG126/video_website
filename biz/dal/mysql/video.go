package mysql

import (
	"context"
	stdErrors "errors"
	"strconv"
	"time"
	"video_website/biz/dal/mysql/entity"
	"video_website/config"

	pkgErrors "github.com/pkg/errors"
	"gorm.io/gorm"
)

func CreateVideo(ctx context.Context, video *entity.Video) error {
	err := DB.WithContext(ctx).Create(video).Error
	if err != nil {
		return pkgErrors.Wrapf(err, "CreateVideo failed, videoID=%s", video.ID)
	}
	return nil
}

func GetVideoByUserID(ctx context.Context, userID string, pageNum, pageSize int32) ([]*entity.Video, int64, error) {
	var videos []*entity.Video
	var total int64
	offset := (pageNum - 1) * pageSize
	db := DB.WithContext(ctx).Model(&entity.Video{}).Where("user_id = ?", userID)
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, pkgErrors.Wrapf(err, "GetVideoByUserID count failed, userID=%s", userID)
	}
	err := db.Offset(int(offset)).Limit(int(pageSize)).Order("created_at desc").Find(&videos).Error
	if err != nil {
		return nil, 0, pkgErrors.Wrapf(err, "GetVideoByUserID query failed, userID=%s", userID)
	}
	return videos, total, nil
}

func GetPopularVideos(ctx context.Context, pageNum, pageSize int32) ([]*entity.Video, error) {
	offset := (pageNum - 1) * pageSize
	var videos []*entity.Video
	err := DB.WithContext(ctx).Order("like_count desc,visit_count desc,created_at desc").Offset(int(offset)).Limit(int(pageSize)).Find(&videos).Error
	if err != nil {
		return nil, pkgErrors.Wrapf(err, "GetPopularVideos query failed")
	}
	return videos, nil
}

func GetVideoByID(ctx context.Context, id string) (*entity.Video, error) {
	var video entity.Video
	err := DB.WithContext(ctx).Where("id = ?", id).First(&video).Error
	if err != nil {
		if stdErrors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, pkgErrors.Wrapf(err, "GetVideoByID failed, id=%s", id)
	}
	return &video, nil
}

func SearchVideos(ctx context.Context, keyword string, pageNum, pageSize int32) ([]*entity.Video, int64, error) {
	var videos []*entity.Video
	var total int64
	db := DB.WithContext(ctx).Model(&entity.Video{}).Where("title LIKE ? OR description LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, pkgErrors.Wrapf(err, "SearchVideos count failed, keyword=%s", keyword)
	}
	offset := (pageNum - 1) * pageSize
	err := db.Offset(int(offset)).Limit(int(pageSize)).Find(&videos).Error
	if err != nil {
		return nil, 0, pkgErrors.Wrapf(err, "SearchVideos query failed, keyword=%s", keyword)
	}
	return videos, total, nil
}

func GetFeedVideos(ctx context.Context, latestTime string) ([]*entity.Video, error) {
	var videos []*entity.Video
	feedLimit := int32(config.Conf.Video.FeedLimit)
	query := DB.WithContext(ctx).Model(&entity.Video{}).Order("created_at desc").Limit(int(feedLimit))
	if latestTime != "" {
		if len(latestTime) != 13 {
			return nil, pkgErrors.New("latest_time must be 13-digit timestamp")
		}
		timestamp, err := strconv.ParseInt(latestTime, 10, 64)
		if err != nil {
			return nil, pkgErrors.Wrapf(err, "latest_time must be valid 13-digit timestamp")
		}
		parsedTime := time.UnixMilli(timestamp)
		query = query.Where("created_at > ?", parsedTime)
	}
	err := query.Find(&videos).Error
	if err != nil {
		return nil, pkgErrors.Wrapf(err, "GetFeedVideos query failed")
	}
	return videos, nil
}
