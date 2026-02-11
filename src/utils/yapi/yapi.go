package yapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"hi-go/src/config"
	"hi-go/src/utils/logger"
	"io"
	"net/http"
	"os"
	"time"

	"go.uber.org/zap"
)

// YApiImportRequest YApi 导入请求结构
type YApiImportRequest struct {
	Type      string `json:"type"`  // 导入类型：swagger
	JSON      string `json:"json"`  // Swagger JSON 内容
	Token     string `json:"token"` // 项目 Token
	MergeMode string `json:"merge"` // 合并模式：normal, good, merge
}

// YApiResponse YApi 响应结构
type YApiResponse struct {
	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
}

// SyncToYApi 同步 Swagger 文档到 YApi
func SyncToYApi() error {
	// 检查是否启用 YApi 同步
	if !config.Config.YApi.Enabled {
		logger.Debug("YApi 同步未启用，跳过")
		return nil
	}

	// 验证配置
	if config.Config.YApi.ServerURL == "" || config.Config.YApi.Token == "" {
		logger.Warn("YApi 配置不完整，跳过同步",
			zap.String("server_url", config.Config.YApi.ServerURL),
			zap.String("token", config.Config.YApi.Token))
		return nil
	}

	logger.Info("开始同步 Swagger 文档到 YApi",
		zap.String("server_url", config.Config.YApi.ServerURL))

	// 读取 swagger.json 文件
	swaggerData, err := os.ReadFile("./docs/swagger.json")
	if err != nil {
		logger.Error("读取 swagger.json 失败", zap.Error(err))
		return fmt.Errorf("读取 swagger.json 失败: %w", err)
	}

	// 构建导入请求
	importReq := YApiImportRequest{
		Type:      "swagger",
		JSON:      string(swaggerData),
		Token:     config.Config.YApi.Token,
		MergeMode: "merge", // 使用合并模式，不会删除已有接口
	}

	// 序列化请求数据
	reqData, err := json.Marshal(importReq)
	if err != nil {
		logger.Error("序列化请求数据失败", zap.Error(err))
		return fmt.Errorf("序列化请求数据失败: %w", err)
	}

	// 构建 YApi 导入接口 URL
	apiURL := fmt.Sprintf("%s/api/open/import_data", config.Config.YApi.ServerURL)

	// 发送 HTTP 请求
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Post(apiURL, "application/json", bytes.NewReader(reqData))
	if err != nil {
		logger.Error("发送请求到 YApi 失败", zap.Error(err))
		return fmt.Errorf("发送请求到 YApi 失败: %w", err)
	}
	defer resp.Body.Close()

	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Error("读取 YApi 响应失败", zap.Error(err))
		return fmt.Errorf("读取 YApi 响应失败: %w", err)
	}

	// 解析响应
	var yapiResp YApiResponse
	if err := json.Unmarshal(body, &yapiResp); err != nil {
		logger.Error("解析 YApi 响应失败",
			zap.Error(err),
			zap.String("response", string(body)))
		return fmt.Errorf("解析 YApi 响应失败: %w", err)
	}

	// 检查响应状态
	if yapiResp.ErrCode != 0 {
		logger.Error("YApi 同步失败",
			zap.Int("errcode", yapiResp.ErrCode),
			zap.String("errmsg", yapiResp.ErrMsg))
		return fmt.Errorf("YApi 同步失败: %s", yapiResp.ErrMsg)
	}

	logger.Info("成功同步 Swagger 文档到 YApi",
		zap.String("server_url", config.Config.YApi.ServerURL))
	return nil
}
