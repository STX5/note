# Note📝
## 简介
本项目是一个简单的Golang微服务架构web应用，总共有三个微服务：api、note、user。其中 api 服务是对外暴露的 HTTP 接口，采用 Hertz 框架。user 服务与 note 服务为内部的微服务。微服务之间通过 kitex 框架的 RPC 功能进行交互。并且接入Opentelemetry和Jaeger进行观测与链路追踪。
``` bash
├── cmd
│   ├── api             # 对外暴露的 HTTP 接口服务
│   │   └── main.go     # api 服务的入口文件
        └── ...
│   ├── note            # 内部的笔记微服务
│   │   └── main.go     # note 服务的入口文件
        └── ...
│   ├── user            # 内部的用户微服务
│   │   └── main.go     # user 服务的入口文件
        └── ...

```
| ServiceName | Usage                     | Path     | IDL             |
| ----------- | -------------------------| -------- | ---------------|
| api         | 对外 HTTP 服务接口        | cmd/api  | idl/api.thrift |
| note        | 内部的笔记微服务         | cmd/note | idl/note.thrift|
| user        | 内部的用户微服务         | cmd/user | idl/user.thrift|

在该项目中，note 服务和 user 服务先通过 ETCD 进行服务注册，API 服务再通过 ETCD 解析出它所依赖的服务的地址。微服务之间通过 RPC 进行通信，API 服务则通过 Hertz 框架对外提供 HTTP 接口。这种架构使得服务间通信更加高效可靠，同时也提高了系统的可扩展性和可维护性。note与user服务都通过gorm操作Mysql数据库进行CRUD。
## 安装与运行
本项目采用了hz与kitex进行代码生成
``` shell
go install github.com/cloudwego/hertz/cmd/hz@latest
go install github.com/cloudwego/kitex/tool/cmd/kitex@latest
go install github.com/cloudwego/thriftgo@v0.13.0
```
通过Makefile进行代码生成，下表为命令对照
| Catalog              | Command                             |
| --------------------| ----------------------------------- |
| hertz_api_model     | make hertz_gen_model                |
| hertz_api_client    | make hertz_gen_client               |
| kitex_user_client   | make kitex_gen_user                 |
| kitex_note_client   | make kitex_gen_note                 |
| hertz_api_new       | cd cmd/api && make hertz_new_api     |
| hertz_api_update    | cd cmd/api && make hertz_update_api  |
| kitex_user_server   | cd cmd/user && make kitex_gen_server |
| kitex_note_server   | cd cmd/note && make kitex_gen_server |

启动依赖环境
```bash
docker-compose up
```
启动user服务
```bash
cd cmd/user
make run
```
启动note服务
```bash
cd cmd/note
make run
```
启动api服务
```bash
cd cmd/api
make run
```
参考 `api_request/api_service/api_service_test.go` 文件，构建客户端，通过HTTP接口进行访问。该测试文件包含IDL中定义的所有api接口服务。
``` thrift
service ApiService {
    CreateUserResponse CreateUser(1: CreateUserRequest req) (api.post="/v1/user/register")

    CheckUserResponse CheckUser(1: CheckUserRequest req) (api.post="/v1/user/login")

    CreateNoteResponse CreateNote(1: CreateNoteRequest req) (api.post="/v1/note")

    QueryNoteResponse QueryNote(1: QueryNoteRequest req) (api.get="/v1/note/query")

    UpdateNoteResponse UpdateNote(1: UpdateNoteRequest req) (api.put="/v1/note/:note_id")

    DeleteNoteResponse DeleteNote(1: DeleteNoteRequest req) (api.delete="/v1/note/:note_id")
}
```
## 开发指南：以User服务为例
```sh
cd cmd/user
make kitex_gen_server
```
即可看到生成的handler.go文件，内含user微服务的所有接口

在 `cmd/user/service/` 目录下对每个接口进行具体实现。
```sh
service/
├── check_user.go           # 用户登录、鉴权
├── create_user.go          # 创建用户
└── mget_user.go            # 获得多个用户信息 

```
查看 `cmd/user/dal/db/` 目录，可以看见user服务通过Gorm，从`consts`包中获得连接DSN，进行数据库连接。并且启用了`gormlogrus`作为logger，`opentelemetry`作为数据库访问的链路追踪。

由于gorm定义的model与kitex生成的rpc model并不完全相同，因此在 `cmd/user/pack` 包对user这个model进行了封装
```go
func User(u *db.User) *user.User {
	if u == nil {
		return nil
	}

	return &user.User{UserId: int64(u.ID), Username: u.Username, Avatar: "test"}
}
```
`main.go`
```golang
// 通过ETCD进行服务注册
r, err := etcd.NewEtcdRegistry([]string{consts.ETCDAddress})

// 设置OPTL云原生链路监测
p := provider.NewOpenTelemetryProvider(
	provider.WithServiceName(consts.UserServiceName),
        // Exporter地址
	provider.WithExportEndpoint(consts.ExportEndpoint),
	provider.WithInsecure(),
)
```
## 链路追踪、观测
### OPTL
Opentelemetry 是一个开源的分布式跟踪和指标收集框架，它可以帮助开发人员收集、分析和可视化分布式应用程序中的各种指标、日志和跟踪数据。Opentelemetry 可以与各种语言和框架集成，包括 Java、Python、Go、Node.js 等，并支持多种数据格式和协议。
### Jaeger
Jaeger 是一个开源的分布式跟踪系统，它可以帮助开发人员追踪应用程序中的请求流程，并记录每个请求经过的服务和调用链路，以及每个调用的时间和性能数据。Jaeger 可以与 Opentelemetry 集成，作为其后端存储和可视化的组件之一。
### Grafana
Grafana 是一个开源的数据可视化工具，它可以将收集到的指标、日志和跟踪数据以图表和面板的形式展示出来，以便开发人员更好地理解应用程序的状态和性能。Grafana 可以与多种数据源集成，包括 Prometheus、Jaeger、Opentelemetry 等，并提供丰富的可视化和报警功能。