namespace go user

include "model.thrift"

// 用户登录
struct LoginReq {
  required string user_name;
  required string password;
}

struct LoginResp {
  required model.BaseResp baseResponse;
  optional model.User user;
}

// 用户注册
struct RegisterReq {
  required string user_name;
  required string password;
}

struct RegisterResp {
  required model.BaseResp baseResponse;
}

// 提交反馈
struct FeedbackReq {
  required i64 consult_id;
  required string content;
}

struct FeedbackResp {
  required model.BaseResp baseResponse;
}

service UserService {
  LoginResp login(1: LoginReq req) (api.post="/api/user/login"),
  RegisterResp register(1: RegisterReq req) (api.post="/api/user/register"),
  FeedbackResp feedback(1: FeedbackReq req) (api.post="/api/user/feedback")
}