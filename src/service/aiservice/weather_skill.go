package aiservice

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"hi-go/src/config"
	"hi-go/src/utils/cache"
)

// WeatherSkill 天气查询技能
type WeatherSkill struct {
	config config.WeatherSkillConfig
	cache  cache.Cache
	client *http.Client
}

// NewWeatherSkill 创建天气查询技能
func NewWeatherSkill(cfg config.WeatherSkillConfig, c cache.Cache) *WeatherSkill {
	return &WeatherSkill{
		config: cfg,
		cache:  c,
		client: &http.Client{
			Timeout: time.Duration(cfg.Timeout) * time.Second,
		},
	}
}

// Name 返回技能名称
func (w *WeatherSkill) Name() string {
	return "get_weather"
}

// Description 返回技能描述
func (w *WeatherSkill) Description() string {
	return "获取指定城市的实时天气信息，包括温度、湿度、天气状况等。支持中国主要城市名称，如：北京、上海、深圳、广州等"
}

// Parameters 返回参数定义
func (w *WeatherSkill) Parameters() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"location": map[string]interface{}{
				"type":        "string",
				"description": "城市名称，如：北京、上海、深圳、广州",
			},
		},
		"required": []string{"location"},
	}
}

// IsEnabled 是否启用
func (w *WeatherSkill) IsEnabled() bool {
	return w.config.Enabled
}

// Execute 执行技能
func (w *WeatherSkill) Execute(ctx context.Context, params map[string]interface{}) (interface{}, error) {
	location, ok := params["location"].(string)
	if !ok || location == "" {
		return nil, fmt.Errorf("invalid location parameter")
	}

	// 检查缓存
	cacheKey := fmt.Sprintf("weather:%s", location)
	if cached, err := w.cache.Get(ctx, cacheKey); err == nil && cached != "" {
		var result WeatherResponse
		if err := json.Unmarshal([]byte(cached), &result); err == nil {
			return result, nil
		}
	}

	// 调用天气API
	weather, err := w.fetchWeather(ctx, location)
	if err != nil {
		return nil, err
	}

	// 缓存结果
	if data, err := json.Marshal(weather); err == nil {
		_ = w.cache.Set(ctx, cacheKey, string(data), time.Duration(w.config.CacheTTL)*time.Second)
	}

	return weather, nil
}

// fetchWeather 获取天气数据
func (w *WeatherSkill) fetchWeather(ctx context.Context, location string) (*WeatherResponse, error) {
	switch w.config.Provider {
	case "qweather":
		return w.fetchQWeather(ctx, location)
	case "openweather":
		return w.fetchOpenWeather(ctx, location)
	default:
		return nil, fmt.Errorf("unsupported weather provider: %s", w.config.Provider)
	}
}

// fetchQWeather 从和风天气获取数据
func (w *WeatherSkill) fetchQWeather(ctx context.Context, location string) (*WeatherResponse, error) {
	// 先获取城市ID
	cityID, err := w.getCityID(ctx, location)
	if err != nil {
		return nil, fmt.Errorf("获取城市ID失败: %w", err)
	}

	// 使用城市ID查询天气
	apiURL := fmt.Sprintf("%s/weather/now?location=%s&key=%s",
		w.config.BaseURL, cityID, w.config.APIKey)

	req, err := http.NewRequestWithContext(ctx, "GET", apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}

	resp, err := w.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求天气API失败: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}

	// 检查 HTTP 状态码
	if resp.StatusCode != http.StatusOK {
		// 尝试解析为和风天气的错误响应格式
		var errResp QWeatherErrorResponse
		if err := json.Unmarshal(body, &errResp); err == nil && errResp.Error.Status != 0 {
			// 新版API错误格式
			return nil, fmt.Errorf("和风天气API错误(HTTP %d): %s - %s",
				errResp.Error.Status, errResp.Error.Title, errResp.Error.Detail)
		}

		// 如果不是标准错误格式，返回原始响应
		preview := string(body)
		if len(preview) > 200 {
			preview = preview[:200] + "..."
		}
		return nil, fmt.Errorf("和风天气API返回HTTP %d，响应: %s。请检查API Key是否有效", resp.StatusCode, preview)
	}

	var qweatherResp QWeatherResponse
	if err := json.Unmarshal(body, &qweatherResp); err != nil {
		preview := string(body)
		if len(preview) > 200 {
			preview = preview[:200] + "..."
		}
		return nil, fmt.Errorf("解析天气数据失败: %w，响应内容: %s", err, preview)
	}

	// 检查API业务状态码
	if qweatherResp.Code != "200" {
		// 提供更友好的错误信息
		var errMsg string
		switch qweatherResp.Code {
		case "204":
			errMsg = fmt.Sprintf("未找到城市 '%s' 的天气数据", location)
		case "400":
			errMsg = "请求参数错误"
		case "401":
			errMsg = "API Key无效，请检查配置"
		case "402":
			errMsg = "API调用次数已超出限制"
		case "403":
			errMsg = "无访问权限"
		case "404":
			errMsg = "数据不存在"
		case "429":
			errMsg = "请求过于频繁"
		case "500":
			errMsg = "和风天气服务器错误"
		default:
			errMsg = fmt.Sprintf("未知错误，代码: %s", qweatherResp.Code)
		}
		return nil, fmt.Errorf("和风天气API错误: %s", errMsg)
	}

	// 转换为统一格式
	return &WeatherResponse{
		Location:    location,
		Temperature: qweatherResp.Now.Temp + "°C",
		FeelsLike:   qweatherResp.Now.FeelsLike + "°C",
		Humidity:    qweatherResp.Now.Humidity + "%",
		Weather:     qweatherResp.Now.Text,
		WindSpeed:   qweatherResp.Now.WindSpeed + "km/h",
		UpdateTime:  qweatherResp.UpdateTime,
	}, nil
}

