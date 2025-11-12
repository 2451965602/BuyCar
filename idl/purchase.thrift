namespace go purchase

include "model.thrift"

// 购车咨询请求
struct PurchaseConsultReq {
  optional string budget_range;
  optional string preferred_type;
  optional string use_case;
  optional string fuel_type;
  optional string brand_preference;
}

struct PurchaseConsultResp {
  required model.BaseResp baseResponse;
  optional model.Consult consult;
}

// 查询咨询记录请求
struct QueryConsultReq {
  optional string consult_id;
}

struct QueryConsultResp {
  required model.BaseResp baseResponse;
  optional model.Consult consult;
}

service ConsultService {
  PurchaseConsultResp purchaseConsult(1: PurchaseConsultReq req) (api.post="/api/consult/purchase"),
  QueryConsultResp queryConsult(1: QueryConsultReq req) (api.get="/api/consult/query")
}