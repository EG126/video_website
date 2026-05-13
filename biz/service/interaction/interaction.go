package interaction

import (
	"context"
	"strings"
	"time"
	"video_website/biz/dal/mysql"
	"video_website/biz/dal/mysql/entity"
	"video_website/biz/dal/redis"
	"video_website/pkg/errno"
	"video_website/pkg/utils"

	interaction "video_website/biz/model/interaction"

	"github.com/cloudwego/hertz/pkg/common/hlog"
)

func LikeAction(ctx context.Context, userID, videoID, commentID, actionType string) error {
	actionTypeInt := int32(0)
	if actionType == "1" {
		actionTypeInt = 1
	} else if actionType == "2" {
		actionTypeInt = 2
	} else {
		hlog.CtxErrorf(ctx, "无效的action_type: %s", actionType)
		return errno.ParamError
	}

	if commentID != "" {
		// 对评论点赞
		hlog.CtxInfof(ctx, "评论点赞操作: userID=%s, commentID=%s, actionType=%d", userID, commentID, actionTypeInt)
		if err := mysql.CommentLikeAction(ctx, userID, commentID, actionTypeInt); err != nil {
			hlog.CtxErrorf(ctx, "评论点赞操作失败: %v", err)
			return errno.DBError
		}
		hlog.CtxInfof(ctx, "评论点赞操作成功: userID=%s, commentID=%s", userID, commentID)
	} else if videoID != "" {
		// 对视频点赞
		cleanVideoID := strings.Trim(videoID, "\"")
		hlog.CtxInfof(ctx, "视频点赞操作: userID=%s, videoID=%s, actionType=%d", userID, cleanVideoID, actionTypeInt)

		if err := mysql.LikeAction(ctx, userID, cleanVideoID, actionTypeInt); err != nil {
			hlog.CtxErrorf(ctx, "视频点赞操作失败: %v", err)
			return errno.DBError
		}

		if _, err := redis.RDB.Del(ctx, "popular:videos").Result(); err != nil {
			hlog.CtxErrorf(ctx, "删除热门视频缓存失败: %v", err)
		} else {
			hlog.CtxInfof(ctx, "热门视频缓存已清除")
		}

		hlog.CtxInfof(ctx, "视频点赞操作成功: userID=%s, videoID=%s", userID, cleanVideoID)
	} else {
		hlog.CtxErrorf(ctx, "必须提供 video_id 或 comment_id")
		return errno.ParamError
	}

	return nil
}

