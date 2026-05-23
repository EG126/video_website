package video

import (
	"context"
	"io"
	"os"
	"path/filepath"
	"time"
	"video_website/biz/dal/mysql"
	"video_website/biz/dal/mysql/entity"
	"video_website/biz/dal/redis"
	"video_website/config"
	"video_website/pkg/constants"
	"video_website/pkg/errno"
	"video_website/pkg/utils"

	video "video_website/biz/model/video"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/cloudwego/hertz/pkg/common/json"
	pkgErrors "github.com/pkg/errors"
)

func Publish(ctx context.Context, c *app.RequestContext, userID, title, description string) error {
	file, err := c.FormFile("data")
	if err != nil {
		return errno.ParamError
	}

	src, err := file.Open()
	if err != nil {
		return errno.ParamError
	}
	defer src.Close()

	dataBytes, err := io.ReadAll(src)
	if err != nil {
		return pkgErrors.Wrap(err, "Publish: ReadFile failed")
	}

	if len(dataBytes) == 0 {
		return errno.ParamError
	}

	fileName := utils.GenerateID() + constants.VideoExt
	filePath := filepath.Join("./static/videos", fileName)
	if err := os.WriteFile(filePath, dataBytes, 0644); err != nil {
		return pkgErrors.Wrap(err, "Publish: WriteFile failed")
	}
	videoURL := config.Conf.Static.BaseURL + config.Conf.Static.VideoPath + fileName

	now := time.Now()
	newVideo := &entity.Video{
		ID:          utils.GenerateID(),
		UserID:      userID,
		Title:       title,
		Description: description,
		VideoURL:    videoURL,
		CoverURL:    "",
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	if err := mysql.CreateVideo(ctx, newVideo); err != nil {
		return pkgErrors.Wrap(err, "Publish: CreateVideo failed")
	}

	hlog.CtxInfof(ctx, "视频发布成功, userID=%s, videoID=%s", userID, newVideo.ID)
	return nil
}

func List(ctx context.Context, userID string, pageNum, pageSize int32) ([]*video.VideoItemsResp, int32, error) {
	videos, total, err := mysql.GetVideoByUserID(ctx, userID, pageNum, pageSize)
	if err != nil {
		return nil, 0, pkgErrors.Wrap(err, "List: GetVideoByUserID failed")
	}

	items := make([]*video.VideoItemsResp, 0, len(videos))
	for _, v := range videos {
		items = append(items, &video.VideoItemsResp{
			ID:           v.ID,
			UserID:       v.UserID,
			VideoURL:     v.VideoURL,
			CoverURL:     v.CoverURL,
			Title:        v.Title,
			Description:  v.Description,
			VisitCount:   int32(v.VisitCount),
			LikeCount:    int32(v.LikeCount),
			CommentCount: int32(v.CommentCount),
			CreatedAt:    v.CreatedAt.Format(constants.DateTimeFormat),
			UpdatedAt:    v.UpdatedAt.Format(constants.DateTimeFormat),
			DeletedAt:    "",
		})
	}

	return items, int32(total), nil
}

func Popular(ctx context.Context, pageNum, pageSize int32) ([]*video.VideoItemsResp, error) {
	cacheKey := constants.PopularVideosKey

	cached, err := redis.RDB.Get(ctx, cacheKey).Result()
	if err == nil {
		var allItems []*video.VideoItemsResp
		if json.Unmarshal([]byte(cached), &allItems) == nil {
			start := (pageNum - 1) * pageSize
			end := start + pageSize
			if int(start) < len(allItems) {
				if int(end) > len(allItems) {
					end = int32(len(allItems))
				}
				return allItems[start:end], nil
			}
			return []*video.VideoItemsResp{}, nil
		}
	}

	hlog.CtxInfof(ctx, "缓存未命中, 从数据库查询热门视频")

	videos, err := mysql.GetPopularVideos(ctx, 1, int32(config.Conf.Video.PopularLimit))
	if err != nil {
		return nil, pkgErrors.Wrap(err, "Popular: GetPopularVideos failed")
	}

	items := make([]*video.VideoItemsResp, 0, len(videos))
	for _, v := range videos {
		items = append(items, &video.VideoItemsResp{
			ID:           v.ID,
			UserID:       v.UserID,
			VideoURL:     v.VideoURL,
			CoverURL:     v.CoverURL,
			Title:        v.Title,
			Description:  v.Description,
			VisitCount:   int32(v.VisitCount),
			LikeCount:    int32(v.LikeCount),
			CommentCount: int32(v.CommentCount),
			CreatedAt:    v.CreatedAt.Format(constants.DateTimeFormat),
			UpdatedAt:    v.UpdatedAt.Format(constants.DateTimeFormat),
			DeletedAt:    "",
		})
	}

	if jsonData, err := json.Marshal(items); err == nil {
		redis.RDB.Set(ctx, cacheKey, jsonData, 5*time.Minute)
	}

	start := (pageNum - 1) * pageSize
	end := start + pageSize
	if int(start) < len(items) {
		if int(end) > len(items) {
			end = int32(len(items))
		}
		return items[start:end], nil
	}
	return []*video.VideoItemsResp{}, nil
}

func Search(ctx context.Context, keywords string, pageNum, pageSize int32) ([]*video.VideoItemsResp, int32, error) {
	videos, total, err := mysql.SearchVideos(ctx, keywords, pageNum, pageSize)
	if err != nil {
		return nil, 0, pkgErrors.Wrap(err, "Search: SearchVideos failed")
	}

	items := make([]*video.VideoItemsResp, 0, len(videos))
	for _, v := range videos {
		items = append(items, &video.VideoItemsResp{
			ID:           v.ID,
			UserID:       v.UserID,
			VideoURL:     v.VideoURL,
			CoverURL:     v.CoverURL,
			Title:        v.Title,
			Description:  v.Description,
			VisitCount:   int32(v.VisitCount),
			LikeCount:    int32(v.LikeCount),
			CommentCount: int32(v.CommentCount),
			CreatedAt:    v.CreatedAt.Format(constants.DateTimeFormat),
			UpdatedAt:    v.UpdatedAt.Format(constants.DateTimeFormat),
			DeletedAt:    "",
		})
	}

	return items, int32(total), nil
}

func Feed(ctx context.Context, latestTime string, userID string) ([]*video.VideoItemsResp, error) {
	videos, err := mysql.GetFeedVideos(ctx, latestTime)
	if err != nil {
		if err.Error() == "latest_time must be 13-digit timestamp" || err.Error() == "latest_time must be valid 13-digit timestamp" {
			return nil, errno.ParamError
		}
		return nil, pkgErrors.Wrap(err, "Feed: GetFeedVideos failed")
	}

	items := make([]*video.VideoItemsResp, 0, len(videos))
	for _, v := range videos {
		items = append(items, &video.VideoItemsResp{
			ID:           v.ID,
			UserID:       v.UserID,
			VideoURL:     v.VideoURL,
			CoverURL:     v.CoverURL,
			Title:        v.Title,
			Description:  v.Description,
			VisitCount:   int32(v.VisitCount),
			LikeCount:    int32(v.LikeCount),
			CommentCount: int32(v.CommentCount),
			CreatedAt:    v.CreatedAt.Format(constants.DateTimeFormat),
			UpdatedAt:    v.UpdatedAt.Format(constants.DateTimeFormat),
			DeletedAt:    "",
		})
	}

	return items, nil
}
