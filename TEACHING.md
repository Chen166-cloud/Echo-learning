# Echo-learning 后端技术栈教学文档

> 面向第一次接触 Go Web 后端的读者。文中所有代码片段都直接来自本项目，并标注了出处文件，建议对照源码阅读。
>
> 只想把服务跑起来：直接看[第二章](#二启动命令教学先把服务跑起来)。想搞懂每一层在干什么：看第三章。

## 目录

- [一、技术栈总览](#一技术栈总览)
- [二、启动命令教学（先把服务跑起来）](#二启动命令教学先把服务跑起来)
- [三、技术栈详解](#三技术栈详解)
- [四、常用命令速查](#四常用命令速查)
- [五、常见问题排查](#五常见问题排查)
- [六、进阶练习](#六进阶练习)

---

## 一、技术栈总览

本项目后端共用到 4 项核心技术（外加 MySQL 数据库本身）：

| 技术 | 本项目版本 | 角色 | 官网 |
|---|---|---|---|
| Go | 1.26 | 编程语言：编译型、自带并发支持 | https://go.dev |
| Echo | v4.16.0 | Web 框架：路由分发、中间件、参数绑定 | https://echo.labstack.com |
| GORM | v1.31.2 | ORM：用 Go 结构体操作数据库，少写 SQL | https://gorm.io |
| gorm.io/driver/mysql | v1.6.0 | GORM 的 MySQL 方言适配器 | https://gorm.io/docs/connecting_to_the_database.html |
| github.com/go-sql-driver/mysql | v1.10.1 | 底层 MySQL 协议驱动，实现标准库 `database/sql` 接口 | https://github.com/go-sql-driver/mysql |
| MySQL | 8.x | 关系型数据库，存数据的地方 | https://www.mysql.com |

它们的关系是一条清晰的调用链：

```
HTTP 请求
  → Echo（接住请求，路由到对应处理函数）
    → GORM（把 Go 结构体操作翻译成 SQL）
      → go-sql-driver（真正和 MySQL 通信）
        → MySQL（执行 SQL，返回数据）
```

所有依赖声明在 [go.mod](go.mod) 里：

```go
module echo-learning          // 本项目的模块名，也是所有 import 路径的前缀

go 1.26                       // 要求的最低 Go 版本

require (
	github.com/go-sql-driver/mysql v1.10.1  // 底层驱动
	github.com/labstack/echo/v4 v4.16.0     // Web 框架
	gorm.io/driver/mysql v1.6.0             // GORM 的 MySQL 驱动
	gorm.io/gorm v1.31.2                    // ORM
)
```

`go.sum` 则记录每个依赖的版本和校验哈希，保证任何人任何时候拉到的依赖完全一致——这两个文件合起来就是 Go 的依赖管理（详见 [3.1 节](#31-go-modulesgomod--gosum)）。

---

## 二、启动命令教学（先把服务跑起来）

### 2.1 环境准备

**① Go（≥ 1.26）**

```bash
go version        # 已装则显示版本号，例如：go version go1.26.5 darwin/arm64
```

没装的话（macOS）：

```bash
brew install go
```

**② MySQL（本机用的是 Homebrew 的 mysql@8.4）**

```bash
brew install mysql@8.4          # 安装（只需一次）
brew services start mysql@8.4   # 启动，并设为开机自启
brew services list | grep mysql # 确认状态是 started
mysql -u root -p                # 能登录说明服务正常
```

代码里的默认连接账号是 `root / 1234`。如果你的 root 还没有密码，先设一个：

```sql
ALTER USER 'root'@'localhost' IDENTIFIED BY '1234';
```

（也可以不改密码，用下面的环境变量 `DB_PASS` 传入你自己的密码。）

### 2.2 启动三步

```bash
cd Echo-learning   # 1. 进入后端项目目录（即本文件所在目录）

go mod tidy        # 2. 按 go.mod 拉取依赖到本地缓存（首次必做；以后改了 import 也要重跑）

go run .           # 3. 编译并立刻运行当前目录的 main 包
```

看到下面这行就说明启动成功了：

```
2025/xx/xx xx:xx:xx server start and listen port=8080
```

**首次启动会自动做两件事，不用手动建库建表**（实现见 [services/mysql/mysql.go](services/mysql/mysql.go) 和 [initialize/init.go](initialize/init.go)）：

1. 自动创建数据库 `echo_learning`（`CREATE DATABASE IF NOT EXISTS ...`）；
2. 通过 GORM 的 `AutoMigrate` 自动创建 `courses` 表。

启动后每次请求，终端还会打印 GORM 执行的 SQL——这是刻意配置的教学功能，方便你观察「每条 Go 代码到底执行了什么 SQL」。

### 2.3 验证服务真的活了

开一个**新终端**（原来的终端正在跑服务，别动它），依次执行：

```bash
# 新增一条课程
curl -X POST http://localhost:8080/api/courses \
  -H "Content-Type: application/json" \
  -d '{"title":"Go 入门","teacher":"张三","description":"Echo + GORM 基础","price":99.5}'

# 查询列表（应能看到刚插入的数据）
curl http://localhost:8080/api/courses

# 也可以直接查数据库
mysql -u root -p echo_learning -e "SELECT * FROM courses;"
```

curl 返回统一格式 `{"code":0,"msg":"success","data":...}`，`code=0` 即成功。

### 2.4 用环境变量改配置

连接参数全部支持环境变量覆盖，见 [services/mysql/mysql.go](services/mysql/mysql.go) 中的 `getEnv`：

| 变量 | 默认值 | 说明 |
|---|---|---|
| `DB_USER` | root | 数据库用户名 |
| `DB_PASS` | 1234 | 数据库密码 |
| `DB_HOST` | 127.0.0.1 | 数据库地址 |
| `DB_PORT` | 3306 | 数据库端口 |
| `DB_NAME` | echo_learning | 库名 |
| `PORT` | 8080 | HTTP 服务监听端口 |

用法（多个变量写在同一条命令前面，空格分隔）：

```bash
PORT=9000 DB_PASS=mypass DB_NAME=echo_dev go run .
```

注意一个实现细节：`getEnv` 的判断是「值为空串就用默认值」，所以 `DB_PASS=`（空）等于没设置。

### 2.5 `go run` 和 `go build` 的区别

| 命令 | 行为 | 适用场景 |
|---|---|---|
| `go run .` | 编译到临时目录，立刻运行，**不留下可执行文件** | 日常开发 |
| `go build .` | 在当前目录生成可执行文件 `echo-learning` | 交付 / 部署 |

```bash
go build .          # 生成二进制文件 echo-learning
./echo-learning     # 之后直接运行，目标机器上不需要装 Go 和源码
```

因为本项目是纯 Go（无 cgo 依赖），交叉编译一条命令就能产出给其他系统的二进制：

```bash
GOOS=linux GOARCH=amd64 go build .    # 给 Linux x86 服务器编译（常见部署目标）
GOOS=darwin GOARCH=arm64 go build .   # 给 Apple 芯片的 Mac
```

### 2.6 停止、重启、端口占用

- **停止**：在跑服务的终端按 `Ctrl+C`。
- **改了代码**：`go run` 没有热重载，必须 `Ctrl+C` 停掉后重新 `go run .` 才生效（前端 Vite 才有热更新，Go 后端没有）。
- **端口被占**（报错 `bind: address already in use`）：

```bash
lsof -i :8080          # 看是谁占了 8080（会列出 PID）
kill <PID>             # 杀掉它；或者干脆换个端口：PORT=8081 go run .
```

---

## 三、技术栈详解

### 3.1 Go Modules（go.mod / go.sum）

Go 从 1.16 起默认用 **Modules** 管理依赖，对标前端的 npm：

| Go Modules | npm 类比 |
|---|---|
| `go.mod` | `package.json`（声明依赖及版本） |
| `go.sum` | `package-lock.json`（锁定精确版本 + 哈希校验） |
| `~/go/pkg/mod` | `node_modules`（全局共享的依赖缓存） |
| `go mod tidy` | `npm install`（按声明拉齐依赖） |

`go.mod` 第一行 `module echo-learning` 定义了**模块名**，它就是本项目所有包的 import 路径前缀。所以代码里的导入长这样：

```go
import (
	"echo-learning/dao"             // 本项目内的包：模块名 + 目录路径
	"echo-learning/logic/course_logic"
	"github.com/labstack/echo/v4"   // 第三方包：直接用它的仓库路径
	"gorm.io/gorm"
)
```

这也是为什么目录名即包名：`dao/` 目录下所有文件 `package dao`，外界通过 `dao.Course`、`dao.CourseTable` 访问它导出的东西（大写开头 = 导出，小写 = 包内私有）。

### 3.2 Echo：Web 框架

Echo 负责一切和 HTTP 有关的事：监听端口、把 URL 分发给对应函数、解析参数、返回 JSON。

#### 创建实例与全局中间件

[initialize/init.go](initialize/init.go)：

```go
e := echo.New()
e.HideBanner = true

// 顺序注册，请不要打乱（中间件按注册顺序嵌套执行）
e.Use(middleware.Recover())   // 兜住 panic，防止一个请求崩掉整个进程
e.Use(middleware.Logger())    // 打印每个请求的访问日志（方法、路径、状态码、耗时）
e.Use(middleware.CORS())      // 允许跨域，前端 localhost:5173 才能直接调 8080
```

`Use` 注册的是**全局中间件**，对之后注册的所有路由生效。

#### 路由与分组

[router/router.go](router/router.go) 先挂一个 `/api` 根分组：

```go
rRoot := echoIns.Group("/api")
DispatchCourseRouters(rRoot)
```

[router/course_router.go](router/course_router.go) 再挂 `/courses` 子分组并注册 5 个路由：

```go
courseGrp := echoIns.Group("/courses")
courseGrp.GET("",      middleware.WrapResponse(controller.ListCourse))    // GET /api/courses
courseGrp.GET("/:id",  middleware.WrapResponse(controller.GetCourse))     // GET /api/courses/:id
courseGrp.POST("",     middleware.WrapResponse(controller.CreateCourse))  // POST /api/courses
courseGrp.PUT("/:id",  middleware.WrapResponse(controller.UpdateCourse))  // PUT /api/courses/:id
courseGrp.DELETE("/:id", middleware.WrapResponse(controller.DeleteCourse))// DELETE /api/courses/:id
```

两个要点：

- `:id` 是**路径参数**，匹配 `/api/courses/1`、`/api/courses/42` 等，值通过结构体 tag `param:"id"` 绑定（见下）。
- `Group` 的作用是复用前缀。以后加用户模块，只需 `rRoot.Group("/users")`，不用每条路由都写 `/api/users/...`。

#### handler 的固定签名

Echo 的处理函数签名固定为：

```go
func(echo.Context) error
```

本项目进一步约定 controller 统一返回 `(interface{}, error)`，再由 `WrapResponse` 适配成 Echo 要的签名（见 [3.6 节](#36-统一响应与错误处理)）。

#### 参数绑定：ctx.Bind

[controller/base.go](controller/base.go) 封装了绑定：

```go
func ParseRequest(ctx echo.Context, req interface{}) error {
	return ctx.Bind(req)
}
```

`Bind` 会根据 tag 自动从**多个来源**取值填充结构体。以更新课程的请求体为例（[model/course.go](model/course.go)）：

```go
type UpdateCourseReq struct {
	ID          int64   `json:"id" param:"id"`   // param: 从 URL 路径 /:id 取
	Title       string  `json:"title"`           // json: 从请求 body 取
	Teacher     string  `json:"teacher"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
}
```

一次 `Bind` 同时完成了「路径参数 + JSON body」的解析，前端只需正常发 `PUT /api/courses/1` + JSON。三个常用 tag 记住即可：`json`（body）、`param`（路径）、`query`（查询串，如 `?page=1`）。

#### 输出 JSON

```go
ctx.JSON(http.StatusOK, result)   // 把任意结构体按 json tag 序列化成响应
```

### 3.3 GORM：ORM

ORM（对象关系映射）让你**用 Go 结构体和链式调用代替手写 SQL**，同时保留「想写 SQL 时随时能写」的能力。

#### 定义模型

[dao/define.go](dao/define.go) 中的 `Course` 结构体和 `courses` 表一一对应：

```go
type Course struct {
	ID          int64     `gorm:"column:id;primaryKey" json:"id"`          // 主键，自增
	Title       string    `gorm:"column:title;size:128;not null" json:"title"`
	Teacher     string    `gorm:"column:teacher;size:64" json:"teacher"`
	Description string    `gorm:"column:description;size:512" json:"description"`
	Price       float64   `gorm:"column:price;default:0" json:"price"`
	CreatedAt   time.Time `gorm:"column:created_at" json:"created_at"`    // GORM 自动维护
	UpdatedAt   time.Time `gorm:"column:updated_at" json:"updated_at"`    // GORM 自动维护
}
```

常用 `gorm` tag 速查：

| tag | 含义 |
|---|---|
| `column:xxx` | 映射到哪个列（不写则按驼峰转下划线推断） |
| `primaryKey` | 主键 |
| `size:128` | 字符串长度，即 `VARCHAR(128)` |
| `not null` | 非空约束 |
| `default:0` | 默认值 |

两个「魔法字段」：只要字段名叫 `CreatedAt` / `UpdatedAt`，GORM 会在插入和更新时**自动填写时间戳**，业务代码完全不用管。

#### 建立连接

[services/mysql/mysql.go](services/mysql/mysql.go)：

```go
dsn := fmt.Sprintf(
	"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
	user, pass, host, port, name,
)
DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
	Logger: gormLogger.Default.LogMode(gormLogger.Info), // 打印 SQL，方便学习
})
```

DSN（Data Source Name，连接串）格式为 `用户:密码@tcp(地址:端口)/库名?参数`，三个参数各司其职：

- `charset=utf8mb4`：完整 UTF-8，能存 emoji；
- `parseTime=True`：把 MySQL 的 `DATETIME` 转成 Go 的 `time.Time`（不开会导致时间字段扫描失败）；
- `loc=Local`：按本地时区解析时间。

连接后还配置了**连接池**（复用连接，避免频繁握手）：

```go
sqlDB.SetMaxIdleConns(5)                 // 最多保留 5 个空闲连接
sqlDB.SetMaxOpenConns(20)                // 最多同时开 20 个连接
sqlDB.SetConnMaxLifetime(time.Hour)      // 单个连接最多活 1 小时后换新
```

#### AutoMigrate 自动建表

[initialize/init.go](initialize/init.go)：

```go
mysql.GetDB().AutoMigrate(&dao.Course{})
```

表不存在就按结构体定义建表；存在则补齐缺失的列。教学项目用它省去手动 DDL；**生产环境慎用**（它不做删列/改列这类破坏性变更，长期能力和人工迁移工具相比有限）。

#### 链式 CRUD

本项目 dao 层（[dao/impl_course.go](dao/impl_course.go)）的每条调用和它实际执行的 SQL 对照：

| Go 代码 | 实际执行的 SQL |
|---|---|
| `Table("courses").Where("id = ?", 1).Find(&course)` | `SELECT * FROM courses WHERE id = 1` |
| `Table("courses").Order("id DESC").Find(&courses)` | `SELECT * FROM courses ORDER BY id DESC` |
| `Table("courses").Create(course)` | `INSERT INTO courses (...) VALUES (...)` |
| `Table("courses").Save(course)` | `UPDATE courses SET ... WHERE id = 1` |
| `Table("courses").Where("id = ?", 1).Delete(&Course{})` | `DELETE FROM courses WHERE id = 1` |

阅读链式调用的技巧：**从左往右读就是 SQL 的书写顺序**——指定表 → 拼条件 → 执行动作。

两个容易踩的差异点：

- `Find` 查不到时**不报错**，返回零值结构体（本项目利用这点，在 logic 层用 `course.ID == 0` 判断「记录不存在」）；而 `First` 查不到会返回 `ErrRecordNotFound` 错误。
- `Save` 按**主键更新整行**（所有字段都会写回）；`Update`/`Updates` 只更新指定字段。本项目更新前先查出了原记录再覆盖字段，所以用 `Save` 是安全的。

#### WithContext：让 SQL 跟着请求走

[dao/session.go](dao/session.go)：

```go
func GetSession(ctx echo.Context) *gorm.DB {
	return mysql.GetDB().WithContext(ctx.Request().Context())
}
```

每次请求都从 `ctx` 取出**请求级 context** 传给 GORM。收益：HTTP 请求被取消/超时（比如客户端断开了），正在执行的 SQL 也会跟着取消，不会白白占着数据库连接。

### 3.4 MySQL 驱动与自动建库

你可能注意到项目 import 了**两个** mysql 相关的包，它们的分工不同：

```go
_ "github.com/go-sql-driver/mysql" // 底层驱动：注册 "mysql" 驱动给标准库 database/sql
"gorm.io/driver/mysql"              // GORM 适配层：把 GORM 的调用翻译成 database/sql 调用
```

第一行是**匿名导入**（`_` 前缀）：不直接使用该包，只为触发它 `init()` 里的驱动注册。这是 Go database/sql 的标准玩法——GORM 之下，真正干活的是它。

自动建库的逻辑在 [services/mysql/mysql.go](services/mysql/mysql.go)：先用**不带库名**的 DSN 连上 MySQL 服务器，执行 `CREATE DATABASE IF NOT EXISTS`，再用带库名的 DSN 建立 GORM 连接。为什么要两步？因为 MySQL 不允许连接一个不存在的库——先「上楼」，再「进房间」。

### 3.5 分层架构与一次请求的完整旅程

项目按 `router → controller → logic → dao` 四层组织（模仿真实生产项目 pop 的结构），每层职责单一：

| 层 | 目录 | 只负责 | 不负责 |
|---|---|---|---|
| router | [router/](router/) | 把 URL 映射到 controller 函数 | 任何业务 |
| controller | [controller/](controller/) | 绑定参数、调 logic、组装响应 | 业务规则、SQL |
| logic | [logic/](logic/) | 业务规则与编排（校验、查-改-存的顺序） | HTTP 细节、SQL 细节 |
| dao | [dao/](dao/) | 数据库读写（接口 + 实现） | 业务规则 |

入口 [main.go](main.go) 只有一行，把启动顺序交给 [initialize/init.go](initialize/init.go)：**连数据库 → 建表 → 建 Echo 实例 → 注册路由 → 监听端口**。

以「**更新课程**」`PUT /api/courses/1` 为例，走一遍完整链路：

```
① Echo 收到请求，按路由表匹配到 controller.UpdateCourse
      router/course_router.go: courseGrp.PUT("/:id", middleware.WrapResponse(controller.UpdateCourse))

② WrapResponse 中间件先包一层（统一响应，见 3.6）
      middleware/wrap_response.go

③ controller 绑定参数：路径里的 1 → req.ID，body 的 JSON → req.Title 等
      controller/controller_course.go: req := &model.UpdateCourseReq{}; ParseRequest(ctx, req)

④ logic 做业务：校验 title 非空 → 先查记录是否存在 → 覆盖字段 → 调 dao 更新
      logic/course_logic/update_course.go

⑤ dao 执行 SQL：GetSession(ctx).Table("courses").Save(course)
      dao/impl_course.go → GORM → MySQL

⑥ 结果原路返回：logic 返回 *dao.Course → controller 包成 Response{code,msg,data}
   → WrapResponse 调 ctx.JSON 输出给前端
```

**为什么值得分层？** 每层可独立替换和测试：换前端协议时 logic 不动；换数据库（MySQL → PostgreSQL）时只动 dao 和连接；业务规则集中在 logic，不用在 HTTP 处理代码里翻找。本项目只有一个 courses 表，分层的好处要等业务变复杂才明显——但结构从小搭对，比以后重构容易得多。

dao 层还用了**接口 + 实现**的经典模式：

```go
type CourseDao interface { ... }                              // 能力声明（define.go）
var CourseTable CourseDao = &courseDaoImpl{tableName: "courses"} // 全局唯一实例
```

logic 层只依赖 `CourseTable` 这个接口。将来想 mock 数据库做单元测试，换一个实现即可，logic 代码一行不用改。

### 3.6 统一响应与错误处理

所有接口的返回都是同一个结构（[model/common.go](model/common.go)）：

```go
type Response struct {
	Code int         `json:"code"` // 0 成功，1 失败
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}
```

而 Echo 要求 handler 返回 `error`，两边的「翻译」由 [middleware/wrap_response.go](middleware/wrap_response.go) 完成：

```go
func WrapResponse(handler func(echo.Context) (interface{}, error)) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		result, err := handler(ctx)
		if err != nil {
			return ctx.JSON(http.StatusOK, model.Response{Code: model.CodeFailure, Msg: err.Error(), Data: nil})
		}
		return ctx.JSON(http.StatusOK, result)
	}
}
```

这就是**自定义中间件**的标准写法：一个「接收 handler、返回新 handler」的函数，像洋葱一样把业务 handler 包在中间。本项目里它承担两个约定：

1. **HTTP 状态码恒为 200**，成败靠 `code` 字段区分——前端只需写一套解析逻辑；
2. controller 里区分两类错误：**参数绑定失败**（`return nil, err`，交给 WrapResponse 兜底）和**业务失败**（主动组装 `CodeFailure` 的 Response 返回）。

错误值本身也有讲究。[logic/course_logic/common.go](logic/course_logic/common.go)：

```go
func errRecordNotFound(ID int64) error {
	return fmt.Errorf("%w: course id=%d not found", dao.ErrNotFound, ID)
}
```

`%w` 是**错误包装**：生成的错误既带上下文信息，又能被 `errors.Is(err, dao.ErrNotFound)` 识别出根因。调用方想区分「记录不存在」和其他数据库错误时，不必做字符串匹配。

---

## 四、常用命令速查

### go 命令

| 命令 | 作用 |
|---|---|
| `go run .` | 编译并运行当前目录（开发主力命令） |
| `go build .` | 编译生成可执行文件 `echo-learning` |
| `go mod tidy` | 按 import 拉齐 go.mod 依赖（加了/删了 import 后跑） |
| `go mod download` | 只下载依赖不整理 |
| `go fmt ./...` | 格式化全部代码（Go 没有风格之争，统一 fmt） |
| `go vet ./...` | 静态检查，发现可疑代码 |
| `go test ./...` | 跑全部测试（本项目暂无测试文件） |
| `go doc github.com/labstack/echo/v4.Group` | 在终端看某个包/类型的文档 |

### curl 测试命令（API 速查）

统一前缀 `http://localhost:8080/api/courses`：

```bash
curl http://localhost:8080/api/courses                                # 列表
curl http://localhost:8080/api/courses/1                              # 详情
curl -X POST http://localhost:8080/api/courses \
  -H "Content-Type: application/json" \
  -d '{"title":"Go 入门","teacher":"张三","description":"基础","price":99.5}'   # 新增
curl -X PUT http://localhost:8080/api/courses/1 \
  -H "Content-Type: application/json" \
  -d '{"title":"Go 进阶","teacher":"李四","description":"深入","price":199}'    # 更新
curl -X DELETE http://localhost:8080/api/courses/1                    # 删除
```

### mysql 命令

```bash
mysql -u root -p                          # 登录
mysql -u root -p echo_learning            # 直接进库
mysql -u root -p echo_learning -e "SELECT * FROM courses;"   # 不进交互界面，执行一条 SQL 就退出
```

---

## 五、常见问题排查

启动失败时，**先看终端第一行报错**，按关键词对号入座：

| 报错关键词 | 原因 | 解决 |
|---|---|---|
| `dial tcp 127.0.0.1:3306: connect: connection refused` | MySQL 没启动 | `brew services start mysql@8.4` |
| `Access denied for user 'root'@'localhost'` | 密码不对 | 用 `DB_PASS=真实密码 go run .`，或 `ALTER USER` 改密码对齐 |
| `bind: address already in use` | 8080 端口被占 | `lsof -i :8080` 找到 PID 后 `kill`，或 `PORT=8081 go run .` |
| `unknown command` / 环境变量不生效 | 命令写法问题 | 环境变量要和 `go run .` 写在**同一条命令**里 |
| 前端页面调接口报错 | 后端没在 8080 跑 | 前端 Vite 代理写死了转发到 8080，先把后端跑起来 |
| 改了后端代码没生效 | `go run` 无热重载 | `Ctrl+C` 后重新 `go run .` |

---

## 六、进阶练习

跑通本项目后，可以按难度依次做这些改造（都能只用本文讲到的知识完成）：

1. **给列表接口加分页**：`GET /api/courses?page=1&size=10`。提示：请求结构体加 `query` tag，dao 里用 `Offset().Limit()` 链式拼 SQL。
2. **加一张新表**：仿照 Course 全链路（define → dao 实现 → logic → controller → router）加 `teachers` 表，体会「分层后加一张表要动几个文件」。
3. **补单元测试**：给 logic 层写测试，用接口替换 `CourseTable` 注入 mock 实现，验证「记录不存在」分支。
4. **热重载**：安装 [air](https://github.com/air-verse/air)（`go install github.com/air-verse/air@latest`），改代码自动重启，省掉手动 `Ctrl+C`。
5. **部署练习**：`GOOS=linux GOARCH=amd64 go build .` 交叉编译，把二进制和 `DB_*` 环境变量搬到一台 Linux 机器上跑起来。

---

> 配套阅读：项目概览见 [README.md](README.md)；前端（Vue3 + Vite）的教学说明见 [../web-Echo-learning/README.md](../web-Echo-learning/README.md)。
