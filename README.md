# Echo-learning 后端

Echo + GORM + MySQL 的教学项目：一张 `courses`（课程）表，实现完整的增删改查 REST API。

分层结构模仿 /Users/didi/pop 项目的组织方式：router → controller → logic → dao。

## 目录结构

```
Echo-learning/
├── main.go                          # 入口：一行调用 initialize.ServerInit()
├── initialize/init.go               # 服务初始化：数据库 → echo 实例 → 路由 → 启动
├── router/
│   ├── router.go                    # DispatchRouters：挂 /api 根 Group
│   └── course_router.go             # course 业务路由
├── controller/
│   ├── base.go                      # ParseRequest：绑定 + 校验请求
│   └── controller_course.go         # 参数绑定 → 调 logic → 组装 Response
├── logic/course_logic/              # 业务逻辑，每个功能一个文件
│   ├── list_course.go
│   ├── get_course.go
│   ├── create_course.go
│   ├── update_course.go
│   ├── delete_course.go
│   └── common.go
├── dao/
│   ├── define.go                    # 实体、Dao 接口、全局表实例 CourseTable
│   ├── session.go                   # GetSession：获取带请求上下文的 gorm 会话
│   └── impl_course.go               # courses 表的增删改查实现
├── model/
│   ├── common.go                    # 统一响应 Response{code,msg,data}
│   └── course.go                    # 各接口的请求结构体
├── services/mysql/mysql.go          # MySQL 连接（自动建库）
├── middleware/wrap_response.go      # WrapResponse：统一 JSON 输出
├── go.mod / go.sum
└── README.md
```

## 一次请求的链路

```
HTTP 请求
  → router（course_router.go 注册路由）
  → middleware.WrapResponse（统一响应包装）
  → controller（绑定参数，调 logic）
  → logic（业务校验与编排，每个功能一个文件）
  → dao（CourseTable 提供的 QueryXxx / Insert / Update / Delete）
  → GORM → MySQL
```

## 运行

### Docker Compose 一键启动（包含前端和 MySQL）

只需安装 Docker 和 Docker Compose，无需在宿主机安装 Go、Node.js 或 MySQL。
默认目录结构如下；如果前端在其他位置，修改 `.env` 中的 `FRONTEND_PATH`。

```text
Project/
├── Echo-learning/       # 本项目，包含 compose.yaml
└── web-Echo-learning/   # 前端，包含 Dockerfile 和 nginx.conf
```

在本项目目录执行：

```bash
cp .env.example .env
# 编辑 .env，修改 DB_PASS 和 MYSQL_ROOT_PASSWORD
docker compose up -d --build --wait
```

打开 http://localhost:8080，API 同样通过 http://localhost:8080/api/courses 访问。
端口被占用时，在 `.env` 中修改 `WEB_PORT`（例如 `18080`）。
若拉取 Docker Hub 镜像超时，可在 `.env` 中设置
`IMAGE_REGISTRY=docker.m.daocloud.io/library` 或其他可访问的镜像源前缀，再执行启动命令。
该设置同时用于 MySQL 镜像和前后端构建基础镜像，不修改全局 Docker 配置。
如果构建时无法访问 `proxy.golang.org`，可设置 `GOPROXY=https://goproxy.cn,direct`。
如果 `npm ci` 下载依赖超时，可设置 `NPM_CONFIG_REGISTRY=https://registry.npmmirror.com`。

三个服务通过 Compose 网络通信：浏览器 → Nginx（前端静态文件与 `/api` 代理）→ Go 后端 → MySQL。
前后端镜像采用多阶段构建，分别使用 `npm ci` 和 `go mod download` 安装锁定依赖后编译。
MySQL 和后端通过健康检查后才启动下一层服务，参考 [Docker 启动顺序说明](https://docs.docker.com/compose/how-tos/startup-order/)。
仅 Nginx 的端口映射到宿主机；后端与 MySQL 使用容器内端口 8080、3306。
MySQL 镜像首次启动时创建数据库和专用用户，后端自动创建 `courses` 表。
容器内的 `DB_HOST` 固定为服务名 `db`；前端 `/api` 请求由 Nginx 转发到 `backend:8080`。

```bash
docker compose ps                     # 查看服务与健康状态
docker compose logs -f backend        # 查看后端日志
docker compose logs -f db             # 查看数据库日志
docker compose up -d --build --wait    # 修改代码后重新构建并启动
docker compose down                   # 停止并移除容器，保留数据库数据
```

数据保存在命名卷 `echo-learning_mysql_data`，重建容器或执行 `down` 后仍然保留。
`MYSQL_*` 初始化设置仅对空数据卷生效；已有数据时修改 `.env` 不会自动修改数据库用户、密码或数据库名。
如确需清空数据重新初始化，执行 `docker compose down -v`（会删除该项目的全部数据库数据）。
`.env` 不进入镜像，也已加入 `.gitignore`。
搬到其他机器时，将两个源码目录一起复制，保持上述结构；在后端目录准备 `.env` 后运行相同命令即可。

### 本地开发

```bash
# 依赖：本地 MySQL 已启动（默认 root/1234，可用环境变量覆盖）
go mod tidy
go run .
```

连接配置通过环境变量覆盖，默认值如下：

| 变量 | 默认值 |
|---|---|
| `DB_USER` | root |
| `DB_PASS` | 1234 |
| `DB_HOST` | 127.0.0.1 |
| `DB_PORT` | 3306 |
| `DB_NAME` | echo_learning |
| `PORT` | 8080 |

首次启动会自动创建 `echo_learning` 数据库和 `courses` 表，无需手动建表。

## API 一览

统一响应格式：`{"code": 0, "msg": "success", "data": ...}`，`code=0` 成功，非 0 失败。

| 方法 | 路径 | 说明 |
|---|---|---|
| GET | `/api/courses` | 课程列表 |
| GET | `/api/courses/:id` | 课程详情 |
| POST | `/api/courses` | 新增课程 |
| PUT | `/api/courses/:id` | 更新课程 |
| DELETE | `/api/courses/:id` | 删除课程 |

### 示例

```bash
# 新增
curl -X POST http://localhost:8080/api/courses \
  -H "Content-Type: application/json" \
  -d '{"title":"Go 入门","teacher":"张三","description":"Echo + GORM 基础","price":99.5}'

# 列表
curl http://localhost:8080/api/courses

# 更新
curl -X PUT http://localhost:8080/api/courses/1 \
  -H "Content-Type: application/json" \
  -d '{"title":"Go 进阶","teacher":"李四","description":"深入 Go","price":199}'

# 删除
curl -X DELETE http://localhost:8080/api/courses/1
```

## 学习要点

1. **分层**（模仿 pop）：`router` 只管路由分发；`controller` 只管参数绑定和响应组装；`logic` 承载业务规则；`dao` 隔离 SQL 细节，接口 + 实现分离（`CourseTable` 全局实例）。
2. **统一响应**：controller 统一返回 `(interface{}, error)`，`WrapResponse` 中间件兜底转换成 `{code, msg, data}`。
3. **GORM**：`AutoMigrate` 自动建表；`Find / Create / Save / Delete` 对应增删改查；`WithContext` 让 SQL 随请求超时取消。

> 想系统学习本项目用到的技术栈（Echo / GORM / Go Modules）和启动命令，看 [TEACHING.md](TEACHING.md)。
