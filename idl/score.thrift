namespace go score

include "model.thrift"

// 查看用户积分
struct GetUserScoreReq {
}

struct GetUserScoreResp {
  required model.BaseResp baseResponse;
  required i64 score;
}

// 兑换汽车周边
struct PurchaseGiftReq {
  required i64 gift_id;
}

struct PurchaseGiftResp {
  required model.BaseResp baseResponse;
}

// 查看汽车周边列表
struct QueryGiftListReq {
}

struct QueryGiftListResp {
  required model.BaseResp baseResponse;
  required list<model.Gift> gifts;
}

service ScoreService {
  GetUserScoreResp getUserScore(1: GetUserScoreReq req) (api.get="/api/score"),
  PurchaseGiftResp purchaseGift(1: PurchaseGiftReq req) (api.post="/api/score/gift/purchase"),
  QueryGiftListResp queryGiftList(1: QueryGiftListReq req) (api.get="/api/score/gift/query")
}