# 编译与准备阶段
## 编译
tinygo build -o main.wasm -target=wasi -buildmode=c-shared -scheduler=none -no-debug .
## 准备配置文件
```
{
    "max_body_size": "50m"
}
```
*可填写的单位：*
- "k": 千字节（1024 B）
- "m": 兆字节（1024 KB）
- "g": 吉字节（1024 MB）
- "": 空代表字节

 # 上传自定义插件
 ## 登录阿里云MSE微服务引擎控制台
 ## 进入网关实例
 ## 进入插件市场
 ### 在自定义插件页上传.wasm后缀插件
 ### 填写信息：
*名称：如request-body-limiter*
*描述：限制请求体大小，支持标准和分块传输*

# 启用插件
- 插件配置新建规则
- 配置规则
```
{
    "max_body_size": "50m"
}
```

# 验证
- 发送一个Body大于50MB的请求，网关应立即返回 413 Request Entity Too Large 状态码
- 分块传输：使用curl -T 发送分块传输的大文件，验证是否在接收过程及时阻断并返回413

## 本地测试
### 运行docker compose
docker compose -f docker-compose.yaml up --force-recreate
### 停滞docker compose
docker compose -f docker-compose.yaml down -v 
### 运行测试脚本
bash test/test.sh