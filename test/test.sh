#!/usr/bin/env bash
set -euo pipefail
BASE_URL="${BASE_URL:-http://127.0.0.1:11000}"

pass() { echo "✅ $1"; }
fail() { echo "❌ $1"; exit 1; }

assert_status() {
  local name="$1"
  local expected="$2"
  local actual="$3"
  if [[ "$actual" == "$expected" ]]; then
    pass "$name (status=$actual)"
  else
    fail "$name (expected=$expected, actual=$actual)"
  fi
}

echo "== body-size-limit smoke tests =="
echo "BASE_URL=$BASE_URL"
# 1) GET 放行
status_get=$(curl -sS -o /dev/null -w "%{http_code}" "$BASE_URL/get")
assert_status "GET /get should pass" "200" "$status_get"

# 2) 小包 Content-Length 放行（512KB）
small_file="$(mktemp)"
dd if=/dev/zero of="$small_file" bs=1k count=512 status=none
status_small=$(curl -sS -o /dev/null -w "%{http_code}" -X POST --data-binary @"$small_file" "$BASE_URL/post")
rm -f "$small_file"
assert_status "Small POST should pass" "200" "$status_small"

# 3) 大包 Content-Length 拦截（2MB，max_body_size=1m 时应 413）
large_file="$(mktemp)"
dd if=/dev/zero of="$large_file" bs=1k count=2048 status=none
status_large=$(curl -sS -o /dev/null -w "%{http_code}" -X POST --data-binary @"$large_file" "$BASE_URL/post")
rm -f "$large_file"
assert_status "Large POST(Content-Length) should block" "413" "$status_large"

# 4) 分块传输大包拦截（2MB chunked）
# 用 Python 手动构造合法 chunked 帧，避免 curl 管道场景下出现 400（非法 chunked 体）
status_chunked=$(
  BASE_URL="$BASE_URL" python3 - <<'PY'
import os
import socket
import ssl
import urllib.parse

base = os.environ.get("BASE_URL", "http://127.0.0.1:11000")
u = urllib.parse.urlparse(base)
scheme = u.scheme or "http"
host = u.hostname or "127.0.0.1"
port = u.port or (443 if scheme == "https" else 80)
path_prefix = u.path.rstrip("/") if u.path else ""
path = f"{path_prefix}/post"

sock = socket.create_connection((host, port), timeout=30)
if scheme == "https":
    ctx = ssl.create_default_context()
    sock = ctx.wrap_socket(sock, server_hostname=host)

host_header = f"{host}:{port}"
req_headers = [
    f"POST {path} HTTP/1.1",
    f"Host: {host_header}",
    "User-Agent: chunked-tester/1.0",
    "Content-Type: application/octet-stream",
    "Transfer-Encoding: chunked",
    "Connection: close",
    "",
    "",
]
sock.sendall("\r\n".join(req_headers).encode())

total = 2 * 1024 * 1024
chunk = b"A" * 65536
sent = 0
while sent < total:
    n = min(len(chunk), total - sent)
    sock.sendall(f"{n:x}\r\n".encode())
    sock.sendall(chunk[:n])
    sock.sendall(b"\r\n")
    sent += n

sock.sendall(b"0\r\n\r\n")

resp_head = b""
while b"\r\n" not in resp_head:
    part = sock.recv(1)
    if not part:
        break
    resp_head += part

sock.close()
line = resp_head.decode(errors="replace").strip()
parts = line.split()
status = parts[1] if len(parts) >= 2 else "000"
print(status)
PY
)
assert_status "Large POST(chunked) should block" "413" "$status_chunked"

echo "🎉 all tests passed"