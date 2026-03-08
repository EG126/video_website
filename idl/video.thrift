include "base.thrift"

namespace go video

namespace go items

struct VideoItemsResp{
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

struct VideoPublishReq{
    1:string access_token (api.header="Access-Token");
    2:string refresh_token (api.header="Refresh-Token");
    3:binary data (api.form="data");
    4:string title (api.form="title");
    5:string description (api.form="description");
}

struct VideoPublishResp{
    1:base.BaseResp Base;
}

struct VideoListReq{
    1:string user_id (api.query="user_id");
    2:i32 page_num (api.query="page_num");
    3:i32 page_size (api.query="page_size");
    4:string access_token (api.header="Access-Token");
    5:string refresh_token (api.header="Refresh-Token");
}

struct VideoListResp{
    1:base.BaseResp Base;
    2:list<VideoItemsResp> Items;
    3:i32 total;
}

struct VideoPopularReq{
    1:i32 page_num (api.query="page_num");
    2:i32 page_size (api.query="page_size");
    3:string access_token (api.header="Access-Token");
    4:string refresh_token (api.header="Refresh-Token");
}

struct VideoPopularResp{
    1:base.BaseResp Base;
    2:list<VideoItemsResp> Items;
}

struct VideoSearchReq{
    1:string access_token (api.header="Access-Token");
    2:string refresh_token (api.header="Refresh-Token");
    3:string keywords (api.form="keywords");
    4:i32 page_num (api.form="page_num");
    5:i32 page_size (api.form="page_size");
    6:i32 from_date (api.form="from_date");
    7:i32 to_date (api.form="to_date");
    8:string username (api.form="username");
}

struct VideoSearchResp{
    1:base.BaseResp Base;
    2:list<VideoItemsResp> Items;
    3:i32 total;
}

service VideoService{
    VideoPublishResp Publish(1:VideoPublishReq req) (api.post="/videos/publish");
    VideoListResp List(1:VideoListReq req) (api.get="/videos/list");
    VideoPopularResp Popular(1:VideoPopularReq req) (api.get="/videos/popular");
    VideoSearchResp Search(1:VideoSearchReq req) (api.post="/videos/search");
}
