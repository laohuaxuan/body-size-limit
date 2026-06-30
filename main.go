package main

import (
	"github.com/higress-group/wasm-go/pkg/wrapper"
)

func init() {
	wrapper.SetCtx(
		"body-size-limit",
		wrapper.ParseConfig(parseConfig),
		wrapper.ProcessRequestHeaders(onHttpRequestHeaders),
		wrapper.ProcessRequestBody(onHttpRequestBody),
	)
}

func main() {}
