package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/tidwall/gjson"
)

// Config 定义插件的配置结构体
type Config struct {
	MaxBodySizeStr string
	MaxBodySize    uint64
}

const defaultMaxBodySizeStr = "10m"

func parseConfig(json gjson.Result, cfg *Config) error {
	sizeStr := strings.TrimSpace(json.Get("max_body_size").String())
	if sizeStr == "" {
		sizeStr = defaultMaxBodySizeStr
	}

	maxBytes, err := parseSizeToBytes(sizeStr)
	if err != nil {
		return fmt.Errorf("invalid max_body_size %q: %w", sizeStr, err)
	}
	if maxBytes == 0 {
		maxBytes = 10 * 1024 * 1024
	}

	cfg.MaxBodySizeStr = sizeStr
	cfg.MaxBodySize = maxBytes
	return nil
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

