package social

import (
	"context"
	"strings"
	"video_website/biz/dal/mysql"
	"video_website/pkg/constants"
	"video_website/pkg/errno"

	social "video_website/biz/model/social"

	"github.com/cloudwego/hertz/pkg/common/hlog"
	pkgErrors "github.com/pkg/errors"
)

func RelationAction(ctx context.Context, followerID, toUserID string, actionType int32) error {
	cleanToUserID := strings.Trim(toUserID, "\"“”")

	targetUser, err := mysql.GetUserByID(ctx, cleanToUserID)
	if err != nil {
		return pkgErrors.Wrap(err, "RelationAction: GetUserByID failed")
	}
	if targetUser == nil {
		return errno.UserNotExist
	}
	targetUserID := targetUser.ID

	if err := mysql.FollowAction(ctx, followerID, targetUserID, actionType); err != nil {
		return pkgErrors.Wrap(err, "RelationAction: FollowAction failed")
	}

	actionName := map[int32]string{constants.ActionLike: "关注", constants.ActionUnlike: "取消关注"}[actionType]
	hlog.CtxInfof(ctx, "关注操作成功: userID=%s %s targetID=%s", followerID, actionName, targetUserID)
	return nil
}

func FollowingList(ctx context.Context, userID string, pageNum, pageSize int32) ([]*social.SocialItemsResp, int32, error) {
	followingIDs, total, err := mysql.GetFollowing(ctx, userID, pageNum, pageSize)
	if err != nil {
		return nil, 0, pkgErrors.Wrap(err, "FollowingList: GetFollowing failed")
	}

	items := make([]*social.SocialItemsResp, 0, len(followingIDs))
	for _, id := range followingIDs {
		u, err := mysql.GetUserByID(ctx, id)
		if err != nil || u == nil {
			continue
		}
		items = append(items, &social.SocialItemsResp{
			ID:        u.ID,
			Username:  u.Username,
			AvatarURL: u.AvatarURL,
		})
	}
	return items, int32(total), nil
}

func FollowerList(ctx context.Context, userID string, pageNum, pageSize int32) ([]*social.SocialItemsResp, int32, error) {
	followerIDs, total, err := mysql.GetFollowers(ctx, userID, pageNum, pageSize)
	if err != nil {
		return nil, 0, pkgErrors.Wrap(err, "FollowerList: GetFollowers failed")
	}

	items := make([]*social.SocialItemsResp, 0, len(followerIDs))
	for _, id := range followerIDs {
		u, err := mysql.GetUserByID(ctx, id)
		if err != nil || u == nil {
			continue
		}
		items = append(items, &social.SocialItemsResp{
			ID:        u.ID,
			Username:  u.Username,
			AvatarURL: u.AvatarURL,
		})
	}
	return items, int32(total), nil
}

func FriendsList(ctx context.Context, userID string, pageNum, pageSize int32) ([]*social.SocialItemsResp, int32, error) {
	friendIDs, total, err := mysql.GetFriends(ctx, userID, pageNum, pageSize)
	if err != nil {
		return nil, 0, pkgErrors.Wrap(err, "FriendsList: GetFriends failed")
	}

	items := make([]*social.SocialItemsResp, 0, len(friendIDs))
	for _, id := range friendIDs {
		u, err := mysql.GetUserByID(ctx, id)
		if err != nil || u == nil {
			continue
		}
		items = append(items, &social.SocialItemsResp{
			ID:        u.ID,
			Username:  u.Username,
			AvatarURL: u.AvatarURL,
		})
	}
	return items, int32(total), nil
}