func LikeList(ctx context.Context, userID string, pageNum, pageSize int32) ([]*interaction.InteractionItemsResp, error) {
	hlog.CtxInfof(ctx, "点赞列表请求: userID=%s, page=%d, size=%d", userID, pageNum, pageSize)

	videoIDs, total, err := mysql.GetUserLikedVideoIDs(ctx, userID, pageNum, pageSize)
	if err != nil {
		hlog.CtxErrorf(ctx, "查询点赞列表失败: %v", err)
		return nil, errno.DBError
	}
	hlog.CtxInfof(ctx, "获取到点赞视频ID列表: %v, 总数: %d", videoIDs, total)

	if len(videoIDs) == 0 {
		hlog.CtxInfof(ctx, "点赞列表为空")
		return []*interaction.InteractionItemsResp{}, nil
	}

	cleanedIDs := make([]string, 0, len(videoIDs))
	for _, id := range videoIDs {
		cleanedIDs = append(cleanedIDs, strings.Trim(id, "\""))
	}
	hlog.CtxInfof(ctx, "清洗后的视频ID列表: %v", cleanedIDs)

	videos, err := mysql.GetVideosByIDs(ctx, cleanedIDs)
	if err != nil {
		hlog.CtxErrorf(ctx, "查询视频信息失败: %v", err)
		return nil, errno.DBError
	}
	hlog.CtxInfof(ctx, "查询到视频数量: %d", len(videos))

	videoMap := make(map[string]*entity.Video, len(videos))
	for _, v := range videos {
		videoMap[v.ID] = v
	}

	items := make([]*interaction.InteractionItemsResp, 0, len(videoIDs))
	for _, originalID := range videoIDs {
		cleanID := strings.Trim(originalID, "\"")
		v, ok := videoMap[cleanID]
		if !ok {
			hlog.CtxInfof(ctx, "视频ID %s 不存在，仅返回ID", cleanID)
			items = append(items, &interaction.InteractionItemsResp{
				ID: cleanID,
			})
			continue
		}
		items = append(items, &interaction.InteractionItemsResp{
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
	hlog.CtxInfof(ctx, "构造点赞列表响应, 实际返回数量: %d", len(items))
	return items, nil
}

func CommentPublish(ctx context.Context, userID, videoID, commentID, content string) error {
	if commentID != "" {
		// 对评论的评论（回复评论）
		hlog.CtxInfof(ctx, "回复评论: userID=%s, commentID=%s, content=%s", userID, commentID, content)

		comment := &entity.Comment{
			ID:        utils.GenerateID(),
			UserID:    userID,
			VideoID:   videoID,
			CommentID: commentID,
			Content:   content,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		if err := mysql.CreateCommentReply(ctx, comment); err != nil {
			hlog.CtxErrorf(ctx, "回复评论失败: %v", err)
			return errno.DBError
		}
		hlog.CtxInfof(ctx, "回复评论成功: commentID=%s", comment.ID)
	} else if videoID != "" {
		// 对视频的评论
		hlog.CtxInfof(ctx, "发表评论: userID=%s, videoID=%s, content=%s", userID, videoID, content)

		comment := &entity.Comment{
			ID:        utils.GenerateID(),
			UserID:    userID,
			VideoID:   videoID,
			Content:   content,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		if err := mysql.CreateComment(ctx, userID, videoID, comment); err != nil {
			hlog.CtxErrorf(ctx, "发表评论失败: %v", err)
			return errno.DBError
		}
		hlog.CtxInfof(ctx, "发表评论成功: commentID=%s", comment.ID)
	} else {
		hlog.CtxErrorf(ctx, "必须提供 video_id 或 comment_id")
		return errno.ParamError
	}

	return nil
}

func CommentList(ctx context.Context, videoID string, pageNum, pageSize int32) ([]*interaction.CommentItemResp, error) {
	hlog.CtxInfof(ctx, "评论列表请求: videoID=%s, page=%d, size=%d", videoID, pageNum, pageSize)

	comments, total, err := mysql.GetCommentsByVideoID(ctx, videoID, pageNum, pageSize)
	if err != nil {
		hlog.CtxErrorf(ctx, "查询评论列表失败: %v", err)
		return nil, errno.DBError
	}
	hlog.CtxInfof(ctx, "查询到评论数量: %d, 总数: %d", len(comments), total)

	items := make([]*interaction.CommentItemResp, 0, len(comments))
	for _, comment := range comments {
		items = append(items, &interaction.CommentItemResp{
			ID:        comment.ID,
			UserID:    comment.UserID,
			Content:   comment.Content,
			CreatedAt: comment.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: comment.UpdatedAt.Format("2006-01-02 15:04:05"),
			DeletedAt: "",
		})
	}
	hlog.CtxInfof(ctx, "构造评论列表响应, 实际返回数量: %d", len(items))
	return items, nil
}

func CommentDelete(ctx context.Context, commentID, userID string) error {
	hlog.CtxInfof(ctx, "删除评论: commentID=%s, userID=%s", commentID, userID)

	if err := mysql.DeleteComment(ctx, commentID, userID); err != nil {
		hlog.CtxErrorf(ctx, "删除评论失败: %v", err)
		return errno.DBError
	}
	hlog.CtxInfof(ctx, "删除评论成功: commentID=%s", commentID)
	return nil
}
