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
	"video_website/pkg/errno"
	"video_website/pkg/utils"

	video "video_website/biz/model/video"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/cloudwego/hertz/pkg/common/json"
)

func Publish(ctx context.Context, c *app.RequestContext, userID, title, description string) error {
	hlog.CtxInfof(ctx, "开始处理投稿请求, 用户ID: %s", userID)

	file, err := c.FormFile("data")
	if err != nil {
		hlog.CtxErrorf(ctx, "获取文件表单失败: %v", err)
		return errno.ParamError
	}

	src, err := file.Open()
	if err != nil {
		hlog.CtxErrorf(ctx, "打开文件失败: %v", err)
		return errno.ParamError
	}
	defer src.Close()

	dataBytes, err := io.ReadAll(src)
	if err != nil {
		hlog.CtxErrorf(ctx, "读取文件失败: %v", err)
		return errno.ParamError
	}
	hlog.CtxInfof(ctx, "获取投稿文件成功, 文件名: %s, 大小: %d bytes", file.Filename, file.Size)

	if len(dataBytes) == 0 {
		return errno.ParamError
	}

	fileName := utils.GenerateID() + ".mp4"
	filePath := filepath.Join("./static/videos", fileName)
	if err := os.WriteFile(filePath, dataBytes, 0644); err != nil {
		hlog.CtxErrorf(ctx, "保存视频文件失败: %v", err)
		return errno.DBError
	}
	videoURL := "http://127.0.0.1:8888/static/videos/" + fileName
	hlog.CtxInfof(ctx, "视频文件保存成功, 路径: %s", filePath)

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
		hlog.CtxErrorf(ctx, "创建视频记录失败: %v", err)
		return errno.DBError
	}
	hlog.CtxInfof(ctx, "视频记录创建成功, 视频ID: %s", newVideo.ID)

	hlog.CtxInfof(ctx, "投稿处理完成, 用户ID: %s", userID)
	return nil
}

func List(ctx context.Context, userID string, pageNum, pageSize int32) ([]*video.VideoItemsResp, int32, error) {
	hlog.CtxInfof(ctx, "开始获取发布列表, 用户ID: %s, 页码: %d, 每页数量: %d", userID, pageNum, pageSize)

	videos, total, err := mysql.GetVideoByUserID(ctx, userID, pageNum, pageSize)
	if err != nil {
		hlog.CtxErrorf(ctx, "查询视频列表失败: %v", err)
		return nil, 0, errno.DBError
	}
	hlog.CtxInfof(ctx, "查询到视频数量: %d, 总数: %d", len(videos), total)

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
			CreatedAt:    v.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:    v.UpdatedAt.Format("2006-01-02 15:04:05"),
			DeletedAt:    "",
		})
	}

	hlog.CtxInfof(ctx, "发布列表返回成功, 用户ID: %s", userID)
	return items, int32(total), nil
}

func Popular(ctx context.Context, pageNum, pageSize int32) ([]*video.VideoItemsResp, error) {
	hlog.CtxInfof(ctx, "开始获取热门视频, 页码: %d, 每页数量: %d", pageNum, pageSize)

	cacheKey := "popular:videos"

	cached, err := redis.RDB.Get(ctx, cacheKey).Result()
	if err == nil {
		var allItems []*video.VideoItemsResp
		if json.Unmarshal([]byte(cached), &allItems) == nil {
			hlog.CtxInfof(ctx, "命中缓存, 热门视频总数: %d", len(allItems))
			start := (pageNum - 1) * pageSize
			end := start + pageSize
			if int(start) < len(allItems) {
				if int(end) > len(allItems) {
					end = int32(len(allItems))
				}
				hlog.CtxInfof(ctx, "从缓存返回热门视频, 数量: %d", end-start)
				return allItems[start:end], nil
			}
			hlog.CtxInfof(ctx, "缓存数据超出范围, 返回空列表")
			return []*video.VideoItemsResp{}, nil
		}
	}

	hlog.CtxInfof(ctx, "缓存未命中, 从数据库查询热门视频")
	videos, err := mysql.GetPopularVideos(ctx, 1, 100)
	if err != nil {
		hlog.CtxErrorf(ctx, "查询热门视频失败: %v", err)
		return nil, errno.DBError
	}
	hlog.CtxInfof(ctx, "从数据库查询到热门视频数量: %d", len(videos))

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
			CreatedAt:    v.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:    v.UpdatedAt.Format("2006-01-02 15:04:05"),
			DeletedAt:    "",
		})
	}

	if jsonData, err := json.Marshal(items); err == nil {
		redis.RDB.Set(ctx, cacheKey, jsonData, 5*time.Minute)
		hlog.CtxInfof(ctx, "热门视频缓存更新成功, 数量: %d", len(items))
	}

	start := (pageNum - 1) * pageSize
	end := start + pageSize
	if int(start) < len(items) {
		if int(end) > len(items) {
			end = int32(len(items))
		}
		hlog.CtxInfof(ctx, "返回热门视频, 数量: %d", end-start)
		return items[start:end], nil
	}
	hlog.CtxInfof(ctx, "热门视频列表为空")
	return []*video.VideoItemsResp{}, nil
}

func Search(ctx context.Context, keywords string, pageNum, pageSize int32) ([]*video.VideoItemsResp, int32, error) {
	hlog.CtxInfof(ctx, "开始搜索视频, 关键词: %s, 页码: %d, 每页数量: %d", keywords, pageNum, pageSize)

	videos, total, err := mysql.SearchVideos(ctx, keywords, pageNum, pageSize)
	if err != nil {
		hlog.CtxErrorf(ctx, "搜索视频失败: %v", err)
		return nil, 0, errno.DBError
	}
	hlog.CtxInfof(ctx, "搜索到视频数量: %d, 总数: %d", len(videos), total)

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
			CreatedAt:    v.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:    v.UpdatedAt.Format("2006-01-02 15:04:05"),
			DeletedAt:    "",
		})
	}

	hlog.CtxInfof(ctx, "搜索视频返回成功, 关键词: %s", keywords)
	return items, int32(total), nil
}
