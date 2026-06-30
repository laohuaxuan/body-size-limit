package main

import (
	"strconv"
	"strings"

	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm"
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/types"
)

// OnHttpRequestHeaders 在接收请求头时触发
func (ctx *httpContext) OnHttpRequestHeaders(numHeaders int, endOfStream bool) types.Action {
	//1. 检查Content-Length头
	contentLengthStr, err := proxywasm.GetHttpRequestHeader("Content-Length")
	if err == nil {
		size, err := strconv.ParseUint(contentLengthStr, 10, 64)
		if err == nil && size > ctx.maxBodySize {
			proxywasm.LogWarnf("Request body too large(Content-Length): %d bytes", size)
			ctx.send413()
			return types.ActionPause
		}
	}
	//2. 检查是否分块传输
	transferEncoding, err := proxywasm.GetHttpRequestHeader("Transfer-Encoding")
	if err == nil && strings.Contains(strings.ToLower(transferEncoding), "chunked") {
		ctx.isChunked = true
	}
	// 未超限，继续接收下一个 chunk
	return types.ActionContinue
}

// OnHttpRequestBody 在接收请求体时触发。
// 这里统一处理“无/有 Content-Length”的场景，避免仅依赖 Transfer-Encoding 识别导致漏拦截。
func (ctx *httpContext) OnHttpRequestBody(bodySize int, endOfStream bool) types.Action {
	// bodySize 在不同运行时里可能是“累计值”或“当前增量值”，兼容两种语义。
	current := uint64(bodySize)
	if current >= ctx.currentBodyLen {
		ctx.currentBodyLen = current
	} else {
		ctx.currentBodyLen += current
	}

	if ctx.currentBodyLen > ctx.maxBodySize {
		proxywasm.LogWarnf("Request body exceeded limit: %d bytes", ctx.currentBodyLen)
		ctx.send413()
		return types.ActionPause
	}
	// 未超限，继续接收下一个 chunk
	return types.ActionContinue
}

// send413 封装发送413响应的逻辑
func (ctx *httpContext) send413() {
	err := proxywasm.SendHttpResponse(
		413,
		[][2]string{{"Content-Type", "text/plain"}},
		[]byte("Request Entity Too Large"),
		-1,
	)
	if err != nil {
		proxywasm.LogErrorf("Failed to send 413 response: %v", err)
	}
}
