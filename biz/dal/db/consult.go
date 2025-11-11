package db

import (
    "buycar/pkg/errno"
    "context"
    "errors"

    "gorm.io/gorm"
)

// CreateConsult 创建咨询记录
func CreateConsult(ctx context.Context, consult *Consult) error {
    if err := DB.WithContext(ctx).Create(consult).Error; err != nil {
        return errno.NewErrNo(errno.InternalDatabaseErrorCode, "创建咨询失败: "+err.Error())
    }
    return nil
}

// GetConsultByID 通过 consult_id 查询咨询记录
func GetConsultByID(ctx context.Context, id int64) (*Consult, error) {
    var consult Consult
    err := DB.WithContext(ctx).Where("consult_id = ?", id).First(&consult).Error
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, errno.NewErrNo(errno.ServiceConsultNotExist, "咨询记录不存在")
        }
        return nil, errno.NewErrNo(errno.InternalDatabaseErrorCode, "查询咨询失败: "+err.Error())
    }
    return &consult, nil
}

// UpdateConsultStatus 更新咨询的状态
func UpdateConsultStatus(ctx context.Context, id int64, status string) error {
    if err := DB.WithContext(ctx).Model(&Consult{}).Where("consult_id = ?", id).Updates(map[string]interface{}{
        "status": status,
    }).Error; err != nil {
        return errno.NewErrNo(errno.InternalDatabaseErrorCode, "更新咨询状态失败: "+err.Error())
    }
    return nil
}

// UpdateConsultLLMFields 更新 LLM 相关字段（可传 nil 表示不修改）
func UpdateConsultLLMFields(ctx context.Context, id int64, llmModel, llmPrompt, llmResponse, recommendations *string) error {
    updates := map[string]interface{}{}
    if llmModel != nil {
        updates["llm_model"] = *llmModel
    }
    if llmPrompt != nil {
        updates["llm_prompt"] = *llmPrompt
    }
    if llmResponse != nil {
        updates["llm_response"] = *llmResponse
    }
    if recommendations != nil {
        updates["recommendations"] = *recommendations
    }
    if len(updates) == 0 {
        return nil
    }
    if err := DB.WithContext(ctx).Model(&Consult{}).Where("consult_id = ?", id).Updates(updates).Error; err != nil {
        return errno.NewErrNo(errno.InternalDatabaseErrorCode, "更新咨询字段失败: "+err.Error())
    }
    return nil
}