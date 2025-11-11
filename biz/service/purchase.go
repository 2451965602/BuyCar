package service

import (
    "buycar/biz/dal/db"
    "buycar/biz/model/module"
    purchase "buycar/biz/model/purchase"
    "buycar/pkg/AIAgent"
    "buycar/pkg/constants"
    "buycar/pkg/errno"
    "context"
    "strconv"
    "time"

    "github.com/cloudwego/hertz/pkg/app"
)

type ConsultService struct {
    ctx context.Context
    c   *app.RequestContext
}

func NewConsultService(ctx context.Context, c *app.RequestContext) *ConsultService {
    return &ConsultService{ctx: ctx, c: c}
}

// PurchaseConsult 创建购车咨询记录
func (s *ConsultService) PurchaseConsult(req *purchase.PurchaseConsultReq) (*module.Consult, error) {
    // 尝试从上下文读取用户ID（可选）
    var userIDPtr *int64
    if uidVal, ok := s.c.Get(constants.ContextUid); ok && uidVal != nil {
        if uid, err := convertToInt64(uidVal); err == nil {
            userIDPtr = &uid
        }
    }

    consult := &db.Consult{
        ConsultId:       time.Now().UnixNano(),
        UserId:          userIDPtr,
        BudgetRange:     req.BudgetRange,
        PreferredType:   req.PreferredType,
        UseCase:         req.UseCase,
        FuelType:        req.FuelType,
        BrandPreference: req.BrandPreference,
        Status:          "created",
    }

    if err := db.CreateConsult(s.ctx, consult); err != nil {
        return nil, err
    }

    // 异步触发 LLM 生成与更新，不阻塞请求返回
    go s.generateAndUpdate(*consult)

    return consult.ToModuleStruct(), nil
}

// QueryConsult 根据 consult_id 查询咨询详情
func (s *ConsultService) QueryConsult(req *purchase.QueryConsultReq) (*module.Consult, error) {
    if req.ConsultID == nil || *req.ConsultID == "" {
        return nil, errno.ParamMissingError.WithMessage("consult_id 不能为空")
    }

    id, err := strconv.ParseInt(*req.ConsultID, 10, 64)
    if err != nil {
        return nil, errno.ParamVerifyError.WithMessage("consult_id 格式不正确")
    }

    consult, err := db.GetConsultByID(s.ctx, id)
    if err != nil {
        return nil, err
    }

    return consult.ToModuleStruct(), nil
}

// generateAndUpdate 在后台调用 LLM 生成，并更新数据库记录
func (s *ConsultService) generateAndUpdate(consult db.Consult) {
    // 将状态切换为 processing
    _ = db.UpdateConsultStatus(context.Background(), consult.ConsultId, "processing")

    prompt := buildPrompt(&consult)

    // 使用独立的超时上下文，避免受请求生命周期影响
    ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
    defer cancel()

    // 仅使用 OpenAI 兼容接口：构造消息并调用
    model := "gpt-4o-mini" // 可根据配置或需求调整
    messages := []AIAgent.Message{
        {Role: "system", Content: "你是一位资深汽车顾问，提供简洁、准确的车型推荐。"},
        {Role: "user", Content: prompt},
    }
    content, err := AIAgent.CallOpenAICompat(ctx, model, messages)
    if err != nil {
        msg := "LLM 生成失败: " + err.Error()
        _ = db.UpdateConsultLLMFields(context.Background(), consult.ConsultId, &model, &prompt, &msg, nil)
        _ = db.UpdateConsultStatus(context.Background(), consult.ConsultId, "failed")
        return
    }

    // 写入生成结果（将 content 作为推荐文本保存；llm_response 记录同样内容便于查询）
    _ = db.UpdateConsultLLMFields(context.Background(), consult.ConsultId, &model, &prompt, &content, &content)
    _ = db.UpdateConsultStatus(context.Background(), consult.ConsultId, "completed")
}

// buildPrompt 构造推荐提示词
func buildPrompt(c *db.Consult) string {
    // 将用户偏好组装为简洁的中文提示，要求输出简明推荐列表
    // 注意：为了最大兼容性，我们要求返回纯文本，服务端直接存储于 recommendations
    var uid string
    if c.UserId != nil {
        uid = strconv.FormatInt(*c.UserId, 10)
    } else {
        uid = "匿名用户"
    }
    // 使用安全的指针解引用
    s := func(p *string) string {
        if p == nil {
            return "未提供"
        }
        return *p
    }

    return "你是一位资深汽车顾问，请基于以下偏好，给出3-5条适合的车型推荐，每条包含车型名与简短理由，输出为纯文本即可：\n" +
        "用户:" + uid + "\n" +
        "预算范围:" + s(c.BudgetRange) + "\n" +
        "偏好类型:" + s(c.PreferredType) + "\n" +
        "使用场景:" + s(c.UseCase) + "\n" +
        "能源类型:" + s(c.FuelType) + "\n" +
        "品牌偏好:" + s(c.BrandPreference)
}