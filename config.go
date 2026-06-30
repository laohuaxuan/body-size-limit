package main

import (
	"encoding/json"
	"errors"
	"strconv"
	"strings"

	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm"
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/types"
)

// Config 定义插件的配置结构体
type Config struct {
	MaxBodySizeStr string `json:"max_body_size" yaml:"max_body_size"`
}

const defaultMaxBodySizeStr = "10m"

func parseConfigFromData(data []byte) (Config, error) {
	cfg := Config{MaxBodySizeStr: defaultMaxBodySizeStr}
	if len(data) > 0 {
		if err := json.Unmarshal(data, &cfg); err != nil {
			return Config{}, err
		}
	}

	cfg.MaxBodySizeStr = strings.TrimSpace(cfg.MaxBodySizeStr)
	if cfg.MaxBodySizeStr == "" {
		cfg.MaxBodySizeStr = defaultMaxBodySizeStr
	}
	return cfg, nil
}

func parseMaxBodySizeFromConfig(cfg Config) (uint64, error) {
	maxBytes, err := parseSizeToBytes(cfg.MaxBodySizeStr)
	if err != nil {
		return 0, err
	}
	if maxBytes == 0 {
		maxBytes = 10 * 1024 * 1024
	}
	return maxBytes, nil
}

// parseSizeToBytes 将字符串转换为字节数
func parseSizeToBytes(sizeStr string) (uint64, error) {
	sizeStr = strings.TrimSpace(sizeStr)
	if sizeStr == "" {
		return 0, errors.New("size string is empty")
	}

	sizeStr = strings.ToLower(sizeStr)

	//获取最后一个字符作为单位
	unit := sizeStr[len(sizeStr)-1:]
	numStr := sizeStr

	//如果最后一个字符不是数字，则视为单位
	if _, err := strconv.Atoi(unit); err != nil {
		numStr = sizeStr[:len(sizeStr)-1]
	}

	num, err := strconv.ParseUint(numStr, 10, 64)
	if err != nil {
		return 0, err
	}

	//根据单位转换为字节数
	switch unit {
	case "k":
		return num * 1024, nil
	case "m":
		return num * 1024 * 1024, nil
	case "g":
		return num * 1024 * 1024 * 1024, nil
	default:
		//如果没有单位或单位不识别，默认当做字节处理
		return num, nil
	}
}

// initConfig 初始化并解析插件配置
func initConfig() (*Config, types.OnPluginStartStatus) {
	data, err := proxywasm.GetPluginConfiguration()
	if err != nil {
		proxywasm.LogErrorf("failed to get plugin configuration: %v", err)
		return nil, types.OnPluginStartStatusFailed
	}

	cfg, err := parseConfigFromData(data)
	if err != nil {
		proxywasm.LogErrorf("failed to parse plugin configuration: %v", err)
		return nil, types.OnPluginStartStatusFailed
	}

	maxBytes, err := parseMaxBodySizeFromConfig(cfg)
	if err != nil {
		proxywasm.LogErrorf("failed to parse max body size value: %v", err)
		return nil, types.OnPluginStartStatusFailed
	}

	proxywasm.LogInfof("startup max_body_size(default or global): %s (%d bytes)", cfg.MaxBodySizeStr, maxBytes)
	return &cfg, types.OnPluginStartStatusOK
}
