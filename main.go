package main

import (
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm"
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/types"
)

func init() {
	//设置全局的VMContext
	proxywasm.SetVMContext(&vmContext{})
}

func main() {}

// vmContext wasm虚拟机上下文实现
type vmContext struct {
	types.DefaultVMContext
}

func (*vmContext) NewPluginContext(contextID uint32) types.PluginContext {
	return &pluginContext{}
}
