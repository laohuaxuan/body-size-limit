# body-size-limit

`body-size-limit` 是一个用于限制 HTTP 请求体大小的 Wasm 插件，支持：
- 普通请求（`Content-Length`）
- 分块传输请求（chunked）

超出阈值后，插件会直接返回 `413 Request Entity Too Large`。

## 配置说明

插件配置为 JSON，示例：

```json
{
  "max_body_size": "50m"
}
```

`max_body_size` 支持以下单位（不区分大小写）：
- `k`：KB（1024 B）
- `m`：MB（1024 KB）
- `g`：GB（1024 MB）
- 无单位：按字节处理

未配置或空值时，默认使用 `10m`。

## 编译 Wasm

使用 TinyGo 以 reactor 模式编译（MSE/Envoy 必需）：

```bash
tinygo build -o main.wasm -target=wasi -buildmode=c-shared -scheduler=none -no-debug .
编译后：main.wasm
```

## 本地联调（Docker Compose）

启动：

```bash
docker compose -f docker-compose.yaml up --force-recreate
```

停止并清理：

```bash
docker compose -f docker-compose.yaml down -v
```

运行测试脚本：

```bash
bash test/test.sh
```

## 部署到阿里云 MSE 网关

1. 登录 MSE 控制台，进入目标网关实例。
2. 打开插件市场，在自定义插件中上传 `main.wasm`。
3. 创建插件规则并填写配置，例如：

```json
{
  "max_body_size": "10m"
}
```
4. 部署后截图
<img width="2340" height="373" alt="image" src="https://github.com/user-attachments/assets/127d2231-65b8-4889-857b-97c373e42094" />


4. 发布规则并等待生效。

## 验证建议

- 小于阈值请求应返回 `200`。
- 大于阈值请求应返回 `413`。
- 分块传输大包也应返回 `413`。
- 拦截截图
  <img width="1477" height="581" alt="image" src="https://github.com/user-attachments/assets/1623647c-8a55-4bcb-91d7-9e1796b548fc" />


## 常见问题

- 出现 `_start ... restricted_callback`：
  通常是 Wasm 编译模式错误，请确认使用了 `-buildmode=c-shared`。