// getCityID 从城市名称获取城市ID
func (w *WeatherSkill) getCityID(ctx context.Context, location string) (string, error) {
	// 先检查缓存
	cacheKey := fmt.Sprintf("city_id:%s", location)
	if cityID, err := w.cache.Get(ctx, cacheKey); err == nil && cityID != "" {
		return cityID, nil
	}

	// 使用和风天气的城市查询接口（注意：使用 geoapi 而不是 devapi）
	// v2版本的城市查询接口
	apiURL := fmt.Sprintf("https://geoapi.qweather.com/v2/city/lookup?location=%s&key=%s",
		url.QueryEscape(location), w.config.APIKey)

	req, err := http.NewRequestWithContext(ctx, "GET", apiURL, nil)
	if err != nil {
		return "", fmt.Errorf("创建城市查询请求失败: %w", err)
	}

	resp, err := w.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("城市查询请求失败: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("读取城市查询响应失败: %w", err)
	}

	// 检查 HTTP 状态码
	if resp.StatusCode != http.StatusOK {
		// 尝试解析错误响应
		var errResp QWeatherErrorResponse
		if err := json.Unmarshal(body, &errResp); err == nil && errResp.Error.Status != 0 {
			return "", fmt.Errorf("城市查询API错误(HTTP %d): %s - %s",
				errResp.Error.Status, errResp.Error.Title, errResp.Error.Detail)
		}

		preview := string(body)
		if len(preview) > 200 {
			preview = preview[:200] + "..."
		}
		return "", fmt.Errorf("城市查询返回HTTP %d，响应: %s", resp.StatusCode, preview)
	}

	var cityResp QWeatherCityResponse
	if err := json.Unmarshal(body, &cityResp); err != nil {
		preview := string(body)
		if len(preview) > 200 {
			preview = preview[:200] + "..."
		}
		return "", fmt.Errorf("解析城市查询响应失败: %w，响应内容: %s", err, preview)
	}

	// 检查API业务状态码
	if cityResp.Code != "200" {
		return "", fmt.Errorf("城市查询失败，状态码: %s", cityResp.Code)
	}

	if len(cityResp.Location) == 0 {
		return "", fmt.Errorf("未找到城市: %s", location)
	}

	// 使用第一个匹配的城市
	cityID := cityResp.Location[0].ID

	// 缓存城市ID（永久缓存，城市ID不会改变）
	_ = w.cache.Set(ctx, cacheKey, cityID, 0)

	return cityID, nil
}

// fetchOpenWeather 从OpenWeather获取数据
func (w *WeatherSkill) fetchOpenWeather(ctx context.Context, location string) (*WeatherResponse, error) {
	// OpenWeather API 实现（可选）
	return nil, fmt.Errorf("openweather provider not implemented yet")
}

// WeatherResponse 统一的天气响应格式
type WeatherResponse struct {
	Location    string `json:"location"`
	Temperature string `json:"temperature"`
	FeelsLike   string `json:"feels_like"`
	Humidity    string `json:"humidity"`
	Weather     string `json:"weather"`
	WindSpeed   string `json:"wind_speed"`
	UpdateTime  string `json:"update_time"`
}

// QWeather API 响应结构
type QWeatherResponse struct {
	Code       string `json:"code"`
	UpdateTime string `json:"updateTime"`
	Now        struct {
		Temp      string `json:"temp"`
		FeelsLike string `json:"feelsLike"`
		Text      string `json:"text"`
		Humidity  string `json:"humidity"`
		WindSpeed string `json:"windSpeed"`
	} `json:"now"`
}

// QWeather API 错误响应结构（新版API）
type QWeatherErrorResponse struct {
	Error struct {
		Status int    `json:"status"`
		Type   string `json:"type"`
		Title  string `json:"title"`
		Detail string `json:"detail"`
	} `json:"error"`
}

// QWeather 城市查询响应结构
type QWeatherCityResponse struct {
	Code     string `json:"code"`
	Location []struct {
		Name    string `json:"name"`
		ID      string `json:"id"`
		Lat     string `json:"lat"`
		Lon     string `json:"lon"`
		Country string `json:"country"`
		Adm1    string `json:"adm1"` // 省份
		Adm2    string `json:"adm2"` // 城市
	} `json:"location"`
}
