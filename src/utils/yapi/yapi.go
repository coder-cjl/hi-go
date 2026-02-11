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
	"strings"
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

	logger.Info("开始同步 Swagger 文档到 YApi（覆盖模式）",
		zap.String("server_url", config.Config.YApi.ServerURL))

	// 读取 swagger.json 文件
	swaggerData, err := os.ReadFile("./docs/swagger.json")
	if err != nil {
		logger.Error("读取 swagger.json 失败", zap.Error(err))
		return fmt.Errorf("读取 swagger.json 失败: %w", err)
	}

	// 预处理 Swagger JSON，展开 allOf 引用，避免 YApi 混淆
	processedSwagger, err := PreprocessSwaggerForYApi(swaggerData)
	if err != nil {
		logger.Error("预处理 Swagger 文档失败", zap.Error(err))
		return fmt.Errorf("预处理 Swagger 文档失败: %w", err)
	}

	// 构建导入请求
	importReq := YApiImportRequest{
		Type:      "swagger",
		JSON:      processedSwagger,
		Token:     config.Config.YApi.Token,
		MergeMode: "normal", // 使用普通模式，完全覆盖旧接口，确保文档最新
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

// PreprocessSwaggerForYApi 预处理 Swagger JSON，展开 allOf 引用
// 解决 YApi 导入时将所有 allOf 定义混合的问题
func PreprocessSwaggerForYApi(swaggerData []byte) (string, error) {
	var swagger map[string]interface{}
	if err := json.Unmarshal(swaggerData, &swagger); err != nil {
		return "", err
	}

	paths, ok := swagger["paths"].(map[string]interface{})
	if !ok {
		return string(swaggerData), nil
	}

	definitions, ok := swagger["definitions"].(map[string]interface{})
	if !ok {
		return string(swaggerData), nil
	}

	// 遍历所有路径和响应，展开 allOf 引用
	for pathKey, pathItem := range paths {
		pathMap, ok := pathItem.(map[string]interface{})
		if !ok {
			continue
		}

		for methodKey, method := range pathMap {
			methodMap, ok := method.(map[string]interface{})
			if !ok {
				continue
			}

			responses, ok := methodMap["responses"].(map[string]interface{})
			if !ok {
				continue
			}

			for statusCode, response := range responses {
				respMap, ok := response.(map[string]interface{})
				if !ok {
					continue
				}

				schema, ok := respMap["schema"].(map[string]interface{})
				if !ok {
					continue
				}

				// 展开 allOf 并为每个接口创建唯一的定义
				if hasAllOf(schema) {
					// 生成唯一的定义名称：路径_方法_状态码_Response
					// 例如：UserLogin200Response, UserProfile200Response
					defName := generateUniqueDefName(pathKey, methodKey, statusCode)

					// 展开 allOf 生成完整定义
					expandedSchema := expandAllOfInline(schema, definitions)

					// 将展开后的定义添加到 definitions 中
					definitions[defName] = expandedSchema

					// 替换响应 schema 为引用
					respMap["schema"] = map[string]interface{}{
						"$ref": "#/definitions/" + defName,
					}
				}
			}
		}
	}

	// 重新序列化
	processedData, err := json.Marshal(swagger)
	if err != nil {
		return "", err
	}

	return string(processedData), nil
}

// generateUniqueDefName 为接口生成唯一的定义名称
func generateUniqueDefName(path, method, statusCode string) string {
	// 清理路径，生成可读的名称
	// /user/login -> UserLogin
	// /home/list -> HomeList
	name := path
	name = strings.ReplaceAll(name, "/", " ")
	name = strings.TrimSpace(name)
	name = strings.ReplaceAll(name, " ", "")
	name = strings.Title(name)

	// 添加方法和状态码
	return name + strings.Title(method) + statusCode + "Response"
}

// hasAllOf 检查 schema 是否包含 allOf
func hasAllOf(schema map[string]interface{}) bool {
	_, ok := schema["allOf"]
	return ok
}

// expandAllOfInline 展开 allOf 返回完整的 schema
func expandAllOfInline(schema map[string]interface{}, definitions map[string]interface{}) map[string]interface{} {
	allOf, ok := schema["allOf"].([]interface{})
	if !ok || len(allOf) == 0 {
		return schema
	}

	// 合并所有 allOf 中的属性
	mergedProps := make(map[string]interface{})
	var mergedRequired []interface{}
	mergedType := "object"

	for _, item := range allOf {
		itemMap, ok := item.(map[string]interface{})
		if !ok {
			continue
		}

		// 处理 $ref 引用
		if ref, hasRef := itemMap["$ref"].(string); hasRef {
			// 解析引用，如 "#/definitions/model.Response"
			refName := strings.TrimPrefix(ref, "#/definitions/")
			refDef, ok := definitions[refName].(map[string]interface{})
			if ok {
				if props, ok := refDef["properties"].(map[string]interface{}); ok {
					for k, v := range props {
						mergedProps[k] = v
					}
				}
				if req, ok := refDef["required"].([]interface{}); ok {
					mergedRequired = append(mergedRequired, req...)
				}
				if t, ok := refDef["type"].(string); ok && t != "" {
					mergedType = t
				}
			}
		}

		// 合并直接定义的属性
		if props, ok := itemMap["properties"].(map[string]interface{}); ok {
			for k, v := range props {
				mergedProps[k] = v
			}
		}
		if req, ok := itemMap["required"].([]interface{}); ok {
			mergedRequired = append(mergedRequired, req...)
		}
		if t, ok := itemMap["type"].(string); ok && t != "" {
			mergedType = t
		}
	}

	// 构建新的 schema
	result := map[string]interface{}{
		"type": mergedType,
	}
	if len(mergedProps) > 0 {
		result["properties"] = mergedProps
	}
	if len(mergedRequired) > 0 {
		result["required"] = mergedRequired
	}

	return result
}
