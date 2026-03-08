include "base.thrift"

namespace go interaction

struct InteractionItemsResp{
    1:string id;
    2:string user_id;
    3:string video_url;
    4:string cover_url;
    5:string title;
    6:string description;
    7:i32 visit_count;
    8:i32 like_count;
    9:i32 comment_count;
    10:string created_at;
    11:string updated_at;
    12:string deleted_at;
}

struct LikeActionReq{
    1:string access_token (api.header="Access-Token");
    2:string refresh_token (api.header="Refresh-Token");
    3:string video_id (api.form="video_id");
    4:string comment_id (api.form="comment_id");
    5:string action_type (api.form="action_type");
}

struct LikeActionResp{
    1:base.BaseResp Base;
}

struct LikeListReq{
    1:string user_id (api.query="user_id");
    2:i32 page_size (api.query="page_size");
    3:i32 page_num (api.query="page_num");
    4:string access_token (api.header="Access-Token");
    5:string refresh_token (api.header="Refresh-Token");
}

struct LikeListResp{
    1:base.BaseResp Base;
    2:list<InteractionItemsResp> Items;
}

struct CommentPublishReq{

    3:string video_id (api.form="video_id");
    4:string comment_id (api.form="comment_id");
    5:string content (api.form="content");
}

struct CommentPublishResp{
    1:base.BaseResp Base;
}

struct CommentListReq{
    1:string video_id (api.query="video_id");
    2:string comment_id (api.query="comment_id");
    3:i32 page_size (api.query="page_size");
    4:i32 page_num (api.query="page_num");
    5:string access_token (api.header="Access-Token");
    6:string refresh_token (api.header="Refresh-Token");
}

struct CommentListResp{
    1:base.BaseResp Base;
    2:list<InteractionItemsResp> Items;
}

struct CommentDeleteReq{
    1:string access_token (api.header="Access-Token");
    2:string refresh_token (api.header="Refresh-Token");
    3:string video_id (api.form="video_id");
    4:string comment_id (api.form="comment_id");
}

struct CommentDeleteResp{
    1:base.BaseResp Base;
}

service InteractionService{
    LikeActionResp LikeAction(1:LikeActionReq req) (api.post="/like/action");
    LikeListResp LikeList(1:LikeListReq req) (api.post="/like/list");
    CommentPublishResp CommentPublish(1:CommentPublishReq req) (api.post="/comment/publish");
    CommentListResp CommentList(1:CommentListReq req) (api.post="/comment/list");
    CommentDeleteResp CommentDelete(1:CommentDeleteReq req) (api.post="/comment/delete");
}
