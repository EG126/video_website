package social

import (
	"context"
	"strings"
	"video_website/biz/dal/mysql"
	"video_website/pkg/errno"

	social "video_website/biz/model/social"

	"github.com/cloudwego/hertz/pkg/common/hlog"
)

func RelationAction(ctx context.Context, followerID, toUserID string, actionType int32) error {
	cleanToUserID := strings.Trim(toUserID, "\"“”")
	hlog.CtxInfof(ctx, "关注操作请求: 目标用户ID(清洗前)=%s, 清洗后=%s, 操作类型=%d", toUserID, cleanToUserID, actionType)
	hlog.CtxInfof(ctx, "发起用户ID: %s", followerID)

	targetUser, err := mysql.GetUserByID(ctx, cleanToUserID)
	if err != nil {
		hlog.CtxErrorf(ctx, "查询目标用户失败: %v", err)
		return errno.DBError
	}
	if targetUser == nil {
		return errno.UserNotExist
	}
	targetUserID := targetUser.ID

	if err := mysql.FollowAction(ctx, followerID, targetUserID, actionType); err != nil {
		hlog.CtxErrorf(ctx, "关注操作失败: %v", err)
		return errno.DBError
	}
	hlog.CtxInfof(ctx, "关注操作成功: 用户 %s %s 用户 %s", followerID, map[int32]string{1: "关注", 2: "取消关注"}[actionType], targetUserID)
	return nil
}

func FollowingList(ctx context.Context, userID string, pageNum, pageSize int32) ([]*social.SocialItemsResp, int32, error) {
	hlog.CtxInfof(ctx, "关注列表请求: user_id=%s, page=%d, size=%d", userID, pageNum, pageSize)

	followingIDs, total, err := mysql.GetFollowing(ctx, userID, pageNum, pageSize)
	if err != nil {
		hlog.CtxErrorf(ctx, "查询关注列表失败: %v", err)
		return nil, 0, errno.DBError
	}
	hlog.CtxInfof(ctx, "查询到关注ID数量: %d, 总数: %d", len(followingIDs), total)

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
	hlog.CtxInfof(ctx, "构造关注列表响应, 实际返回数量: %d", len(items))
	return items, int32(total), nil
}

func FollowerList(ctx context.Context, userID string, pageNum, pageSize int32) ([]*social.SocialItemsResp, int32, error) {
	hlog.CtxInfof(ctx, "粉丝列表请求: user_id=%s, page=%d, size=%d", userID, pageNum, pageSize)

	followerIDs, total, err := mysql.GetFollowers(ctx, userID, pageNum, pageSize)
	if err != nil {
		hlog.CtxErrorf(ctx, "查询粉丝列表失败: %v", err)
		return nil, 0, errno.DBError
	}
	hlog.CtxInfof(ctx, "查询到粉丝ID数量: %d, 总数: %d", len(followerIDs), total)

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
	hlog.CtxInfof(ctx, "构造粉丝列表响应, 实际返回数量: %d", len(items))
	return items, int32(total), nil
}

func FriendsList(ctx context.Context, userID string, pageNum, pageSize int32) ([]*social.SocialItemsResp, int32, error) {
	hlog.CtxInfof(ctx, "好友列表请求: userID=%s, page=%d, size=%d", userID, pageNum, pageSize)

	friendIDs, total, err := mysql.GetFriends(ctx, userID, pageNum, pageSize)
	if err != nil {
		hlog.CtxErrorf(ctx, "查询好友列表失败: %v", err)
		return nil, 0, errno.DBError
	}
	hlog.CtxInfof(ctx, "查询到好友ID数量: %d, 总数: %d", len(friendIDs), total)

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
	hlog.CtxInfof(ctx, "构造好友列表响应, 实际返回数量: %d", len(items))
	return items, int32(total), nil
}
