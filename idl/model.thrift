namespace go module

struct BaseResp {
    required i32 code,
    required string message,
}

struct User {
    required i64 user_id,
    required string username,
    required string name,
    required string telephone,
    required string email,
    optional string password_hash,
    required bool is_admin,
    required i64 score,
    required string created_at,
    required string updated_at,
}

struct Consult {
    required string consult_id,
    optional string user_id,
    optional string budget_range,
    optional string preferred_type,
    optional string use_case,
    optional string fuel_type,
    optional string brand_preference,
    optional string llm_model,
    optional string llm_prompt,
    optional string llm_response,
    optional string recommendations,
    required i64 created_at,
    required i64 updated_at,
}

struct Feedback {
    required i64 id,
    required i64 user_id,
    required i64 consult_id,
    required string content,
    required string created_at,
}

struct Gift {
    required i64 gift_id,
    required string name,
    required string description,
    required i64 score_cost,
    required i64 stock,
    required string status,
    required string created_at,
    required string updated_at,
}

struct GiftRedemption {
    required i64 id,
    required i64 gift_id,
    required string user_id,
    required i32 quantity,
    required i32 points_spent,
    required string status,
    required i64 created_at,
}

struct ScoreTransaction {
    required i64 id,
    required i64 user_id,
    required i64 change_amount,
    required string reason,
    required i64 ref_id,
    required string description,
    required string created_at,
}