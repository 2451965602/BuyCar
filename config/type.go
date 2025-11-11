package config

type mySQL struct {
	Addr     string
	Database string
	Username string
	Password string
	Charset  string
}

type server struct {
	Addr string
	Port int
}

type aiEndpoint struct {
	Url    string
	ApiKey string
}

// LLM 配置
type llm struct {
    // 当前使用的提供商，例如: "tongyi"
    Provider string
    // 是否启用 LLM 生成功能
    Enabled  bool
    // 通义千问配置
    Tongyi   tongyi `mapstructure:"tongyi"`
}

type tongyi struct {
    // 通义千问 API Key（DashScope 密钥）
    ApiKey         string `mapstructure:"api_key"`
    // 使用的模型名称，示例: qwen-plus
    Model          string `mapstructure:"model"`
    // Chat Completions 端点（可选，默认官方地址）
    Endpoint       string `mapstructure:"endpoint"`
    // 请求超时时间（秒）
    TimeoutSeconds int    `mapstructure:"timeout_seconds"`
}

type config struct {
	MySQL      mySQL
	Server     server
	AiEndpoint aiEndpoint
}
