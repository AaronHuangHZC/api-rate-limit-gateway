# 动手练习 (Hands-on Exercises)

这些练习帮助你通过实际操作理解代码。

## 🚀 快速开始：运行代码

### 1. 启动基础设施（需要 Docker）

```bash
make docker-up
```

这会启动 Redis 和 Postgres。如果 Docker 没运行，先启动 Docker Desktop。

### 2. 编译代码

```bash
make build
# 或者
go build -o bin/gateway ./cmd/gateway
```

### 3. 运行 Gateway（在一个终端）

```bash
./bin/gateway
# 或者
make run
```

你应该看到：
```
{"level":"info","message":"Starting API Rate Limiting Gateway","time":"2024-..."}
{"level":"info","message":"Server listening","address":"0.0.0.0:8080","time":"2024-..."}
```

### 4. 测试 Health Endpoint（在另一个终端）

```bash
curl http://localhost:8080/health
```

应该返回：
```json
{"status":"healthy"}
```

### 5. 观察结构化日志

发送请求时，你会看到类似这样的日志：
```json
{
  "level": "info",
  "request_id": "abc-123-def",
  "message": "HTTP request completed",
  "method": "GET",
  "path": "/health",
  "status_code": 200,
  "duration_ms": 2
}
```

---

## 📖 阅读任务清单

### Task 1: 理解 main.go 的执行流程

**目标：** 画出代码执行流程图

**步骤：**
1. 打开 `cmd/gateway/main.go`
2. 按顺序阅读每一行，理解它在做什么
3. 画出执行流程图（用纸笔或者 draw.io）

**检查点：**
- [ ] 能解释为什么用 `go func()` 启动 server
- [ ] 能解释 `<-quit` 这行代码的作用
- [ ] 能解释 `srv.Shutdown(ctx)` 和直接 `os.Exit(1)` 的区别

---

### Task 2: 理解配置加载

**目标：** 理解环境变量如何被读取和使用

**步骤：**
1. 阅读 `internal/config/config.go`
2. 运行测试：`go test ./internal/config/... -v`
3. 尝试修改环境变量，看配置如何变化：

```bash
export SERVER_PORT=9000
export GATEWAY_FAILURE_POLICY=fail-open
go run ./cmd/gateway
# 观察日志中的端口号
```

**问题思考：**
- [ ] 为什么 `getEnv()` 函数要有 `defaultValue` 参数？
- [ ] 如果环境变量值是无效的（比如 `SERVER_PORT=abc`），会发生什么？
- [ ] 为什么 `PostgresDSN()` 是一个方法而不是字段？

---

### Task 3: 理解中间件链

**目标：** 理解请求如何在中间件间传递

**步骤：**
1. 阅读 `internal/gateway/middleware.go`
2. 在 `LoggingMiddleware` 中添加一行日志，看看什么时候执行：

```go
func LoggingMiddleware(logger *observability.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logger.Info("Middleware: BEFORE calling next")  // 添加这行
			start := time.Now()
			// ... 原有代码 ...
			next.ServeHTTP(wrapped, r)  // 这会调用下一个中间件或 handler
			logger.Info("Middleware: AFTER calling next")   // 添加这行
			// ... 原有代码 ...
		})
	}
}
```

3. 重新编译运行，发送请求，观察日志顺序

**问题思考：**
- [ ] 为什么 "BEFORE" 和 "AFTER" 的顺序是这样？
- [ ] 中间件的执行顺序和注册顺序有什么关系？

---

### Task 4: 理解 Request ID 追踪

**目标：** 理解如何在请求中传递 context

**步骤：**
1. 在 `handler.go` 的 `Health` 函数中，尝试获取 request ID：

```go
func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	requestID := r.Context().Value("request_id").(string)  // 添加这行
	h.logger.Infof("Health check called", map[string]interface{}{
		"request_id": requestID,  // 添加这行
	})
	// ... 原有代码 ...
}
```

2. 发送请求时，在 header 中指定 request ID：
```bash
curl -H "X-Request-ID: my-custom-id-123" http://localhost:8080/health
```

3. 观察日志，看看 request_id 是什么

**问题思考：**
- [ ] 如果没有 `X-Request-ID` header，request_id 是什么？
- [ ] Context 是什么？为什么用它传递数据而不是全局变量？

---

## 🛠️ 小练习：修改代码

### 练习 1: 添加一个新的配置项

**任务：** 添加一个 `SERVER_SHUTDOWN_TIMEOUT` 配置，控制优雅关闭的超时时间

**提示：**
1. 在 `config.go` 的 `ServerConfig` 结构体中添加字段
2. 在 `Load()` 函数中读取环境变量，默认值 10 秒
3. 在 `main.go` 中使用这个配置

**检查你的实现：**
- [ ] 默认值是合理的（10秒）
- [ ] 如果环境变量是无效的，会使用默认值
- [ ] `main.go` 中使用了这个配置

---

### 练习 2: 添加一个简单的中间件

**任务：** 创建一个 `CORSMiddleware`，添加 CORS headers

**提示：**
1. 在 `middleware.go` 中添加新函数
2. 添加这些 headers：
   - `Access-Control-Allow-Origin: *`
   - `Access-Control-Allow-Methods: GET, POST, PUT, DELETE, OPTIONS`
   - `Access-Control-Allow-Headers: Content-Type, X-API-Key`

**检查你的实现：**
- [ ] 在 `main.go` 中注册这个中间件
- [ ] 用 curl 测试：
  ```bash
  curl -H "Origin: http://example.com" -v http://localhost:8080/health
  ```
- [ ] 能看到 CORS headers 在响应中

---

### 练习 3: 改进 Health Check

**任务：** 让 `/health` 返回更多信息（版本、启动时间等）

**提示：**
1. 在 `main.go` 中定义一个版本常量
2. 在启动时记录启动时间
3. 将这些信息传递给 Handler（可能需要修改 Handler 结构）
4. 在 `Health` handler 中返回这些信息

**检查你的实现：**
- [ ] `/health` 返回 JSON 包含版本和运行时间
- [ ] 代码结构清晰（不要用全局变量）

---

## 🧪 运行测试

理解测试也很重要！

```bash
# 运行所有测试
make test

# 运行特定包的测试
go test ./internal/config/... -v

# 查看测试覆盖率
make test-coverage
```

**任务：** 阅读 `config_test.go`，理解每个测试在测什么

---

## 💡 下一步

完成这些练习后，你对代码应该有了深入理解。然后我们可以继续：

1. **Phase 1.4**: 数据库迁移（我会给你更多指导，让你写更多代码）
2. **Phase 2**: 速率限制逻辑（核心功能，你会写更多）

准备好了告诉我！

