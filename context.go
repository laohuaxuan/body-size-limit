package main

import "github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/types"

// pluginContext 插件级上下文，所有请求共享
type pluginContext struct {
	types.DefaultPluginContext
	maxBodySize uint64
}

// OnPluginStart 插件启动时调用,初始化配置
func (ctx *pluginContext) OnPluginStart(pluginConfigurationSize int) types.OnPluginStartStatus {
	cfg, status := initConfig()
	if status != types.OnPluginStartStatusOK {
		return status
	}

	maxBytes, _ := parseSizeToBytes(cfg.MaxBodySizeStr)
	if maxBytes == 0 {
		maxBytes = 10 * 1024 * 1024
	}

	ctx.maxBodySize = maxBytes
	return types.OnPluginStartStatusOK
}

// NewHttpContext 为每个HTTP请求创建独立的上下文
func (ctx *pluginContext) NewHttpContext(contextID uint32) types.HttpContext {
	return &httpContext{
		pluginCtx: ctx,
		maxBodySize: ctx.maxBodySize,
	}
}

// httpContext 请求及上下文，记录单个请求的状态
type httpContext struct {
	types.DefaultHttpContext
	pluginCtx      *pluginContext
	maxBodySize    uint64
	currentBodyLen uint64
	isChunked      bool
}
