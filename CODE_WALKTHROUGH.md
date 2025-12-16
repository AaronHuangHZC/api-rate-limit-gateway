# 代码学习指南 (Code Walkthrough)

这个文档会引导你理解现有的代码，帮助你掌握每个模块的设计思路。

## 📚 学习顺序建议

### Step 1: 从入口点开始 (`cmd/gateway/main.go`)

这是整个应用的入口，理解它就能理解整个应用的启动流程。

**关键概念：**
- **依赖注入**：先创建依赖（config, logger, handler），再组装
- **中间件链**：请求的处理顺序是自下而上的（先注册的中间件先执行）
- **优雅关闭（Graceful Shutdown）**：用 channel 监听系统信号（Ctrl+C），然后优雅地关闭连接

**自己动手试试：**
1. 阅读 `main.go`，理解每个步骤
2. 尝试修改日志级别，看看输出有什么不同
3. 思考：为什么要用 goroutine 启动 server？为什么要用 channel 等待信号？

---

### Step 2: 理解配置模块 (`internal/config/config.go`)

这个模块负责从环境变量读取配置。

**关键设计点：**

1. **为什么用环境变量？**
   - 12-Factor App 原则：配置和代码分离
   - 同一个 binary 可以在不同环境运行
   - 容器化部署的标准做法

2. **为什么有默认值？**
   - 开发环境可以直接运行，不需要设置很多环境变量
   - 但生产环境应该显式设置所有配置

3. **配置验证：**
   - `FailurePolicy` 必须是 "fail-open" 或 "fail-closed"
   - 启动时立即发现配置错误，而不是运行时

**自己动手试试：**
1. 阅读 `config.go`，看看每个配置项的作用
2. 运行 `go test ./internal/config/... -v`，理解测试覆盖了什么
3. 尝试设置 `GATEWAY_FAILURE_POLICY=invalid`，看看会发生什么

---

### Step 3: 理解日志模块 (`internal/observability/logger.go`)

这个模块提供结构化日志功能。

**关键设计点：**

1. **为什么用 zerolog？**
   - 性能好（零分配）
   - API 简单
   - 支持结构化日志（JSON）

2. **结构化日志 vs 普通日志：**
   ```
   普通日志: "User 123 accessed /api/users"
   结构化: {"level":"info","user_id":123,"path":"/api/users","time":"2024-01-01T00:00:00Z"}
   ```
   - 结构化日志可以被日志系统（ELK, Loki）解析和查询
   - 可以按字段筛选：`user_id == 123 AND path == "/api/users"`

3. **Request ID 追踪：**
   - `WithRequestID()` 方法可以给日志添加 request_id 字段
   - 这样可以在分布式系统中追踪一个请求的所有日志

**自己动手试试：**
1. 阅读 `logger.go`，理解每个方法的作用
2. 在代码中添加一些日志，设置 `LOG_FORMAT=console` 和 `LOG_FORMAT=json`，看看输出区别
3. 思考：为什么 `WithRequestID()` 返回一个新的 Logger 而不是修改原来的？

---

### Step 4: 理解 HTTP 处理 (`internal/gateway/handler.go` 和 `middleware.go`)

这部分处理 HTTP 请求。

**关键设计点：**

1. **Handler 结构：**
   ```go
   type Handler struct {
       logger *observability.Logger  // 依赖注入
   }
   ```
   - Handler 持有依赖（logger），而不是全局变量
   - 这样便于测试（可以注入 mock logger）

2. **中间件模式：**
   ```
   请求 → RequestIDMiddleware → LoggingMiddleware → Handler
   ```
   - 每个中间件负责一个职责
   - 可以组合、重用

3. **Response Writer 包装：**
   - `responseWriter` 包装了 `http.ResponseWriter`
   - 目的是捕获 status code（标准库的 ResponseWriter 不提供这个功能）
   - 然后在 LoggingMiddleware 中记录 status code

**自己动手试试：**
1. 阅读 `handler.go` 和 `middleware.go`
2. 尝试添加一个中间件，比如记录请求体大小
3. 思考：为什么 `Health` 和 `Ready` 是两个不同的端点？（Hint: K8s liveness vs readiness probe）

---

## 🔍 关键代码片段解析

### 1. Graceful Shutdown (main.go:73-89)

```go
quit := make(chan os.Signal, 1)
signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
<-quit  // 阻塞等待信号
```

**为什么重要？**
- 容器编排系统（K8s）会发送 SIGTERM 给容器
- 优雅关闭让正在处理的请求完成，而不是直接杀掉进程
- 这是生产环境的标准做法

**自己试试：**
1. 启动 server
2. 发送一些请求
3. 按 Ctrl+C，观察日志输出

---

### 2. 中间件链的执行顺序

```go
router.Use(middleware.Recoverer)           // 1. 最先执行（处理 panic）
router.Use(middleware.RequestID)           // 2. 
router.Use(gateway.RequestIDMiddleware())  // 3. 我们的自定义 RequestID
router.Use(gateway.LoggingMiddleware())    // 4. 最后执行（记录最终结果）
```

**执行顺序（请求时）：**
- Recoverer → RequestID → RequestIDMiddleware → LoggingMiddleware → Handler

**返回时（响应）：**
- Handler → LoggingMiddleware → RequestIDMiddleware → RequestID → Recoverer

**自己试试：**
在 LoggingMiddleware 中添加日志，看看执行顺序

---

### 3. Context 传递 Request ID

```go
ctx := context.WithValue(r.Context(), "request_id", requestID)
next.ServeHTTP(w, r.WithContext(ctx))
```

**为什么用 Context？**
- Context 是 Go 中传递请求级别数据的标准方式
- 任何 handler 都可以从 context 中获取 request_id
- 避免修改 `*http.Request` 结构

**自己试试：**
在 Handler 中尝试从 context 获取 request_id：
```go
requestID := r.Context().Value("request_id").(string)
```

---

## 📝 练习题（加深理解）

### 练习 1: 添加一个新的配置项
在 `config.go` 中添加一个新配置 `LOG_SAMPLE_RATE`（日志采样率），默认值 1.0（100% 采样）

### 练习 2: 添加一个新的中间件
创建一个 `TimingMiddleware`，在响应头中添加 `X-Response-Time` 字段

### 练习 3: 修改 Health Check
让 `/ready` 端点检查 Redis 连接是否正常（现在先返回 mock 结果）

### 练习 4: 理解测试
阅读 `config_test.go`，理解测试的结构，尝试写一个测试来验证 `PostgresDSN()` 的输出格式

---

## 🎯 接下来你要做什么？

从现在开始，我会**更多引导你写代码**，而不是直接给你完整代码。我会：

1. **解释设计思路**：为什么这样设计，tradeoff 是什么
2. **提供代码框架**：给你函数签名和注释，你填充实现
3. **Code Review**：你写完代码后，我会帮你 review，指出可以改进的地方
4. **解释错误**：如果代码有问题，我会解释为什么错，怎么改

这样你才能真正学到东西！

---

## ❓ 问题时间

如果看代码时有任何问题，随时问我！比如：
- 为什么这里要用 goroutine？
- 这个设计有什么缺点？
- 有没有更好的实现方式？

