package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	RedisChannelPrivate = "chat:private:%s"
	RedisChannelRoom    = "chat:room:%s"
	RedisKeyOnlineUsers = "chat:online_users"
	RedisKeyUnreadCount = "chat:unread:%s:%s"
)

func SetUserOnline(userID string) error {
	ctx := context.Background()
	pipe := RDB.Pipeline()
	pipe.SAdd(ctx, RedisKeyOnlineUsers, userID)
	pipe.HSet(ctx, fmt.Sprintf("chat:user:%s:online", userID), "online", "1", "last_seen", time.Now().Unix())
	_, err := pipe.Exec(ctx)
	return err
}

func SetUserOffline(userID string) error {
	ctx := context.Background()
	pipe := RDB.Pipeline()
	pipe.SRem(ctx, RedisKeyOnlineUsers, userID)
	pipe.HSet(ctx, fmt.Sprintf("chat:user:%s:online", userID), "online", "0", "last_seen", time.Now().Unix())
	_, err := pipe.Exec(ctx)
	return err
}

func IsUserOnline(userID string) bool {
	return RDB.SIsMember(context.Background(), RedisKeyOnlineUsers, userID).Val()
}

func IncrementUnreadCount(toUserID, fromUserID string) error {
	ctx := context.Background()
	key := fmt.Sprintf(RedisKeyUnreadCount, toUserID, fromUserID)
	pipe := RDB.Pipeline()
	pipe.Incr(ctx, key)
	pipe.Expire(ctx, key, 7*24*time.Hour)
	_, err := pipe.Exec(ctx)
	return err
}

func GetUnreadCount(toUserID, fromUserID string) int64 {
	key := fmt.Sprintf(RedisKeyUnreadCount, toUserID, fromUserID)
	count, err := RDB.Get(context.Background(), key).Int64()
	if err != nil {
		return 0
	}
	return count
}

func ClearUnreadCount(toUserID, fromUserID string) error {
	key := fmt.Sprintf(RedisKeyUnreadCount, toUserID, fromUserID)
	return RDB.Del(context.Background(), key).Err()
}

func PublishPrivateMessage(toUserID string, data []byte) error {
	channel := fmt.Sprintf(RedisChannelPrivate, toUserID)
	return RDB.Publish(context.Background(), channel, data).Err()
}

func PublishRoomMessage(roomID string, data []byte) error {
	channel := fmt.Sprintf(RedisChannelRoom, roomID)
	return RDB.Publish(context.Background(), channel, data).Err()
}

func SubscribePrivateMessages(ctx context.Context) *redis.PubSub {
	return RDB.PSubscribe(ctx, "chat:private:*")
}

func SubscribeRoomMessages(ctx context.Context) *redis.PubSub {
	return RDB.PSubscribe(ctx, "chat:room:*")
}
