include "base.thrift"

namespace go user

struct UserDataResp{
    1:string id;
    2:string username;
    3:string avatar_url;
    4:string created_at;
    5:string updated_at;
    6:string deleted_at;
}

struct UserRegisterReq{
    1:string username (api.form="username");
    2:string password (api.form="password");
}

struct UserRegisterResp{
    1:base.BaseResp Base;
}

struct UserLoginReq{
    1:string username (api.form="username");
    2:string password (api.form="password");
    3:string code (api.form="code");
}

struct UserLoginResp{
    1:base.BaseResp Base;
    2:list<UserDataResp> Data;
    3:string access_token;
    4:string refresh_token;
}

struct RefreshReq{
    1:string refresh_token (api.form="refresh_token");
}

struct RefreshResp{
    1:base.BaseResp Base;
    2:string access_token;
    3:string refresh_token;
}

struct UserInfoReq{
    1:string user_id (api.query="user_id");
    2:string access_token (api.header="Access-Token");
    3:string refresh_token (api.header="Refresh-Token");
}

struct UserInfoResp{
    1:base.BaseResp Base;
    2:list<UserDataResp> Data;
}

struct UserAvatarReq{
    1:string access_token (api.header="Access-Token");
    2:string refresh_token (api.header="Refresh-Token");
    3:binary data (api.form="data");
}

struct UserAvatarResp{
    1:base.BaseResp Base;
    2:list<UserDataResp> Data;
}

service UserService{
    UserRegisterResp Register(1:UserRegisterReq req) (api.post="/user/register");
    UserLoginResp Login(1:UserLoginReq req) (api.post="/user/login");
    UserInfoResp Info(1:UserInfoReq req) (api.get="/user/info");
    UserAvatarResp UploadAvatar(1:UserAvatarReq req) (api.put="/user/avatar/upload");
    RefreshResp Refresh(1:RefreshReq req)  (api.post="/user/refresh");
}



