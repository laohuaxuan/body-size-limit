package main

import (
	"strconv"

	"github.com/higress-group/proxy-wasm-go-sdk/proxywasm"
	"github.com/higress-group/proxy-wasm-go-sdk/proxywasm/types"
	"github.com/higress-group/wasm-go/pkg/wrapper"
)

func onHttpRequestHeaders(ctx wrapper.HttpContext, cfg Config) types.Action {
	contentLengthStr, err := proxywasm.GetHttpRequestHeader("Content-Length")
	if err == nil {
		size, err := strconv.ParseUint(contentLengthStr, 10, 64)
		if err == nil && size > cfg.MaxBodySize {
			proxywasm.LogWarnf("Request body too large(Content-Length): %d bytes", size)
			send413()
			return types.ActionPause
		}
	}
	return types.ActionContinue
}

func onHttpRequestBody(ctx wrapper.HttpContext, cfg Config, body []byte) types.Action {
	if uint64(len(body)) > cfg.MaxBodySize {
		proxywasm.LogWarnf("Request body exceeded limit: %d bytes", len(body))
		send413()
		return types.ActionPause
	}
	return types.ActionContinue
}

func send413() {
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
