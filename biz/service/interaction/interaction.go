package interaction

import (
	"context"
	"strings"
	"time"
	"video_website/biz/dal/mysql"
	"video_website/biz/dal/mysql/entity"
	"video_website/biz/dal/redis"
	"video_website/pkg/constants"
	"video_website/pkg/errno"
	"video_website/pkg/utils"

	interaction "video_website/biz/model/interaction"

	pkgErrors "github.com/pkg/errors"
)

func LikeAction(ctx context.Context, userID, videoID, commentID, actionType string) error {
	actionTypeInt := int32(0)
	switch actionType {
	case "1":
		actionTypeInt = constants.ActionLike
	case "2":
		actionTypeInt = constants.ActionUnlike
	default:
		return errno.ParamError
	}

	switch {
	case commentID != "":
		if err := mysql.CommentLikeAction(ctx, userID, commentID, actionTypeInt); err != nil {
			return pkgErrors.Wrap(err, "LikeAction: CommentLikeAction failed")
		}
	case videoID != "":
		cleanVideoID := strings.Trim(videoID, "\"")

		if err := mysql.LikeAction(ctx, userID, cleanVideoID, actionTypeInt); err != nil {
			return pkgErrors.Wrap(err, "LikeAction: LikeAction failed")
		}

		redis.RDB.Del(ctx, constants.PopularVideosKey)
	default:
		return errno.ParamError
	}

	return nil
}

func LikeList(ctx context.Context, userID string, pageNum, pageSize int32) ([]*interaction.InteractionItemsResp, error) {
	videoIDs, _, err := mysql.GetUserLikedVideoIDs(ctx, userID, pageNum, pageSize)
	if err != nil {
		return nil, pkgErrors.Wrap(err, "LikeList: GetUserLikedVideoIDs failed")
	}

	if len(videoIDs) == 0 {
		return []*interaction.InteractionItemsResp{}, nil
	}

	cleanedIDs := make([]string, 0, len(videoIDs))
	for _, id := range videoIDs {
		cleanedIDs = append(cleanedIDs, strings.Trim(id, "\""))
	}

	videos, err := mysql.GetVideosByIDs(ctx, cleanedIDs)
	if err != nil {
		return nil, pkgErrors.Wrap(err, "LikeList: GetVideosByIDs failed")
	}

	videoMap := make(map[string]*entity.Video, len(videos))
	for _, v := range videos {
		videoMap[v.ID] = v
	}

	items := make([]*interaction.InteractionItemsResp, 0, len(videoIDs))
	for _, originalID := range videoIDs {
		cleanID := strings.Trim(originalID, "\"")
		v, ok := videoMap[cleanID]
		if !ok {
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
			CreatedAt:    v.CreatedAt.Format(constants.DateTimeFormat),
			UpdatedAt:    v.UpdatedAt.Format(constants.DateTimeFormat),
			DeletedAt:    "",
		})
	}
	return items, nil
}

func CommentPublish(ctx context.Context, userID, videoID, commentID, content string) error {
	switch {
	case commentID != "":
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
			return pkgErrors.Wrap(err, "CommentPublish: CreateCommentReply failed")
		}
	case videoID != "":
		comment := &entity.Comment{
			ID:        utils.GenerateID(),
			UserID:    userID,
			VideoID:   videoID,
			Content:   content,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		if err := mysql.CreateComment(ctx, userID, videoID, comment); err != nil {
			return pkgErrors.Wrap(err, "CommentPublish: CreateComment failed")
		}
	default:
		return errno.ParamError
	}

	return nil
}

func CommentList(ctx context.Context, videoID string, pageNum, pageSize int32) ([]*interaction.CommentItemResp, error) {
	comments, _, err := mysql.GetCommentsByVideoID(ctx, videoID, pageNum, pageSize)
	if err != nil {
		return nil, pkgErrors.Wrap(err, "CommentList: GetCommentsByVideoID failed")
	}

	items := make([]*interaction.CommentItemResp, 0, len(comments))
	for _, comment := range comments {
		items = append(items, &interaction.CommentItemResp{
			ID:        comment.ID,
			UserID:    comment.UserID,
			Content:   comment.Content,
			CreatedAt: comment.CreatedAt.Format(constants.DateTimeFormat),
			UpdatedAt: comment.UpdatedAt.Format(constants.DateTimeFormat),
			DeletedAt: "",
		})
	}
	return items, nil
}

func CommentDelete(ctx context.Context, commentID, userID string) error {
	if err := mysql.DeleteComment(ctx, commentID, userID); err != nil {
		return pkgErrors.Wrap(err, "CommentDelete: DeleteComment failed")
	}
	return nil
}
