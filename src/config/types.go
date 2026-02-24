package config

// AppConfig 应用配置结构体
type AppConfig struct {
	Server        ServerConfig        `mapstructure:"server"`
	JWT           JWTConfig           `mapstructure:"jwt"`
	Database      DatabaseConfig      `mapstructure:"database"`
	Redis         RedisConfig         `mapstructure:"redis"`
	Snowflake     SnowflakeConfig     `mapstructure:"snowflake"`
	Business      BusinessConfig      `mapstructure:"business"`
	Log           LogConfig           `mapstructure:"log"`
	Elasticsearch ElasticsearchConfig `mapstructure:"elasticsearch"`
	Logstash      LogstashConfig      `mapstructure:"logstash"`
	YApi          YApiConfig          `mapstructure:"yapi"`
	AI            AIConfig            `mapstructure:"ai"`
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Port string `mapstructure:"port"`
	Mode string `mapstructure:"mode"` // debug, release, test
}

// JWTConfig JWT配置
type JWTConfig struct {
	SecretKey            string `mapstructure:"secret_key"`
	Issuer               string `mapstructure:"issuer"`
	AccessTokenDuration  int    `mapstructure:"access_token_duration"`  // 秒
	RefreshTokenDuration int    `mapstructure:"refresh_token_duration"` // 秒
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Host            string `mapstructure:"host"`
	Port            int    `mapstructure:"port"`
	Username        string `mapstructure:"username"`
	Password        string `mapstructure:"password"`
	DBName          string `mapstructure:"dbname"`
	Charset         string `mapstructure:"charset"`
	MaxOpenConns    int    `mapstructure:"max_open_conns"`
	MaxIdleConns    int    `mapstructure:"max_idle_conns"`
	ConnMaxLifetime int    `mapstructure:"conn_max_lifetime"`  // 秒
	ConnMaxIdleTime int    `mapstructure:"conn_max_idle_time"` // 秒
}

// RedisConfig Redis配置
type RedisConfig struct {
	Host       string `mapstructure:"host"`
	Port       int    `mapstructure:"port"`
	Password   string `mapstructure:"password"`
	DB         int    `mapstructure:"db"`
	TokenTTL   int    `mapstructure:"token_ttl"`   // 秒
	SessionTTL int    `mapstructure:"session_ttl"` // 秒
}

// SnowflakeConfig 雪花ID配置
type SnowflakeConfig struct {
	MachineID int64 `mapstructure:"machine_id"`
}

// BusinessConfig 业务配置
type BusinessConfig struct {
	DefaultPageSize   int `mapstructure:"default_page_size"`
	MaxPageSize       int `mapstructure:"max_page_size"`
	PasswordMinLength int `mapstructure:"password_min_length"`
	UsernameMinLength int `mapstructure:"username_min_length"`
}

// LogConfig 日志配置
type LogConfig struct {
	Level      string `mapstructure:"level"`
	Filename   string `mapstructure:"filename"`
	MaxSize    int    `mapstructure:"max_size"`
	MaxBackups int    `mapstructure:"max_backups"`
	MaxAge     int    `mapstructure:"max_age"`
	Compress   bool   `mapstructure:"compress"`
}

// ElasticsearchConfig Elasticsearch配置
type ElasticsearchConfig struct {
	Enabled  bool     `mapstructure:"enabled"`   // 是否启用 Elasticsearch 日志
	Addrs    []string `mapstructure:"addrs"`     // ES 集群地址
	Username string   `mapstructure:"username"`  // 用户名（可选）
	Password string   `mapstructure:"password"`  // 密码（可选）
	Index    string   `mapstructure:"index"`     // 索引名称
	MaxRetry int      `mapstructure:"max_retry"` // 最大重试次数
}

// LogstashConfig Logstash配置
type LogstashConfig struct {
	Enabled    bool   `mapstructure:"enabled"`     // 是否启用 Logstash
	Host       string `mapstructure:"host"`        // Logstash 服务器地址
	Port       int    `mapstructure:"port"`        // Logstash TCP 端口
	Protocol   string `mapstructure:"protocol"`    // 协议：tcp 或 udp
	Timeout    int    `mapstructure:"timeout"`     // 连接超时（秒）
	Reconnect  bool   `mapstructure:"reconnect"`   // 是否自动重连
	BufferSize int    `mapstructure:"buffer_size"` // 缓冲区大小
}

// YApiConfig YApi配置
type YApiConfig struct {
	Enabled   bool   `mapstructure:"enabled"`    // 是否启用同步
	ServerURL string `mapstructure:"server_url"` // YApi服务器地址
	ProjectID string `mapstructure:"project_id"` // 项目ID
	Token     string `mapstructure:"token"`      // 项目 Token
}

// AIConfig AI配置
type AIConfig struct {
	Enabled      bool           `mapstructure:"enabled"`       // 是否启用AI功能
	Provider     string         `mapstructure:"provider"`      // AI提供商：openai, deepseek, anthropic等
	Model        string         `mapstructure:"model"`         // 模型名称
	APIKey       string         `mapstructure:"api_key"`       // API密钥
	BaseURL      string         `mapstructure:"base_url"`      // API基础URL
	Timeout      int            `mapstructure:"timeout"`       // 请求超时（秒）
	MaxTokens    int            `mapstructure:"max_tokens"`    // 最大token数
	Temperature  float64        `mapstructure:"temperature"`   // 温度参数
	SystemPrompt string         `mapstructure:"system_prompt"` // 系统提示词
	Skills       AISkillsConfig `mapstructure:"skills"`        // AI技能配置
}

// AISkillsConfig AI技能配置
type AISkillsConfig struct {
	Weather WeatherSkillConfig `mapstructure:"weather"` // 天气查询技能
}

// WeatherSkillConfig 天气查询技能配置
type WeatherSkillConfig struct {
	Enabled    bool   `mapstructure:"enabled"`     // 是否启用
	Provider   string `mapstructure:"provider"`    // 天气API提供商：qweather, openweather
	APIKey     string `mapstructure:"api_key"`     // API密钥
	BaseURL    string `mapstructure:"base_url"`    // API基础URL
	Timeout    int    `mapstructure:"timeout"`     // 请求超时（秒）
	CacheTTL   int    `mapstructure:"cache_ttl"`   // 缓存时间（秒）
	MaxRetries int    `mapstructure:"max_retries"` // 最大重试次数
}
