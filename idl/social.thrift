include "base.thrift"

namespace go social

struct SocialItemsResp{
    1:string id;
    2:string username;
    3:string avatar_url;
}

struct RelationActionReq{
    1:string access_token (api.header="Access-Token");
    2:string refresh_token (api.header="Refresh-Token");
    3:string to_user_id (api.form="to_user_id");
    4:i32 action_type (api.form="action_type");
}

struct RelationActionResp{
    1:base.BaseResp Base;
}

struct FollowingListReq{
    1:string user_id (api.query="user_id");
    2:i32 page_num (api.query="page_num");
    3:i32 page_size (api.query="page_size");
    4:string access_token (api.header="Access-Token");
    5:string refresh_token (api.header="Refresh-Token");
}

struct FollowingListResp{
    1:base.BaseResp Base;
    2:list<SocialItemsResp> Items;
    3:i32 total;
}

struct FollowerListReq{
    1:string user_id (api.query="user_id");
    2:i32 page_num (api.query="page_num");
    3:i32 page_size (api.query="page_size");
    4:string access_token (api.header="Access-Token");
    5:string refresh_token (api.header="Refresh-Token");
}

struct FollowerListResp{
    1:base.BaseResp Base;
    2:list<SocialItemsResp> Items;
    3:i32 total;
}

struct FriendsListReq{
    1:i32 page_num (api.query="page_num");
    2:i32 page_size (api.query="page_size");
    3:string access_token (api.header="Access-Token");
    4:string refresh_token (api.header="Refresh-Token");
}

struct FriendsListResp{
    1:base.BaseResp Base;
    2:list<SocialItemsResp> Items;
    3:i32 total;
}

service SocialService{
    RelationActionResp RelationAction(1:RelationActionReq req) (api.post="/relation/action");
    FollowingListResp FollowingList(1:FollowingListReq req) (api.get="/following/list");
    FollowerListResp FollowerList(1:FollowerListReq req) (api.get="/follower/list");
    FriendsListResp FriendsList(1:FriendsListReq req) (api.get="/friends/list");
}
